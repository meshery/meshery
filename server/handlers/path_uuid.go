package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/meshery/meshery/server/models"
)

// requirePathUUID reads the named path variable and checks that it is a
// well-formed, non-nil UUID. A malformed or zero-value UUID is a client error,
// so it is answered here with 400 and models.ErrInvalidUUID instead of being passed
// to the provider or database layer. ok is false when the response has already been
// written and the caller must return.
func (h *Handler) requirePathUUID(w http.ResponseWriter, req *http.Request, name string) (id string, ok bool) {
	raw := mux.Vars(req)[name]
	parsed, err := uuid.FromString(raw)
	if err != nil || parsed == uuid.Nil {
		var invalidIDErr error
		if err != nil {
			invalidIDErr = models.ErrInvalidUUID(err)
		} else {
			invalidIDErr = models.ErrInvalidUUID(fmt.Errorf("invalid %s: zero or nil UUID is not allowed", name))
		}
		if h != nil && h.log != nil {
			h.log.Error(invalidIDErr)
		}
		writeMeshkitError(w, invalidIDErr, http.StatusBadRequest)
		return "", false
	}
	return raw, true
}
