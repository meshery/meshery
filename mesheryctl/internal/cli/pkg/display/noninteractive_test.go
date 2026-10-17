package display

import (
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
)

type paginationTestItem struct {
	ID   string
	Name string
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

// TestPromptPageHandlerExactNameMatch verifies that when multiple candidates match
// a search query but exactly one has an exact case-insensitive name match, promptPageHandler
// auto-selects it without failing as ambiguous or prompting interactively.
func TestPromptPageHandlerExactNameMatch(t *testing.T) {
	kubernetesCandidates := []paginationTestItem{
		{ID: "1", Name: "aws-eks-controller"},
		{ID: "2", Name: "aerospike-kubernetes-operator"},
		{ID: "3", Name: "kube-arangodb"},
		{ID: "4", Name: "azure-kubernetes-configuration"},
		{ID: "5", Name: "kubernetes"},
		{ID: "6", Name: "kubernetes-secret-generator"},
		{ID: "7", Name: "kube-ui-server"},
	}

	tests := []struct {
		name         string
		interactive  bool
		searchTerm   string
		rows         []paginationTestItem
		totalCount   int64
		pageSize     int
		extractNames func(any) []string
		expectError  bool
		errCode      string
		expectID     string
	}{
		{
			name:        "unique exact match auto-selects without terminal",
			interactive: false,
			searchTerm:  "kubernetes",
			rows:        kubernetesCandidates,
			totalCount:  7,
			pageSize:    10,
			expectError: false,
			expectID:    "5",
		},
		{
			name:        "unique exact match is case-insensitive",
			interactive: false,
			searchTerm:  "Kubernetes",
			rows:        kubernetesCandidates,
			totalCount:  7,
			pageSize:    10,
			expectError: false,
			expectID:    "5",
		},
		{
			name:        "unique exact match auto-selects with terminal",
			interactive: true,
			searchTerm:  "kubernetes",
			rows:        kubernetesCandidates,
			totalCount:  7,
			pageSize:    10,
			expectError: false,
			expectID:    "5",
		},
		{
			name:        "multiple exact matches fails as ambiguous without terminal",
			interactive: false,
			searchTerm:  "Pod",
			rows: []paginationTestItem{
				{ID: "pod-1", Name: "Pod"},
				{ID: "pod-2", Name: "Pod"},
				{ID: "pod-3", Name: "PodDisruptionBudget"},
			},
			totalCount:  3,
			pageSize:    10,
			expectError: true,
			errCode:     ErrAmbiguousSelectionCode,
		},
		{
			name:        "zero exact matches fails as ambiguous without terminal",
			interactive: false,
			searchTerm:  "kube",
			rows: []paginationTestItem{
				{ID: "1", Name: "kube-arangodb"},
				{ID: "2", Name: "kube-ui-server"},
			},
			totalCount:  2,
			pageSize:    10,
			expectError: true,
			errCode:     ErrAmbiguousSelectionCode,
		},
		{
			name:        "exact match on partial page does not auto-select without terminal",
			interactive: false,
			searchTerm:  "kubernetes",
			rows:        []paginationTestItem{{ID: "5", Name: "kubernetes"}},
			totalCount:  7,
			pageSize:    1,
			expectError: true,
			errCode:     ErrAmbiguousSelectionCode,
		},
		{
			name:        "custom ExtractNames function is respected",
			interactive: false,
			searchTerm:  "custom-alias",
			rows: []paginationTestItem{
				{ID: "alias-1", Name: "something-else"},
				{ID: "alias-2", Name: "another-item"},
			},
			totalCount: 2,
			pageSize:   10,
			extractNames: func(item any) []string {
				if p, ok := item.(paginationTestItem); ok && p.ID == "alias-1" {
					return []string{"custom-alias"}
				}
				return nil
			},
			expectError: false,
			expectID:    "alias-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTerminal(t, tt.interactive)

			var selected paginationTestItem
			handler := promptPageHandler(
				DisplayDataAsync{
					DataType:     "models",
					SearchTerm:   tt.searchTerm,
					ExtractNames: tt.extractNames,
				},
				func(rows []paginationTestItem) []string {
					labels := make([]string, len(rows))
					for i, r := range rows {
						labels[i] = r.Name
					}
					return labels
				},
				func(_ *[]paginationTestItem) ([]paginationTestItem, int64) {
					return tt.rows, tt.totalCount
				},
				&selected,
			)

			shouldContinue, err := handler(&[]paginationTestItem{}, 0, tt.pageSize)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errCode != "" {
					if got := errors.GetCode(err); got != tt.errCode {
						t.Errorf("error code = %q, want %q (err: %v)", got, tt.errCode, err)
					}
				}
				if selected.ID != "" {
					t.Errorf("expected nothing selected, got ID %q", selected.ID)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if shouldContinue {
					t.Errorf("expected shouldContinue=false, got true")
				}
				if selected.ID != tt.expectID {
					t.Errorf("selected ID = %q, want %q", selected.ID, tt.expectID)
				}
			}
		})
	}
}
