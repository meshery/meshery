package filter

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// TestDeleteCmdArgs pins the argument validator of `filter delete`:
// the command requires exactly 1 positional argument (filter-name or ID).
// Zero arguments or more than one argument must be rejected before the handler runs.
func TestDeleteCmdArgs(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectErr   bool
		expectedErr string
	}{
		{
			name:        "zero arguments rejected",
			args:        []string{},
			expectErr:   true,
			expectedErr: "accepts 1 arg(s), received 0",
		},
		{
			name:        "one argument accepted",
			args:        []string{"filter-name"},
			expectErr:   false,
			expectedErr: "",
		},
		{
			name:        "multiple arguments rejected",
			args:        []string{"filter1", "filter2"},
			expectErr:   true,
			expectedErr: "accepts 1 arg(s), received 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteCmd.Args(deleteCmd, tc.args)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error for args %v, got nil", tc.args)
				}
				if err.Error() != tc.expectedErr {
					t.Errorf("expected error %q for args %v, got %q", tc.expectedErr, tc.args, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error for args %v, got %v", tc.args, err)
				}
			}
		})
	}
}

// TestDeleteCmdExecutionRejectsInvalidArgs tests that executing the command with
// zero or more than one argument fails at the Cobra argument validation layer
// with the exact argument-count validation error, while restoring shared command state.
func TestDeleteCmdExecutionRejectsInvalidArgs(t *testing.T) {
	testCases := []struct {
		name              string
		args              []string
		expectedErrSubstr string
	}{
		{
			name:              "reject zero arguments",
			args:              []string{"delete"},
			expectedErrSubstr: "accepts 1 arg(s), received 0",
		},
		{
			name:              "reject multiple arguments",
			args:              []string{"delete", "filter1", "filter2"},
			expectedErrSubstr: "accepts 1 arg(s), received 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer
			FilterCmd.SetOut(&outBuf)
			FilterCmd.SetErr(&errBuf)
			FilterCmd.SetArgs(tc.args)

			t.Cleanup(func() {
				FilterCmd.SetArgs(nil)
				FilterCmd.SetOut(nil)
				FilterCmd.SetErr(nil)
			})

			err := FilterCmd.Execute()
			if err == nil {
				t.Fatalf("expected error executing 'filter %v', got nil", tc.args)
			}
			if !strings.Contains(err.Error(), tc.expectedErrSubstr) {
				t.Errorf("expected error containing %q, got %q", tc.expectedErrSubstr, err.Error())
			}
		})
	}
}

func TestDeleteCmd(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	testcase := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "Delete Kuma-Test",
			Args: []string{"delete", "c0c6035a-b1b9-412d-aab2-4ed1f1d51f84"},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/filter/c0c6035a-b1b9-412d-aab2-4ed1f1d51f84",
					Response:     "delete.kuma.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/filter?page_size=10000",
					Response:     "filter.list.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "delete.kuma.output.golden",
			ExpectError:      false,
		}, {
			Name: "Delete RolloutAndIstio",
			Args: []string{"delete", "RolloutAndIstio"},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/filter/d0e09134-acb6-4c71-b051-3d5611653f70",
					Response:     "delete.rollout.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/filter?page_size=10000",
					Response:     "filter.list.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "delete.rollout.output.golden",
			ExpectError:      false,
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, FilterCmd, testcase, currDir, "filter", func() {})
}
