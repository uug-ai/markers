// Package markers persists video annotations to MongoDB.
//
// Deprecated: use github.com/uug-ai/ingest/pkg/markers instead. This module is
// no longer maintained and nothing on the platform imports it; cli was the last
// consumer and now writes through ingest.
//
// It is not merely superseded, it is unsafe against a multi-tenant database.
// The media-tagging writes in AddMarkerToMongodb bound only on a device label
// and a time window (mongodb.go:332 and :385) with no organisation clause at
// all, so a device key reused across organisations tags media belonging to
// another tenant — and the second of the two is an UpdateMany. The option and
// range collections are organisation-scoped, which is what made the gap easy to
// miss. There is no project axis anywhere: the models pin (v1.4.26) predates
// Marker.ProjectId, so every marker this writes is unreachable from a
// project-scoped read.
//
// github.com/uug-ai/ingest/pkg/markers is the same API — New, Create,
// AddMarkerToMongodb — with both tenant clauses applied to every filter.
package markers

import (
	"context"
	"errors"

	"github.com/uug-ai/models/pkg/models"
	"github.com/uug-ai/trace/pkg/opentelemetry"
	"go.mongodb.org/mongo-driver/mongo"
)

type Marker struct {
	// Define marker fields here
}

// New returns a Marker writer.
//
// Deprecated: use github.com/uug-ai/ingest/pkg/markers.New. See the package
// comment for why.
func New() *Marker {
	return &Marker{}
}

// Create validates the marker and persists it.
//
// Deprecated: use github.com/uug-ai/ingest/pkg/markers.(*Marker).Create. See the
// package comment for why.
func (m *Marker) Create(ctxTracer context.Context, tracer *opentelemetry.Tracer, client *mongo.Client, marker models.Marker, mediaIds ...string) (models.Marker, error) {

	// We require a marker name to be set, as this is used to identify the marker.
	if marker.Name == "" {
		return models.Marker{}, errors.New("marker name is required")
	}

	// Set the duration, difference between start and end time
	marker.Duration = marker.EndTimestamp - marker.StartTimestamp

	// Add the marker to the database
	insertedMarker, err := AddMarkerToMongodb(ctxTracer, tracer, client, marker, mediaIds...)
	if err != nil {
		return models.Marker{}, err
	}

	return insertedMarker, nil
}
