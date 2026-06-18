package core

import (
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	v1beta3comp "github.com/meshery/schemas/models/v1beta3/component"
)

func TestDesignNameFromFileName(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		expectedName string
	}{
		{
			name:         "given regular yaml file when DesignNameFromFileName then return name",
			fileName:     "my-deployment.yaml",
			expectedName: "my-deployment",
		},
		{
			name:         "given regular yml file when DesignNameFromFileName then return name",
			fileName:     "my-deployment.yml",
			expectedName: "my-deployment",
		},
		{
			name:         "given tar.gz file when DesignNameFromFileName then strip compound extension",
			fileName:     "my-chart.tar.gz",
			expectedName: "my-chart",
		},
		{
			name:         "given json file when DesignNameFromFileName then return name",
			fileName:     "config.json",
			expectedName: "config",
		},
		{
			name:         "given filename with multiple dots when DesignNameFromFileName then strip only last extension",
			fileName:     "my.k8s.deployment.yaml",
			expectedName: "my.k8s.deployment",
		},
		{
			name:         "given empty filename when DesignNameFromFileName then return empty",
			fileName:     "",
			expectedName: "",
		},
		{
			name:         "given filename without extension when DesignNameFromFileName then return as is",
			fileName:     "mydesign",
			expectedName: "mydesign",
		},
		{
			name:         "given tgz file when DesignNameFromFileName then strip extension",
			fileName:     "helm-chart.tgz",
			expectedName: "helm-chart",
		},
		{
			name:         "given extension only when DesignNameFromFileName then return empty",
			fileName:     ".yaml",
			expectedName: "",
		},
		{
			name:         "given tar.gz extension only when DesignNameFromFileName then return empty",
			fileName:     ".tar.gz",
			expectedName: "",
		},
		{
			name:         "given unsupported xls extension when DesignNameFromFileName then still strip extension",
			fileName:     "spreadsheet.xls",
			expectedName: "spreadsheet",
		},
		{
			name:         "given unsupported zip extension when DesignNameFromFileName then still strip extension",
			fileName:     "archive.zip",
			expectedName: "archive",
		},
		{
			name:         "given unsupported exe extension when DesignNameFromFileName then still strip extension",
			fileName:     "installer.exe",
			expectedName: "installer",
		},
		{
			name:         "given unsupported tar extension when DesignNameFromFileName then strip last suffix only",
			fileName:     "archive.tar",
			expectedName: "archive",
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

func TestNewPatternFileFromK8sManifestIgnoreErrors(t *testing.T) {
	registryManager := newTestRegistryManager(t)

	tests := []struct {
		name               string
		manifest           string
		ignoreErrors       bool
		expectedCode       string
		expectedComponents int
	}{
		{
			name: "given unresolved manifest and ignoreErrors true when NewPatternFileFromK8sManifest then return no resolved components error",
			manifest: `
apiVersion: v1
kind: Service
metadata:
  name: unresolved-service
`,
			ignoreErrors:       true,
			expectedCode:       ErrNoResolvedComponentsCode,
			expectedComponents: 0,
		},
		{
			name: "given unresolved manifest and ignoreErrors false when NewPatternFileFromK8sManifest then return no resolved components error",
			manifest: `
apiVersion: v1
kind: Service
metadata:
  name: unresolved-service
`,
			ignoreErrors: false,
			expectedCode: ErrNoResolvedComponentsCode,
		},
		{
			name: "given empty yaml documents when NewPatternFileFromK8sManifest then return empty manifest error",
			manifest: `
---
---
`,
			ignoreErrors: true,
			expectedCode: ErrParseK8sManifestCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := NewPatternFileFromK8sManifest(tt.manifest, "manifest.yaml", tt.ignoreErrors, registryManager)

			if tt.expectedCode != "" && err == nil {
				t.Fatal("expected error, got nil")
			}

			if tt.expectedCode == "" && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			if tt.expectedCode != "" && errors.GetCode(err) != tt.expectedCode {
				t.Fatalf("expected error code %q, got %q", tt.expectedCode, errors.GetCode(err))
			}

			if len(pattern.Components) != tt.expectedComponents {
				t.Fatalf("expected %d resolved components, got %d", tt.expectedComponents, len(pattern.Components))
			}
		})
	}

	t.Run("given resolvable manifest when NewPatternFileFromK8sManifest then return resolved component", func(t *testing.T) {
		rm := newTestRegistryManager(t)
		seedDeploymentComponent(t, rm)

		pattern, err := NewPatternFileFromK8sManifest(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
`, "manifest.yaml", true, rm)

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		if len(pattern.Components) != 1 {
			t.Fatalf("expected 1 resolved component, got %d", len(pattern.Components))
		}

		if pattern.Components[0].DisplayName != "Deployment" {
			t.Fatalf("expected DisplayName 'Deployment', got %q", pattern.Components[0].DisplayName)
		}
	})
}

func seedDeploymentComponent(t *testing.T, rm *registry.RegistryManager) {
	t.Helper()

	conn := connection.Connection{
		Name:    "test-registrant",
		Kind:    "kubernetes",
		Type:    "platform",
		SubType: "orchestration",
		Status:  connection.ConnectionStatusConnected,
	}
	enabled := v1beta3comp.Enabled
	bgColor := "#123456"
	shape := core.Shape("round-rectangle")
	comp := v1beta3comp.ComponentDefinition{
		DisplayName:   "Deployment",
		SchemaVersion: "core.meshery.io/v1beta1",
		Status:        &enabled,
		Component: v1beta3comp.Component{
			Kind:    "Deployment",
			Version: "apps/v1",
			Schema:  `{"properties":{}}`,
		},
		Model: &model.ModelDefinition{
			Name:          "kubernetes",
			DisplayName:   "Kubernetes",
			SchemaVersion: "models.meshery.io/v1beta1",
			Version:       "v1.25.0",
			Model:         model.Model{Version: "v1.25.0"},
			Category:      category.CategoryDefinition{Name: "Orchestration"},
			Status:        model.Enabled,
		},
		Styles: &core.ComponentStyles{
			BackgroundColor: &bgColor,
			PrimaryColor:    "#123456",
			Shape:           &shape,
		},
	}
	id, err := comp.GenerateID()
	if err != nil {
		t.Fatalf("seedDeploymentComponent: GenerateID failed: %v", err)
	}
	comp.ID = id
	_, _, err = rm.RegisterEntity(conn, &comp)
	if err != nil {
		t.Fatalf("seedDeploymentComponent: RegisterEntity failed: %v", err)
	}
}

func newTestRegistryManager(t *testing.T) *registry.RegistryManager {
	t.Helper()

	db, err := database.New(database.Options{
		Filename: ":memory:",
		Engine:   "sqlite",
	})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	t.Cleanup(func() {
		if err := db.DBClose(); err != nil {
			t.Logf("failed to close database: %v", err)
		}
	})

	registryManager, err := registry.NewRegistryManager(&db)
	if err != nil {
		t.Fatalf("failed to create registry manager: %v", err)
	}

	return registryManager
}
