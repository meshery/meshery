package models

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestPublishResultsAnonymousSharing verifies that the "Send Anonymous Performance
// Results" preference only controls whether a result is sent. The result is saved
// locally either way, so a run always appears in the profile's results.
func TestPublishResultsAnonymousSharing(t *testing.T) {
	tests := []struct {
		name     string
		sharing  bool
		wantSent int32
	}{
		{name: "sharing off keeps the result without sending it", sharing: false, wantSent: 0},
		{name: "sharing on keeps and sends the result", sharing: true, wantSent: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var sent atomic.Int32
			remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				sent.Add(1)
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":"` + uuid.Must(uuid.NewV4()).String() + `"}`))
			}))
			defer remote.Close()

			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			require.NoError(t, err)
			require.NoError(t, db.AutoMigrate(&MesheryResult{}))

			log, err := logger.New("test", logger.Options{})
			require.NoError(t, err)

			preferences, err := NewMapPreferencePersister()
			require.NoError(t, err)

			provider := &DefaultLocalProvider{
				MapPreferencePersister: preferences,
				ProviderBaseURL:        remote.URL,
				ResultPersister:        &MesheryResultsPersister{DB: &database.Handler{DB: db}},
				Log:                    log,
			}

			preference := NewDefaultPreference()
			preference.AnonymousPerfResults = tc.sharing
			require.NoError(t, provider.WriteToPersister(LocalProviderUser().ID.String(), preference))

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			result := &MesheryResult{Name: "test", Result: map[string]interface{}{"qps": 1.0}}

			resultID, err := provider.PublishResults(req, result, uuid.Must(uuid.NewV4()).String())
			require.NoError(t, err)

			var saved int64
			require.NoError(t, db.Table("meshery_results").Where("id = ?", resultID).Count(&saved).Error)
			assert.Equal(t, int64(1), saved, "the result must be saved locally")
			assert.Equal(t, tc.wantSent, sent.Load(), "the preference controls whether the result is sent")
		})
	}
}
