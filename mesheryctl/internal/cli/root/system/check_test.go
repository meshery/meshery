package system

import (
	"flag"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSystemCheckFlagTagsMatchCobraFlags(t *testing.T) {
	flagStructType := reflect.TypeOf(cmdSystemCheckFlags{})
	for i := 0; i < flagStructType.NumField(); i++ {
		field := flagStructType.Field(i)
		jsonTag := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		assert.NotNil(t, checkCmd.Flags().Lookup(jsonTag), "expected Cobra flag --%s for struct field %s", jsonTag, field.Name)
	}
}

var update = flag.Bool("update", false, "update golden files")

// This is an Integration test
func TestPreflightCmdIntegration(t *testing.T) {
	// skipping this integration test with --short flag
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	setupContextTestEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "check.output.golden",
		},
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--preflight"},
			ExpectedResponse: "check.output.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// setting up log to grab logs
			buf := setupSystemOutCmdTest(t)
			defer func() {
				buf.Reset()
				// resetCmdFlags(SystemCmd, t)
			}()

			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			output := buf.String()
			actualResponse := output

			// get current directory
			_, filename, _, ok := runtime.Caller(0)
			if !ok {
				t.Fatal("Not able to get current working directory")
			}

			currDir := filepath.Dir(filename)
			testdataDir := filepath.Join(currDir, "testdata")
			// testdataDir := filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, filepath.Join(testdataDir, "check"))

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)
		})
		t.Log("PreflightCmdIntegration Test Passed")
	}
}
