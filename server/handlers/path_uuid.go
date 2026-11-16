package handlers

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/meshery/meshery/server/models"
)

// requirePathUUID reads the named path variable and checks that it is a
// well-formed UUID. A malformed value is a client error, so it is answered
// here with 400 and models.ErrInvalidUUID instead of being passed to the
// provider, where it used to surface as 502 Bad Gateway. ok is false when the
// response has already been written and the caller must return.
func (h *Handler) requirePathUUID(w http.ResponseWriter, req *http.Request, name string) (id string, ok bool) {
	raw := mux.Vars(req)[name]
	if _, err := uuid.FromString(raw); err != nil {
		invalidIDErr := models.ErrInvalidUUID(err)
		h.log.Error(invalidIDErr)
		writeMeshkitError(w, invalidIDErr, http.StatusBadRequest)
		return "", false
	}
	return raw, true
}
