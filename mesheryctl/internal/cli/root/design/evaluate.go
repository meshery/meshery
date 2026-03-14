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

package design

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models/pattern/core"
	encoding "github.com/meshery/meshkit/encoding"
	meshkitutils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type evaluateDesignFlag struct {
	File         string `json:"file"          validate:"omitempty"`
	OutputFile   string `json:"output"        validate:"required"`
	OutputFormat string `json:"output-format" validate:"required,oneof=json yaml"`
}

var evaluateFlags = &evaluateDesignFlag{}

var evaluateCmd = &cobra.Command{
	Use:   "evaluate [ID]",
	Short: "Evaluate a design",
	Long: `Evaluate a design by running relationship evaluation policies.
The evaluated design is saved to the specified output file while an overview
of evaluation actions is printed to the terminal.`,
	Example: `
// Evaluate a design from a file and save the result
mesheryctl design evaluate -f design.yaml -o evaluated-design.yaml

// Evaluate a design by ID
mesheryctl design evaluate 12345678-abcd-efgh-ijkl-123456789012 -o result.yaml

// Evaluate and save as JSON
mesheryctl design evaluate -f design.yaml --output-format json -o evaluated-design.json
	`,
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := mesheryctlflags.ValidateCmdFlags(cmd, evaluateFlags); err != nil {
			return err
		}

		// Require exactly one of -f or positional ID
		if evaluateFlags.File == "" && len(args) == 0 {
			return ErrEvaluateDesign()
		}
		if evaluateFlags.File != "" && len(args) > 0 {
			return ErrEvaluateDesign()
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		baseURL := mctlCfg.GetBaseMesheryURL()

		var designPayload pattern.PatternFile

		if evaluateFlags.File != "" {
			designPayload, err = readDesignFromFile(evaluateFlags.File)
			if err != nil {
				return err
			}
		} else {
			designPayload, err = fetchDesignByID(baseURL, args[0])
			if err != nil {
				return err
			}
		}

		evalResponse, err := sendEvaluationRequest(designPayload)
		if err != nil {
			return err
		}

		if err := saveEvaluatedDesign(evalResponse.Design, evaluateFlags.OutputFile, evaluateFlags.OutputFormat); err != nil {
			return err
		}

		printActionsOverview(evalResponse)

		return nil
	},
}

func readDesignFromFile(filePath string) (pattern.PatternFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return pattern.PatternFile{}, utils.ErrFileRead(err)
	}

	pf, err := core.NewPatternFile(content)
	if err != nil {
		return pattern.PatternFile{}, ErrParseDesignFile(err)
	}
	return pf, nil
}

func sendEvaluationRequest(designPayload pattern.PatternFile) (*pattern.EvaluationResponse, error) {
	evalReq := pattern.EvaluationRequest{
		Design: designPayload,
	}

	payloadBytes, err := json.Marshal(evalReq)
	if err != nil {
		return nil, utils.ErrMarshal(err)
	}

	s := utils.CreateDefaultSpinner("Evaluating design", "")
	s.Start()
	resp, err := api.Add("api/meshmodels/relationships/evaluate", bytes.NewBuffer(payloadBytes), nil)
	s.Stop()
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.ErrReadFromBody(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrEvaluateDesignResponse(fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body)))
	}

	var evalResponse pattern.EvaluationResponse
	if err := json.Unmarshal(body, &evalResponse); err != nil {
		return nil, utils.ErrUnmarshal(err)
	}

	return &evalResponse, nil
}

func saveEvaluatedDesign(design pattern.PatternFile, outputFile string, format string) error {
	var content []byte
	var err error

	switch format {
	case "json":
		content, err = encoding.Marshal(design)
	case "yaml":
		content, err = yaml.Marshal(design)
	default:
		content, err = yaml.Marshal(design)
	}
	if err != nil {
		return utils.ErrMarshal(err)
	}

	if err := meshkitutils.WriteToFile(outputFile, string(content)); err != nil {
		return err
	}

	utils.Log.Infof("Evaluated design saved to %s", outputFile)
	return nil
}

func printActionsOverview(evalResponse *pattern.EvaluationResponse) {
	if len(evalResponse.Actions) == 0 {
		utils.Log.Info("No evaluation actions were taken. The design is already up to date.")
		return
	}

	if evalResponse.Timestamp != nil {
		utils.Log.Infof("Evaluation completed at %s", evalResponse.Timestamp.Format("2006-01-02 15:04:05 MST"))
	} else {
		utils.Log.Info("Evaluation completed")
	}
	utils.Log.Infof("Total actions: %d\n", len(evalResponse.Actions))

	// Group actions by operation type
	groups := groupActionsByOp(evalResponse.Actions)

	// Display order for action categories
	displayOrder := []struct {
		op    string
		label string
	}{
		{"add_component", "Components Added"},
		{"delete_component", "Components Removed"},
		{"update_component", "Components Updated"},
		{"update_component_configuration", "Component Configurations Updated"},
		{"add_relationship", "Relationships Added"},
		{"delete_relationship", "Relationships Removed"},
		{"update_relationship", "Relationships Updated"},
	}

	for _, entry := range displayOrder {
		actions, ok := groups[entry.op]
		if !ok || len(actions) == 0 {
			continue
		}
		printActionGroup(entry.label, entry.op, actions)
	}

	// Print any unknown operation types not covered above
	known := map[string]bool{
		"add_component":                  true,
		"delete_component":               true,
		"update_component":               true,
		"update_component_configuration": true,
		"add_relationship":               true,
		"delete_relationship":            true,
		"update_relationship":            true,
	}
	for op, actions := range groups {
		if !known[op] {
			printActionGroup(op, op, actions)
		}
	}
}

func groupActionsByOp(actions []pattern.Action) map[string][]pattern.Action {
	groups := make(map[string][]pattern.Action)
	for _, a := range actions {
		groups[a.Op] = append(groups[a.Op], a)
	}
	return groups
}

func printActionGroup(label, op string, actions []pattern.Action) {
	fmt.Printf("%s (%d)\n", label, len(actions))
	fmt.Println(strings.Repeat("-", 40))

	for _, action := range actions {
		switch op {
		case "add_component":
			printAddComponentAction(action)
		case "delete_component":
			printDeleteComponentAction(action)
		case "update_component", "update_component_configuration":
			printUpdateComponentAction(action)
		case "add_relationship":
			printAddRelationshipAction(action)
		case "delete_relationship":
			printDeleteRelationshipAction(action)
		case "update_relationship":
			printUpdateRelationshipAction(action)
		default:
			printGenericAction(action)
		}
	}
	fmt.Println()
}

func printAddComponentAction(action pattern.Action) {
	item := action.Value["item"]
	kind, displayName := extractComponentInfo(item)
	if kind != "" {
		fmt.Printf("  + %s %q\n", kind, displayName)
	} else {
		fmt.Printf("  + component (id: %s)\n", extractID(action.Value))
	}
}

func printDeleteComponentAction(action pattern.Action) {
	kind := extractStringField(action.Value, "component")
	id := extractID(action.Value)
	if kind != "" {
		fmt.Printf("  - %s (id: %s)\n", kind, id)
	} else {
		fmt.Printf("  - component (id: %s)\n", id)
	}
}

func printUpdateComponentAction(action pattern.Action) {
	id := extractID(action.Value)
	pathStr := formatActionPath(action.Value["path"])
	valStr := formatActionValue(action.Value["value"])
	if valStr != "" {
		fmt.Printf("  ~ component %s: %s -> %s\n", id, pathStr, valStr)
	} else {
		fmt.Printf("  ~ component %s: %s\n", id, pathStr)
	}
}

func printAddRelationshipAction(action pattern.Action) {
	item := action.Value["item"]
	kind, subType, from, to := extractRelationshipInfo(item)
	if kind != "" {
		fmt.Printf("  + %s-%s relationship from %s to %s\n", kind, subType, from, to)
	} else {
		fmt.Printf("  + relationship (id: %s)\n", extractID(action.Value))
	}
}

func printDeleteRelationshipAction(action pattern.Action) {
	relType := extractStringField(action.Value, "relationship")
	id := extractID(action.Value)
	if relType != "" {
		fmt.Printf("  - %s (id: %s)\n", relType, id)
	} else {
		fmt.Printf("  - relationship (id: %s)\n", id)
	}
}

func printUpdateRelationshipAction(action pattern.Action) {
	id := extractID(action.Value)
	pathStr := formatActionPath(action.Value["path"])
	valStr := formatActionValue(action.Value["value"])
	if valStr != "" {
		fmt.Printf("  ~ relationship %s: %s -> %s\n", id, pathStr, valStr)
	} else {
		fmt.Printf("  ~ relationship %s: %s\n", id, pathStr)
	}
}

func printGenericAction(action pattern.Action) {
	id := extractID(action.Value)
	if id != "" {
		fmt.Printf("  ? %s (id: %s)\n", action.Op, id)
	} else {
		fmt.Printf("  ? %s\n", action.Op)
	}
}

// extractID returns the "id" string from an action's Value map.
func extractID(value map[string]interface{}) string {
	id, _ := value["id"].(string)
	return id
}

// extractStringField returns a named string field from the value map.
func extractStringField(value map[string]interface{}, key string) string {
	s, _ := value[key].(string)
	return s
}

// extractComponentInfo pulls kind and displayName from an add_component action's "item" field.
func extractComponentInfo(item interface{}) (kind, displayName string) {
	m, ok := item.(map[string]interface{})
	if !ok {
		return "", ""
	}
	if comp, ok := m["component"].(map[string]interface{}); ok {
		kind, _ = comp["kind"].(string)
	}
	displayName, _ = m["displayName"].(string)
	if displayName == "" {
		displayName, _ = m["display_name"].(string)
	}
	return kind, displayName
}

// extractRelationshipInfo pulls kind, subType, from, and to from an add_relationship action's "item" field.
func extractRelationshipInfo(item interface{}) (kind, subType, from, to string) {
	m, ok := item.(map[string]interface{})
	if !ok {
		return "", "", "unknown", "unknown"
	}
	kind, _ = m["kind"].(string)
	subType, _ = m["subType"].(string)
	if subType == "" {
		subType, _ = m["sub_type"].(string)
	}
	from = "unknown"
	to = "unknown"

	selectors, ok := m["selectors"].([]interface{})
	if !ok || len(selectors) == 0 {
		return kind, subType, from, to
	}
	selector, ok := selectors[0].(map[string]interface{})
	if !ok {
		return kind, subType, from, to
	}
	allow, ok := selector["allow"].(map[string]interface{})
	if !ok {
		return kind, subType, from, to
	}

	if fromArr, ok := allow["from"].([]interface{}); ok && len(fromArr) > 0 {
		if fromItem, ok := fromArr[0].(map[string]interface{}); ok {
			if k, ok := fromItem["kind"].(string); ok {
				from = k
			}
		}
	}
	if toArr, ok := allow["to"].([]interface{}); ok && len(toArr) > 0 {
		if toItem, ok := toArr[0].(map[string]interface{}); ok {
			if k, ok := toItem["kind"].(string); ok {
				to = k
			}
		}
	}
	return kind, subType, from, to
}

// formatActionPath converts an action path (string or []interface{}) to a dot-separated string.
func formatActionPath(raw interface{}) string {
	switch v := raw.(type) {
	case string:
		return v
	case []interface{}:
		parts := make([]string, 0, len(v))
		for _, p := range v {
			parts = append(parts, fmt.Sprintf("%v", p))
		}
		return strings.Join(parts, ".")
	default:
		return fmt.Sprintf("%v", raw)
	}
}

// formatActionValue returns a quoted string representation of a value for display.
func formatActionValue(raw interface{}) string {
	if raw == nil {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func init() {
	evaluateCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	evaluateCmd.Flags().StringVarP(&evaluateFlags.File, "file", "f", "", "Path to design file")
	evaluateCmd.Flags().StringVarP(&evaluateFlags.OutputFile, "output", "o", "", "Path to save the evaluated design")
	evaluateCmd.Flags().StringVar(&evaluateFlags.OutputFormat, "output-format", "yaml", "Output format for the evaluated design [json|yaml]")
}
