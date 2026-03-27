package handlers

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"
)

type trackingReadCloser struct {
	*bytes.Reader
}

func (t *trackingReadCloser) Close() error {
	return nil
}

func TestReadBodyWithLimitDoesNotReplaceRequestBody(t *testing.T) {
	originalBody := &trackingReadCloser{Reader: bytes.NewReader([]byte(`{"name":"meshery"}`))}
	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Body = originalBody

	body, err := readBodyWithLimit(req, DefaultMaxBodySize)
	if err != nil {
		t.Fatalf("readBodyWithLimit() error = %v", err)
	}

	if string(body) != `{"name":"meshery"}` {
		t.Fatalf("readBodyWithLimit() body = %q, want %q", string(body), `{"name":"meshery"}`)
	}

	if req.Body != originalBody {
		t.Fatal("readBodyWithLimit() unexpectedly replaced request body")
	}
}

func TestReadBodyWithLimitReturnsHelpfulErrorWhenBodyTooLarge(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/test", strings.NewReader("abcdef"))

	_, err := readBodyWithLimit(req, 3)
	if err == nil {
		t.Fatal("readBodyWithLimit() error = nil, want body too large error")
	}

	expected := "request body exceeds maximum allowed size of 3 bytes"
	if err.Error() != expected {
		t.Fatalf("readBodyWithLimit() error = %q, want %q", err.Error(), expected)
	}
}
