package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/meshery/meshkit/database"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newK8sContextPersister(t *testing.T, migrate bool) *MesheryK8sContextPersister {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}
	// Each pooled connection to ":memory:" gets its own empty database.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to access underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })

	if migrate {
		if err := db.AutoMigrate(&K8sContext{}); err != nil {
			t.Fatalf("failed to migrate k8s contexts: %v", err)
		}
	}
	return &MesheryK8sContextPersister{DB: &database.Handler{DB: db}}
}

func TestGetMesheryK8sContexts_ReturnsPage(t *testing.T) {
	persister := newK8sContextPersister(t, true)
	for _, c := range []K8sContext{{ID: "1", Name: "alpha"}, {ID: "2", Name: "beta"}, {ID: "3", Name: "alphabet"}} {
		if err := persister.DB.Create(&c).Error; err != nil {
			t.Fatalf("failed to seed context %s: %v", c.Name, err)
		}
	}

	resp, err := persister.GetMesheryK8sContexts("ALPHA", "name asc", 0, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var page MesheryK8sContextPage
	if err := json.Unmarshal(resp, &page); err != nil {
		t.Fatalf("failed to decode page: %v", err)
	}
	if page.TotalCount != 2 {
		t.Fatalf("expected total count 2, got %d", page.TotalCount)
	}
	if len(page.Contexts) != 1 || page.Contexts[0].Name != "alpha" {
		t.Fatalf("expected only context alpha on the page, got %+v", page.Contexts)
	}
}

func TestGetMesheryK8sContexts_PropagatesCountError(t *testing.T) {
	persister := newK8sContextPersister(t, false)

	resp, err := persister.GetMesheryK8sContexts("", "", 0, 10)
	if err == nil {
		t.Fatalf("expected an error, got response %s", resp)
	}
	if resp != nil {
		t.Fatalf("expected no response body, got %s", resp)
	}
	if code := meshkiterrors.GetCode(err); code != ErrDBReadCode {
		t.Fatalf("expected error code %s, got %s", ErrDBReadCode, code)
	}
}

func TestGetMesheryK8sContexts_PropagatesFindError(t *testing.T) {
	persister := newK8sContextPersister(t, true)
	// COUNT never decodes a row, so an undecodable column fails only in Find.
	if err := persister.DB.Exec(`INSERT INTO k8s_contexts (id, name, auth) VALUES ('1', 'corrupt', 'not-json')`).Error; err != nil {
		t.Fatalf("failed to seed corrupt context: %v", err)
	}

	resp, err := persister.GetMesheryK8sContexts("", "", 0, 10)
	if err == nil {
		t.Fatalf("expected an error, got response %s", resp)
	}
	if resp != nil {
		t.Fatalf("expected no response body, got %s", resp)
	}
	if code := meshkiterrors.GetCode(err); code != ErrDBReadCode {
		t.Fatalf("expected error code %s, got %s", ErrDBReadCode, code)
	}
	if !strings.Contains(err.Error(), "Unable to read data") {
		t.Fatalf("expected a read error, got %v", err)
	}
}
