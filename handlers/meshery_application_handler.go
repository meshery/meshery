package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
)

// MesheryApplicationRequestBody refers to the type of request body that
// SaveMesheryApplication would receive
type MesheryApplicationRequestBody struct {
	URL             string                     `json:"url,omitempty"`
	Path            string                     `json:"path,omitempty"`
	Save            bool                       `json:"save,omitempty"`
	ApplicationData *models.MesheryApplication `json:"application_data,omitempty"`
}

// ApplicationFileHandler handles the requested related to application files
func (h *Handler) ApplicationFileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	// Application files are just pattern files
	h.PatternFileHandler(rw, r, prefObj, user, provider)
}

// ApplicationFileRequestHandler will handle requests of both type GET and POST
// on the route /api/experimental/application
func (h *Handler) ApplicationFileRequestHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	if r.Method == http.MethodGet {
		h.GetMesheryApplicationsHandler(rw, r, prefObj, user, provider)
		return
	}

	if r.Method == http.MethodPost {
		h.handleApplicationPOST(rw, r, prefObj, user, provider)
		return
	}
}

func (h *Handler) handleApplicationPOST(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var parsedBody *MesheryApplicationRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRetrieveData(err))
		http.Error(rw, ErrRetrieveData(err).Error(), http.StatusBadRequest)
		// rw.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		return
	}

	format := r.URL.Query().Get("output")

	// If Content is not empty then assume it's a local upload
	if parsedBody.ApplicationData != nil {
		applicationName, err := models.GetApplicationName(parsedBody.ApplicationData.ApplicationFile)
		if err != nil {
			obj := "save"
			h.log.Error(ErrApplicationFailure(err, obj))
			http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusBadRequest)
			return
		}

		// Assign a name if no name is provided
		if parsedBody.ApplicationData.Name == "" {
			parsedBody.ApplicationData.Name = applicationName
		}
		// Assign a location if no location is specified
		if parsedBody.ApplicationData.Location == nil {
			parsedBody.ApplicationData.Location = map[string]interface{}{
				"host":   "",
				"path":   "",
				"type":   "local",
				"branch": "",
			}
		}

		mesheryApplication := parsedBody.ApplicationData

		if parsedBody.Save {
			resp, err := provider.SaveMesheryApplication(token, mesheryApplication)
			if err != nil {
				obj := "save"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
				return
			}

			h.formatApplicationOutput(rw, resp, format)
			return
		}

		byt, err := json.Marshal([]models.MesheryApplication{*mesheryApplication})
		if err != nil {
			obj := "application"
			h.log.Error(ErrEncoding(err, obj))
			http.Error(rw, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
			return
		}

		h.formatApplicationOutput(rw, byt, format)
		return
	}

	if parsedBody.URL != "" {
		resp, err := provider.RemoteApplicationFile(r, parsedBody.URL, parsedBody.Path, parsedBody.Save)

		if err != nil {
			obj := "import"
			h.log.Error(ErrApplicationFailure(err, obj))
			http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
			return
		}

		h.formatApplicationOutput(rw, resp, format)
		return
	}
}

// GetMesheryApplicationsHandler returns the list of all the applications saved by the current user
func (h *Handler) GetMesheryApplicationsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()

	resp, err := provider.GetMesheryApplications(r, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"))
	if err != nil {
		obj := "fetch"
		h.log.Error(ErrApplicationFailure(err, obj))
		http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// DeleteMesheryApplicationHandler deletes a application with the given id
func (h *Handler) DeleteMesheryApplicationHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	applicationID := mux.Vars(r)["id"]

	resp, err := provider.DeleteMesheryApplication(r, applicationID)
	if err != nil {
		obj := "delete"
		h.log.Error(ErrApplicationFailure(err, obj))
		http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// GetMesheryApplicationHandler fetched the application with the given id
func (h *Handler) GetMesheryApplicationHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	applicationID := mux.Vars(r)["id"]

	resp, err := provider.GetMesheryApplication(r, applicationID)
	if err != nil {
		obj := "get"
		h.log.Error(ErrApplicationFailure(err, obj))
		http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

func (h *Handler) formatApplicationOutput(rw http.ResponseWriter, content []byte, format string) {
	contentMesheryApplicationSlice := make([]models.MesheryApplication, 0)

	if err := json.Unmarshal(content, &contentMesheryApplicationSlice); err != nil {
		obj := "application data into go slice"
		h.log.Error(ErrDecoding(err, obj))
		http.Error(rw, ErrDecoding(err, obj).Error(), http.StatusInternalServerError)
		// rw.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprintf(rw, "failed to decode applications data into go slice: %s", err)
		return
	}

	result := []models.MesheryApplication{}

	data, err := json.Marshal(&result)
	if err != nil {
		obj := "application file"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(rw, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		//rw.WriteHeader(http.StatusInternalServerError)
		//fmt.Fprintf(rw, "failed to marshal application file: %s", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(data))
}
