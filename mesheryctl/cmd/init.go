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
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func Init() {
	if err := exec.Command("chmod", "+x", dockerComposeBinary).Run(); err != nil {
		log.Fatal(dockerComposeWarningMessage, err)
		ostype, osarch := Prereq()
		osdetails := strings.TrimRight(string(ostype), "\r\n") + "-" + strings.TrimRight(string(osarch), "\r\n")
		dockerComposeBinaryUrl := dockerComposeBinaryUrl + "-" + osdetails
		if err := DownloadFile(dockerComposeBinary, dockerComposeBinaryUrl); err != nil {
			log.Fatal(err)
		}
		if err := exec.Command("chmod", "+x", dockerComposeBinary).Run(); err != nil {
			log.Fatal(err)
		}
	}

}

var cleanup bool

// startCmd represents the start command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Meshery",
	Long:  `Check and installs docker and docker-compose if not exists`,
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		if cleanup {

			Cleanup()
			os.Exit(0)
		}
		log.Printf(mesheryPostInitMessage)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&cleanup, "cleanup", false, "Delete Mesheryctl")
}
