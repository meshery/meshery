package system

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestContextViewCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/TestContext.yaml")

	tests := []utils.CmdTestInput{
		{
			Name:             "view for default context",
			Args:             []string{"context", "view"},
			ExpectedResponse: "view.golden",
		},
		{
			Name:             "view with specified context through context flag",
			Args:             []string{"context", "view", "--context", "local2"},
			ExpectedResponse: "viewWithContextExpected.golden",
		},
		{
			Name:             "Error for viewing a non-existing context",
			Args:             []string{"context", "view", "local3"},
			IsOutputGolden:   false,
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    ErrContextNotExists(fmt.Errorf("context `local3` does not exist")),
		},
		{
			Name:             "view with specified context as argument",
			Args:             []string{"context", "view", "local2"},
			ExpectedResponse: "viewWithContextExpected.golden",
		},
		{
			Name:             "view with all flag set",
			Args:             []string{"context", "view", "--all"},
			ExpectedResponse: "viewAllExpected.golden",
		},
	}

	mesheryctlflags.InitValidators(SystemCmd)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected error, got nil")
			}

			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			actualResponse := buf.String()

			if *update {
				golden.Write(actualResponse)
			}

			expectedResponse := golden.Load()
			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedExceptedResponse, cleanedActualResponse)

		})
		t.Logf("List %s test", "context")
	}
}

func TestContextListCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/TestContext.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "list all contexts",
			Args:             []string{"context", "list"},
			ExpectedResponse: "listExpected.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			//t.Logf("Expected response:\n%s", expectedResponse)
			//t.Logf("Actual response:\n%s", actualResponse)

			assert.Equal(t, expectedResponse, actualResponse)
		})
		t.Log("ListContextCmd test passed")
	}
}

func TestContextDeleteCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedDelete.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "delete given context",
			Args:             []string{"context", "delete", "local2"},
			ExpectedResponse: "delete.context.golden",
		},
	}

	mesheryctlflags.InitValidators(SystemCmd)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			if expectedResponse != actualResponse {
				t.Error("Expected response not obtained")
				t.Errorf("Expected: %v", expectedResponse)
				t.Errorf("Actual: %v", actualResponse)
			}
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedDelete.yaml"

			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, "deleteExpected.golden", testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			deleteExpected := golden.Load()
			if actualResponse != deleteExpected {
				t.Error("Contexts are mismatched")
				t.Error("Expected:")
				t.Errorf("%v", deleteExpected)
				t.Error("Actual:")
				t.Errorf("%v", actualResponse)
			}

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("DeleteContextCmd test Passed")
	}
}

func TestContextCreateCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedAdd.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:                 "given context name provided when system context create then context is created",
			Args:                 []string{"context", "create", "local3"},
			ExpectedResponse:     "createContext.golden",
			ExpectedResponseYaml: "addExpected.golden",
		},
		{
			Name:           "given no context provided when system context create  then thorw error",
			Args:           []string{"context", "create"},
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errArgMsg, contextCreateUsageMsg)),
			IsOutputGolden: false,
		},
		{
			Name:           "given multiple context name provided when system context create then throw error",
			Args:           []string{"context", "create", "local1", "local2"},
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errArgMsg, contextCreateUsageMsg)),
			IsOutputGolden: false,
		},
	}

	mesheryctlflags.InitValidators(SystemCmd)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected error, got nil")
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedAdd.yaml"

			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, tt.ExpectedResponseYaml, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			addExpected := golden.Load()
			assert.Equal(t, addExpected, actualResponse)

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("CreateContextCmd test Passed")
	}

}

func TestContextAddUppercaseCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedAdd.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:                 "given context name which contains uppercase provided when system context create then context is created in lowercase",
			Args:                 []string{"context", "create", "Local3"},
			ExpectedResponse:     "createContext.golden",
			ExpectedResponseYaml: "addExpected.golden",
		},
	}

	mesheryctlflags.InitValidators(SystemCmd)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected error, got nil")
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedAdd.yaml"

			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, tt.ExpectedResponseYaml, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			addExpected := golden.Load()
			assert.Equal(t, addExpected, actualResponse)

			// To check context is lowercase
			contentFile, err := os.ReadFile(filepath)
			if err != nil {
				t.Fatal(err)
			}

			var data map[string]interface{}
			err = yaml.Unmarshal(contentFile, &data)
			if err != nil {
				t.Fatal(err)
			}

			contexts := data["contexts"].(map[string]interface{})

			if len(tt.Args) > 2 {
				_, existsUpper := contexts[tt.Args[2]]
				if existsUpper {
					t.Fatalf("uppercase context should not exist")
				}
			}

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("CreateUppercaseContextCmd test Passed")
	}

}

func TestContextSwitchCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedSwitch.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "switch to a different context",
			Args:             []string{"context", "switch", "local2"},
			ExpectedResponse: "switch.context.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)

			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedSwitch.yaml"
			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, "switchExpected.golden", testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			switchExpected := golden.Load()
			assert.Equal(t, switchExpected, actualResponse)

			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("SwitchContextCmd test passed")
	}
}

func TestContextPingCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/TestContext.yaml")

	// TestContext.yaml's contexts reference stored tokens ("default",
	// "default2") that attachContextAuthDetails must resolve to a real file
	// on disk; point MesheryFolder at the fixture directory containing
	// auth.json so that resolution succeeds. SetFileLocationTesting also
	// mutates DockerComposeFile and AuthConfigFile as a side effect, so all
	// three package-level vars must be restored, not just MesheryFolder, or
	// later tests in this binary (e.g. TestResetCmd) inherit this test's
	// fixture paths.
	origMesheryFolder := utils.MesheryFolder
	origDockerComposeFile := utils.DockerComposeFile
	origAuthConfigFile := utils.AuthConfigFile
	utils.SetFileLocationTesting(currDir)
	t.Cleanup(func() {
		utils.MesheryFolder = origMesheryFolder
		utils.DockerComposeFile = origDockerComposeFile
		utils.AuthConfigFile = origAuthConfigFile
	})

	utils.StartMockery(t)
	defer utils.StopMockery(t)

	mesheryctlflags.InitValidators(SystemCmd)

	t.Run("valid token reports reachable and valid", func(t *testing.T) {
		contextPingFlags.Context = ""
		httpmock.RegisterResponder("GET", "http://localhost:9081/api/user",
			httpmock.NewStringResponder(200, `{"email":"alice@example.com"}`))

		buf := utils.SetupMeshkitLoggerTesting(t, false)
		SystemCmd.SetOut(buf)
		SystemCmd.SetErr(buf)
		SystemCmd.SetArgs([]string{"context", "ping"})
		err := SystemCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error for a valid token, got: %v", err)
		}

		out := buf.String()
		for _, want := range []string{
			"Context: local",
			"Endpoint: http://localhost:9081",
			"✅ Server reachable",
			"✅ Token valid (user: alice@example.com)",
		} {
			if !strings.Contains(out, want) {
				t.Fatalf("expected output to contain %q, got: %s", want, out)
			}
		}
	})

	t.Run("expired or invalid token is reported, not mistaken for unreachable", func(t *testing.T) {
		contextPingFlags.Context = ""
		httpmock.RegisterResponder("GET", "http://localhost:32242/api/user",
			httpmock.NewStringResponder(401, `{"error":"unauthorized"}`))

		buf := utils.SetupMeshkitLoggerTesting(t, false)
		SystemCmd.SetOut(buf)
		SystemCmd.SetErr(buf)
		SystemCmd.SetArgs([]string{"context", "ping", "local2"})
		err := SystemCmd.Execute()
		if err == nil {
			t.Fatal("expected a non-nil error for an invalid/expired token")
		}

		out := buf.String()
		for _, want := range []string{
			"Context: local2",
			"Endpoint: http://localhost:32242",
			"✅ Server reachable",
			"⚠️  Token invalid/expired",
		} {
			if !strings.Contains(out, want) {
				t.Fatalf("expected output to contain %q, got: %s", want, out)
			}
		}
	})

	t.Run("unreachable server is reported without a token verdict", func(t *testing.T) {
		contextPingFlags.Context = ""
		httpmock.RegisterResponder("GET", "http://localhost:9081/api/user",
			httpmock.NewErrorResponder(errors.New("dial tcp: connect: connection refused")))

		buf := utils.SetupMeshkitLoggerTesting(t, false)
		SystemCmd.SetOut(buf)
		SystemCmd.SetErr(buf)
		SystemCmd.SetArgs([]string{"context", "ping"})
		err := SystemCmd.Execute()
		if err == nil {
			t.Fatal("expected a non-nil error for an unreachable server")
		}

		out := buf.String()
		for _, want := range []string{
			"Context: local",
			"Endpoint: http://localhost:9081",
			"❌ Server unreachable: dial tcp: connect: connection refused",
		} {
			if !strings.Contains(out, want) {
				t.Fatalf("expected output to contain %q, got: %s", want, out)
			}
		}
		if strings.Contains(out, "Token valid") || strings.Contains(out, "Token invalid") {
			t.Fatalf("did not expect a token verdict when the server is unreachable, got: %s", out)
		}
	})

	t.Run("token configured but file missing surfaces an error without contacting the server", func(t *testing.T) {
		contextPingFlags.Context = ""

		// Point MesheryFolder at a location with no auth.json, so the
		// context's configured token ("default") cannot be resolved to a
		// real file, without touching the fixture used by other subtests.
		missingTokenDir := t.TempDir()
		origMesheryFolder := utils.MesheryFolder
		utils.MesheryFolder = missingTokenDir
		defer func() { utils.MesheryFolder = origMesheryFolder }()

		callsBefore := httpmock.GetTotalCallCount()

		buf := utils.SetupMeshkitLoggerTesting(t, false)
		SystemCmd.SetOut(buf)
		SystemCmd.SetErr(buf)
		SystemCmd.SetArgs([]string{"context", "ping"})
		err := SystemCmd.Execute()
		if err == nil {
			t.Fatal("expected a non-nil error when the stored token file is missing")
		}
		if !strings.Contains(err.Error(), "unable to attach stored token") {
			t.Fatalf("expected a token-attachment error, got: %v", err)
		}

		if got := httpmock.GetTotalCallCount(); got != callsBefore {
			t.Fatalf("expected no network call when the token file could not be attached, call count went from %d to %d", callsBefore, got)
		}

		out := buf.String()
		for _, want := range []string{
			"Context: local",
			"Endpoint: http://localhost:9081",
		} {
			if !strings.Contains(out, want) {
				t.Fatalf("expected output to contain %q, got: %s", want, out)
			}
		}
		if strings.Contains(out, "Server reachable") {
			t.Fatalf("did not expect a reachability verdict when the token could not be attached, got: %s", out)
		}
	})

	t.Run("non-loopback http endpoint refuses to send the stored token", func(t *testing.T) {
		contextPingFlags.Context = ""

		// Switch to a fixture whose only context is a non-loopback http://
		// endpoint, so the new scheme guard in attachContextAuthDetails is
		// the thing that stops the request - not a missing token or file.
		// Restore the fixture used by the rest of this test afterward.
		utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/NonLoopbackContext.yaml")
		defer utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/TestContext.yaml")

		callsBefore := httpmock.GetTotalCallCount()

		buf := utils.SetupMeshkitLoggerTesting(t, false)
		SystemCmd.SetOut(buf)
		SystemCmd.SetErr(buf)
		SystemCmd.SetArgs([]string{"context", "ping"})
		err := SystemCmd.Execute()
		if err == nil {
			t.Fatal("expected a non-nil error for a non-loopback http endpoint")
		}
		if !strings.Contains(err.Error(), "refusing to send stored token over plaintext HTTP to non-loopback endpoint") {
			t.Fatalf("expected a plaintext-HTTP refusal error, got: %v", err)
		}

		if got := httpmock.GetTotalCallCount(); got != callsBefore {
			t.Fatalf("expected no network call when the endpoint is refused, call count went from %d to %d", callsBefore, got)
		}

		out := buf.String()
		if strings.Contains(out, "Server reachable") {
			t.Fatalf("did not expect a reachability verdict when the endpoint was refused, got: %s", out)
		}
	})

	t.Run("nonexistent context name errors out", func(t *testing.T) {
		contextPingFlags.Context = ""
		buf := utils.SetupMeshkitLoggerTesting(t, false)
		SystemCmd.SetOut(buf)
		SystemCmd.SetErr(buf)
		SystemCmd.SetArgs([]string{"context", "ping", "does-not-exist"})
		err := SystemCmd.Execute()
		utils.AssertMeshkitErrorsEqual(t, err, ErrContextNotExists(fmt.Errorf("context `does-not-exist` does not exist")))
	})
}

func resetVariables() {
	//reset context before tests
	tempCntxt = ""
}
