package model

import (
	"flag"
	"testing"

	"github.com/spf13/cobra"
)

var update = flag.Bool("update", false, "Update the model")

func TestGenerateModelDataToDisplay_NilResponse(t *testing.T) {
	rows, total := generateModelDataToDisplay(nil)
	if len(rows) != 0 || total != 0 {
		t.Fatalf("expected 0 rows and 0 total on nil response, got %d rows and %d total", len(rows), total)
	}
}

// TestModelArgs verifies that invalid subcommands are rejected regardless of
// --count, while count-only usage remains valid and a missing subcommand
// without --count remains invalid.
func TestModelArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		count   bool
		wantErr bool
	}{
		{
			name:    "rejects missing subcommand without count",
			wantErr: true,
		},
		{
			name:  "accepts count without positional arguments",
			count: true,
		},
		{
			name:    "rejects invalid subcommand without count",
			args:    []string{"invalidsubcommand"},
			wantErr: true,
		},
		{
			name:    "rejects invalid subcommand with count",
			args:    []string{"invalidsubcommand"},
			count:   true,
			wantErr: true,
		},
		{
			name:    "rejects empty subcommand with count",
			args:    []string{""},
			count:   true,
			wantErr: true,
		},
		{
			name:    "rejects multiple invalid arguments with count",
			args:    []string{"invalidsubcommand1", "invalidsubcommand2"},
			count:   true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use:  "resource",
				Args: ModelCmd.Args,
				Run:  func(*cobra.Command, []string) {},
			}
			cmd.Flags().Bool("count", false, "")
			args := append([]string{}, tt.args...)
			if tt.count {
				args = append([]string{"--count"}, args...)
			}
			cmd.SetArgs(args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
