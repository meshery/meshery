// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"encoding/json"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils/format"
	"gopkg.in/yaml.v2"
)

// TestVersionCmdFlags verifies that the output-format flag and shorthand are registered.
func TestVersionCmdFlags(t *testing.T) {
	flag := versionCmd.Flags().Lookup("output-format")
	if flag == nil {
		t.Fatal("flag --output-format should be registered")
	}
	if flag.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got '%s'", flag.Shorthand)
	}
	if flag.DefValue != "" {
		t.Errorf("expected default value '', got '%s'", flag.DefValue)
	}
}

// TestVersionCmdFlagParsing tests parsing valid output-format flags.
func TestVersionCmdFlagParsing(t *testing.T) {
	t.Cleanup(func() {
		_ = versionCmd.Flags().Set("output-format", "")
	})

	if err := versionCmd.Flags().Set("output-format", "json"); err != nil {
		t.Fatalf("unexpected error setting json format: %v", err)
	}
	if outputFormat != "json" {
		t.Errorf("expected outputFormat 'json', got '%s'", outputFormat)
	}

	if err := versionCmd.Flags().Set("output-format", "yaml"); err != nil {
		t.Fatalf("unexpected error setting yaml format: %v", err)
	}
	if outputFormat != "yaml" {
		t.Errorf("expected outputFormat 'yaml', got '%s'", outputFormat)
	}
}

// TestVersionOutputSchema validates serialization using production types.
func TestVersionOutputSchema(t *testing.T) {
	dummy := versionOutput{
		Client: versionInfo{Version: "v1.0.69", GitSHA: "7b53ef517e3c632a3bee63fa7c21b031430943ab"},
		Server: versionInfo{Version: "unavailable", GitSHA: "unavailable"},
	}

	jsonData, err := json.Marshal(dummy)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}

	var jsonMap map[string]map[string]string
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}
	if _, ok := jsonMap["client"]; !ok {
		t.Error("expected key 'client' in json map")
	}
	if _, ok := jsonMap["server"]; !ok {
		t.Error("expected key 'server' in json map")
	}
	if jsonMap["client"]["version"] != "v1.0.69" {
		t.Errorf("expected client version 'v1.0.69', got '%s'", jsonMap["client"]["version"])
	}
	if jsonMap["client"]["git_sha"] != "7b53ef517e3c632a3bee63fa7c21b031430943ab" {
		t.Errorf("expected client git_sha '7b53ef517e3c632a3bee63fa7c21b031430943ab', got '%s'", jsonMap["client"]["git_sha"])
	}
	if jsonMap["server"]["version"] != "unavailable" {
		t.Errorf("expected server version 'unavailable', got '%s'", jsonMap["server"]["version"])
	}
	if jsonMap["server"]["git_sha"] != "unavailable" {
		t.Errorf("expected server git_sha 'unavailable', got '%s'", jsonMap["server"]["git_sha"])
	}

	yamlData, err := yaml.Marshal(dummy)
	if err != nil {
		t.Fatalf("failed to marshal yaml: %v", err)
	}

	var yamlMap map[string]map[string]string
	if err := yaml.Unmarshal(yamlData, &yamlMap); err != nil {
		t.Fatalf("failed to unmarshal yaml: %v", err)
	}
	if _, ok := yamlMap["client"]; !ok {
		t.Error("expected key 'client' in yaml map")
	}
	if _, ok := yamlMap["server"]; !ok {
		t.Error("expected key 'server' in yaml map")
	}
	if yamlMap["client"]["version"] != "v1.0.69" {
		t.Errorf("expected client version 'v1.0.69', got '%s'", yamlMap["client"]["version"])
	}
	if yamlMap["client"]["git_sha"] != "7b53ef517e3c632a3bee63fa7c21b031430943ab" {
		t.Errorf("expected git_sha '7b53ef517e3c632a3bee63fa7c21b031430943ab', got '%s'", yamlMap["client"]["git_sha"])
	}
	if yamlMap["server"]["version"] != "unavailable" {
		t.Errorf("expected server version 'unavailable', got '%s'", yamlMap["server"]["version"])
	}
	if yamlMap["server"]["git_sha"] != "unavailable" {
		t.Errorf("expected server git_sha 'unavailable', got '%s'", yamlMap["server"]["git_sha"])
	}
}

// TestVersionFormatOutput verifies format helpers.
func TestVersionFormatOutput(t *testing.T) {
	dummy := versionOutput{
		Client: versionInfo{Version: "v1.0.69", GitSHA: "7b53ef517e3c632a3bee63fa7c21b031430943ab"},
		Server: versionInfo{Version: "unavailable", GitSHA: "unavailable"},
	}

	if err := format.OutputJson(dummy); err != nil {
		t.Errorf("unexpected error in OutputJson: %v", err)
	}
	if err := format.OutputYaml(dummy); err != nil {
		t.Errorf("unexpected error in OutputYaml: %v", err)
	}
}
