package planner

import (
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

func testLogger(t *testing.T) logger.Handler {
	t.Helper()

	log, err := logger.New("planner-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	return log
}

// newComponent builds a design component named displayName that declares a
// dependency on each of dependsOn. dependsOn is passed in the []interface{}
// shape a design carries once it has been round-tripped over the wire.
func newComponent(displayName string, dependsOn ...string) *component.ComponentDefinition {
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

func newDesign(name string, components ...*component.ComponentDefinition) pattern.PatternFile {
	return pattern.PatternFile{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       name,
		Components: components,
	}
}

// execution records the order in which components were dispatched and which
// ones were withheld, as the provision stage observes them.
type execution struct {
	mx         sync.Mutex
	dispatched []string
	withheld   map[string]string
}

func newExecution() *execution {
	return &execution{withheld: map[string]string{}}
}

// run executes the plan, failing the components named in fails.
func (e *execution) run(t *testing.T, plan *Plan, fails map[string]bool) {
	t.Helper()

	names := map[string]string{}
	for _, comp := range plan.Data.Components {
		names[comp.ID.String()] = comp.DisplayName
	}

	done := make(chan error, 1)

	go func() {
		done <- plan.Execute(
			func(id string, _ component.ComponentDefinition) bool {
				e.mx.Lock()
				e.dispatched = append(e.dispatched, names[id])
				e.mx.Unlock()

				return !fails[names[id]]
			},
			func(id string, _ component.ComponentDefinition, failedDependency string) {
				e.mx.Lock()
				e.withheld[names[id]] = names[failedDependency]
				e.mx.Unlock()
			},
			testLogger(t),
		)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Plan.Execute() returned an unexpected error: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Plan.Execute() did not return; the traversal is stuck waiting on a dependency signal")
	}
}

func (e *execution) order() []string {
	e.mx.Lock()
	defer e.mx.Unlock()

	return append([]string(nil), e.dispatched...)
}

func indexOf(names []string, name string) int {
	for i, n := range names {
		if n == name {
			return i
		}
	}

	return -1
}

func errorCode(t *testing.T, err error) string {
	t.Helper()

	if err == nil {
		t.Fatal("expected an error, got none")
	}

	return errors.GetCode(err)
}

func TestCreatePlanOrdersDeclaredDependencies(t *testing.T) {
	a := newComponent("a")
	b := newComponent("b", "a")
	design := newDesign("ordered", a, b)

	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	feasible, err := plan.IsFeasible()
	if err != nil {
		t.Fatalf("Plan.IsFeasible() returned an unexpected error: %v", err)
	}

	if !feasible {
		t.Fatal("Plan.IsFeasible() = false, want true")
	}

	run := newExecution()
	run.run(t, plan, nil)

	order := run.order()
	if len(order) != 2 {
		t.Fatalf("dispatched %v, want both components dispatched", order)
	}

	if indexOf(order, "a") > indexOf(order, "b") {
		t.Errorf("dispatched %v, want \"a\" before \"b\"", order)
	}
}

func TestCreatePlanInvertsOrderForUndeploy(t *testing.T) {
	a := newComponent("a")
	b := newComponent("b", "a")
	design := newDesign("inverted", a, b)

	plan, err := CreatePlan(design, nil, true)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	run := newExecution()
	run.run(t, plan, nil)

	order := run.order()
	if indexOf(order, "b") > indexOf(order, "a") {
		t.Errorf("dispatched %v, want \"b\" before \"a\" when the plan is inverted", order)
	}
}

func TestCreatePlanRejectsUnresolvableDependency(t *testing.T) {
	design := newDesign("unresolvable", newComponent("a"), newComponent("b", "nowhere"))

	plan, err := CreatePlan(design, nil, false)
	if plan != nil {
		t.Error("CreatePlan() returned a plan for an unresolvable dependency, want none")
	}

	if got := errorCode(t, err); got != ErrUnresolvedDependencyCode {
		t.Errorf("CreatePlan() error code = %s, want %s", got, ErrUnresolvedDependencyCode)
	}
}

func TestCreatePlanRejectsAmbiguousDependency(t *testing.T) {
	design := newDesign("ambiguous", newComponent("a"), newComponent("a"), newComponent("b", "a"))

	plan, err := CreatePlan(design, nil, false)
	if plan != nil {
		t.Error("CreatePlan() returned a plan for an ambiguous dependency, want none")
	}

	if got := errorCode(t, err); got != ErrAmbiguousDependencyCode {
		t.Errorf("CreatePlan() error code = %s, want %s", got, ErrAmbiguousDependencyCode)
	}
}

func TestCreatePlanTreatsDependencyOnAnnotationAsSatisfied(t *testing.T) {
	note := newComponent("a note")
	note.Metadata.IsAnnotation = true

	b := newComponent("b", "a note")
	design := newDesign("annotated", b)

	plan, err := CreatePlan(design, []*component.ComponentDefinition{note}, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	if edges := len(plan.Edges); edges != 0 {
		t.Errorf("plan carries %d dependency edges, want 0 for a dependency on an annotation", edges)
	}

	run := newExecution()
	run.run(t, plan, nil)

	if order := run.order(); len(order) != 1 || order[0] != "b" {
		t.Errorf("dispatched %v, want [b]", order)
	}
}

func TestCreatePlanRejectsMalformedDependencyDeclaration(t *testing.T) {
	b := newComponent("b")
	b.Metadata.AdditionalProperties = map[string]interface{}{"dependsOn": "a"}

	design := newDesign("malformed", newComponent("a"), b)

	if _, err := CreatePlan(design, nil, false); errorCode(t, err) != ErrInvalidDependencyCode {
		t.Errorf("CreatePlan() error code = %s, want %s", errorCode(t, err), ErrInvalidDependencyCode)
	}
}

func TestCreatePlanAcceptsAbsentDependencyDeclaration(t *testing.T) {
	a := newComponent("a")
	a.Metadata.AdditionalProperties = map[string]interface{}{"dependsOn": nil}

	design := newDesign("no-dependencies", a)

	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error for an empty \"dependsOn\": %v", err)
	}

	run := newExecution()
	run.run(t, plan, nil)

	if order := run.order(); len(order) != 1 || order[0] != "a" {
		t.Errorf("dispatched %v, want [a]", order)
	}
}

func TestExecuteWithholdsDependentsOfAFailedComponent(t *testing.T) {
	a := newComponent("a")
	b := newComponent("b", "a")
	c := newComponent("c")
	design := newDesign("failing", a, b, c)

	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	run := newExecution()
	run.run(t, plan, map[string]bool{"a": true})

	order := run.order()
	if indexOf(order, "b") != -1 {
		t.Errorf("dispatched %v, want \"b\" withheld because \"a\" failed", order)
	}

	if indexOf(order, "a") == -1 {
		t.Errorf("dispatched %v, want \"a\" dispatched", order)
	}

	// A component that declared no dependency on the failed one is unaffected.
	if indexOf(order, "c") == -1 {
		t.Errorf("dispatched %v, want \"c\" dispatched regardless of \"a\" failing", order)
	}

	if got := run.withheld["b"]; got != "a" {
		t.Errorf("withheld[\"b\"] = %q, want \"a\"", got)
	}

	if _, reported := run.withheld["c"]; reported {
		t.Error("\"c\" was reported as withheld, want it deployed")
	}
}

func TestExecutePropagatesWithholdingAlongTheChain(t *testing.T) {
	a := newComponent("a")
	b := newComponent("b", "a")
	c := newComponent("c", "b")
	design := newDesign("chain", a, b, c)

	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	run := newExecution()
	run.run(t, plan, map[string]bool{"a": true})

	if order := run.order(); len(order) != 1 || order[0] != "a" {
		t.Errorf("dispatched %v, want only [a]", order)
	}

	if got := run.withheld["c"]; got != "b" {
		t.Errorf("withheld[\"c\"] = %q, want \"b\"", got)
	}
}

// A node with one failed and one successful dependency settles on the failure
// rather than waiting for a signal that never arrives.
func TestExecuteSettlesOnMixedDependencyOutcomes(t *testing.T) {
	a := newComponent("a")
	b := newComponent("b")
	c := newComponent("c", "a", "b")
	design := newDesign("mixed", a, b, c)

	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	run := newExecution()
	run.run(t, plan, map[string]bool{"a": true})

	if indexOf(run.order(), "c") != -1 {
		t.Errorf("dispatched %v, want \"c\" withheld", run.order())
	}

	if got := run.withheld["c"]; got != "a" {
		t.Errorf("withheld[\"c\"] = %q, want \"a\"", got)
	}
}

func TestNewParallelProcessGraphRejectsAnEdgeToAnUnknownNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("known", component.ComponentDefinition{})
	g.AddEdge("known", "missing")

	pg, err := NewParallelProcessGraph(g)
	if pg != nil {
		t.Error("NewParallelProcessGraph() returned a graph for a dangling edge, want none")
	}

	if got := errorCode(t, err); got != ErrUnknownPlanNodeCode {
		t.Errorf("NewParallelProcessGraph() error code = %s, want %s", got, ErrUnknownPlanNodeCode)
	}
}

func TestNewParallelProcessGraphRejectsAnEdgeFromAnUnknownNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("known", component.ComponentDefinition{})
	g.AddEdge("missing", "known")

	if _, err := NewParallelProcessGraph(g); errorCode(t, err) != ErrUnknownPlanNodeCode {
		t.Errorf("NewParallelProcessGraph() error code = %s, want %s", errorCode(t, err), ErrUnknownPlanNodeCode)
	}
}

func TestDetectCycleRejectsAnEdgeToAnUnknownNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("known", component.ComponentDefinition{})
	g.AddEdge("known", "missing")

	if _, err := g.DetectCycle(); errorCode(t, err) != ErrUnknownPlanNodeCode {
		t.Errorf("Graph.DetectCycle() error code = %s, want %s", errorCode(t, err), ErrUnknownPlanNodeCode)
	}
}

func TestCreatePlanRejectsCyclicDependencies(t *testing.T) {
	a := newComponent("a", "b")
	b := newComponent("b", "a")
	design := newDesign("cyclic", a, b)

	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	feasible, err := plan.IsFeasible()
	if err != nil {
		t.Fatalf("Plan.IsFeasible() returned an unexpected error: %v", err)
	}

	if feasible {
		t.Error("Plan.IsFeasible() = true, want false for a dependency cycle")
	}
}
