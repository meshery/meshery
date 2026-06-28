package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/modeloci"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	meshkitUtils "github.com/meshery/meshkit/utils"
)

func publishModelOCIArtifacts() error {
	if err := os.MkdirAll(ociOutputPath, 0o755); err != nil {
		return err
	}

	totalPublished := 0
	failures := []string{}

	for _, model := range models {
		modelName := meshkitUtils.FormatName(model.Model)
		modelDir, version, err := resolveModelOCISource(ociSourcePath, modelName, model)
		if err != nil {
			if ociFailFast {
				return err
			}
			failures = append(failures, err.Error())
			continue
		}

		utils.Log.Infof("Building OCI artifact for model %s version %s", modelName, version)
		artifactPath, err := modeloci.BuildModelOCIArtifactFromFolder(filepath.Join(modelDir, version), ociOutputPath, modelName, version)
		if err != nil {
			err = fmt.Errorf("failed to build OCI artifact for model %s version %s: %w", modelName, version, err)
			if ociFailFast {
				return err
			}
			failures = append(failures, err.Error())
			continue
		}

		utils.Log.Infof("Published OCI artifact for model %s: %s", modelName, artifactPath)
		totalPublished++
	}

	utils.Log.Info("Total model OCI artifacts published: ", totalPublished)
	if len(failures) > 0 {
		return fmt.Errorf("encountered %d OCI publishing failure(s): %s", len(failures), strings.Join(failures, "; "))
	}

	return nil
}

func resolveModelOCIVersion(sourcePath, modelName string, model meshkitRegistryUtils.ModelCSV) (string, error) {
	_, version, err := resolveModelOCISource(sourcePath, modelName, model)
	return version, err
}

func resolveModelOCISource(sourcePath, modelName string, model meshkitRegistryUtils.ModelCSV) (string, string, error) {
	modelDir, err := resolveModelOCISourceDir(sourcePath, modelName, model.Model)
	if err != nil {
		return "", "", err
	}

	versions, err := readModelVersionDirectories(modelDir)
	if err != nil {
		return "", "", err
	}
	if len(versions) == 0 {
		return "", "", fmt.Errorf("no model version directories found for model %s in %s", modelName, modelDir)
	}
	if len(versions) == 1 {
		return modelDir, versions[0], nil
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return compareVersionStrings(versions[i], versions[j]) > 0
	})
	return modelDir, versions[0], nil
}

func resolveModelOCISourceDir(sourcePath, modelName, rawModelName string) (string, error) {
	candidates := []string{modelName, rawModelName}
	seen := map[string]bool{}
	for _, candidate := range candidates {
		if candidate == "" || seen[candidate] {
			continue
		}
		seen[candidate] = true
		candidatePath := filepath.Join(sourcePath, candidate)
		info, err := os.Stat(candidatePath)
		if err == nil && info.IsDir() {
			return candidatePath, nil
		}
	}

	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if entry.IsDir() && meshkitUtils.FormatName(entry.Name()) == modelName {
			return filepath.Join(sourcePath, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("model source directory not found for model %s in %s", modelName, sourcePath)
}

func readModelVersionDirectories(modelDir string) ([]string, error) {
	entries, err := os.ReadDir(modelDir)
	if err != nil {
		return nil, err
	}

	versions := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			containsDigit := false
			for _, r := range entry.Name() {
				if unicode.IsDigit(r) {
					containsDigit = true
					break
				}
			}
			if containsDigit {
				versions = append(versions, entry.Name())
			}
		}
	}

	return versions, nil
}

func compareVersionStrings(left, right string) int {
	leftParts := versionNumberParts(left)
	rightParts := versionNumberParts(right)
	maxLen := len(leftParts)
	if len(rightParts) > maxLen {
		maxLen = len(rightParts)
	}

	for i := 0; i < maxLen; i++ {
		leftPart := 0
		rightPart := 0
		if i < len(leftParts) {
			leftPart = leftParts[i]
		}
		if i < len(rightParts) {
			rightPart = rightParts[i]
		}
		if leftPart > rightPart {
			return 1
		}
		if leftPart < rightPart {
			return -1
		}
	}

	leftHasDash := strings.Contains(left, "-")
	rightHasDash := strings.Contains(right, "-")
	if leftHasDash && !rightHasDash {
		return -1
	}
	if !leftHasDash && rightHasDash {
		return 1
	}
	return strings.Compare(left, right)
}

func versionNumberParts(version string) []int {
	version = strings.TrimPrefix(version, "v")
	if idx := strings.IndexAny(version, "-+"); idx != -1 {
		version = version[:idx]
	}
	fields := strings.FieldsFunc(version, func(r rune) bool {
		return !unicode.IsDigit(r)
	})

	parts := make([]int, 0, len(fields))
	for _, field := range fields {
		if field == "" {
			continue
		}
		part, err := strconv.Atoi(field)
		if err != nil {
			continue
		}
		parts = append(parts, part)
	}

	return parts
}
