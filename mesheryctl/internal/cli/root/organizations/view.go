package organizations

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/organization"
	"github.com/spf13/cobra"
)

type orgViewFlags struct {
	OutputFormat string `json:"output-format" validate:"required,oneof=json yaml"`
	Save         bool   `json:"save" validate:"boolean"`
}

var orgViewFlagsProvided orgViewFlags

func formatOrgLabel(rows []organization.Organization) []string {
	labels := []string{}
	for _, o := range rows {
		labels = append(labels, fmt.Sprintf("%s (ID: %s)", o.Name, o.ID.String()))
	}
	return labels
}

var viewOrgCmd = &cobra.Command{
	Use:   "view [organization-name|organization-id]",
	Short: "View an organization",
	Long: `View an organization by its ID or name.
Find more information at: https://docs.meshery.io/reference/mesheryctl/organizations/view`,
	Example: `// View details of a specific organization by ID
mesheryctl organization view [organization-id]

// View details of a specific organization by name
mesheryctl organization view [organization-name]

// View details of a specific organization in JSON format
mesheryctl organization view [organization-id] --output-format json

// View details of a specific organization and save it to a file
mesheryctl organization view [organization-id] --output-format json --save`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &orgViewFlagsProvided)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl organization view [organization-name|organization-id]\nRun 'mesheryctl organization view --help' to see detailed help message"
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one organization name or ID\n\n%v", errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		orgNameOrID := args[0]
		selectedOrg := new(organization.Organization)
		var urlPath string
		var displayData display.DisplayDataAsync

		if utils.IsUUID(orgNameOrID) {
			urlPath = fmt.Sprintf("%s/%s", organizationsApiPath, url.PathEscape(orgNameOrID))
			displayData = display.DisplayDataAsync{UrlPath: urlPath}
			err := display.PromptAsyncPagination(
				displayData,
				formatOrgLabel,
				func(data *organization.Organization) ([]organization.Organization, int64) {
					return []organization.Organization{*data}, 1
				},
				selectedOrg,
			)
			if err != nil {
				return err
			}
		} else {
			viewUrlValue := url.Values{}
			viewUrlValue.Add("search", orgNameOrID)
			urlPath = fmt.Sprintf("%s?%s", organizationsApiPath, viewUrlValue.Encode())
			displayData = display.DisplayDataAsync{
				UrlPath:    urlPath,
				SearchTerm: orgNameOrID,
			}
			err := display.PromptAsyncPagination(
				displayData,
				formatOrgLabel,
				func(data *models.OrganizationsPage) ([]organization.Organization, int64) {
					orgs := make([]organization.Organization, len(data.Organizations))
					for i, o := range data.Organizations {
						orgs[i] = *o
					}
					return orgs, int64(data.TotalCount)
				},
				selectedOrg,
			)
			if err != nil {
				return err
			}
		}

		outputFormatterFactory := display.OutputFormatterFactory[organization.Organization]{}
		outputFormatter, err := outputFormatterFactory.New(orgViewFlagsProvided.OutputFormat, *selectedOrg)
		if err != nil {
			return err
		}

		outputFormatter = outputFormatter.WithOutput(cmd.OutOrStdout())
		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		if orgViewFlagsProvided.Save {
			orgString := strings.ReplaceAll(selectedOrg.Name, " ", "_")
			if orgString == "" {
				orgString = selectedOrg.ID.String()
			}
			fileName := filepath.Join(utils.MesheryFolder, fmt.Sprintf("organization_%s.%s", orgString, orgViewFlagsProvided.OutputFormat))
			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[organization.Organization]{}
			outputFormatterSaver, err := outputFormatterSaverFactory.New(orgViewFlagsProvided.OutputFormat, outputFormatter)
			if err != nil {
				return err
			}
			outputFormatterSaver = outputFormatterSaver.WithFilePath(fileName)
			err = outputFormatterSaver.Save()
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	viewOrgCmd.Flags().StringVarP(&orgViewFlagsProvided.OutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewOrgCmd.Flags().BoolVarP(&orgViewFlagsProvided.Save, "save", "s", false, "(optional) save output as a JSON/YAML file")
}
