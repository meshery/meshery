package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils/kubernetes/kompose"
)

// MesheryApplicationRequestBody refers to the type of request body that
// SaveMesheryApplication would receive
type MesheryApplicationRequestBody struct {
	URL             string                     `json:"url,omitempty"`
	Path            string                     `json:"path,omitempty"`
	Save            bool                       `json:"save,omitempty"`
	ApplicationData *models.MesheryApplication `json:"application_data,omitempty"`
}

// swagger:route POST /api/application/deploy ApplicationsAPI idPostDeployApplicationFile
// Handle POST request for Application File Deploy
//
// Deploy an attached application file with the request
// responses:
//  200: applicationFilesResponseWrapper

// swagger:route DELETE /api/application/deploy ApplicationsAPI idDeleteApplicationFile
// Handle DELETE request for Application File Deploy
//
// Delete a deployed application file with the request
// responses:
//  200:

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

// swagger:route POST /api/application/ ApplicationsAPI idPostApplicationFileRequest
// Handle POST request for Application Files
//
// Save attached Meshery Application File
// responses:
//  200: mesheryApplicationResponseWrapper

// swagger:route GET /api/application/ ApplicationsAPI idGetApplicationFileRequest
// Handle GET request for Application Files
//
// Returns requests for all Meshery Applications
// responses:
//  200: mesheryApplicationsResponseWrapper

// ApplicationFileRequestHandler will handle requests of both type GET and POST
// on the route /api/application
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

		bytApplication := []byte(mesheryApplication.ApplicationFile)
		// check whether the uploaded file is a docker compose file
		if kompose.IsManifestADockerCompose(bytApplication) {
			if err := kompose.VaildateDockerComposeFile(bytApplication); err != nil {
				obj := "validate"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
				return
			}
			if err := kompose.FormatComposeFile(&bytApplication); err != nil {
				obj := "format"
				h.log.Warn(ErrApplicationFailure(err, obj))
			}
			res, err := kompose.Convert(bytApplication) // convert the docker compose file into kubernetes manifest
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
				return
			}
			mesheryApplication.ApplicationFile = res
		}

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

// swagger:route GET /api/application/{id} ApplicationsAPI idGetMesheryApplication
// Handle GET request for Meshery Application with the given id
//
// Fetches the list of all applications saved by the current user
// responses:
//  200: mesheryApplicationResponseWrapper

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

// swagger:route DELETE /api/application/{id} ApplicationsAPI idDeleteMesheryApplicationFile
// Handle Delete for a Meshery Application File
//
// Deletes a meshery application file with ID: id
// responses:
//  200: noContentWrapper

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
