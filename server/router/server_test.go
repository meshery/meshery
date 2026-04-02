package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/handlers"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	schemasConnection "github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/spf13/viper"
)

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	t.Log("Need to run Close() skipping")
	//err := r.Close()
	//if err != nil {
	//	t.Errorf("Close() failed with error: %s", err)
	//}
}

func TestFileEndpointsRequireProviderAuth(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	instanceID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to create instance id: %v", err)
	}

	previousInstanceID := viper.Get("INSTANCE_ID")
	viper.Set("INSTANCE_ID", &instanceID)
	t.Cleanup(func() {
		viper.Set("INSTANCE_ID", previousInstanceID)
	})

	handler := handlers.NewHandlerInstance(
		&models.HandlerConfig{
			ProviderCookieName: "meshery-provider",
			Providers:          map[string]models.Provider{},
		},
		nil,
		log,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		"",
		nil,
		nil,
		schemasConnection.MeshsyncDeploymentMode(""),
	)

	router := NewRouter(context.Background(), handler, 0, http.NotFoundHandler(), http.NotFoundHandler())

	testCases := []struct {
		name string
		path string
	}{
		{
			name: "view endpoint",
			path: "/api/system/fileView?file=/tmp/registry-logs.log",
		},
		{
			name: "download endpoint",
			path: "/api/system/fileDownload?file=/tmp/registry-logs.log",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			router.S.ServeHTTP(rec, req)

			resp := rec.Result()
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			})

			if resp.StatusCode != http.StatusFound {
				t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
			}

			if location := resp.Header.Get("Location"); !strings.HasPrefix(location, "/provider") {
				t.Fatalf("expected redirect to provider login, got %q", location)
			}
		})
	}
}
