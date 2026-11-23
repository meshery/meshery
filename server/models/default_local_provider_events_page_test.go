package models

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
)

// GetEvents must echo the requested page; the handler derives the offset from it,
// so a response that reports page 0 contradicts the rows it returns.
func TestDefaultLocalProviderGetEventsEchoesRequestedPage(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := AutoMigrateSystemTables(&db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	l := &DefaultLocalProvider{EventsPersister: &EventsPersister{DB: &db}}
	resp, err := l.GetEvents("", &events.EventsFilter{Limit: 10, Offset: 20, SortOn: "created_at", Order: "desc"}, 2, uuid.Nil, uuid.Nil)
	if err != nil {
		t.Fatalf("GetEvents: %v", err)
	}
	if resp.Page != 2 {
		t.Errorf("Page = %d, want 2", resp.Page)
	}
}
