//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/gofrs/uuid"
	duration "github.com/golang/protobuf/ptypes/duration"
	"github.com/layer5io/meshery/models"
	SMPS "github.com/layer5io/service-mesh-performance-specification/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LoadTestPrefencesHandler is used for persisting load test preferences
func (h *Handler) LoadTestPrefencesHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	q := req.FormValue("qps")
	qps, err := strconv.Atoi(q)
	if err != nil {
		err = errors.Wrap(err, "unable to parse qps")
		logrus.Error(err)
		http.Error(w, "please provide a valid value for qps", http.StatusBadRequest)
		return
	}
	if qps < 0 {
		http.Error(w, "please provide a valid value for qps", http.StatusBadRequest)
		return
	}

	dur := req.FormValue("t")
	durT, err := time.ParseDuration(dur)
	if err != nil {
		err = errors.Wrap(err, "unable to parse t as a duration")
		logrus.Error(err)
		http.Error(w, "please provide a valid value for t", http.StatusBadRequest)
		return
	}

	cu := req.FormValue("c")
	c, err := strconv.Atoi(cu)
	if err != nil {
		err = errors.Wrap(err, "unable to parse c")
		logrus.Error(err)
		http.Error(w, "please provide a valid value for c", http.StatusBadRequest)
		return
	}
	if c < 0 {
		http.Error(w, "please provide a valid value for c", http.StatusBadRequest)
		return
	}

	gen := req.FormValue("gen")
	// TODO: after we have interfaces for load generators in place, we need to make a generic check, for now using a hard coded one
	var loadGen SMPS.PerformanceTestConfig_Client_LoadGenerator
	switch gen {
	case "fortio":
		loadGen = SMPS.PerformanceTestConfig_Client_fortio
	case "wrk2":
		loadGen = SMPS.PerformanceTestConfig_Client_wrk2
	default:
		logrus.Error("invalid value for gen")
		http.Error(w, "please provide a valid value for gen (load generator)", http.StatusBadRequest)
		return
	}

	proto := req.FormValue("protocol")
	var protocol SMPS.PerformanceTestConfig_Client_Protocol
	switch proto {
	case "HTTP":
		protocol = SMPS.PerformanceTestConfig_Client_http
	case "gRPC":
		protocol = SMPS.PerformanceTestConfig_Client_grpc
	case "TCP":
		protocol = SMPS.PerformanceTestConfig_Client_tcp
	default:
		logrus.Error("invalid value for protocol")
		http.Error(w, "please provide a valid value for gen (load generator)", http.StatusBadRequest)
		return
	}

	testUUID := req.FormValue("uuid")

	headersString := req.FormValue("headers")
	cookiesString := req.FormValue("cookies")
	contentType := req.FormValue("contentType")
	bodyString := req.FormValue("reqBody")

	headersPtr := h.jsonToMap(headersString)
	cookiesPtr := h.jsonToMap(cookiesString)

	headers := make(map[string]string)
	cookies := make(map[string]string)

	if headersPtr != nil {
		headers = *headersPtr
	}
	if cookiesPtr != nil {
		cookies = *cookiesPtr
	}

	testConfig := &SMPS.PerformanceTestConfig{}
	testConfig.Id = testUUID
	testConfig.Clients = []*SMPS.PerformanceTestConfig_Client{}
	testConfig.Clients = append(testConfig.Clients, &SMPS.PerformanceTestConfig_Client{
		Protocol:      protocol,
		Connections:   int32(c),
		Rps:           int64(qps),
		Body:          bodyString,
		ContentType:   contentType,
		Headers:       headers,
		Cookies:       cookies,
		LoadGenerator: loadGen,
	})
	testConfig.Duration = &duration.Duration{
		Seconds: int64(durT.Seconds()),
		Nanos:   int32(durT.Nanoseconds()),
	}

	newUid, _ := uuid.NewV4()
	prefObj.LoadTestPreferences[newUid.String()] = testConfig

	if err = provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
		logrus.Errorf("unable to save user preferences: %v", err)
		http.Error(w, "unable to save user preferences", http.StatusInternalServerError)
		return
	}

	respBody, _ := json.Marshal(&testConfig)
	_, _ = w.Write(respBody)
}
