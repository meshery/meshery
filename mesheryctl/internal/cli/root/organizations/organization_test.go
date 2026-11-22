package organizations

import (
	goerrors "errors"
	"flag"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update golden files")

func TestOrganization(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)
	testdataDir := filepath.Join(currentDirectory, "testdata")

	tests := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		ExpectError      bool
		ExpectedError    error
		IsOutputGolden   bool
	}{
		{
			Name:             "Launch organization without args",
			Args:             []string{},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New("please provide a subcommand with the command")),
		},
		{
			Name:             "Launch organization with invalid subcommand name",
			Args:             []string{"test-invalid-subcommand"},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New("'test-invalid-subcommand' is an invalid command. Use 'mesheryctl organization --help' to display usage guide")),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			golden := utils.NewGoldenFile(t, test.ExpectedResponse, testdataDir)

			// setting up log to grab logs
			b := utils.SetupMeshkitLoggerTesting(t, false)
			OrgCmd.SetOut(b)
			OrgCmd.SetArgs(test.Args)
			err := OrgCmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if test.ExpectError {
					if test.IsOutputGolden {
						// write it in file
						if *update {
							golden.Write(err.Error())
						}
						expectedResponse := golden.Load()

						utils.Equals(t, expectedResponse, err.Error())
						return
					}
					assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(test.ExpectedError))
					assert.Equal(t, errors.GetCode(err), errors.GetCode(test.ExpectedError))
					assert.Equal(t, errors.GetLDescription(err), errors.GetLDescription(test.ExpectedError))
					return
				}
				t.Error(err)
			}

			// response being printed in console
			actualResponse := b.String()

			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)
		})
	}
}

// TestOrganizationCountArgs verifies that count-only usage is valid,
// missing subcommands without --count are rejected, and invalid
// subcommands are rejected regardless of --count.
func TestOrganizationCountArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		count   bool
		wantErr string
	}{
		{
			name:    "rejects missing subcommand without count",
			wantErr: "please provide a subcommand",
		},
		{
			name:  "accepts count without arguments",
			count: true,
		},
		{
			name:    "rejects invalid subcommand without count",
			args:    []string{"invalid"},
			wantErr: "invalid command",
		},
		{
			name:    "rejects invalid subcommand with count",
			args:    []string{"invalid"},
			count:   true,
			wantErr: "invalid command",
		},
		{
			name:    "rejects empty subcommand with count",
			args:    []string{""},
			count:   true,
			wantErr: "invalid command",
		},
		{
			name:    "rejects multiple invalid arguments with count",
			args:    []string{"invalid", "extra"},
			count:   true,
			wantErr: "invalid command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use:  "resource",
				Args: OrgCmd.Args,
				Run:  func(*cobra.Command, []string) {},
			}
			cmd.Flags().Bool("count", false, "")
			args := append([]string{}, tt.args...)
			if tt.count {
				args = append([]string{"--count"}, args...)
			}
			cmd.SetArgs(args)

			err := cmd.Execute()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Execute() error = %v", err)
				}
			} else if err == nil || !strings.Contains(fmt.Sprint(errors.GetLDescription(err)), tt.wantErr) {
				t.Fatalf("Execute() error = %+v, want %q", err, tt.wantErr)
			}
		})
	}
}
