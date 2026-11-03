package machines

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

func TestExtractActorIDs(t *testing.T) {
	testUserID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}
	testSysID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}

	tests := []struct {
		name         string
		ctx          context.Context
		expectedUser core.Uuid
		expectedSys  *core.Uuid
	}{
		{
			name:         "nil context",
			ctx:          nil,
			expectedUser: core.Uuid(uuid.Nil),
			expectedSys:  nil,
		},
		{
			name:         "empty context",
			ctx:          context.Background(),
			expectedUser: core.Uuid(uuid.Nil),
			expectedSys:  nil,
		},
		{
			name: "valid user and pointer system ID",
			ctx: func() context.Context {
				sID := core.Uuid(testSysID)
				c := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(testUserID)})
				return context.WithValue(c, models.SystemIDKey, &sID)
			}(),
			expectedUser: core.Uuid(testUserID),
			expectedSys:  func() *core.Uuid { s := core.Uuid(testSysID); return &s }(),
		},
		{
			name: "valid user and value system ID",
			ctx: func() context.Context {
				c := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(testUserID)})
				return context.WithValue(c, models.SystemIDKey, core.Uuid(testSysID))
			}(),
			expectedUser: core.Uuid(testUserID),
			expectedSys:  func() *core.Uuid { s := core.Uuid(testSysID); return &s }(),
		},
		{
			name: "typed nil *models.User",
			ctx: func() context.Context {
				var nilUser *models.User
				return context.WithValue(context.Background(), models.UserCtxKey, nilUser)
			}(),
			expectedUser: core.Uuid(uuid.Nil),
			expectedSys:  nil,
		},
		{
			name: "typed nil *core.Uuid",
			ctx: func() context.Context {
				var nilSys *core.Uuid
				return context.WithValue(context.Background(), models.SystemIDKey, nilSys)
			}(),
			expectedUser: core.Uuid(uuid.Nil),
			expectedSys:  nil,
		},
		{
			name: "mismatched types in context",
			ctx: func() context.Context {
				c := context.WithValue(context.Background(), models.UserCtxKey, "not-a-user")
				return context.WithValue(c, models.SystemIDKey, 12345)
			}(),
			expectedUser: core.Uuid(uuid.Nil),
			expectedSys:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, sys := ExtractActorIDs(tt.ctx)
			if user != tt.expectedUser {
				t.Errorf("ExtractActorIDs user = %v, want %v", user, tt.expectedUser)
			}
			if (sys == nil) != (tt.expectedSys == nil) {
				t.Fatalf("ExtractActorIDs sys nilness = %v, want %v", sys == nil, tt.expectedSys == nil)
			}
			if sys != nil && *sys != *tt.expectedSys {
				t.Errorf("ExtractActorIDs sys = %v, want %v", *sys, *tt.expectedSys)
			}
		})
	}
}

func TestExtractProvider(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		wantOk   bool
		nilCheck bool
	}{
		{
			name:   "nil context",
			ctx:    nil,
			wantOk: false,
		},
		{
			name:   "empty context",
			ctx:    context.Background(),
			wantOk: false,
		},
		{
			name: "nil interface provider",
			ctx: func() context.Context {
				var p models.Provider
				return context.WithValue(context.Background(), models.ProviderCtxKey, p)
			}(),
			wantOk: false,
		},
		{
			name: "typed nil concrete pointer implementing Provider",
			ctx: func() context.Context {
				var fp *fakeProvider
				return context.WithValue(context.Background(), models.ProviderCtxKey, fp)
			}(),
			wantOk: false,
		},
		{
			name: "typed nil concrete pointer assigned to Provider interface",
			ctx: func() context.Context {
				var fp *fakeProvider
				var p models.Provider = fp
				return context.WithValue(context.Background(), models.ProviderCtxKey, p)
			}(),
			wantOk: false,
		},
		{
			name: "valid provider",
			ctx: func() context.Context {
				fp := &fakeProvider{}
				return context.WithValue(context.Background(), models.ProviderCtxKey, fp)
			}(),
			wantOk: true,
		},
		{
			name: "wrong type in ProviderCtxKey",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), models.ProviderCtxKey, "not-a-provider")
			}(),
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, ok := ExtractProvider(tt.ctx)
			if ok != tt.wantOk {
				t.Errorf("ExtractProvider ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && p == nil {
				t.Errorf("ExtractProvider returned ok=true but nil provider")
			}
		})
	}
}

func TestApplyEventAttribution(t *testing.T) {
	testUserID := core.Uuid(uuid.Must(uuid.NewV4()))
	testSysID := core.Uuid(uuid.Must(uuid.NewV4()))

	t.Run("nil event builder does not panic", func(t *testing.T) {
		res := ApplyEventAttribution(nil, testUserID, &testSysID)
		if res != nil {
			t.Errorf("expected nil result for nil event builder")
		}
	})

	t.Run("both user and system present", func(t *testing.T) {
		eb := events.NewEvent().WithCategory("test")
		event := ApplyEventAttribution(eb, testUserID, &testSysID).Build()
		if event.Owner == nil || *event.Owner != testUserID {
			t.Errorf("event.Owner = %v, want %v", event.Owner, testUserID)
		}
		if event.SystemID != testSysID {
			t.Errorf("event.SystemID = %v, want %v", event.SystemID, testSysID)
		}
	})

	t.Run("only user present", func(t *testing.T) {
		eb := events.NewEvent().WithCategory("test")
		event := ApplyEventAttribution(eb, testUserID, nil).Build()
		if event.Owner == nil || *event.Owner != testUserID {
			t.Errorf("event.Owner = %v, want %v", event.Owner, testUserID)
		}
		if event.SystemID != core.Uuid(uuid.Nil) {
			t.Errorf("event.SystemID = %v, want nil uuid", event.SystemID)
		}
	})

	t.Run("only system present", func(t *testing.T) {
		eb := events.NewEvent().WithCategory("test")
		event := ApplyEventAttribution(eb, core.Uuid(uuid.Nil), &testSysID).Build()
		if event.Owner != nil {
			t.Errorf("expected nil Owner, got %v", event.Owner)
		}
		if event.SystemID != testSysID {
			t.Errorf("event.SystemID = %v, want %v", event.SystemID, testSysID)
		}
	})

	t.Run("neither ID present does not panic", func(t *testing.T) {
		eb := events.NewEvent().WithCategory("test")
		event := ApplyEventAttribution(eb, core.Uuid(uuid.Nil), nil).Build()
		if event.Owner != nil {
			t.Errorf("expected nil Owner, got %v", event.Owner)
		}
		if event.SystemID != core.Uuid(uuid.Nil) {
			t.Errorf("expected nil SystemID, got %v", event.SystemID)
		}
	})
}
