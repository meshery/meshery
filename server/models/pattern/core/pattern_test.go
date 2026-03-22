package core

import (
	"testing"
)

func TestDesignNameFromFileName(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		expectedName string
	}{
		{
			name:         "regular yaml file",
			fileName:     "my-deployment.yaml",
			expectedName: "my-deployment",
		},
		{
			name:         "regular yml file",
			fileName:     "my-deployment.yml",
			expectedName: "my-deployment",
		},
		{
			name:         "tar.gz file",
			fileName:     "my-chart.tar.gz",
			expectedName: "my-chart",
		},
		{
			name:         "json file",
			fileName:     "config.json",
			expectedName: "config",
		},
		{
			name:         "file with multiple dots",
			fileName:     "my.k8s.deployment.yaml",
			expectedName: "my.k8s.deployment",
		},
		{
			name:         "empty filename",
			fileName:     "",
			expectedName: "",
		},
		{
			name:         "filename without extension",
			fileName:     "mydesign",
			expectedName: "mydesign",
		},
		{
			name:         "tgz file",
			fileName:     "helm-chart.tgz",
			expectedName: "helm-chart",
		},
		{
			name:         "extension only",
			fileName:     ".yaml",
			expectedName: "",
		},
		{
			name:         "tar.gz extension only",
			fileName:     ".tar.gz",
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DesignNameFromFileName(tt.fileName)
			if result != tt.expectedName {
				t.Errorf("expected %q, got %q", tt.expectedName, result)
			}
		})
	}
}
