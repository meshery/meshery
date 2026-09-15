// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package printer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

type sampleData struct {
	Name  string `json:"name" yaml:"name"`
	Count int    `json:"count" yaml:"count"`
}

func TestAddFormatFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	var format string

	AddFormatFlag(cmd, &format)

	flag := cmd.Flags().Lookup("output-format")
	if flag == nil {
		t.Fatal("expected 'output-format' flag to be added")
	}
	if flag.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got '%s'", flag.Shorthand)
	}
}

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		writer      io.Writer
		expectErr   bool
		expectedFmt string
	}{
		{
			name:      "nil writer returns error",
			format:    "json",
			writer:    nil,
			expectErr: true,
		},
		{
			name:        "valid format json",
			format:      "json",
			writer:      &bytes.Buffer{},
			expectErr:   false,
			expectedFmt: "json",
		},
		{
			name:        "valid format yaml uppercase trimmed",
			format:      "  YAML ",
			writer:      &bytes.Buffer{},
			expectErr:   false,
			expectedFmt: "yaml",
		},
		{
			name:        "empty format defaults to empty/table",
			format:      "",
			writer:      &bytes.Buffer{},
			expectErr:   false,
			expectedFmt: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(tc.format, tc.writer)
			if tc.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p != nil && p.format != tc.expectedFmt {
				t.Errorf("expected format %q, got %q", tc.expectedFmt, p.format)
			}
		})
	}
}

func TestPrinter_Print(t *testing.T) {
	data := sampleData{Name: "test-model", Count: 42}

	t.Run("renders JSON with indentation", func(t *testing.T) {
		var buf bytes.Buffer
		p, _ := New("json", &buf)

		if err := p.Print(data, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "{\n  \"name\": \"test-model\",\n  \"count\": 42\n}\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("renders YAML with indent", func(t *testing.T) {
		var buf bytes.Buffer
		p, _ := New("yaml", &buf)

		if err := p.Print(data, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "name: test-model\ncount: 42\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("delegates to humanRender for table/empty format", func(t *testing.T) {
		var buf bytes.Buffer
		p, _ := New("table", &buf)

		called := false
		err := p.Print(data, func(w io.Writer) error {
			called = true
			_, err := fmt.Fprint(w, "rendered table")
			return err
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Errorf("expected humanRender to be called")
		}
		if buf.String() != "rendered table" {
			t.Errorf("got %q, want %q", buf.String(), "rendered table")
		}
	})

	t.Run("returns error on invalid format", func(t *testing.T) {
		var buf bytes.Buffer
		p, _ := New("unsupported_format", &buf)

		err := p.Print(data, nil)
		if err == nil {
			t.Fatalf("expected error for invalid format, got nil")
		}
		if !strings.Contains(err.Error(), "invalid value for --output-format") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("passes through humanRender error", func(t *testing.T) {
		var buf bytes.Buffer
		p, _ := New("", &buf)

		customErr := errors.New("render failed")
		err := p.Print(data, func(w io.Writer) error {
			return customErr
		})

		if !errors.Is(err, customErr) {
			t.Errorf("got error %v, want %v", err, customErr)
		}
	})
}
