package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// savedProfileProvider returns a saved performance profile and accepts the
// test result, so LoadTestHandler can run a real load test without a database.
type savedProfileProvider struct {
	models.Provider
	profile []byte
}

func (p *savedProfileProvider) GetPerformanceProfile(_ *http.Request, _ string) ([]byte, error) {
	return p.profile, nil
}

func (p *savedProfileProvider) PublishResults(_ *http.Request, _ *models.MesheryResult, _ string) (string, error) {
	return uuid.Must(uuid.NewV4()).String(), nil
}

// methodRecorder is a load test target that records the method of each request.
type methodRecorder struct {
	mu      sync.Mutex
	methods map[string]bool
}

func (m *methodRecorder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	m.methods[r.Method] = true
	m.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

func (m *methodRecorder) seen() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	methods := make([]string, 0, len(m.methods))
	for method := range m.methods {
		methods = append(methods, method)
	}
	sort.Strings(methods)
	return methods
}

// TestLoadTestHandlerAppliesSavedAdditionalOptions runs a saved profile against
// a local server and checks that the saved options reach the load generator,
// whichever shape they were saved in.
func TestLoadTestHandlerAppliesSavedAdditionalOptions(t *testing.T) {
	const putOptions = `{"MethodOverride":"PUT"}`

	tests := []struct {
		name    string
		options interface{}
		want    string
	}{
		// The UI saves the options as a one-element list.
		{name: "saved by the UI", options: []interface{}{putOptions}, want: http.MethodPut},
		// mesheryctl saves them as a string.
		{name: "saved by mesheryctl", options: putOptions, want: http.MethodPut},
		// The UI saves an empty field as a list holding an empty string.
		{name: "left empty in the UI", options: []interface{}{""}, want: http.MethodGet},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := &methodRecorder{methods: map[string]bool{}}
			server := httptest.NewServer(target)
			defer server.Close()

			profile, err := json.Marshal(models.PerformanceProfile{
				Name:     "saved-profile",
				Metadata: map[string]interface{}{"additional_options": tt.options},
			})
			require.NoError(t, err)

			h := &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t)}
			query := url.Values{
				"name": {"additional-options-test"},
				"url":  {server.URL},
				"t":    {"1"},
				"dur":  {"s"},
				"c":    {"1"},
				"qps":  {"5"},
			}
			profileID := uuid.Must(uuid.NewV4()).String()
			req := httptest.NewRequest(http.MethodGet, "/api/user/performance/profiles/"+profileID+"/run?"+query.Encode(), nil)
			req = mux.SetURLVars(req, map[string]string{"id": profileID})
			rec := httptest.NewRecorder()

			h.LoadTestHandler(rec, req, &models.Preference{}, &models.User{}, &savedProfileProvider{profile: profile})

			assert.Equal(t, []string{tt.want}, target.seen(), "request methods the target received")
		})
	}
}
