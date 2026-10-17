package handlers

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshery/server/models/pattern/planner"
	"github.com/meshery/meshery/server/models/pattern/stages"
	"github.com/meshery/meshkit/logger"
	events "github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

func mutateTestComponent(displayName, kind string) *component.ComponentDefinition {
	return &component.ComponentDefinition{
		ID:          uuid.Must(uuid.NewV4()),
		DisplayName: displayName,
		Component:   component.Component{Kind: kind},
	}
}

func newMutateActionProvider(t *testing.T) *serviceActionProvider {
	t.Helper()

	log, err := logger.New("mutate-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	return &serviceActionProvider{log: log}
}

// The dependency Mutate writes has to be one the execution plan can resolve: a
// dependency names another component of the same design by its name, and no
// component depends on itself.
func TestMutateOrdersComponentsBehindACustomResourceDefinition(t *testing.T) {
	crd := mutateTestComponent("widgets.example.com", "CustomResourceDefinition")
	widget := mutateTestComponent("a widget", "Widget")

	design := pattern.PatternFile{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       "mutate-test",
		Components: []*component.ComponentDefinition{crd, widget},
	}

	newMutateActionProvider(t).Mutate(&design)

	dependsOn, err := planner.DeclaredDependencies(design.Name, widget)
	if err != nil {
		t.Fatalf("the dependency Mutate wrote cannot be read back: %v", err)
	}

	if len(dependsOn) != 1 || dependsOn[0] != crd.DisplayName {
		t.Errorf("%q depends on %v, want [%q]", widget.DisplayName, dependsOn, crd.DisplayName)
	}

	for _, dep := range dependsOn {
		if dep == widget.DisplayName || dep == widget.ID.String() {
			t.Errorf("%q was made to depend on itself through %q", widget.DisplayName, dep)
		}
	}

	if crdDependsOn, err := planner.DeclaredDependencies(design.Name, crd); err != nil || len(crdDependsOn) != 0 {
		t.Errorf("the custom resource definition depends on %v (err %v), want nothing", crdDependsOn, err)
	}

	// The plan is what consumes the dependency, so building one is what proves
	// Mutate wrote a value that is actually usable.
	plan, err := planner.CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("no plan could be built from the design Mutate produced: %v", err)
	}

	log, err := logger.New("mutate-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	order := []string{}
	if err := plan.Execute(func(id string, comp component.ComponentDefinition) bool {
		order = append(order, comp.DisplayName)
		return true
	}, func(string, component.ComponentDefinition, string) {}, log); err != nil {
		t.Fatalf("the plan could not be executed: %v", err)
	}

	if len(order) != 2 || order[0] != crd.DisplayName {
		t.Errorf("deployed in the order %v, want the custom resource definition first", order)
	}
}

// A design that already declares the dependency Mutate would add is left alone,
// rather than accumulating the same name twice.
func TestMutateDoesNotRepeatADependencyTheDesignAlreadyDeclares(t *testing.T) {
	crd := mutateTestComponent("widgets.example.com", "CustomResourceDefinition")
	widget := mutateTestComponent("a widget", "Widget")
	widget.Metadata.AdditionalProperties = map[string]interface{}{
		"dependsOn": []interface{}{crd.DisplayName},
	}

	design := pattern.PatternFile{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       "mutate-test",
		Components: []*component.ComponentDefinition{crd, widget},
	}

	newMutateActionProvider(t).Mutate(&design)

	dependsOn, err := planner.DeclaredDependencies(design.Name, widget)
	if err != nil {
		t.Fatalf("the dependency Mutate wrote cannot be read back: %v", err)
	}

	if len(dependsOn) != 1 || dependsOn[0] != crd.DisplayName {
		t.Errorf("%q depends on %v, want [%q]", widget.DisplayName, dependsOn, crd.DisplayName)
	}
}

func componentOutcome(name string, success bool) []patterns.DeploymentMessagePerContext {
	return []patterns.DeploymentMessagePerContext{
		{
			Summary: []patterns.DeploymentMessagePerComp{
				{Kind: "Widget", CompName: name, Success: success},
			},
		},
	}
}

// A deploy whose components did not all make it has to say so, rather than
// reporting a completed deployment the way a fully successful one does.
func TestDeploymentOutcomeReportsFailedComponents(t *testing.T) {
	failedPrerequisite := componentOutcome("with-prerequisite", true)
	failedPrerequisite[0].Summary = append(failedPrerequisite[0].Summary, patterns.DeploymentMessagePerComp{
		CompName:       "with-prerequisite",
		Success:        false,
		IsPrerequisite: true,
	})

	tests := []struct {
		name            string
		action          string
		response        map[string]interface{}
		wantSeverity    events.EventSeverity
		wantDescription string
	}{
		{
			name:   "every component deployed",
			action: "deploy",
			response: map[string]interface{}{
				"a": componentOutcome("a", true),
				"b": componentOutcome("b", true),
			},
			wantSeverity:    events.Informational,
			wantDescription: "Deployment completed for design 'outcome-test'.",
		},
		{
			name:   "some components failed to deploy",
			action: "deploy",
			response: map[string]interface{}{
				"a": componentOutcome("a", true),
				"b": componentOutcome("b", false),
				"c": componentOutcome("c", true),
			},
			wantSeverity:    events.Warning,
			wantDescription: "Failed to deploy 1 of 3 components in design 'outcome-test'.",
		},
		{
			name:   "every component failed to deploy",
			action: "deploy",
			response: map[string]interface{}{
				"a": componentOutcome("a", false),
				"b": componentOutcome("b", false),
			},
			wantSeverity:    events.Error,
			wantDescription: "Failed to deploy 2 of 2 components in design 'outcome-test'.",
		},
		{
			name:   "some components failed to undeploy",
			action: "undeploy",
			response: map[string]interface{}{
				"a": componentOutcome("a", false),
				"b": componentOutcome("b", true),
			},
			wantSeverity:    events.Warning,
			wantDescription: "Failed to undeploy 1 of 2 components in design 'outcome-test'.",
		},
		{
			name:   "every component undeployed",
			action: "undeploy",
			response: map[string]interface{}{
				"a": componentOutcome("a", true),
			},
			wantSeverity:    events.Success,
			wantDescription: "undeployed design 'outcome-test'.",
		},
		{
			// Installing prerequisites is fail-forward, so a failed one does not
			// make the component it was installed for a failure.
			name:   "only a prerequisite failed",
			action: "deploy",
			response: map[string]interface{}{
				"with-prerequisite": failedPrerequisite,
			},
			wantSeverity:    events.Informational,
			wantDescription: "Deployment completed for design 'outcome-test'.",
		},
		{
			name:   "dry run",
			action: "dry-run",
			response: map[string]interface{}{
				stages.DryRunResponseKey: map[string]interface{}{},
			},
			wantSeverity:    events.Informational,
			wantDescription: "Dry Run completed for design 'outcome-test'.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			severity, description := deploymentOutcome(tt.action, "outcome-test", tt.response)

			if severity != tt.wantSeverity {
				t.Errorf("severity is %q, want %q", severity, tt.wantSeverity)
			}

			if description != tt.wantDescription {
				t.Errorf("description is %q, want %q", description, tt.wantDescription)
			}
		})
	}
}
