package policies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// A relationship definition's mutatorRef/mutatedRef paths address a component's
// `configuration` document. When a path names a field the component's own JSON
// schema does not define, the evaluation engine happily writes it and the value
// never reaches the rendered resource - meshery/meshery#21482, where the Ingress
// -> Service edge patched the pre-1.22 extensions/v1beta1 `backend.serviceName`
// shape into a networking.k8s.io/v1 component.
//
// These tests resolve every `configuration`-rooted mutation path declared by the
// in-scope relationship definitions against the component schema shipped beside
// it, so that class of defect fails here instead of silently in a user's design.

// relationshipScopeDirs are the model version directories whose relationship
// definitions are held to schema conformance, relative to the repo root.
var relationshipScopeDirs = []string{
	"models/kubernetes/v1.37.0-rc.0/v1.0.0",
	"models/cert-manager/v1.21.1/v1.0.0",
}

// componentSchemaDirs maps a selector item's model name to the model version
// directory whose components resolve its paths. A selector referencing a model
// absent from this map is out of scope and skipped.
var componentSchemaDirs = map[string]string{
	"kubernetes":   "models/kubernetes/v1.37.0-rc.0/v1.0.0",
	"cert-manager": "models/cert-manager/v1.21.1/v1.0.0",
}

// knownUnresolvedMutationPaths are definitions that already carry the
// meshery/meshery#21482 defect and predate this guard. They are pinned rather
// than ignored: each entry must still fail to resolve, so fixing one forces its
// removal here. Tracked for repair as a follow-up to #21482.
var knownUnresolvedMutationPaths = map[string]bool{
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-aatvf.json|ServiceAccount|configuration.spec.signerName":                                  true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-amdty.json|WatchEvent|configuration.object.kind":                                          true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-cyynq.json|StorageVersion|configuration.spec.resource.group":                              true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-djklo.json|Node|configuration.spec.nodeName":                                              true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-gbkch.json|ClusterRoleBinding|configuration.spec.roleRef.name":                            true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-keceb.json|RoleBinding|configuration.spec.roleRef.name":                                   true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-lnrwq.json|Node|configuration.spec.nodeSelector":                                          true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-mszif.json|PriorityLevelConfiguration|configuration.spec.priorityLevelConfiguration.name": true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-tvdia.json|ServiceAccount|configuration.spec.audiences":                                   true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-tvyop.json|Job|configuration.spec.jobTemplate":                                            true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-ulrah.json|ResourceClaim|configuration.spec.deviceClassName":                              true,
	"models/kubernetes/v1.37.0-rc.0/v1.0.0/relationships/edge-non-binding-reference-wasdd.json|StorageClass|configuration.spec.storageClassName":                              true,
}

type relationshipSelectorItem struct {
	Kind  string `json:"kind"`
	Model struct {
		Name string `json:"name"`
	} `json:"model"`
	Patch struct {
		MutatorRef [][]string `json:"mutatorRef"`
		MutatedRef [][]string `json:"mutatedRef"`
	} `json:"patch"`
}

type relationshipSelector struct {
	Allow struct {
		From []relationshipSelectorItem `json:"from"`
		To   []relationshipSelectorItem `json:"to"`
	} `json:"allow"`
}

type relationshipDefinitionFile struct {
	Selectors []relationshipSelector `json:"selectors"`
}

type componentDefinitionFile struct {
	Component struct {
		Kind string `json:"kind"`
		// Schema is a JSON document carried as an embedded string.
		Schema string `json:"schema"`
	} `json:"component"`
}

// componentSchemaLoader resolves a (model name, kind) pair to the parsed JSON
// schema shipped for that component, caching each read.
type componentSchemaLoader struct {
	root  string
	cache map[string]map[string]any
}

func newComponentSchemaLoader(root string) *componentSchemaLoader {
	return &componentSchemaLoader{root: root, cache: map[string]map[string]any{}}
}

// schema returns the component's parsed schema. ok is false when the model is
// out of scope, the component ships no JSON file, or its schema is empty - all
// cases where conformance cannot be decided rather than cases that fail.
func (l *componentSchemaLoader) schema(modelName, kind string) (schema map[string]any, ok bool, err error) {
	modelDir, inScope := componentSchemaDirs[modelName]
	if !inScope {
		return nil, false, nil
	}
	key := modelName + "/" + kind
	if cached, hit := l.cache[key]; hit {
		return cached, cached != nil, nil
	}

	path := filepath.Join(l.root, modelDir, "components", kind+".json")
	raw, readErr := os.ReadFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			l.cache[key] = nil
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("reading %s: %w", path, readErr)
	}

	var component componentDefinitionFile
	if err := json.Unmarshal(raw, &component); err != nil {
		return nil, false, fmt.Errorf("parsing %s: %w", path, err)
	}
	if strings.TrimSpace(component.Component.Schema) == "" {
		l.cache[key] = nil
		return nil, false, nil
	}

	parsed := map[string]any{}
	if err := json.Unmarshal([]byte(component.Component.Schema), &parsed); err != nil {
		return nil, false, fmt.Errorf("parsing embedded schema of %s: %w", path, err)
	}
	l.cache[key] = parsed
	return parsed, true, nil
}

var arrayIndexSegment = regexp.MustCompile(`^\d+$`)

// schemaAlternatives flattens allOf/oneOf/anyOf wrappers into the set of nodes a
// segment may be looked up in.
func schemaAlternatives(node map[string]any) []map[string]any {
	alternatives := []map[string]any{node}
	for _, combiner := range []string{"allOf", "oneOf", "anyOf"} {
		entries, isList := node[combiner].([]any)
		if !isList {
			continue
		}
		for _, entry := range entries {
			sub, isObject := entry.(map[string]any)
			if !isObject {
				continue
			}
			alternatives = append(alternatives, schemaAlternatives(sub)...)
		}
	}
	return alternatives
}

// descendSchema resolves one mutation path segment. A numeric segment or the
// wildcard `_` descends into an array's item schema; anything else is a property
// name, falling back to additionalProperties for open maps.
func descendSchema(node map[string]any, segment string) (map[string]any, bool) {
	descendsIntoItem := segment == "_" || arrayIndexSegment.MatchString(segment)
	for _, alternative := range schemaAlternatives(node) {
		if descendsIntoItem {
			if items, isObject := alternative["items"].(map[string]any); isObject {
				return items, true
			}
			continue
		}
		if properties, isObject := alternative["properties"].(map[string]any); isObject {
			if next, isObject := properties[segment].(map[string]any); isObject {
				return next, true
			}
		}
		switch additional := alternative["additionalProperties"].(type) {
		case map[string]any:
			return additional, true
		case bool:
			if additional {
				return map[string]any{}, true
			}
		}
	}
	return nil, false
}

func resolveSchemaPath(schema map[string]any, segments []string) bool {
	node := schema
	for _, segment := range segments {
		next, resolved := descendSchema(node, segment)
		if !resolved {
			return false
		}
		node = next
	}
	return true
}

// schemaCheckedSegments reports the segments of a mutation path that the
// component schema is expected to define, and whether the path is checkable at
// all. Paths outside `configuration` (displayName, component.kind) address the
// component envelope, and the whole `configuration.metadata` subtree is
// Kubernetes ObjectMeta, which Meshery strips from the shipped schema.
func schemaCheckedSegments(path []string) ([]string, bool) {
	if len(path) < 2 || path[0] != "configuration" {
		return nil, false
	}
	segments := path[1:]
	if segments[0] == "metadata" {
		return nil, false
	}
	return segments, true
}

func TestRelationshipMutationPathsResolveAgainstComponentSchemas(t *testing.T) {
	root := repoRoot(t)
	loader := newComponentSchemaLoader(root)

	unresolved := map[string]bool{}
	checked := 0

	for _, modelDir := range relationshipScopeDirs {
		files, err := filepath.Glob(filepath.Join(root, modelDir, "relationships", "*.json"))
		if err != nil {
			t.Fatalf("globbing relationships in %s: %v", modelDir, err)
		}
		if len(files) == 0 {
			t.Fatalf("no relationship definitions found under %s - the scope list is stale", modelDir)
		}

		for _, file := range files {
			raw, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("reading %s: %v", file, err)
			}
			var definition relationshipDefinitionFile
			if err := json.Unmarshal(raw, &definition); err != nil {
				t.Fatalf("parsing %s: %v", file, err)
			}

			relFile, err := filepath.Rel(root, file)
			if err != nil {
				t.Fatalf("relativizing %s: %v", file, err)
			}
			relFile = filepath.ToSlash(relFile)

			for _, selector := range definition.Selectors {
				items := append(append([]relationshipSelectorItem{}, selector.Allow.From...), selector.Allow.To...)
				for _, item := range items {
					if item.Kind == "" || item.Kind == "*" {
						continue
					}
					paths := append(append([][]string{}, item.Patch.MutatorRef...), item.Patch.MutatedRef...)
					for _, path := range paths {
						segments, checkable := schemaCheckedSegments(path)
						if !checkable {
							continue
						}
						schema, ok, err := loader.schema(item.Model.Name, item.Kind)
						if err != nil {
							t.Fatalf("%s: %v", relFile, err)
						}
						if !ok {
							continue
						}
						checked++
						if resolveSchemaPath(schema, segments) {
							continue
						}
						key := fmt.Sprintf("%s|%s|%s", relFile, item.Kind, strings.Join(path, "."))
						unresolved[key] = true
						if !knownUnresolvedMutationPaths[key] {
							t.Errorf("%s: %s mutates %q, which the shipped %s component schema does not define",
								relFile, item.Kind, strings.Join(path, "."), item.Kind)
						}
					}
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal("no configuration-rooted mutation paths were checked - the scope or model map is wrong")
	}

	stale := make([]string, 0, len(knownUnresolvedMutationPaths))
	for key := range knownUnresolvedMutationPaths {
		if !unresolved[key] {
			stale = append(stale, key)
		}
	}
	sort.Strings(stale)
	for _, key := range stale {
		t.Errorf("knownUnresolvedMutationPaths entry no longer fails and must be deleted: %s", key)
	}
}

// TestIngressBackendMutationUsesNetworkingV1Shape pins meshery/meshery#21482:
// the Ingress -> Service edge patched the pre-1.22 extensions/v1beta1 backend
// shape, which the networking.k8s.io/v1 component schema does not define.
func TestIngressBackendMutationUsesNetworkingV1Shape(t *testing.T) {
	root := repoRoot(t)
	loader := newComponentSchemaLoader(root)

	schema, ok, err := loader.schema("kubernetes", "Ingress")
	if err != nil {
		t.Fatalf("loading Ingress component schema: %v", err)
	}
	if !ok {
		t.Fatal("kubernetes Ingress component ships no schema")
	}

	resolving := [][]string{
		{"spec", "rules", "0", "http", "paths", "0", "backend", "service", "name"},
		{"spec", "rules", "0", "http", "paths", "0", "backend", "service", "port", "number"},
		{"spec", "defaultBackend", "service", "name"},
		{"spec", "defaultBackend", "service", "port", "number"},
		{"spec", "tls", "0", "secretName"},
		{"spec", "ingressClassName"},
	}
	for _, path := range resolving {
		if !resolveSchemaPath(schema, path) {
			t.Errorf("networking.k8s.io/v1 path %q does not resolve against the shipped Ingress schema", strings.Join(path, "."))
		}
	}

	superseded := [][]string{
		{"spec", "rules", "0", "http", "paths", "0", "backend", "serviceName"},
		{"spec", "rules", "0", "http", "paths", "0", "backend", "servicePort"},
	}
	for _, path := range superseded {
		if resolveSchemaPath(schema, path) {
			t.Errorf("extensions/v1beta1 path %q unexpectedly resolves - the guard for #21482 no longer bites", strings.Join(path, "."))
		}
	}
}
