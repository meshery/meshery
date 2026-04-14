// Package context provides utilities for building LLM context from Meshery's
// component registry, including schema injection and token budget management.
package context

import (
	"encoding/json"
	"strings"

	"github.com/meshery/meshkit/models/meshmodel/entity"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/component"
)

// SchemaBuilder builds LLM-ready schema context from the Meshery component
// registry. It performs a two-pass strategy:
//
//  1. Keyword pass: find component names/kinds that appear in the user prompt.
//  2. Schema pass: fetch the full definitions for those components and format
//     them as compact JSON summaries for injection into the system prompt.
type SchemaBuilder struct {
	registry  *meshmodel.RegistryManager
	maxTokens int // approximate token budget for schema injection
}

// NewSchemaBuilder creates a SchemaBuilder backed by the given registry.
// maxTokens controls the approximate token budget allocated for schema text
// (use 0 for no limit, which is only appropriate in tests).
func NewSchemaBuilder(reg *meshmodel.RegistryManager, maxTokens int) *SchemaBuilder {
	return &SchemaBuilder{registry: reg, maxTokens: maxTokens}
}

// ComponentSummary is a compact representation of a Meshery component that
// can be embedded in a system prompt without excessive token cost.
type ComponentSummary struct {
	Kind        string         `json:"kind"`
	Model       string         `json:"model,omitempty"`
	DisplayName string         `json:"displayName,omitempty"`
	Description string         `json:"description,omitempty"`
	Properties  map[string]any `json:"properties,omitempty"`
}

// BuildContext returns a formatted string with component schema summaries that
// are relevant to the user prompt. The string is suitable for embedding in a
// system prompt. If no matching components are found, an empty string is returned.
func (b *SchemaBuilder) BuildContext(userPrompt string) string {
	components := b.findRelevantComponents(userPrompt)
	if len(components) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## Available Meshery Component Schemas\n\n")
	sb.WriteString("Use the following component definitions when generating Meshery designs:\n\n")

	for _, c := range components {
		summary := toSummary(c)
		bs, err := json.Marshal(summary)
		if err != nil {
			continue
		}
		sb.WriteString("```json\n")
		sb.Write(bs)
		sb.WriteString("\n```\n\n")

		// Respect token budget: stop adding schemas if we've exceeded the limit.
		// A rough approximation: 4 characters ≈ 1 token.
		if b.maxTokens > 0 && sb.Len()/4 > b.maxTokens {
			sb.WriteString("_[additional schemas omitted to fit context window]_\n")
			break
		}
	}

	return sb.String()
}

// findRelevantComponents searches the registry for components whose name or
// kind matches keywords extracted from the user prompt. It returns at most
// 20 components to keep context size manageable.
func (b *SchemaBuilder) findRelevantComponents(prompt string) []component.ComponentDefinition {
	keywords := ExtractKeywords(prompt)
	if len(keywords) == 0 {
		return nil
	}

	var comps []component.ComponentDefinition
	seen := map[string]bool{}

	for _, kw := range keywords {
		entities, _, _, err := b.registry.GetEntities(&regv1beta1.ComponentFilter{
			DisplayName: kw,
			Greedy:      true,
			Limit:       10,
		})
		if err != nil {
			continue
		}
		for _, e := range entities {
			c, ok := e.(*component.ComponentDefinition)
			if !ok || seen[c.Component.Kind] {
				continue
			}
			seen[c.Component.Kind] = true
			comps = append(comps, *c)
		}
		if len(comps) >= 20 {
			break
		}
	}

	return comps
}

// toSummary converts a full ComponentDefinition into a compact summary.
func toSummary(c component.ComponentDefinition) ComponentSummary {
	s := ComponentSummary{
		Kind:        c.Component.Kind,
		DisplayName: c.DisplayName,
	}
	if c.Model != nil {
		s.Model = c.Model.Name
	}
	if c.Description != "" {
		s.Description = c.Description
	}

	// Extract a trimmed properties map from the component schema if present.
	if c.Component.Schema != "" {
		var schema map[string]any
		if err := json.Unmarshal([]byte(c.Component.Schema), &schema); err == nil {
			if props, ok := schema["properties"].(map[string]any); ok {
				// Keep only the top-level property names and their types to
				// avoid ballooning the token count.
				compact := make(map[string]any, len(props))
				for k, v := range props {
					if vm, ok := v.(map[string]any); ok {
						compact[k] = map[string]any{
							"type":        vm["type"],
							"description": vm["description"],
						}
					}
				}
				s.Properties = compact
			}
		}
	}

	return s
}

// ExtractKeywords extracts meaningful tokens from a natural-language prompt
// that could match Meshery component names (e.g. "nginx", "deployment",
// "service", "kubernetes"). Exported so prompts.go can reuse it.
func ExtractKeywords(prompt string) []string {
	// A small set of well-known Meshery/Kubernetes keywords that signal the
	// user wants specific component types. Order matters: more specific first.
	knownKeywords := []string{
		"deployment", "statefulset", "daemonset", "replicaset", "cronjob", "job",
		"service", "ingress", "configmap", "secret", "namespace", "pod",
		"persistentvolumeclaim", "serviceaccount",
		"nginx", "redis", "postgres", "mysql", "mongodb", "kafka", "rabbitmq",
		"prometheus", "grafana", "jaeger",
		"gateway", "virtualservice", "destinationrule",
	}

	lower := strings.ToLower(prompt)
	var found []string
	seen := map[string]bool{}
	for _, kw := range knownKeywords {
		if strings.Contains(lower, kw) && !seen[kw] {
			found = append(found, kw)
			seen[kw] = true
		}
	}

	// If no known keywords match, fall back to splitting the prompt into words
	// and returning the first few non-trivial tokens as search terms.
	if len(found) == 0 {
		words := strings.Fields(prompt)
		for _, w := range words {
			w = strings.Trim(w, ".,!?;:")
			if len(w) > 4 {
				found = append(found, strings.ToLower(w))
			}
			if len(found) >= 3 {
				break
			}
		}
	}

	return found
}

// entity.Filter compatibility: the registry's GetEntities accepts entity.Filter
// but ComponentFilter implements it. This file imports entity to satisfy the
// compiler even though it's used indirectly via regv1beta1.ComponentFilter.
var _ entity.Filter = (*regv1beta1.ComponentFilter)(nil)
