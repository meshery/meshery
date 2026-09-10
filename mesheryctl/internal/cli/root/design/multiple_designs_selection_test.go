package design

import (
	"os"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
)

// withStdin points os.Stdin at the given text for the duration of fn. Passing
// an empty string gives an immediate EOF, which is what a script, a pipe or a
// CI job looks like from inside the process.
func withStdin(t *testing.T, text string, fn func()) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	original := os.Stdin
	os.Stdin = r
	defer func() {
		os.Stdin = original
		_ = r.Close()
	}()

	if text != "" {
		if _, err := w.WriteString(text); err != nil {
			t.Fatalf("write to pipe: %v", err)
		}
	}
	// Closing the writer is what produces EOF once the text is consumed.
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}

	fn()
}

func designs(n int) []models.MesheryPattern {
	out := make([]models.MesheryPattern, 0, n)
	for i := 0; i < n; i++ {
		id, _ := uuid.NewV4()
		out = append(out, models.MesheryPattern{
			ID:          &id,
			Name:        "My Design",
			PatternFile: "name: My Design\n",
		})
	}
	return out
}

// A name matching several designs used to select index 0 and carry on when
// there was nobody to answer the prompt. Both commands share the behaviour, so
// both are covered here.
func TestMultipleDesignsSelectionIsNotInteractive(t *testing.T) {
	cases := []struct {
		name    string
		matches int
		call    func([]models.MesheryPattern) (int, error)
	}{
		{"deploy", 2, func(p []models.MesheryPattern) (int, error) {
			return multiplepatternsConfirmation(p, "My Design")
		}},
		{"deploy, five matches", 5, func(p []models.MesheryPattern) (int, error) {
			return multiplepatternsConfirmation(p, "My Design")
		}},
		{"apply", 2, func(p []models.MesheryPattern) (int, error) {
			return multiplePatternsConfirmation(p, "My Design")
		}},
		{"apply, five matches", 5, func(p []models.MesheryPattern) (int, error) {
			return multiplePatternsConfirmation(p, "My Design")
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_ = utils.SetupMeshkitLoggerTesting(t, false)

			var (
				index int
				err   error
			)

			withStdin(t, "", func() {
				index, err = tc.call(designs(tc.matches))
			})

			if err == nil {
				t.Fatalf("expected an error with no interactive stdin, got index %d and nil", index)
			}
			if index != 0 {
				t.Errorf("index should be the zero value alongside an error, got %d", index)
			}
			if code := errors.GetCode(err); code != ErrDesignSelectNotInteractiveCode {
				t.Errorf("error code = %s, want %s", code, ErrDesignSelectNotInteractiveCode)
			}
			// The message has to name the design, or the operator cannot tell
			// which invocation failed from a CI log.
			if !strings.Contains(err.Error(), "My Design") {
				t.Errorf("error should name the design, got: %v", err)
			}
		})
	}
}

// The interactive path must be unchanged.
func TestMultipleDesignsSelectionAcceptsAnAnswer(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stdin string
		want  int
	}{
		{"first", "0\n", 0},
		{"middle", "2\n", 2},
		{"last", "4\n", 4},
		{"surrounding whitespace", "  3  \n", 3},
		{"retries past an unparseable answer", "banana\n1\n", 1},
		{"retries past an out-of-range answer", "99\n1\n", 1},
		{"retries past a negative answer", "-1\n2\n", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_ = utils.SetupMeshkitLoggerTesting(t, false)

			var (
				index int
				err   error
			)

			withStdin(t, tc.stdin, func() {
				index, err = multiplepatternsConfirmation(designs(5), "My Design")
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if index != tc.want {
				t.Errorf("index = %d, want %d", index, tc.want)
			}
		})
	}
}

// An answer that never becomes valid still terminates, because the pipe closes.
func TestMultipleDesignsSelectionGivesUpWhenInputRunsOut(t *testing.T) {
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	var err error

	withStdin(t, "banana\n", func() {
		_, err = multiplepatternsConfirmation(designs(3), "My Design")
	})

	if err == nil {
		t.Fatal("expected an error once stdin is exhausted, got nil")
	}
}
