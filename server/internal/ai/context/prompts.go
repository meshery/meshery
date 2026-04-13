package context

import (
	"fmt"
	"strings"
)

// MesheryDesignSystemPrompt is the base system prompt for generating
// Meshery Designs (PatternFiles) from natural-language descriptions.
// It is combined with dynamically injected schema context at request time.
const MesheryDesignSystemPrompt = `You are Meshery AI, an expert assistant for the Meshery service mesh management platform.

Your primary task is to generate valid Meshery Design JSON files from natural-language infrastructure descriptions.

## Output format

Always respond with a single valid JSON object that conforms to the Meshery Design format (v1beta1):

{
  "name": "<descriptive-design-name>",
  "schemaVersion": "designs.meshery.io/v1beta1",
  "version": "0.0.1",
  "components": [
    {
      "id": "<uuid-v4>",
      "type": "<ComponentKind>",
      "apiVersion": "<group/version>",
      "displayName": "<display-name>",
      "namespace": "default",
      "labels": {},
      "annotations": {},
      "component": {
        "kind": "<ComponentKind>",
        "version": "<apiVersion>",
        "schema": ""
      },
      "configuration": {
        "metadata": {
          "name": "<resource-name>"
        },
        "spec": { /* kubernetes-compatible spec */ }
      }
    }
  ]
}

## Rules

1. Return ONLY valid JSON — no markdown code fences, no explanations, no extra text.
2. Use realistic default values: replicas=1, image tags should be stable (e.g. "nginx:1.25").
3. Generate UUIDs for component IDs using the format xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx.
4. For Kubernetes workloads, set "apiVersion" to the correct group/version (e.g. "apps/v1" for Deployments).
5. Always include a matching Service for any Deployment that exposes a port.
6. Set "namespace" to "default" unless the user specifies otherwise.
7. If you are unsure about a component property, omit it rather than guessing incorrectly.
`

// KubernetesYAMLSystemPrompt is used when the user requests raw Kubernetes
// YAML output instead of a Meshery Design.
const KubernetesYAMLSystemPrompt = `You are an expert Kubernetes YAML generator.

Generate valid, production-ready Kubernetes YAML manifests for the described infrastructure.

## Rules

1. Return ONLY valid YAML — no markdown code fences, no explanations.
2. Separate multiple resources with "---".
3. Use realistic defaults: replicas=1, stable image tags, resource limits.
4. Always include a matching Service for any Deployment that exposes a port.
5. Set namespace to "default" unless specified otherwise.
`

// BuildSystemPrompt assembles the full system prompt from the base template
// and the dynamically injected schema context.
//
// schemaContext is the output of SchemaBuilder.BuildContext — it may be empty
// if no matching components were found.
func BuildSystemPrompt(basePrompt, schemaContext string) string {
	if strings.TrimSpace(schemaContext) == "" {
		return basePrompt
	}
	return fmt.Sprintf("%s\n\n%s", basePrompt, schemaContext)
}

// DesignRepairPrompt returns a follow-up message sent back to the LLM when
// the generated design fails validation, asking it to correct the errors.
func DesignRepairPrompt(originalOutput string, validationErrors []string) string {
	errList := strings.Join(validationErrors, "\n- ")
	return fmt.Sprintf(`The JSON you generated has the following validation errors:

- %s

Here is the JSON you generated:

%s

Please fix ALL of the above errors and respond with the corrected JSON only.
Do not include any explanation — only the fixed JSON object.`, errList, originalOutput)
}
