package planner

import (
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

func testComponent(t *testing.T, displayName string, dependsOn ...string) *component.ComponentDefinition {
	t.Helper()

	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate component id: %v", err)
	}

	comp := &component.ComponentDefinition{
		ID:          id,
		DisplayName: displayName,
	}
	if len(dependsOn) > 0 {
		comp.Metadata.AdditionalProperties = map[string]interface{}{
			"dependsOn": dependsOn,
		}
	}

	return comp
}

func testLogger(t *testing.T) logger.Handler {
	t.Helper()

	log, err := logger.New("planner-test", logger.Options{
		Format: logger.SyslogLogFormat,
	})
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	return log
}

// executionOrder runs the plan and returns component display names in the
// order the visit callback was invoked.
func executionOrder(t *testing.T, plan *Plan) []string {
	t.Helper()

	var mu sync.Mutex
	order := []string{}

	err := plan.Execute(func(_ string, comp component.ComponentDefinition) bool {
		mu.Lock()
		defer mu.Unlock()
		order = append(order, string(comp.DisplayName))
		return true
	}, testLogger(t))
	if err != nil {
		t.Fatalf("plan execution failed: %v", err)
	}

	return order
}

func indexOf(items []string, target string) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return -1
}

func TestCreatePlanDependsOnByDisplayName(t *testing.T) {
	db := testComponent(t, "db")
	app := testComponent(t, "app", "db")

	plan, err := CreatePlan(pattern.PatternFile{Components: []*component.ComponentDefinition{app, db}}, false)
	if err != nil {
		t.Fatalf("CreatePlan returned an error: %v", err)
	}

	if !plan.IsFeasible() {
		t.Fatal("expected an acyclic plan to be feasible")
	}

	order := executionOrder(t, plan)
	if len(order) != 2 {
		t.Fatalf("expected 2 components to be visited, got %d: %v", len(order), order)
	}
	if indexOf(order, "db") > indexOf(order, "app") {
		t.Fatalf("expected db to be provisioned before its dependent app, got order %v", order)
	}
}

func TestCreatePlanInvertReversesOrder(t *testing.T) {
	db := testComponent(t, "db")
	app := testComponent(t, "app", "db")

	plan, err := CreatePlan(pattern.PatternFile{Components: []*component.ComponentDefinition{app, db}}, true)
	if err != nil {
		t.Fatalf("CreatePlan returned an error: %v", err)
	}

	order := executionOrder(t, plan)
	if indexOf(order, "app") > indexOf(order, "db") {
		t.Fatalf("expected app to be removed before its dependency db on inverted plan, got order %v", order)
	}
}

func TestCreatePlanDetectsDisplayNameCycle(t *testing.T) {
	a := testComponent(t, "a", "b")
	b := testComponent(t, "b", "a")

	plan, err := CreatePlan(pattern.PatternFile{Components: []*component.ComponentDefinition{a, b}}, false)
	if err != nil {
		t.Fatalf("CreatePlan returned an error: %v", err)
	}

	if plan.IsFeasible() {
		t.Fatal("expected a cyclic dependsOn chain to be reported as infeasible")
	}
}

func TestCreatePlanDependsOnByComponentID(t *testing.T) {
	db := testComponent(t, "db")
	app := testComponent(t, "app", db.ID.String())

	plan, err := CreatePlan(pattern.PatternFile{Components: []*component.ComponentDefinition{app, db}}, false)
	if err != nil {
		t.Fatalf("CreatePlan returned an error: %v", err)
	}

	order := executionOrder(t, plan)
	if indexOf(order, "db") > indexOf(order, "app") {
		t.Fatalf("expected db to be provisioned before its dependent app, got order %v", order)
	}
}

func TestCreatePlanAmbiguousDependencyErrors(t *testing.T) {
	db1 := testComponent(t, "db")
	db2 := testComponent(t, "db")
	app := testComponent(t, "app", "db")

	_, err := CreatePlan(pattern.PatternFile{Components: []*component.ComponentDefinition{app, db1, db2}}, false)
	if err == nil {
		t.Fatal("expected an error when a dependsOn entry matches multiple components")
	}
}

func TestCreatePlanUnknownDependencyErrors(t *testing.T) {
	app := testComponent(t, "app", "ghost")

	_, err := CreatePlan(pattern.PatternFile{Components: []*component.ComponentDefinition{app}}, false)
	if err == nil {
		t.Fatal("expected an error for a dependsOn entry that matches no component")
	}
}
