package registry // Adjust if the package is different, e.g., "cmd"

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3" // For YAML parsing; add to go.mod if needed: go get gopkg.in/yaml.v3
)

// CRD represents a detected CustomResourceDefinition with path and basic metadata.
type CRD struct {
	Path     string
	Group    string
	Version  string
	Kind     string
	Plural   string
	Singular string
}

// findCmd represents the registry find command
var findCmd = &cobra.Command{
	Use:   "find [image-reference]",
	Short: "Scan an OCI image for bundled CRDs and discover potential components",
	Long: `Scans the specified OCI image (e.g., docker.io/crossplane/crossplane:v1.14.0) for YAML files containing CustomResourceDefinitions (CRDs).
This can help discover new components to register in Meshery's registry by pointing to a container image in a registry like Docker Hub.
Example: mesheryctl registry find crossplane/crossplane:v1.14.0`,
	Args: cobra.ExactArgs(1), // Requires exactly one argument: the image ref
	RunE: func(cmd *cobra.Command, args []string) error {
		imageRef := args[0]
		crds, err := FindCRDsInImage(imageRef)
		if err != nil {
			return fmt.Errorf("failed to scan image: %w", err)
		}

		if len(crds) == 0 {
			fmt.Println("No CRDs found in the image.")
			return nil
		}

		fmt.Println("Discovered CRDs:")
		for _, crd := range crds {
			fmt.Printf("- Path: %s\n  Group: %s, Version: %s, Kind: %s, Plural: %s, Singular: %s\n", crd.Path, crd.Group, crd.Version, crd.Kind, crd.Plural, crd.Singular)
		}

		// TODO: Integrate with Meshery registry - e.g., call internal APIs to register these as models/components.
		// For example, if Meshery has a registry client: registryClient.RegisterComponents(crds)

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
		defer rc.Close()

		gr, err := gzip.NewReader(rc)
		if err != nil && err != gzip.ErrHeader {
			return nil, fmt.Errorf("creating gzip reader: %w", err)
		}
		var tr *tar.Reader
		if err == gzip.ErrHeader {
			// Not gzipped, reset and use raw
			if _, err := rc.Seek(0, io.SeekStart); err != nil {
				return nil, err
			}
			tr = tar.NewReader(rc)
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

// scanTarForCRDs scans the tar for .yaml files, reads content, and checks if it's a CRD.
func scanTarForCRDs(tr *tar.Reader, crds []CRD) ([]CRD, error) {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar header: %w", err)
		}

		if hdr.Typeflag == tar.TypeReg && strings.HasSuffix(hdr.Name, ".yaml") {
			// Read file content
			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, tr); err != nil {
				return nil, fmt.Errorf("reading file %s: %w", hdr.Name, err)
			}

			// Parse YAML
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
			if err := yaml.Unmarshal(buf.Bytes(), &doc); err != nil {
				continue // Skip invalid YAML
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
	return crds, nil
}

func init() {
	// Add flags if needed, e.g., for filters or registry URL
	findCmd.Flags().StringP("filter", "f", "**/*.yaml", "Glob pattern for files to scan")
}
