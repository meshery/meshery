package planner

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

// TestParallelProcess_HighConcurrencyFanIn verifies that a dependent node
// with a large number of upstream prerequisites receives all completion signals
// concurrently without channel buffer stalls or lock contention.
func TestParallelProcess_HighConcurrencyFanIn(t *testing.T) {
	const numPrereqs = 20
	prereqNames := make([]string, numPrereqs)
	components := make([]*component.ComponentDefinition, 0, numPrereqs+1)

	for i := 0; i < numPrereqs; i++ {
		name := fmt.Sprintf("prereq-%d", i)
		prereqNames[i] = name
		components = append(components, newComponent(name))
	}

	target := newComponent("target", prereqNames...)
	components = append(components, target)

	design := newDesign("fan-in", components...)
	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	feasible, err := plan.IsFeasible()
	if err != nil || !feasible {
		t.Fatalf("Plan.IsFeasible() = (%v, %v), want (true, nil)", feasible, err)
	}

	var prereqCompletedCount int64
	var targetExecuted int64
	var targetSawAllPrereqs bool

	done := make(chan error, 1)
	go func() {
		done <- plan.Execute(
			func(id string, comp component.ComponentDefinition) bool {
				if comp.DisplayName == "target" {
					completed := atomic.LoadInt64(&prereqCompletedCount)
					if completed == int64(numPrereqs) {
						targetSawAllPrereqs = true
					}
					atomic.AddInt64(&targetExecuted, 1)
					return true
				}

				// Simulate variable async work duration
				time.Sleep(time.Duration(1+rand.Intn(5)) * time.Millisecond)
				atomic.AddInt64(&prereqCompletedCount, 1)
				return true
			},
			func(string, component.ComponentDefinition, string) {},
			testLogger(t),
		)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Plan.Execute() returned an unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Plan.Execute() timed out; channel buffer stall or deadlock during fan-in")
	}

	if atomic.LoadInt64(&targetExecuted) != 1 {
		t.Errorf("target executed %d times, want 1", targetExecuted)
	}
	if !targetSawAllPrereqs {
		t.Errorf("target executed before all %d prerequisites completed", numPrereqs)
	}
}

// TestParallelProcess_DiamondTopology verifies that parallel branch execution
// correctly converges at the sink node.
func TestParallelProcess_DiamondTopology(t *testing.T) {
	root := newComponent("root")
	branchA := newComponent("branchA", "root")
	branchB := newComponent("branchB", "root")
	sink := newComponent("sink", "branchA", "branchB")

	design := newDesign("diamond", root, branchA, branchB, sink)
	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	var mx sync.Mutex
	dispatched := []string{}

	done := make(chan error, 1)
	go func() {
		done <- plan.Execute(
			func(id string, comp component.ComponentDefinition) bool {
				mx.Lock()
				dispatched = append(dispatched, comp.DisplayName)
				mx.Unlock()
				time.Sleep(2 * time.Millisecond)
				return true
			},
			func(string, component.ComponentDefinition, string) {},
			testLogger(t),
		)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Plan.Execute() returned an unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Plan.Execute() timed out on diamond topology")
	}

	if len(dispatched) != 4 {
		t.Fatalf("dispatched %v, want 4 components", dispatched)
	}

	rootIdx := indexOf(dispatched, "root")
	aIdx := indexOf(dispatched, "branchA")
	bIdx := indexOf(dispatched, "branchB")
	sinkIdx := indexOf(dispatched, "sink")

	if rootIdx > aIdx || rootIdx > bIdx {
		t.Errorf("root must execute before branches, got order: %v", dispatched)
	}
	if aIdx > sinkIdx || bIdx > sinkIdx {
		t.Errorf("both branches must execute before sink, got order: %v", dispatched)
	}
}

// TestParallelProcess_MultiDependentFanOutNonBlocking tests that an upstream node
// with multiple downstream dependents sends signals concurrently without blocking.
func TestParallelProcess_MultiDependentFanOutNonBlocking(t *testing.T) {
	producer := newComponent("producer")
	const numConsumers = 10
	components := []*component.ComponentDefinition{producer}

	for i := 0; i < numConsumers; i++ {
		consumer := newComponent(fmt.Sprintf("consumer-%d", i), "producer")
		components = append(components, consumer)
	}

	design := newDesign("fan-out", components...)
	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned an unexpected error: %v", err)
	}

	var consumersExecuted int64
	done := make(chan error, 1)
	go func() {
		done <- plan.Execute(
			func(id string, comp component.ComponentDefinition) bool {
				if comp.DisplayName != "producer" {
					atomic.AddInt64(&consumersExecuted, 1)
				}
				return true
			},
			func(string, component.ComponentDefinition, string) {},
			testLogger(t),
		)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Plan.Execute() returned an unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Plan.Execute() timed out during fan-out dispatch")
	}

	if consumersExecuted != numConsumers {
		t.Errorf("executed %d consumers, want %d", consumersExecuted, numConsumers)
	}
}

// TestPlanExecute_RejectsCyclicDependency verifies that Plan.Execute self-guards
// against cyclic plans and returns ErrCyclicDependency immediately without hanging.
func TestPlanExecute_RejectsCyclicDependency(t *testing.T) {
	g := NewGraph()
	aID := uuid.Must(uuid.NewV4()).String()
	bID := uuid.Must(uuid.NewV4()).String()

	g.AddNode(aID, component.ComponentDefinition{DisplayName: "a"})
	g.AddNode(bID, component.ComponentDefinition{DisplayName: "b"})
	g.AddEdge(aID, bID)
	g.AddEdge(bID, aID) // Cycle

	plan := &Plan{
		Data: pattern.PatternFile{
			Name: "cyclic-plan",
		},
		Graph: g,
	}

	done := make(chan error, 1)
	go func() {
		done <- plan.Execute(
			func(string, component.ComponentDefinition) bool { return true },
			func(string, component.ComponentDefinition, string) {},
			testLogger(t),
		)
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Plan.Execute() expected error for cyclic graph, got nil")
		}
		if got := errors.GetCode(err); got != ErrCyclicDependencyCode {
			t.Errorf("Plan.Execute() error code = %s, want %s", got, ErrCyclicDependencyCode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Plan.Execute() hung indefinitely on cyclic dependency; missing cycle defense")
	}
}

// TestParallelProcess_MixedOutcomesWithManyPrereqs verifies that when multiple
// upstream prerequisites fail and succeed concurrently, the dependent is withheld
// and all channels drain properly without blocking senders.
func TestParallelProcess_MixedOutcomesWithManyPrereqs(t *testing.T) {
	const numPrereqs = 8
	prereqNames := make([]string, numPrereqs)
	components := make([]*component.ComponentDefinition, 0, numPrereqs+1)

	for i := 0; i < numPrereqs; i++ {
		name := fmt.Sprintf("p-%d", i)
		prereqNames[i] = name
		components = append(components, newComponent(name))
	}

	target := newComponent("target", prereqNames...)
	components = append(components, target)

	design := newDesign("mixed-fan-in", components...)
	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned unexpected error: %v", err)
	}

	var targetVisited bool
	var targetWithheld bool
	var failedDepReported string
	var mx sync.Mutex

	// Make odd prerequisites fail, even prerequisites succeed
	done := make(chan error, 1)
	go func() {
		done <- plan.Execute(
			func(id string, comp component.ComponentDefinition) bool {
				mx.Lock()
				defer mx.Unlock()

				if comp.DisplayName == "target" {
					targetVisited = true
					return true
				}

				var idx int
				_, _ = fmt.Sscanf(comp.DisplayName, "p-%d", &idx)
				return idx%2 == 0
			},
			func(id string, comp component.ComponentDefinition, failedDep string) {
				mx.Lock()
				defer mx.Unlock()
				if comp.DisplayName == "target" {
					targetWithheld = true
					failedDepReported = failedDep
				}
			},
			testLogger(t),
		)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Plan.Execute() returned unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Plan.Execute() timed out during mixed outcome fan-in")
	}

	mx.Lock()
	defer mx.Unlock()

	if targetVisited {
		t.Error("target was visited, want it withheld because prerequisites failed")
	}
	if !targetWithheld {
		t.Error("target was not reported as withheld")
	}
	if failedDepReported == "" {
		t.Error("target did not record which dependency failed")
	}
}

// TestParallelProcess_PanicInPrerequisiteDoesNotDeadlockDependents verifies that when
// a prerequisite node panics during visitor execution, downstream dependents are
// not permanently starved on channel reads, and the plan execution does not deadlock.
func TestParallelProcess_PanicInPrerequisiteDoesNotDeadlockDependents(t *testing.T) {
	prereq := newComponent("prereq")
	dependent := newComponent("dependent", "prereq")

	design := newDesign("panic-recovery", prereq, dependent)
	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned unexpected error: %v", err)
	}

	var (
		mu              sync.Mutex
		dependentVisited bool
		dependentWithheld bool
		reportedDep     string
	)

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			_ = recover()
		}()
		_ = plan.Execute(
			func(id string, comp component.ComponentDefinition) bool {
				if comp.DisplayName == "prereq" {
					panic("simulated fatal panic in component visitor")
				}
				mu.Lock()
				if comp.DisplayName == "dependent" {
					dependentVisited = true
				}
				mu.Unlock()
				return true
			},
			func(id string, comp component.ComponentDefinition, failedDep string) {
				mu.Lock()
				if comp.DisplayName == "dependent" {
					dependentWithheld = true
					reportedDep = failedDep
				}
				mu.Unlock()
			},
			testLogger(t),
		)
	}()

	select {
	case <-done:
		// Plan execution completed without deadlocking
	case <-time.After(2 * time.Second):
		t.Fatal("Plan.Execute() deadlocked; prerequisite panic bypassed signal propagation, permanently wedging dependent channel reads")
	}

	mu.Lock()
	defer mu.Unlock()

	if dependentVisited {
		t.Error("dependent node was visited despite prerequisite panic")
	}
	if !dependentWithheld {
		t.Error("dependent node was not marked as withheld following prerequisite panic")
	}
	if reportedDep != prereq.ID.String() {
		t.Errorf("dependent reported failed dependency %q, want %q", reportedDep, prereq.ID.String())
	}
}

// TestParallelProcess_PanicInMiddleOfDeepChainPropagatesWithholding verifies that when
// a middle node in a multi-hop dependency chain panics (A -> B -> C -> D, B panics),
// subsequent nodes C and D are sequentially withheld without deadlock or starvation.
func TestParallelProcess_PanicInMiddleOfDeepChainPropagatesWithholding(t *testing.T) {
	nodeA := newComponent("nodeA")
	nodeB := newComponent("nodeB", "nodeA")
	nodeC := newComponent("nodeC", "nodeB")
	nodeD := newComponent("nodeD", "nodeC")

	design := newDesign("deep-chain-panic", nodeA, nodeB, nodeC, nodeD)
	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned unexpected error: %v", err)
	}

	var (
		mu         sync.Mutex
		aVisited   bool
		cVisited   bool
		dVisited   bool
		cWithheld  bool
		dWithheld  bool
		cFailedDep string
		dFailedDep string
	)

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			_ = recover()
		}()
		_ = plan.Execute(
			func(id string, comp component.ComponentDefinition) bool {
				mu.Lock()
				defer mu.Unlock()
				switch comp.DisplayName {
				case "nodeA":
					aVisited = true
					return true
				case "nodeB":
					panic("simulated fatal panic in nodeB")
				case "nodeC":
					cVisited = true
					return true
				case "nodeD":
					dVisited = true
					return true
				}
				return true
			},
			func(id string, comp component.ComponentDefinition, failedDep string) {
				mu.Lock()
				defer mu.Unlock()
				switch comp.DisplayName {
				case "nodeC":
					cWithheld = true
					cFailedDep = failedDep
				case "nodeD":
					dWithheld = true
					dFailedDep = failedDep
				}
			},
			testLogger(t),
		)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Plan.Execute() deadlocked; panic in middle of deep chain hung downstream nodes")
	}

	mu.Lock()
	defer mu.Unlock()

	if !aVisited {
		t.Error("nodeA was not visited, want it visited before nodeB")
	}
	if cVisited {
		t.Error("nodeC was visited, want it withheld because nodeB panicked")
	}
	if dVisited {
		t.Error("nodeD was visited, want it withheld because nodeC was withheld")
	}
	if !cWithheld || cFailedDep != nodeB.ID.String() {
		t.Errorf("nodeC withheld=%v, failedDep=%q, want true and %q", cWithheld, cFailedDep, nodeB.ID.String())
	}
	if !dWithheld || dFailedDep != nodeC.ID.String() {
		t.Errorf("nodeD withheld=%v, failedDep=%q, want true and %q", dWithheld, dFailedDep, nodeC.ID.String())
	}
}

// TestParallelProcess_DiamondTopologyPanicOneBranchWithholdsSink verifies that in a diamond
// DAG (Root -> BranchA, BranchB -> Sink), when BranchA panics while BranchB succeeds,
// the Sink receives both signals, correctly determines it cannot proceed, and is withheld.
func TestParallelProcess_DiamondTopologyPanicOneBranchWithholdsSink(t *testing.T) {
	root := newComponent("root")
	branchA := newComponent("branchA", "root")
	branchB := newComponent("branchB", "root")
	sink := newComponent("sink", "branchA", "branchB")

	design := newDesign("diamond-panic", root, branchA, branchB, sink)
	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned unexpected error: %v", err)
	}

	var (
		mu           sync.Mutex
		sinkVisited  bool
		sinkWithheld bool
		sinkFailedDep string
	)

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			_ = recover()
		}()
		_ = plan.Execute(
			func(id string, comp component.ComponentDefinition) bool {
				switch comp.DisplayName {
				case "branchA":
					panic("simulated fatal panic in branchA")
				case "branchB":
					return true
				case "sink":
					mu.Lock()
					sinkVisited = true
					mu.Unlock()
					return true
				}
				return true
			},
			func(id string, comp component.ComponentDefinition, failedDep string) {
				if comp.DisplayName == "sink" {
					mu.Lock()
					sinkWithheld = true
					sinkFailedDep = failedDep
					mu.Unlock()
				}
			},
			testLogger(t),
		)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Plan.Execute() deadlocked; panic in diamond branch hung sink node")
	}

	mu.Lock()
	defer mu.Unlock()

	if sinkVisited {
		t.Error("sink node was visited, want it withheld because branchA panicked")
	}
	if !sinkWithheld {
		t.Error("sink node was not reported as withheld")
	}
	if sinkFailedDep != branchA.ID.String() {
		t.Errorf("sink reported failed dependency %q, want %q", sinkFailedDep, branchA.ID.String())
	}
}

// TestParallelProcess_FanInHalfPanickingPrereqsDoesNotDeadlock verifies high-concurrency fan-in
// when 10 upstream prerequisites succeed while 10 prerequisites panic concurrently. The target
// node must receive all 20 signals without channel buffer deadlock or dropped signals.
func TestParallelProcess_FanInHalfPanickingPrereqsDoesNotDeadlock(t *testing.T) {
	const totalPrereqs = 20
	prereqNames := make([]string, totalPrereqs)
	components := make([]*component.ComponentDefinition, 0, totalPrereqs+1)

	for i := 0; i < totalPrereqs; i++ {
		name := fmt.Sprintf("fan-prereq-%d", i)
		prereqNames[i] = name
		components = append(components, newComponent(name))
	}

	target := newComponent("fan-target", prereqNames...)
	components = append(components, target)

	design := newDesign("fan-in-panics", components...)
	plan, err := CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("CreatePlan() returned unexpected error: %v", err)
	}

	var (
		mu             sync.Mutex
		targetVisited  bool
		targetWithheld bool
	)

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			_ = recover()
		}()
		_ = plan.Execute(
			func(id string, comp component.ComponentDefinition) bool {
				if comp.DisplayName == "fan-target" {
					mu.Lock()
					targetVisited = true
					mu.Unlock()
					return true
				}

				// Even indices succeed, odd indices panic
				var idx int
				_, _ = fmt.Sscanf(comp.DisplayName, "fan-prereq-%d", &idx)
				if idx%2 == 1 {
					panic(fmt.Sprintf("panic in %s", comp.DisplayName))
				}
				return true
			},
			func(id string, comp component.ComponentDefinition, failedDep string) {
				if comp.DisplayName == "fan-target" {
					mu.Lock()
					targetWithheld = true
					mu.Unlock()
				}
			},
			testLogger(t),
		)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Plan.Execute() deadlocked; high-concurrency mixed panic fan-in hung target node")
	}

	mu.Lock()
	defer mu.Unlock()

	if targetVisited {
		t.Error("target was visited, want it withheld because half prerequisites panicked")
	}
	if !targetWithheld {
		t.Error("target was not reported as withheld")
	}
}


