package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newApplicationPersisterTestDB(t *testing.T) (*gorm.DB, *MesheryApplicationPersister) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&MesheryApplication{}); err != nil {
		t.Fatalf("migrate meshery_applications: %v", err)
	}
	return db, &MesheryApplicationPersister{DB: &database.Handler{DB: db}}
}

func TestSaveAndGetMesheryApplications(t *testing.T) {
	_, persister := newApplicationPersisterTestDB(t)

	now := time.Now()
	apps := []MesheryApplication{
		{
			Name:            "App Beta",
			ApplicationFile: "app-beta.yaml",
			SourceContent:   []byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: beta"),
			UpdatedAt:       &now,
			CreatedAt:       &now,
		},
		{
			Name:            "App Alpha",
			ApplicationFile: "app-alpha.yaml",
			SourceContent:   []byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: alpha"),
			UpdatedAt:       &now,
			CreatedAt:       &now,
		},
	}

	savedBytes, err := persister.SaveMesheryApplications(apps)
	if err != nil {
		t.Fatalf("unexpected error saving applications: %v", err)
	}
	var saved []MesheryApplication
	if err := json.Unmarshal(savedBytes, &saved); err != nil {
		t.Fatalf("unmarshal saved applications: %v", err)
	}
	if len(saved) != 2 {
		t.Fatalf("saved count = %d, want 2", len(saved))
	}
	if saved[0].ID == nil || saved[1].ID == nil {
		t.Fatalf("expected non-nil generated IDs for saved applications")
	}

	// Test GetMesheryApplications pagination and search
	pageBytes, err := persister.GetMesheryApplications("", "name asc", 0, 10, "")
	if err != nil {
		t.Fatalf("unexpected error getting applications: %v", err)
	}
	var page MesheryApplicationPage
	if err := json.Unmarshal(pageBytes, &page); err != nil {
		t.Fatalf("unmarshal application page: %v", err)
	}
	if page.TotalCount != 2 {
		t.Fatalf("totalCount = %d, want 2", page.TotalCount)
	}
	if len(page.Applications) != 2 {
		t.Fatalf("applications len = %d, want 2", len(page.Applications))
	}
	if page.Applications[0].Name != "App Alpha" {
		t.Errorf("expected first application to be 'App Alpha' under asc sort, got %q", page.Applications[0].Name)
	}

	// Test search filter
	searchBytes, err := persister.GetMesheryApplications("Alpha", "", 0, 10, "")
	if err != nil {
		t.Fatalf("unexpected error searching applications: %v", err)
	}
	var searchPage MesheryApplicationPage
	if err := json.Unmarshal(searchBytes, &searchPage); err != nil {
		t.Fatalf("unmarshal search page: %v", err)
	}
	if searchPage.TotalCount != 1 || len(searchPage.Applications) != 1 {
		t.Fatalf("search totalCount = %d, want 1", searchPage.TotalCount)
	}
	if searchPage.Applications[0].Name != "App Alpha" {
		t.Errorf("search result = %q, want 'App Alpha'", searchPage.Applications[0].Name)
	}
}

func TestGetMesheryApplication_Existing(t *testing.T) {
	_, persister := newApplicationPersisterTestDB(t)

	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	apps := []MesheryApplication{
		{
			ID:              &id,
			Name:            "Single App",
			ApplicationFile: "single.yaml",
			SourceContent:   []byte("source content bytes"),
			UpdatedAt:       &now,
			CreatedAt:       &now,
		},
	}

	if _, err := persister.SaveMesheryApplications(apps); err != nil {
		t.Fatalf("saving application: %v", err)
	}

	appBytes, err := persister.GetMesheryApplication(id)
	if err != nil {
		t.Fatalf("unexpected error getting application: %v", err)
	}
	var app MesheryApplication
	if err := json.Unmarshal(appBytes, &app); err != nil {
		t.Fatalf("unmarshal application: %v", err)
	}
	if app.ID == nil || *app.ID != id {
		t.Errorf("retrieved ID = %v, want %v", app.ID, id)
	}
	if app.Name != "Single App" {
		t.Errorf("retrieved Name = %q, want 'Single App'", app.Name)
	}
}

func TestGetMesheryApplication_NotFound(t *testing.T) {
	_, persister := newApplicationPersisterTestDB(t)

	nonExistentID := uuid.Must(uuid.NewV4())
	_, err := persister.GetMesheryApplication(nonExistentID)
	if err == nil {
		t.Fatalf("expected error getting nonexistent application, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrResultNotFoundCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrResultNotFoundCode, err)
	}
}

func TestGetMesheryApplicationSource_Existing(t *testing.T) {
	_, persister := newApplicationPersisterTestDB(t)

	id := uuid.Must(uuid.NewV4())
	expectedContent := []byte("source content test")
	apps := []MesheryApplication{
		{
			ID:              &id,
			Name:            "Source App",
			ApplicationFile: "source.yaml",
			SourceContent:   expectedContent,
		},
	}

	if _, err := persister.SaveMesheryApplications(apps); err != nil {
		t.Fatalf("saving application: %v", err)
	}

	content, err := persister.GetMesheryApplicationSource(id)
	if err != nil {
		t.Fatalf("unexpected error getting application source: %v", err)
	}
	if string(content) != string(expectedContent) {
		t.Errorf("content = %q, want %q", string(content), string(expectedContent))
	}
}

func TestGetMesheryApplicationSource_NotFound(t *testing.T) {
	_, persister := newApplicationPersisterTestDB(t)

	nonExistentID := uuid.Must(uuid.NewV4())
	_, err := persister.GetMesheryApplicationSource(nonExistentID)
	if err == nil {
		t.Fatalf("expected error getting nonexistent application source, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrResultNotFoundCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrResultNotFoundCode, err)
	}
}

func TestDeleteMesheryApplication_Existing(t *testing.T) {
	_, persister := newApplicationPersisterTestDB(t)

	id := uuid.Must(uuid.NewV4())
	apps := []MesheryApplication{
		{
			ID:              &id,
			Name:            "Delete Me",
			ApplicationFile: "delete.yaml",
		},
	}

	if _, err := persister.SaveMesheryApplications(apps); err != nil {
		t.Fatalf("saving application: %v", err)
	}

	resp, err := persister.DeleteMesheryApplication(id)
	if err != nil {
		t.Fatalf("unexpected error deleting application: %v", err)
	}
	if len(resp) == 0 {
		t.Errorf("expected non-empty response bytes, got empty")
	}

	// Deleting again should return ErrResultNotFound
	_, err = persister.DeleteMesheryApplication(id)
	if err == nil {
		t.Fatalf("expected error deleting already-deleted application, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrResultNotFoundCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrResultNotFoundCode, err)
	}
}

func TestDeleteMesheryApplication_NotFound(t *testing.T) {
	_, persister := newApplicationPersisterTestDB(t)

	nonExistentID := uuid.Must(uuid.NewV4())
	_, err := persister.DeleteMesheryApplication(nonExistentID)
	if err == nil {
		t.Fatalf("expected error deleting nonexistent application, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrResultNotFoundCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrResultNotFoundCode, err)
	}
}
