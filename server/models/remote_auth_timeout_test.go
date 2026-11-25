package models

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newUnresponsiveProvider returns a RemoteProvider whose remote provider accepts a
// request and then never answers it, and shortens the client timeout so the test
// does not wait for the production value.
func newUnresponsiveProvider(t *testing.T) *RemoteProvider {
	t.Helper()

	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		<-release
	}))
	// Cleanups run last in, first out: free the blocked handler before Close waits for it.
	t.Cleanup(srv.Close)
	t.Cleanup(func() { close(release) })

	previous := remoteAuthClient.Timeout
	remoteAuthClient.Timeout = 200 * time.Millisecond
	t.Cleanup(func() { remoteAuthClient.Timeout = previous })

	return newTestRemoteProvider(t, srv.URL)
}

// callWithin runs fn and fails the test if it has not returned within d.
func callWithin(t *testing.T, d time.Duration, what string, fn func() error) error {
	t.Helper()

	done := make(chan error, 1)
	go func() { done <- fn() }()
	select {
	case err := <-done:
		return err
	case <-time.After(d):
		t.Fatalf("%s did not return within %s: the remote provider request has no timeout", what, d)
		return nil
	}
}

// refreshToken holds TokenStoreMut while it calls the remote provider. Without a
// timeout, one unresponsive provider would hold the lock forever and block every
// other token refresh.
func TestRefreshToken_UnresponsiveProviderTimesOut(t *testing.T) {
	provider := newUnresponsiveProvider(t)

	err := callWithin(t, 5*time.Second, "refreshToken", func() error {
		_, err := provider.refreshToken("old-token")
		return err
	})
	if err == nil {
		t.Fatal("expected an error from an unresponsive remote provider")
	}
	if !provider.TokenStoreMut.TryLock() {
		t.Fatal("TokenStoreMut is still held after refreshToken failed")
	}
	provider.TokenStoreMut.Unlock()
}

func TestUpdateJWKs_UnresponsiveProviderTimesOut(t *testing.T) {
	provider := newUnresponsiveProvider(t)

	err := callWithin(t, 5*time.Second, "UpdateJWKs", provider.UpdateJWKs)
	if err == nil {
		t.Fatal("expected an error from an unresponsive remote provider")
	}
}
