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
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	testURL  = ""
	testName = ""
	testMesh = ""
	// testFile         = ""
	qps                = ""
	concurrentRequests = ""
	testDuration       = ""
	loadGenerator      = ""
	testCookie         = ""
)

// perfTerminalFormatter declearation
type perfTerminalFormatter struct{}

var perfcfgFile string

var perfcmdDetails = `
Mesheryctl performance management
Usage:
    mesheryctl perf [flag]
Example Usage:
    mesheryctl perf --name "a quick stress test" --url http://192.168.1.15/productpage --qps 300 --concurrent-requests 2 --duration 30s --load-generator wrk2
Available Commands:
    perf     Performance Tests and Benchmarking

Flags:
    –name   (optional) A memorable name for the test.(default) a random string
    –mesh optional) Name of the service mesh. (default) empty string
    –file (optional)    URI to the service mesh performance test configuration file.(default) empty string
    –url (required) URL of the endpoint send load to during testing 
    –qps (optional) Queries per second (default) 0
    –concurrent-requests (optional) Number of concurrent requests(default) 1
    –duration (optional)    Duration of the test.
    –load-generator (optional)  choice of load generator: fortio (OR) wrk2 (default) fortio
`

//Format is exported
func (f *perfTerminalFormatter) Format(entry *log.Entry) ([]byte, error) {
	return append([]byte(entry.Message), '\n'), nil
}

// perfrootCmd represents the base command when called without any subcommands
var perfrootCmd = &cobra.Command{
	Use:   "mesheryctl perf",
	Short: "Meshery Command Line tool",
	Long:  `Mesheryctl performance management`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		b, _ := cmd.Flags().GetBool("version")
		if b {
			versionCmd.Run(nil, nil)
			return
		}
		if len(args) == 0 {
			log.Print(perfcmdDetails)
		}
	},
}

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset generates a random string with a given length
func StringWithCharset(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// perfCmd represents the Performance command
var perfCmd = &cobra.Command{
	Use:   "perf",
	Short: "Performance testing and benchmarking",
	Long:  `Performance testing and benchmarking.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check prerequisite
		preReqCheck()

		const mesheryURL string = "http://localhost:9081/api/load-test-smps?"
		postData := ""

		startTime := time.Now()
		duration, err := time.ParseDuration(testDuration)
		if err != nil {
			println("Error: Test duration invalid")
			return
		}
		endTime := startTime.Add(duration)

		postData = postData + "start_time: " + startTime.Format(time.RFC3339)
		postData = postData + "\nend_time: " + endTime.Format(time.RFC3339)

		if len(testURL) > 0 {
			println("Test name used : ", testName)
			postData = postData + "\nendpoint_url: " + testURL
		} else {
			println(perfcmdDetails)
			return
		}

		postData = postData + "\nclient:"
		postData = postData + "\n connections: " + concurrentRequests
		postData = postData + "\n rps: " + qps

		req, err := http.NewRequest("POST", mesheryURL, bytes.NewBuffer([]byte(postData)))
		if err != nil {
			println("Error in building the request")
			return
		}
		cookieConf := strings.SplitN(testCookie, "=", 2)
		cookieName := cookieConf[0]
		cookieValue := cookieConf[1]
		req.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})
		q := req.URL.Query()
		q.Add("name", testName)
		q.Add("loadGenerator", loadGenerator)
		if len(testMesh) > 0 {
			q.Add("mesh", testMesh)
		}
		req.URL.RawQuery = q.Encode()

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}

		buf := make([]byte, 4)
		for {
			n, err := resp.Body.Read(buf)
			fmt.Print(string(buf[:n]))
			if err == io.EOF {
				break
			}
		}
		println("\nTest Completed Successfully!")
	},
}

func init() {
	perfCmd.Flags().StringVar(&testURL, "url", "", "(required) URL of the endpoint to use for the test")
	perfCmd.Flags().StringVar(&testName, "name", StringWithCharset(8), "(optional) A memorable name for the test.")
	perfCmd.Flags().StringVar(&testMesh, "mesh", "", "(optional) Name of the service mesh.")
	// perfCmd.Flags().StringVar(&testFile, "file", "", "DESCRIPTION")
	perfCmd.Flags().StringVar(&qps, "qps", "0", "(optional) Queries per second")
	perfCmd.Flags().StringVar(&concurrentRequests, "concurrent-requests", "1", "DESCRIPTION")
	perfCmd.Flags().StringVar(&testDuration, "duration", "30s", "(optional) Duration of the test like 10s, 5m, 2h. We are following the convention described at https://golang.org/pkg/time/#ParseDuration")
	perfCmd.Flags().StringVar(&testCookie, "cookie", "meshery-provider=Default Local Provider", "(required) identification of choice of provider.")
	perfCmd.Flags().StringVar(&loadGenerator, "load-generator", "fortio", " (optional) choice of load generator: fortio (OR) wrk2")
	rootCmd.AddCommand(perfCmd)
}
