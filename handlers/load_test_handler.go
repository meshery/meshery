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
	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/apps/v1"
)

// LoadTestUsingSMPSHandler runs the load test with the given parameters and SMPS
func (h *Handler) LoadTestUsingSMPSHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method != http.MethodPost && req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer func() {
		_ = req.Body.Close()
	}()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		msg := "unable to read request body"
		err = errors.Wrapf(err, msg)
		logrus.Error(err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	benchMark := &models.BenchmarkSpec{}
	if err := yaml.Unmarshal(body, benchMark); err != nil {
		msg := "unable to parse the provided input"
		err = errors.Wrapf(err, msg)
		logrus.Error(err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	q := req.URL.Query()

	// testName - should be loaded from the file and updated with a random string appended to the end of the name
	testName := q.Get("name")
	if testName == "" {
		logrus.Errorf("Error: name field is blank")
		http.Error(w, "Provide a name for the test.", http.StatusForbidden)
		return
	}
	meshName := q.Get("mesh")
	testUUID := uuid.Must(uuid.NewV4()).String()

	loadTestOptions := &models.LoadTestOptions{}

	loadTestOptions.Duration = benchMark.EndTime.Sub(benchMark.StartTime)

	if loadTestOptions.Duration.Seconds() <= 0 {
		loadTestOptions.Duration = time.Second
	}

	loadTestOptions.IsGRPC = false

	loadTestOptions.URL = benchMark.EndpointURL
	if benchMark.Client != nil {
		loadTestOptions.HTTPNumThreads = benchMark.Client.Connections
		loadTestOptions.HTTPQPS = benchMark.Client.Rps

	}

	if loadTestOptions.HTTPNumThreads < 1 {
		loadTestOptions.HTTPNumThreads = 1
	}

	ltURL, err := url.Parse(loadTestOptions.URL)
	if err != nil || !ltURL.IsAbs() {
		logrus.Errorf("unable to parse the provided load test url: %v", err)
		http.Error(w, "invalid load test URL", http.StatusBadRequest)
		return
	}
	loadTestOptions.Name = testName

	if loadTestOptions.HTTPQPS < 0 {
		loadTestOptions.HTTPQPS = 0
	}

	// loadGenerator := q.Get("loadGenerator")

	// switch loadGenerator {
	// case "wrk2":
	// 	loadTestOptions.LoadGenerator = models.Wrk2LG
	// default:
	loadTestOptions.LoadGenerator = models.FortioLG
	// }

	h.loadTestHelperHandler(w, req, testName, meshName, testUUID, prefObj, loadTestOptions, provider)
}

// LoadTestHandler runs the load test with the given parameters
func (h *Handler) LoadTestHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method != http.MethodPost && req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := req.ParseForm()
	if err != nil {
		logrus.Errorf("Error: unable to parse form: %v", err)
		http.Error(w, "unable to process the received data", http.StatusForbidden)
		return
	}
	q := req.URL.Query()

	testName := q.Get("name")
	if testName == "" {
		logrus.Errorf("Error: name field is blank")
		http.Error(w, "Provide a name for the test.", http.StatusForbidden)
		return
	}
	meshName := q.Get("mesh")
	testUUID := q.Get("uuid")

	loadTestOptions := &models.LoadTestOptions{}

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
		logrus.Errorf("Error: unable to parse load test duration: %v", err)
		http.Error(w, "unable to process the received data", http.StatusForbidden)
		return
	}

	loadTestOptions.IsGRPC = false

	cc, _ := strconv.Atoi(q.Get("c"))
	if cc < 1 {
		cc = 1
	}
	loadTestOptions.HTTPNumThreads = cc

	loadTestURL := q.Get("url")
	ltURL, err := url.Parse(loadTestURL)
	if err != nil || !ltURL.IsAbs() {
		logrus.Errorf("unable to parse the provided load test url: %v", err)
		http.Error(w, "invalid load test URL", http.StatusBadRequest)
		return
	}
	loadTestOptions.URL = loadTestURL
	loadTestOptions.Name = testName

	qps, _ := strconv.ParseFloat(q.Get("qps"), 64)
	if qps < 0 {
		qps = 0
	}
	loadTestOptions.HTTPQPS = qps

	loadGenerator := q.Get("loadGenerator")

	switch loadGenerator {
	case models.Wrk2LG.Name():
		loadTestOptions.LoadGenerator = models.Wrk2LG
	default:
		loadTestOptions.LoadGenerator = models.FortioLG
	}

	// q.Set("json", "on")

	// client := &http.Client{}
	// fortioURL, err := url.Parse(h.config.FortioURL)
	// if err != nil {
	// 	logrus.Errorf("unable to parse the provided fortio url: %v", err)
	// 	http.Error(w, "error while running load test", http.StatusInternalServerError)
	// 	return
	// }
	// fortioURL.RawQuery = q.Encode()
	// logrus.Infof("load test constructed url: %s", fortioURL.String())
	// fortioResp, err := client.Get(fortioURL.String())
	h.loadTestHelperHandler(w, req, testName, meshName, testUUID, prefObj, loadTestOptions, provider)
}

func (h *Handler) loadTestHelperHandler(w http.ResponseWriter, req *http.Request, testName, meshName, testUUID string,
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

	notify := w.(http.CloseNotifier).CloseNotify()
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
				logrus.Errorf("error: unable to marshal meshery result for shipping: %v", err)
				http.Error(w, "error while running load test", http.StatusInternalServerError)
				return
			}

			log.Debug("received new data on response channel")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", bd)
			if flusher != nil {
				flusher.Flush()
				log.Debugf("Flushed the messages on the wire...")
			}
		}
		endChan <- struct{}{}
		log.Debug("response channel closed")
	}()
	go func() {
		h.executeLoadTest(req, testName, meshName, testUUID, prefObj, provider, loadTestOptions, respChan)
		close(respChan)
	}()
	select {
	case <-notify:
		log.Debugf("received signal to close connection and channels")
		break
	case <-endChan:
		log.Debugf("load test completed")
	}
}

func (h *Handler) executeLoadTest(req *http.Request, testName, meshName, testUUID string, prefObj *models.Preference, provider models.Provider, loadTestOptions *models.LoadTestOptions, respChan chan *models.LoadTestResponse) {
	respChan <- &models.LoadTestResponse{
		Status:  models.LoadTestInfo,
		Message: "Initiating load test . . . ",
	}
	// resultsMap, resultInst, err := helpers.FortioLoadTest(loadTestOptions)
	var (
		resultsMap map[string]interface{}
		resultInst *periodic.RunnerResults
		err        error
	)
	if loadTestOptions.LoadGenerator == models.Wrk2LG {
		resultsMap, resultInst, err = helpers.WRK2LoadTest(loadTestOptions)
	} else {
		resultsMap, resultInst, err = helpers.FortioLoadTest(loadTestOptions)
	}
	if err != nil {
		msg := "error: unable to perform load test"
		err = errors.Wrap(err, msg)
		logrus.Error(err)
		respChan <- &models.LoadTestResponse{
			Status:  models.LoadTestError,
			Message: msg,
		}
		return
	}

	respChan <- &models.LoadTestResponse{
		Status:  models.LoadTestInfo,
		Message: "Load test completed, fetching metadata now",
	}

	if prefObj.K8SConfig != nil {
		nodesChan := make(chan []*models.K8SNode)
		versionChan := make(chan string)
		installedMeshesChan := make(chan map[string][]v1.Deployment)

		go func() {
			var nodes []*models.K8SNode
			var err error
			if len(prefObj.K8SConfig.Nodes) == 0 {
				nodes, err = helpers.FetchKubernetesNodes(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
				if err != nil {
					err = errors.Wrap(err, "unable to ping kubernetes")
					// logrus.Error(err)
					logrus.Warn(err)
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
					err = errors.Wrap(err, "unable to ping kubernetes")
					// logrus.Error(err)
					logrus.Warn(err)
					// return
				}
			}
			versionChan <- serverVersion
		}()
		go func() {
			installedMeshes, err := helpers.ScanKubernetes(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
			if err != nil {
				err = errors.Wrap(err, "unable to scan kubernetes")
				logrus.Warn(err)
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
	// // defer fortioResp.Body.Close()
	// // bd, err := ioutil.ReadAll(fortioResp.Body)
	// bd, err := json.Marshal(resp)
	// if err != nil {
	// 	logrus.Errorf("Error: unable to parse response from fortio: %v", err)
	// 	http.Error(w, "error while running load test", http.StatusInternalServerError)
	// 	return
	// }

	result := &models.MesheryResult{
		Name:   testName,
		Mesh:   meshName,
		Result: resultsMap,
	}

	resultID, err := provider.PublishResults(req, result)
	if err != nil {
		// http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		// return
		msg := "error: unable to persist the load test results"
		err = errors.Wrap(err, msg)
		logrus.Error(err)
		// http.Error(w, "error while running load test", http.StatusInternalServerError)
		respChan <- &models.LoadTestResponse{
			Status:  models.LoadTestError,
			Message: msg,
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

	logrus.Debugf("promURL: %s, testUUID: %s, resultID: %s", promURL, testUUID, resultID)
	if promURL != "" && testUUID != "" && resultID != "" &&
		(provider.GetProviderType() == models.RemoteProviderType ||
			(provider.GetProviderType() == models.LocalProviderType && prefObj.AnonymousPerfResults)) {
		_ = h.task.Call(&models.SubmitMetricsConfig{
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
	// w.Write(bd)
	respChan <- &models.LoadTestResponse{
		Status: models.LoadTestSuccess,
		Result: result,
	}
}

// CollectStaticMetrics is used for collecting static metrics from prometheus and submitting it to SaaS
func (h *Handler) CollectStaticMetrics(config *models.SubmitMetricsConfig) error {
	logrus.Debugf("initiating collecting prometheus static board metrics for test id: %s", config.TestUUID)
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
			// sd, _ := json.Marshal(seriesData)
			// sd, _ := json.Marshal(queryResponse)
			// logrus.Debugf("Retrieved series data: %s", sd)
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
		logrus.Error(errors.Wrap(err, "error parsing result uuid"))
		return err
	}
	result := &models.MesheryResult{
		ID:                resultUUID,
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
