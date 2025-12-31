// # Copyright Meshery Authors
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
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	mutils "github.com/meshery/meshkit/utils"
	"gopkg.in/yaml.v3"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
)

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
	componentSpredsheetGID         int64
	relationshipSpredsheetGID      int64
	outputLocation                 string
	pathToRegistrantConnDefinition string
	pathToRegistrantCredDefinition string
	GoogleSpreadSheetURL           = "https://docs.google.com/spreadsheets/d/"
	srv                            *sheets.Service

	// current working directory location
	cwd string

	registryLocation    string
	totalAggregateModel int
	defVersion          = "v1.0.0"
	checkFlag           bool   // Triggers the check mode
	checkFormat         string // Holds "json", "yaml", or "table"
	checkCount          bool   // Triggers "count only" mode
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Models",
	Long: `Prerequisite: Excecute this command from the root of a meshery/meshery repo fork.\n\nGiven a Google Sheet with a list of model names and source locations, generate models and components any Registrant (e.g. GitHub, Artifact Hub) repositories.\n\nGenerated Model files are written to local filesystem under "/server/models/<model-name>".
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/registry/generate`,
	Example: `
// Generate Meshery Models from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred "$CRED"

// Directly generate models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition
mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]

// Generate a specific Model from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred --model "[model-name]"

// Generate Meshery Models and Component from csv files in a local directory.
mesheryctl registry generate --directory [DIRECTORY_PATH]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Prerequisite check is needed - https://github.com/meshery/meshery/issues/10369
		// TODO: Include a prerequisite check to confirm that this command IS being the executED from within a fork of the Meshery repo, and is being executed at the root of that fork.
		if checkFlag {
			return nil
		}
		const errorMsg = "[ Spreadsheet ID | Registrant Connection Definition Path | Local Directory ] isn't specified\n\nUsage: \nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED --model \"[model-name]\"\nRun 'mesheryctl registry generate --help' to see detailed help message"

		spreadsheetIdFlag, _ := cmd.Flags().GetString("spreadsheet-id")
		registrantDefFlag, _ := cmd.Flags().GetString("registrant-def")
		directory, _ := cmd.Flags().GetString("directory")

		if spreadsheetIdFlag == "" && registrantDefFlag == "" && directory == "" {
			return errors.New(utils.RegistryError(errorMsg, "generate"))
		}

		spreadsheetCredFlag, _ := cmd.Flags().GetString("spreadsheet-cred")
		registrantCredFlag, _ := cmd.Flags().GetString("registrant-cred")

		if spreadsheetIdFlag != "" && spreadsheetCredFlag == "" {
			return errors.New(utils.RegistryError("Spreadsheet Credentials is required\n\nUsage: \nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED --model \"[model-name]\"\nRun 'mesheryctl registry generate --help'", "generate"))
		}

		if registrantDefFlag != "" && registrantCredFlag == "" {
			return errors.New(utils.RegistryError("Registrant Credentials is required\n\nUsage: mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]\nRun 'mesheryctl registry generate --help'", "generate"))
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if checkFlag {
			// 1. Validation
			if len(args) == 0 {
				return errors.New("please provide an image to scan.\nUsage: mesheryctl registry generate --check <image> [--format json|yaml]")
			}
			imageRef := args[0]

			// 2. Validate Format
			if checkFormat != "" && checkFormat != "json" && checkFormat != "yaml" && checkFormat != "table" {
				return fmt.Errorf("invalid format '%s'. Allowed: json, yaml, table", checkFormat)
			}

			utils.Log.Info(fmt.Sprintf("Scanning OCI image: %s", imageRef))

			crds, err := FindCRDsInImage(imageRef)
			if err != nil {
				return fmt.Errorf("failed to scan image: %w", err)
			}

			// 4. Handle "Count Only" Mode
			if checkCount {
				utils.Log.Info(fmt.Sprintf("Found %d CRD(s) in image %s", len(crds), imageRef))
				return nil
			}

			if len(crds) == 0 {
				utils.Log.Info("No CRDs found in the image.")
				return nil
			}

			// 5. Build Result Struct
			// (Ensure FindResult struct is defined in this file!)
			result := FindResult{
				ImageRef: imageRef,
				CRDs:     crds,
				Count:    len(crds),
			}

			// 6. Output Formatting (JSON / YAML / Table)
			switch checkFormat {
			case "json":
				output, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(output))
			case "yaml":
				output, err := yaml.Marshal(result)
				if err != nil {
					return fmt.Errorf("failed to marshal YAML: %w", err)
				}
				fmt.Println(string(output))
			default: // "table"
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

			// 7. STOP. Do not run the generation logic below.
			return nil
		}
		var wg sync.WaitGroup
		cwd, _ = os.Getwd()
		registryLocation = filepath.Join(cwd, outputLocation)

		if pathToRegistrantConnDefinition != "" {
			utils.Log.Info("Model generation from Registrant definitions not yet supported.")
			return nil
		}
		var err error

		if csvDirectory == "" {
			srv, err = mutils.NewSheetSRV(spreadsheeetCred)
			if err != nil {
				return errors.New(utils.RegistryError("Invalid JWT Token: Ensure the provided token is a base64-encoded, valid Google Spreadsheets API token.", "generate"))
			}

			resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
			if err != nil || resp.HTTPStatusCode != 200 {
				utils.LogError.Error(ErrUpdateRegistry(err, outputLocation))
				return nil
			}

			sheetGID = GetSheetIDFromTitle(resp, "Models")
			componentSpredsheetGID = GetSheetIDFromTitle(resp, "Components")
			// Collect list of corresponding relationship by name from spreadsheet
			relationshipSpredsheetGID = GetSheetIDFromTitle(resp, "Relationships")
		} else {
			modelCSVFilePath, componentCSVFilePath, relationshipCSVFilePath, err = meshkitRegistryUtils.GetCsv(csvDirectory)
			if err != nil {
				return fmt.Errorf("error reading the directory: %v", err)
			}
			if modelCSVFilePath == "" || componentCSVFilePath == "" || relationshipCSVFilePath == "" {
				return fmt.Errorf("ModelCSV, ComponentCSV and RelationshipCSV files must be present in the directory")
			}
		}

		err = meshkitRegistryUtils.InvokeGenerationFromSheet(&wg, registryLocation, sheetGID, componentSpredsheetGID, spreadsheeetID, modelName, modelCSVFilePath, componentCSVFilePath, spreadsheeetCred, relationshipCSVFilePath, relationshipSpredsheetGID, srv)
		if err != nil {
			// meshkit
			utils.LogError.Error(err)
			return nil
		}
		_ = logFile.Close()
		_ = errorLogFile.Close()

		utils.Log.UpdateLogOutput(os.Stdout)
		utils.LogError.UpdateLogOutput(os.Stdout)
		return err
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
		crds, err = processLayer(layer, crds)
		if err != nil {
			return nil, err
		}
	}

	return crds, nil
}

func processLayer(layer v1.Layer, existingCRDs []CRD) ([]CRD, error) {
	rc, err := layer.Compressed()
	if err != nil {
		return nil, fmt.Errorf("getting layer compressed reader: %w", err)
	}
	defer rc.Close()

	// 1. Use bufio to peek at the bytes without reading everything into memory
	br := bufio.NewReader(rc)
	peek, _ := br.Peek(2)

	var tr *tar.Reader

	// 2. Check for GZIP Magic Bytes (0x1f, 0x8b)
	if len(peek) >= 2 && peek[0] == 0x1f && peek[1] == 0x8b {
		gr, err := gzip.NewReader(br)
		if err != nil {
			return nil, fmt.Errorf("creating gzip reader: %w", err)
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else {
		// Not gzipped, treat as standard tar
		tr = tar.NewReader(br)
	}

	return scanTarForCRDs(tr, existingCRDs)
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

			fileContent := buf.String()

			if !strings.Contains(fileContent, "kind: CustomResourceDefinition") {
				continue
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
	generateCmd.PersistentFlags().StringVar(&spreadsheeetID, "spreadsheet-id", "", "spreadsheet ID for the integration spreadsheet")
	generateCmd.PersistentFlags().StringVar(&spreadsheeetCred, "spreadsheet-cred", "", "base64 encoded credential to download the spreadsheet")

	generateCmd.MarkFlagsRequiredTogether("spreadsheet-id", "spreadsheet-cred")

	generateCmd.PersistentFlags().StringVar(&pathToRegistrantConnDefinition, "registrant-def", "", "path pointing to the registrant connection definition")
	generateCmd.PersistentFlags().StringVar(&pathToRegistrantCredDefinition, "registrant-cred", "", "path pointing to the registrant credential definition")

	generateCmd.MarkFlagsRequiredTogether("registrant-def", "registrant-cred")

	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-id", "registrant-def")
	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-cred", "registrant-cred")
	generateCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "", "specific model name to be generated")
	generateCmd.PersistentFlags().StringVarP(&outputLocation, "output", "o", "../server/meshmodel", "location to output generated models, defaults to ../server/meshmodels")

	generateCmd.PersistentFlags().StringVarP(&csvDirectory, "directory", "d", "", "Directory containing the Model and Component CSV files")
	generateCmd.Flags().BoolVarP(&checkFlag, "check", "c", false, "Scan image for CRDs without generating")
	generateCmd.Flags().StringVar(&checkFormat, "format", "table", "Output format (json|yaml|table)")
	generateCmd.Flags().BoolVar(&checkCount, "count", false, "Only output the number of CRDs found")
}
