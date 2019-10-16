package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LoadTestHandler runs the load test with the given parameters
func (h *Handler) LoadTestHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodPost && req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tokenVal, _ := session.Values[h.config.SaaSTokenName].(string)
	err := req.ParseForm()
	if err != nil {
		logrus.Errorf("Error: unable to parse form: %v", err)
		http.Error(w, "unable to process the received data", http.StatusForbidden)
		return
	}
	q := req.Form

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

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("Unable to read session from the session persister. Starting a new session.")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	resultsMap, resultInst, err := helpers.FortioLoadTest(loadTestOptions)
	if err != nil {
		logrus.Errorf("error: unable to perform load test: %v", err)
		http.Error(w, "error while running load test", http.StatusInternalServerError)
		return
	}

	if sessObj.K8SConfig != nil {
		nodesChan := make(chan []*models.K8SNode)
		versionChan := make(chan string)
		installedMeshesChan := make(chan map[string]string)

		go func() {
			var nodes []*models.K8SNode
			var err error
			if len(sessObj.K8SConfig.Nodes) == 0 {
				nodes, err = helpers.FetchKubernetesNodes(sessObj.K8SConfig.Config, sessObj.K8SConfig.ContextName)
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
			if sessObj.K8SConfig.ServerVersion == "" {
				serverVersion, err = helpers.FetchKubernetesVersion(sessObj.K8SConfig.Config, sessObj.K8SConfig.ContextName)
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
			installedMeshes, err := helpers.ScanKubernetes(sessObj.K8SConfig.Config, sessObj.K8SConfig.ContextName)
			if err != nil {
				err = errors.Wrap(err, "unable to scan kubernetes")
				logrus.Warn(err)
			}
			installedMeshesChan <- installedMeshes
		}()

		nodes := <-nodesChan
		if len(nodes) > 0 {
			sessObj.K8SConfig.Nodes = nodes
		}
		serverVersion := <-versionChan
		if serverVersion != "" {
			sessObj.K8SConfig.ServerVersion = serverVersion
		}
		if sessObj.K8SConfig.ServerVersion != "" && len(sessObj.K8SConfig.Nodes) > 0 {
			resultsMap["kubernetes"] = map[string]interface{}{
				"server_version": sessObj.K8SConfig.ServerVersion,
				"nodes":          sessObj.K8SConfig.Nodes,
			}
		}
		installedMeshes := <-installedMeshesChan
		if len(installedMeshes) > 0 {
			resultsMap["detected-meshes"] = installedMeshes
		}
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
	bd, err := json.Marshal(result)
	if err != nil {
		logrus.Errorf("error: unable to marshal meshery result for shipping: %v", err)
		http.Error(w, "error while running load test", http.StatusInternalServerError)
		return
	}

	resultID, err := h.publishResultsToSaaS(h.config.SaaSTokenName, tokenVal, bd)
	if err != nil {
		// logrus.Errorf("Error: unable to parse response from fortio: %v", err)
		http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		return
	}

	var promURL string
	if sessObj.Prometheus != nil {
		promURL = sessObj.Prometheus.PrometheusURL
	}

	logrus.Debugf("promURL: %s, testUUID: %s, resultID: %s", promURL, testUUID, resultID)
	if promURL != "" && testUUID != "" && resultID != "" {
		h.task.Call(&models.SubmitMetricsConfig{
			TestUUID:  testUUID,
			ResultID:  resultID,
			PromURL:   promURL,
			StartTime: resultInst.StartTime,
			EndTime:   resultInst.StartTime.Add(resultInst.ActualDuration),
			TokenKey:  h.config.SaaSTokenName,
			TokenVal:  tokenVal,
		})
	}

	w.Write(bd)
}

// CollectStaticMetrics is used for collecting static metrics from prometheus and submitting it to SaaS
func (h *Handler) CollectStaticMetrics(config *models.SubmitMetricsConfig) error {
	logrus.Debugf("initiating collecting prometheus static board metrics for test id: %s", config.TestUUID)
	ctx := context.Background()
	queries := h.config.QueryTracker.GetQueriesForUUID(ctx, config.TestUUID)
	promClient, err := helpers.NewPrometheusClient(ctx, config.PromURL, false) // probably don't need to validate here
	if err != nil {
		return err
	}
	queryResults := map[string]map[string]interface{}{}
	step := promClient.ComputeStep(ctx, config.StartTime, config.EndTime)
	for query, flag := range queries {
		if !flag {
			seriesData, err := promClient.QueryRangeUsingClient(ctx, query, config.StartTime, config.EndTime, step)
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

	prometheusClient, err := helpers.NewPrometheusClient(ctx, config.PromURL, false)
	if err != nil {
		return err
	}

	board, err := prometheusClient.GetClusterStaticBoard(ctx)
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
	sd, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return err
	}

	logrus.Debugf("Result: %s, size: %d", sd, len(sd))

	if err = h.publishMetricsToSaaS(config.TokenKey, config.TokenVal, sd); err != nil {
		return err
	}
	// now to remove all the queries for the uuid
	h.config.QueryTracker.RemoveUUID(ctx, config.TestUUID)
	return nil
}

func (h *Handler) publishMetricsToSaaS(tokenKey, tokenVal string, bd []byte) error {
	logrus.Infof("attempting to publish metrics to SaaS")
	bf := bytes.NewBuffer(bd)
	saasURL, _ := url.Parse(h.config.SaaSBaseURL + "/result/metrics")
	req, _ := http.NewRequest(http.MethodPut, saasURL.String(), bf)
	req.AddCookie(&http.Cookie{
		Name:     tokenKey,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		logrus.Errorf("unable to send metrics: %v", err)
		return err
	}
	if resp.StatusCode == http.StatusOK {
		logrus.Infof("metrics successfully published to SaaS")
		return nil
	}
	defer resp.Body.Close()
	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return err
	}
	logrus.Errorf("error while sending metrics: %s", bdr)
	return fmt.Errorf("error while sending metrics - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

func (h *Handler) publishResultsToSaaS(tokenKey, tokenVal string, bd []byte) (string, error) {
	logrus.Infof("attempting to publish results to SaaS")
	bf := bytes.NewBuffer(bd)
	saasURL, _ := url.Parse(h.config.SaaSBaseURL + "/result")
	req, _ := http.NewRequest(http.MethodPost, saasURL.String(), bf)
	req.AddCookie(&http.Cookie{
		Name:     tokenKey,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		logrus.Errorf("unable to send results: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return "", err
	}
	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("results successfully published to SaaS")
		idMap := map[string]string{}
		if err = json.Unmarshal(bdr, &idMap); err != nil {
			logrus.Errorf("unable to unmarshal body: %v", err)
			return "", err
		}
		resultID, ok := idMap["id"]
		if ok {
			return resultID, nil
		}
		return "", nil
	}
	logrus.Errorf("error while sending results: %s", bdr)
	return "", fmt.Errorf("error while sending results - Status code: %d, Body: %s", resp.StatusCode, bdr)
}
