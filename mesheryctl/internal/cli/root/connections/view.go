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
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type connectionViewFlags struct {
	outputFormat string
	save         bool
}

var connectionViewFlagsProvided connectionViewFlags

var (
	validOutputFormats = []string{"json", "yaml"}
)

var viewConnectionCmd = &cobra.Command{
	Use:   "view",
	Short: "View a connection",
	Long: `View a connection by its ID or name.
Documentation for viewing connection can be found at https://docs.meshery.io/reference/mesheryctl/connection/view`,
	Example: `
// View details of a specific connection in default format (yaml)
mesheryctl connection view [connection-name|connection-id]

// View details of a specific connection in JSON format
mesheryctl connection view [connection-name|connection-id] --output-format json

// View details of a specific connection in json format and save it to a file
mesheryctl connection view [connection-name|connection-id] --output-format json --save
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl connection view [connection-name]\nRun 'mesheryctl connection view --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("connection name or ID isn't specified\n\n%v", errMsg))
		}

		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments\n\n%v", errMsg))
		}

		if !slices.Contains(validOutputFormats, strings.ToLower(connectionViewFlagsProvided.outputFormat)) {
			return utils.ErrInvalidArgument(errors.New(invalidOutputFormatMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		connectionNameOrID := args[0]

		var selectedConnection *connection.Connection

		if isArgumentUUID(connectionNameOrID) {
			fetchedConnection, err := fetchConnectionByID(connectionNameOrID)
			if err != nil {
				return err
			}
			selectedConnection = fetchedConnection
		} else {
			fetchedConnection, err := fetchConnectionByName(connectionNameOrID)
			if err != nil {
				return err
			}

			if fetchedConnection == nil {
				fmt.Println("No connection(s) found for the given name: ", connectionNameOrID)
				return nil
			}

			selectedConnection = fetchedConnection
		}

		// Get the home directory of the user to save the output file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return utils.ErrRetrieveHomeDir(errors.Wrap(err, "failed to determine user home directory"))
		}

		outputFormatterFactory := display.OutputFormatterFactory[connection.Connection]{}
		outputFormatter, err := outputFormatterFactory.New(connectionViewFlagsProvided.outputFormat, *selectedConnection)
		if err != nil {
			return err
		}

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		if connectionViewFlagsProvided.save {
			// Prepare the connection string for file naming since connection from local provider
			// can be created without a name.
			connectionString := func(c connection.Connection) string {
				if c.Name == "" {
					return c.ID.String()
				}
				return strings.ReplaceAll(fmt.Sprintf("%v", c.Name), " ", "_")
			}(*selectedConnection)

			fileName := fmt.Sprintf("connection_%s.%s", connectionString, strings.ToLower(connectionViewFlagsProvided.outputFormat))
			file := filepath.Join(homeDir, ".meshery", fileName)
			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[connection.Connection]{}
			outputFormatterSaver, err := outputFormatterSaverFactory.New(connectionViewFlagsProvided.outputFormat, outputFormatter)
			if err != nil {
				return err
			}
			outputFormatterSaver = outputFormatterSaver.WithFilePath(file)
			err = outputFormatterSaver.Save()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

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

func isArgumentUUID(arg string) bool {
	_, err := uuid.FromString(arg)
	return err == nil
}

func fetchConnectionByID(connectionID string) (*connection.Connection, error) {
	url := fmt.Sprintf("%s/%s", connectionApiPath, connectionID)
	fetchedConnection, err := api.Fetch[connection.Connection](url)
	if err != nil {
		return nil, err
	}
	return fetchedConnection, nil
}

func fetchConnectionByName(connectionName string) (*connection.Connection, error) {
	viewUrlValue := url.Values{}
	viewUrlValue.Add("search", connectionName)
	viewUrlValue.Add("pagesize", "all")

	urlPath := fmt.Sprintf("%s?%s", connectionApiPath, viewUrlValue.Encode())

	connectionsResponse, err := api.Fetch[connection.ConnectionPage](urlPath)

	if err != nil {
		return nil, err
	}

	if connectionsResponse.TotalCount == 0 {
		return nil, nil
	}

	if connectionsResponse.TotalCount > 1 {
		return selectConnectionPrompt(connectionsResponse.Connections), nil
	}

	return connectionsResponse.Connections[0], nil
}

func init() {
	viewConnectionCmd.Flags().StringVarP(&connectionViewFlagsProvided.outputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewConnectionCmd.Flags().BoolVarP(&connectionViewFlagsProvided.save, "save", "s", false, "(optional) save output as a JSON/YAML file")
}
