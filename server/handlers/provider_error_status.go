package handlers

import (
	"net/http"

	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

func providerErrStatusCode(err error, fallback int) int {
	if meshkiterrors.GetCode(err) == models.ErrProviderInvalidUUIDCode {
		return http.StatusBadRequest
	}

	return fallback
}
