package appoptics

import (
	"strings"

	ao "github.com/appoptics/appoptics-api-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Org - struct for representing AppOptics Org
type Org struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AppOptics client instance
type AppOptics struct {
	token  string
	client *ao.Client
	spaces []*ao.Space
	charts []*ao.Chart
}

// AODash - struct used to collect token and chart ids in a dashboard. Used mainly to send data back to the UI.
type AODash struct {
	Token    string
	ChartIds []int
}

// NewAOClient - creates a new AppOptics client instance with the metadata
func NewAOClient(token, spaceName string) (*AppOptics, error) {
	logrus.Infof("AO: received token: %s", token)
	a := &AppOptics{
		token:  token,
		client: ao.NewClient(token),
	}
	logrus.Info("retrieving spaces")
	if err := a.loadAllRelevantSpaces(spaceName); err != nil {
		return nil, err
	}
	logrus.Info("retrieving charts")
	if err := a.loadAllRelevantCharts(); err != nil {
		return nil, err
	}
	logrus.Infof("created AO client: %+#v", a)
	return a, nil
}

func (a *AppOptics) loadAllRelevantSpaces(spaceName string) error {
	ss := a.client.SpacesService()
	spaces, err := ss.List(nil)
	if err != nil {
		err = errors.Wrap(err, "unable to get all spaces")
		logrus.Error(err)
		return err
	}
	a.spaces = []*ao.Space{}
	for _, sp := range spaces {
		if strings.Contains(strings.ToLower(sp.Name), spaceName) {
			a.spaces = append(a.spaces, sp)
		}
	}
	return nil
}

func (a *AppOptics) loadAllRelevantCharts() error {
	var err error
	for _, sp := range a.spaces {
		a.charts, err = a.client.ChartsService().List(sp.ID)
		if err != nil {
			err = errors.Wrapf(err, "unable to get all charts for space: %s (%d)", sp.Name, sp.ID)
			logrus.Error(err)
			return err
		}
	}
	return nil
}

// GenerateDataForTemplate - used for injecting data into struct
func (a *AppOptics) GenerateDataForTemplate() *AODash {
	if a != nil {
		d := &AODash{
			Token: a.token,
		}
		chartIds := []int{}
		if a.charts != nil && len(a.charts) > 0 {
			for _, ch := range a.charts {
				chartIds = append(chartIds, *ch.ID)
			}
		}
		d.ChartIds = chartIds
		return d
	}
	return nil
}
