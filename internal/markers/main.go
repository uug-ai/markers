
package markers

type Marker struct {
	// Define marker fields here
}


func New () *Marker {
	return &Marker{}
}

fund c (m *Marker) Create(){
	// Bind JSON to marker struct
	if err := c.ShouldBindJSON(&marker); err != nil {
		metaData.Error = err.Error()
		errorResponse := api.CreateError(
			api.HttpStatusBadRequest,
			api.ApplicationStatusAddFailed,
			api.MarkerBindingFailed,
			metaData,
		)
		api.LogError(logger, errorResponse)
		c.JSON(api.HttpStatusBadRequest, api.AddMarkerErrorResponse{ErrorResponse: errorResponse})
		return
	}

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

	// Add marker to database
	if user.Role == "application" {
		// For application users, we will take the organisation ID from the marker payload.
		// @TODO check if the application has access to the organisation ID.
		if marker.OrganisationId == "" {
			metaData.Error = "Organisation ID is required for application users"
			metaData.MissingFields = []string{"organisationId"}
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
	} else {
		// Verify that the organisation ID and user organisation match
		if marker.OrganisationId != "" && marker.OrganisationId != user.Id.Hex() {
			metaData.Error = "Organisation ID does not match user organisation"
			errorResponse := api.CreateError(
				api.HttpStatusBadRequest,
				api.ApplicationStatusAddFailed,
				api.MarkerAddFailed,
				metaData,
			)
			api.LogError(logger, errorResponse)
			c.JSON(api.HttpStatusBadRequest, api.AddMarkerErrorResponse{ErrorResponse: errorResponse})
			return
		}
		marker.OrganisationId = user.Id.Hex() // Set the organisation ID for the marker
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

	// Success - return created marker
	successResponse := api.CreateSuccess(
		api.HttpStatusOK,
		api.ApplicationStatusAddSuccess,
		api.MarkerAddSuccess,
		metaData,
	)
	data := api.AddMarkerResponse{
		Marker: marker,
	}
	api.LogDebug(
		logger,
		api.CreateDebug(
			api.HttpStatusOK,
			api.ApplicationStatusAddSuccess,
			api.MarkerAddSuccess,
			metaData,
		),
	)
	c.JSON(api.HttpStatusOK, api.AddMarkerSuccessResponse{
		SuccessResponse: successResponse,
		Data:            data,
	})
}
