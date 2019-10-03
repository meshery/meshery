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
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

const (
	url                         = "http://localhost:9081"
	fileURL                     = "https://raw.githubusercontent.com/layer5io/meshery/master/docker-compose.yaml"
	dockerComposeWebURL         = "https://api.github.com/repos/docker/compose/releases/latest"
	defaultDockerComposeVersion = "1.24.1/docker-compose"
	// dockerComposeBinaryURL = "https://github.com/docker/compose/releases/download/1.24.1/docker-compose"
	dockerComposeBinaryURL = "https://github.com/docker/compose/releases/download/"
	dockerComposeBinary    = "/usr/local/bin/docker-compose"
)

// See setFileLocation function below.
var (
	mesheryFolder     = ".meshery"
	dockerComposeFile = "/meshery.yaml"
)

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func prereq() ([]byte, []byte) {
	ostype, err := exec.Command("uname", "-s").Output()
	if err != nil {
		log.Fatal("[ERROR] Please, install docker-compose. The error message: \n", err)
	}
	//fmt.Printf("%s\n", ostype)

	osarch, err := exec.Command("uname", "-m").Output()
	if err != nil {
		log.Fatal("[ERROR] Please, install docker-compose. The error message: \n", err)
	}
	//	fmt.Printf("%s\n", arch)
	return ostype, osarch
}

func setFileLocation() {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("[ERROR] Cannot determine location of $HOME")
	}
	mesheryFolder = path.Join(home, mesheryFolder)
	dockerComposeFile = path.Join(mesheryFolder, dockerComposeFile)
}
