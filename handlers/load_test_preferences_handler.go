//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/protobuf/jsonpb"

	duration "github.com/golang/protobuf/ptypes/duration"
	"github.com/layer5io/meshery/models"
	SMPS "github.com/layer5io/service-mesh-performance-specification/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type custTestConf struct {
	Val *SMPS.PerformanceTestConfig
}

func (c *custTestConf) MarshalJSON() ([]byte, error) {
	m := jsonpb.Marshaler{
		EmitDefaults: true,
	}
	val, err := m.MarshalToString(c.Val)
	return []byte(val), err
}

// LoadTestPrefencesHandler is used for persisting load test preferences
func (h *Handler) LoadTestPrefencesHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {

	if req.Method == http.MethodPost || req.Method == http.MethodPut {

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

		endpoint := req.FormValue("endpoint")
		if endpoint == "" {
			logrus.Error("unable to find endpoint")
			http.Error(w, "please provide a value for endpoint", http.StatusBadRequest)
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

		headersString := req.FormValue("headers")
		cookiesString := req.FormValue("cookies")
		labelsString := req.FormValue("labels")
		contentType := req.FormValue("contentType")
		bodyString := req.FormValue("reqBody")

		headersPtr := h.jsonToMap(headersString)
		cookiesPtr := h.jsonToMap(cookiesString)
		labelsPtr := h.jsonToMap(labelsString)

		headers := make(map[string]string)
		cookies := make(map[string]string)
		labels := make(map[string]string)

		if headersPtr != nil {
			headers = *headersPtr
		}
		if cookiesPtr != nil {
			cookies = *cookiesPtr
		}
		if labelsPtr != nil {
			labels = *labelsPtr
		}

		testConfig := &SMPS.PerformanceTestConfig{}
		testConfig.Labels = labels
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
			EndpointUrl:   endpoint,
		})
		testConfig.Duration = &duration.Duration{
			Seconds: int64(durT.Seconds()),
			// Nanos:   int32(durT.Nanoseconds()),
		}

		tid, err := prefObj.CreateUpdateLoadTestConfig(testConfig)
		if err != nil {
			logrus.Errorf("unable to save user preferences: %v", err)
			http.Error(w, "unable to save user preferences", http.StatusInternalServerError)
			return
		}
		if err = provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
			logrus.Errorf("unable to save user preferences: %v", err)
			http.Error(w, "unable to save user preferences", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte(tid))
	} else if req.Method == http.MethodDelete {
		testUUID := req.URL.Query().Get("uuid")
		if testUUID == "" {
			logrus.Error("field uuid not found")
			http.Error(w, "field uuid not found", http.StatusBadRequest)
			return
		}
		prefObj.DeleteLoadTestConfig(testUUID)
	} else if req.Method == http.MethodGet {
		testUUID := req.URL.Query().Get("uuid")
		if testUUID == "" {
			testObj := prefObj.ReadAllLoadTestConfig()
			custTestObjs := []*custTestConf{}
			for _, tst := range testObj {
				custTestObjs = append(custTestObjs, &custTestConf{
					Val: tst,
				})
			}
			body, err := json.Marshal(&custTestObjs)
			if err != nil {
				logrus.Error("error reading database")
				http.Error(w, "error reading database", http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(body)
		} else {
			testObj := prefObj.ReadLoadTestConfig(testUUID)
			if testObj == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			m := jsonpb.Marshaler{
				EmitDefaults: true,
			}
			if err := m.Marshal(w, testObj); err != nil {
				logrus.Error("error reading database: %v", err)
				http.Error(w, "error reading database", http.StatusInternalServerError)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

}
