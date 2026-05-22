package policies

import (
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// Shared fixtures and builders for policy engine tests. These keep the typed
// PatternFile/component/relationship boilerplate in one place so individual
// tests describe only what is unique about their design, not how to assemble
// the schema structs.

// strPtr returns a pointer to the given string.
func strPtr(s string) *string { return &s }

// testUUID returns a deterministic UUID from a short hex suffix, left-padded
// into the reserved all-zero test prefix: testUUID("1") is
// 00000000-0000-0000-0000-000000000001 and testUUID("dead") ends in ...0000dead.
func testUUID(suffix string) uuid.UUID {
	if len(suffix) > 12 {
		panic("testUUID: suffix longer than 12 hex digits: " + suffix)
	}
	id, err := uuid.FromString("00000000-0000-0000-0000-" + strings.Repeat("0", 12-len(suffix)) + suffix)
	if err != nil {
		panic("testUUID: invalid suffix " + suffix)
	}
	return id
}

// comp builds a ComponentDefinition of the given kind and id. Pass a nil
// config for components that do not need one.
func comp(kind string, id uuid.UUID, config map[string]interface{}) *component.ComponentDefinition {
	c := &component.ComponentDefinition{
		Component:     component.Component{Kind: kind},
		Configuration: config,
	}
	c.ID = id
	return c
}

// k8sComp is comp with the kubernetes model reference set, the common case.
func k8sComp(kind string, id uuid.UUID, config map[string]interface{}) *component.ComponentDefinition {
	c := comp(kind, id, config)
	c.ModelReference = modelv1beta1.ModelReference{Name: "kubernetes"}
	return c
}

// k8sModel returns the kubernetes model reference used inside selectors.
func k8sModel() *modelv1beta1.ModelReference {
	return &modelv1beta1.ModelReference{Name: "kubernetes"}
}

// mutatorRef and mutatedRef build single-path ref pointers.
func mutatorRef(path ...string) *relationship.MutatorRef {
	r := relationship.MutatorRef{path}
	return &r
}

func mutatedRef(path ...string) *relationship.MutatedRef {
	r := relationship.MutatedRef{path}
	return &r
}

// mutatorPatch and mutatedPatch build a selector patch carrying a single
// mutator or mutated path.
func mutatorPatch(path ...string) *relationship.RelationshipDefinitionSelectorsPatch {
	return &relationship.RelationshipDefinitionSelectorsPatch{MutatorRef: mutatorRef(path...)}
}

func mutatedPatch(path ...string) *relationship.RelationshipDefinitionSelectorsPatch {
	return &relationship.RelationshipDefinitionSelectorsPatch{MutatedRef: mutatedRef(path...)}
}

// makePatternFile builds a PatternFile from typed components and relationships.
func makePatternFile(comps []*component.ComponentDefinition, rels []*relationship.RelationshipDefinition) *pattern.PatternFile {
	return &pattern.PatternFile{
		Components:    comps,
		Relationships: rels,
	}
}
