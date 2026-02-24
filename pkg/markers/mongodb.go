package markers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/uug-ai/models/pkg/models"
	"github.com/uug-ai/trace/pkg/opentelemetry"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MARKERS_COLLECTION                    = "markers"
	MARKER_OPTIONS_COLLECTION             = "marker_options"
	MARKER_OPTION_RANGES_COLLECTION       = "marker_option_ranges"
	MARKER_TAG_OPTIONS_COLLECTION         = "marker_tag_options"
	MARKER_TAG_OPTION_RANGES_COLLECTION   = "marker_tag_option_ranges"
	MARKER_EVENT_OPTIONS_COLLECTION       = "marker_event_options"
	MARKER_EVENT_OPTION_RANGES_COLLECTION = "marker_event_option_ranges"
	MARKER_CATEGORY_OPTIONS_COLLECTION    = "marker_category_options"
	MEDIA_COLLECTION                      = "media"

	DatabaseName = "Kerberos"
	TIMEOUT      = 10 * time.Second
)

func AddMarkerToMongodb(ctxTracer context.Context, tracer *opentelemetry.Tracer, client *mongo.Client, marker models.Marker, mediaIds ...string) (models.Marker, error) {

	ctxAddMarkerToMongodb, span := tracer.CreateSpan(ctxTracer, map[string]string{})
	defer span.End()

	ctx, cancel := context.WithTimeout(ctxAddMarkerToMongodb, TIMEOUT)
	defer cancel()

	// Open markers collection
	db := client.Database(DatabaseName)
	c := db.Collection(MARKERS_COLLECTION)

	// Generate new ID for the marker
	marker.Id = primitive.NewObjectID()

	res, err := c.InsertOne(ctx, marker)
	if err != nil {
		return models.Marker{}, err
	}

	// Check if the inserted ID is of type primitive.ObjectID
	if res.InsertedID == nil {
		return models.Marker{}, errors.New("Inserted ID is nil")
	}

	if _, ok := res.InsertedID.(primitive.ObjectID); !ok {
		return models.Marker{}, errors.New("Inserted ID is not of type primitive.ObjectID")
	}

	marker.Id = res.InsertedID.(primitive.ObjectID)

	// As part of the marker we also need to insert into some other collections for performance reasons.
	// For example on the media page we have marker options, marker event options, marker tag options, marker category options.

	// Collections for tracking unique entries
	nameSet := make(map[string]struct{})
	tagSet := make(map[string]struct{})
	eventSet := make(map[string]struct{})
	categorySet := make(map[string]struct{})

	// Slices for bulk operations
	var markerOptUpserts []mongo.WriteModel
	var tagOptUpserts []mongo.WriteModel
	var eventOptUpserts []mongo.WriteModel
	var categoryOptUpserts []mongo.WriteModel

	// Slices for range documents
	var markerRangeDocs []any
	var tagRangeDocs []any
	var eventRangeDocs []any

	now := time.Now().Unix()

	// marker option upsert
	if marker.Name != "" {
		if _, exists := nameSet[marker.Name]; !exists {
			nameSet[marker.Name] = struct{}{}
			var categoryNamesList []string
			for _, cat := range marker.Categories {
				if cat.Name != "" {
					categoryNamesList = append(categoryNamesList, cat.Name)
				}
			}
			up := mongo.NewUpdateOneModel()
			up.SetFilter(bson.M{"value": marker.Name, "organisationId": marker.OrganisationId})
			up.SetUpdate(bson.M{
				"$setOnInsert": bson.M{
					"value":          marker.Name,
					"text":           marker.Name,
					"organisationId": marker.OrganisationId,
					"createdAt":      now,
				},
				"$set": bson.M{
					"updatedAt": now,
				},
				"$addToSet": bson.M{
					"categories": bson.M{"$each": categoryNamesList},
				},
			})
			up.SetUpsert(true)
			markerOptUpserts = append(markerOptUpserts, up)
		}
		markerRangeDocs = append(markerRangeDocs, bson.M{
			"value":          marker.Name,
			"text":           marker.Name,
			"organisationId": marker.OrganisationId,
			"start":          marker.StartTimestamp,
			"end":            marker.EndTimestamp,
			"deviceId":       marker.DeviceId,
			"groupId":        marker.GroupId,
			"createdAt":      now,
		})
	}

	// tags
	for _, tag := range marker.Tags {
		if tag.Name == "" {
			continue
		}
		if _, exists := tagSet[tag.Name]; !exists {
			tagSet[tag.Name] = struct{}{}
			up := mongo.NewUpdateOneModel()
			up.SetFilter(bson.M{"value": tag.Name, "organisationId": marker.OrganisationId})
			up.SetUpdate(bson.M{
				"$setOnInsert": bson.M{
					"value":          tag.Name,
					"text":           tag.Name,
					"organisationId": marker.OrganisationId,
					"createdAt":      now,
				},
				"$set": bson.M{
					"updatedAt": now,
				},
			})
			up.SetUpsert(true)
			tagOptUpserts = append(tagOptUpserts, up)
		}
		tagRangeDocs = append(tagRangeDocs, bson.M{
			"value":          tag.Name,
			"text":           tag.Name,
			"organisationId": marker.OrganisationId,
			"start":          marker.StartTimestamp,
			"end":            marker.EndTimestamp,
			"deviceId":       marker.DeviceId,
			"groupId":        marker.GroupId,
			"createdAt":      now,
		})
	}

	// events
	for _, event := range marker.Events {
		if event.Name == "" {
			continue
		}
		if _, exists := eventSet[event.Name]; !exists {
			eventSet[event.Name] = struct{}{}
			up := mongo.NewUpdateOneModel()
			up.SetFilter(bson.M{"value": event.Name, "organisationId": marker.OrganisationId})
			up.SetUpdate(bson.M{
				"$setOnInsert": bson.M{
					"value":          event.Name,
					"text":           event.Name,
					"organisationId": marker.OrganisationId,
					"createdAt":      now,
				},
				"$set": bson.M{
					"updatedAt": now,
				},
			})
			up.SetUpsert(true)
			eventOptUpserts = append(eventOptUpserts, up)
		}
		eventRangeDocs = append(eventRangeDocs, bson.M{
			"value":          event.Name,
			"text":           event.Name,
			"organisationId": marker.OrganisationId,
			"start":          event.StartTimestamp,
			"end":            event.EndTimestamp,
			"deviceId":       marker.DeviceId,
			"groupId":        marker.GroupId,
			"createdAt":      now,
			"updatedAt":      now,
		})
	}

	// categories
	for _, category := range marker.Categories {
		if category.Name == "" {
			continue
		}
		if _, exists := categorySet[category.Name]; !exists {
			categorySet[category.Name] = struct{}{}
			up := mongo.NewUpdateOneModel()
			up.SetFilter(bson.M{"value": category.Name, "organisationId": marker.OrganisationId})
			up.SetUpdate(bson.M{
				"$setOnInsert": bson.M{
					"value":          category.Name,
					"text":           category.Name,
					"organisationId": marker.OrganisationId,
					"createdAt":      now,
				},
				"$set": bson.M{
					"updatedAt": now,
				},
			})
			up.SetUpsert(true)
			categoryOptUpserts = append(categoryOptUpserts, up)
		}
	}

	// Execute bulk operations for marker options
	if len(markerOptUpserts) > 0 {
		markerOptCol := db.Collection(MARKER_OPTIONS_COLLECTION)
		if _, err := markerOptCol.BulkWrite(ctx, markerOptUpserts); err != nil {
			return marker, fmt.Errorf("failed to upsert marker options: %w", err)
		}
	}

	// Insert marker option ranges
	if len(markerRangeDocs) > 0 {
		markerRangeCol := db.Collection(MARKER_OPTION_RANGES_COLLECTION)
		if _, err := markerRangeCol.InsertMany(ctx, markerRangeDocs); err != nil {
			return marker, fmt.Errorf("failed to insert marker ranges: %w", err)
		}
	}

	// Execute bulk operations for tag options
	if len(tagOptUpserts) > 0 {
		tagOptCol := db.Collection(MARKER_TAG_OPTIONS_COLLECTION)
		if _, err := tagOptCol.BulkWrite(ctx, tagOptUpserts); err != nil {
			return marker, fmt.Errorf("failed to upsert tag options: %w", err)
		}
	}

	// Insert tag option ranges
	if len(tagRangeDocs) > 0 {
		tagRangeCol := db.Collection(MARKER_TAG_OPTION_RANGES_COLLECTION)
		if _, err := tagRangeCol.InsertMany(ctx, tagRangeDocs); err != nil {
			return marker, fmt.Errorf("failed to insert tag ranges: %w", err)
		}
	}

	// Execute bulk operations for event options
	if len(eventOptUpserts) > 0 {
		eventOptCol := db.Collection(MARKER_EVENT_OPTIONS_COLLECTION)
		if _, err := eventOptCol.BulkWrite(ctx, eventOptUpserts); err != nil {
			return marker, fmt.Errorf("failed to upsert event options: %w", err)
		}
	}

	// Insert event option ranges
	if len(eventRangeDocs) > 0 {
		eventRangeCol := db.Collection(MARKER_EVENT_OPTION_RANGES_COLLECTION)
		if _, err := eventRangeCol.InsertMany(ctx, eventRangeDocs); err != nil {
			return marker, fmt.Errorf("failed to insert event ranges: %w", err)
		}
	}

	// Execute bulk operations for category options
	if len(categoryOptUpserts) > 0 {
		categoryOptCol := db.Collection(MARKER_CATEGORY_OPTIONS_COLLECTION)
		if _, err := categoryOptCol.BulkWrite(ctx, categoryOptUpserts); err != nil {
			return marker, fmt.Errorf("failed to upsert category options: %w", err)
		}
	}

	// If mediaIds are provided, update the media documents with marker names, tag names, and event names
	for _, mediaId := range mediaIds {
		if mediaId == "" {
			continue
		}

		mediaObjectId, err := primitive.ObjectIDFromHex(mediaId)
		if err != nil {
			return marker, fmt.Errorf("invalid mediaId format: %w", err)
		}

		// Collect unique marker names, tag names, and event names
		var markerNames []string
		if marker.Name != "" {
			markerNames = append(markerNames, marker.Name)
		}

		var tagNames []string
		for _, tag := range marker.Tags {
			if tag.Name != "" {
				tagNames = append(tagNames, tag.Name)
			}
		}

		var eventNames []string
		for _, event := range marker.Events {
			if event.Name != "" {
				eventNames = append(eventNames, event.Name)
			}
		}

		// Build update document using $addToSet with $each to ensure uniqueness
		updateDoc := bson.M{}
		if len(markerNames) > 0 {
			updateDoc["markerNames"] = bson.M{"$each": markerNames}
		}
		if len(tagNames) > 0 {
			updateDoc["tagNames"] = bson.M{"$each": tagNames}
		}
		if len(eventNames) > 0 {
			updateDoc["eventNames"] = bson.M{"$each": eventNames}
		}

		if len(updateDoc) > 0 {
			mediaCol := db.Collection(MEDIA_COLLECTION)
			filter := bson.M{
				"_id":            mediaObjectId,
				"startTimestamp": bson.M{"$lte": marker.StartTimestamp},
				"endTimestamp":   bson.M{"$gte": marker.StartTimestamp},
			}
			update := bson.M{"$addToSet": updateDoc}
			_, err := mediaCol.UpdateOne(ctx, filter, update)
			if err != nil {
				return marker, fmt.Errorf("failed to update media with marker data: %w", err)
			}
		}
	}

	return marker, nil
}
