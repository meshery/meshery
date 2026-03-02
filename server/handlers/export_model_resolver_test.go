package handlers

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
)

func TestResolveModelForExport_PrefersExactNameMatch(t *testing.T) {
	aero := &modelv1beta1.ModelDefinition{Name: "aerospike-kubernetes-operator"}
	k8s := &modelv1beta1.ModelDefinition{Name: "kubernetes"}

	entities := []entity.Entity{aero, k8s}

	got, status, msg := resolveModelForExport(entities, "", "Kubernetes")
	if status != 0 {
		t.Fatalf("expected success, got status=%d msg=%q", status, msg)
	}
	if got == nil || got.Name != "kubernetes" {
		t.Fatalf("expected kubernetes, got %#v", got)
	}
}

func TestResolveModelForExport_AmbiguousPartialMatches(t *testing.T) {
	aero := &modelv1beta1.ModelDefinition{Name: "aerospike-kubernetes-operator"}
	k8s := &modelv1beta1.ModelDefinition{Name: "kubernetes"}

	entities := []entity.Entity{aero, k8s}

	got, status, _ := resolveModelForExport(entities, "", "kube")
	if got != nil {
		t.Fatalf("expected nil model, got %#v", got)
	}
	if status != http.StatusConflict {
		t.Fatalf("expected %d, got %d", http.StatusConflict, status)
	}
}

func TestResolveModelForExport_IDMustBeUnique(t *testing.T) {
	aero := &modelv1beta1.ModelDefinition{Name: "aerospike-kubernetes-operator"}
	k8s := &modelv1beta1.ModelDefinition{Name: "kubernetes"}

	entities := []entity.Entity{aero, k8s}

	got, status, _ := resolveModelForExport(entities, "82f90070-9745-6d97-414a-3fe5343dd7b6", "")
	if got != nil {
		t.Fatalf("expected nil model, got %#v", got)
	}
	if status != http.StatusConflict {
		t.Fatalf("expected %d, got %d", http.StatusConflict, status)
	}
}

func TestResolveModelForExport_IDExactMatch(t *testing.T) {
	k8s := &modelv1beta1.ModelDefinition{Name: "kubernetes"}
	entities := []entity.Entity{k8s}

	got, status, msg := resolveModelForExport(entities, "82f90070-9745-6d97-414a-3fe5343dd7b6", "")
	if status != 0 {
		t.Fatalf("expected success, got status=%d msg=%q", status, msg)
	}
	if got == nil || got.Name != "kubernetes" {
		t.Fatalf("expected kubernetes, got %#v", got)
	}
}

func TestResolveModelForExport_NotFound(t *testing.T) {
	got, status, _ := resolveModelForExport(nil, "", "kubernetes")
	if got != nil {
		t.Fatalf("expected nil model, got %#v", got)
	}
	if status != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, status)
	}
}

func TestResolveModelForExport_DedupesExactNameByID(t *testing.T) {
	k8s := &modelv1beta1.ModelDefinition{
		Id:   uuid.FromStringOrNil("8b34062d-6599-2805-acff-117735854b73"),
		Name: "kubernetes",
	}
	dup := &modelv1beta1.ModelDefinition{
		Id:   uuid.FromStringOrNil("8b34062d-6599-2805-acff-117735854b73"),
		Name: "kubernetes",
	}
	entities := []entity.Entity{k8s, dup}

	got, status, msg := resolveModelForExport(entities, "", "kubernetes")
	if status != 0 {
		t.Fatalf("expected success, got status=%d msg=%q", status, msg)
	}
	if got == nil || got.Id.String() != "8b34062d-6599-2805-acff-117735854b73" {
		t.Fatalf("expected deduped kubernetes model, got %#v", got)
	}
}

func TestResolveModelForExport_IDAllowsDuplicateRowsSameID(t *testing.T) {
	k8s := &modelv1beta1.ModelDefinition{
		Id:   uuid.FromStringOrNil("8b34062d-6599-2805-acff-117735854b73"),
		Name: "kubernetes",
	}
	dup := &modelv1beta1.ModelDefinition{
		Id:   uuid.FromStringOrNil("8b34062d-6599-2805-acff-117735854b73"),
		Name: "kubernetes",
	}
	entities := []entity.Entity{k8s, dup}

	got, status, msg := resolveModelForExport(entities, "8b34062d-6599-2805-acff-117735854b73", "")
	if status != 0 {
		t.Fatalf("expected success, got status=%d msg=%q", status, msg)
	}
	if got == nil || got.Name != "kubernetes" {
		t.Fatalf("expected kubernetes, got %#v", got)
	}
}

func TestResolveModelForExport_AmbiguousExactMatchesDifferentIDs(t *testing.T) {
	first := &modelv1beta1.ModelDefinition{
		Id:      uuid.FromStringOrNil("8b34062d-6599-2805-acff-117735854b73"),
		Name:    "kubernetes",
		Version: "v1.0.0",
		Model:   modelv1beta1.Model{Version: "v1.35.0-rc.1"},
	}
	second := &modelv1beta1.ModelDefinition{
		Id:      uuid.FromStringOrNil("fcf07220-c5c9-4ad4-b5d5-d91d74549d96"),
		Name:    "kubernetes",
		Version: "v1.0.0",
		Model:   modelv1beta1.Model{Version: "v1.34.0"},
	}
	entities := []entity.Entity{first, second}

	got, status, _ := resolveModelForExport(entities, "", "kubernetes")
	if got != nil {
		t.Fatalf("expected nil model, got %#v", got)
	}
	if status != http.StatusConflict {
		t.Fatalf("expected %d, got %d", http.StatusConflict, status)
	}
}
