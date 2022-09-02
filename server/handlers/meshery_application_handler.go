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

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
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

// swagger:route PUT /api/application/{sourcetype} ApplicationsAPI idPutApplicationFileRequest
// Handle PUT request for Application Files
//
// Updates the design for the provided application
// responses:
//  200: mesheryApplicationResponseWrapper

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
	sourcetype := mux.Vars(r)["sourcetype"]
	if sourcetype == "" {
		http.Error(rw, "missing route variable \"source-type\"", http.StatusBadRequest)
		return
	}
	var parsedBody *MesheryApplicationRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		http.Error(rw, ErrRetrieveData(err).Error(), http.StatusBadRequest)
		// rw.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
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

			pattern, err := core.NewPatternFileFromK8sManifest(k8sres, false)
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
				return
			}
			response, err := yaml.Marshal(pattern)
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
				return
			}
			mesheryApplication.ApplicationFile = string(response)
		} else {
			obj := "convert"
			h.log.Error(ErrApplicationFailure(fmt.Errorf("invalid source type"), obj))
			http.Error(rw, ErrApplicationFailure(fmt.Errorf("invalid source type"), obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
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
				return
			}
			sourceContent, err := io.ReadAll(helmSourceResp.Body)
			if err != nil {
				http.Error(rw, "error read body", http.StatusInternalServerError)
				return
			}

			resp, err := kubernetes.ConvertHelmChartToK8sManifest(kubernetes.ApplyHelmChartConfig{
				URL: parsedBody.URL,
			})
			if err != nil {
				obj := "import"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
				return
			}
			result := string(resp)
			pattern, err := core.NewPatternFileFromK8sManifest(result, false)
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
				return
			}
			response, err := yaml.Marshal(pattern)
			if err != nil {
				obj := "convert"
				h.log.Error(ErrApplicationFailure(err, obj))
				http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
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

				pfs, err := githubRepoApplicationScan(owner, repo, path, branch, sourcetype)
				if err != nil {
					http.Error(rw, ErrRemoteApplication(err).Error(), http.StatusInternalServerError)
					return
				}

				mesheryApplication = &pfs[0]
			} else {
				// Fallback to generic HTTP import
				pfs, err := genericHTTPApplicationFile(parsedBody.URL, sourcetype)
				if err != nil {
					http.Error(rw, ErrRemoteApplication(err).Error(), http.StatusInternalServerError)
					return
				}
				mesheryApplication = &pfs[0]
			}
		} else {
			obj := "convert"
			h.log.Error(ErrApplicationFailure(fmt.Errorf("invalid source type"), obj))
			http.Error(rw, ErrApplicationFailure(fmt.Errorf("invalid source type"), obj).Error(), http.StatusInternalServerError) // sending a 500 when we cannot convert the file into kuberentes manifest
			return
		}
	}
	if parsedBody.Save {
		resp, err := provider.SaveMesheryApplication(token, mesheryApplication)
		if err != nil {
			obj := "save"
			h.log.Error(ErrApplicationFailure(err, obj))
			http.Error(rw, ErrApplicationFailure(err, obj).Error(), http.StatusInternalServerError)
			return
		}

		var mesheryApplicationContent []models.MesheryApplication
		err = json.Unmarshal(resp, &mesheryApplicationContent)
		if err != nil {
			obj := "application"
			h.log.Error(ErrEncoding(err, obj))
			http.Error(rw, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
			return
		}

		err = provider.SaveApplicationSourceContent(token, (mesheryApplicationContent[0].ID).String(), mesheryApplication.SourceContent)

		if err != nil {
			obj := "upload"
			h.log.Error(ErrApplicationSourceContent(err, obj))
			http.Error(rw, ErrApplicationSourceContent(err, obj).Error(), http.StatusInternalServerError)
			return
		}

		go h.config.ConfigurationChannel.PublishApplications()
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

func (h *Handler) handleApplicationUpdate(rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider) {
	defer func() {
		_ = r.Body.Close()
	}()
	sourcetype := mux.Vars(r)["sourcetype"]
	if sourcetype == "" {
		http.Error(rw, "missing route variable \"source-type\"", http.StatusBadRequest)
		return
	}

	var parsedBody *MesheryApplicationRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		http.Error(rw, ErrRetrieveData(err).Error(), http.StatusBadRequest)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		return
	}

	format := r.URL.Query().Get("output")
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
		return
	}

	go h.config.ConfigurationChannel.PublishApplications()

	h.formatApplicationOutput(rw, resp, format)
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
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetMesheryApplications(tokenString, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"))
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

	go h.config.ConfigurationChannel.PublishApplications()
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

// swagger:route GET /api/application/types ApplicationsAPI typeGetMesheryApplication
// Handle GET request for Meshery Application with the provided type
//
// Get application file type
// responses:
//  200: mesheryApplicationResponseWrapper

// GetMesheryApplicationHandler fetched the application with the given id
func (h *Handler) GetMesheryApplicationTypesHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
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

// swagger:route GET /api/application/download/{id}/{sourcetype} ApplicationsAPI typeGetApplication
// Handle GET request for Meshery Application with of source content
//
// Get the application source-content
// responses:
//  200: applicationFileParamsWrapper

// GetMesheryApplicationHandler fetched the application with the given id
func (h *Handler) GetMesheryApplicationSourceHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
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

	data, err := json.Marshal(&contentMesheryApplicationSlice)
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

func githubRepoApplicationScan(
	owner,
	repo,
	path,
	branch,
	sourceType string,
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
						return err
					}
				}
				pattern, err := core.NewPatternFileFromK8sManifest(k8sres, false)
				if err != nil {
					return err
				}
				response, err := yaml.Marshal(pattern)
				if err != nil {
					return err
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

	return result, err
}

func genericHTTPApplicationFile(fileURL, sourceType string) ([]models.MesheryApplication, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file not found")
	}

	defer models.SafeClose(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	k8sres := string(body)

	if sourceType == string(models.DockerCompose) {
		k8sres, err = kompose.Convert(body)
		if err != nil {
			return nil, err
		}
	}

	pattern, err := core.NewPatternFileFromK8sManifest(k8sres, false)
	if err != nil {
		return nil, err
	}
	response, err := yaml.Marshal(pattern)

	if err != nil {
		return nil, err
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
