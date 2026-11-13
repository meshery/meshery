package models

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models/httputil"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

// TestErrMarshalYAML is a regression test for issue #22681: the short
// description carried a bare "%s" format verb that was never fed through
// Sprintf, so the obj value was dropped and the literal "%s" surfaced to the
// user. The rendered message must now contain the interpolated obj value and
// no literal format verb.
func TestErrMarshalYAML(t *testing.T) {
	underlying := errors.New("yaml: line 3: mapping values are not allowed")
	const obj = "my-design.yaml"

	err := ErrMarshalYAML(underlying, obj)
	if err == nil {
		t.Fatal("ErrMarshalYAML returned nil, expected a wrapped error")
	}

	if meshkiterrors.GetCode(err) != ErrMarshalYAMLCode {
		t.Errorf("expected error code %q, got %q", ErrMarshalYAMLCode, meshkiterrors.GetCode(err))
	}

	sdesc := meshkiterrors.GetSDescription(err)
	if !strings.Contains(sdesc, obj) {
		t.Errorf("short description %q does not contain interpolated obj %q", sdesc, obj)
	}
	if strings.Contains(sdesc, "%s") || strings.Contains(sdesc, "%v") {
		t.Errorf("short description %q still contains a literal format verb", sdesc)
	}

	if msg := err.Error(); !strings.Contains(msg, obj) {
		t.Errorf("Error() %q does not contain interpolated obj %q", msg, obj)
	}
}

func assertProviderErrorNilSafeBehavior(t *testing.T, fn func(error, string, int) error, wantCode, obj string) {
	t.Helper()

	t.Run("nil error with standard HTTP status", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("unexpected panic with nil error: %v", r)
			}
		}()

		err := fn(nil, obj, http.StatusInternalServerError)
		if err == nil {
			t.Fatal("expected non-nil error, got nil")
		}
		if got := meshkiterrors.GetCode(err); got != wantCode {
			t.Errorf("error code = %q, want %q", got, wantCode)
		}
		if gotStatus, ok := httputil.ProviderStatusCode(err); !ok || gotStatus != http.StatusInternalServerError {
			t.Errorf("ProviderStatusCode = (%d, %v), want (%d, true)", gotStatus, ok, http.StatusInternalServerError)
		}
		if sdesc := meshkiterrors.GetSDescription(err); !strings.Contains(sdesc, obj) {
			t.Errorf("short description %q does not contain %q", sdesc, obj)
		}

		mkErr := assertMeshKitErrorStruct(t, err)
		longDesc := strings.Join(mkErr.LongDescription, " ")
		if !strings.Contains(longDesc, "500") || !strings.Contains(longDesc, http.StatusText(http.StatusInternalServerError)) {
			t.Errorf("long description %q does not contain status code and status text", longDesc)
		}
	})

	t.Run("nil error with custom HTTP status", func(t *testing.T) {
		const customStatus = 599
		err := fn(nil, obj, customStatus)
		if err == nil {
			t.Fatal("expected non-nil error, got nil")
		}
		if got := meshkiterrors.GetCode(err); got != wantCode {
			t.Errorf("error code = %q, want %q", got, wantCode)
		}
		if gotStatus, ok := httputil.ProviderStatusCode(err); !ok || gotStatus != customStatus {
			t.Errorf("ProviderStatusCode = (%d, %v), want (%d, true)", gotStatus, ok, customStatus)
		}

		mkErr := assertMeshKitErrorStruct(t, err)
		longDesc := strings.Join(mkErr.LongDescription, " ")
		if !strings.Contains(longDesc, "599") || !strings.Contains(longDesc, "unexpected HTTP status 599") {
			t.Errorf("long description %q does not contain unexpected HTTP status %d fallback", longDesc, customStatus)
		}
	})

	t.Run("non-nil error preserves cause", func(t *testing.T) {
		cause := errors.New("upstream gateway timeout")
		err := fn(cause, obj, http.StatusBadGateway)
		if err == nil {
			t.Fatal("expected non-nil error, got nil")
		}
		if got := meshkiterrors.GetCode(err); got != wantCode {
			t.Errorf("error code = %q, want %q", got, wantCode)
		}
		if gotStatus, ok := httputil.ProviderStatusCode(err); !ok || gotStatus != http.StatusBadGateway {
			t.Errorf("ProviderStatusCode = (%d, %v), want (%d, true)", gotStatus, ok, http.StatusBadGateway)
		}

		mkErr := assertMeshKitErrorStruct(t, err)
		longDesc := strings.Join(mkErr.LongDescription, " ")
		if !strings.Contains(longDesc, cause.Error()) {
			t.Errorf("long description %q does not contain cause %q", longDesc, cause.Error())
		}
	})
}

// TestErrFetch_NilError verifies ErrFetch handles nil underlying errors safely
// without panicking while preserving its MeshKit code and provider status tag.
func TestErrFetch_NilError(t *testing.T) {
	assertProviderErrorNilSafeBehavior(t, ErrFetch, ErrFetchCode, "Users Data")
}

// TestErrPost_NilError verifies ErrPost handles nil underlying errors safely
// without panicking while preserving its MeshKit code and provider status tag.
func TestErrPost_NilError(t *testing.T) {
	assertProviderErrorNilSafeBehavior(t, ErrPost, ErrPostCode, "Perf Result")
}

// TestHTTPStatusFallbackAndUnexpectedError protects the shared status-text
// fallback helper used by providerErrorCause and unexpectedHTTPStatusError.
func TestHTTPStatusFallbackAndUnexpectedError(t *testing.T) {
	if got := httpStatusFallback(http.StatusBadGateway, false); got != "Bad Gateway" {
		t.Errorf("httpStatusFallback(502, false) = %q, want %q", got, "Bad Gateway")
	}
	if got := httpStatusFallback(http.StatusBadGateway, true); got != "unexpected HTTP status 502 Bad Gateway" {
		t.Errorf("httpStatusFallback(502, true) = %q, want %q", got, "unexpected HTTP status 502 Bad Gateway")
	}
	if got := httpStatusFallback(599, false); got != "unexpected HTTP status 599" {
		t.Errorf("httpStatusFallback(599, false) = %q, want %q", got, "unexpected HTTP status 599")
	}
	if got := httpStatusFallback(599, true); got != "unexpected HTTP status 599" {
		t.Errorf("httpStatusFallback(599, true) = %q, want %q", got, "unexpected HTTP status 599")
	}

	if got := providerErrorCause(errors.New("dial timeout"), 502); got != "dial timeout" {
		t.Errorf("providerErrorCause(non-nil) = %q, want %q", got, "dial timeout")
	}
	if got := unexpectedHTTPStatusError(http.StatusForbidden).Error(); got != "unexpected HTTP status 403 Forbidden" {
		t.Errorf("unexpectedHTTPStatusError(403) = %q, want %q", got, "unexpected HTTP status 403 Forbidden")
	}
	if got := unexpectedHTTPStatusError(599).Error(); got != "unexpected HTTP status 599" {
		t.Errorf("unexpectedHTTPStatusError(599) = %q, want %q", got, "unexpected HTTP status 599")
	}
}
