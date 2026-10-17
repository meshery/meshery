package display

import (
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
)

type paginationTestItem struct {
	ID string
}

// withTerminal swaps the interactive-terminal probe for the duration of a test,
// and gives utils.Log somewhere to write - the CLI installs it at startup, so
// it is nil in a unit test.
func withTerminal(t *testing.T, interactive bool) {
	t.Helper()

	utils.SetupMeshkitLoggerTesting(t, false)

	original := utils.IsInteractiveTerminal
	utils.IsInteractiveTerminal = func() bool { return interactive }
	t.Cleanup(func() { utils.IsInteractiveTerminal = original })
}

// TestListPageHandlerStopsWithoutTerminal pins the behaviour that makes
// mesheryctl usable in scripts, pipelines and CI: with more results than fit on
// one page and no terminal to page on, the handler stops after the page it has
// printed and reports success.
//
// It used to fall through to the keyboard library, which reads the controlling
// terminal directly and fails with "open /dev/tty: no such device or address"
// where there isn't one - so `mesheryctl connection list` printed a correct
// first page and then exited non-zero, which is how it failed in CI for an
// account holding more than one page of connections.
func TestListPageHandlerStopsWithoutTerminal(t *testing.T) {
	withTerminal(t, false)

	handler := listPageHandler(
		DisplayDataAsync{DataType: "connections", Header: []string{"ID"}},
		func(_ *[]paginationTestItem) ([][]string, int64) {
			// One page of rows out of a much larger total, so the handler
			// reaches the "wait for a keypress" decision.
			return [][]string{{"1"}, {"2"}}, 294
		},
	)

	shouldContinue, err := handler(&[]paginationTestItem{}, 0, 2)
	if err != nil {
		t.Fatalf("expected no error without a terminal, got %v", err)
	}
	if shouldContinue {
		t.Error("expected paging to stop without a terminal, got shouldContinue=true")
	}
}

// TestPromptPageHandlerFailsClearlyWithoutTerminal covers the other half: an
// ambiguous match genuinely cannot be resolved without a terminal, because
// picking for the user would act on a resource they did not name. It must fail
// - but with an error that names the ambiguity, not with the selection
// library's "open /dev/tty".
func TestPromptPageHandlerFailsClearlyWithoutTerminal(t *testing.T) {
	withTerminal(t, false)

	var selected paginationTestItem
	handler := promptPageHandler(
		DisplayDataAsync{DataType: "connections"},
		func(rows []paginationTestItem) []string { return []string{rows[0].ID, rows[1].ID} },
		func(_ *[]paginationTestItem) ([]paginationTestItem, int64) {
			return []paginationTestItem{{ID: "1"}, {ID: "2"}}, 93
		},
		&selected,
	)

	shouldContinue, err := handler(&[]paginationTestItem{}, 0, 2)
	if err == nil {
		t.Fatal("expected an error for an ambiguous match without a terminal, got nil")
	}
	if got := errors.GetCode(err); got != ErrAmbiguousSelectionCode {
		t.Errorf("error code = %q, want %q (err: %v)", got, ErrAmbiguousSelectionCode, err)
	}
	if shouldContinue {
		t.Error("expected the handler to stop, got shouldContinue=true")
	}
	if selected.ID != "" {
		t.Errorf("expected nothing to be selected on the user's behalf, got %q", selected.ID)
	}
}

// TestPromptPageHandlerAutoSelectsSingleMatch guards the boundary: a single
// match is unambiguous, so it is still resolved without a terminal.
func TestPromptPageHandlerAutoSelectsSingleMatch(t *testing.T) {
	withTerminal(t, false)

	var selected paginationTestItem
	handler := promptPageHandler(
		DisplayDataAsync{DataType: "connections"},
		func(rows []paginationTestItem) []string { return []string{rows[0].ID} },
		func(_ *[]paginationTestItem) ([]paginationTestItem, int64) {
			return []paginationTestItem{{ID: "only"}}, 1
		},
		&selected,
	)

	if _, err := handler(&[]paginationTestItem{}, 0, 2); err != nil {
		t.Fatalf("expected a single match to resolve without a terminal, got %v", err)
	}
	if selected.ID != "only" {
		t.Errorf("selected = %q, want %q", selected.ID, "only")
	}
}

// TestPromptPageHandlerDoesNotAutoSelectOnePageOfMany is the other side of that
// boundary. A page holding one row is not the same as one match: the caller
// picks the page size, and `--pagesize 1` makes every page look like a single
// match. Auto-selecting there would act on the first of several candidates as
// though the search had been unambiguous - the exact outcome the terminal check
// exists to prevent.
func TestPromptPageHandlerDoesNotAutoSelectOnePageOfMany(t *testing.T) {
	withTerminal(t, false)

	var selected paginationTestItem
	handler := promptPageHandler(
		DisplayDataAsync{DataType: "connections"},
		func(rows []paginationTestItem) []string { return []string{rows[0].ID} },
		func(_ *[]paginationTestItem) ([]paginationTestItem, int64) {
			// One row on this page, 93 matches overall - `connection view
			// minikube` against the CI provider account.
			return []paginationTestItem{{ID: "first-of-93"}}, 93
		},
		&selected,
	)

	if _, err := handler(&[]paginationTestItem{}, 0, 1); err == nil {
		t.Fatal("expected an ambiguous-match error, got nil")
	}
	if selected.ID != "" {
		t.Errorf("expected nothing to be selected, got %q", selected.ID)
	}
}
