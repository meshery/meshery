package format

import (
	"testing"
)

type sampleData struct {
	Name string   `json:"name" yaml:"name" toon:"name"`
	Type string   `json:"type" yaml:"type" toon:"type"`
	Tags []string `json:"tags" yaml:"tags" toon:"tags"`
}

func TestOutputFormatters(t *testing.T) {
	data := sampleData{
		Name: "test-component",
		Type: "sample",
		Tags: []string{"tag1", "tag2"},
	}

	t.Run("OutputJson", func(t *testing.T) {
		err := OutputJson(data)
		if err != nil {
			t.Fatalf("OutputJson failed: %v", err)
		}
	})

	t.Run("OutputYaml", func(t *testing.T) {
		err := OutputYaml(data)
		if err != nil {
			t.Fatalf("OutputYaml failed: %v", err)
		}
	})

	t.Run("OutputToon", func(t *testing.T) {
		err := OutputToon(data)
		if err != nil {
			t.Fatalf("OutputToon failed: %v", err)
		}
	})
}
