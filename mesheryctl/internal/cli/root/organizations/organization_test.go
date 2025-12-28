package organizations

import (
	goerrors "errors"
	"flag"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
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
		ExpectedError    error `default:"nil"`
		IsOutputGolden   bool  `default:"true"`
	}{
		{
			Name:             "Launch organization without args",
			Args:             []string{},
			ExpectedResponse: "organization.no.args.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New("please provide a subcommand with the command")),
		},
		{
			Name:             "Launch organization with invalid subcommand name",
			Args:             []string{"test-invalid-subcommand"},
			ExpectedResponse: "organization.invalid.subcommand.golden",
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
