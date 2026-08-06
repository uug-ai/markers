package markers

import (
	"context"

	ingestmarkers "github.com/uug-ai/ingest/pkg/markers"
	"github.com/uug-ai/models/pkg/models"
	"github.com/uug-ai/trace/pkg/opentelemetry"
	"go.mongodb.org/mongo-driver/mongo"
)

// Marker is the authoring handle retained for the direct callers (analysers,
// alert pipelines, cli seeders) that create markers outside the ingest core. It
// delegates to github.com/uug-ai/ingest/pkg/markers.
type Marker struct {
	// Define marker fields here
}

func New() *Marker {
	return &Marker{}
}

// Create validates and inserts a single marker, computing its duration. It
// forwards the package-level collection overrides and delegates to the shared
// ingest markers implementation (insert path: each call creates a distinct
// marker document).
func (m *Marker) Create(ctxTracer context.Context, tracer *opentelemetry.Tracer, client *mongo.Client, marker models.Marker, mediaIds ...string) (models.Marker, error) {
	syncConfig()
	return ingestmarkers.New().Create(ctxTracer, tracer, client, marker, mediaIds...)
}
