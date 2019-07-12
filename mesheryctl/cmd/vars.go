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
)

const (
	url                         = "http://localhost:9081"
	fileUrl                     = "https://raw.githubusercontent.com/layer5io/meshery/master/docker-compose.yaml"
	mesheryLocalFolder          = ".meshery"
	dockerComposeFile           = mesheryLocalFolder + "/meshery.yaml"
	dockerComposeBinaryUrl      = "https://github.com/docker/compose/releases/download/1.24.1/docker-compose"
	dockerComposeBinary         = "/usr/local/bin/docker-compose"
	dockerComposeWarningMessage = "Prerequisite Docker Compose not available. Attempting Docker Compose installation… \n"
	ostypeWarningMessage        = "[Warning] Cannot detect Operating System type \n"
	osarchWarningMessage        = "[Warning] Cannot detect Operating System architecture \n"
	mesheryPostStartMessage     = "Opening Meshery in your broswer. If Meshery does not open, please point your browser to http://localhost:9081 to access Meshery."
	mesheryPreStartMessage      = "Starting Meshery . . ."
	mesheryPreStopMessage       = "Stopping Meshery now . . ."
	mesheryPostStoptMessage     = "Meshery is stopped"
	mesheryPreUpdateMessage     = "Updating Meshery now . . ."
	mesheryPostUpdateMessage    = "Meshery is now up-to-date"
	mesheryStatusMessage        = "Meshery containers status . . .\n"
	mesheryPreLogMessage        = "Starting Meshery logging . . ."
	mesheryPreCleanupMessage    = "Cleaning old Meshery config . . ."
	mesheryPostCleanupMessage   = "Meshery config is now cleaned up. \n"
	mesheryPostInitMessage      = "Prerequisite Docker Compose is installed. \n"
)

func DownloadFile(filepath string, url string) error {
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

func Prereq() ([]byte, []byte) {
	ostype, err := exec.Command("uname", "-s").Output()
	if err != nil {
		log.Fatal(ostypeWarningMessage, err)
	}
	//fmt.Printf("%s\n", ostype)

	osarch, err := exec.Command("uname", "-m").Output()
	if err != nil {
		log.Fatal(osarchWarningMessage, err)
	}
	//	fmt.Printf("%s\n", arch)
	return ostype, osarch
}
