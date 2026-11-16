package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/require"
)

// postUserPrefs sends POST /api/user/prefs with body for a user whose saved
// load test preferences use generator.
func postUserPrefs(t *testing.T, generator, body string) *httptest.ResponseRecorder {
	t.Helper()
	persister, err := models.NewMapPreferencePersister()
	require.NoError(t, err)
	provider := &models.DefaultLocalProvider{MapPreferencePersister: persister}
	saved := &models.Preference{
		AnonymousUsageStats: true,
		LoadTestPreferences: &models.LoadTestPreferences{ConcurrentRequests: 1, Duration: "30s", LoadGenerator: generator},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs", strings.NewReader(body))
	rec := httptest.NewRecorder()
	(&Handler{log: newTestLogger(t)}).UserPrefsHandler(rec, req, saved, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

// A preference update that does not touch load test settings succeeds when the
// saved generator is one Meshery has removed: wrk2 and nighthawk both run on
// fortio now.
func TestUserPrefsHandlerAcceptsRemovedLoadGenerators(t *testing.T) {
	for _, generator := range []string{"wrk2", "nighthawk"} {
		t.Run(generator, func(t *testing.T) {
			rec := postUserPrefs(t, generator, `{"anonymousUsageStats":false}`)

			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
		})
	}
}

// A generator Meshery never supported is still rejected.
func TestUserPrefsHandlerRejectsUnknownLoadGenerator(t *testing.T) {
	rec := postUserPrefs(t, "fortio", `{"loadTestPrefs":{"gen":"jmeter"}}`)

	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
}
