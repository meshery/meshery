package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"fortio.org/fortio/fgrpc"
	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/periodic"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/layer5io/gowrk2/api"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/utils"
	nighthawk_client "github.com/layer5io/nighthawk-go/pkg/client"
	nighthawk_proto "github.com/layer5io/nighthawk-go/pkg/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

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
			Service:            opts.GRPCHealthSvc,
			Streams:            opts.GRPCStreamsCount,
			AllowInitialErrors: opts.AllowInitialErrors,
			Payload:            httpOpts.PayloadUTF8(),
			Delay:              opts.GRPCPingDelay,
			UsePing:            opts.GRPCDoPing,
		}

		o.TLSOptions = fhttp.TLSOptions{
			CACert:    opts.CACert,
			UnixDomainSocket: httpOpts.UnixDomainSocket,
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

		logrus.Debugf("options string: %s", opts.Options)
		if opts.Options != "" {
			logrus.Debugf("Fortio config: %+#v", o)
			err := json.Unmarshal([]byte(opts.Options), &o)
			if err != nil {
				return nil, nil, ErrUnmarshal(err, "options string")
			}
			logrus.Debugf("Fortio config with options: %+#v", o)
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

	logrus.Debugf("options string: %s", opts.Options)
	if opts.Options != "" {
		logrus.Debugf("GoWrk2 config: %+#v", ro)
		err := json.Unmarshal([]byte(opts.Options), &ro)
		if err != nil {
			return nil, nil, ErrUnmarshal(err, "options string")
		}
		logrus.Debugf("GoWrk2 config with options: %+#v", ro)
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

type NighthawkCliOptions nighthawk_proto.CommandLineOptions

func (opt *NighthawkCliOptions) UnmarshalJSON(data []byte) error {

	type Duration struct {
		Seconds int64 `json:"seconds"`
		Nanos   int32 `json:"nanos"`
	}


	o := &struct {	
		RequestsPerSecond 						uint32 				`json:"requests_per_second"`
		Connections 	  						uint32 				`json:"connections"`
		Timeout           						Duration 			`json:"timeout"`
		Concurrency		  						string 				`json:"concurrency"`
		Verbosity	  	  						int32 				`json:"verbosity"`
		OutputFormat	  						int32 	    		`json:"output_format"`
		PrefetchConnections 					bool 				`json:"prefetch_connections"`
		BurstSize								uint32				`json:"burst_size"`
		AddressFamily							int32	    		`json:"address_family"`
		MaxPendingRequests						uint32				`json:"max_pending_requests"`
		MaxActiveRequests   					uint32				`json:"max_active_requests"`
		MaxRequestsPerConnection 				uint32				`json:"max_requests_per_connection"`
		SequencerIdleStrategy 					int32 				`json:"sequencer_idle_strategy"`
		Trace 									string 				`json:"trace"`
		ExperimentalH1ConnectionReuseStrategy 	int32 			  	`json:"experimental_h1_connection_reuse_strategy"`
		TerminationPredicates 					map[string]uint64 	`json:"termination_predicates"`
		FailurePredicates 						map[string]uint64 	`json:"failure_predicates"`
		OpenLoop							 	bool 			  	`json:"open_loop"`
		JitterUniform 							Duration 		  	`json:"jitter_uniform"`
		NighthawkService 						string 			  	`json:"nighthawk_service"`
		MaxConcurrentStreams 					uint32 			  	`json:"max_concurrent_streams"`	
		Labels 									[]string 		  	`json:"labels"`
		SimpleWarmup 							bool 			  	`json:"simple_warmup"`
		StatsFlushInterval 						uint32 			  	`json:"stats_flush_interval"`
		LatencyResponseHeaderName 				string 			  	`json:"latency_response_header_name"`
		ScheduledStart 							Duration 		  	`json:"scheduled_start"`
		ExecutionId 							string 			  	`json:"execution_id"`
	}{}

	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}
	
	opt.RequestsPerSecond = &wrappers.UInt32Value{Value: o.RequestsPerSecond}
	opt.Connections = &wrappers.UInt32Value{Value: o.Connections}
	opt.Timeout = durationpb.New(time.Duration(o.Timeout.Seconds) * time.Second + time.Duration(o.Timeout.Nanos) * time.Nanosecond)
	opt.Concurrency = &wrappers.StringValue{Value: o.Concurrency}
	verbosity := nighthawk_proto.Verbosity_VerbosityOptions(o.Verbosity)
	opt.Verbosity = &nighthawk_proto.Verbosity{Value: verbosity}
	outputFormat := nighthawk_proto.OutputFormat_OutputFormatOptions(o.OutputFormat)
	opt.OutputFormat = &nighthawk_proto.OutputFormat{Value: outputFormat}
	opt.PrefetchConnections = &wrappers.BoolValue{Value: o.PrefetchConnections}
	opt.BurstSize = &wrappers.UInt32Value{Value: o.BurstSize}
	addressFamily := nighthawk_proto.AddressFamily_AddressFamilyOptions(o.AddressFamily)
	opt.AddressFamily = &nighthawk_proto.AddressFamily{Value: addressFamily}
	opt.MaxPendingRequests = &wrappers.UInt32Value{Value: o.MaxPendingRequests}
	opt.MaxActiveRequests = &wrappers.UInt32Value{Value: o.MaxActiveRequests}
	opt.MaxRequestsPerConnection = &wrappers.UInt32Value{Value: o.MaxRequestsPerConnection}
	sequencerIdleStrategy := nighthawk_proto.SequencerIdleStrategy_SequencerIdleStrategyOptions(o.SequencerIdleStrategy)
	opt.SequencerIdleStrategy = &nighthawk_proto.SequencerIdleStrategy{Value: sequencerIdleStrategy}
	opt.Trace = &wrappers.StringValue{Value: o.Trace}
	experimentalH1ConnectionReuseStrategy := nighthawk_proto.H1ConnectionReuseStrategy_H1ConnectionReuseStrategyOptions(o.ExperimentalH1ConnectionReuseStrategy)
	opt.ExperimentalH1ConnectionReuseStrategy = &nighthawk_proto.H1ConnectionReuseStrategy{Value: experimentalH1ConnectionReuseStrategy}
	opt.TerminationPredicates = o.TerminationPredicates
	opt.FailurePredicates = o.FailurePredicates
	opt.OpenLoop = &wrappers.BoolValue{Value: o.OpenLoop}
	opt.JitterUniform = durationpb.New(time.Duration(o.JitterUniform.Seconds) * time.Second + time.Duration(o.JitterUniform.Nanos) * time.Nanosecond)
	opt.NighthawkService = &wrappers.StringValue{Value: o.NighthawkService}
	opt.MaxConcurrentStreams = &wrappers.UInt32Value{Value: o.MaxConcurrentStreams}
	opt.Labels = o.Labels
	opt.SimpleWarmup = &wrappers.BoolValue{Value: o.SimpleWarmup}
	opt.StatsFlushInterval = &wrappers.UInt32Value{Value: o.StatsFlushInterval}
	opt.LatencyResponseHeaderName = &wrappers.StringValue{Value: o.LatencyResponseHeaderName}
	opt.ScheduledStart = timestamppb.New(time.Unix(o.ScheduledStart.Seconds, int64(o.ScheduledStart.Nanos)))
	opt.ExecutionId = &wrappers.StringValue{Value: o.ExecutionId}

	return nil
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

	ro:= NighthawkCliOptions{
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

	logrus.Debugf("options string: %s", opts.Options)
	if opts.Options != "" {
		// logrus.Debugf("Nighthawk CommandLineOptions: %+#v", ro)
		err := json.Unmarshal([]byte(opts.Options), &ro)
		if err != nil {
			return nil, nil, ErrUnmarshal(err, "options string")
		}
		// logrus.Debugf("Nighthawk CommandLineOptions with options: %+#v", ro)
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

	opt := nighthawk_proto.CommandLineOptions(ro)
	err = client.Send(&nighthawk_proto.ExecutionRequest{
		CommandSpecificOptions: &nighthawk_proto.ExecutionRequest_StartRequest{
			StartRequest: &nighthawk_proto.StartRequest{
				Options: &opt,
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
		hres := &HTTPRunnerResults{}
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

type HTTPRunnerResults fhttp.HTTPRunnerResults

func (r *HTTPRunnerResults) UnmarshalJSON(data []byte) error {
	type HTTPRunnerResultsAlias HTTPRunnerResults
	rr := &struct {
		*HTTPRunnerResultsAlias
		RequestedQPS      int
		ActualDuration    interface{}
		RetCodes          interface{}
		Sizes             interface{}
		HeaderSizes       interface{}
		DurationHistogram interface{}
	}{
		HTTPRunnerResultsAlias: (*HTTPRunnerResultsAlias)(r),
	}

	if err := json.Unmarshal(data, &rr); err != nil {
		return err
	}
	r.RequestedQPS = fmt.Sprint(rr.RequestedQPS)

	updatedMap := make(map[int]int64)
	marshalledMap := rr.RetCodes.(map[string]interface{})

	for k, v := range marshalledMap {
		key, _ := strconv.Atoi(k)
		val, ok := v.(uint64)
		if ok {
			updatedMap[key] = int64(val)
		}
	}

	duration := rr.ActualDuration.(float64)
	r.ActualDuration = time.Duration(duration * 1e9)
	r.RetCodes = updatedMap
	return nil
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
	httpOpts.CACert = opts.CACert

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
