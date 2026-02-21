package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

func TestApplyDefaultsToConfigurationSimpleDefault(t *testing.T) {
	schema := `{
		"type":"object",
		"properties":{
			"name":{"type":"string"},
			"age":{"type":"integer","default":25}
		}
	}`

	cfg := map[string]interface{}{
		"name": "meshery",
	}

	if err := ApplyDefaultsToConfiguration(schema, &cfg); err != nil {
		t.Fatalf("unexpected error while applying defaults: %v", err)
	}

	ageVal, ok := cfg["age"].(float64)
	if !ok {
		t.Fatalf("expected default age to be set")
	}

	if ageVal != 25 {
		t.Fatalf("expected age default to be 25, got %v", ageVal)
	}
}

func TestApplyDefaultsToConfigurationNestedArrayDefault(t *testing.T) {
	schema := `{
		"type":"object",
		"properties":{
			"spec":{
				"type":"object",
				"properties":{
					"ports":{
						"type":"array",
						"items":{
							"allOf":[
								{
									"type":"object",
									"properties":{
										"port":{"type":"integer"},
										"protocol":{"type":"string","default":"TCP"}
									}
								}
							]
						}
					}
				}
			}
		}
	}`

	cfg := map[string]interface{}{
		"spec": map[string]interface{}{
			"ports": []interface{}{
				map[string]interface{}{
					"port": 8080,
				},
			},
		},
	}

	if err := ApplyDefaultsToConfiguration(schema, &cfg); err != nil {
		t.Fatalf("unexpected error while applying defaults: %v", err)
	}

	spec, ok := cfg["spec"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected spec to be an object")
	}

	ports, ok := spec["ports"].([]interface{})
	if !ok || len(ports) != 1 {
		t.Fatalf("expected one port entry")
	}

	portDef, ok := ports[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected port definition to be object")
	}

	protocol, ok := portDef["protocol"].(string)
	if !ok {
		t.Fatalf("expected protocol default to be set")
	}

	if protocol != "TCP" {
		t.Fatalf("expected protocol to be TCP, got %s", protocol)
	}
}

func TestApplyDefaultsToConfigurationDoesNotOverride(t *testing.T) {
	schema := `{
		"type":"object",
		"properties":{
			"spec":{
				"type":"object",
				"properties":{
					"ports":{
						"type":"array",
						"items":{
							"type":"object",
							"properties":{
								"protocol":{"type":"string","default":"TCP"}
							}
						}
					}
				}
			}
		}
	}`

	cfg := map[string]interface{}{
		"spec": map[string]interface{}{
			"ports": []interface{}{
				map[string]interface{}{
					"protocol": "UDP",
				},
			},
		},
	}

	if err := ApplyDefaultsToConfiguration(schema, &cfg); err != nil {
		t.Fatalf("unexpected error while applying defaults: %v", err)
	}

	spec, _ := cfg["spec"].(map[string]interface{})
	ports, _ := spec["ports"].([]interface{})
	portDef, _ := ports[0].(map[string]interface{})

	protocol, _ := portDef["protocol"].(string)
	if protocol != "UDP" {
		t.Fatalf("expected protocol to remain UDP, got %s", protocol)
	}
}

func TestApplyDefaultsToPatternComponents(t *testing.T) {
	compID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to create component id: %v", err)
	}

	pt := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{
			{
				Id:          compID,
				DisplayName: "service",
				Component: component.Component{
					Kind: "Service",
					Schema: `{
						"type":"object",
						"properties":{
							"spec":{
								"type":"object",
								"properties":{
									"ports":{
										"type":"array",
										"items":{
											"type":"object",
											"properties":{
												"protocol":{"type":"string","default":"TCP"}
											}
										}
									}
								}
							}
						}
					}`,
				},
				Configuration: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{},
						},
					},
				},
			},
		},
	}

	errs := ApplyDefaultsToPatternComponents(pt)
	if len(errs) > 0 {
		t.Fatalf("expected defaults hydration for pattern to succeed, got errors: %v", errs)
	}

	compCfg := pt.Components[0].Configuration
	spec, ok := compCfg["spec"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected spec to be object")
	}
	ports, ok := spec["ports"].([]interface{})
	if !ok || len(ports) != 1 {
		t.Fatalf("expected one port in configuration")
	}
	portDef, ok := ports[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected port definition to be object")
	}

	protocol, _ := portDef["protocol"].(string)
	if protocol != "TCP" {
		t.Fatalf("expected protocol default TCP, got %q", protocol)
	}
}

func TestApplyDefaultsToPatternComponentsWithRealKubernetesServiceFixture(t *testing.T) {
	serviceFixture := loadKubernetesServiceFixture(t)

	compID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to create component id: %v", err)
	}

	serviceFixture.Id = compID
	serviceFixture.DisplayName = "real-service-fixture"
	serviceFixture.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"ports": []interface{}{
				map[string]interface{}{
					"port":       8080,
					"targetPort": 8080,
				},
			},
		},
	}

	pt := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{&serviceFixture},
	}

	errs := ApplyDefaultsToPatternComponents(pt)
	if len(errs) > 0 {
		t.Fatalf("expected defaults hydration for real fixture to succeed, got errors: %v", errs)
	}

	spec, ok := pt.Components[0].Configuration["spec"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected spec to be object")
	}
	ports, ok := spec["ports"].([]interface{})
	if !ok || len(ports) != 1 {
		t.Fatalf("expected one port in configuration")
	}
	portDef, ok := ports[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected port definition to be object")
	}

	protocol, _ := portDef["protocol"].(string)
	if protocol != "TCP" {
		t.Fatalf("expected protocol default TCP from real fixture, got %q", protocol)
	}
}

func loadKubernetesServiceFixture(t *testing.T) component.ComponentDefinition {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to resolve current test file path")
	}

	// defaults_test.go is in server/models/pattern/utils.
	// Move to server/ and read real meshmodel fixtures.
	serverRoot := filepath.Join(filepath.Dir(currentFile), "..", "..", "..")
	pattern := filepath.Join(serverRoot, "meshmodel", "kubernetes", "*", "v1.0.0", "components", "Service.json")

	paths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to enumerate kubernetes service fixtures: %v", err)
	}
	if len(paths) == 0 {
		t.Fatalf("no kubernetes service fixture found in meshmodel directory")
	}

	sort.Strings(paths)

	var lastErr error
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			lastErr = err
			continue
		}

		var comp component.ComponentDefinition
		if err := json.Unmarshal(content, &comp); err != nil {
			lastErr = err
			continue
		}

		if strings.TrimSpace(comp.Component.Schema) == "" {
			lastErr = fmt.Errorf("component schema is empty in fixture: %s", path)
			continue
		}

		// Ensure we pick a schema fixture that actually defines a protocol default.
		if strings.Contains(comp.Component.Schema, `"protocol":{"default":"TCP"`) {
			return comp
		}
	}

	t.Fatalf("failed to load kubernetes service fixture with protocol default: %v", lastErr)
	return component.ComponentDefinition{}
}
