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

func New() *Marker {
	return &Marker{}
}

func (m *Marker) Create(ctxTracer context.Context, tracer *opentelemetry.Tracer, client *mongo.Client, marker models.Marker) (models.Marker, error) {

	// We require a marker name to be set, as this is used to identify the marker.
	if marker.Name == "" {
		return models.Marker{}, errors.New("marker name is required")
	}

	// Set the duration, difference between start and end time
	marker.Duration = marker.EndTimestamp - marker.StartTimestamp

	// Add the marker to the database
	insertedMarker, err := AddMarkerToMongodb(ctxTracer, tracer, client, marker)
	if err != nil {
		return models.Marker{}, err
	}

	return insertedMarker, nil
}
