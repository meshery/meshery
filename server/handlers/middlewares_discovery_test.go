package handlers

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

// recordingEventsProvider records the events persisted through it.
type recordingEventsProvider struct {
	models.Provider
	persisted []events.Event
}

func (p *recordingEventsProvider) PersistEvent(event events.Event, _ string) error {
	p.persisted = append(p.persisted, event)
	return nil
}

// A background Discovery reports its failure and the resolution of the
// connection's unreachable condition, and nothing for a routine success.
func TestReportBackgroundDiscovery(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
		log:    newTestLogger(t),
	}
	user := uuid.Must(uuid.NewV4())
	key := models.ConditionKey("k8s-connection", uuid.Must(uuid.NewV4()).String(), "unreachable")

	routine := events.NewEvent().WithSeverity(events.Informational).WithDescription("kubernetes connection changed to connected").Build()
	failure := events.NewEvent().WithSeverity(events.Error).WithCorrelationKey(key).Build()
	resolution := events.NewEvent().WithSeverity(events.Success).Resolves(key, time.Now()).Build()

	cases := []struct {
		name        string
		event       *events.Event
		err         error
		wantPersist bool
	}{
		{"routine success", routine, nil, false},
		{"suppressed repeat", nil, nil, false},
		{"failure", failure, errors.New("unreachable"), true},
		{"resolution", resolution, nil, true},
	}
	for _, c := range cases {
		provider := &recordingEventsProvider{}
		h.reportBackgroundDiscovery(provider, "token", user, c.event, c.err)
		if got := len(provider.persisted) == 1; got != c.wantPersist {
			t.Errorf("%s: persisted %d events, want persisted=%v", c.name, len(provider.persisted), c.wantPersist)
		}
	}
}
