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
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Meshery",
	Long:  `Run 'docker-compose' to start Meshery and each of its service mesh adapters.`,
	Run: func(cmd *cobra.Command, args []string) {
		var out bytes.Buffer
		var stderr bytes.Buffer

		if _, err := os.Stat(mesheryLocalFolder); os.IsNotExist(err) {
			os.Mkdir(mesheryLocalFolder, 0777)
		}

		if err := downloadFile(dockerComposeFile, fileURL); err != nil {
			log.Fatal(err)
		}

		log.Info("Starting Meshery...")
		start := exec.Command("docker-compose", "-f", dockerComposeFile, "up", "-d")
		start.Stdout = &out
		start.Stderr = &stderr

		if err := start.Run(); err != nil {
			log.Println(stderr.String())
			return
		}

		checkFlag := 0 //flag to check

		//connection to docker-client
		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		//check for container meshery_meshery_1 running status
		for _, container := range containers {
			if "/meshery_meshery_1" == container.Names[0] {
				log.Info("Opening Meshery in your broswer. If Meshery does not open, please point your browser to http://localhost:9081 to access Meshery.")
				ostype, err := exec.Command("uname", "-s").Output()
				if err != nil {
					log.Fatal("[WARNING] Unable to detect OS type. Warning message: \n", err)
				}
				os := strings.TrimSpace(string(ostype))
				// Link to Meshery User Interface
				url := "http://localhost:9081"
				if os == "Linux" {
					// Meshery running on Linux host
					exec.Command("xdg-open", url).Start()
				} else {
					// Asssume Meshery running on MacOS host
					exec.Command("open", url).Start()
				}
				checkFlag = 0
				break
			} else {
				checkFlag = 1
			}
		}

		//if meshery_meshery_1 failed to start showing logs
		//code for logs
		if checkFlag == 1 {
			log.Info("Starting Meshery logging . . .")
			cmdlog := exec.Command("docker-compose", "-f", dockerComposeFile, "logs", "-f")
			cmdReader, err := cmdlog.StdoutPipe()
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(cmdReader)
			go func() {
				for scanner.Scan() {
					log.Println(scanner.Text())
				}
			}()
			if err := cmdlog.Start(); err != nil {
				log.Fatal(err)
			}
			if err := cmdlog.Wait(); err != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
