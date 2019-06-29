// Copyright 2019 The Meshery Authors
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

package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Meshery status",
	Long:  `Check status of Meshery and Meshery adapters.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Meshery containers status . . .\n")
		out, err := exec.Command("docker-compose", "-f", dockerComposeFile, "ps").Output()
		if err != nil {
			log.Fatal("[ERROR] Please, install docker-compose. The error message: \n", err)
		}
		fmt.Printf("%s\n", out)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
