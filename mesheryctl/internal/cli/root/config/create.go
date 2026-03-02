package config

import (
	"os"
)

func CreateConfigFile(path string) error {
	// Create file with minimal valid YAML structure
	initialConfig := `contexts:
tokens:
current-context: ""
`
	err := os.WriteFile(path, []byte(initialConfig), 0644)
	if err != nil {
		return err
	}
	return nil
}
