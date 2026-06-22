package models

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/meshery/meshkit/logger"
)

func TestExtensionProxy_PanicFix(t *testing.T) {
	log, err := logger.New("meshery", logger.Options{Format: logger.SyslogLogFormat, LogLevel: 5})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	l := &RemoteProvider{Log: log}

	// This path does not contain /api/extensions, which previously caused an index out of bounds panic
	req := &http.Request{
		URL: &url.URL{
			Path: "/invalid/path",
		},
	}

	// The function should now return an error instead of panicking
	_, err = l.ExtensionProxy(req)
	if err == nil {
		t.Errorf("expected error, got nil")
	} else if !strings.Contains(err.Error(), "invalid extension proxy path") {
		t.Errorf("unexpected error: %v", err)
	}
}
