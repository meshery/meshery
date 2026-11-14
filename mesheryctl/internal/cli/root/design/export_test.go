package design

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

// TestErrDoRequestNilError ensures ErrDoRequest does not panic when err is nil.
func TestErrDoRequestNilError(t *testing.T) {
	err := models.ErrDoRequest(nil, "GET", "http://example.com/api/test")
	if err == nil {
		t.Fatal("expected non-nil error from ErrDoRequest, got nil")
	}
	if meshkiterrors.GetCode(err) != models.ErrDoRequestCode {
		t.Fatalf("expected code %s, got %s", models.ErrDoRequestCode, meshkiterrors.GetCode(err))
	}
}

// TestFetchPatternDataNonOKStatus ensures fetchPatternData handles non-200 responses
// gracefully without panicking.
func TestFetchPatternDataNonOKStatus(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(DesignCmd, t)
	utils.TokenFlag = utils.GetToken(t)

	t.Run("404 Not Found", func(t *testing.T) {
		targetURL := testContext.BaseURL + "/api/pattern/test-id-404"
		httpmock.RegisterResponder("GET", targetURL, func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, "pattern not found"), nil
		})

		pattern, err := fetchPatternData(targetURL)
		if err == nil {
			t.Fatal("expected error on 404 response, got nil")
		}
		if pattern != nil {
			t.Fatalf("expected nil pattern on error, got %+v", pattern)
		}
		if meshkiterrors.GetCode(err) != utils.ErrNotFoundCode {
			t.Fatalf("expected ErrNotFoundCode (%s), got %s", utils.ErrNotFoundCode, meshkiterrors.GetCode(err))
		}
	})

	t.Run("201 Non-OK Status without MakeRequest error", func(t *testing.T) {
		targetURL := testContext.BaseURL + "/api/pattern/test-id-201"
		httpmock.RegisterResponder("GET", targetURL, func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, "created"), nil
		})

		pattern, err := fetchPatternData(targetURL)
		if err == nil {
			t.Fatal("expected error on 201 response, got nil")
		}
		if pattern != nil {
			t.Fatalf("expected nil pattern on error, got %+v", pattern)
		}
		if meshkiterrors.GetCode(err) != utils.ErrResponseStatusCode {
			t.Fatalf("expected ErrResponseStatusCode (%s), got %s", utils.ErrResponseStatusCode, meshkiterrors.GetCode(err))
		}
	})
}

// TestFetchPatternIDByNameNonOKStatus ensures fetchPatternIDByName handles non-200 responses
// gracefully without panicking.
func TestFetchPatternIDByNameNonOKStatus(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(DesignCmd, t)
	utils.TokenFlag = utils.GetToken(t)

	t.Run("500 Internal Server Error", func(t *testing.T) {
		patternName := "test-pattern-500"
		searchURL := fmt.Sprintf("%s/api/pattern?search=%s", testContext.BaseURL, url.QueryEscape(patternName))
		httpmock.RegisterResponder("GET", searchURL, func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusInternalServerError, "internal server error"), nil
		})

		id, err := fetchPatternIDByName(testContext.BaseURL, patternName)
		if err == nil {
			t.Fatal("expected error on 500 response, got nil")
		}
		if id != "" {
			t.Fatalf("expected empty id on error, got %s", id)
		}
		if meshkiterrors.GetCode(err) != utils.ErrMesheryServerInternalErrorCode {
			t.Fatalf("expected ErrMesheryServerInternalErrorCode (%s), got %s", utils.ErrMesheryServerInternalErrorCode, meshkiterrors.GetCode(err))
		}
	})

	t.Run("201 Non-OK Status without MakeRequest error", func(t *testing.T) {
		patternName := "test-pattern-201"
		searchURL := fmt.Sprintf("%s/api/pattern?search=%s", testContext.BaseURL, url.QueryEscape(patternName))
		httpmock.RegisterResponder("GET", searchURL, func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, "created"), nil
		})

		id, err := fetchPatternIDByName(testContext.BaseURL, patternName)
		if err == nil {
			t.Fatal("expected error on 201 response, got nil")
		}
		if id != "" {
			t.Fatalf("expected empty id on error, got %s", id)
		}
		if meshkiterrors.GetCode(err) != utils.ErrResponseStatusCode {
			t.Fatalf("expected ErrResponseStatusCode (%s), got %s", utils.ErrResponseStatusCode, meshkiterrors.GetCode(err))
		}
	})
}

// TestGetOwnerNameNonOKStatus ensures getOwnerName handles non-200 responses
// gracefully without panicking.
func TestGetOwnerNameNonOKStatus(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(DesignCmd, t)
	utils.TokenFlag = utils.GetToken(t)

	t.Run("401 Unauthorized", func(t *testing.T) {
		ownerID := "owner-401"
		profileURL := fmt.Sprintf("%s/api/user/profile/%s", testContext.BaseURL, ownerID)
		httpmock.RegisterResponder("GET", profileURL, func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusUnauthorized, "unauthorized"), nil
		})

		owner, err := getOwnerName(ownerID, testContext.BaseURL)
		if err == nil {
			t.Fatal("expected error on 401 response, got nil")
		}
		if owner != "" {
			t.Fatalf("expected empty owner name on error, got %s", owner)
		}
		if meshkiterrors.GetCode(err) != utils.ErrFailReqStatusCode {
			t.Fatalf("expected ErrFailReqStatusCode (%s), got %s", utils.ErrFailReqStatusCode, meshkiterrors.GetCode(err))
		}
	})

	t.Run("201 Non-OK Status without MakeRequest error", func(t *testing.T) {
		ownerID := "owner-201"
		profileURL := fmt.Sprintf("%s/api/user/profile/%s", testContext.BaseURL, ownerID)
		httpmock.RegisterResponder("GET", profileURL, func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, "created"), nil
		})

		owner, err := getOwnerName(ownerID, testContext.BaseURL)
		if err == nil {
			t.Fatal("expected error on 201 response, got nil")
		}
		if owner != "" {
			t.Fatalf("expected empty owner name on error, got %s", owner)
		}
		if meshkiterrors.GetCode(err) != utils.ErrResponseStatusCode {
			t.Fatalf("expected ErrResponseStatusCode (%s), got %s", utils.ErrResponseStatusCode, meshkiterrors.GetCode(err))
		}
	})
}

// TestFetchPatternDataSuccess ensures fetchPatternData decodes valid response.
func TestFetchPatternDataSuccess(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(DesignCmd, t)
	utils.TokenFlag = utils.GetToken(t)

	expectedPattern := models.MesheryPattern{
		Name: "my-design",
	}

	targetURL := testContext.BaseURL + "/api/pattern/my-design-id"
	httpmock.RegisterResponder("GET", targetURL, func(req *http.Request) (*http.Response, error) {
		respBytes, _ := json.Marshal(expectedPattern)
		return httpmock.NewBytesResponse(http.StatusOK, respBytes), nil
	})

	pattern, err := fetchPatternData(targetURL)
	if err != nil {
		t.Fatalf("unexpected error from fetchPatternData: %v", err)
	}
	if pattern == nil || pattern.Name != "my-design" {
		t.Fatalf("expected pattern name 'my-design', got %+v", pattern)
	}
}

// TestExportDesignNonOKDownloadStatus ensures exportDesign returns a status-based error
// when the download endpoint responds with a non-200 status that MakeRequest accepts.
func TestExportDesignNonOKDownloadStatus(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(DesignCmd, t)
	utils.TokenFlag = utils.GetToken(t)

	design := "test-design-201"
	patternBytes, _ := json.Marshal(models.MesheryPattern{Name: design})
	patternURL := fmt.Sprintf("%s/api/pattern/%s", testContext.BaseURL, design)
	httpmock.RegisterResponder("GET", patternURL, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewBytesResponse(http.StatusOK, patternBytes), nil
	})

	downloadURL := fmt.Sprintf("%s/api/pattern/download/%s", testContext.BaseURL, design)
	httpmock.RegisterResponder("GET", downloadURL, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusCreated, "created"), nil
	})

	err := exportDesign(testContext.BaseURL, design, "")
	if err == nil {
		t.Fatal("expected error on non-200 download response, got nil")
	}
	if meshkiterrors.GetCode(err) != utils.ErrResponseStatusCode {
		t.Fatalf("expected ErrResponseStatusCode (%s), got %s", utils.ErrResponseStatusCode, meshkiterrors.GetCode(err))
	}
}

// TestFetchPatternIDByName verifies the lookup against the response the server
// sends, built from its own type so the test can't drift from the wire format.
func TestFetchPatternIDByName(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	utils.TokenFlag = utils.GetToken(t)

	patternID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name     string
		patterns []models.MesheryPattern
		wantID   string
	}{
		{name: "one match returns its ID", patterns: []models.MesheryPattern{{ID: &patternID, Name: "my-design"}}, wantID: patternID.String()},
		{name: "no match returns an error", patterns: []models.MesheryPattern{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(models.PatternsAPIResponse{TotalCount: uint(len(tt.patterns)), Patterns: tt.patterns})
			if err != nil {
				t.Fatal(err)
			}
			searchURL := fmt.Sprintf("%s/api/pattern?search=%s", testContext.BaseURL, url.QueryEscape("my-design"))
			httpmock.RegisterResponder("GET", searchURL, httpmock.NewStringResponder(http.StatusOK, string(body)))

			id, err := fetchPatternIDByName(testContext.BaseURL, "my-design")
			if tt.wantID == "" {
				wantCode := meshkiterrors.GetCode(ErrDesignNotFound("my-design"))
				if meshkiterrors.GetCode(err) != wantCode {
					t.Fatalf("expected error code %s, got %v", wantCode, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if id != tt.wantID {
				t.Fatalf("expected id %s, got %s", tt.wantID, id)
			}
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	const id = "a12b3c4d-5e6f-4890-abcd-ef1234567890"
	tests := []struct {
		designType, sourceType, want string
	}{
		{"current", "", "my-design_ef1234567890.yaml"},
		{"oci", "", "my-design_ef1234567890.tar"},
		{"original", string(models.HelmChart), "my-design_ef1234567890.tar.gz"},
		{"original", string(models.K8sManifest), "my-design_ef1234567890.yaml"},
		{"original", string(models.DockerCompose), "my-design_ef1234567890.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.designType+"/"+tt.sourceType, func(t *testing.T) {
			if got := generateFilename("my-design", id, tt.designType, tt.sourceType); got != tt.want {
				t.Fatalf("generateFilename(%q, %q) = %q, want %q", tt.designType, tt.sourceType, got, tt.want)
			}
		})
	}
}

// TestExportDesignFileExtension runs exportDesign end to end and checks the
// name of the file it writes for each export type.
func TestExportDesignFileExtension(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	utils.TokenFlag = utils.GetToken(t)

	const design = "a12b3c4d-5e6f-4890-abcd-ef1234567890"
	tests := []struct {
		name, designType, sourceType, downloadPath, wantFile string
	}{
		{"oci image is a tar", "oci", string(models.Design), "?oci=true", "my-design_ef1234567890.tar"},
		{"original kubernetes manifest is yaml", "original", string(models.K8sManifest), "/Kubernetes%20Manifest", "my-design_ef1234567890.yaml"},
		{"original helm chart is a gzipped tar", "original", string(models.HelmChart), "/Helm%20Chart", "my-design_ef1234567890.tar.gz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := models.MesheryPattern{Name: "my-design"}
			pattern.Type.String, pattern.Type.Valid = tt.sourceType, true
			patternBytes, err := json.Marshal(pattern)
			if err != nil {
				t.Fatal(err)
			}
			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/pattern/%s", testContext.BaseURL, design),
				httpmock.NewBytesResponder(http.StatusOK, patternBytes))
			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/pattern/download/%s%s", testContext.BaseURL, design, tt.downloadPath),
				httpmock.NewStringResponder(http.StatusOK, "content"))

			prevOutputDir := outputDir
			outputDir = t.TempDir()
			defer func() { outputDir = prevOutputDir }()

			if err := exportDesign(testContext.BaseURL, design, tt.designType); err != nil {
				t.Fatalf("exportDesign: %v", err)
			}
			if _, err := os.Stat(filepath.Join(outputDir, tt.wantFile)); err != nil {
				entries, _ := os.ReadDir(outputDir)
				var got []string
				for _, e := range entries {
					got = append(got, e.Name())
				}
				t.Fatalf("expected %s to be written, directory has %v", tt.wantFile, got)
			}
		})
	}
}

// TestExportCmdRejectsUnknownType checks that an unknown --type fails before
// any request is made, instead of silently exporting the current design.
func TestExportCmdRejectsUnknownType(t *testing.T) {
	_ = utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(DesignCmd, t)
	utils.TokenFlag = utils.GetToken(t)

	DesignCmd.SetArgs([]string{"export", "my-design", "--type", "helm"})
	err := DesignCmd.Execute()
	if err == nil {
		t.Fatal("expected an error for --type helm, got nil")
	}
	if meshkiterrors.GetCode(err) != utils.ErrInvalidArgumentCode {
		t.Fatalf("expected ErrInvalidArgument (%s), got %v", utils.ErrInvalidArgumentCode, err)
	}
	if n := httpmock.GetTotalCallCount(); n != 0 {
		t.Fatalf("expected no requests before rejecting --type, got %d", n)
	}
}
