// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"strconv"
	"strings"
	"sync"
	"time"

	"fortio.org/fortio/periodic"
	yaml "github.com/ghodss/yaml"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

// LoadTestUsingSMPHandler runs the load test with the given parameters and SMP
func (h *Handler) LoadTestUsingSMPHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// Read the SMP File
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	if req.Header.Get("Content-Type") == "application/json" {
		body, err = yaml.JSONToYAML(body)
		if err != nil {
			h.log.Error(ErrPatternFile(err))
			writeMeshkitError(w, ErrPatternFile(err), http.StatusInternalServerError)
			return
		}
	}

	jsonBytes, _ := yaml.YAMLToJSON(body)

	perfTest := &models.PerformanceTestConfigFile{}
	if err := json.Unmarshal(jsonBytes, perfTest); err != nil {
		h.log.Error(ErrParseBool(err, "provided input"))
		writeMeshkitError(w, ErrParseBool(err, "provided input"), http.StatusBadRequest)
		return
	}

	// Guard against a missing/empty "test" config before dereferencing it
	// below (name, id, duration, clients[0]); a malformed body must not panic.
	if perfTest == nil || perfTest.Config == nil || len(perfTest.Config.Clients) == 0 {
		err := fmt.Errorf("request body is missing the required \"test\" performance configuration or load-test clients")
		h.log.Error(models.ErrUnmarshal(err, "performance test config"))
		writeMeshkitError(w, models.ErrUnmarshal(err, "performance test config"), http.StatusBadRequest)
		return
	}

	// testName - should be loaded from the file and updated with a random string appended to the end of the name
	testName := perfTest.Config.Name
	if testName == "" {
		h.log.Error(ErrBlankName(err))
		writeMeshkitError(w, ErrBlankName(err), http.StatusForbidden)
		return
	}

	// The technology under test is identified by a Meshery Registry model
	// name (free-form string), replacing the constraining SMP service-mesh enum.
	meshName := perfTest.ServiceMesh

	profileID := perfTest.Config.ID

	loadTestOptions := &models.LoadTestOptions{}

	testDuration, err := time.ParseDuration(perfTest.Config.Duration)
	if err != nil {
		h.log.Error(ErrParseDuration)
		writeMeshkitError(w, ErrParseDuration, http.StatusBadRequest)
		return
	}
	loadTestOptions.Duration = testDuration
	if loadTestOptions.Duration.Seconds() <= 0 {
		loadTestOptions.Duration = time.Second
	}

	// TODO: check multiple clients in case of distributed perf test
	testClient := perfTest.Config.Clients[0]
	if len(testClient.EndpointUrls) == 0 {
		err := fmt.Errorf("request body is missing load-test endpoint URLs")
		h.log.Error(models.ErrUnmarshal(err, "performance test config"))
		writeMeshkitError(w, models.ErrUnmarshal(err, "performance test config"), http.StatusBadRequest)
		return
	}

	// TODO: consider the multiple endpoints
	loadTestOptions.URL = testClient.EndpointUrls[0]
	loadTestOptions.HTTPNumThreads = int(testClient.Connections)
	loadTestOptions.HTTPQPS = float64(testClient.Rps)

	if loadTestOptions.HTTPNumThreads < 1 {
		loadTestOptions.HTTPNumThreads = 1
	}

	ltURL, err := url.Parse(loadTestOptions.URL)
	if err != nil {
		obj := "the provided load test"
		h.log.Error(ErrParseBool(err, obj))
		writeMeshkitError(w, ErrParseBool(err, obj), http.StatusBadRequest)
		return
	}
	if !ltURL.IsAbs() {
		h.log.Error(ErrInvalidLTURL(ltURL.String()))
		writeMeshkitError(w, ErrInvalidLTURL(ltURL.String()), http.StatusBadRequest)
		return
	}
	loadTestOptions.Name = testName

	if loadTestOptions.HTTPQPS < 0 {
		loadTestOptions.HTTPQPS = 0
	}

	// fortio is the only supported load generator. Any generator carried by
	// an existing profile (including the removed "wrk2") runs on fortio.
	loadTestOptions.LoadGenerator = models.FortioLG
	loadTestOptions.AllowInitialErrors = true

	h.loadTestHelperHandler(w, req, profileID, testName, meshName, "", prefObj, loadTestOptions, provider)
}

func (h *Handler) jsonToMap(headersString string) *map[string]string {
	headers := make(map[string]string)
	err := json.Unmarshal([]byte(headersString), &headers)
	if err != nil {
		return nil
	}
	return &headers
}

// LoadTestHandler runs the load test with the given parameters
func (h *Handler) LoadTestHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	cleanUpFiles := make([]string, 0)

	defer func() {
		for _, file := range cleanUpFiles {
			err := os.Remove(file)
			if err != nil {
				h.log.Error(ErrCleanupCertificate(err, file))
			}
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	isSSLCertificateProvided := req.URL.Query().Get("cert") == queryParamTrue

	/*
	 When "cert" query param is set the body contains self-signed certs
	 and not the SMP config, hence we shouldn't use SMP Handler,
	 if query param is unset/not present presence of body
	 if values have been passed as body we run test using SMP Handler
	*/
	if !isSSLCertificateProvided && string(body) != "" {
		h.log.Info("Running test with SMP config")
		req.Body = io.NopCloser(strings.NewReader(string(body)))
		h.LoadTestUsingSMPHandler(w, req, prefObj, user, provider)
		return
	}

	loadTestOptions := &models.LoadTestOptions{}

	profileID := mux.Vars(req)["id"]

	performanceProfileData, err := provider.GetPerformanceProfile(req, profileID)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, ErrFetchProfile(err), http.StatusInternalServerError)
		return
	}

	performanceProfile := models.PerformanceProfile{}
	err = json.Unmarshal(performanceProfileData, &performanceProfile)
	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, "performance profile"))
		writeMeshkitError(w, models.ErrUnmarshal(err, "performance profile"), http.StatusBadRequest)
		return
	}

	options, ok := performanceProfile.Metadata["additional_options"].(string)

	if ok {
		loadTestOptions.Options = options
	}

	if isSSLCertificateProvided {
		var isErrLoadingCertificate bool
		caCertificate, certificateOk := performanceProfile.Metadata["ca_certificate"].(map[string]interface{})

		if certificateOk {
			certificateName, ok := caCertificate["name"].(string)
			if ok {
				filePath := utils.SanitizeFileName(certificateName)

				file, err := os.CreateTemp(os.TempDir(), filePath)
				if err != nil {
					h.log.Error(ErrCreateFile(err, certificateName))
					writeMeshkitError(w, ErrCreateFile(err, certificateName), http.StatusInternalServerError)
					return
				}
				cleanUpFiles = append(cleanUpFiles, file.Name())
				certificateData := caCertificate["file"].(string)

				_, err = file.Write([]byte(certificateData))
				if err != nil {
					h.log.Error(ErrCreateFile(err, certificateName))
					writeMeshkitError(w, ErrCreateFile(err, certificateName), http.StatusInternalServerError)
					return
				}
				assignCertificatePath("ca_certificate", file.Name(), loadTestOptions)
			} else {
				isErrLoadingCertificate = true
			}

		} else {
			isErrLoadingCertificate = true
		}

		if isErrLoadingCertificate {
			h.log.Info("unable to load SSL certificate, skipping")
		}
	}

	err = req.ParseForm()
	if err != nil {
		obj := "form"
		h.log.Error(ErrParseBool(err, obj))
		writeMeshkitError(w, ErrParseBool(err, obj), http.StatusForbidden)
		return
	}
	q := req.URL.Query()

	testName := q.Get("name")
	if testName == "" {
		h.log.Error(ErrBlankName(err))
		writeMeshkitError(w, ErrBlankName(err), http.StatusForbidden)
		return
	}
	meshName := q.Get("mesh")
	testUUID := q.Get("uuid")
	// getting profile id from URL

	headersString := q.Get("headers")
	cookiesString := q.Get("cookies")
	contentType := q.Get("contentType")
	bodyString := q.Get("reqBody")

	headers := h.jsonToMap(headersString)
	cookies := h.jsonToMap(cookiesString)
	body = []byte(bodyString)
	h.log.Debug("Headers : ", headers)

	loadTestOptions.Headers = headers
	loadTestOptions.Cookies = cookies
	loadTestOptions.Body = body
	loadTestOptions.ContentType = contentType

	tt, _ := strconv.Atoi(q.Get("t"))
	if tt < 1 {
		tt = 1
	}
	dur := ""
	switch strings.ToLower(q.Get("dur")) {
	case "h":
		dur = "h"
	case "m":
		dur = "m"
	// case "s":
	default:
		dur = "s"
	}
	loadTestOptions.Duration, err = time.ParseDuration(fmt.Sprintf("%d%s", tt, dur))
	if err != nil {
		obj := "load test duration"
		h.log.Error(ErrParseBool(err, obj))
		writeMeshkitError(w, ErrParseBool(err, obj), http.StatusForbidden)
		return
	}

	cc, _ := strconv.Atoi(q.Get("c"))
	if cc < 1 {
		cc = 1
	}
	loadTestOptions.HTTPNumThreads = cc

	loadTestURL := q.Get("url")
	ltURL, err := url.Parse(loadTestURL)
	if err != nil || !ltURL.IsAbs() {
		obj := "the provided load test url"
		h.log.Error(ErrParseBool(err, obj))
		writeMeshkitError(w, ErrParseBool(err, obj), http.StatusBadRequest)
		return
	}
	loadTestOptions.URL = loadTestURL
	loadTestOptions.Name = testName
	loadTestOptions.AllowInitialErrors = true

	qps, _ := strconv.ParseFloat(q.Get("qps"), 64)
	if qps < 0 {
		qps = 0
	}
	loadTestOptions.HTTPQPS = qps

	// fortio is the only supported load generator. Any generator passed by
	// the client (including the removed "wrk2") runs on fortio.
	loadTestOptions.LoadGenerator = models.FortioLG
	h.log.Info("perf test with config: ", loadTestOptions)
	h.loadTestHelperHandler(w, req, profileID, testName, meshName, testUUID, prefObj, loadTestOptions, provider)
}

func (h *Handler) loadTestHelperHandler(w http.ResponseWriter, req *http.Request, profileID, testName, meshName, testUUID string,
	prefObj *models.Preference, loadTestOptions *models.LoadTestOptions, provider models.Provider) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		h.log.Error(ErrEventStreamingNotSupported)
		writeMeshkitError(w, ErrEventStreamingNotSupported, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	notify := req.Context()

	respChan := make(chan *models.LoadTestResponse, 100)
	endChan := make(chan struct{})
	defer close(endChan)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.log.Error(ErrPanicRecovery(r))
			}
		}()
		for data := range respChan {
			bd, err := json.Marshal(data)
			if err != nil {
				// We are inside an active SSE connection — Content-Type:
				// text/event-stream has already been committed above. Calling
				// http.Error here would silently corrupt the stream (it tries
				// to set text/plain after headers are flushed, then writes a
				// trailing newline that breaks SSE framing). Instead, emit an
				// in-protocol error event using the same LoadTestResponse
				// envelope the client already understands (status="error" is
				// handled by the EventSource onmessage consumer in the UI).
				//
				// We deliberately do NOT return here: the worker goroutine
				// will keep writing into respChan until executeLoadTest
				// exits, and an early return would block the worker once
				// the buffered slots fill. Continue draining so the worker
				// can finish and close(respChan) terminates this loop.
				h.log.Error(models.ErrMarshal(err, "meshery result for shipping"))
				errPayload, mErr := json.Marshal(&models.LoadTestResponse{
					Status:  models.LoadTestError,
					Message: models.ErrMarshal(err, "meshery result for shipping").Error(),
				})
				if mErr != nil {
					// Last-resort: a static, hand-rolled JSON literal so the
					// stream stays well-formed even if envelope marshal fails.
					errPayload = []byte(`{"status":"error","message":"failed to marshal load test result"}`)
				}
				_, _ = fmt.Fprintf(w, "data: %s\n\n", errPayload)
				flusher.Flush()
				continue
			}

			h.log.Debug("received new data on response channel")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", bd)
			flusher.Flush()
			h.log.Debug("Flushed the messages on the wire...")
		}
		endChan <- struct{}{}
		h.log.Debug("response channel closed")
	}()
	go func() {
		ctx := context.Background()
		h.executeLoadTest(ctx, req, profileID, testName, meshName, testUUID, prefObj, provider, loadTestOptions, respChan)
		close(respChan)
	}()
	select {
	case <-notify.Done():
		h.log.Debug("received signal to close connection and channels")
		break
	case <-endChan:
		h.log.Debug("load test completed")
		_ = req.Body.Close()
	}
}

func (h *Handler) executeLoadTest(ctx context.Context, req *http.Request, profileID, testName, meshName, testUUID string, prefObj *models.Preference, provider models.Provider, loadTestOptions *models.LoadTestOptions, respChan chan *models.LoadTestResponse) {
	respChan <- &models.LoadTestResponse{
		Status:  models.LoadTestInfo,
		Message: "Initiating load test . . . ",
	}
	var (
		resultsMap map[string]interface{}
		resultInst *periodic.RunnerResults
		err        error
	)
	// fortio is the only supported load generator.
	resultsMap, resultInst, err = helpers.FortioLoadTest(loadTestOptions, h.log)
	if err != nil {
		h.log.Error(ErrLoadTest(err, "unable to perform"))
		respChan <- &models.LoadTestResponse{
			Status:  models.LoadTestError,
			Message: "unable to perform",
		}
		return
	}

	respChan <- &models.LoadTestResponse{
		Status:  models.LoadTestInfo,
		Message: "Load test completed, fetching metadata now",
	}

	resultsMap["load-generator"] = loadTestOptions.LoadGenerator

	mk8sContexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
	if !ok {
		h.log.Warn(ErrInvalidK8SConfigNil)
	}
	var wg sync.WaitGroup

	// If no Kubernetes context is selected, skip cluster discovery
	if len(mk8sContexts) == 0 {
		result := &models.MesheryResult{
			Name:   testName,
			Mesh:   meshName,
			Result: resultsMap,
		}

		h.persistPerformanceTestResult(ctx, req, result, testUUID, profileID, resultInst, prefObj, provider, respChan)
		return
	}

	var resultsMx sync.Mutex
	for _, k8context := range mk8sContexts {
		wg.Add(1)
		go func(mk8scontext models.K8sContext) {
			defer wg.Done()
			var serverVersion string
			var nodes []*models.K8SNode
			var installedMeshes map[string][]corev1.Pod
			// Get the k8sconfig
			k8sconfig, err := mk8scontext.GenerateKubeConfig()
			if err == nil {
				nodesChan := make(chan []*models.K8SNode)
				versionChan := make(chan string)
				installedMeshesChan := make(chan map[string][]corev1.Pod)

				// Fire goroutines to get cluster information
				go func() {
					var nodes []*models.K8SNode
					var err error
					nodes, err = helpers.FetchKubernetesNodes(k8sconfig, mk8scontext.Name, h.log)
					if err != nil {
						err = errors.Wrap(err, "unable to ping kubernetes for context: "+mk8scontext.ID)
						h.log.Warn(ErrFetchKubernetes(err))
					}

					nodesChan <- nodes
				}()
				go func() {
					var serverVersion string
					var err error
					serverVersion, err = helpers.FetchKubernetesVersion(k8sconfig, mk8scontext.Name, h.log)
					if err != nil {
						h.log.Error(ErrFetchKubernetes(err))
					}

					versionChan <- serverVersion
				}()
				go func() {
					installedMeshes, err := helpers.ScanKubernetes(k8sconfig, mk8scontext.Name, h.log)
					if err != nil {
						h.log.Warn(ErrFetchKubernetes(err))
					}

					installedMeshesChan <- installedMeshes
				}()

				// Retrieve cluster information from the goroutines
				serverVersion = <-versionChan
				nodes = <-nodesChan
				installedMeshes = <-installedMeshesChan
			}

			// Making sure different context goroutines don't interfere with each others' results Map
			resultsMx.Lock()
			defer resultsMx.Unlock()

			resultsMap["kubernetes"] = nil
			if serverVersion != "" {
				resultsMap["kubernetes"] = map[string]interface{}{
					"server_version": serverVersion,
					"nodes":          nodes,
				}
			}

			resultsMap["detected-meshes"] = nil
			if len(installedMeshes) > 0 {
				resultsMap["detected-meshes"] = installedMeshes
			}

			respChan <- &models.LoadTestResponse{
				Status:  models.LoadTestInfo,
				Message: "Obtained the needed metadatas, attempting to persist the result for cluster " + mk8scontext.Name,
			}

			result := &models.MesheryResult{
				Name:   testName,
				Mesh:   meshName,
				Result: resultsMap,
			}

			h.persistPerformanceTestResult(ctx, req, result, testUUID, profileID, resultInst, prefObj, provider, respChan)
		}(k8context)
	}
	wg.Wait()
}

// persistPerformanceTestResult takes the test result and saves it on the provider
func (h *Handler) persistPerformanceTestResult(ctx context.Context, req *http.Request, result *models.MesheryResult, testUUID, profileID string, resultInst *periodic.RunnerResults, prefObj *models.Preference, provider models.Provider, respChan chan *models.LoadTestResponse) {
	resultID, err := provider.PublishResults(req, result, profileID)
	if err != nil {
		h.log.Error(ErrLoadTest(err, "unable to persist in cluster"))
		respChan <- &models.LoadTestResponse{
			Status:  models.LoadTestError,
			Message: "unable to persist",
		}
		return
	}
	respChan <- &models.LoadTestResponse{
		Status:  models.LoadTestInfo,
		Message: "Done persisting the load test results.",
	}

	key := uuid.FromStringOrNil(resultID)
	if key == uuid.Nil {
		key = uuid.Must(uuid.NewV4())
	}
	result.ID = key
	respChan <- &models.LoadTestResponse{
		Status: models.LoadTestSuccess,
		Result: result,
	}
}

func assignCertificatePath(key, path string, loadTestOptions *models.LoadTestOptions) {
	loadTestOptions.CACert = path
}
