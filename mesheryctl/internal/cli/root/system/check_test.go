package system

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestPreflightCmd(t *testing.T) {
	SetupContextEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "default-check-output.golden",
			MockResponse:     "default-check-input.golden",
		},
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--preflight"},
			ExpectedResponse: "default-check-output.golden",
			MockResponse:     "default-check-input.golden",
		},
		// {
		// 	Name:             "Docker not available",
		// 	Args:             []string{"check", "--pre"},
		// 	ExpectedResponse: "docker-not-check-output.golden",
		// 	MockResponse:     "docker-not-check-input.golden",
		// },
		// {
		// 	Name:             "Kubernetes not available",
		// 	Args:             []string{"check", "--pre"},
		// 	ExpectedResponse: "k8s-not-check-output.golden",
		// 	MockResponse:     "k8s-not-check-input.golden",
		// },
		// {
		// 	Name:             "No requirments met",
		// 	Args:             []string{"check", "--pre"},
		// 	ExpectedResponse: "no-req-met-check-output.golden",
		// 	MockResponse:     "no-req-met-check-input.golden",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			SystemCmd.SetArgs(tt.Args)
			SystemCmd.SetOut(rescueStdout)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			actualResponse := string(out)

			// get current directory
			_, filename, _, ok := runtime.Caller(0)
			if !ok {
				t.Fatal("problems recovering caller information")
			}

			currDir := filepath.Dir(filename)
			testdataDir := filepath.Join(currDir, "testdata")
			// testdataDir := filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, filepath.Join(testdataDir, "check"))

			if *update {
				golden.Write(actualResponse)
			}
			t.Log(actualResponse)
			expectedResponse := golden.Load()

			if expectedResponse != actualResponse {
				t.Fatalf("diff: %v", utils.Diff(expectedResponse, actualResponse))
				// t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}
		})
	}
}
