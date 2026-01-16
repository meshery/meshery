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

package connections

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	outputFormatFlag string
	saveFlag         bool
)

var viewConnectionCmd = &cobra.Command{
	Use:   "view",
	Short: "View a connection",
	Long:  `View a connection by its ID or name`,
	Example: `
// View details of a specific connection
mesheryctl connection view [connection-name]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl connection view [connection-name]\nRun 'mesheryctl connection view --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New("connection name or ID isn't specified\n\n" + errMsg))
		} else if len(args) > 1 {
			return utils.ErrInvalidArgument(errors.New("too many arguments\n\n" + errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		connectionNameOrID := args[0]

		var selectedConnection *connection.Connection

		// Check if the argument is a valid UUID
		if _, err := uuid.FromString(connectionNameOrID); err == nil {
			// Fetch connection directly by ID
			url := fmt.Sprintf("%s/%s", connectionApiPath, connectionNameOrID)
			conn, err := api.Fetch[connection.Connection](url)
			if err != nil {
				utils.Log.Error(err)
				return err
			}
			selectedConnection = conn
		} else {
			// Search by name
			url := fmt.Sprintf("%s?search=%s&pagesize=all", connectionApiPath, connectionNameOrID)
			connectionsResponse, err := api.Fetch[connection.ConnectionPage](url)
			if err != nil {
				utils.Log.Error(err)
				return err
			}

			switch connectionsResponse.TotalCount {
			case 0:
				fmt.Println("No connection(s) found for the given name: ", connectionNameOrID)
				return nil
			case 1:
				selectedConnection = connectionsResponse.Connections[0]
			default:
				selectedConnection = selectConnectionPrompt(connectionsResponse.Connections)
			}
		}

		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outputFormatFlag = strings.ToLower(outputFormatFlag)

		if outputFormatFlag != "json" && outputFormatFlag != "yaml" {
			return errors.New("output-format choice is invalid or not provided, use [json|yaml]")
		}
		// Get the home directory of the user to save the output file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return errors.Wrap(err, "failed to determine user home directory")
		}
		connectionString := strings.ReplaceAll(fmt.Sprintf("%v", selectedConnection.Name), " ", "_")

		if outputFormatFlag == "yaml" {
			if output, err = yaml.Marshal(selectedConnection); err != nil {
				return errors.Wrap(err, "failed to format output in YAML")
			}
			if saveFlag {
				fmt.Println("Saving output as YAML file")
				err = os.WriteFile(homeDir+"/.meshery/connection_"+connectionString+".yaml", output, 0644)
				if err != nil {
					return errors.Wrap(err, "failed to save output as YAML file")
				}
				fmt.Println("Output saved as YAML file in ~/.meshery/connection_" + connectionString + ".yaml")
			} else {
				fmt.Print(string(output))
			}
		} else {
			if saveFlag {
				fmt.Println("Saving output as JSON file")
				output, err = json.MarshalIndent(selectedConnection, "", "  ")
				if err != nil {
					return errors.Wrap(err, "failed to format output in JSON")
				}
				err = os.WriteFile(homeDir+"/.meshery/connection_"+connectionString+".json", output, 0644)
				if err != nil {
					return errors.Wrap(err, "failed to save output as JSON file")
				}
				fmt.Println("Output saved as JSON file in ~/.meshery/connection_" + connectionString + ".json")
				return nil
			}
			return outputConnectionJSON(selectedConnection)
		}

		return nil
	},
}

// selectConnectionPrompt lets user to select a connection if connections are more than one
func selectConnectionPrompt(connectionsList []*connection.Connection) *connection.Connection {
	connectionNames := []string{}

	for _, conn := range connectionsList {
		connectionName := fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn.ID.String(), conn.Name, conn.Type)
		connectionNames = append(connectionNames, connectionName)
	}

	prompt := promptui.Select{
		Label: "Select connection",
		Items: connectionNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return connectionsList[i]
	}
}

// outputConnectionJSON outputs the connection in JSON format
func outputConnectionJSON(conn *connection.Connection) error {
	// Create a new JSON encoder that writes to the standard output (os.Stdout).
	enc := json.NewEncoder(os.Stdout)
	// Configure the JSON encoder settings.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	return enc.Encode(conn)
}

func init() {
	viewConnectionCmd.Flags().StringVarP(&outputFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewConnectionCmd.Flags().BoolVarP(&saveFlag, "save", "s", false, "(optional) save output as a JSON/YAML file")
}
