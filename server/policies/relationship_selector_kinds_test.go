package policies

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// knownUnmodelledSelectorKinds pins selectors, keyed `<model>/<kind>`, whose
// kind is not modelled by the model they name at all, as opposed to being a
// misspelling of a kind it does model. Whether to drop these selectors or add
// the missing components is a maintainer decision tracked in
// meshery/meshery#23093. They are pinned rather than ignored: each pinned kind
// must still be missing from its model, so modelling one forces its removal
// here. Keys use the registered model name, so they stay valid whether a
// selector spells that name exactly or matches it through matchName.
var knownUnmodelledSelectorKinds = map[string]bool{
	"aws-applicationautoscaling-controller/AdoptedResource": true,
	"aws-cloudfront-controller/AdoptedResource":             true,
	"aws-cloudtrail-controller/AdoptedResource":             true,
	"aws-cloudwatch-controller/AdoptedResource":             true,
	"aws-cloudwatchlogs-controller/AdoptedResource":         true,
	"aws-documentdb-controller/AdoptedResource":             true,
	"aws-dynamodb-controller/AdoptedResource":               true,
	"aws-ec2-controller/SpotFleetRequest":                   true,
	"aws-mq-controller/AdoptedResource":                     true,
	"aws-opensearchservice-controller/AdoptedResource":      true,
	"aws-prometheusservice-controller/AdoptedResource":      true,
	"aws-route53-controller/AdoptedResource":                true,
	"aws-route53resolver-controller/AdoptedResource":        true,
	"aws-s3-controller/AdoptedResource":                     true,
	"aws-secretsmanager-controller/AdoptedResource":         true,
	"aws-sfn-controller/AdoptedResource":                    true,
	"aws-sqs-controller/AdoptedResource":                    true,
	"aws-ssm-controller/AdoptedResource":                    true,
	"aws-vpc-cni/PolicyEndpoint":                            true,
	"azure-compute/KeyVault":                                true,
	"kubernetes/KubernetesConnection":                       true,
	"kubernetes/KubernetesCredential":                       true,
	"kubernetes/PodSecurityPolicy":                          true,
	"meshery-shapes/Section":                                true,
}

// seededComponentKinds returns, per seeded model name, the set of component
// kinds shipped in its seeded definition directories.
func seededComponentKinds(t *testing.T, seeded []seededModelDefDir) map[string]map[string]bool {
	t.Helper()
	kinds := map[string]map[string]bool{}
	for _, s := range seeded {
		if kinds[s.name] == nil {
			kinds[s.name] = map[string]bool{}
		}
		files, err := filepath.Glob(filepath.Join(s.dir, "components", "*.json"))
		if err != nil {
			t.Fatalf("globbing components in %s: %v", s.dir, err)
		}
		for _, file := range files {
			raw, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("reading %s: %v", file, err)
			}
			var definition struct {
				Component struct {
					Kind string `json:"kind"`
				} `json:"component"`
			}
			if err := json.Unmarshal(raw, &definition); err != nil {
				t.Fatalf("parsing %s: %v", file, err)
			}
			if definition.Component.Kind != "" {
				kinds[s.name][definition.Component.Kind] = true
			}
		}
	}
	return kinds
}

// resolveSelectorModels returns the seeded models a selector's model.name
// refers to: the model of that exact name when one exists, otherwise every
// model matchName accepts, as isRelationshipFeasible would.
func resolveSelectorModels(name string, kinds map[string]map[string]bool) []string {
	if _, exact := kinds[name]; exact {
		return []string{name}
	}
	resolved := []string{}
	for model := range kinds {
		if matchName(model, name) {
			resolved = append(resolved, model)
		}
	}
	sort.Strings(resolved)
	return resolved
}

// A selector matches a component only when matchName accepts both the
// component's model name and its kind (isRelationshipFeasible). A selector
// whose kind is not a component of the model it names can therefore never
// match. Every selector kind in the seeded relationships must be a component
// kind of the model the selector names (meshery/meshery#23093). Selectors
// naming a model that does not exist are left to
// TestRelationshipSelectorsNameRegisteredModels.
func TestRelationshipSelectorKindsExistInTheirModel(t *testing.T) {
	root := repoRoot(t)
	seeded := seededModelDefDirs(t, root)
	kinds := seededComponentKinds(t, seeded)

	unknown := map[string]int{}
	for _, s := range seeded {
		files, err := filepath.Glob(filepath.Join(s.dir, "relationships", "*.json"))
		if err != nil {
			t.Fatalf("globbing relationships in %s: %v", s.dir, err)
		}
		for _, file := range files {
			definition, err := loadRelationshipDefinition(file)
			if err != nil {
				t.Fatal(err)
			}
			for _, item := range allSelectorItems(definition) {
				kind := selectorItemKind(item)
				if kind == "" || kind == "*" || item.Model == nil {
					continue
				}
				name := item.Model.Name
				if name == "" || name == "*" {
					continue
				}
				models := resolveSelectorModels(name, kinds)
				if len(models) == 0 {
					continue
				}
				matched := false
				for _, model := range models {
					for componentKind := range kinds[model] {
						if matchName(componentKind, kind) {
							matched = true
							break
						}
					}
					if matched {
						break
					}
				}
				if matched {
					continue
				}
				key := models[0] + "/" + kind
				if knownUnmodelledSelectorKinds[key] {
					continue
				}
				unknown[key]++
			}
		}
	}

	if len(unknown) > 0 {
		total, summary := summarizeCounts(unknown)
		t.Errorf("%d selectors use a kind their model does not have: %s", total, summary)
	}
	for key := range knownUnmodelledSelectorKinds {
		model, kind, _ := strings.Cut(key, "/")
		for componentKind := range kinds[model] {
			if matchName(componentKind, kind) {
				t.Errorf("%s is pinned in knownUnmodelledSelectorKinds but %s now models %s; remove the pin", key, model, kind)
				break
			}
		}
	}
}
