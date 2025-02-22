package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"helm.sh/helm/v3/pkg/chart/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func main() {
	// filePath := "" // Define filePath

	// Test for Kubernetes manifest
	identifyYAML("clusterlogging.yaml")

	// Test for Helm chart
	identifyYAML("Chart.yaml")

	// Test for Kustomize file
	identifyYAML("kustomization.yaml")

	// Test for Docker Compose file
	identifyYAML("docker-compose.yaml")

	// Test for Meshery Design file
	identifyYAML("weaveworks-gitops-enterprise-design.yaml")

	// switch detectYAML(filePath) {

	// case "DockerCompose":
	// 	fmt.Println("Docker Compose file detected.")
	// case "Kubernetes":
	// 	fmt.Println("Kubernetes manifest detected.")
	// case "Helm":
	// 	fmt.Println("Helm chart detected.")
	// case "Kustomize":
	// 	fmt.Println("Kustomize file detected.")
	// default:
	// 	fmt.Println("Unknown YAML format.")
	// }
}

func identifyYAML(filePath string) string {
	var fileContent []byte
	var err error
	fileContent, err = ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "" // Return an empty string on error
	}
	// Convert fileContent to string for debugging
	// content := string(fileContent)
	// fmt.Printf("File content: %s\n", content)

	// Check if it's a Kubernetes manifest
	if isKubernetesManifest(fileContent) {
		return "Kubernetes"
	}

	// Check if it's a Helm chart
	if isHelmChart(filePath) {
		return "HelmChart"
	}

	// Check if it's a Kustomize file
	if isKustomizeFile(filePath) {
		return "Kustomize"
	}

	// Check if it's a Meshery Design file
	if isMesheryDesign(filePath) {
		return "MesheryDesign"
	}

	// // Check if it's a Meshery Design file
	// if isDockerCompose(filePath) {
	// 	return "DockerCompose"
	// }

	// If none of the above, assume it's a Docker Compose file
	// (especially if it has a.yml or.yaml extension and is in the project root)
	return "default"
}

func isKubernetesManifest(fileContent []byte) bool {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(fileContent), 4096)

	var typeMeta metav1.TypeMeta
	if err := decoder.Decode(&typeMeta); err != nil {
		return false
	}

	// Check for required Kubernetes fields
	if typeMeta.APIVersion != "" && typeMeta.Kind != "" {
		// Print the identified API version and kind
		fmt.Printf("Kubernetes Manifest Identified: APIVersion=%s, Kind=%s\n", typeMeta.APIVersion, typeMeta.Kind) // Print for debugging
		return true                                                                                                // Return true if valid Kubernetes manifest
	}
	return false
}

func isHelmChart(filePath string) bool {
	// Verification Criteria:
	// 1. Check if there is a Chart.yaml file in the directory
	// 2. Attempt to load the chart using loader.Load
	chartPath := filepath.Dir(filePath)
	if _, err := os.Stat(filepath.Join(chartPath, "Chart.yaml")); err == nil {
		// Use loader.Load instead of chart.Load
		_, err := loader.Load(chartPath)
		if err == nil {
			fmt.Println("Helm Chart Identified") // Print for debugging
			return true
		}
	}
	return false
}

func isKustomizeFile(filePath string) bool {
	// Verification Criteria:
	// 1. Check if the filename is kustomization.yaml or kustomization.yml
	// 2. Creates a krusty.Kustomizer instance with default options using krusty.MakeDefaultOptions().
	// 2.a. Calls the k.Run() method with the file system and file path to attempt to load and process the file as a Kustomize file.
	// 2.b. If k.Run() returns an error, it indicates that the file is not a valid Kustomize file. Otherwise, it's considered a Kustomize file.
	if filepath.Base(filePath) == "kustomization.yaml" || filepath.Base(filePath) == "kustomization.yml" {
		// Attempt to decode the kustomization file
		fSys := filesys.MakeFsOnDisk()
		k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
		_, err := k.Run(fSys, filepath.Dir(filePath))
		if err == nil {
			fmt.Println("Kustomize File Identified") // Print for debugging
			return true
		}
	}
	return false
}

func isMesheryDesign(filePath string) bool {
	// Check if the file is a Meshery Design file
	// Verification Criteria:
	// 1. Check if the file contains required Meshery Design fields
	// 2. Validate against Meshery Design schema v1beta1
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	var design struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
		Spec struct {
			Pattern map[string]interface{} `yaml:"pattern"`
		} `yaml:"spec"`
	}

	decoder := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(content), 4096)
	if err := decoder.Decode(&design); err != nil {
		return false
	}

	// Validate required fields for Meshery Design v1beta1
	if design.APIVersion != "design.meshery.io/v1beta1" ||
		design.Kind != "Design" ||
		design.Metadata.Name == "" ||
		design.Spec.Pattern == nil {
		return false
	}
	fmt.Println("Meshery Design Identified") // Print for debugging
	return true
}

func TestIdentifyYAML(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"clusterlogging.yaml", "Kubernetes"},
		// {"weaveworks-gitops-enterprise-design.yaml", "HelmChart"},
		// {"kustomization.yaml", "Kustomize"},
		// {"docker-compose.yaml", "default"},                   // Assuming this is the expected output
	}

	for _, test := range tests {
		result := identifyYAML(test.fileName)
		if result != test.expected {
			t.Errorf("identifyYAML(%s) = %s; want %s", test.fileName, result, test.expected)
		}
	}
}
