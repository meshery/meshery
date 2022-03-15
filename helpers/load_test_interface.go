package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/durationpb"

	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/config/metrics/v3"
)

var (
	nighthawkStatus  sync.Mutex
	nighthawkRunning = false
)

// FortioLoadTest is the actual code which invokes Fortio to run the load test
func FortioLoadTest(opts *models.LoadTestOptions) (map[string]interface{}, *periodic.RunnerResults, error) {
	defaults := &periodic.DefaultRunnerOptions
	httpOpts, err := sharedHTTPOptions(opts)
	if err != nil {
		return nil, nil, ErrGeneratingLoadTest(err)
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
		return nil, nil, ErrRunningTest(err)
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
		return nil, nil, ErrConvertingResultToMap(err)
	}

	resultsMap := map[string]interface{}{}
	err = json.Unmarshal(bd, &resultsMap)
	if err != nil {
		return nil, nil, ErrUnmarshal(err, "data to map")
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
		return nil, nil, ErrGrpcSupport(err, "Wrk2")
	}
	var gres *api.GoWRK2
	gres, err = api.WRKRun(ro)
	if err == nil {
		logrus.Debugf("WRK Result: %+v", gres)
		res, err = api.TransformWRKToFortio(gres, ro)
	}

	if err != nil {
		return nil, nil, ErrRunningTest(err)
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
		return nil, nil, ErrConvertingResultToMap(err)
	}

	resultsMap := map[string]interface{}{}
	err = json.Unmarshal(bd, &resultsMap)
	if err != nil {
		return nil, nil, ErrUnmarshal(err, "data to map")
	}
	logrus.Debugf("Mapped version of the test: %+#v", resultsMap)
	return resultsMap, result, nil
}

func startNighthawkServer(timeout int64) error {
	nighthawkStatus.Lock()
	defer nighthawkStatus.Unlock()

	curDir, err := os.Getwd()
	if err != nil {
		return ErrStartingNighthawkServer(err)
	}

	command := filepath.Join(curDir, "nighthawk_service")
	transformCommand := filepath.Join(curDir, "nighthawk_output_transform")

	_, err = os.Stat(command)
	if err != nil {
		return ErrStartingNighthawkServer(err)
	}
	cmd := exec.Command(command)
	if !nighthawkRunning {
		err := cmd.Start()
		if err != nil {
			nighthawkStatus.Unlock()
			return ErrStartingNighthawkServer(err)
		}
		nighthawkRunning = true
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			nighthawkRunning = false
			return
		}
	}()

	_, err = os.Stat(transformCommand)
	if err != nil {
		nighthawkStatus.Unlock()
		return ErrStartingNighthawkServer(err)
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
	return ErrStartingNighthawkServer(err)
}

// NighthawkLoadTest is the actual code which invokes nighthawk to run the load test
func NighthawkLoadTest(opts *models.LoadTestOptions) (map[string]interface{}, *periodic.RunnerResults, error) {
	err := startNighthawkServer(int64(opts.Duration))
	if err != nil {
		return nil, nil, ErrRunningNighthawkServer(err)
	}

	u, err := url.Parse(opts.URL)
	if err != nil {
		return nil, nil, ErrRunningTest(err)
	}
	rURL := u.Host
	if u.Hostname() == "localhost" {
		if u.Port() != "" {
			rURL = fmt.Sprintf("0.0.0.0:%s%s", u.Port(), u.Path)
		} else {
			rURL = fmt.Sprintf("0.0.0.0%s", u.Path)
		}
	}

	if u.Port() == "" {
		if u.Scheme == "http" {
			rURL = fmt.Sprintf("http://%s:80%s", u.Hostname(), u.Path)
		} else {
			rURL = fmt.Sprintf("https://%s:443%s", u.Hostname(), u.Path)
		}
		// Add support for more protocols here
	}

	// advanced request options supported by nighthawk
	requestOptions := &nighthawk_proto.RequestOptions{}

	headers := make([]*v3.HeaderValueOption, 0)
	if opts.Headers != nil {
		for key, value := range *opts.Headers {
			headers = append(headers, &v3.HeaderValueOption{
				Header: &v3.HeaderValue{
					Key:   key,
					Value: value,
				},
			})
		}
	}

	if opts.Cookies != nil {
		cookies := ""
		for key, val := range *opts.Cookies {
			cookies += fmt.Sprintf(" %s=%s;", key, val)
		}
		headers = append(headers, &v3.HeaderValueOption{
			Header: &v3.HeaderValue{
				Key:   "Cookie",
				Value: cookies,
			},
		})
	}

	if len(opts.ContentType) > 0 {
		headers = append(headers, &v3.HeaderValueOption{
			Header: &v3.HeaderValue{
				Key:   "Content-Type",
				Value: opts.ContentType,
			},
		})
	}

	requestOptions.RequestHeaders = headers

	// Nighthawk doesn't send the specified request payload but instead sends
	// 'a' characters corresponding to that body's size
	// TODO: Request method should be specifiable through UI
	if reqBodyLength := len(opts.Body); reqBodyLength > 0 {
		requestOptions.RequestBodySize = &wrappers.UInt32Value{Value: uint32(len(opts.Body))}
		requestOptions.RequestMethod = v3.RequestMethod_POST
	} else {
		requestOptions.RequestBodySize = &wrappers.UInt32Value{Value: uint32(0)}
		requestOptions.RequestMethod = v3.RequestMethod_GET
	}

	ro := &nighthawk_proto.CommandLineOptions{
		OneofDurationOptions: &nighthawk_proto.CommandLineOptions_Duration{
			Duration: durationpb.New(opts.Duration),
		},
		Timeout: durationpb.New(10 * time.Second),
		// TODO: support multiple http versions
		Concurrency:         &wrappers.StringValue{Value: fmt.Sprint(opts.HTTPNumThreads)},
		Verbosity:           &nighthawk_proto.Verbosity{Value: nighthawk_proto.Verbosity_INFO},
		OutputFormat:        &nighthawk_proto.OutputFormat{Value: nighthawk_proto.OutputFormat_FORTIO},
		PrefetchConnections: &wrappers.BoolValue{Value: false},
		BurstSize:           &wrappers.UInt32Value{Value: uint32(0)},
		AddressFamily:       &nighthawk_proto.AddressFamily{Value: nighthawk_proto.AddressFamily_AUTO},
		OneofRequestOptions: &nighthawk_proto.CommandLineOptions_RequestOptions{
			RequestOptions: requestOptions,
		},
		// use default values for nighthawk's configuration to avoid any unexpected
		// failures until there is a dynamic way to specify this
		// MaxPendingRequests:       &wrappers.UInt32Value{Value: uint32(10)},
		// MaxActiveRequests:        &wrappers.UInt32Value{Value: uint32(100)},
		// MaxRequestsPerConnection: &wrappers.UInt32Value{Value: uint32(100)},
		SequencerIdleStrategy: &nighthawk_proto.SequencerIdleStrategy{Value: nighthawk_proto.SequencerIdleStrategy_DEFAULT},
		OneofUri: &nighthawk_proto.CommandLineOptions_Uri{
			Uri: &wrappers.StringValue{Value: rURL},
		},
		ExperimentalH1ConnectionReuseStrategy: &nighthawk_proto.H1ConnectionReuseStrategy{
			Value: nighthawk_proto.H1ConnectionReuseStrategy_DEFAULT,
		},
		TerminationPredicates: make(map[string]uint64),
		// Used for specifying parameters for failing execution. Use defailt for now
		//FailurePredicates:                    make(map[string]uint64),
		OpenLoop:                             &wrappers.BoolValue{Value: false},
		JitterUniform:                        durationpb.New(0 * time.Second),
		ExperimentalH2UseMultipleConnections: &wrappers.BoolValue{Value: false},
		Labels:                               []string{opts.Name, " -_- ", rURL},
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

	qps := opts.HTTPQPS
	// set QPS only if given QPS > 0
	// user nighthawk's default QPS otherwise
	if qps > 0 {
		ro.RequestsPerSecond = &wrappers.UInt32Value{Value: uint32(qps)}
	}

	if opts.SupportedLoadTestMethods == 2 {
		return nil, nil, ErrGrpcSupport(err, "Nighthawk")
	}

	c, err := nighthawk_client.New(nighthawk_client.Options{
		ServerHost: "0.0.0.0",
		ServerPort: 8443,
	})
	if err != nil {
		return nil, nil, ErrRunningTest(err)
	}

	logrus.Info("starting test")

	client, err := c.Handler.ExecutionStream(context.TODO())
	if err != nil {
		return nil, nil, ErrRunningTest(err)
	}

	err = client.Send(&nighthawk_proto.ExecutionRequest{
		CommandSpecificOptions: &nighthawk_proto.ExecutionRequest_StartRequest{
			StartRequest: &nighthawk_proto.StartRequest{
				Options: ro,
			},
		},
	})
	if err != nil {
		return nil, nil, ErrRunningTest(err)
	}

	var res1 *nighthawk_proto.ExecutionResponse
	logrus.Info("listening to test events")
	for {
		var err error
		res1, err = client.Recv()
		if err != nil {
			return nil, nil, ErrRunningTest(err)
		}
		if res1 != nil {
			break
		}
	}

	d, err := nighthawk_client.Transform(res1)
	if err != nil {
		return nil, nil, ErrTransformingData(err)
	}

	var result *periodic.RunnerResults
	if opts.SupportedLoadTestMethods == 2 {
		gres := &fgrpc.GRPCRunnerResults{}
		err := json.Unmarshal(d, gres)
		if err != nil {
			return nil, nil, ErrUnmarshal(err, "data to object")
		}
		result = gres.Result()
	} else {
		hres := &fhttp.HTTPRunnerResults{}
		err := json.Unmarshal(d, hres)
		if err != nil {
			return nil, nil, ErrUnmarshal(err, "data to object")
		}
		result = hres.Result()
	}
	if err != nil {
		return nil, nil, ErrUnmarshal(err, "results to map")
	}

	resultsMap := map[string]interface{}{}
	err = json.Unmarshal(d, &resultsMap)
	if err != nil {
		return nil, nil, ErrUnmarshal(err, "data to map")
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
				return nil, ErrAddAndValidateExtraHeader(err)
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
			return nil, ErrAddAndValidateExtraHeader(err)
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
