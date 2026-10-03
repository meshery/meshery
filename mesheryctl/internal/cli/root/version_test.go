package root

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func TestVersionCmdRegistersOutputFormatFlag(t *testing.T) {
	flag := versionCmd.Flag("output-format")
	if flag == nil {
		t.Fatal("version command does not register output-format")
	}
	if flag.Shorthand != "o" {
		t.Fatalf("output-format shorthand = %q, want %q", flag.Shorthand, "o")
	}
}

func TestDisplayVersionOutputJSON(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	err := displayVersionOutput(cmd, "json", "1.2.3", "client-sha", config.Version{
		Build:     "4.5.6",
		CommitSHA: "server-sha",
	}, nil, nil)
	if err != nil {
		t.Fatalf("displayVersionOutput() error = %v", err)
	}

	var actual map[string]map[string]string
	if err := json.Unmarshal(output.Bytes(), &actual); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if !reflect.DeepEqual(actual, expectedVersionOutput()) {
		t.Errorf("JSON output = %#v, want %#v", actual, expectedVersionOutput())
	}
}

func TestDisplayVersionOutputYAML(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	err := displayVersionOutput(cmd, "yaml", "1.2.3", "client-sha", config.Version{
		Build:     "4.5.6",
		CommitSHA: "server-sha",
	}, nil, nil)
	if err != nil {
		t.Fatalf("displayVersionOutput() error = %v", err)
	}

	var actual map[string]map[string]string
	if err := yaml.Unmarshal(output.Bytes(), &actual); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}
	if !reflect.DeepEqual(actual, expectedVersionOutput()) {
		t.Errorf("YAML output = %#v, want %#v", actual, expectedVersionOutput())
	}
}

func TestDisplayVersionOutputTableFallback(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	stdout := captureStdout(t, func() {
		err := displayVersionOutput(
			&cobra.Command{},
			"",
			"1.2.3",
			"client-sha",
			config.Version{Build: "4.5.6", CommitSHA: "server-sha"},
			[]string{"", "Version", "GitSHA"},
			[][]string{{"Client", "1.2.3", "client-sha"}, {"Server", "4.5.6", "server-sha"}},
		)
		if err != nil {
			t.Fatalf("displayVersionOutput() error = %v", err)
		}
	})

	for _, expected := range []string{"Client", "Server", "1.2.3", "4.5.6", "client-sha", "server-sha"} {
		if !strings.Contains(stdout, expected) {
			t.Errorf("table output does not contain %q: %q", expected, stdout)
		}
	}
}

func TestVersionCmdRejectsInvalidOutputFormat(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	original := versionOutputFormat
	versionOutputFormat = "invalid"
	t.Cleanup(func() { versionOutputFormat = original })

	err := versionCmd.PreRunE(versionCmd, nil)
	if err == nil {
		t.Fatal("version command accepted invalid output format")
	}
	if !strings.Contains(err.Error(), display.ErrInvalidOutputFormat("invalid").Error()) {
		t.Errorf("error = %v, want invalid output format error", err)
	}
}

func expectedVersionOutput() map[string]map[string]string {

	return map[string]map[string]string{
		"client": {"version": "1.2.3", "git_sha": "client-sha"},
		"server": {"version": "4.5.6", "git_sha": "server-sha"},
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = writer
	t.Cleanup(func() { os.Stdout = original })

	fn()
	if err := writer.Close(); err != nil {
		t.Fatalf("stdout writer close error = %v", err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("stdout read error = %v", err)
	}
	return string(output)
}
