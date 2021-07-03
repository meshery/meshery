//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fortio.org/fortio/periodic"
	yamlj "github.com/ghodss/yaml"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	corev1 "k8s.io/api/core/v1"
)

// LoadTestUsingSMPHandler runs the load test with the given parameters and SMP
func (h *Handler) LoadTestUsingSMPHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	defer func() {
		_ = req.Body.Close()
	}()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err) 

        http.Error(w, ErrRequestBody(err),Error(), http.StatusInternalServerError)
        return
	}
	jsonBody, err := yamlj.YAMLToJSON(body)
	if err != nil {
		h.log.Error(ErrConversion(err))
        http.Error(w, ErrConversion(err).Error(), http.StatusInternalServerError)
        return
	}
	perfTest := &SMP.PerformanceTestConfig{}
	if err := protojson.Unmarshal(jsonBody, perfTest); err != nil {
		obj := "provided input"
        h.log.Error(ErrParseBool(err, obj))
        http.Error(w, ErrParseBool(err, obj).Error(), http.StatusBadRequest)
        return
	}

	// testName - should be loaded from the file and updated with a random string appended to the end of the name
	testName := perfTest.Name
	if testName == "" {
		h.log.Error(ErrBlankName(err))
        http.Error(w, ErrBlankName(err).Error(), http.StatusForbidden)
        return
	}
	// meshName := q.Get("mesh")
	testUUID := perfTest.Id

	loadTestOptions := &models.LoadTestOptions{}

	testDuration, err := time.ParseDuration(perfTest.Duration)
	if err != nil {
		msg := "error parsing test duration, please refer to: https://docs.meshery.io/guides/mesheryctl#performance-management"
		err = errors.Wrapf(err, msg)
		h.log.Error(ErrParseDuration(err))
        http.Error(w, ErrParseDuration(err).Error(), http.StatusBadRequest)
        return
	}
	loadTestOptions.Duration = testDuration
	if loadTestOptions.Duration.Seconds() <= 0 {
		loadTestOptions.Duration = time.Second
	}

	// TODO: check multiple clients in case of distributed perf test
	testClient := perfTest.Clients[0]

	// TODO: consider the multiple endpoints
	loadTestOptions.URL = testClient.EndpointUrls[0]
	loadTestOptions.HTTPNumThreads = int(testClient.Connections)
	loadTestOptions.HTTPQPS = float64(testClient.Rps)

	if loadTestOptions.HTTPNumThreads < 1 {
		loadTestOptions.HTTPNumThreads = 1
	}

	ltURL, err := url.Parse(loadTestOptions.URL)
	if err != nil || !ltURL.IsAbs() {
		h.log.Error(ErrParseBool(err,obj))
        http.Error(w, ErrParseBool(err,obj).Error(), http.StatusBadRequest)
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

	h.loadTestHelperHandler(w, req, testName, "", testUUID, prefObj, loadTestOptions, provider)
}

func (h *Handler) jsonToMap(headersString string) *map[string]string {
	headers := make(map[string]string)
	err := json.Unmarshal([]byte(headersString), &headers)
	if err != nil {
		return nil
	}
	return &headers
}

// swagger:route GET /api/user/performance/profiles/{id}/run PerformanceAPI idRunPerformanceTest
// Handle GET request to run a performance test
//
// Runs the load test with the given parameters
// responses:
// 	200:

// LoadTestHandler runs the load test with the given parameters
func (h *Handler) LoadTestHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	err := req.ParseForm()
	if err != nil {
		obj := "form"
        h.log.Errorf(ErrParseBool(err,obj))
        http.Error(w, ErrParseBool(err,obj),Error(), http.StatusForbidden)
        return

	}
	q := req.URL.Query()

	testName := q.Get("name")
	if testName == "" {
		h.log.Error(ErrBlankName(err))
        http.Error(w, ErrBlankName(err),Error(), http.StatusForbidden)
        return
	}
	meshName := q.Get("mesh")
	testUUID := q.Get("uuid")

	headersString := q.Get("headers")
	cookiesString := q.Get("cookies")
	contentType := q.Get("contentType")
	bodyString := q.Get("reqBody")

	headers := h.jsonToMap(headersString)
	cookies := h.jsonToMap(cookiesString)
	body := []byte(bodyString)
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
        h.log.Error(ErrParseBool(err,obj))
        http.Error(w, ErrParseBool(err,obj).Error(), http.StatusBadRequest)
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
	h.log.Infof("perf test with config: %v", loadTestOptions)
	h.loadTestHelperHandler(w, req, testName, meshName, testUUID, prefObj, loadTestOptions, provider)
}

func (h *Handler) loadTestHelperHandler(w http.ResponseWriter, req *http.Request, testName, meshName, testUUID string,
	prefObj *models.Preference, loadTestOptions *models.LoadTestOptions, provider models.Provider) {
	log := h.log.WithField("file", "load_test_handler")

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
				log.Errorf("Recovered from panic: %v.", r)
			}
		}()
		for data := range respChan {
			bd, err := json.Marshal(data)
			if err != nil {
				obj:= "meshery result for shipping"
                h.log.Error(ErrMarshal(err,obj))
                http.Error(w, ErrMarshal(err,obj).Error(), http.StatusInternalServerError)
                return
			}

			log.Debug("received new data on response channel")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", bd)
			if flusher != nil {
				flusher.Flush()
				log.Debug("Flushed the messages on the wire...")
			}
		}
		endChan <- struct{}{}
		log.Debug("response channel closed")
	}()
	go func() {
		ctx := context.Background()
		h.executeLoadTest(ctx, req, testName, meshName, testUUID, prefObj, provider, loadTestOptions, respChan)
		close(respChan)
	}()
	select {
	case <-notify.Done():
		log.Debug("received signal to close connection and channels")
		break
	case <-endChan:
		log.Debug("load test completed")
	}
}

func (h *Handler) executeLoadTest(ctx context.Context, req *http.Request, testName, meshName, testUUID string, prefObj *models.Preference, provider models.Provider, loadTestOptions *models.LoadTestOptions, respChan chan *models.LoadTestResponse) {
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
		obj :=  "unable to perform"
	
		h.log.Error(ErrLoadTest(err,obj))
        respChan <- &models.LoadTestResponse{
            Status:  models.LoadTestError,
			Message: obj
          
		}
		return
	}

	respChan <- &models.LoadTestResponse{
		Status:  models.LoadTestInfo,
		Message: "Load test completed, fetching metadata now",
	}

	resultsMap["load-generator"] = loadTestOptions.LoadGenerator

	if prefObj.K8SConfig != nil {
		nodesChan := make(chan []*models.K8SNode)
		versionChan := make(chan string)
		installedMeshesChan := make(chan map[string][]corev1.Pod)

		go func() {
			var nodes []*models.K8SNode
			var err error
			if len(prefObj.K8SConfig.Nodes) == 0 {
				nodes, err = helpers.FetchKubernetesNodes(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
				if err != nil {
		
					err = errors.Wrap(err, "unable to ping kubernetes")
					// logrus.Error(err)
					h.log.Warn(ErrFetchKubernetes(err)
					// return
				}
			}
			nodesChan <- nodes
		}()
		go func() {
			var serverVersion string
			var err error
			if prefObj.K8SConfig.ServerVersion == "" {
				serverVersion, err = helpers.FetchKubernetesVersion(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
				if err != nil {
                    
                    h.log.Error(ErrFetchKubernetes(err,obj)
                }

			}
			versionChan <- serverVersion
		}()
		go func() {
			installedMeshes, err := helpers.ScanKubernetes(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
			if err != nil {
			
                
                h.log.Warn(ErrFetchKubernetes(err,obj)

			}
			installedMeshesChan <- installedMeshes
		}()

		prefObj.K8SConfig.Nodes = <-nodesChan
		prefObj.K8SConfig.ServerVersion = <-versionChan

		if prefObj.K8SConfig.ServerVersion != "" && len(prefObj.K8SConfig.Nodes) > 0 {
			resultsMap["kubernetes"] = map[string]interface{}{
				"server_version": prefObj.K8SConfig.ServerVersion,
				"nodes":          prefObj.K8SConfig.Nodes,
			}
		}
		installedMeshes := <-installedMeshesChan
		if len(installedMeshes) > 0 {
			resultsMap["detected-meshes"] = installedMeshes
		}
	}
	respChan <- &models.LoadTestResponse{
		Status:  models.LoadTestInfo,
		Message: "Obtained the needed metadatas, attempting to persist the result",
	}

	result := &models.MesheryResult{
		Name:   testName,
		Mesh:   meshName,
		Result: resultsMap,
	}

	resultID, err := provider.PublishResults(req, result, mux.Vars(req)["id"])
	if err != nil {
		obj:="unable to persist"
        
        h.log.Error(ErrLoadTest(err,obj))

		respChan <- &models.LoadTestResponse{
			Status:  models.LoadTestError,
			ErrLoadTest(err,obj).Error()

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
		obj := "result uuid"
        h.log.Error(ErrParseBool(err, obj))

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
