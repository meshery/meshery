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
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

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

		if err := DownloadFile(dockerComposeFile, fileUrl); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Starting Meshery...")
		start := exec.Command("docker-compose", "-f", dockerComposeFile, "up", "-d")
		time.Sleep(2 * time.Second)
		//testing purpose
		//exec.Command("docker", "kill", "meshery_meshery_1")
		start.Stdout = &out
		start.Stderr = &stderr

		if err := start.Run(); err != nil {
			fmt.Println(stderr.String())
			return
		}
		// Name of all running container
		cmda := exec.Command("docker", "ps", "--format", "{{.Names}}")
		outa, err := cmda.Output()
		s := strings.Split(string(outa), "\n")
		fmt.Println(s)
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		check(s, 0)
	},
}

func check(s []string, atempt int) {
	count := 0
	for _, x := range s {
		if strings.TrimRight(x, " ") == "meshery_meshery_1" {
			count++
		}
	}
	if count < 1 && atempt < 1 {
		//checks for meshery_meshery_1 container running or not if not then restarting meshery
		fmt.Println("service cannot be started, attempting to restart the service!")
		exec.Command("mesheryctl", "stop")
		exec.Command("mesheryctl", "start")
		atempt++
		check(s, atempt)
	} else if atempt >= 1 {
		//showing logs if restart doesn't work
		fmt.Println("service cannot be started. Showing meshery logs")
		//\
		fmt.Println("Starting Meshery logging . . .")
		cmd_log := exec.Command("docker-compose", "-f", dockerComposeFile, "logs", "-f")
		cmdReader, err := cmd_log.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()
		if err := cmd_log.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd_log.Wait(); err != nil {
			log.Fatal(err)
		}
		//

	} else {
		// If meshery_meshery_1 container runs fine then open web browser
		fmt.Println("Opening Meshery in your broswer. If Meshery does not open, please point your browser to http://localhost:9081 to access Meshery.")
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
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
