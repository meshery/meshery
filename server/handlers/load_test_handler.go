// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"fortio.org/fortio/periodic"
	yaml "github.com/ghodss/yaml"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/models"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// LoadTestUsingSMPHandler runs the load test with the given parameters and SMP
func (h *Handler) LoadTestUsingSMPHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// Read the SMP File
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to read request body: %s", err)
		return
	}

	if req.Header.Get("Content-Type") == "application/json" {
		body, err = yaml.JSONToYAML(body)
		if err != nil {
			h.log.Error(ErrPatternFile(err))
			http.Error(w, ErrPatternFile(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	jsonBytes, _ := yaml.YAMLToJSON(body)

	perfTest := &models.PerformanceTestConfigFile{}
	if err := json.Unmarshal(jsonBytes, perfTest); err != nil {
		h.log.Error(ErrParseBool(err, "provided input"))
		http.Error(w, ErrParseBool(err, "provided input").Error(), http.StatusBadRequest)
		return
	}

	// testName - should be loaded from the file and updated with a random string appended to the end of the name
	testName := perfTest.Config.Name
	if testName == "" {
		h.log.Error(ErrBlankName(err))
		http.Error(w, ErrBlankName(err).Error(), http.StatusForbidden)
		return
	}

	meshType := perfTest.ServiceMesh.Type
	meshName := SMP.ServiceMesh_Type_name[int32(meshType)]

	profileID := perfTest.Config.Id

	loadTestOptions := &models.LoadTestOptions{}

	testDuration, err := time.ParseDuration(perfTest.Config.Duration)
	if err != nil {
		h.log.Error(ErrParseDuration)
		http.Error(w, ErrParseDuration.Error(), http.StatusBadRequest)
		return
	}
	loadTestOptions.Duration = testDuration
	if loadTestOptions.Duration.Seconds() <= 0 {
		loadTestOptions.Duration = time.Second
	}

	// TODO: check multiple clients in case of distributed perf test
	testClient := perfTest.Config.Clients[0]

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
		http.Error(w, ErrParseBool(err, obj).Error(), http.StatusBadRequest)
		return
	}
	if !ltURL.IsAbs() {
		h.log.Error(ErrInvalidLTURL(ltURL.String()))
		http.Error(w, ErrInvalidLTURL(ltURL.String()).Error(), http.StatusBadRequest)
		return
	}
	loadTestOptions.Name = testName

	if loadTestOptions.HTTPQPS < 0 {
		loadTestOptions.HTTPQPS = 0
	}

	loadGenerator := testClient.LoadGenerator

	switch loadGenerator {
	case models.Wrk2LG.Name():
		loadTestOptions.LoadGenerator = models.Wrk2LG
	case models.NighthawkLG.Name():
		loadTestOptions.LoadGenerator = models.NighthawkLG
	default:
		loadTestOptions.LoadGenerator = models.FortioLG
	}
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

// swagger:route GET /api/perf/profile PerfAPI idRunPerfTest
// Handle GET request to run a test
//
// Runs the load test with the given parameters
// responses:
// 	200:

// swagger:route GET /api/user/performance/profiles/{id}/run PerformanceAPI idRunPerformanceTest
// Handle GET request to run a performance test
//
// Runs the load test with the given parameters
// responses:
// 	200:

// LoadTestHandler runs the load test with the given parameters
func (h *Handler) LoadTestHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		msg := "unable to read request body"
		err = errors.Wrapf(err, msg)
		logrus.Error(err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// if values have been passed as body we run test using SMP Handler
	if string(body) != "" {
		logrus.Info("Running test with SMP config")
		req.Body = io.NopCloser(strings.NewReader(string(body)))
		h.LoadTestUsingSMPHandler(w, req, prefObj, user, provider)
		return
	}

	err = req.ParseForm()
	if err != nil {
		obj := "form"
		h.log.Error(ErrParseBool(err, obj))
		http.Error(w, ErrParseBool(err, obj).Error(), http.StatusForbidden)
		return
	}
	q := req.URL.Query()

	testName := q.Get("name")
	if testName == "" {
		h.log.Error(ErrBlankName(err))
		http.Error(w, ErrBlankName(err).Error(), http.StatusForbidden)
		return
	}
	meshName := q.Get("mesh")
	testUUID := q.Get("uuid")
	// getting profile id from URL
	profileID := mux.Vars(req)["id"]

	headersString := q.Get("headers")
	cookiesString := q.Get("cookies")
	contentType := q.Get("contentType")
	bodyString := q.Get("reqBody")

	headers := h.jsonToMap(headersString)
	cookies := h.jsonToMap(cookiesString)
	body = []byte(bodyString)
	h.log.Debug("Headers : ", headers)

	loadTestOptions := &models.LoadTestOptions{}
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
		http.Error(w, ErrParseBool(err, obj).Error(), http.StatusForbidden)
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
		http.Error(w, ErrParseBool(err, obj).Error(), http.StatusBadRequest)
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

	loadGenerator := q.Get("loadGenerator")

	switch loadGenerator {
	case models.Wrk2LG.Name():
		loadTestOptions.LoadGenerator = models.Wrk2LG
	case models.NighthawkLG.Name():
		loadTestOptions.LoadGenerator = models.NighthawkLG
	default:
		loadTestOptions.LoadGenerator = models.FortioLG
	}
	h.log.Info("perf test with config: ", loadTestOptions)
	h.loadTestHelperHandler(w, req, profileID, testName, meshName, testUUID, prefObj, loadTestOptions, provider)
}

func (h *Handler) loadTestHelperHandler(w http.ResponseWriter, req *http.Request, profileID, testName, meshName, testUUID string,
	prefObj *models.Preference, loadTestOptions *models.LoadTestOptions, provider models.Provider) {
	log := logrus.WithField("file", "load_test_handler")

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Error("Event streaming not supported.")
		http.Error(w, "Event streaming is not supported at the moment.", http.StatusInternalServerError)
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
				h.log.Error(ErrMarshal(err, "meshery result for shipping"))
				http.Error(w, ErrMarshal(err, "meshery result for shipping").Error(), http.StatusInternalServerError)
				return
			}

			h.log.Debug("received new data on response channel")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", bd)
			if flusher != nil {
				flusher.Flush()
				h.log.Debug("Flushed the messages on the wire...")
			}
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
	if loadTestOptions.LoadGenerator == models.Wrk2LG {
		resultsMap, resultInst, err = helpers.WRK2LoadTest(loadTestOptions)
	} else if loadTestOptions.LoadGenerator == models.NighthawkLG {
		resultsMap, resultInst, err = helpers.NighthawkLoadTest(loadTestOptions)
	} else {
		resultsMap, resultInst, err = helpers.FortioLoadTest(loadTestOptions)
	}
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
	if !ok || len(mk8sContexts) == 0 {
		h.log.Error(ErrInvalidK8SConfig)
		return
	}
	var wg sync.WaitGroup
	for _, k8context := range mk8sContexts {
		wg.Add(1)
		go func(mk8scontext *models.K8sContext) {
			defer wg.Done()
			// Get the k8sconfig
			k8sconfig, err := mk8scontext.GenerateKubeConfig()
			if err == nil {
				nodesChan := make(chan []*models.K8SNode)
				versionChan := make(chan string)
				installedMeshesChan := make(chan map[string][]corev1.Pod)

				go func() {
					var nodes []*models.K8SNode
					var err error
					nodes, err = helpers.FetchKubernetesNodes(k8sconfig, mk8scontext.Name)
					if err != nil {
						err = errors.Wrap(err, "unable to ping kubernetes for context: "+mk8scontext.ID)
						h.log.Warn(ErrFetchKubernetes(err))
					}

					nodesChan <- nodes
				}()
				go func() {
					var serverVersion string
					var err error
					serverVersion, err = helpers.FetchKubernetesVersion(k8sconfig, mk8scontext.Name)
					if err != nil {
						h.log.Error(ErrFetchKubernetes(err))
					}

					versionChan <- serverVersion
				}()
				go func() {
					installedMeshes, err := helpers.ScanKubernetes(k8sconfig, mk8scontext.Name)
					if err != nil {
						h.log.Warn(ErrFetchKubernetes(err))
					}
					installedMeshesChan <- installedMeshes
				}()

				serverVersion := <-versionChan
				nodes := <-nodesChan

				resultsMap["kubernetes"] = map[string]interface{}{
					"server_version": serverVersion,
					"nodes":          nodes,
				}

				installedMeshes := <-installedMeshesChan
				if len(installedMeshes) > 0 {
					resultsMap["detected-meshes"] = installedMeshes
				}
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

			var promURL string
			if prefObj.Prometheus != nil {
				promURL = prefObj.Prometheus.PrometheusURL
			}

			tokenVal, _ := provider.GetProviderToken(req)

			h.log.Debug("promURL: , testUUID: , resultID: ", promURL, testUUID, resultID)
			if promURL != "" && testUUID != "" && resultID != "" &&
				(provider.GetProviderType() == models.RemoteProviderType ||
					(provider.GetProviderType() == models.LocalProviderType && prefObj.AnonymousPerfResults)) {
				_ = h.task.WithArgs(ctx, &models.SubmitMetricsConfig{
					TestUUID:  testUUID,
					ResultID:  resultID,
					PromURL:   promURL,
					StartTime: resultInst.StartTime,
					EndTime:   resultInst.StartTime.Add(resultInst.ActualDuration),
					TokenVal:  tokenVal,
					Provider:  provider,
				})
			}

			key := uuid.FromStringOrNil(resultID)
			if key == uuid.Nil {
				key, _ = uuid.NewV4()
			}
			result.ID = key
			respChan <- &models.LoadTestResponse{
				Status: models.LoadTestSuccess,
				Result: result,
			}

			if h.config.PerformanceChannel != nil {
				h.config.PerformanceChannel <- struct{}{}
			}

			if h.config.PerformanceResultChannel != nil {
				h.config.PerformanceResultChannel <- struct{}{}
			}
		}(&k8context)
	}
	wg.Wait()
}

// CollectStaticMetrics is used for collecting static metrics from prometheus and submitting it to Remote Provider
func (h *Handler) CollectStaticMetrics(config *models.SubmitMetricsConfig) error {
	h.log.Debug("initiating collecting prometheus static board metrics for test id: ", config.TestUUID)
	ctx := context.Background()
	queries := h.config.QueryTracker.GetQueriesForUUID(ctx, config.TestUUID)
	queryResults := map[string]map[string]interface{}{}
	step := h.config.PrometheusClient.ComputeStep(ctx, config.StartTime, config.EndTime)
	for query, flag := range queries {
		if !flag {
			seriesData, err := h.config.PrometheusClient.QueryRangeUsingClient(ctx, config.PromURL, query, config.StartTime, config.EndTime, step)
			if err != nil {
				return err
			}
			queryResults[query] = map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"resultType": seriesData.Type(),
					"result":     seriesData,
				},
			}
			h.config.QueryTracker.AddOrFlagQuery(ctx, config.TestUUID, query, true)
		}
	}

	board, err := h.config.PrometheusClient.GetClusterStaticBoard(ctx, config.PromURL)
	if err != nil {
		return err
	}
	// TODO: we are NOT persisting the Node metrics for now

	resultUUID, err := uuid.FromString(config.ResultID)
	if err != nil {
		h.log.Error(ErrParseBool(err, "result uuid"))
		return err
	}
	result := &models.MesheryResult{
		ID:                resultUUID,
		TestID:            config.TestUUID,
		ServerMetrics:     queryResults,
		ServerBoardConfig: board,
	}

	if err = config.Provider.PublishMetrics(config.TokenVal, result); err != nil {
		return err
	}
	// now to remove all the queries for the uuid
	h.config.QueryTracker.RemoveUUID(ctx, config.TestUUID)
	return nil
}
