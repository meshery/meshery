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
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Meshery",
	Long:  `Run 'docker-compose' to spin up Meshery and each of each service mesh adapters.`,
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := os.Stat(mesheryLocalFolder); os.IsNotExist(err) {
			os.Mkdir(mesheryLocalFolder, 0777)
		}

		if err := DownloadFile(dockerComposeFile, fileUrl); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Starting Meshery . . .")
		if err := exec.Command("docker-compose", "-f", dockerComposeFile, "up", "-d").Run(); err != nil {
			log.Fatal("[ERROR] Prerequisite Docker Compose not available. Error message: \n", err)
		}

		if err := exec.Command("docker-compose", "-f", dockerComposeFile, "up", "-d").Run(); err != nil {
			log.Fatal("[ERROR] Prerequisite Docker Compose not available. Error message: \n", err)
		}

		fmt.Println("Opening Meshery in your broswer. If Meshery does not open, please point your browser to http://localhost:9081 to access Meshery.")
		exec.Command("open", url).Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
