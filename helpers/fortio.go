package helpers

import (
	"encoding/json"
	"os"
	"strings"

	"fortio.org/fortio/fgrpc"
	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/periodic"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// SharedHTTPOptions is the flag->httpoptions transfer code shared between
// fortio_main and fcurl.
func sharedHTTPOptions(opts *models.LoadTestOptions) *fhttp.HTTPOptions {
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
	// httpOpts.Payload = fnet.GeneratePayload(*PayloadFileFlag, *PayloadSizeFlag, *PayloadFlag)
	// httpOpts.UnixDomainSocket = *unixDomainSocketFlag
	if false { // *followRedirectsFlag {
		httpOpts.FollowRedirects = true
		httpOpts.DisableFastClient = true
	}
	return &httpOpts
}

// FortioLoadTest is the actual code which invokes Fortio to run the load test
func FortioLoadTest(opts *models.LoadTestOptions) (map[string]interface{}, error) {
	defaults := &periodic.DefaultRunnerOptions
	// httpOpts := bincommon.SharedHTTPOptions()
	httpOpts := sharedHTTPOptions(opts)
	if opts.IsInsecure {
		httpOpts.Insecure = true
	}
	// if justCurl {
	// 	bincommon.FetchURL(httpOpts)
	// 	return
	// }
	url := httpOpts.URL
	out := os.Stdout
	qps := opts.HTTPQPS // TODO possibly use translated <=0 to "max" from results/options normalization in periodic/
	// _, _ = fmt.Fprintf(out, "Fortio %s running at %g queries per second, %d->%d procs",
	// 	version.Short(), qps, prevGoMaxProcs, runtime.GOMAXPROCS(0))
	// if *exactlyFlag > 0 {
	// 	_, _ = fmt.Fprintf(out, ", for %d calls: %s\n", *exactlyFlag, url)
	// } else {
	// 	if *durationFlag <= 0 {
	// 		// Infinite mode is determined by having a negative duration value
	// 		*durationFlag = -1
	// 		_, _ = fmt.Fprintf(out, ", until interrupted: %s\n", url)
	// 	} else {
	// 		_, _ = fmt.Fprintf(out, ", for %v: %s\n", *durationFlag, url)
	// 	}
	// }
	if qps <= 0 {
		qps = -1 // 0==unitialized struct == default duration, -1 (0 for flag) is max
	}
	labels := opts.Name + " - " + url
	// if labels == "" {
	// 	// hname, _ := os.Hostname()
	// 	shortURL := url
	// 	for _, p := range []string{"https://", "http://"} {
	// 		if strings.HasPrefix(url, p) {
	// 			shortURL = url[len(p):]
	// 			break
	// 		}
	// 	}
	// 	labels = shortURL
	// }
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
	var err error
	if opts.IsGRPC {
		o := fgrpc.GRPCRunnerOptions{
			RunnerOptions:      ro,
			Destination:        url,
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
		return nil, err
	}
	logrus.Debugf("original version of the test: %+#v", res)
	resultsMap := map[string]interface{}{}
	var bd []byte
	if opts.IsGRPC {
		gres, _ := res.(*fgrpc.GRPCRunnerResults)
		bd, err = json.Marshal(gres)
	} else {
		hres, _ := res.(*fhttp.HTTPRunnerResults)
		bd, err = json.Marshal(hres)
	}
	if err != nil {
		err = errors.Wrap(err, "error while marshalling results")
		logrus.Error(err)
		return nil, err
	}
	err = json.Unmarshal(bd, &resultsMap)
	if err != nil {
		err = errors.Wrap(err, "error while unmarshaling data to map")
		logrus.Error(err)
		return nil, err
	}
	logrus.Debugf("Mapped version of the test: %+#v", resultsMap)
	return resultsMap, nil
}
