package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/meshery/schemas/models/core"
	componentv1beta1 "github.com/meshery/schemas/models/v1beta1/component"
	componentv1beta3 "github.com/meshery/schemas/models/v1beta3/component"
)

// TestComponentV1beta3ToV1beta1_NilSource guards against the nil-source
// crash path that was observed in the relationship-policy evaluator before
// this bridge was introduced.
func TestComponentV1beta3ToV1beta1_NilSource(t *testing.T) {
	if got := ComponentV1beta3ToV1beta1(nil); got != nil {
		t.Fatalf("expected nil for nil source, got %#v", got)
	}
}

// TestComponentV1beta3ToV1beta1_FieldsCopied is a regression test for the
// SIGSEGV in policy_relationship_handler.go:447. The handler asserts the
// meshkit-registry result to *v1beta1/component.ComponentDefinition, but
// the registry now stores entities as *v1beta3/component.ComponentDefinition.
// The assertion silently failed and `_component` was nil, panicking on the
// next field write. This bridge is what the handler now uses instead, so
// every field the handler later mutates must transfer cleanly.
func TestComponentV1beta3ToV1beta1_FieldsCopied(t *testing.T) {
	id := core.Uuid(uuid.New())
	modelID := core.Uuid(uuid.New())
	src := &componentv1beta3.ComponentDefinition{
		ID:          id,
		DisplayName: "kubernetes-pod-abc",
		Description: "A pod",
		Format:      "JSON",
		Metadata: componentv1beta3.ComponentDefinition_Metadata{
			IsAnnotation: true,
			Genealogy:    "deployment",
		},
		Configuration: map[string]interface{}{"replicas": 3},
		Component: componentv1beta3.Component{
			Kind:    "Pod",
			Version: "v1",
			Schema:  `{"type":"object"}`,
		},
		ModelID: &modelID,
	}

	dst := ComponentV1beta3ToV1beta1(src)
	if dst == nil {
		t.Fatal("expected non-nil destination")
	}
	if dst.ID != id {
		t.Errorf("ID: want %v, got %v", id, dst.ID)
	}
	if dst.DisplayName != "kubernetes-pod-abc" {
		t.Errorf("DisplayName: want %q, got %q", "kubernetes-pod-abc", dst.DisplayName)
	}
	if !dst.Metadata.IsAnnotation {
		t.Error("Metadata.IsAnnotation: want true, got false")
	}
	if dst.Metadata.Genealogy != "deployment" {
		t.Errorf("Metadata.Genealogy: want %q, got %q", "deployment", dst.Metadata.Genealogy)
	}
	if dst.Configuration["replicas"] != 3 {
		t.Errorf("Configuration[replicas]: want 3, got %v", dst.Configuration["replicas"])
	}
	if dst.Component.Kind != "Pod" {
		t.Errorf("Component.Kind: want Pod, got %q", dst.Component.Kind)
	}
	if dst.ModelId == nil || *dst.ModelId != modelID {
		t.Errorf("ModelId: want %v, got %v", modelID, dst.ModelId)
	}
	// Status round-trips through string-typed enums in both directions.
	enabled := componentv1beta3.ComponentDefinitionStatus("enabled")
	src.Status = &enabled
	dst = ComponentV1beta3ToV1beta1(src)
	if dst.Status == nil || *dst.Status != componentv1beta1.ComponentDefinitionStatus("enabled") {
		t.Errorf("Status: want enabled, got %v", dst.Status)
	}
}

// TestComponentV1beta3ToV1beta1_AliasingPreserved documents the contract
// that inner reference-typed fields alias the source — mirrors the
// existing ComponentV1beta3ToV1beta2 contract.
func TestComponentV1beta3ToV1beta1_AliasingPreserved(t *testing.T) {
	cfg := map[string]interface{}{"k": "v"}
	src := &componentv1beta3.ComponentDefinition{Configuration: cfg}
	dst := ComponentV1beta3ToV1beta1(src)
	dst.Configuration["k"] = "mutated"
	if cfg["k"] != "mutated" {
		t.Errorf("Configuration aliasing broken: source.k = %v, expected mutated", cfg["k"])
	}
}
