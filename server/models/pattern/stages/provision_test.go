package stages

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	patterncore "github.com/meshery/meshery/server/models/pattern/core"
	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta2/component"
	componentv1beta3 "github.com/meshery/schemas/models/v1beta3/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

type stubInfoProvider struct{ isDelete bool }

func (s *stubInfoProvider) GetMesheryPatternResource(_, _, _, _ string) (*core.Uuid, error) {
	return nil, nil
}

func (s *stubInfoProvider) IsDelete() bool { return s.isDelete }

// stubActionProvider stands in for the deployment fulfillment paths. Components
// named in applyFailures come back reported as failed (the shape both
// fulfillment paths use for a component that could not be applied), and
// components named in dispatchFailures come back as an error (the shape used
// when Meshery could not dispatch the component at all).
//
// Components named in noOutcomes come back as neither - no outcome and no
// error. No fulfillment path does that today, which is precisely why the stage
// has to be asked to demonstrate what it does with it.
//
// A component is acted on once per cluster, so an outcome is reported for each
// of clusters. clusterFailures names the single cluster a component fails on,
// which is how a deployment to several clusters at once partly fails.
type stubActionProvider struct {
	registry *registry.RegistryManager

	clusters             []string
	applyFailures        map[string]bool
	clusterFailures      map[string]string
	dispatchFailures     map[string]bool
	noOutcomes           map[string]bool
	prerequisiteFailures map[string]bool

	mx         sync.Mutex
	dispatched []string
	terminated error
}

func (s *stubActionProvider) Terminate(err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.terminated = err
}

func (s *stubActionProvider) Log(_ string) {}

func (s *stubActionProvider) GetRegistry() *registry.RegistryManager { return s.registry }

func (s *stubActionProvider) DryRun(_ []*component.ComponentDefinition) (map[string]map[string]patterncore.DryRunResponseWrapper, error) {
	return nil, nil
}

func (s *stubActionProvider) Mutate(_ *pattern.PatternFile) {}

func (s *stubActionProvider) Provision(ccp CompConfigPair) ([]patterns.DeploymentMessagePerContext, error) {
	name := ccp.Component.DisplayName

	s.mx.Lock()
	s.dispatched = append(s.dispatched, name)
	s.mx.Unlock()

	if s.dispatchFailures[name] {
		return nil, errors.New("meshery-server-test", errors.Alert, []string{"dispatch failed"}, []string{"dispatch failed"}, []string{}, []string{})
	}

	if s.noOutcomes[name] {
		return nil, nil
	}

	clusters := s.clusters
	if len(clusters) == 0 {
		clusters = []string{"test-cluster"}
	}

	msgs := []patterns.DeploymentMessagePerContext{}

	for _, cluster := range clusters {
		summary := []patterns.DeploymentMessagePerComp{}

		// patterns.Process reports a component's prerequisites on their own entry,
		// appended to the same Summary as the component's own apply entry and ahead
		// of it. Installing them is fail-forward, so the component is applied either
		// way and both entries are reported for it.
		if s.prerequisiteFailures[name] {
			summary = append(summary, patterns.DeploymentMessagePerComp{
				CompName:       name,
				Success:        false,
				IsPrerequisite: true,
				Message:        fmt.Sprintf("could not install what %s needs", name),
				Error:          errors.New("meshery-server-test", errors.Alert, []string{"prerequisite failed"}, []string{"prerequisite failed"}, []string{}, []string{}),
			})
		}

		applied := !s.applyFailures[name] && s.clusterFailures[name] != cluster

		summary = append(summary, patterns.DeploymentMessagePerComp{
			CompName: name,
			Success:  applied,
			Message:  fmt.Sprintf("deployed %s on %s", name, cluster),
		})

		msgs = append(msgs, patterns.DeploymentMessagePerContext{
			SystemName: cluster,
			Summary:    summary,
		})
	}

	return msgs, nil
}

func (s *stubActionProvider) order() []string {
	s.mx.Lock()
	defer s.mx.Unlock()

	return append([]string(nil), s.dispatched...)
}

func newStubActionProvider(t *testing.T) *stubActionProvider {
	t.Helper()

	// A file-backed database rather than ":memory:" so that the schema created
	// by NewRegistryManager is visible to every pooled connection. The registry
	// holds no rows, so every component resolves to an empty registrant, which
	// is what routes it to the stub below.
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "registry.db")})
	if err != nil {
		t.Fatalf("failed to open a database: %v", err)
	}

	reg, err := registry.NewRegistryManager(&db)
	if err != nil {
		t.Fatalf("failed to build a registry manager: %v", err)
	}

	return &stubActionProvider{
		registry:             reg,
		applyFailures:        map[string]bool{},
		clusterFailures:      map[string]string{},
		dispatchFailures:     map[string]bool{},
		noOutcomes:           map[string]bool{},
		prerequisiteFailures: map[string]bool{},
	}
}

func newDesignComponent(displayName string, dependsOn ...string) *component.ComponentDefinition {
	comp := &component.ComponentDefinition{
		ID:          uuid.Must(uuid.NewV4()),
		DisplayName: displayName,
	}

	if len(dependsOn) > 0 {
		deps := make([]interface{}, 0, len(dependsOn))
		for _, dep := range dependsOn {
			deps = append(deps, dep)
		}

		comp.Metadata.AdditionalProperties = map[string]interface{}{"dependsOn": deps}
	}

	return comp
}

// runProvision runs the provision stage over a design and returns the outcome
// reported for each component, keyed by the component's display name.
func runProvision(t *testing.T, act *stubActionProvider, components ...*component.ComponentDefinition) map[string][]patterns.DeploymentMessagePerContext {
	t.Helper()

	log, err := logger.New("provision-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	design := pattern.PatternFile{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       "provision-test",
		Version:    "1.0.0",
		Components: components,
	}

	names := map[string]string{}
	mapping := map[core.Uuid]component.ComponentDefinition{}
	for _, comp := range components {
		names[comp.ID.String()] = comp.DisplayName
		mapping[comp.ID] = *comp
	}

	data := &Data{
		Pattern:                       &design,
		Other:                         map[string]interface{}{},
		DeclartionToDefinitionMapping: mapping,
	}

	Provision(&stubInfoProvider{}, act, log)(data, nil, nil)

	reported := map[string][]patterns.DeploymentMessagePerContext{}

	data.Lock.Lock()
	defer data.Lock.Unlock()

	for key, value := range data.Other {
		id := key[:len(key)-len(ProvisionSuffixKey)]

		msgs, ok := value.([]patterns.DeploymentMessagePerContext)
		if !ok {
			t.Fatalf("data.Other[%q] holds %T, want []patterns.DeploymentMessagePerContext", key, value)
		}

		reported[names[id]] = msgs
	}

	return reported
}

func succeeded(msgs []patterns.DeploymentMessagePerContext) bool {
	for _, msg := range msgs {
		for _, summary := range msg.Summary {
			if !summary.Success {
				return false
			}
		}
	}

	return len(msgs) > 0
}

func dispatched(order []string, name string) bool {
	for _, n := range order {
		if n == name {
			return true
		}
	}

	return false
}

func TestProvisionWithholdsDependentsOfAComponentThatFailedToApply(t *testing.T) {
	act := newStubActionProvider(t)
	act.applyFailures["a"] = true

	reported := runProvision(t, act,
		newDesignComponent("a"),
		newDesignComponent("b", "a"),
		newDesignComponent("c"),
	)

	order := act.order()

	if !dispatched(order, "a") {
		t.Errorf("dispatched %v, want \"a\" dispatched", order)
	}

	if dispatched(order, "b") {
		t.Errorf("dispatched %v, want \"b\" withheld because \"a\" failed to apply", order)
	}

	if !dispatched(order, "c") {
		t.Errorf("dispatched %v, want \"c\" dispatched: it declares no dependency on \"a\"", order)
	}

	if succeeded(reported["a"]) {
		t.Error("\"a\" was reported as successful, want it reported as failed")
	}

	withheld, reportedAtAll := reported["b"]
	if !reportedAtAll {
		t.Fatal("\"b\" produced no deployment message, want it reported as withheld")
	}

	if succeeded(withheld) {
		t.Error("\"b\" was reported as successful, want it reported as withheld")
	}

	if got := errors.GetCode(withheld[0].Summary[0].Error); got != ErrDependencyNotSatisfiedCode {
		t.Errorf("\"b\" was reported with error code %s, want %s", got, ErrDependencyNotSatisfiedCode)
	}

	if !succeeded(reported["c"]) {
		t.Error("\"c\" was not reported as successful, want it deployed regardless of \"a\" failing")
	}
}

// A component deployed to several clusters at once that fails on one of them is
// treated as failed outright, so its dependents are withheld everywhere rather
// than only on the cluster where it failed.
//
// This is the conservative choice the deployment-engine documentation describes:
// a dependent is held back on the clusters where its dependency did come up, in
// exchange for never applying it against a cluster where its dependency is
// absent.
func TestProvisionWithholdsDependentsWhenAComponentFailedOnOneOfSeveralClusters(t *testing.T) {
	act := newStubActionProvider(t)
	act.clusters = []string{"cluster-a", "cluster-b"}
	act.clusterFailures["a"] = "cluster-b"

	reported := runProvision(t, act,
		newDesignComponent("a"),
		newDesignComponent("b", "a"),
	)

	if order := act.order(); dispatched(order, "b") {
		t.Errorf("dispatched %v, want \"b\" withheld on every cluster: \"a\" failed on one of them", order)
	}

	if _, ok := reported["b"]; !ok {
		t.Error("\"b\" produced no deployment message, want it reported as withheld")
	}

	var succeededOn, failedOn []string
	for _, msg := range reported["a"] {
		for _, entry := range msg.Summary {
			if entry.Success {
				succeededOn = append(succeededOn, msg.SystemName)
				continue
			}

			failedOn = append(failedOn, msg.SystemName)
		}
	}

	if len(succeededOn) != 1 || succeededOn[0] != "cluster-a" {
		t.Errorf("\"a\" succeeded on %v, want it reported as deployed on cluster-a", succeededOn)
	}

	if len(failedOn) != 1 || failedOn[0] != "cluster-b" {
		t.Errorf("\"a\" failed on %v, want it reported as failed on cluster-b", failedOn)
	}
}

// Installing what a component needs is deliberately fail-forward - what it
// needs may already be present in the cluster - and the component is applied
// afterwards regardless. A component that goes on to apply successfully has
// therefore deployed, and must not withhold anything that depends on it.
func TestProvisionDeploysDependentsWhenOnlyAPrerequisiteFailed(t *testing.T) {
	act := newStubActionProvider(t)
	act.prerequisiteFailures["a"] = true

	reported := runProvision(t, act,
		newDesignComponent("a"),
		newDesignComponent("b", "a"),
	)

	if order := act.order(); !dispatched(order, "b") {
		t.Errorf("dispatched %v, want \"b\" deployed: \"a\" itself applied successfully", order)
	}

	summary, ok := reported["a"]
	if !ok {
		t.Fatal("\"a\" produced no deployment message, want the prerequisite failure reported against it")
	}

	var warned, applied bool
	for _, msg := range summary {
		for _, entry := range msg.Summary {
			if entry.IsPrerequisite && !entry.Success && entry.Error != nil {
				warned = true
			}

			if !entry.IsPrerequisite && entry.Success {
				applied = true
			}
		}
	}

	if !warned {
		t.Error("\"a\" did not carry the prerequisite failure, want it reported as a warning against the component")
	}

	if !applied {
		t.Error("\"a\" did not report its own successful apply")
	}
}

// A component that fails before it can be handed to a fulfillment path is still
// reported in its own right.
//
// Withholding names the dependency that caused it, so a component that fails
// and is never recorded still appears in the summary - but only inside its
// dependents' messages, as the thing that held them back, with no outcome of
// its own anywhere.
func TestProvisionReportsAComponentThatCouldNotBeEnriched(t *testing.T) {
	act := newStubActionProvider(t)

	original := enrichComponent
	t.Cleanup(func() { enrichComponent = original })
	enrichComponent = func(comp *componentv1beta3.ComponentDefinition, designID, designVersion string) error {
		if comp.DisplayName == "a" {
			return errors.New("meshery-server-test", errors.Alert, []string{"enrichment failed"}, []string{"enrichment failed"}, []string{}, []string{})
		}

		return original(comp, designID, designVersion)
	}

	reported := runProvision(t, act,
		newDesignComponent("a"),
		newDesignComponent("b", "a"),
	)

	if order := act.order(); dispatched(order, "a") {
		t.Errorf("dispatched %v, want \"a\" never dispatched: it could not be enriched", order)
	}

	summary, ok := reported["a"]
	if !ok {
		t.Fatal("\"a\" produced no deployment message of its own, want the failure reported against it")
	}

	if succeeded(summary) {
		t.Error("\"a\" was reported as successful, want it reported as failed")
	}

	if _, ok := reported["b"]; !ok {
		t.Error("\"b\" produced no deployment message, want it reported as withheld")
	}
}

// A component that could not be handed to a fulfillment path at all still
// withholds its dependents, and is still named in the summary rather than
// vanishing from it.
func TestProvisionReportsAComponentThatCouldNotBeDispatched(t *testing.T) {
	act := newStubActionProvider(t)
	act.dispatchFailures["a"] = true

	reported := runProvision(t, act,
		newDesignComponent("a"),
		newDesignComponent("b", "a"),
	)

	summary, ok := reported["a"]
	if !ok {
		t.Fatal("\"a\" produced no deployment message, want the dispatch failure reported against it")
	}

	if succeeded(summary) {
		t.Error("\"a\" was reported as successful, want it reported as failed")
	}
}

func TestProvisionWithholdsDependentsOfAComponentThatCouldNotBeDispatched(t *testing.T) {
	act := newStubActionProvider(t)
	act.dispatchFailures["a"] = true

	reported := runProvision(t, act,
		newDesignComponent("a"),
		newDesignComponent("b", "a"),
	)

	if order := act.order(); dispatched(order, "b") {
		t.Errorf("dispatched %v, want \"b\" withheld because \"a\" could not be dispatched", order)
	}

	if _, ok := reported["b"]; !ok {
		t.Error("\"b\" produced no deployment message, want it reported as withheld")
	}
}

// A component a fulfillment path said nothing about is reported rather than
// left out of the summary, and it is not taken for a success: an empty outcome
// reads exactly like a successful one, so treating it as success would deploy
// the components that depend on it on a success nothing ever reported.
func TestProvisionReportsAComponentThatReportedNoOutcome(t *testing.T) {
	act := newStubActionProvider(t)
	act.noOutcomes["a"] = true

	reported := runProvision(t, act,
		newDesignComponent("a"),
		newDesignComponent("b", "a"),
	)

	summary, ok := reported["a"]
	if !ok {
		t.Fatal("\"a\" produced no deployment message, want the missing outcome reported against it")
	}

	if succeeded(summary) {
		t.Error("\"a\" was reported as successful, want it reported as failed: nothing was reported for it")
	}

	// succeeded() is false for an empty message too, so the summary is checked
	// for naming the component rather than merely existing.
	named := false
	for _, msg := range summary {
		for _, entry := range msg.Summary {
			if entry.CompName == "a" {
				named = true
			}
		}
	}

	if !named {
		t.Error("\"a\" is not named in the outcome reported for it, want it named rather than missing from the summary")
	}

	if order := act.order(); dispatched(order, "b") {
		t.Errorf("dispatched %v, want \"b\" withheld: nothing was reported for \"a\"", order)
	}

	if _, ok := reported["b"]; !ok {
		t.Error("\"b\" produced no deployment message, want it reported as withheld")
	}
}

func TestProvisionDeploysEveryComponentWhenNoneFails(t *testing.T) {
	act := newStubActionProvider(t)

	reported := runProvision(t, act,
		newDesignComponent("a"),
		newDesignComponent("b", "a"),
	)

	order := act.order()
	if len(order) != 2 {
		t.Fatalf("dispatched %v, want both components dispatched", order)
	}

	if order[0] != "a" {
		t.Errorf("dispatched %v, want \"a\" first: \"b\" depends on it", order)
	}

	for _, name := range []string{"a", "b"} {
		if !succeeded(reported[name]) {
			t.Errorf("%q was not reported as successful", name)
		}
	}
}

func TestProvisionTerminatesOnAnUnresolvableDependency(t *testing.T) {
	act := newStubActionProvider(t)

	runProvision(t, act, newDesignComponent("a"), newDesignComponent("b", "nowhere"))

	if order := act.order(); len(order) != 0 {
		t.Errorf("dispatched %v, want nothing deployed when the plan cannot be built", order)
	}

	if act.terminated == nil {
		t.Fatal("the stage did not terminate, want a structured error for the unresolvable dependency")
	}
}

func TestProvisionTerminatesOnCyclicDependencies(t *testing.T) {
	act := newStubActionProvider(t)

	runProvision(t, act, newDesignComponent("a", "b"), newDesignComponent("b", "a"))

	if order := act.order(); len(order) != 0 {
		t.Errorf("dispatched %v, want nothing deployed for a cyclic design", order)
	}

	if got := errors.GetCode(act.terminated); got != "meshery-server-1479" {
		t.Errorf("the stage terminated with error code %s, want the cyclic-dependency code", got)
	}
}

func TestProvisionNeverDispatchesAnnotations(t *testing.T) {
	act := newStubActionProvider(t)

	note := newDesignComponent("a note")
	note.Metadata.IsAnnotation = true

	runProvision(t, act, note, newDesignComponent("b", "a note"))

	order := act.order()
	if dispatched(order, "a note") {
		t.Errorf("dispatched %v, want the annotation never dispatched", order)
	}

	if !dispatched(order, "b") {
		t.Errorf("dispatched %v, want \"b\" dispatched: a dependency on an annotation is always satisfied", order)
	}
}
