package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
)

func saveTestCredential(t *testing.T, provider *DefaultLocalProvider, userID uuid.UUID, name string) *Credential {
	t.Helper()
	cred, err := provider.SaveUserCredential("", &Credential{
		Name:   name,
		Type:   "token",
		UserId: userID,
		Secret: core.Map{"token": "old"},
	})
	if err != nil {
		t.Fatalf("save credential: %v", err)
	}
	return cred
}

func credentialDeletedAt(t *testing.T, provider *DefaultLocalProvider, id uuid.UUID) core.NullTime {
	t.Helper()
	var stored Credential
	if err := provider.GetGenericPersister().Where("id = ?", id).First(&stored).Error; err != nil {
		t.Fatalf("read credential %s: %v", id, err)
	}
	return stored.DeletedAt
}

func TestLocalUpdateUserCredential(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())

	t.Run("updates a live credential", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)
		cred := saveTestCredential(t, provider, userID, "before")

		got, err := provider.UpdateUserCredential(nil, &Credential{ID: cred.ID, UserId: userID, Name: "after"})
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		if got.Name != "after" {
			t.Errorf("name = %q, want %q", got.Name, "after")
		}
	})

	t.Run("unknown id is not found", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)

		_, err := provider.UpdateUserCredential(nil, &Credential{ID: uuid.Must(uuid.NewV4()), UserId: userID, Name: "x"})
		if errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s", err, ErrResultNotFoundCode)
		}
	})

	t.Run("deleted credential is not found and stays unchanged", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)
		cred := saveTestCredential(t, provider, userID, "before")
		if _, err := provider.DeleteUserCredential(nil, cred.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}

		got, err := provider.UpdateUserCredential(nil, &Credential{ID: cred.ID, UserId: userID, Name: "after"})
		if errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v (credential %+v), want code %s", err, got, ErrResultNotFoundCode)
		}
		var stored Credential
		if err := provider.GetGenericPersister().Where("id = ?", cred.ID).First(&stored).Error; err != nil {
			t.Fatalf("read back: %v", err)
		}
		if stored.Name != "before" {
			t.Errorf("deleted credential was renamed to %q", stored.Name)
		}
	})

	t.Run("another user's credential is not found", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)
		cred := saveTestCredential(t, provider, userID, "before")

		_, err := provider.UpdateUserCredential(nil, &Credential{ID: cred.ID, UserId: uuid.Must(uuid.NewV4()), Name: "after"})
		if errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s", err, ErrResultNotFoundCode)
		}
	})
}

func TestLocalDeleteUserCredential(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())

	t.Run("deletes a live credential", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)
		cred := saveTestCredential(t, provider, userID, "c")

		req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials", nil)
		req = req.WithContext(context.WithValue(req.Context(), UserCtxKey, &User{ID: userID}))

		if _, err := provider.DeleteUserCredential(req, cred.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}
		if !credentialDeletedAt(t, provider, cred.ID).Valid {
			t.Fatal("deleted_at was not set")
		}
	})

	t.Run("unknown id is not found", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)

		req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials", nil)
		req = req.WithContext(context.WithValue(req.Context(), UserCtxKey, &User{ID: userID}))

		_, err := provider.DeleteUserCredential(req, uuid.Must(uuid.NewV4()))
		if errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s", err, ErrResultNotFoundCode)
		}
	})

	t.Run("another user's credential is not deleted", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)
		cred := saveTestCredential(t, provider, userID, "c")

		otherUser := &User{ID: uuid.Must(uuid.NewV4())}
		req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials", nil)
		req = req.WithContext(context.WithValue(req.Context(), UserCtxKey, otherUser))

		_, err := provider.DeleteUserCredential(req, cred.ID)
		if errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s", err, ErrResultNotFoundCode)
		}
		if credentialDeletedAt(t, provider, cred.ID).Valid {
			t.Fatal("deleted_at was unexpectedly set for another user's credential")
		}
	})

	t.Run("fails closed when request or user is missing", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)
		cred := saveTestCredential(t, provider, userID, "c")

		// Nil request
		if _, err := provider.DeleteUserCredential(nil, cred.ID); errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s on nil req", err, ErrResultNotFoundCode)
		}
		// Request without UserCtxKey
		emptyReq := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials", nil)
		if _, err := provider.DeleteUserCredential(emptyReq, cred.ID); errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s on unauthenticated req", err, ErrResultNotFoundCode)
		}
		// Request with nil user
		nilUserReq := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials", nil)
		nilUserReq = nilUserReq.WithContext(context.WithValue(nilUserReq.Context(), UserCtxKey, (*User)(nil)))
		if _, err := provider.DeleteUserCredential(nilUserReq, cred.ID); errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s on nil user", err, ErrResultNotFoundCode)
		}
		// Request with uuid.Nil user
		zeroUserReq := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials", nil)
		zeroUserReq = zeroUserReq.WithContext(context.WithValue(zeroUserReq.Context(), UserCtxKey, &User{ID: uuid.Nil}))
		if _, err := provider.DeleteUserCredential(zeroUserReq, cred.ID); errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s on zero user ID", err, ErrResultNotFoundCode)
		}
		// Live credential was not deleted
		if credentialDeletedAt(t, provider, cred.ID).Valid {
			t.Fatal("deleted_at was unexpectedly set for unauthenticated delete attempt")
		}
	})

	t.Run("second delete is not found and keeps the first deleted_at", func(t *testing.T) {
		provider := newTestProviderWithCredentialDB(t)
		cred := saveTestCredential(t, provider, userID, "c")
		req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials", nil)
		req = req.WithContext(context.WithValue(req.Context(), UserCtxKey, &User{ID: userID}))

		if _, err := provider.DeleteUserCredential(req, cred.ID); err != nil {
			t.Fatalf("first delete: %v", err)
		}
		first := credentialDeletedAt(t, provider, cred.ID)

		time.Sleep(10 * time.Millisecond)
		_, err := provider.DeleteUserCredential(req, cred.ID)
		if errors.GetCode(err) != ErrResultNotFoundCode {
			t.Fatalf("err = %v, want code %s", err, ErrResultNotFoundCode)
		}
		if second := credentialDeletedAt(t, provider, cred.ID); !second.Time.Equal(first.Time) {
			t.Errorf("deleted_at moved from %v to %v", first.Time, second.Time)
		}
	})
}
