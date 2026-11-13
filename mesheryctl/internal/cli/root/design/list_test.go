package design

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

func TestDesignListCmd(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no arguments provided when design list then design list is displayed",
			Args:             []string{"list", "--page", "1"},
			ExpectedResponse: "list.design.output.golden",
			Fixture:          "list.design.api.response.golden",
			URL:              "/api/pattern?page=0&pagesize=10",
			ExpectError:      false,
		},
		{
			Name:             "given extra positional arguments when design list then throw error",
			Args:             []string{"list", "extra-arg"},
			ExpectedResponse: "",
			Fixture:          "list.design.api.response.golden",
			URL:              "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg)),
		},
		{
			Name:             "given --page 0 when design list then throw invalid flag error",
			Args:             []string{"list", "--page", "0"},
			ExpectedResponse: "",
			URL:              "",
			Fixture:          "list.design.api.response.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '0'")),
		},
		{
			Name:             "given a negative --page when design list then throw invalid flag error",
			Args:             []string{"list", "--page", "-1"},
			ExpectedResponse: "",
			URL:              "",
			Fixture:          "list.design.api.response.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '-1'")),
		},
		{
			Name:             "given --pagesize 0 when design list then throw invalid flag error",
			Args:             []string{"list", "--pagesize", "0"},
			ExpectedResponse: "",
			URL:              "",
			Fixture:          "list.design.api.response.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '0'")),
		},
		{
			Name:             "given a negative --pagesize when design list then throw invalid flag error",
			Args:             []string{"list", "--pagesize", "-1"},
			ExpectedResponse: "",
			URL:              "",
			Fixture:          "list.design.api.response.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '-1'")),
		},
	}

	mesheryctlflags.InitValidators(DesignCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, DesignCmd, tests, currDir, "design")
}

func TestProcessDesignDataNilFields(t *testing.T) {
	origProvider := provider
	origVerbose := verbose
	defer func() {
		provider = origProvider
		verbose = origVerbose
	}()

	ts := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	designID := core.Uuid{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
	userID := core.Uuid{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99}

	rows, count := processDesignData(nil)
	if rows != nil || count != 0 {
		t.Fatalf("processDesignData(nil) = (%v, %d), want (nil, 0)", rows, count)
	}

	tests := []struct {
		name      string
		provider  string
		isVerbose bool
		pattern   models.MesheryPattern
		wantRow   []string
	}{
		{
			name:      "remote provider handles nil id, userId, and timestamps",
			provider:  "Meshery",
			isVerbose: false,
			pattern:   models.MesheryPattern{Name: "my-design.yaml"},
			wantRow:   []string{"N/A", "null", "my-design", "N/A", "N/A"},
		},
		{
			name:      "remote provider verbose handles nil id, userId, and timestamps",
			provider:  "Meshery",
			isVerbose: true,
			pattern:   models.MesheryPattern{Name: "my-design.yaml"},
			wantRow:   []string{"N/A", "null", "my-design", "N/A", "N/A"},
		},
		{
			name:      "local provider handles nil timestamps with valid id",
			provider:  "None",
			isVerbose: false,
			pattern:   models.MesheryPattern{ID: &designID, Name: "local-design.yml"},
			wantRow:   []string{"12345678-90ab-cdef-1234-567890abcdef", "local-design", "N/A", "N/A"},
		},
		{
			name:      "remote provider verbose formats non-nil fields",
			provider:  "Meshery",
			isVerbose: true,
			pattern:   models.MesheryPattern{ID: &designID, UserID: &userID, Name: "full-design.yaml", CreatedAt: &ts, UpdatedAt: &ts},
			wantRow:   []string{"12345678-90ab-cdef-1234-567890abcdef", "aabbccdd-eeff-0011-2233-445566778899", "full-design", "01-02-2025 03:04:05", "01-02-2025 03:04:05"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider = tt.provider
			verbose = tt.isVerbose
			gotRows, gotCount := processDesignData(&models.PatternsAPIResponse{
				Patterns:   []models.MesheryPattern{tt.pattern},
				TotalCount: 1,
			})
			if gotCount != 1 {
				t.Fatalf("gotCount = %d, want 1", gotCount)
			}
			want := [][]string{tt.wantRow}
			if !reflect.DeepEqual(gotRows, want) {
				t.Errorf("gotRows = %#v, want %#v", gotRows, want)
			}
		})
	}
}
