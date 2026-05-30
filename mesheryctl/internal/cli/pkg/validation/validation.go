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

// Package validation provides shared utilities for mesheryctl validate commands.
package validation

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/encoding"
)

// maxResponseBytes limits URL-fetched file size to 25 MB.
const maxResponseBytes = 25 * 1024 * 1024

// Result represents the outcome of validating a single entity.
type Result struct {
	FilePath   string   `json:"filePath" yaml:"filePath"`
	EntityType string   `json:"entityType" yaml:"entityType"`
	EntityName string   `json:"entityName" yaml:"entityName"`
	IsValid    bool     `json:"isValid" yaml:"isValid"`
	Errors     []string `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// Summary holds the overall validation results.
type Summary struct {
	TotalEntities int      `json:"totalEntities" yaml:"totalEntities"`
	TotalValid    int      `json:"totalValid" yaml:"totalValid"`
	TotalInvalid  int      `json:"totalInvalid" yaml:"totalInvalid"`
	Results       []Result `json:"results" yaml:"results"`
}

// SchemaHeader is used to detect the entity type from the schemaVersion field.
type SchemaHeader struct {
	SchemaVersion string `json:"schemaVersion" yaml:"schemaVersion"`
}

// DetectEntityType determines the entity type from the schemaVersion field.
func DetectEntityType(data []byte) string {
	var header SchemaHeader
	if err := encoding.Unmarshal(data, &header); err != nil {
		return ""
	}
	switch {
	case strings.HasPrefix(header.SchemaVersion, "models.meshery.io"):
		return "model"
	case strings.HasPrefix(header.SchemaVersion, "components.meshery.io"):
		return "component"
	case strings.HasPrefix(header.SchemaVersion, "relationships.meshery.io"):
		return "relationship"
	default:
		return ""
	}
}

// FileFilter is a predicate that determines whether a file (by its contents)
// should be included during directory traversal.
type FileFilter func(data []byte) bool

// AnyEntityFilter includes files with any recognized schemaVersion.
func AnyEntityFilter(data []byte) bool {
	return DetectEntityType(data) != ""
}

// RelationshipOnlyFilter includes only relationship definition files.
func RelationshipOnlyFilter(data []byte) bool {
	return DetectEntityType(data) == "relationship"
}

// ValidateHTTPScheme checks that the given URL uses http or https.
func ValidateHTTPScheme(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %w", rawURL, err)
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported URL scheme %q in %s: only http and https are supported", u.Scheme, rawURL)
	}
	return nil
}

// FetchFileFromURL retrieves file content from an HTTP/HTTPS URL.
// The response is bounded by maxResponseBytes.
func FetchFileFromURL(fetchURL string) ([]byte, error) {
	if err := ValidateHTTPScheme(fetchURL); err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(fetchURL) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", fetchURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: HTTP %d", fetchURL, resp.StatusCode)
	}

	limitedReader := io.LimitReader(resp.Body, maxResponseBytes+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s: %w", fetchURL, err)
	}
	if len(data) > maxResponseBytes {
		return nil, fmt.Errorf("response from %s exceeds maximum size of %d bytes", fetchURL, maxResponseBytes)
	}
	return data, nil
}

// CollectFiles gathers file contents from the given paths (files, directories,
// or URLs). The filter controls which files are included during directory
// traversal; explicitly specified files and URLs are always collected.
func CollectFiles(paths []string, filter FileFilter) (map[string][]byte, error) {
	result := make(map[string][]byte)

	for _, path := range paths {
		if utils.IsValidUrl(path) {
			data, err := FetchFileFromURL(path)
			if err != nil {
				return nil, err
			}
			result[path] = data
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("cannot access %s: %w", path, err)
		}

		if info.IsDir() {
			err := filepath.Walk(path, func(filePath string, fi os.FileInfo, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if fi.IsDir() {
					return nil
				}
				ext := strings.ToLower(filepath.Ext(filePath))
				if ext != ".json" && ext != ".yaml" && ext != ".yml" {
					return nil
				}
				data, readErr := os.ReadFile(filePath)
				if readErr != nil {
					return nil // skip unreadable files during directory walk
				}
				if filter == nil || filter(data) {
					result[filePath] = data
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", path, err)
			}
			// Always include explicitly specified files
			result[path] = data
		}
	}

	return result, nil
}

// SortedPaths returns the keys of a map sorted lexicographically.
func SortedPaths(files map[string][]byte) []string {
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// DisplayResults renders validation results in the specified output format.
// entityLabel is used in the summary line (e.g. "entities" or "relationships").
func DisplayResults(summary Summary, outputFormat string, entityLabel string) error {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "json":
		formatter := display.NewJSONOutputFormatter(summary)
		if err := formatter.Display(); err != nil {
			return err
		}
	case "yaml":
		formatter := display.NewYAMLOutputFormatter(summary)
		if err := formatter.Display(); err != nil {
			return err
		}
	default:
		header := []string{"FILE", "TYPE", "NAME", "STATUS"}
		rows := make([][]string, 0, len(summary.Results))
		for _, result := range summary.Results {
			status := "PASS"
			if !result.IsValid {
				status = "FAIL"
			}
			rows = append(rows, []string{
				result.FilePath,
				result.EntityType,
				result.EntityName,
				status,
			})
		}
		utils.PrintToTable(header, rows, nil)

		// Show errors for failed entities
		for _, result := range summary.Results {
			if !result.IsValid {
				fmt.Printf("\n  %s %s:\n", utils.BoldString("Errors in"), result.FilePath)
				for _, e := range result.Errors {
					fmt.Printf("    - %s\n", e)
				}
			}
		}

		fmt.Printf("\n%s: %d %s validated, %d passed, %d failed\n",
			utils.BoldString("SUMMARY"),
			summary.TotalEntities,
			entityLabel,
			summary.TotalValid,
			summary.TotalInvalid,
		)
	}

	return nil
}
