package filter

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/core"
)

func TestListCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenrios for fetching data
	tests := []utils.MesheryCommandTest{
		{
			Name:             "Fetch Filter List",
			Args:             []string{"list"},
			ExpectedResponse: "list.filter.output.golden",
			Fixture:          "list.filter.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
			HttpMethod:       "GET",
			HttpStatusCode:   200,
		},
		{
			Name:             "Fetch Filter List with pagesize",
			Args:             []string{"list", "--page", "1", "--pagesize", "10"},
			ExpectedResponse: "list.filter.output.golden",
			Fixture:          "list.filter.api.response.golden",
			URL:              "/api/filter?pagesize=10&page=0",
			ExpectError:      false,
			HttpMethod:       "GET",
			HttpStatusCode:   200,
		},
		{
			Name:             "List filters with invalid page number",
			Args:             []string{"list", "--page", "0"},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '0'")),
		},
		{
			Name:             "List filters with invalid page size",
			Args:             []string{"list", "--pagesize", "0"},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '0'")),
		},
		{
			Name:             "List filters with too many arguments",
			Args:             []string{"list", "filter1", "filter2"},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg)),
		},
	}
	mesheryctlflags.InitValidators(FilterCmd)
	// Run tests
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, tests, currDir, "filter")
}

func TestListCmdCount(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Fetch Filter Count",
			Args:             []string{"list", "--count"},
			ExpectedResponse: "list.filter.count.output.golden",
			Fixture:          "filter.list.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
		},
		{
			Name:             "Fetch Filter Count with no matching filters",
			Args:             []string{"list", "nonexistent", "--count"},
			ExpectedResponse: "list.filter.count.empty.output.golden",
			Fixture:          "filter.list.empty.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(FilterCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, FilterCmd, tests, currDir, "filter")
}

func TestProcessFilterData(t *testing.T) {
	ts := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	filterID := core.Uuid{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
	longOwner := "user1234-5678-90ab-cdef-1234567890ab"
	shortOwner := "usr"

	header, data, footer := processFilterData(nil, false, true)
	if header != nil || data != nil || footer != nil {
		t.Fatalf("processFilterData(nil) = (%v, %v, %v), want (nil, nil, nil)", header, data, footer)
	}

	tests := []struct {
		name       string
		isVerbose  bool
		isLocal    bool
		filter     models.MesheryFilter
		wantHeader []string
		wantRow    []string
		wantFooter []string
	}{
		{
			name:       "local non-verbose handles nil id and timestamps",
			isVerbose:  false,
			isLocal:    true,
			filter:     models.MesheryFilter{Name: "my-filter.wasm"},
			wantHeader: []string{"FILTER ID", "NAME", "CREATED", "UPDATED"},
			wantRow:    []string{"N/A", "my-filter", "N/A", "N/A"},
			wantFooter: []string{"Total", "1", "", ""},
		},
		{
			name:       "local verbose handles nil id and timestamps",
			isVerbose:  true,
			isLocal:    true,
			filter:     models.MesheryFilter{Name: "my-filter.wasm"},
			wantHeader: []string{"FILTER ID", "NAME", "CREATED", "UPDATED"},
			wantRow:    []string{"N/A", "my-filter", "N/A", "N/A"},
			wantFooter: []string{"Total", "1", "", ""},
		},
		{
			name:       "remote non-verbose handles nil id, owner, and timestamps",
			isVerbose:  false,
			isLocal:    false,
			filter:     models.MesheryFilter{Name: "remote-filter.wasm"},
			wantHeader: []string{"FILTER ID", "OWNER", "NAME", "CREATED", "UPDATED"},
			wantRow:    []string{"N/A", "null", "remote-filter", "N/A", "N/A"},
			wantFooter: []string{"Total", "1", "", "", ""},
		},
		{
			name:       "remote non-verbose handles short owner without slice panic",
			isVerbose:  false,
			isLocal:    false,
			filter:     models.MesheryFilter{ID: &filterID, Owner: &shortOwner, Name: "short-owner.wasm", CreatedAt: &ts, UpdatedAt: &ts},
			wantHeader: []string{"FILTER ID", "OWNER", "NAME", "CREATED", "UPDATED"},
			wantRow:    []string{"12345678", "usr", "short-owner", "1-2-2025", "1-2-2025"},
			wantFooter: []string{"Total", "1", "", "", ""},
		},
		{
			name:       "remote verbose formats full owner and timestamp",
			isVerbose:  true,
			isLocal:    false,
			filter:     models.MesheryFilter{ID: &filterID, Owner: &longOwner, Name: "full-filter.wasm", CreatedAt: &ts, UpdatedAt: &ts},
			wantHeader: []string{"FILTER ID", "OWNER", "NAME", "CREATED", "UPDATED"},
			wantRow:    []string{"12345678", longOwner, "full-filter", "1-2-2025 3:4:5", "1-2-2025 3:4:5"},
			wantFooter: []string{"Total", "1", "", "", ""},
		},
		{
			name:       "local verbose keeps full id and timestamp",
			isVerbose:  true,
			isLocal:    true,
			filter:     models.MesheryFilter{ID: &filterID, Name: "full-local.wasm", CreatedAt: &ts, UpdatedAt: &ts},
			wantHeader: []string{"FILTER ID", "NAME", "CREATED", "UPDATED"},
			wantRow:    []string{"12345678-90ab-cdef-1234-567890abcdef", "full-local", "1-2-2025 3:4:5", "1-2-2025 3:4:5"},
			wantFooter: []string{"Total", "1", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHeader, gotData, gotFooter := processFilterData(&models.FiltersAPIResponse{
				Filters:    []models.MesheryFilter{tt.filter},
				TotalCount: 1,
			}, tt.isVerbose, tt.isLocal)
			if !reflect.DeepEqual(gotHeader, tt.wantHeader) {
				t.Errorf("gotHeader = %#v, want %#v", gotHeader, tt.wantHeader)
			}
			if !reflect.DeepEqual(gotData, [][]string{tt.wantRow}) {
				t.Errorf("gotData = %#v, want %#v", gotData, [][]string{tt.wantRow})
			}
			if !reflect.DeepEqual(gotFooter, tt.wantFooter) {
				t.Errorf("gotFooter = %#v, want %#v", gotFooter, tt.wantFooter)
			}
		})
	}
}
