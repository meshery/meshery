package modeloci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	meshkitOci "github.com/meshery/meshkit/models/oci"
)

func ParseModelInput(input string) (name, version string) {
	parts := strings.Split(input, "/")
	name = parts[0]
	if len(parts) > 1 {
		version = parts[1]
	}
	return
}

func CompileFolderName(path string, name string, version string) string {
	dirParts := make([]string, 0, 3)
	if path != "" {
		dirParts = append(dirParts, path)
	}
	dirParts = append(dirParts, name)
	if version != "" {
		dirParts = append(dirParts, version)
	}
	return filepath.Join(dirParts...)
}

func CompileImageName(name string, version string, extension string) string {
	if version != "" {
		return fmt.Sprintf(
			"%s-%s.%s",
			name,
			strings.ReplaceAll(version, ".", "-"),
			extension,
		)
	}
	return fmt.Sprintf(
		"%s.%s",
		name,
		extension,
	)
}

func BuildModelOCIArtifact(sourcePath, outputPath, name, version string) (string, error) {
	return BuildModelOCIArtifactFromFolder(CompileFolderName(sourcePath, name, version), outputPath, name, version)
}

func BuildModelOCIArtifactFromFolder(folder, outputPath, name, version string) (string, error) {
	if _, err := os.Stat(folder); err != nil {
		return "", err
	}

	if outputPath != "" {
		if err := os.MkdirAll(outputPath, 0o755); err != nil {
			return "", err
		}
	}

	img, err := meshkitOci.BuildImage(folder)
	if err != nil {
		return "", err
	}

	artifactPath := filepath.Join(outputPath, CompileImageName(name, version, "tar"))
	if err := meshkitOci.SaveOCIArtifact(img, artifactPath, name); err != nil {
		return "", err
	}

	return artifactPath, nil
}
