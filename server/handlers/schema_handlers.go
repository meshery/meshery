package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/schemas"
)

func (h *Handler) HandleResourceSchemas(rw http.ResponseWriter, r *http.Request) {
	rscName := mux.Vars(r)["resourceName"]
	rjsfSchema, uiSchema, err := schemas.ServeJSonFile(rscName)
	if err != nil {
		writeMeshkitError(rw, ErrServeSchema(err), http.StatusInternalServerError)
		return
	}

	var rjsfSchemaJSON map[string]interface{}
	var uiSchemaJSON map[string]interface{}

	if err := json.Unmarshal(rjsfSchema, &rjsfSchemaJSON); err != nil {
		writeMeshkitError(rw, models.ErrUnmarshal(err, "RJSF schema"), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(uiSchema, &uiSchemaJSON); err != nil {
		uiSchemaJSON = map[string]interface{}{}
	}

	// merge the two data
	mergedData := map[string]interface{}{
		"rjsfSchema": rjsfSchemaJSON,
		"uiSchema":   uiSchemaJSON,
	}

	mergedJSON, err := json.Marshal(mergedData)
	if err != nil {
		writeMeshkitError(rw, models.ErrMarshal(err, "merged RJSF + UI schema response"), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, writeErr := rw.Write(mergedJSON); writeErr != nil {
		writeMeshkitError(rw, ErrWriteResponse(writeErr), http.StatusInternalServerError)
		return
	}
}
