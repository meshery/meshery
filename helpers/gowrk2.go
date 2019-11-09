package helpers

import (
	"encoding/json"
	"strings"

	"fortio.org/fortio/fgrpc"
	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/periodic"
	"github.com/layer5io/gowrk2/api"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// FortioLoadTest is the actual code which invokes Fortio to run the load test
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
		// Connection:        10,
		RQPS:        qps,
		URL:         rURL,
		Labels:      labels,
		Percentiles: []float64{50, 75, 90, 99, 99.99, 99.999},
	}
	var res periodic.HasRunnerResult
	var err error
	if opts.IsGRPC {
		err := errors.New("wrk2 does not support gRPC at the moment")
		logrus.Error(err)
		return nil, nil, err
	} else {
		var gres *api.GoWRK2
		gres, err = api.WRKRun(ro)
		if err == nil {
			logrus.Debugf("WRK Result: %+v", gres)
			res, err = api.TransformWRKToFortio(gres, ro)
		}
	}
	if err != nil {
		err = errors.Wrap(err, "error while running tests")
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Debugf("original version of the test: %+#v", res)

	var result *periodic.RunnerResults
	var bd []byte
	if opts.IsGRPC {
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
