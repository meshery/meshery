package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"
)

var (
	modelName    string
	outputFormat string
	version      string
	pathDir      string
)

var modelInitCmd = &cobra.Command{
	Use:   "init [model-name]",
	Short: "Initialize a new Meshery model",
	Long:  `Creates a scaffold directory structure for a new Meshery model with optional flags for format, version and path.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("model name is required")
		}
		modelName = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		modelPath := filepath.Join(pathDir, modelName, version)
		dirs := []string{
			modelPath,
			filepath.Join(modelPath, "components"),
			filepath.Join(modelPath, "relationships"),
			filepath.Join(modelPath, "connections"),
			filepath.Join(modelPath, "credentials"),
		}

		for _, dir := range dirs {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}

		if outputFormat == "yaml" {
			if err := generateOrderedYAMLModel(modelPath); err != nil {
				return err
			}
		} else if outputFormat == "json" {
			if err := generateJSONModel(modelPath); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}

		cmd.Printf("Created new Meshery model: %s at %s\n", modelName, modelPath)
		return nil
	},
}

func init() {
	ModelCmd.AddCommand(modelInitCmd)
	modelInitCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "yaml", "Format of generated model [json|yaml]")
	modelInitCmd.Flags().StringVar(&version, "version", "1.0.0", "Version of the model")
	modelInitCmd.Flags().StringVarP(&pathDir, "path", "p", ".", "Path where model should be initialized")
}

func generateOrderedYAMLModel(outputPath string) error {
	id := uuid.New().String()

	om := orderedmap.New[string, any]()
	om.Set("id", id)
	om.Set("schemaVersion", "v1beta1")
	om.Set("version", version)
	om.Set("name", modelName)
	om.Set("displayName", "My Model Display Name")
	om.Set("description", "This is an example description.")
	om.Set("status", "enabled")
	om.Set("registrant", map[string]string{"id": "<connection-id>"})
	om.Set("registrantId", "<connection-id>")
	om.Set("categoryId", "<category-id>")
	om.Set("category", map[string]string{"name": "Networking"})
	om.Set("subCategory", map[string]string{"name": "CNI"})
	om.Set("metadata", map[string]any{
		"svgWhite": "<svg-white>",
		"svgColor": "<svg-color>",
	})
	om.Set("model", map[string]any{})
	om.Set("relationships", []any{})
	om.Set("components", []any{})
	om.Set("componentsCount", 0)
	om.Set("relationshipsCount", 0)

	bytes, err := orderedMapToYAML(om)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputPath, "model.yaml"), bytes, 0644)
}

func orderedMapToYAML(m *orderedmap.OrderedMap[string, any]) ([]byte, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}
	for pair := m.Oldest(); pair != nil; pair = pair.Next() {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: pair.Key}
		valNode := &yaml.Node{}
		if err := valNode.Encode(pair.Value); err != nil {
			return nil, err
		}
		node.Content = append(node.Content, keyNode, valNode)
	}
	return yaml.Marshal(node)
}

func generateJSONModel(outputPath string) error {
	model := map[string]any{
		"id":                 uuid.New().String(),
		"schemaVersion":      "v1beta1",
		"version":            version,
		"name":               modelName,
		"displayName":        "My Model Display Name",
		"description":        "This is an example description.",
		"status":             "enabled",
		"registrant":         map[string]string{"id": "<connection-id>"},
		"registrantId":       "<connection-id>",
		"categoryId":         "<category-id>",
		"category":           map[string]string{"name": "Networking"},
		"subCategory":        map[string]string{"name": "CNI"},
		"metadata":           map[string]string{"svgWhite": "<svg-white>", "svgColor": "<svg-color>"},
		"model":              map[string]any{},
		"relationships":      []any{},
		"components":         []any{},
		"componentsCount":    0,
		"relationshipsCount": 0,
	}

	bytes, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputPath, "model.json"), bytes, 0644)
}
