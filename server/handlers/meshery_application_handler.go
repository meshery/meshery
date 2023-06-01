package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	guid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshery/server/models"
	pCore "github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/kubernetes/kompose"
	"github.com/layer5io/meshkit/utils/walker"
	"gopkg.in/yaml.v2"
)

// MesheryApplicationRequestBody refers to the type of request body that
// SaveMesheryApplication would receive
type MesheryApplicationRequestBody struct {
	URL             string                     `json:"url,omitempty"`
	Path            string                     `json:"path,omitempty"`
	Save            bool                       `json:"save,omitempty"`
	ApplicationData *models.MesheryApplication `json:"application_data,omitempty"`
	CytoscapeJSON   string                     `json:"cytoscape_json,omitempty"`
	Name            string                     `json:"name,omitempty"`
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

// swagger:route POST /api/application/{sourcetype} ApplicationsAPI idPostApplicationFileRequest
// Handle POST request for Application Files
//
// Creates a new application with source-content
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

	if r.Method == http.MethodPut {
		h.handleApplicationUpdate(rw, r, prefObj, user, provider)
	}
}

// swagger:route POST /api/application/{sourcetype} ApplicationsAPI idPutApplicationFileRequest
// Handle POST request for Application Files
//
// Updates the design for the provided application
// responses:
//  200: mesheryApplicationResponseWrapper

func (h *Handler) handleApplicationPOST(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()
	res := meshes.EventsResponse{
		Component:     "core",
		ComponentName: "Application",
		OperationId:   guid.NewString(),
		EventType:     meshes.EventType_INFO,
	}
	sourcetype := mux.Vars(r)["sourcetype"]
	if sourcetype == "" {
		http.Error(rw, "missing route variable \"source-type\"", http.StatusBadRequest)
		addMeshkitErr(&res, ErrSaveApplication(fmt.Errorf("missing route variable \"source-type\"")))
		go h.EventsBuffer.Publish(&res)
		return
	}
	var parsedBody *MesheryApplicationRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		http.Error(rw, ErrRetrieveData(err).Error(), http.StatusBadRequest)
		addMeshkitErr(&res, ErrRetrieveData(err))
		go h.EventsBuffer.Publish(&res)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		addMeshkitErr(&res, ErrRetrieveData(err))
		go h.EventsBuffer.Publish(&res)
		return
	}

	format := r.URL.Query().Get("output")
	var mesheryApplication *models.MesheryApplication
	// If Content is not empty then assume it's a local upload
	//Note: The Application data will not be present in case of helm charts as we do not support local helm upload.
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

		mesheryApplication = parsedBody.ApplicationData

		bytApplication := []byte(mesheryApplication.ApplicationFile)
		mesheryApplication.SourceContent = bytApplication
		if sourcetype == string(models.DockerCompose) || sourcetype == string(models.K8sManifest) {
			var k8sres string
			if sourcetype == string(models.DockerCompose) {
				k8sres, err = kompose.Convert(bytApplication) // convert the docker compose file into kubernetes manifest
				if err != nil {
					obj := "convert"
					h.log.Error(ErrApplicationFailure(err, obj))
					http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
					addMeshkitErr(&res, ErrApplicationFailure(err, obj))
					go h.EventsBuffer.Publish(&res)
					return
				}
				mesheryApplication.Type = sql.NullString{
					String: string(models.DockerCompose),
					Valid:  true,
				}
			} else if sourcetype == string(models.K8sManifest) {
				k8sres = string(bytApplication)
				mesheryApplication.Type = sql.NullString{
					String: string(models.K8sManifest),
					Valid:  true,
				}
			}
			pattern, err := pCore.NewPatternFileFromK8sManifest(k8sres, false, h.registryManager)
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
				addMeshkitErr(&res, err)                                                                //this error is already a meshkit error so no further wrapping required
				go h.EventsBuffer.Publish(&res)
				return
			}
			response, err := yaml.Marshal(pattern)
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
				addMeshkitErr(&res, ErrApplicationFailure(err, obj))
				go h.EventsBuffer.Publish(&res)
				return
			}
			mesheryApplication.ApplicationFile = string(response)
		} else {
			obj := "convert"
			h.log.Error(ErrApplicationFailure(fmt.Errorf("invalid source type"), obj))
			http.Error(rw, ErrApplicationFailure(fmt.Errorf("invalid source type"), obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
			addMeshkitErr(&res, ErrApplicationFailure(err, obj))
			go h.EventsBuffer.Publish(&res)
			return
		}
	}

	if parsedBody.URL != "" {
		if sourcetype == string(models.HelmChart) {
			helmSourceResp, err := http.Get(parsedBody.URL)
			defer func() {
				_ = helmSourceResp.Body.Close()
			}()
			if err != nil {
				obj := "import"
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
				addMeshkitErr(&res, ErrApplicationFailure(err, obj))
				go h.EventsBuffer.Publish(&res)
				return
			}
			sourceContent, err := io.ReadAll(helmSourceResp.Body)
			if err != nil {
				http.Error(rw, "error read body", http.StatusInternalServerError)
				addMeshkitErr(&res, ErrSaveApplication(fmt.Errorf("error reading body")))
				go h.EventsBuffer.Publish(&res)
				return
			}

			resp, err := kubernetes.ConvertHelmChartToK8sManifest(kubernetes.ApplyHelmChartConfig{
				URL: parsedBody.URL,
			})
			if err != nil {
				obj := "import"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
				addMeshkitErr(&res, ErrApplicationFailure(err, obj))
				go h.EventsBuffer.Publish(&res)
				return
			}
			result := string(resp)
			pattern, err := pCore.NewPatternFileFromK8sManifest(result, false, h.registryManager)
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
				addMeshkitErr(&res, err)
				go h.EventsBuffer.Publish(&res)
				return
			}
			response, err := yaml.Marshal(pattern)
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
				addMeshkitErr(&res, ErrApplicationFailure(err, obj))
				go h.EventsBuffer.Publish(&res)
				return
			}
			url := strings.Split(parsedBody.URL, "/")
			mesheryApplication = &models.MesheryApplication{
				Name:            strings.TrimSuffix(url[len(url)-1], ".tgz"),
				ApplicationFile: string(response),
				Type: sql.NullString{
					String: string(models.HelmChart),
					Valid:  true,
				},
				Location: map[string]interface{}{
					"type":   "http",
					"host":   parsedBody.URL,
					"path":   "",
					"branch": "",
				},
				SourceContent: sourceContent,
			}
		} else if sourcetype == string(models.DockerCompose) || sourcetype == string(models.K8sManifest) {
			parsedURL, err := url.Parse(parsedBody.URL)
			if err != nil {
				http.Error(rw, "error parsing provided URL", http.StatusInternalServerError)
				addMeshkitErr(&res, ErrSaveApplication(fmt.Errorf("error parsing URL")))
				go h.EventsBuffer.Publish(&res)
				return
			}

			// Check if hostname is github
			if parsedURL.Host == "github.com" {
				parsedPath := strings.Split(parsedURL.Path, "/")
				if parsedPath[3] == "tree" {
					parsedPath = append(parsedPath[0:3], parsedPath[4:]...)
				}
				if len(parsedPath) < 3 {
					http.Error(rw, "malformed URL: url should be of type github.com/<owner>/<repo>/[branch]", http.StatusNotAcceptable)
				}

				owner := parsedPath[1]
				repo := parsedPath[2]
				branch := "master"
				path := parsedBody.Path
				if len(parsedPath) == 4 {
					branch = parsedPath[3]
				}
				if path == "" && len(parsedPath) > 4 {
					path = strings.Join(parsedPath[4:], "/")
				}

				pfs, err := githubRepoApplicationScan(owner, repo, path, branch, sourcetype, h.registryManager)
				if err != nil {
					http.Error(rw, ErrRemoteApplication(err).Error(), http.StatusInternalServerError)
					addMeshkitErr(&res, err) //error guaranteed to be meshkit error
					go h.EventsBuffer.Publish(&res)
					return
				}

				mesheryApplication = &pfs[0]
			} else {
				// Fallback to generic HTTP import
				pfs, err := genericHTTPApplicationFile(parsedBody.URL, sourcetype, h.registryManager)
				if err != nil {
					http.Error(rw, ErrRemoteApplication(err).Error(), http.StatusInternalServerError)
					addMeshkitErr(&res, err) //error guaranteed to be meshkit error
					go h.EventsBuffer.Publish(&res)
					return
				}
				mesheryApplication = &pfs[0]
			}
		} else {
			obj := "convert"
			h.log.Error(ErrApplicationFailure(fmt.Errorf("invalid source type"), obj))
			http.Error(rw, ErrApplicationFailure(fmt.Errorf("invalid source type"), obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
			addMeshkitErr(&res, ErrApplicationFailure(fmt.Errorf("error parsing URL"), obj))
			go h.EventsBuffer.Publish(&res)
			return
		}
	}

	var savedApplicationID *uuid.UUID

	if parsedBody.Save {
		resp, err := provider.SaveMesheryApplication(token, mesheryApplication)
		if err != nil {
			obj := "save"
			h.log.Error(ErrApplicationFailure(err, obj))
			http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
			addMeshkitErr(&res, ErrApplicationFailure(err, obj))
			go h.EventsBuffer.Publish(&res)
			return
		}

		var mesheryApplicationContent []models.MesheryApplication
		err = json.Unmarshal(resp, &mesheryApplicationContent)
		if err != nil {
			obj := "application"
			h.log.Error(ErrEncoding(err, obj))
			http.Error(rw, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
			addMeshkitErr(&res, ErrEncoding(err, obj))
			go h.EventsBuffer.Publish(&res)
			return
		}
		savedApplicationID = mesheryApplicationContent[0].ID
		err = provider.SaveApplicationSourceContent(token, (savedApplicationID).String(), mesheryApplication.SourceContent)

		if err != nil {
			obj := "upload"
			h.log.Error(ErrApplicationSourceContent(err, obj))
			http.Error(rw, ErrApplicationSourceContent(err, obj).Error(), http.StatusInternalServerError)
			addMeshkitErr(&res, ErrApplicationSourceContent(err, obj))
			go h.EventsBuffer.Publish(&res)
			return
		}

		go h.config.ConfigurationChannel.PublishApplications()
	}

	mesheryApplication.ID = savedApplicationID
	byt, err := json.Marshal([]models.MesheryApplication{*mesheryApplication})
	if err != nil {
		obj := "application"
		h.log.Error(ErrEncoding(err, obj))
		http.Error(rw, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		addMeshkitErr(&res, ErrEncoding(err, obj))
		go h.EventsBuffer.Publish(&res)
		return
	}

	h.formatApplicationOutput(rw, byt, format, &res)
	return
}

func (h *Handler) handleApplicationUpdate(rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider) {
	defer func() {
		_ = r.Body.Close()
	}()
	res := meshes.EventsResponse{
		Component:     "core",
		ComponentName: "Design",
		OperationId:   guid.NewString(),
		EventType:     meshes.EventType_INFO,
	}
	sourcetype := mux.Vars(r)["sourcetype"]
	if sourcetype == "" {
		http.Error(rw, "missing route variable \"source-type\"", http.StatusBadRequest)
		addMeshkitErr(&res, ErrSaveApplication(fmt.Errorf("missing route \"source-type\"")))
		go h.EventsBuffer.Publish(&res)
		return
	}

	var parsedBody *MesheryApplicationRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		http.Error(rw, ErrRetrieveData(err).Error(), http.StatusBadRequest)
		addMeshkitErr(&res, ErrRetrieveData(err))
		go h.EventsBuffer.Publish(&res)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		addMeshkitErr(&res, ErrRetrieveUserToken(err))
		go h.EventsBuffer.Publish(&res)
		return
	}
	format := r.URL.Query().Get("output")

	if parsedBody.CytoscapeJSON != "" {
		pf, err := pCore.NewPatternFileFromCytoscapeJSJSON(parsedBody.Name, []byte(parsedBody.CytoscapeJSON))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "%s", err)
			addMeshkitErr(&res, ErrSavePattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		pfByt, err := pf.ToYAML()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "%s", err)
			addMeshkitErr(&res, ErrSavePattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		patternName, err := models.GetPatternName(string(pfByt))
		if err != nil {
			h.log.Error(ErrGetPattern(err))
			http.Error(rw, ErrGetPattern(err).Error(), http.StatusBadRequest)
			addMeshkitErr(&res, ErrGetPattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		mesheryApplication := &models.MesheryApplication{
			Name:            patternName,
			ApplicationFile: string(pfByt),
			Location: map[string]interface{}{
				"host": "",
				"path": "",
				"type": "local",
			},
			Type: sql.NullString{
				String: sourcetype,
				Valid:  true,
			},
		}
		if parsedBody.ApplicationData != nil {
			mesheryApplication.ID = parsedBody.ApplicationData.ID
		}
		if parsedBody.Save {
			resp, err := provider.SaveMesheryApplication(token, mesheryApplication)
			if err != nil {
				h.log.Error(ErrSavePattern(err))
				http.Error(rw, ErrSavePattern(err).Error(), http.StatusInternalServerError)
				addMeshkitErr(&res, ErrSavePattern(err))
				go h.EventsBuffer.Publish(&res)
				return
			}

			go h.config.ConfigurationChannel.PublishApplications()
			h.formatApplicationOutput(rw, resp, format, &res)
			return
		}

		byt, err := json.Marshal([]models.MesheryApplication{*mesheryApplication})
		if err != nil {
			h.log.Error(ErrEncodePattern(err))
			http.Error(rw, ErrEncodePattern(err).Error(), http.StatusInternalServerError)
			addMeshkitErr(&res, ErrEncodePattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		h.formatApplicationOutput(rw, byt, format, &res)
		return
	}
	mesheryApplication := parsedBody.ApplicationData
	mesheryApplication.Type = sql.NullString{
		String: sourcetype,
		Valid:  true,
	}
	resp, err := provider.SaveMesheryApplication(token, mesheryApplication)
	if err != nil {
		obj := "save"
		h.log.Error(ErrApplicationFailure(err, obj))
		http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
		addMeshkitErr(&res, ErrApplicationFailure(err, obj))
		go h.EventsBuffer.Publish(&res)
		return
	}

	go h.config.ConfigurationChannel.PublishApplications()

	h.formatApplicationOutput(rw, resp, format, &res)
}

// swagger:route GET /api/application/{id} ApplicationsAPI idGetMesheryApplication
// Handle GET request for Meshery Application with the given id
//
// Fetches the list of all applications saved by the current user
//
// ```?updated_after=<timestamp>``` timestamp should be of the format "YYYY-MM-DD HH:MM:SS"
//
// ```?order={field}``` orders on the passed field
//
// ```?search=<application name>``` A string matching is done on the specified application name
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 10
//
// responses:
//
//  200: mesheryApplicationResponseWrapper

// GetMesheryApplicationsHandler returns the list of all the applications saved by the current user
func (h *Handler) GetMesheryApplicationsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetMesheryApplications(tokenString, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"), q.Get("updated_after"))
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
	_ *models.Preference,
	_ *models.User,
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

	go h.config.ConfigurationChannel.PublishApplications()
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// GetMesheryApplicationHandler fetched the application with the given id
func (h *Handler) GetMesheryApplicationHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
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

// swagger:route GET /api/application/types ApplicationsAPI typeGetMesheryApplication
// Handle GET request for Meshery Application with the provided type
//
// Get application file type
// responses:
//  200: mesheryApplicationResponseWrapper

// GetMesheryApplicationHandler fetched the application with the given id
func (h *Handler) GetMesheryApplicationTypesHandler(
	rw http.ResponseWriter,
	_ *http.Request,
	_ *models.Preference,
	_ *models.User,
	_ models.Provider,
) {
	response := models.GetApplicationTypes()
	b, err := json.Marshal(response)
	if err != nil {
		obj := "available types"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(rw, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(b))
}

// swagger: route GET /api/application/download/{id} ApplicationsAPI idGetApplication
// Handle GET request for Meshery Application with the given id
//
// Get the application file with the given id
// responses:
//  200

// GetMesheryApplicationFile returns the application file with the given id
func (h *Handler) GetMesheryApplicationFile(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	applicationID := mux.Vars(r)["id"]
	resp, err := provider.GetMesheryApplication(r, applicationID)

	if err != nil {
		obj := "download"
		h.log.Error(ErrApplicationFailure(err, obj))
		http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusNotFound)
		return
	}

	application := &models.MesheryApplication{}

	err = json.Unmarshal(resp, &application)
	if err != nil {
		h.log.Error(ErrApplicationFailure(err, "parse failure"))
		http.Error(rw, ErrApplicationFailure(err, "parse failure").Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/x-yaml")
	if _, err := io.Copy(rw, strings.NewReader(application.ApplicationFile)); err != nil {
		h.log.Error(ErrApplicationSourceContent(err, "download"))
		http.Error(rw, ErrApplicationSourceContent(err, "download").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/application/download/{id}/{sourcetype} ApplicationsAPI typeGetApplication
// Handle GET request for Meshery Application with of source content
//
// Get the application source-content
// responses:
//  200

// GetMesheryApplicationHandler fetched the application with the given id
func (h *Handler) GetMesheryApplicationSourceHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	applicationID := mux.Vars(r)["id"]
	resp, err := provider.GetApplicationSourceContent(r, applicationID)
	if err != nil {
		obj := "download"
		h.log.Error(ErrApplicationFailure(err, obj))
		http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusNotFound)
		return
	}

	var mimeType string
	sourcetype := mux.Vars(r)["sourcetype"]

	if models.ApplicationType(sourcetype) == models.HelmChart { //serve the content in a tgz file
		mimeType = "application/x-tar"
	} else { // serve the content in yaml file
		mimeType = "application/x-yaml"
	}
	reader := bytes.NewReader(resp)
	rw.Header().Set("Content-Type", mimeType)
	_, err = io.Copy(rw, reader)
	if err != nil {
		h.log.Error(ErrApplicationSourceContent(err, "download"))
		http.Error(rw, ErrApplicationSourceContent(err, "download").Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) formatApplicationOutput(rw http.ResponseWriter, content []byte, _ string, res *meshes.EventsResponse) {
	contentMesheryApplicationSlice := make([]models.MesheryApplication, 0)
	names := []string{}

	if err := json.Unmarshal(content, &contentMesheryApplicationSlice); err != nil {
		obj := "application data into go slice"
		h.log.Error(ErrDecoding(err, obj))
		http.Error(rw, ErrDecoding(err, obj).Error(), http.StatusInternalServerError)
		addMeshkitErr(res, ErrDecoding(err, obj))
		go h.EventsBuffer.Publish(res)
		return
	}

	data, err := json.Marshal(&contentMesheryApplicationSlice)
	if err != nil {
		obj := "application file"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(rw, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		addMeshkitErr(res, ErrMarshal(err, obj))
		go h.EventsBuffer.Publish(res)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(data))
	for _, app := range contentMesheryApplicationSlice {
		names = append(names, app.Name)
	}
	res.Details = "\"" + strings.Join(names, ",") + "\" application saved"
	res.Summary = "Changes to the \"" + strings.Join(names, ",") + "\" application have been saved."
	// go h.EventsBuffer.Publish(res)
}

// Note: This function is guaranteed to return meshkit errors
func githubRepoApplicationScan(
	owner,
	repo,
	path,
	branch,
	sourceType string,
	reg *meshmodel.RegistryManager,
) ([]models.MesheryApplication, error) {
	var mu sync.Mutex
	ghWalker := walker.NewGit()
	result := make([]models.MesheryApplication, 0)
	err := ghWalker.
		Owner(owner).
		Repo(repo).
		Branch(branch).
		Root(path).
		RegisterFileInterceptor(func(f walker.File) error {
			ext := filepath.Ext(f.Name)
			var k8sres string
			var err error
			k8sres = f.Content
			if ext == ".yml" || ext == ".yaml" {
				if sourceType == string(models.DockerCompose) {
					k8sres, err = kompose.Convert([]byte(f.Content))
					if err != nil {
						return ErrRemoteApplication(err)
					}
				}
				pattern, err := pCore.NewPatternFileFromK8sManifest(k8sres, false, reg)
				if err != nil {
					return err //always a meshkit error
				}
				response, err := yaml.Marshal(pattern)
				if err != nil {
					return ErrMarshal(err, string(response))
				}

				af := models.MesheryApplication{
					Name:            strings.TrimSuffix(f.Name, ext),
					ApplicationFile: string(response),
					Location: map[string]interface{}{
						"type":   "github",
						"host":   fmt.Sprintf("github.com/%s/%s", owner, repo),
						"path":   f.Path,
						"branch": branch,
					},
					Type: sql.NullString{
						String: string(sourceType),
						Valid:  true,
					},
					SourceContent: []byte(f.Content),
				}

				mu.Lock()
				result = append(result, af)
				mu.Unlock()
			}

			return nil
		}).
		Walk()

	return result, ErrRemoteApplication(err)
}

// Note: Always return meshkit error from this function
func genericHTTPApplicationFile(fileURL, sourceType string, reg *meshmodel.RegistryManager) ([]models.MesheryApplication, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, ErrRemoteApplication(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrRemoteApplication(fmt.Errorf("file not found"))
	}

	defer models.SafeClose(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrRemoteApplication(err)
	}

	k8sres := string(body)

	if sourceType == string(models.DockerCompose) {
		k8sres, err = kompose.Convert(body)
		if err != nil {
			return nil, ErrRemoteApplication(err)
		}
	}

	pattern, err := pCore.NewPatternFileFromK8sManifest(k8sres, false, reg)
	if err != nil {
		return nil, err //This error is already a meshkit error
	}
	response, err := yaml.Marshal(pattern)

	if err != nil {
		return nil, ErrMarshal(err, string(response))
	}

	url := strings.Split(fileURL, "/")
	af := models.MesheryApplication{
		Name:            url[len(url)-1],
		ApplicationFile: string(response),
		Location: map[string]interface{}{
			"type":   "http",
			"host":   fileURL,
			"path":   "",
			"branch": "",
		},
		Type: sql.NullString{
			String: string(sourceType),
			Valid:  true,
		},
		SourceContent: body,
	}
	return []models.MesheryApplication{af}, nil
}
