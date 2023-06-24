package models

import (
	"encoding/json"
	"strconv"
	"time"

	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/periodic"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// LoadGenerator - represents the load generator type
type LoadGenerator string

const (
	// FortioLG - represents the Fortio load generator
	FortioLG LoadGenerator = "fortio"

	// Wrk2LG - represents the wrk2 load generator
	Wrk2LG LoadGenerator = "wrk2"

	// NighthawkLG - represents the nighthawk load generator
	NighthawkLG LoadGenerator = "nighthawk"
)

// Name - retrieves a string value for the generator
func (l LoadGenerator) Name() string {
	return string(l)
}

// SupportedLoadTestMethods - represents the load test method type
type SupportedLoadTestMethods int

const (
	// HTTP Load Test
	HTTP SupportedLoadTestMethods = 1

	// gRPC Load Test
	// gRPC SupportedLoadTestMethods = 2

	// TCP Load Test
	TCP SupportedLoadTestMethods = 3
)

// LoadTestOptions represents the load test options
type LoadTestOptions struct {
	Name     string
	Location string
	URL      string

	HTTPQPS float64

	HTTPNumThreads int

	Headers     *map[string]string
	Cookies     *map[string]string
	Body        []byte
	ContentType string

	IsInsecure bool
	Duration   time.Duration

	LoadGenerator LoadGenerator

	SupportedLoadTestMethods SupportedLoadTestMethods

	Cert, Key, CACert string

	AllowInitialErrors bool

	// Values required for fortio gRPC health & ping test
	GRPCStreamsCount int
	GRPCDoHealth     bool
	GRPCHealthSvc    string
	GRPCDoPing       bool
	GRPCPingDelay    time.Duration
}

// LoadTestStatus - used for representing load test status
type LoadTestStatus string

const (
	// LoadTestError - respresents an error status
	LoadTestError LoadTestStatus = "error"

	// LoadTestInfo - represents a info status
	LoadTestInfo LoadTestStatus = "info"

	// LoadTestSuccess - represents a success status
	LoadTestSuccess LoadTestStatus = "success"
)

// LoadTestResponse - used to bundle the response with status to the client
type LoadTestResponse struct {
	Status  LoadTestStatus `json:"status,omitempty"`
	Message string         `json:"message,omitempty"`
	Result  *MesheryResult `json:"result,omitempty"`
}

// MesheryResult - represents the results from Meshery test run to be shipped
type MesheryResult struct {
	ID                 uuid.UUID              `json:"meshery_id,omitempty"`
	Name               string                 `json:"name,omitempty"`
	Mesh               string                 `json:"mesh,omitempty"`
	PerformanceProfile *uuid.UUID             `json:"performance_profile,omitempty"`
	TestID             string                 `json:"test_id"`
	Result             map[string]interface{} `json:"runner_results,omitempty" gorm:"type:JSONB"`

	ServerMetrics     interface{} `json:"server_metrics,omitempty" gorm:"type:JSONB"`
	ServerBoardConfig interface{} `json:"server_board_config,omitempty" gorm:"type:JSONB"`

	TestStartTime          *time.Time         `json:"test_start_time,omitempty"`
	PerformanceProfileInfo PerformanceProfile `json:"-,omitempty" gorm:"constraint:OnDelete:SET NULL;foreignKey:PerformanceProfile"`

	UpdatedAt string `json:"updated_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UserID    string `json:"user_id,omitempty"`
}

// ConvertToSpec - converts meshery result to SMP
func (m *MesheryResult) ConvertToSpec() (*PerformanceSpec, error) {
	b := &PerformanceSpec{
		Latencies: &LatenciesMs{},
	}
	var (
		results periodic.HasRunnerResult
	)
	retcodesString, _ := m.Result["RetCodes"].(map[string]interface{})
	logrus.Debugf("retcodes: %+v, %T", m.Result["RetCodes"], m.Result["RetCodes"])
	retcodes := map[int]int64{}
	for k, v := range retcodesString {
		k1, _ := strconv.Atoi(k)
		retcodes[k1], _ = v.(int64)
	}
	m.Result["RetCodes"] = retcodes
	// loadGenerator := m.Result["load-generator"].(string)
	logrus.Debugf("result to be converted: %+v", m)
	if m.Result["RunType"].(string) == "HTTP" {
		httpResults := &fhttp.HTTPRunnerResults{}
		resJ, err := json.Marshal(m.Result)
		if err != nil {
			return nil, ErrMarshal(err, "Perf Results")
		}
		err = json.Unmarshal(resJ, httpResults)
		if err != nil {
			return nil, ErrUnmarshal(err, "Perf Results")
		}

		results = httpResults
		logrus.Debugf("httpresults: %+v", httpResults)
	}

	result := results.Result()
	// b.SMPVersion = MoreDetails.SmiVersion
	b.SMPVersion = "test_smp_version"
	b.id = "test_id"
	b.labels = map[string]string{"test_label": "test_value"}
	b.StartTime = result.StartTime
	b.EndTime = result.StartTime.Add(result.ActualDuration)
	b.Latencies = &LatenciesMs{
		Min:     result.DurationHistogram.Min,
		Max:     result.DurationHistogram.Max,
		Average: result.DurationHistogram.Avg,
	}
	for _, p := range result.DurationHistogram.Percentiles {
		switch p.Percentile {
		case 50:
			b.Latencies.P50 = p.Value
		case 90:
			b.Latencies.P90 = p.Value
		case 99:
			b.Latencies.P99 = p.Value
		}
	}
	b.ActualQPS = result.ActualQPS
	b.DetailsURI = "test_details"
	b.TestID = "test_test_id"
	b.MeshConfigID = "test_meshconfigID"
	b.EnvID = "test_envID"
	// b.LoadGenerator = loadGenerator
	// b.Client.Connections = result.NumThreads
	// b.Client.Internal = false

	// k8sI, ok := m.Result["kubernetes"]
	// if ok {
	// 	k8s, _ := k8sI.(map[string]interface{})
	// 	b.Env.Kubernetes, _ = k8s["server_version"].(string)
	// 	nodesI, okk := k8s["nodes"]
	// 	if okk {
	// 		nodes, okkk := nodesI.([]*K8SNode)
	// 		if okkk {
	// 			b.Env.NodeCount = len(nodes)
	// 		}
	// 	}
	// }
	return b, nil
}
