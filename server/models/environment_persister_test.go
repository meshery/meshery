package models

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta3/environment"
)

func TestEnvironmentPersisterUpdateEnvironmentByID(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&environment.Environment{}); err != nil {
		t.Fatalf("migrate environments: %v", err)
	}
	persister := &EnvironmentPersister{DB: &db}

	t.Run("unknown ID returns not found without inserting", func(t *testing.T) {
		unknown := &environment.Environment{
			ID:   uuid.Must(uuid.NewV4()),
			Name: "must not be inserted",
		}

		var before int64
		if err := db.Model(&environment.Environment{}).Count(&before).Error; err != nil {
			t.Fatalf("count before update: %v", err)
		}

		_, err := persister.UpdateEnvironmentByID(unknown)
		if err == nil {
			t.Fatal("update of an unknown environment returned success")
		}
		if code := meshkiterrors.GetCode(err); code != ErrResultNotFoundCode {
			t.Fatalf("error code = %q, want %q", code, ErrResultNotFoundCode)
		}

		var after int64
		if err := db.Model(&environment.Environment{}).Count(&after).Error; err != nil {
			t.Fatalf("count after update: %v", err)
		}
		if after != before {
			t.Fatalf("row count = %d, want %d", after, before)
		}
	})

	t.Run("known ID is updated", func(t *testing.T) {
		id := uuid.Must(uuid.NewV4())
		original := &environment.Environment{ID: id, Name: "before", Description: "old"}
		if err := db.Create(original).Error; err != nil {
			t.Fatalf("create environment: %v", err)
		}

		original.Name = "after"
		original.Description = ""
		updated, err := persister.UpdateEnvironmentByID(original)
		if err != nil {
			t.Fatalf("update known environment: %v", err)
		}
		if updated.Name != "after" || updated.Description != "" {
			t.Fatalf("updated environment = {%q, %q}, want {%q, %q}", updated.Name, updated.Description, "after", "")
		}

		var stored environment.Environment
		if err := db.Where("id = ?", id).First(&stored).Error; err != nil {
			t.Fatalf("read updated environment: %v", err)
		}
		if stored.Name != "after" || stored.Description != "" {
			t.Fatalf("stored environment = {%q, %q}, want {%q, %q}", stored.Name, stored.Description, "after", "")
		}
	})
}
