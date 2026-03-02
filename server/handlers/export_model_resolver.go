package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
)

// resolveModelForExport deterministically selects a single model from a set of candidates.
//
// Rules:
// - Duplicate rows for the same model are ignored (deduped).
// - If modelID is provided, there must be exactly one unique candidate.
// - If name is provided:
//   - Prefer a single exact (case-insensitive) name match, if present.
//   - If no exact match exists:
//     - If exactly one unique candidate exists, select it.
//     - If multiple candidates exist, return 409 (do not guess).
// - If neither modelID nor name is provided, select the first candidate.
//
// It returns (model, 0, "") on success. On failure, it returns (nil, httpStatus, message).
func resolveModelForExport(entities []entity.Entity, modelID, name string) (*modelv1beta1.ModelDefinition, int, string) {
	candidates := uniqueModelCandidates(entities)
	if len(candidates) == 0 {
		return nil, http.StatusNotFound, "model not found"
	}

	if modelID != "" {
		if len(candidates) != 1 {
			return nil, http.StatusConflict, fmt.Sprintf("multiple models match id %s; please report this as a bug", modelID)
		}
		return candidates[0], 0, ""
	}

	if name != "" {
		exact := make([]*modelv1beta1.ModelDefinition, 0, 1)
		matches := make([]string, 0, len(candidates))

		for _, md := range candidates {
			matches = append(matches, md.Name)
			if strings.EqualFold(md.Name, name) {
				exact = append(exact, md)
			}
		}

		switch {
		case len(exact) == 1:
			return exact[0], 0, ""
		case len(exact) > 1:
			return nil, http.StatusConflict, fmt.Sprintf("multiple models exactly match name %s; please specify --version or export by id", name)
		case len(candidates) == 1:
			return candidates[0], 0, ""
		default:
			return nil, http.StatusConflict, fmt.Sprintf("multiple models match name %s: %s", name, strings.Join(matches, ", "))
		}
	}

	return candidates[0], 0, ""
}

func uniqueModelCandidates(entities []entity.Entity) []*modelv1beta1.ModelDefinition {
	candidates := make([]*modelv1beta1.ModelDefinition, 0, len(entities))
	seen := make(map[string]struct{}, len(entities))

	for _, ent := range entities {
		md, ok := ent.(*modelv1beta1.ModelDefinition)
		if !ok || md == nil {
			continue
		}

		key := modelCandidateKey(md)
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		candidates = append(candidates, md)
	}

	return candidates
}

func modelCandidateKey(md *modelv1beta1.ModelDefinition) string {
	if md.Id != uuid.Nil {
		return "id:" + md.Id.String()
	}

	registrantKind := ""
	if md.Registrant.Kind != "" {
		registrantKind = strings.ToLower(string(md.Registrant.Kind))
	}

	return strings.ToLower(fmt.Sprintf(
		"name:%s|version:%s|modelVersion:%s|registrant:%s",
		md.Name,
		md.Version,
		md.Model.Version,
		registrantKind,
	))
}
