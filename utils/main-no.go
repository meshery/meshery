package maisn

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

var fileTypeCache = cache.New(5*time.Minute, 10*time.Minute)

func maisn() {
	// Command-line argument parsing
	filePath := flag.String("file", "", "Path to the YAML file")
	flag.Parse()

	if *filePath == "" {
		os.Exit(1) // Exit if no file is provided
	}

	content, err := os.ReadFile(*filePath)
	if err != nil {
		os.Exit(1) // Exit if file cannot be read
	}

	fileType := identifyYAMLTypeConcurrent(*filePath, content)
	// Output the identified file type
	println("Identified file type:", fileType)
}

func identifyYAMLTypeConcurrent(filePath string, content []byte) string {
	results := make(chan string, 4)

	// Run checks concurrently

	go func() {
		if isKubernetesManifest(content) {
			results <- "Kubernetes"
		}
	}()

	go func() {
		if isHelmChart(filePath) {
			results <- "Helm"
		}
	}()

	go func() {
		if isKustomize(filePath) {
			results <- "Kustomize"
		}
	}()

	// ... similar for Kubernetes and Docker Compose

	// Use select with timeout
	select {
	case result := <-results:
		return result
	case <-time.After(2 * time.Second):
		return "Unknown"
	}
}

func getCachedFileType(filePath string, content []byte) string {
	if cachedType, found := fileTypeCache.Get(filePath); found {
		return cachedType.(string)
	}

	fileType := identifyYAMLTypeConcurrent(filePath, content)
	fileTypeCache.Set(filePath, fileType, cache.DefaultExpiration)
	return fileType
}

func isKubernetesManifest(content []byte) bool {
	decoder := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(content), 4096)
	var obj runtime.Object
	if err := decoder.Decode(&obj); err != nil {
		return false
	}

	// Additional validation using k8s.io/apimachinery
	_, err := meta.Accessor(obj)
	return err == nil
}

func isHelmChart(filePath string) bool {
	if _, err := os.Stat(filepath.Join(filepath.Dir(filePath), "Chart.yaml")); err == nil {
		_, err := loader.Load(filepath.Dir(filePath))
		return err == nil
	}
	return false
}

func isKustomize(filePath string) bool {
	if base := filepath.Base(filePath); base == "kustomization.yaml" || base == "kustomization.yml" {
		fSys := filesys.MakeFsOnDisk()
		k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
		_, err := k.Run(fSys, filepath.Dir(filePath))
		return err == nil
	}
	return false
}
