package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshkit/schemas"
)

// HandleResourceSchemas handles the request to retrieve and merge resource schemas.
// swagger:route GET /api/schema/resource/{resourceName} HandleResourceSchemas
//
// Handles the request to retrieve and merge resource JSON schema and UI schema.
//
// Responses:
//   200
//   500
func (h *Handler) HandleResourceSchemas(rw http.ResponseWriter, r *http.Request) {
	rscName := mux.Vars(r)["resourceName"]
	rjsfSchema, uiSchema, err := schemas.ServeJSonFile(rscName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}


	var rjsfSchemaJSON map[string]interface{}
	var uiSchemaJSON map[string]interface{}

	if err := json.Unmarshal(rjsfSchema, &rjsfSchemaJSON); err != nil {
		http.Error(rw, ErrUnmarshal(err, "RJSF schema").Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(uiSchema, &rjsfSchemaJSON); err != nil {
		uiSchemaJSON = map[string]interface{}{}
	}

	// merge the two data
	mergedData := map[string]interface{}{
		"rjsfSchema": rjsfSchemaJSON,
		"uiSchema":   uiSchemaJSON,
	}

	mergedJSON, err := json.Marshal(mergedData)
	if err != nil {
		http.Error(rw, "Error marshalling json", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(mergedJSON)
	if _, writeErr := rw.Write(mergedJSON); writeErr != nil {
		http.Error(rw, "Error writing response", http.StatusInternalServerError)
		return
	}
}
