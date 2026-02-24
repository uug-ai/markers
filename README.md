# Markers

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/github/license/uug-ai/markers.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/uug-ai/markers?status.svg)](https://godoc.org/github.com/uug-ai/markers)
[![Go Report Card](https://goreportcard.com/badge/github.com/uug-ai/markers)](https://goreportcard.com/report/github.com/uug-ai/markers)
[![Release](https://img.shields.io/github/release/uug-ai/markers.svg)](https://github.com/uug-ai/markers/releases/latest)

A Go library for managing video markers in MongoDB with built-in support for tags, events, categories, and OpenTelemetry tracing.

A markers library for the Kerberos video surveillance platform that provides a unified interface for creating and persisting video annotations with optimized bulk operations and multi-tenant support.

## Features

- **MongoDB Integration**: Full MongoDB support with optimized bulk write operations
- **Multi-tenant Support**: Organization-scoped markers with data isolation
- **OpenTelemetry Tracing**: Built-in distributed tracing for observability
- **Rich Metadata**: Support for tags, events, and categories on markers
- **Option Collections**: Automatic management of unique marker/tag/event/category options per organization
- **Range Tracking**: Time range documents for efficient filtering and querying
- **Media Linking**: Automatic linking of markers to media documents
- **Production Ready**: Optimized for high-performance video annotation applications

## Installation

```bash
go get github.com/uug-ai/markers
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/uug-ai/markers/pkg/markers"
    "github.com/uug-ai/models/pkg/models"
    "github.com/uug-ai/trace/pkg/opentelemetry"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    // Connect to MongoDB
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }

    // Initialize tracer
    tracer := opentelemetry.NewTracer("markers-service")
    ctx := context.Background()

    // Create marker instance
    m := markers.New()

    // Define a marker
    marker := models.Marker{
        Name:           "Person Detected",
        StartTimestamp: 1708790400,
        EndTimestamp:   1708790410,
        OrganisationId: "org-123",
        DeviceId:       "camera-001",
        GroupId:        "group-001",
        Tags:           []models.Tag{{Name: "security"}},
        Events:         []models.Event{{Name: "motion", StartTimestamp: 1708790400}},
        Categories:     []models.Category{{Name: "surveillance"}},
    }

    // Create marker in MongoDB (mediaIds are optional)
    result, err := m.Create(ctx, tracer, client, marker, "media-id-here")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Created marker with ID: %s", result.Id.Hex())
}
```

## Core Concepts

### Marker Structure

Markers are annotations applied to video segments containing:

- **Name**: Required identifier for the marker (e.g., "Person Detected", "Vehicle Entry")
- **Timestamps**: Start and end timestamps defining the marker's time range
- **Duration**: Automatically calculated from timestamps
- **Tags**: Additional labels for categorization
- **Events**: Time-stamped events within the marker
- **Categories**: Hierarchical classification
- **Organization/Device/Group**: Multi-tenant identifiers

### Creating a Marker

Each marker creation follows this pattern:

1. **Initialize** the Marker instance using `markers.New()`
2. **Build** the marker model with required fields (Name is mandatory)
3. **Call** `Create()` with context, tracer, MongoDB client, marker, and optional mediaIds (variadic `...string`)
4. **Handle** the returned marker with its generated ID

### MongoDB Collections

The library manages multiple collections for optimal query performance:

| Collection | Purpose |
|------------|---------|
| `markers` | Primary marker storage |
| `marker_options` | Unique marker names per organization |
| `marker_option_ranges` | Time ranges for each marker name |
| `marker_tag_options` | Unique tag names per organization |
| `marker_tag_option_ranges` | Time ranges for each tag |
| `marker_event_options` | Unique event names per organization |
| `marker_event_option_ranges` | Time ranges for each event |
| `marker_category_options` | Unique category names per organization |
| `media` | Updated with marker/tag/event names |

## Usage Examples

### Basic Marker Creation

```go
package main

import (
    "context"
    "log"

    "github.com/uug-ai/markers/pkg/markers"
    "github.com/uug-ai/models/pkg/models"
)

func main() {
    m := markers.New()

    marker := models.Marker{
        Name:           "Motion Detected",
        StartTimestamp: 1708790400,
        EndTimestamp:   1708790420,
        OrganisationId: "org-123",
        DeviceId:       "camera-001",
    }

    // No mediaIds - skip media linking
    result, err := m.Create(ctx, tracer, mongoClient, marker)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Marker created: %s", result.Id.Hex())
}
```

### Marker with Tags and Events

```go
marker := models.Marker{
    Name:           "Security Alert",
    StartTimestamp: 1708790400,
    EndTimestamp:   1708790500,
    OrganisationId: "org-123",
    DeviceId:       "camera-001",
    GroupId:        "entrance-group",
    Tags: []models.Tag{
        {Name: "high-priority"},
        {Name: "entrance"},
    },
    Events: []models.Event{
        {Name: "door-opened", StartTimestamp: 1708790420, EndTimestamp: 1708790425},
        {Name: "person-entered", StartTimestamp: 1708790430, EndTimestamp: 1708790435},
    },
    Categories: []models.Category{
        {Name: "security"},
        {Name: "access-control"},
    },
}

result, err := m.Create(ctx, tracer, client, marker, "media-object-id")
```

### Linking to Media

When mediaIds are provided, the marker's names, tags, and events are automatically added to the corresponding media documents:

```go
// Link to a single media document
result, err := m.Create(ctx, tracer, client, marker, "64a1b2c3d4e5f6789012abcd")

// Link to multiple media documents
result, err := m.Create(ctx, tracer, client, marker, "media-1", "media-2", "media-3")

// Skip media linking (no mediaIds)
result, err := m.Create(ctx, tracer, client, marker)
```

The media document is updated using `$addToSet` to ensure uniqueness, and only if the marker's timestamp falls within the media's time range.

## Project Structure

```
.
├── pkg/
│   └── markers/               # Core markers implementation
│       ├── main.go           # Marker struct and Create() method
│       └── mongodb.go        # MongoDB persistence layer
├── main.go                    # Entry point
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## Configuration

### Database Settings

The library uses the following default configuration:

```go
DatabaseName = "Kerberos"
TIMEOUT      = 10 * time.Second
```

### Environment Variables

You can configure the MongoDB connection using environment variables:

```bash
# MongoDB
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=Kerberos

# Tracing
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
OTEL_SERVICE_NAME=markers-service
```

## Validation

The library validates marker data before persistence:

- **Name**: Required field - markers must have a name
- **Timestamps**: Used to calculate duration automatically

Validation is performed in the `Create()` method:

```go
if marker.Name == "" {
    return models.Marker{}, errors.New("marker name is required")
}
```

## Error Handling

The library provides clear error handling at each stage:

```go
// Create marker (omit mediaIds to skip media linking)
result, err := m.Create(ctx, tracer, client, marker)
if err != nil {
    // Validation error or database error
    log.Printf("Failed to create marker: %v", err)
    return
}

// Marker created successfully
log.Printf("Created marker: %s", result.Id.Hex())
```

Common error scenarios:
- Missing required `Name` field
- MongoDB connection issues
- Invalid `mediaId` format (when provided)
- Bulk write failures for option collections

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests for specific components:

```bash
# Markers tests
go test ./pkg/markers -v
```

## Contributing

Contributions are welcome! When adding new features, please follow the existing patterns demonstrated in this repository.

### Development Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Follow the existing code patterns
4. Add comprehensive tests for your changes
5. Ensure all tests pass: `go test ./...`
6. Commit your changes following [Conventional Commits](https://www.conventionalcommits.org/)
7. Push to your branch (`git push origin feat/amazing-feature`)
8. Open a Pull Request

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `types`

**Scopes**:
- `markers` - Core markers functionality
- `mongodb` - MongoDB persistence layer
- `options` - Marker options collections
- `docs` - Documentation updates
- `tests` - Test updates

**Examples**:

```
feat(markers): add bulk marker creation support
fix(mongodb): correct timeout handling in bulk writes
docs(readme): update usage examples
refactor(options): optimize deduplication logic
test(markers): add validation error tests
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Dependencies

This project uses the following key libraries:

- [mongo-driver](https://github.com/mongodb/mongo-go-driver) - Official MongoDB Go driver
- [uug-ai/models](https://github.com/uug-ai/models) - Shared model types including Marker struct
- [uug-ai/trace](https://github.com/uug-ai/trace) - OpenTelemetry tracing utilities

See [go.mod](go.mod) for the complete list of dependencies.

## Support

- **Issues**: [GitHub Issues](https://github.com/uug-ai/markers/issues)
- **Discussions**: [GitHub Discussions](https://github.com/uug-ai/markers/discussions)
- **Documentation**: See inline code comments and examples above

---

**Built with ❤️ by [UUG.AI](https://github.com/uug-ai)**
