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

	"github.com/Masterminds/semver/v3"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// A relationship definition's mutatorRef/mutatedRef paths address a component's
// `configuration` document. When a path names a field the component's own JSON
// schema does not define, the evaluation engine happily writes it and the value
// never reaches the rendered resource - meshery/meshery#21482, where the Ingress
// -> Service edge patched the pre-1.22 extensions/v1beta1 `backend.serviceName`
// shape into a networking.k8s.io/v1 component.
//
// These tests resolve every `configuration`-rooted mutation path declared by the
// in-scope models against the component schema shipped in the same model version
// directory, so that class of defect fails here instead of in a user's design.

// relationshipConformanceModels are the models whose relationship definitions are
// held to schema conformance. Every version directory of each is checked.
var relationshipConformanceModels = []string{"kubernetes", "cert-manager"}

// knownUnresolvedMutationPaths are definitions that already carry the
// meshery/meshery#21482 defect and predate this guard, keyed
// `<model>/<file>|<kind>|<path>` because the model generator fans every row out
// to all of a model's version directories. They are pinned rather than ignored:
// each entry must still fail to resolve, so fixing one forces its removal here.
// Tracked for repair as a follow-up to #21482.
var knownUnresolvedMutationPaths = map[string]bool{
	"kubernetes/edge-non-binding-reference-aatvf.json|ServiceAccount|configuration.spec.signerName":                                  true,
	"kubernetes/edge-non-binding-reference-amdty.json|WatchEvent|configuration.object.kind":                                          true,
	"kubernetes/edge-non-binding-reference-cyynq.json|StorageVersion|configuration.spec.resource.group":                              true,
	"kubernetes/edge-non-binding-reference-djklo.json|Node|configuration.spec.nodeName":                                              true,
	"kubernetes/edge-non-binding-reference-gbkch.json|ClusterRoleBinding|configuration.spec.roleRef.name":                            true,
	"kubernetes/edge-non-binding-reference-keceb.json|RoleBinding|configuration.spec.roleRef.name":                                   true,
	"kubernetes/edge-non-binding-reference-lnrwq.json|Node|configuration.spec.nodeSelector":                                          true,
	"kubernetes/edge-non-binding-reference-mszif.json|PriorityLevelConfiguration|configuration.spec.priorityLevelConfiguration.name": true,
	"kubernetes/edge-non-binding-reference-tvdia.json|ServiceAccount|configuration.spec.audiences":                                   true,
	"kubernetes/edge-non-binding-reference-tvyop.json|Job|configuration.spec.jobTemplate":                                            true,
	"kubernetes/edge-non-binding-reference-ulrah.json|ResourceClaim|configuration.spec.deviceClassName":                              true,
	"kubernetes/edge-non-binding-reference-wasdd.json|StorageClass|configuration.spec.storageClassName":                              true,
}

// modelVersionDirs returns every `models/<model>/<version>/v1.0.0` directory,
// newest version first. Versions that do not parse as semver sort last.
func modelVersionDirs(root, modelName string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(root, "models", modelName, "*", "v1.0.0"))
	if err != nil {
		return nil, fmt.Errorf("globbing version directories of %s: %w", modelName, err)
	}
	versionOf := func(dir string) string {
		return filepath.Base(filepath.Dir(dir))
	}
	sort.Slice(matches, func(i, j int) bool {
		left, leftErr := semver.NewVersion(versionOf(matches[i]))
		right, rightErr := semver.NewVersion(versionOf(matches[j]))
		switch {
		case leftErr != nil && rightErr != nil:
			return versionOf(matches[i]) > versionOf(matches[j])
		case leftErr != nil:
			return false
		case rightErr != nil:
			return true
		default:
			return left.GreaterThan(right)
		}
	})
	return matches, nil
}

// componentSchemaLoader resolves a (model version directory, kind) pair to the
// parsed JSON schema shipped for that component, caching each read. A selector
// item is resolved against its own model version directory when it references
// the model the relationship belongs to, and against the newest version
// directory of the referenced model otherwise.
type componentSchemaLoader struct {
	root       string
	cache      map[string]map[string]any
	newestDirs map[string]string
}

func newComponentSchemaLoader(root string) *componentSchemaLoader {
	return &componentSchemaLoader{
		root:       root,
		cache:      map[string]map[string]any{},
		newestDirs: map[string]string{},
	}
}

// newestVersionDir returns the newest version directory of a model, or "" when
// the model ships none.
func (l *componentSchemaLoader) newestVersionDir(modelName string) (string, error) {
	if cached, hit := l.newestDirs[modelName]; hit {
		return cached, nil
	}
	dirs, err := modelVersionDirs(l.root, modelName)
	if err != nil {
		return "", err
	}
	newest := ""
	if len(dirs) > 0 {
		newest = dirs[0]
	}
	l.newestDirs[modelName] = newest
	return newest, nil
}

// schema returns the component's parsed schema. ok is false when the model ships
// no version directory, the component ships no JSON file, or its schema is empty
// - all cases where conformance cannot be decided rather than cases that fail.
func (l *componentSchemaLoader) schema(modelDir, kind string) (schema map[string]any, ok bool, err error) {
	key := modelDir + "|" + kind
	if cached, hit := l.cache[key]; hit {
		return cached, cached != nil, nil
	}

	path := filepath.Join(modelDir, "components", kind+".json")
	raw, readErr := os.ReadFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			l.cache[key] = nil
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("reading %s: %w", path, readErr)
	}

	var definition component.ComponentDefinition
	if err := json.Unmarshal(raw, &definition); err != nil {
		return nil, false, fmt.Errorf("parsing %s: %w", path, err)
	}
	if strings.TrimSpace(definition.Component.Schema) == "" {
		l.cache[key] = nil
		return nil, false, nil
	}

	parsed := map[string]any{}
	if err := json.Unmarshal([]byte(definition.Component.Schema), &parsed); err != nil {
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

// loadRelationshipDefinition decodes one relationship definition file.
func loadRelationshipDefinition(path string) (*relationship.RelationshipDefinition, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	definition := &relationship.RelationshipDefinition{}
	if err := json.Unmarshal(raw, definition); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return definition, nil
}

// selectorItems flattens a definition's allow.from and allow.to selector items.
func selectorItems(definition *relationship.RelationshipDefinition) []relationship.SelectorItem {
	if definition.Selectors == nil {
		return nil
	}
	items := []relationship.SelectorItem{}
	for _, selector := range *definition.Selectors {
		items = append(items, selector.Allow.From...)
		items = append(items, selector.Allow.To...)
	}
	return items
}

// mutationPaths returns every mutatorRef and mutatedRef path of a selector item.
func mutationPaths(item relationship.SelectorItem) [][]string {
	patch := item.RelationshipDefinitionSelectorsPatch
	if patch == nil {
		return nil
	}
	paths := [][]string{}
	if patch.MutatorRef != nil {
		paths = append(paths, *patch.MutatorRef...)
	}
	if patch.MutatedRef != nil {
		paths = append(paths, *patch.MutatedRef...)
	}
	return paths
}

func TestRelationshipMutationPathsResolveAgainstComponentSchemas(t *testing.T) {
	root := repoRoot(t)
	loader := newComponentSchemaLoader(root)

	unresolved := map[string]bool{}
	checked := 0
	scanned := 0

	for _, modelName := range relationshipConformanceModels {
		versionDirs, err := modelVersionDirs(root, modelName)
		if err != nil {
			t.Fatal(err)
		}
		if len(versionDirs) == 0 {
			t.Fatalf("model %q ships no version directories - the scope list is stale", modelName)
		}

		for _, modelDir := range versionDirs {
			files, err := filepath.Glob(filepath.Join(modelDir, "relationships", "*.json"))
			if err != nil {
				t.Fatalf("globbing relationships in %s: %v", modelDir, err)
			}
			scanned += len(files)

			for _, file := range files {
				definition, err := loadRelationshipDefinition(file)
				if err != nil {
					t.Fatal(err)
				}

				for _, item := range selectorItems(definition) {
					if item.Kind == nil || *item.Kind == "" || *item.Kind == "*" || item.Model == nil {
						continue
					}
					kind := *item.Kind

					targetDir := modelDir
					if item.Model.Name != modelName {
						targetDir, err = loader.newestVersionDir(item.Model.Name)
						if err != nil {
							t.Fatal(err)
						}
						if targetDir == "" {
							continue
						}
					}

					for _, path := range mutationPaths(item) {
						segments, checkable := schemaCheckedSegments(path)
						if !checkable {
							continue
						}
						schema, ok, err := loader.schema(targetDir, kind)
						if err != nil {
							t.Fatalf("%s: %v", file, err)
						}
						if !ok {
							continue
						}
						checked++
						if resolveSchemaPath(schema, segments) {
							continue
						}
						key := fmt.Sprintf("%s/%s|%s|%s", modelName, filepath.Base(file), kind, strings.Join(path, "."))
						unresolved[key] = true
						if !knownUnresolvedMutationPaths[key] {
							relFile, relErr := filepath.Rel(root, file)
							if relErr != nil {
								relFile = file
							}
							t.Errorf("%s: %s mutates %q, which the %s component schema shipped in %s does not define",
								filepath.ToSlash(relFile), kind, strings.Join(path, "."), kind, filepath.Base(filepath.Dir(targetDir)))
						}
					}
				}
			}
		}
	}

	if scanned == 0 {
		t.Fatal("no relationship definitions were scanned - the scope list is stale")
	}
	if checked == 0 {
		t.Fatal("no configuration-rooted mutation paths were checked - the scope list is stale")
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

// TestIngressToServiceEdgeUsesNetworkingV1BackendShape pins meshery/meshery#21482
// in the artifact the fix changed: the Ingress -> Service edge must mutate the
// networking.k8s.io/v1 backend shape, in every kubernetes version directory the
// definition is fanned out to. patchMutatorsAction pairs mutatorRefs[i] with
// mutatedRefs[i], so both sides are pinned in order - reordering either one alone
// would write the Service port into backend.service.name.
func TestIngressToServiceEdgeUsesNetworkingV1BackendShape(t *testing.T) {
	root := repoRoot(t)

	wantMutated := [][]string{
		{"configuration", "spec", "rules", "0", "http", "paths", "0", "backend", "service", "name"},
		{"configuration", "spec", "rules", "0", "http", "paths", "0", "backend", "service", "port", "number"},
	}
	wantMutator := [][]string{
		{"displayName"},
		{"configuration", "spec", "ports", "0", "port"},
	}

	versionDirs, err := modelVersionDirs(root, "kubernetes")
	if err != nil {
		t.Fatal(err)
	}
	if len(versionDirs) == 0 {
		t.Fatal("kubernetes ships no version directories")
	}

	covered := 0
	for _, modelDir := range versionDirs {
		file := filepath.Join(modelDir, "relationships", "edge-non-binding-network-jccsr.json")
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		covered++

		definition, err := loadRelationshipDefinition(file)
		if err != nil {
			t.Fatal(err)
		}

		gotMutated := [][]string{}
		gotMutator := [][]string{}
		for _, item := range selectorItems(definition) {
			switch {
			case item.Kind == nil:
			case *item.Kind == "Ingress":
				gotMutated = append(gotMutated, mutationPaths(item)...)
			case *item.Kind == "Service":
				gotMutator = append(gotMutator, mutationPaths(item)...)
			}
		}

		version := filepath.Base(filepath.Dir(modelDir))
		assertPathsInOrder(t, version, "Ingress mutated", gotMutated, wantMutated)
		assertPathsInOrder(t, version, "Service mutator", gotMutator, wantMutator)
	}

	if covered == 0 {
		t.Fatal("edge-non-binding-network-jccsr.json is absent from every kubernetes version directory")
	}
}

// assertPathsInOrder compares a selector's declared mutation paths against the
// expected list positionally, because the engine pairs the two sides by index.
func assertPathsInOrder(t *testing.T, version, side string, got, want [][]string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("kubernetes %s: %s selector declares %d mutation paths, want %d (%v)", version, side, len(got), len(want), got)
		return
	}
	for i, path := range want {
		if strings.Join(got[i], ".") != strings.Join(path, ".") {
			t.Errorf("kubernetes %s: %s path %d is %q, want %q",
				version, side, i, strings.Join(got[i], "."), strings.Join(path, "."))
		}
	}
}
