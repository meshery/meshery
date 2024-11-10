package experimental

import (
	"testing"
)

func TestExperimentalList(t *testing.T) {
	// test scenrios
	tests := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		ExpecteError     bool
	}{
		{
			Name:             "List registered relationships",
			Args:             []string{"relationship", "list"},
			ExpectedResponse: "",
			ExpecteError:     false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ExpCmd.SetArgs(tt.Args)
			err := ExpCmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tt.ExpecteError {
					// if the error was expected, we've to check that if the generated error is equal to the "ExpectedResponse"
					return
				}
			}
		})
	}
}
