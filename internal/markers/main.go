
package markers

type Marker struct {
	// Define marker fields here
}


func New () *Marker {
	return &Marker{}
}

fund c (m *Marker) Create() error{


	// We require a marker name to be set, as this is used to identify the marker.
	if marker.Name == "" {
		metaData.Error = "Marker name is required"
		metaData.MissingFields = []string{"name"}
		errorResponse := api.CreateError(
			api.HttpStatusBadRequest,
			api.ApplicationStatusAddFailed,
			api.MarkerMissingInfo,
			metaData,
		)
		api.LogError(logger, errorResponse)
		c.JSON(api.HttpStatusBadRequest, api.AddMarkerErrorResponse{ErrorResponse: errorResponse})
		return
	}

	// Set the duration, difference between start and end time
	marker.Duration = marker.EndTimestamp - marker.StartTimestamp

	// Add the marker to the database
	if insertedMarker, err := database.AddMarker(ctxAddMarker, user, marker); err != nil {
		metaData.Error = err.Error()
		errorResponse := api.CreateError(
			api.HttpStatusBadRequest,
			api.ApplicationStatusAddFailed,
			api.MarkerAddFailed,
			metaData,
		)
		api.LogError(logger, errorResponse)
		c.JSON(api.HttpStatusBadRequest, api.AddMarkerErrorResponse{ErrorResponse: errorResponse})
		return
	} else {
		marker = insertedMarker
	}
}