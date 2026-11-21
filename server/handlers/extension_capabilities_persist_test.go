package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"gorm.io/gorm"
)

// newPersistingLocalProvider returns a local provider whose capabilities are
// stored in a fresh SQLite database, like the one the server runs with.
func newPersistingLocalProvider(t *testing.T) (*models.DefaultLocalProvider, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "meshery.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&models.UserCapabilities{}); err != nil {
		t.Fatalf("migrate user capabilities: %v", err)
	}
	lp := newTestLocalProvider()
	lp.UserCapabilitiesPersister = &models.UserCapabilitiesPersister{DB: &db}
	return lp, &db
}

func installAccountExtension(t *testing.T, h *Handler, lp *models.DefaultLocalProvider, user *models.User, title string) {
	t.Helper()
	body := fmt.Sprintf(`{"extType":"account","extensionMetadata":{"title":%q}}`, title)
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.InstallExtensionHandler(rec, req, nil, user, lp)
	if rec.Code != http.StatusOK {
		t.Errorf("install %q: status %d, body %s", title, rec.Code, rec.Body.String())
	}
}

func persistedAccountTitles(t *testing.T, lp *models.DefaultLocalProvider, user *models.User) []string {
	t.Helper()
	props, err := lp.ReadCapabilitiesForUser(user.ID.String())
	if err != nil {
		t.Fatalf("read persisted capabilities: %v", err)
	}
	titles := make([]string, 0, len(props.Extensions.Account))
	for _, ext := range props.Extensions.Account {
		titles = append(titles, ext.Title)
	}
	sort.Strings(titles)
	return titles
}

// An install that is slow to write its capabilities must not overwrite the
// capabilities a later install already stored. /api/provider/capabilities
// serves the stored copy, so an overwritten extension disappears from the UI
// even though the provider still has it.
func TestInstallExtensionHandler_SlowPersistDoesNotOverwriteNewerCapabilities(t *testing.T) {
	lp, db := newPersistingLocalProvider(t)
	h := newTestHandler(t, map[string]models.Provider{}, "")
	userID, _ := uuid.NewV4()
	user := &models.User{ID: userID}

	// Store an initial row so both installs below go through an update.
	props := lp.GetProviderProperties()
	if err := lp.WriteCapabilitiesForUser(user.ID.String(), &props); err != nil {
		t.Fatalf("seed capabilities: %v", err)
	}

	// Hold the first capabilities write that follows until released.
	var armed atomic.Bool
	paused := make(chan struct{})
	release := make(chan struct{})
	err := db.Callback().Update().Before("gorm:begin_transaction").Register("test:hold_first_capabilities_write", func(tx *gorm.DB) {
		if _, ok := tx.Statement.Model.(*models.UserCapabilities); !ok {
			return
		}
		if armed.CompareAndSwap(true, false) {
			close(paused)
			<-release
		}
	})
	if err != nil {
		t.Fatalf("register callback: %v", err)
	}
	armed.Store(true)

	firstDone := make(chan struct{})
	go func() {
		defer close(firstDone)
		installAccountExtension(t, h, lp, user, "first")
	}()

	select {
	case <-paused:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("first install never reached its capabilities write")
	}

	secondDone := make(chan struct{})
	go func() {
		defer close(secondDone)
		installAccountExtension(t, h, lp, user, "second")
	}()

	// Let the second install finish first if it can. If it has to wait for
	// the first one's write, release that write after a while.
	select {
	case <-secondDone:
	case <-time.After(time.Second):
	}
	close(release)
	<-firstDone
	<-secondDone

	got := persistedAccountTitles(t, lp, user)
	want := []string{"first", "second"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("persisted account extensions = %v, want %v (provider has %d)", got, want, len(lp.GetProviderProperties().Extensions.Account))
	}
}

// Installs and removes running at the same time must leave the stored
// capabilities equal to the provider's own state. Run with -race.
func TestExtensionHandlers_ConcurrentInstallAndRemoveKeepStoredCapabilitiesInSync(t *testing.T) {
	lp, _ := newPersistingLocalProvider(t)
	h := newTestHandler(t, map[string]models.Provider{}, "")
	userID, _ := uuid.NewV4()
	user := &models.User{ID: userID}

	const n = 8
	for i := range n {
		installAccountExtension(t, h, lp, user, fmt.Sprintf("old-%d", i))
	}

	var wg sync.WaitGroup
	for i := range n {
		wg.Add(2)
		go func() {
			defer wg.Done()
			installAccountExtension(t, h, lp, user, fmt.Sprintf("new-%d", i))
		}()
		go func() {
			defer wg.Done()
			body := fmt.Sprintf(`{"extType":"account","title":"old-%d"}`, i)
			req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", strings.NewReader(body))
			rec := httptest.NewRecorder()
			h.RemoveExtensionHandler(rec, req, nil, user, lp)
			if rec.Code != http.StatusOK {
				t.Errorf("remove old-%d: status %d, body %s", i, rec.Code, rec.Body.String())
			}
		}()
	}
	wg.Wait()

	got := persistedAccountTitles(t, lp, user)
	want := make([]string, 0, n)
	for i := range n {
		want = append(want, fmt.Sprintf("new-%d", i))
	}
	sort.Strings(want)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("persisted account extensions = %v, want %v", got, want)
	}
}
