package organizations

import (
	goerrors "errors"
	"reflect"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestOrganization(t *testing.T) {
	tests := []struct {
		Name             string
		Args             []string
		ExpectError      bool
		ExpectedError    error `default:"nil"`
	}{
		{
			Name:             "Launch organization without args",
			Args:             []string{},
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New("please provide a subcommand with the command")),
		},
		{
			Name:             "Launch organization with invalid subcommand name",
			Args:             []string{"test-invalid-subcommand"},
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New("'test-invalid-subcommand' is an invalid command. Use 'mesheryctl organization --help' to display usage guide")),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// setting up log to grab logs
			b := utils.SetupMeshkitLoggerTesting(t, false)
			OrgCmd.SetOut(b)
			OrgCmd.SetArgs(test.Args)
			err := OrgCmd.Execute()
			if err != nil {
				if !test.ExpectError {
					t.Fatal(err)
				}
				assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(test.ExpectedError))
				assert.Equal(t, errors.GetCode(err), errors.GetCode(test.ExpectedError))
				assert.Equal(t, errors.GetLDescription(err), errors.GetLDescription(test.ExpectedError))
				return
			}
			if test.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}
		})
	}
}
