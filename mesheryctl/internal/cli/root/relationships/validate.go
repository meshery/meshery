package relationships

import (
	"fmt"
	"os"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitschema "github.com/meshery/meshkit/schema"
	schemav1alpha3 "github.com/meshery/schemas/models/v1alpha3"
	"github.com/spf13/cobra"
)

type relationshipValidateFlags struct {
	File string `json:"file" validate:"required"`
}

var relationshipValidateFlagsProvided relationshipValidateFlags

const relationshipValidateUsage = "Usage: mesheryctl relationship validate --file <path>\nRun 'mesheryctl relationship validate --help' to see detailed help message"

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a relationship definition file",
	Long:  "Validate a relationship definition file against the Meshery relationship schema.",
	Example: `
// Validate a relationship definition file
mesheryctl relationship validate --file ./relationship.yaml

// Validate a JSON relationship definition file
mesheryctl relationship validate --file ./relationship.json
`,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &relationshipValidateFlagsProvided)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", relationshipValidateUsage))
		}

		return nil
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		relationshipData, err := os.ReadFile(relationshipValidateFlagsProvided.File)
		if err != nil {
			return utils.ErrFileRead(err)
		}

		err = meshkitschema.ValidateWithRef(meshkitschema.Ref{
			SchemaVersion: schemav1alpha3.RelationshipSchemaVersion,
			Type:          meshkitschema.TypeRelationship,
		}, relationshipData)
		if err != nil {
			if details, ok := meshkitschema.ValidationDetailsFromError(err); ok {
				return ErrRelationshipValidationFailed(details)
			}
			return err
		}

		utils.Log.Info("Relationship definition is valid.")
		return nil
	},
}

func init() {
	validateCmd.Flags().StringVarP(&relationshipValidateFlagsProvided.File, "file", "f", "", "(required) path to the relationship definition file")
}
