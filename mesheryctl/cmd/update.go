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
	"encoding/json"
	"fmt"
	"github.com/artdarek/go-unzip"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshery images from Docker Hub and mesheryctl client.",
	Long:  `Pull new Meshery images from Docker Hub. Pulls new mesheryctl client. This command may be executed while Meshery is running.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Updating Meshery now...")

		//Fetch latest version binary of mesheryctl
		version := "" //To store lateset mesheryctl version info
		resp, err := http.Get(mesheryURL)
		downloadMesheryURL := downloadMesheryURL
		if err != nil {
			log.Print("error")
		} else {
			var dat map[string]interface{}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			if err := json.Unmarshal(body, &dat); err != nil {
				log.Fatal(err)
			}
			num := dat["tag_name"]
			os := ""
			arch := ""

			if runtime.GOOS == "windows" {
				// Meshery running on Windows host
				os = "Windows"
			} else if runtime.GOOS == "linux" {
				// Meshery running on Linux host
				os = "Linux"
			} else {
				// Assume Meshery running on MacOS host
				os = "Darwin"
			}

			if runtime.GOARCH == "amd64" {
				// Meshery running on x86_64 host
				arch = "x86_64"
			} else {
				// Meshery running on i386 host
				arch = "i386"
			}

			version = fmt.Sprint("%v", num)
			downloadMesheryURL = fmt.Sprintf(downloadMesheryURL+"/%v/mesheryctl_%v_%v_%v.zip", num, version[3:], os, arch)
		}

		//If sudo is needed to update mesheryctl
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		if user.Username == "root" && os.Getenv("SUDO_USER") != "" {
			if _, err := os.Stat(mesheryFolder); os.IsNotExist(err) {
				_ = os.Mkdir(mesheryFolder, 0777)
			}
		}

		//download mehseryctl binary
		if version != Build {
			log.Info("Downloading latest mesheryctl..")
			if err := downloadFile(mesherybinary, downloadMesheryURL); err != nil {
				log.Warn("Unable to download mesheryctl release bundle: ", err)
			} else {
				//Unzip mesheryctl binary
				uz := unzip.New(mesherybinary, "/usr/local/bin/")
				err = uz.Extract()
				if err != nil {
					log.Warn(err)
					log.Warn("Skipping mesheryctl binary update!")
				}
			}
		} else {
			log.Info("Latest mesheryctl already present!\nSkipping mesheryctl upgradation!")
		}

		//update docker image
		if _, err := os.Stat(dockerComposeFile); os.IsNotExist(err) {
			if err := downloadFile(dockerComposeFile, fileURL); err != nil {
				log.Fatal("Could not update Meshery containers: ", err)
			}
		}

		start := exec.Command("docker-compose", "-f", dockerComposeFile, "pull")
		start.Stdout = os.Stdout
		start.Stderr = os.Stderr
		if err := start.Run(); err != nil {
			log.Fatal(err)
		}

		log.Info("Meshery is now up-to-date")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
