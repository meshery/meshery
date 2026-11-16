package models

import (
	"errors"
	"testing"

	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
)

func TestPatternResourcePersisterNil(t *testing.T) {
	for _, p := range []*PatternResourcePersister{nil, {}, {DB: &database.Handler{}}} {
		id := mustUUID(t)
		if _, err := p.SavePatternResource(nil); !errors.Is(err, ErrDBConnection) {
			t.Fatalf("save: %v", err)
		}
		if err := p.DeletePatternResource(id); !errors.Is(err, ErrDBConnection) {
			t.Fatalf("delete: %v", err)
		}
		if _, err := p.GetPatternResource(id); !errors.Is(err, ErrDBConnection) {
			t.Fatalf("get: %v", err)
		}
		if _, err := p.GetPatternResourceByAttributes("", "", "", ""); !errors.Is(err, ErrDBConnection) {
			t.Fatalf("attributes: %v", err)
		}
		if _, err := p.GetPatternResources("", "", "", "", "", "", 0, 10); !errors.Is(err, ErrDBConnection) {
			t.Fatalf("list: %v", err)
		}
		if p.Exists("", "", "", "") {
			t.Fatal("nil persister reports a resource")
		}
	}
}

func TestPatternResourcePersistenceLifecycle(t *testing.T) {
	p := &PatternResourcePersister{DB: newMigratedDB(t)}
	if _, err := p.SavePatternResource(nil); !errors.Is(err, ErrNilConfigData) {
		t.Fatalf("nil save: %v", err)
	}
	resource := &PatternResource{Name: "resource", Namespace: "ns", Type: "type", OAMType: "oam"}
	if _, err := p.SavePatternResource(resource); err != nil {
		t.Fatal(err)
	}
	if _, err := p.GetPatternResource(*resource.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := p.GetPatternResourceByAttributes("resource", "ns", "type", "oam"); err != nil {
		t.Fatal(err)
	}
	if !p.Exists("resource", "ns", "type", "oam") {
		t.Fatal("resource missing")
	}
	page, err := p.GetPatternResources("", "", "", "", "", "", 0, 10)
	if err != nil || page.TotalCount != 1 || len(page.Resources) != 1 {
		t.Fatalf("list: %#v, %v", page, err)
	}
	if err := p.DeletePatternResource(*resource.ID); err != nil {
		t.Fatal(err)
	}
	_, err = p.GetPatternResource(*resource.ID)
	if mkerrors.GetCode(err) != ErrResultNotFoundCode {
		t.Fatalf("deleted get: %v", err)
	}
	_, err = p.GetPatternResourceByAttributes("resource", "ns", "type", "oam")
	if mkerrors.GetCode(err) != ErrResultNotFoundCode {
		t.Fatalf("deleted lookup: %v", err)
	}
	if err := p.DeletePatternResource(*resource.ID); mkerrors.GetCode(err) != ErrResultNotFoundCode {
		t.Fatalf("repeat delete: %v", err)
	}
	if p.Exists("resource", "ns", "type", "oam") {
		t.Fatal("deleted resource exists")
	}
	if err := p.DB.Migrator().DropTable(&PatternResource{}); err != nil {
		t.Fatal(err)
	}
	if _, err := p.GetPatternResources("", "", "", "", "", "", 0, 10); err == nil {
		t.Fatal("list swallowed database error")
	}
}
