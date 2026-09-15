package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

// newMeshkitError mirrors the shape server/models.ErrPost produces, without
// importing server/models (which would be an import cycle).
func newMeshkitError(code string) error {
	return meshkiterrors.New(
		code,
		meshkiterrors.Alert,
		[]string{"Unable to post data to the Provider.Environment"},
		[]string{"Status Code: 403 ", "failed to save the environment"},
		[]string{},
		[]string{},
	)
}

// TestWithProviderStatus_RoundTrips is the core of the misreporting fix: the
// status a remote provider responded with must survive as a readable value,
// not only as prose inside the MeshKit cause.
func TestWithProviderStatus_RoundTrips(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		wantStatus int
		wantOK     bool
	}{
		{name: "forbidden is propagated", statusCode: http.StatusForbidden, wantStatus: http.StatusForbidden, wantOK: true},
		{name: "conflict is propagated", statusCode: http.StatusConflict, wantStatus: http.StatusConflict, wantOK: true},
		{name: "provider 500 is propagated", statusCode: http.StatusInternalServerError, wantStatus: http.StatusInternalServerError, wantOK: true},
		{name: "not found is propagated", statusCode: http.StatusNotFound, wantStatus: http.StatusNotFound, wantOK: true},
		{name: "success status is not a failure and is dropped", statusCode: http.StatusOK, wantOK: false},
		{name: "zero status is dropped", statusCode: 0, wantOK: false},
		{name: "out-of-range status is dropped", statusCode: 700, wantOK: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tagged := WithProviderStatus(newMeshkitError("meshery-server-1271"), tc.statusCode)

			gotStatus, gotOK := ProviderStatusCode(tagged)
			if gotOK != tc.wantOK {
				t.Fatalf("ProviderStatusCode ok = %v, want %v", gotOK, tc.wantOK)
			}
			if tc.wantOK && gotStatus != tc.wantStatus {
				t.Fatalf("ProviderStatusCode = %d, want %d", gotStatus, tc.wantStatus)
			}
		})
	}
}

// TestWithProviderStatus_PreservesMeshkitMetadata guards the reason errors.Join
// was chosen over a replacement type: every MeshKit accessor must still reach
// through to the original error, or tagging the status would cost us the code.
func TestWithProviderStatus_PreservesMeshkitMetadata(t *testing.T) {
	const code = "meshery-server-1271"
	tagged := WithProviderStatus(newMeshkitError(code), http.StatusForbidden)

	if got := meshkiterrors.GetCode(tagged); got != code {
		t.Errorf("GetCode = %q, want %q - tagging the status must not hide the MeshKit code", got, code)
	}
	if got := meshkiterrors.GetSeverity(tagged); got != meshkiterrors.Alert {
		t.Errorf("GetSeverity = %v, want Alert", got)
	}

	var meshkitErr *meshkiterrors.Error
	if !errors.As(tagged, &meshkitErr) {
		t.Error("errors.As failed to reach the wrapped MeshKit error")
	}
}

// TestWithProviderStatus_NilPassesThrough confirms the helper never manufactures
// an error where none existed.
func TestWithProviderStatus_NilPassesThrough(t *testing.T) {
	if got := WithProviderStatus(nil, http.StatusForbidden); got != nil {
		t.Fatalf("WithProviderStatus(nil) = %v, want nil", got)
	}
}

// TestStatusForProviderError_FallsBackWhenUnknown covers the failures that never
// produced an upstream HTTP status (unreachable provider, unsupported
// capability, marshalling error): the caller's fallback wins, and it is never
// silently turned into a 404.
func TestStatusForProviderError_FallsBackWhenUnknown(t *testing.T) {
	untagged := fmt.Errorf("could not reach remote provider: %w", errors.New("dial tcp: connection refused"))

	if got := StatusForProviderError(untagged, http.StatusBadGateway); got != http.StatusBadGateway {
		t.Fatalf("StatusForProviderError = %d, want %d", got, http.StatusBadGateway)
	}

	tagged := WithProviderStatus(newMeshkitError("meshery-server-1271"), http.StatusForbidden)
	if got := StatusForProviderError(tagged, http.StatusBadGateway); got != http.StatusForbidden {
		t.Fatalf("StatusForProviderError = %d, want %d (the provider's real status)", got, http.StatusForbidden)
	}
}

// TestWriteMeshkitError_StatusTaggedErrorStillSerializes ties the two halves
// together: a status-tagged error written to the wire must still carry its
// MeshKit code, because the UI keys error rendering off that code.
func TestWriteMeshkitError_StatusTaggedErrorStillSerializes(t *testing.T) {
	const code = "meshery-server-1271"
	tagged := WithProviderStatus(newMeshkitError(code), http.StatusForbidden)

	rec := httptest.NewRecorder()
	WriteMeshkitError(rec, tagged, StatusForProviderError(tagged, http.StatusBadGateway))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}

	var decoded struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Fatalf("body did not parse as JSON: %v", err)
	}
	if decoded.Code != code {
		t.Errorf("code = %q, want %q", decoded.Code, code)
	}
	if decoded.Error == "" {
		t.Error("expected a non-empty error message")
	}
}
