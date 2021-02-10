package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"fortio.org/fortio/fgrpc"
	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/periodic"
	"github.com/layer5io/gowrk2/api"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/nighthawk-go/apinighthawk"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

// NighthawkLoadTest is the actual code which invokes nighthawk to run the load test
func NighthawkLoadTest(opts *models.LoadTestOptions) (map[string]interface{}, *periodic.RunnerResults, error) {
	qps := opts.HTTPQPS

	if qps <= 0 {
		qps = -1 // 0==unitialized struct == default duration, -1 (0 for flag) is max
	}

	rURL := strings.TrimLeft(opts.URL, " \t\r\n")

	ro := &apinighthawk.NighthawkConfig{
		DurationInSeconds: opts.Duration.Seconds(),
		Thread:            opts.HTTPNumThreads,
		QPS:               qps,
		URL:               rURL,
	}

	var err error

	if opts.SupportedLoadTestMethods == 2 {
		err := errors.New("nighthawk does not support gRPC load testing")
		logrus.Error(err)
		return nil, nil, err
	}

	res1, err := apinighthawk.NighthawkRun(ro)
	if err != nil {
		err = errors.Wrap(err, "error while running tests")
		logrus.Error(err)
		return nil, nil, err
	}

	res := string(res1)

	logrus.Debugf("original version of the test: %+#v", res)

	var result *periodic.RunnerResults = &periodic.RunnerResults{}

	err = json.Unmarshal([]byte(res1), result)

	if err != nil {
		err = errors.Wrap(err, "Error while unmarshaling  Nighthawk results to the FortioHTTPRunner")
		logrus.Error(err)
		return nil, nil, err
	}

	resultsMap := map[string]interface{}{}
	err = json.Unmarshal(res1, &resultsMap)
	if err != nil {
		err = errors.Wrap(err, "Error while unmarshaling Nighthawk results to map")
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
