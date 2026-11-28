package system

import (
	"fmt"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestConfigSubcommandsArgs(t *testing.T) {
	tests := []struct {
		Name          string
		Cmd           *cobra.Command
		Args          []string
		ExpectError   bool
		ExpectedError error
	}{
		// aks
		{
			Name:        "aks valid with no positional arguments",
			Cmd:         aksConfigCmd,
			Args:        []string{},
			ExpectError: false,
		},
		{
			Name:          "aks invalid with unexpected positional argument",
			Cmd:           aksConfigCmd,
			Args:          []string{"unexpected-arg"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", aksConfigUsageMsg)),
		},
		{
			Name:          "aks invalid with multiple unexpected positional arguments",
			Cmd:           aksConfigCmd,
			Args:          []string{"arg1", "arg2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", aksConfigUsageMsg)),
		},
		// eks
		{
			Name:        "eks valid with no positional arguments",
			Cmd:         eksConfigCmd,
			Args:        []string{},
			ExpectError: false,
		},
		{
			Name:          "eks invalid with unexpected positional argument",
			Cmd:           eksConfigCmd,
			Args:          []string{"unexpected-arg"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", eksConfigUsageMsg)),
		},
		{
			Name:          "eks invalid with multiple unexpected positional arguments",
			Cmd:           eksConfigCmd,
			Args:          []string{"arg1", "arg2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", eksConfigUsageMsg)),
		},
		// gke
		{
			Name:        "gke valid with no positional arguments",
			Cmd:         gkeConfigCmd,
			Args:        []string{},
			ExpectError: false,
		},
		{
			Name:          "gke invalid with unexpected positional argument",
			Cmd:           gkeConfigCmd,
			Args:          []string{"unexpected-arg"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", gkeConfigUsageMsg)),
		},
		{
			Name:          "gke invalid with multiple unexpected positional arguments",
			Cmd:           gkeConfigCmd,
			Args:          []string{"arg1", "arg2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", gkeConfigUsageMsg)),
		},
		// minikube
		{
			Name:        "minikube valid with no positional arguments",
			Cmd:         minikubeConfigCmd,
			Args:        []string{},
			ExpectError: false,
		},
		{
			Name:          "minikube invalid with unexpected positional argument",
			Cmd:           minikubeConfigCmd,
			Args:          []string{"unexpected-arg"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", minikubeConfigUsageMsg)),
		},
		{
			Name:          "minikube invalid with multiple unexpected positional arguments",
			Cmd:           minikubeConfigCmd,
			Args:          []string{"arg1", "arg2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", minikubeConfigUsageMsg)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Cmd.Args(tt.Cmd, tt.Args)
			if tt.ExpectError {
				assert.Error(t, err)
				assert.Equal(t, utils.ErrInvalidArgumentCode, errors.GetCode(err))
				utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
