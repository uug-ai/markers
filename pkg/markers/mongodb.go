package markers

import (
	"context"
	"time"

	ingestmarkers "github.com/uug-ai/ingest/pkg/markers"
	"github.com/uug-ai/models/pkg/models"
	"github.com/uug-ai/trace/pkg/opentelemetry"
	"go.mongodb.org/mongo-driver/mongo"
)

// This package is a thin compatibility shim over github.com/uug-ai/ingest's
// markers package, which is the single implementation of the Hub's marker write
// side effects: the denormalised option/range/category projections, the media
// tagging, and the per-occurrence markerSummary. Keeping the shim lets the
// direct callers (analysers, alert pipelines, cli seeders) keep importing
// github.com/uug-ai/markers unchanged while the logic lives in one place, so a
// change like the markerSummary denormalisation can never again drift between
// two copies.

// Collection-name and connection overrides retained for backward compatibility:
// the cli seeders mutate these package-level vars to point the writer at
// non-default collections. They are forwarded to the ingest markers
// implementation before each write (see syncConfig), so setting e.g.
// markers.DatabaseName still configures the underlying writer.
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

// syncConfig forwards this shim's collection-name and connection overrides to the
// ingest markers implementation so a caller that set e.g. markers.DatabaseName or
// markers.MARKERS_COLLECTION still configures the actual writer. It runs before
// every delegated write, mirroring the original package's mutable-global
// configuration model (the cli seeders set these once at startup).
func syncConfig() {
	ingestmarkers.MARKERS_COLLECTION = MARKERS_COLLECTION
	ingestmarkers.MARKER_OPTIONS_COLLECTION = MARKER_OPTIONS_COLLECTION
	ingestmarkers.MARKER_OPTION_RANGES_COLLECTION = MARKER_OPTION_RANGES_COLLECTION
	ingestmarkers.MARKER_TAG_OPTIONS_COLLECTION = MARKER_TAG_OPTIONS_COLLECTION
	ingestmarkers.MARKER_TAG_OPTION_RANGES_COLLECTION = MARKER_TAG_OPTION_RANGES_COLLECTION
	ingestmarkers.MARKER_EVENT_OPTIONS_COLLECTION = MARKER_EVENT_OPTIONS_COLLECTION
	ingestmarkers.MARKER_EVENT_OPTION_RANGES_COLLECTION = MARKER_EVENT_OPTION_RANGES_COLLECTION
	ingestmarkers.MARKER_CATEGORY_OPTIONS_COLLECTION = MARKER_CATEGORY_OPTIONS_COLLECTION
	ingestmarkers.MEDIA_COLLECTION = MEDIA_COLLECTION
	ingestmarkers.DatabaseName = DatabaseName
	ingestmarkers.TIMEOUT = TIMEOUT
}

// AddMarkerToMongodb inserts a fresh marker (a new _id every call) together with
// its denormalised option/range/category projections and media tagging. It
// delegates to github.com/uug-ai/ingest/pkg/markers, forwarding the package-level
// collection overrides first.
func AddMarkerToMongodb(ctxTracer context.Context, tracer *opentelemetry.Tracer, client *mongo.Client, marker models.Marker, mediaIds ...string) (models.Marker, error) {
	syncConfig()
	return ingestmarkers.AddMarkerToMongodb(ctxTracer, tracer, client, marker, mediaIds...)
}
