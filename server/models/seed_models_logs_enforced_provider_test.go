package models

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/spf13/viper"
)

// stubProvider is a Provider that panics on every call it is not given an
// implementation for. Embedding the interface is deliberate: the tests below
// must fail loudly if RegistryLog ever reaches back into the provider
// registration map for a system-event sink.
type stubProvider struct {
	Provider
}

// enforcedProviderKey is the registration key a pinned deployment keeps. It is
// spelled like a real remote's key (cmd/main.go registers remotes under the
// name from /capabilities, falling back to the URL host) so the fixture matches
// what issue #21584 was reported against: PROVIDER pinned to cloud.meshery.io.
const enforcedProviderKey = "cloud.meshery.io"

// newRegistryLogTestConfig builds the HandlerConfig boot-time seeding runs with,
// plus the database its system events land in.
func newRegistryLogTestConfig(t *testing.T, db *database.Handler) *HandlerConfig {
	t.Helper()

	if err := db.AutoMigrate(events.Event{}); err != nil {
		t.Fatalf("migrate events: %v", err)
	}
	return &HandlerConfig{
		Providers: map[string]Provider{
			LocalProviderName:   stubProvider{},
			enforcedProviderKey: stubProvider{},
		},
		EventBroadcaster:     NewBroadcaster("Events"),
		SystemEventPersister: &EventsPersister{DB: db},
	}
}

// persistedSystemEvents reads back every event RegistryLog persisted.
func persistedSystemEvents(t *testing.T, db *database.Handler) []events.Event {
	t.Helper()

	var persisted []events.Event
	if err := db.Find(&persisted).Error; err != nil {
		t.Fatalf("read events: %v", err)
	}
	return persisted
}

// registryLogTestLogger keeps RegistryLog's output out of the test log while
// still satisfying its logger.Handler argument.
func registryLogTestLogger(t *testing.T) logger.Handler {
	t.Helper()
	return newSeedTestLogger(t)
}

// loggedError is one record of the JSON stream logger.Handler writes to its
// error output. That stream is the operator-facing diagnostic contract: Error
// emits the MeshKit error's code alongside the message, and the code is what an
// operator greps for and what docs/data/errorref is keyed on.
type loggedError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

// errorLogSink captures that stream. logrus serializes its own writes; the
// mutex covers reads, which happen on the test goroutine.
type errorLogSink struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *errorLogSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

// records decodes every error the logger emitted.
func (s *errorLogSink) records(t *testing.T) []loggedError {
	t.Helper()

	s.mu.Lock()
	defer s.mu.Unlock()

	var records []loggedError
	dec := json.NewDecoder(bytes.NewReader(s.buf.Bytes()))
	for {
		var rec loggedError
		err := dec.Decode(&rec)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("decode captured error log stream: %v", err)
		}
		records = append(records, rec)
	}
	return records
}

// reports says whether the given MeshKit error code was emitted.
func (s *errorLogSink) reports(t *testing.T, code string) bool {
	t.Helper()

	for _, rec := range s.records(t) {
		if rec.Code == code {
			return true
		}
	}
	return false
}

// capturingTestLogger returns a logger.Handler whose error output is captured
// rather than written to stderr, so a test can assert on what was reported.
func capturingTestLogger(t *testing.T) (logger.Handler, *errorLogSink) {
	t.Helper()

	log := newSeedTestLogger(t)
	sink := &errorLogSink{}
	log.UpdateErrorLogOutput(sink)
	return log, sink
}

// TestRegistryLogUnderEnforcedProviderSeedsWithoutPanicking is the regression
// test for meshery/meshery#21584.
//
// A deployment that pins PROVIDER has every other registration - the Local
// Provider included - deleted from HandlerConfig.Providers by
// RestrictToEnforcedProvider. RegistryLog used to index that map by
// LocalProviderName and call PersistSystemEvent on the zero Provider it got
// back, which panicked in the unrecovered seeding goroutine and killed the
// server at boot before a single model was seeded.
//
// The assertion is deliberately not just "did not panic": the summary events
// must still be persisted, because dropping them would trade a crash for
// silently missing notifications on exactly the deployments that pin a
// provider.
func TestRegistryLogUnderEnforcedProviderSeedsWithoutPanicking(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t))
	hc := newRegistryLogTestConfig(t, db)

	RestrictToEnforcedProvider(hc.Providers, enforcedProviderKey)
	if _, ok := hc.Providers[LocalProviderName]; ok {
		t.Fatalf("fixture is not exercising the reported condition: Local is still registered after enforcement")
	}

	viper.Set("REGISTRY_LOG_FILE", t.TempDir()+"/registry-errors.log")
	t.Cleanup(func() { viper.Set("REGISTRY_LOG_FILE", "") })

	RegistryLog(registryLogTestLogger(t), hc, regm, NewRegistrationFailureLogHandler())

	persisted := persistedSystemEvents(t, db)
	if len(persisted) == 0 {
		t.Fatal("no system event was persisted; the registrant summaries raised while seeding were dropped")
	}
	var summaries int
	for _, e := range persisted {
		if strings.HasPrefix(e.Description, "For registrant ") {
			summaries++
		}
	}
	if summaries == 0 {
		t.Fatalf("no registrant summary event persisted; got %d event(s): %+v", len(persisted), persisted)
	}
}

// TestRegistryLogDoesNotReadTheProviderRegistry pins the invariant the fix
// establishes: system events raised while seeding are not routed through
// HandlerConfig.Providers at all, so no future change to which providers a
// deployment registers can reintroduce #21584. An empty map is the strongest
// form of "no provider is registered", and it must still seed and persist.
func TestRegistryLogDoesNotReadTheProviderRegistry(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t))
	hc := newRegistryLogTestConfig(t, db)
	hc.Providers = map[string]Provider{}

	viper.Set("REGISTRY_LOG_FILE", t.TempDir()+"/registry-errors.log")
	t.Cleanup(func() { viper.Set("REGISTRY_LOG_FILE", "") })

	RegistryLog(registryLogTestLogger(t), hc, regm, NewRegistrationFailureLogHandler())

	if len(persistedSystemEvents(t, db)) == 0 {
		t.Fatal("no system event was persisted with an empty provider registry")
	}
}

// TestRegistryLogWithoutASystemEventSinkStillLogs covers the one path that
// legitimately drops events: a HandlerConfig built without a
// SystemEventPersister. It must report the omission and keep going rather than
// panic - the boot summaries are still worth logging.
func TestRegistryLogWithoutASystemEventSinkStillLogs(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t))
	hc := newRegistryLogTestConfig(t, db)
	hc.SystemEventPersister = nil

	viper.Set("REGISTRY_LOG_FILE", t.TempDir()+"/registry-errors.log")
	t.Cleanup(func() { viper.Set("REGISTRY_LOG_FILE", "") })

	log, sink := capturingTestLogger(t)
	RegistryLog(log, hc, regm, NewRegistrationFailureLogHandler())

	if got := len(persistedSystemEvents(t, db)); got != 0 {
		t.Fatalf("expected no persisted events without a sink, got %d", got)
	}
	if !sink.reports(t, ErrNoSystemEventSinkCode) {
		t.Fatalf("dropping the events was not reported as %s, making it a silent skip; emitted: %+v", ErrNoSystemEventSinkCode, sink.records(t))
	}
}

// TestRunSeedStageRecoversPanic pins the resiliency half of #21584: a fault
// while seeding must be reported and contained, never allowed to terminate the
// server. cmd/main.go runs each seeding stage through this helper.
func TestRunSeedStageRecoversPanic(t *testing.T) {
	ran := false

	log, sink := capturingTestLogger(t)
	RunSeedStage(log, "models", func() {
		panic("boom")
	})
	RunSeedStage(log, "policies", func() {
		ran = true
	})

	if !ran {
		t.Fatal("a stage that panicked prevented a later stage from running")
	}
	if !sink.reports(t, ErrSeedingStagePanicCode) {
		t.Fatalf("the recovered fault was not reported as %s, so containing it swallowed it; emitted: %+v", ErrSeedingStagePanicCode, sink.records(t))
	}
}
