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

package registry

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CRD represents a detected CustomResourceDefinition with path and basic metadata.
type CRD struct {
	Path     string `json:"path" yaml:"path"`
	Group    string `json:"group" yaml:"group"`
	Version  string `json:"version" yaml:"version"`
	Kind     string `json:"kind" yaml:"kind"`
	Plural   string `json:"plural" yaml:"plural"`
	Singular string `json:"singular" yaml:"singular"`
}

// FindResult represents the complete output of a find operation
type FindResult struct {
	ImageRef string `json:"imageRef" yaml:"imageRef"`
	CRDs     []CRD  `json:"crds" yaml:"crds"`
	Count    int    `json:"count" yaml:"count"`
}

var (
	findOutputFormat string
	findCountOnly    bool
)

// findCmd represents the registry find command
var findCmd = &cobra.Command{
	Use:   "find <oci-image-reference>",
	Short: "Scan an OCI image for CRDs and discover potential components",
	Long: `Scan an OCI container image for CustomResourceDefinitions (CRDs) that can be registered as Meshery components.

This command pulls the specified OCI image and scans its layers for YAML files containing CRDs.
The discovered CRDs can be used to generate new Meshery models and components.

The output can be formatted as JSON or YAML for piping to other commands like 'mesheryctl registry generate'.
Documentation for the registry command can be found at https://docs.meshery.io/reference/mesheryctl/registry`,
	Example: `
// Scan an OCI image for CRDs
mesheryctl registry find docker.io/crossplane/crossplane:v1.14.0

// Output discovered CRDs as JSON (useful for piping to other commands)
mesheryctl registry find crossplane/crossplane:v1.14.0 -o json

// Output discovered CRDs as YAML
mesheryctl registry find crossplane/crossplane:v1.14.0 -o yaml

// Only display the count of discovered CRDs
mesheryctl registry find crossplane/crossplane:v1.14.0 --count
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl registry find <oci-image-reference>\nRun 'mesheryctl registry find --help' to see detailed help message"
		if len(args) == 0 {
			return ErrFindImageRefRequired()
		}
		if len(args) > 1 {
			return fmt.Errorf("accepts 1 arg, received %d\n\n%s", len(args), errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		imageRef := args[0]

		// Validate output format
		if findOutputFormat != "" && findOutputFormat != "json" && findOutputFormat != "yaml" {
			return ErrFindInvalidOutputFormat(findOutputFormat)
		}

		utils.Log.Info(fmt.Sprintf("Scanning OCI image: %s", imageRef))

		crds, err := FindCRDsInImage(imageRef)
		if err != nil {
			return ErrFindScanImage(err, imageRef)
		}

		result := FindResult{
			ImageRef: imageRef,
			CRDs:     crds,
			Count:    len(crds),
		}

		// Handle count-only output
		if findCountOnly {
			utils.Log.Info(fmt.Sprintf("Found %d CRD(s) in image %s", result.Count, imageRef))
			return nil
		}

		if len(crds) == 0 {
			utils.Log.Info("No CRDs found in the image.")
			return nil
		}

		// Handle formatted output (JSON/YAML)
		switch findOutputFormat {
		case "json":
			output, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return utils.ErrMarshalIndent(err)
			}
			fmt.Println(string(output))
		case "yaml":
			output, err := yaml.Marshal(result)
			if err != nil {
				return utils.ErrMarshal(err)
			}
			fmt.Println(string(output))
		default:
			// Default human-readable table output
			utils.Log.Info(fmt.Sprintf("Discovered %d CRD(s):", len(crds)))
			fmt.Println()
			header := []string{"KIND", "GROUP", "VERSION", "PLURAL", "PATH"}
			rows := [][]string{}
			for _, crd := range crds {
				rows = append(rows, []string{
					crd.Kind,
					crd.Group,
					crd.Version,
					crd.Plural,
					crd.Path,
				})
			}
			utils.PrintToTable(header, rows, nil)
		}

		return nil
	},
}

// FindCRDsInImage scans the OCI image for YAML files and parses them to detect CRDs.
func FindCRDsInImage(imageRef string) ([]CRD, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, fmt.Errorf("parsing reference: %w", err)
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, fmt.Errorf("fetching image: %w", err)
	}

	layers, err := img.Layers()
	if err != nil {
		return nil, fmt.Errorf("getting layers: %w", err)
	}

	var crds []CRD
	for _, layer := range layers {
		rc, err := layer.Compressed()
		if err != nil {
			return nil, fmt.Errorf("getting layer compressed reader: %w", err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("reading layer data: %w", err)
		}

		reader := bytes.NewReader(data)
		gr, err := gzip.NewReader(reader)
		var tr *tar.Reader
		if err == gzip.ErrHeader {
			reader = bytes.NewReader(data)
			tr = tar.NewReader(reader)
		} else if err != nil {
			return nil, fmt.Errorf("creating gzip reader: %w", err)
		} else {
			defer gr.Close()
			tr = tar.NewReader(gr)
		}

		crds, err = scanTarForCRDs(tr, crds)
		if err != nil {
			return nil, err
		}
	}

	return crds, nil
}

// scanTarForCRDs scans the tar for .yaml/.yml files, reads content, and checks if it's a CRD.
func scanTarForCRDs(tr *tar.Reader, crds []CRD) ([]CRD, error) {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar header: %w", err)
		}

		// Check for YAML files
		isYAML := strings.HasSuffix(hdr.Name, ".yaml") || strings.HasSuffix(hdr.Name, ".yml")
		if hdr.Typeflag == tar.TypeReg && isYAML {
			// Read file content
			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, tr); err != nil {
				return nil, fmt.Errorf("reading file %s: %w", hdr.Name, err)
			}

			// Parse YAML - handle multi-document YAML files
			decoder := yaml.NewDecoder(bytes.NewReader(buf.Bytes()))
			for {
				var doc struct {
					Kind       string `yaml:"kind"`
					APIVersion string `yaml:"apiVersion"`
					Metadata   struct {
						Name string `yaml:"name"`
					} `yaml:"metadata"`
					Spec struct {
						Group    string `yaml:"group"`
						Versions []struct {
							Name string `yaml:"name"`
						} `yaml:"versions"`
						Names struct {
							Plural   string `yaml:"plural"`
							Singular string `yaml:"singular"`
							Kind     string `yaml:"kind"`
						} `yaml:"names"`
					} `yaml:"spec"`
				}

				if err := decoder.Decode(&doc); err != nil {
					if err == io.EOF {
						break
					}
					continue // Skip invalid YAML documents
				}

				if doc.Kind == "CustomResourceDefinition" {
					version := ""
					if len(doc.Spec.Versions) > 0 {
						version = doc.Spec.Versions[0].Name
					}
					crds = append(crds, CRD{
						Path:     filepath.Clean("/" + hdr.Name),
						Group:    doc.Spec.Group,
						Version:  version,
						Kind:     doc.Spec.Names.Kind,
						Plural:   doc.Spec.Names.Plural,
						Singular: doc.Spec.Names.Singular,
					})
				}
			}
		}
	}
	return crds, nil
}

func init() {
	findCmd.Flags().StringVarP(&findOutputFormat, "output-format", "o", "", "(optional) format to display in [json|yaml]")
	findCmd.Flags().BoolVarP(&findCountOnly, "count", "c", false, "(optional) only display the count of CRDs found")
}
