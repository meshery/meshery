package design

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

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
