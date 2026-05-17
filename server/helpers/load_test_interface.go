package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"fortio.org/fortio/fgrpc"
	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/periodic"
	"github.com/layer5io/gowrk2/api"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
)

// FortioLoadTest is the actual code which invokes Fortio to run the load test
func FortioLoadTest(opts *models.LoadTestOptions, log logger.Handler) (map[string]interface{}, *periodic.RunnerResults, error) {
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
			RunnerOptions: ro,
			TLSOptions:    fhttp.TLSOptions{},
			Destination:   rURL,
			Service:       opts.GRPCHealthSvc,
			// Profiler:           "",
			Payload:            string(httpOpts.Payload),
			Streams:            opts.GRPCStreamsCount,
			Delay:              opts.GRPCPingDelay,
			CertOverride:       opts.CACert,
			AllowInitialErrors: opts.AllowInitialErrors,
			UsePing:            opts.GRPCDoPing,
			// Metadata:           map[string][]string{},
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

		log.Debug("options string: ", opts.Options)
		if opts.Options != "" {
			log.Debug(fmt.Sprintf("Fortio config: %+#v", o))
			err := json.Unmarshal([]byte(opts.Options), &o)
			if err != nil {
				return nil, nil, models.ErrUnmarshal(err, "options string")
			}
			log.Debug(fmt.Sprintf("Fortio config with options: %+#v", o))
		}
		res, err = fhttp.RunHTTPTest(&o)
	}
	if err != nil {
		return nil, nil, ErrRunningTest(err)
	}
	log.Debug(fmt.Sprintf("original version of the test: %+#v", res))

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
		return nil, nil, models.ErrUnmarshal(err, "data to map")
	}
	log.Debug(fmt.Sprintf("Mapped version of the test: %+#v", resultsMap))
	return resultsMap, result, nil
}

// WRK2LoadTest is the actual code which invokes Wrk2 to run the load test
func WRK2LoadTest(opts *models.LoadTestOptions, log logger.Handler) (map[string]interface{}, *periodic.RunnerResults, error) {
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

	log.Debug("options string: ", opts.Options)
	if opts.Options != "" {
		log.Debug(fmt.Sprintf("GoWrk2 config: %+#v", ro))
		err := json.Unmarshal([]byte(opts.Options), &ro)
		if err != nil {
			return nil, nil, models.ErrUnmarshal(err, "options string")
		}
		log.Debug(fmt.Sprintf("GoWrk2 config with options: %+#v", ro))
	}

	var res periodic.HasRunnerResult
	var err error
	if opts.SupportedLoadTestMethods == 2 {
		return nil, nil, ErrGrpcSupport(err, "Wrk2")
	}
	var gres *api.GoWRK2
	gres, err = api.WRKRun(ro)
	if err == nil {
		log.Debug(fmt.Sprintf("WRK Result: %+v", gres))
		res, err = api.TransformWRKToFortio(gres, ro)
	}

	if err != nil {
		return nil, nil, ErrRunningTest(err)
	}
	log.Debug(fmt.Sprintf("original version of the test: %+#v", res))

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
		return nil, nil, models.ErrUnmarshal(err, "data to map")
	}
	log.Debug(fmt.Sprintf("Mapped version of the test: %+#v", resultsMap))
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
