package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"fortio.org/fortio/fgrpc"
	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/periodic"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/layer5io/gowrk2/api"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils"
	nighthawk_client "github.com/layer5io/nighthawk-go/pkg/client"
	nighthawk_proto "github.com/layer5io/nighthawk-go/pkg/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/durationpb"

	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/config/metrics/v3"
)

var (
	nighthawkStatus sync.Mutex
)

// FortioLoadTest is the actual code which invokes Fortio to run the load test
func FortioLoadTest(opts *models.LoadTestOptions) (map[string]interface{}, *periodic.RunnerResults, error) {
	defaults := &periodic.DefaultRunnerOptions
	httpOpts, err := sharedHTTPOptions(opts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "generating load test options failed")
	}
	if opts.IsInsecure {
		httpOpts.Insecure = true
	}
	rURL := httpOpts.URL
	out := os.Stdout
	qps := opts.HTTPQPS // TODO possibly use translated <=0 to "max" from results/options normalization in periodic/
	if qps <= 0 {
		qps = -1 // 0==unitialized struct == default duration, -1 (0 for flag) is max
	}
	labels := opts.Name + " -_- " + rURL
	ro := periodic.RunnerOptions{
		QPS:         qps,
		Duration:    opts.Duration,
		NumThreads:  opts.HTTPNumThreads,
		Percentiles: []float64{50, 75, 90, 99, 99.9},
		Resolution:  defaults.Resolution,
		Out:         out,
		Labels:      labels,
		Exactly:     0,
	}
	var res periodic.HasRunnerResult
	if opts.SupportedLoadTestMethods == 2 {
		o := fgrpc.GRPCRunnerOptions{
			RunnerOptions:      ro,
			Destination:        rURL,
			CACert:             opts.CACert,
			Service:            opts.GRPCHealthSvc,
			Streams:            opts.GRPCStreamsCount,
			AllowInitialErrors: opts.AllowInitialErrors,
			Payload:            httpOpts.PayloadString(),
			Delay:              opts.GRPCPingDelay,
			UsePing:            opts.GRPCDoPing,
			UnixDomainSocket:   httpOpts.UnixDomainSocket,
		}
		res, err = fgrpc.RunGRPCTest(&o)
	} else {
		o := fhttp.HTTPRunnerOptions{
			HTTPOptions:        *httpOpts,
			RunnerOptions:      ro,
			Profiler:           "",
			AllowInitialErrors: opts.AllowInitialErrors,
			AbortOn:            0,
		}
		res, err = fhttp.RunHTTPTest(&o)
	}
	if err != nil {
		err = errors.Wrap(err, "error while running tests")
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Debugf("original version of the test: %+#v", res)

	var result *periodic.RunnerResults
	var bd []byte
	if opts.SupportedLoadTestMethods == 2 {
		gres, _ := res.(*fgrpc.GRPCRunnerResults)
		bd, err = json.Marshal(gres)
		result = gres.Result()
	} else {
		hres, _ := res.(*fhttp.HTTPRunnerResults)
		bd, err = json.Marshal(hres)
		result = hres.Result()
	}
	if err != nil {
		err = errors.Wrap(err, "error while converting results to map")
		logrus.Error(err)
		return nil, nil, err
	}

	resultsMap := map[string]interface{}{}
	err = json.Unmarshal(bd, &resultsMap)
	if err != nil {
		err = errors.Wrap(err, "error while unmarshaling data to map")
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Debugf("Mapped version of the test: %+#v", resultsMap)
	return resultsMap, result, nil
}

// WRK2LoadTest is the actual code which invokes Wrk2 to run the load test
func WRK2LoadTest(opts *models.LoadTestOptions) (map[string]interface{}, *periodic.RunnerResults, error) {
	qps := opts.HTTPQPS // TODO possibly use translated <=0 to "max" from results/options normalization in periodic/
	if qps <= 0 {
		qps = -1 // 0==unitialized struct == default duration, -1 (0 for flag) is max
	}
	rURL := strings.TrimLeft(opts.URL, " \t\r\n")

	labels := opts.Name + " -_- " + rURL
	ro := &api.GoWRK2Config{
		DurationInSeconds: opts.Duration.Seconds(),
		Thread:            opts.HTTPNumThreads,
		RQPS:              qps,
		URL:               rURL,
		Labels:            labels,
		Percentiles:       []float64{50, 75, 90, 99, 99.99, 99.999},
	}
	var res periodic.HasRunnerResult
	var err error
	if opts.SupportedLoadTestMethods == 2 {
		err = errors.New("wrk2 does not support gRPC at the moment")
		logrus.Error(err)
		return nil, nil, err
	}
	var gres *api.GoWRK2
	gres, err = api.WRKRun(ro)
	if err == nil {
		logrus.Debugf("WRK Result: %+v", gres)
		res, err = api.TransformWRKToFortio(gres, ro)
	}

	if err != nil {
		err = errors.Wrap(err, "error while running tests")
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Debugf("original version of the test: %+#v", res)

	var result *periodic.RunnerResults
	var bd []byte
	if opts.SupportedLoadTestMethods == 2 {
		gres, _ := res.(*fgrpc.GRPCRunnerResults)
		bd, err = json.Marshal(gres)
		result = gres.Result()
	} else {
		hres, _ := res.(*fhttp.HTTPRunnerResults)
		bd, err = json.Marshal(hres)
		result = hres.Result()
	}
	if err != nil {
		err = errors.Wrap(err, "error while converting results to map")
		logrus.Error(err)
		return nil, nil, err
	}

	resultsMap := map[string]interface{}{}
	err = json.Unmarshal(bd, &resultsMap)
	if err != nil {
		err = errors.Wrap(err, "error while unmarshaling data to map")
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Debugf("Mapped version of the test: %+#v", resultsMap)
	return resultsMap, result, nil
}

func startNighthawkServer(timeout int64) error {
	nighthawkStatus.Lock()
	command := "./nighthawk_service"
	transformCommand := "./nighthawk_output_transform"
	err := exec.Command(command).Start()
	if err != nil {
		nighthawkStatus.Unlock()
		return err
	}

	_, err = os.Stat(transformCommand)
	if err != nil {
		nighthawkStatus.Unlock()
		return err
	}

	for timeout != 0 {
		if utils.TcpCheck(&utils.HostPort{
			Address: "0.0.0.0",
			Port:    8443,
		}, nil) {
			return nil
		}
		timeout--
		time.Sleep(1 * time.Second)
	}
	nighthawkStatus.Unlock()
	return errors.New("unable to start nighthawk server")
}

// NighthawkLoadTest is the actual code which invokes nighthawk to run the load test
func NighthawkLoadTest(opts *models.LoadTestOptions) (map[string]interface{}, *periodic.RunnerResults, error) {
	err := startNighthawkServer(int64(opts.Duration))
	if err != nil {
		err = errors.Wrap(err, "error while running nighthawk server")
		logrus.Error(err)
		return nil, nil, err
	}

	qps := opts.HTTPQPS

	if qps <= 0 {
		qps = -1 // 0==unitialized struct == default duration, -1 (0 for flag) is max
	}

	u, err := url.Parse(opts.URL)
	if err != nil {
		err = errors.Wrap(err, "error while running tests")
		logrus.Error(err)
		return nil, nil, err
	}
	rURL := u.Host
	if u.Hostname() == "localhost" {
		rURL = fmt.Sprintf("0.0.0.0:%s", u.Port())
	}

	ro := &nighthawk_proto.CommandLineOptions{
		RequestsPerSecond: &wrappers.UInt32Value{Value: uint32(qps)},
		Connections:       &wrappers.UInt32Value{Value: uint32(2)},
		OneofDurationOptions: &nighthawk_proto.CommandLineOptions_Duration{
			Duration: durationpb.New(opts.Duration),
		},
		Timeout:             durationpb.New(10 * time.Second),
		H2:                  &wrappers.BoolValue{Value: false},
		Concurrency:         &wrappers.StringValue{Value: fmt.Sprint(opts.HTTPNumThreads)},
		Verbosity:           &nighthawk_proto.Verbosity{Value: nighthawk_proto.Verbosity_INFO},
		OutputFormat:        &nighthawk_proto.OutputFormat{Value: nighthawk_proto.OutputFormat_FORTIO},
		PrefetchConnections: &wrappers.BoolValue{Value: false},
		BurstSize:           &wrappers.UInt32Value{Value: uint32(0)},
		AddressFamily:       &nighthawk_proto.AddressFamily{Value: nighthawk_proto.AddressFamily_AUTO},
		OneofRequestOptions: &nighthawk_proto.CommandLineOptions_RequestOptions{
			RequestOptions: &nighthawk_proto.RequestOptions{
				RequestMethod:   v3.RequestMethod_GET,
				RequestHeaders:  make([]*v3.HeaderValueOption, 0),
				RequestBodySize: &wrappers.UInt32Value{Value: uint32(10)},
			},
		},
		MaxPendingRequests:       &wrappers.UInt32Value{Value: uint32(10)},
		MaxActiveRequests:        &wrappers.UInt32Value{Value: uint32(100)},
		MaxRequestsPerConnection: &wrappers.UInt32Value{Value: uint32(100)},
		SequencerIdleStrategy:    &nighthawk_proto.SequencerIdleStrategy{Value: nighthawk_proto.SequencerIdleStrategy_DEFAULT},
		OneofUri: &nighthawk_proto.CommandLineOptions_Uri{
			Uri: &wrappers.StringValue{Value: rURL},
		},
		ExperimentalH1ConnectionReuseStrategy: &nighthawk_proto.H1ConnectionReuseStrategy{
			Value: nighthawk_proto.H1ConnectionReuseStrategy_DEFAULT,
		},
		TerminationPredicates:                make(map[string]uint64),
		FailurePredicates:                    make(map[string]uint64),
		OpenLoop:                             &wrappers.BoolValue{Value: false},
		JitterUniform:                        durationpb.New(0 * time.Second),
		ExperimentalH2UseMultipleConnections: &wrappers.BoolValue{Value: false},
		Labels:                               make([]string, 0),
		//TransportSocket: &v3.TransportSocket{
		//	Name: "test",
		//	ConfigType: &v3.TransportSocket_TypedConfig{
		//		TypedConfig: &any.Any{
		//			TypeUrl: "",
		//			Value:   []byte{},
		//		},
		//	},
		//},
		SimpleWarmup:              &wrappers.BoolValue{Value: false},
		StatsSinks:                make([]*v32.StatsSink, 0),
		StatsFlushInterval:        &wrappers.UInt32Value{Value: uint32(5)},
		LatencyResponseHeaderName: &wrappers.StringValue{Value: ""},
		//ScheduledStart: &timestamp.Timestamp{
		//	Seconds: 0,
		//	Nanos:   0,
		//},
		ExecutionId: &wrappers.StringValue{Value: "gg"},
	}

	if opts.SupportedLoadTestMethods == 2 {
		err := errors.New("nighthawk does not support gRPC load testing")
		logrus.Error(err)
		return nil, nil, err
	}

	c, err := nighthawk_client.New(nighthawk_client.Options{
		ServerHost: "0.0.0.0",
		ServerPort: 8443,
	})
	if err != nil {
		err = errors.Wrap(err, "error while running tests")
		logrus.Error(err)
		return nil, nil, err
	}

	logrus.Info("starting test")

	client, err := c.Handler.ExecutionStream(context.TODO())
	if err != nil {
		err = errors.Wrap(err, "error while running tests")
		logrus.Error(err)
		return nil, nil, err
	}

	err = client.Send(&nighthawk_proto.ExecutionRequest{
		CommandSpecificOptions: &nighthawk_proto.ExecutionRequest_StartRequest{
			StartRequest: &nighthawk_proto.StartRequest{
				Options: ro,
			},
		},
	})
	if err != nil {
		err = errors.Wrap(err, "error while running tests")
		logrus.Error(err)
		return nil, nil, err
	}

	var res1 *nighthawk_proto.ExecutionResponse
	logrus.Info("listening to test events")
	for {
		var err error
		res1, err = client.Recv()
		if err != nil {
			err = errors.Wrap(err, "error while running tests")
			logrus.Error(err)
			return nil, nil, err
		}
		if res1 != nil {
			break
		}
	}

	try, _ := json.Marshal(res1)

	logrus.Debugf("original version of the test: %s", string(try))

	d, err := nighthawk_client.Transform(res1, "fortio")
	if err != nil {
		err = errors.Wrap(err, "error while transforming data")
		logrus.Error(err)
		return nil, nil, err
	}

	var result *periodic.RunnerResults
	if opts.SupportedLoadTestMethods == 2 {
		gres := &fgrpc.GRPCRunnerResults{}
		err := json.Unmarshal(d, gres)
		if err != nil {
			err = errors.Wrap(err, "error while unmarshaling data to object")
			logrus.Error(err)
			return nil, nil, err
		}
		result = gres.Result()
	} else {
		hres := &fhttp.HTTPRunnerResults{}
		err := json.Unmarshal(d, hres)
		if err != nil {
			err = errors.Wrap(err, "error while unmarshaling data to object")
			logrus.Error(err)
			return nil, nil, err
		}
		result = hres.Result()
	}
	if err != nil {
		err = errors.Wrap(err, "error while converting results to map")
		logrus.Error(err)
		return nil, nil, err
	}

	resultsMap := map[string]interface{}{}
	err = json.Unmarshal(d, &resultsMap)
	if err != nil {
		err = errors.Wrap(err, "error while unmarshaling data to map")
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Debugf("Mapped version of the test: %+#v", resultsMap)
	return resultsMap, result, nil
}

// sharedHTTPOptions is the flag->httpoptions transfer code shared between
// fortio_main and fcurl.
func sharedHTTPOptions(opts *models.LoadTestOptions) (*fhttp.HTTPOptions, error) {
	url := strings.TrimLeft(opts.URL, " \t\r\n")
	httpOpts := fhttp.HTTPOptions{}
	httpOpts.URL = url
	httpOpts.HTTP10 = false
	httpOpts.DisableFastClient = false
	httpOpts.DisableKeepAlive = false
	httpOpts.AllowHalfClose = false
	httpOpts.Compression = false
	httpOpts.HTTPReqTimeOut = fhttp.HTTPReqTimeOutDefaultValue
	httpOpts.Insecure = opts.IsInsecure
	httpOpts.UserCredentials = ""
	httpOpts.ContentType = ""
	httpOpts.FollowRedirects = true
	httpOpts.DisableFastClient = true

	if opts.Headers != nil {
		for key, val := range *opts.Headers {
			err := httpOpts.AddAndValidateExtraHeader(key + ":" + val)
			if err != nil {
				return nil, err
			}
		}
	}

	if opts.Cookies != nil {
		cookies := ""
		for key, val := range *opts.Cookies {
			cookies += fmt.Sprintf(" %s=%s;", key, val)
		}
		err := httpOpts.AddAndValidateExtraHeader("Cookie" + ":" + cookies)
		if err != nil {
			return nil, err
		}
	}
	if len(opts.Body) > 0 {
		httpOpts.Payload = opts.Body
	}
	if len(opts.ContentType) > 0 {
		httpOpts.ContentType = opts.ContentType
	}

	return &httpOpts, nil
}
