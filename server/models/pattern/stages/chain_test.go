package stages

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshery/server/models/pattern/planner"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta2/component"
	componentv1beta3 "github.com/meshery/schemas/models/v1beta3/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

// The tests below drive the whole deployment chain - Format, Filler, Validator,
// Provision - rather than the Provision stage on its own.
//
// "dependsOn" is not a first-class field of a design, so it survives decoding
// as an untyped entry that holds a []interface{} rather than a []string. A
// suite that constructs a design in memory and enters at Provision never sees
// the stage that has to cope with that, which is why these designs are built by
// decoding JSON, exactly as a design imported or loaded by Meshery Server is.

const (
	chainKind       = "Namespace"
	chainAPIVersion = "v1"
	chainModel      = "kubernetes"
	chainModelVer   = "v1.25.0"
)

// seedChainComponent registers the component definition that the design
// components below declare, so that the Validator stage resolves them and the
// chain reaches Provision.
func seedChainComponent(t *testing.T, reg *registry.RegistryManager) {
	t.Helper()

	enabled := componentv1beta3.Enabled
	comp := componentv1beta3.ComponentDefinition{
		DisplayName:   chainKind,
		SchemaVersion: "core.meshery.io/v1beta1",
		Status:        &enabled,
		Component: componentv1beta3.Component{
			Kind:    chainKind,
			Version: chainAPIVersion,
			Schema:  `{"properties":{}}`,
		},
		Model: &model.ModelDefinition{
			Name:          chainModel,
			DisplayName:   "Kubernetes",
			SchemaVersion: "models.meshery.io/v1beta1",
			Version:       chainModelVer,
			Model:         model.Model{Version: chainModelVer},
			Category:      category.CategoryDefinition{Name: "Orchestration"},
			Status:        model.Enabled,
		},
	}

	id, err := comp.GenerateID()
	if err != nil {
		t.Fatalf("failed to generate a component id: %v", err)
	}
	comp.ID = id

	registrant := connection.Connection{
		Name:    "chain-test-registrant",
		Kind:    chainModel,
		Type:    "platform",
		SubType: "orchestration",
		Status:  connection.ConnectionStatusConnected,
	}

	if _, _, err := reg.RegisterEntity(registry.RegistrantHostToV1beta3(registrant), &comp); err != nil {
		t.Fatalf("failed to register the component definition: %v", err)
	}
}

// chainComponent describes one component of a design the way it is written in a
// design document, before any decoding.
func chainComponent(displayName string, dependsOn ...string) map[string]interface{} {
	metadata := map[string]interface{}{}
	if len(dependsOn) > 0 {
		metadata["dependsOn"] = dependsOn
	}

	return map[string]interface{}{
		"id":            uuid.Must(uuid.NewV4()).String(),
		"displayName":   displayName,
		"configuration": map[string]interface{}{},
		"metadata":      metadata,
		"component": map[string]interface{}{
			"kind":    chainKind,
			"version": chainAPIVersion,
		},
		"model": map[string]interface{}{
			"name":  chainModel,
			"model": map[string]interface{}{"version": chainModelVer},
		},
	}
}

// chainComponentWithRawDependsOn describes a component whose "dependsOn" is
// written in some shape other than a list of names.
func chainComponentWithRawDependsOn(displayName string, dependsOn interface{}) map[string]interface{} {
	comp := chainComponent(displayName)
	comp["metadata"] = map[string]interface{}{"dependsOn": dependsOn}

	return comp
}

// decodeDesign encodes a design and decodes it back into a PatternFile, so that
// the chain is handed a design in the representation it really arrives in.
func decodeDesign(t *testing.T, components ...map[string]interface{}) *pattern.PatternFile {
	t.Helper()

	raw, err := json.Marshal(map[string]interface{}{
		"id":            uuid.Must(uuid.NewV4()).String(),
		"name":          "chain-test",
		"schemaVersion": "designs.meshery.io/v1beta1",
		"version":       "1.0.0",
		"components":    components,
	})
	if err != nil {
		t.Fatalf("failed to encode the design: %v", err)
	}

	design := &pattern.PatternFile{}
	if err := json.Unmarshal(raw, design); err != nil {
		t.Fatalf("failed to decode the design: %v", err)
	}

	for _, comp := range design.Components {
		if _, declared := comp.Metadata.AdditionalProperties["dependsOn"]; !declared {
			continue
		}

		// Guards the premise of these tests: a decoded design carries
		// "dependsOn" as a []interface{}. Should a schemas release start
		// decoding it as a []string, these tests would still pass while no
		// longer covering the shape the chain actually meets.
		if _, isStrings := comp.Metadata.AdditionalProperties["dependsOn"].([]string); isStrings {
			t.Fatalf("a decoded design carried \"dependsOn\" on %q as []string; these tests assume the untyped shape", comp.DisplayName)
		}
	}

	return design
}

// runChain runs the full deployment chain over a design and returns the outcome
// reported for each component, keyed by the component's display name.
func runChain(t *testing.T, act *stubActionProvider, components ...map[string]interface{}) (*pattern.PatternFile, map[string][]patterns.DeploymentMessagePerContext) {
	t.Helper()

	log, err := logger.New("chain-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	design := decodeDesign(t, components...)

	data := &Data{
		Pattern:                       design,
		DeclartionToDefinitionMapping: map[core.Uuid]component.ComponentDefinition{},
		Other:                         map[string]interface{}{},
	}

	CreateChain().
		Add(Format()).
		Add(Filler(true)).
		Add(Validator(&stubInfoProvider{}, act, false)).
		Add(Provision(&stubInfoProvider{}, act, log)).
		Process(data)

	reported := map[string][]patterns.DeploymentMessagePerContext{}

	data.Lock.Lock()
	defer data.Lock.Unlock()

	names := map[string]string{}
	for _, comp := range design.Components {
		names[comp.ID.String()] = comp.DisplayName
	}

	for key, value := range data.Other {
		id := key[:len(key)-len(ProvisionSuffixKey)]

		msgs, ok := value.([]patterns.DeploymentMessagePerContext)
		if !ok {
			t.Fatalf("data.Other[%q] holds %T, want []patterns.DeploymentMessagePerContext", key, value)
		}

		reported[names[id]] = msgs
	}

	return design, reported
}

func newChainActionProvider(t *testing.T) *stubActionProvider {
	t.Helper()

	act := newStubActionProvider(t)
	seedChainComponent(t, act.GetRegistry())

	return act
}

// A design that declares a dependency is deployed in dependency order when it
// is run through the whole chain, not only when the Provision stage is entered
// directly with a design assembled in memory.
func TestChainDeploysInDependencyOrder(t *testing.T) {
	act := newChainActionProvider(t)

	_, reported := runChain(t, act,
		chainComponent("alpha"),
		chainComponent("beta", "alpha"),
		chainComponent("gamma"),
	)

	if act.terminated != nil {
		t.Fatalf("the chain terminated with %v, want it to reach deployment", act.terminated)
	}

	order := act.order()
	if len(order) != 3 {
		t.Fatalf("dispatched %v, want all three components dispatched", order)
	}

	alpha, beta := -1, -1
	for i, name := range order {
		switch name {
		case "alpha":
			alpha = i
		case "beta":
			beta = i
		}
	}

	if alpha == -1 || beta == -1 {
		t.Fatalf("dispatched %v, want both \"alpha\" and \"beta\" dispatched", order)
	}

	if alpha > beta {
		t.Errorf("dispatched %v, want \"alpha\" before \"beta\": \"beta\" depends on it", order)
	}

	for _, name := range []string{"alpha", "beta", "gamma"} {
		if _, ok := reported[name]; !ok {
			t.Errorf("%q produced no deployment message", name)
		}
	}
}

// A reference inside "dependsOn" is resolved by the Filler stage, and the plan
// built afterwards orders on the resolved name.
//
// This is what makes the resolution durable rather than local to the Filler: an
// unresolved reference names no component of the design, so a plan built from
// one is rejected and nothing is deployed at all.
func TestChainResolvesReferencesInDependsOn(t *testing.T) {
	act := newChainActionProvider(t)

	_, _ = runChain(t, act,
		chainComponent("alpha"),
		chainComponent("beta", "$(#ref.components.0.displayName)"),
	)

	if act.terminated != nil {
		t.Fatalf("the chain terminated with %v, want the reference resolved to \"alpha\" and the design deployed", act.terminated)
	}

	order := act.order()
	if len(order) != 2 {
		t.Fatalf("dispatched %v, want both components dispatched", order)
	}

	if order[0] != "alpha" {
		t.Errorf("dispatched %v, want \"alpha\" first: \"beta\" depends on it through a reference", order)
	}
}

// A component whose declared dependency failed to apply is withheld, over the
// whole chain. This is the guarantee the deployment-engine documentation makes.
func TestChainWithholdsDependentsOfAFailedComponent(t *testing.T) {
	act := newChainActionProvider(t)
	act.applyFailures["alpha"] = true

	_, reported := runChain(t, act,
		chainComponent("alpha"),
		chainComponent("beta", "alpha"),
		chainComponent("gamma"),
	)

	if act.terminated != nil {
		t.Fatalf("the chain terminated with %v, want it to reach deployment", act.terminated)
	}

	order := act.order()
	if dispatched(order, "beta") {
		t.Errorf("dispatched %v, want \"beta\" withheld: \"alpha\" failed to apply", order)
	}

	if !dispatched(order, "gamma") {
		t.Errorf("dispatched %v, want \"gamma\" deployed: it declared no dependency on \"alpha\"", order)
	}

	withheld, ok := reported["beta"]
	if !ok {
		t.Fatal("\"beta\" produced no deployment message, want it reported as withheld")
	}

	if succeeded(withheld) {
		t.Error("\"beta\" was reported as successful, want it reported as withheld")
	}

	if !succeeded(reported["gamma"]) {
		t.Error("\"gamma\" was not reported as successful, want it unaffected by \"alpha\"")
	}
}

// A "dependsOn" written in a shape Meshery cannot read is reported as exactly
// that, carrying its own code all the way to where the chain terminates.
//
// The failure is not a reference that could not be resolved, and reporting it
// as one sends the user looking for a reference that was never there.
func TestChainReportsAMalformedDependsOnAsItsOwnError(t *testing.T) {
	act := newChainActionProvider(t)

	runChain(t, act,
		chainComponent("alpha"),
		chainComponentWithRawDependsOn("beta", "alpha"),
	)

	if order := act.order(); len(order) != 0 {
		t.Errorf("dispatched %v, want nothing deployed when the design cannot be read", order)
	}

	if got := errors.GetCode(act.terminated); got != planner.ErrInvalidDependencyCode {
		t.Errorf("the chain terminated with error code %s, want the malformed-dependency code %s", got, planner.ErrInvalidDependencyCode)
	}
}

// Every entry of a "dependsOn" list has to be a component name; a list holding
// something else is reported the same way.
func TestChainReportsANonStringDependsOnEntryAsItsOwnError(t *testing.T) {
	act := newChainActionProvider(t)

	runChain(t, act,
		chainComponent("alpha"),
		chainComponentWithRawDependsOn("beta", []interface{}{"alpha", 7}),
	)

	if got := errors.GetCode(act.terminated); got != planner.ErrInvalidDependencyCode {
		t.Errorf("the chain terminated with error code %s, want the malformed-dependency code %s", got, planner.ErrInvalidDependencyCode)
	}
}

// A design whose dependency names no component of it is rejected before
// anything is deployed, rather than being read as a reference that failed to
// resolve.
func TestChainRejectsAnUnresolvableDependency(t *testing.T) {
	act := newChainActionProvider(t)

	runChain(t, act,
		chainComponent("alpha"),
		chainComponent("beta", "nowhere"),
	)

	if order := act.order(); len(order) != 0 {
		t.Errorf("dispatched %v, want nothing deployed when the plan cannot be built", order)
	}

	if got := errors.GetCode(act.terminated); got != planner.ErrUnresolvedDependencyCode {
		t.Errorf("the chain terminated with error code %s, want the unresolvable-dependency code %s", got, planner.ErrUnresolvedDependencyCode)
	}
}
