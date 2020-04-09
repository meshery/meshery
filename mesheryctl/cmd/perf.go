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

	"github.com/spf13/cobra"
)

var (
	testURL            = ""
	testName           = ""
	testMesh           = ""
	qps                = ""
	concurrentRequests = ""
	testDuration       = ""
	loadGenerator      = ""
	testCookie         = ""
)

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
	Short: "Performance Testing",
	Long:  `Performance Testing & Benchmarking using Meshery CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check prerequisite
		preReqCheck()

		if len(testName) <= 0 {
			println("Test Name not provided \n")
			testName = StringWithCharset(8)
			println("Using random test name: ", testName)
		}

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
			postData = postData + "\nendpoint_url: " + testURL
		} else {
			println("Error: Please enter a TestURL")
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
	perfCmd.Flags().StringVar(&testURL, "url", "", "(required) Endpoint URL to test")
	perfCmd.Flags().StringVar(&testName, "name", "", "(optional) Name of the Test")
	perfCmd.Flags().StringVar(&testMesh, "mesh", "", "(optional) Name of the Service Mesh")
	perfCmd.Flags().StringVar(&qps, "qps", "0", "(optional) Queries per second")
	perfCmd.Flags().StringVar(&concurrentRequests, "concurrent-requests", "1", "(required) Number of Parallel Requests")
	perfCmd.Flags().StringVar(&testDuration, "duration", "30s", "(optional) Duration of the Test")
	perfCmd.Flags().StringVar(&testCookie, "cookie", "meshery-provider=Default Local Provider", "(required) Choice of Provider")
	perfCmd.Flags().StringVar(&loadGenerator, "load-generator", "fortio", "(optional) Load-Generator to be used (fortio/wrk2)")
	rootCmd.AddCommand(perfCmd)
}
