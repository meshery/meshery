package handlers

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models"
	"google.golang.org/grpc"
)

type adapterConfigSpyProvider struct {
	*models.DefaultLocalProvider
	recorded atomic.Int32
}

func newAdapterConfigSpyProvider() *adapterConfigSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &adapterConfigSpyProvider{DefaultLocalProvider: base}
}

func (p *adapterConfigSpyProvider) RecordPreferences(_ *http.Request, _ string, pref *models.Preference) error {
	p.recorded.Add(1)
	return nil
}

type fakeMeshService struct {
	meshes.UnimplementedMeshServiceServer
	supportedOpsCalls  atomic.Int32
	componentInfoCalls atomic.Int32
	supportedOpsErr    error
	componentInfoErr   error
}

func (f *fakeMeshService) SupportedOperations(context.Context, *meshes.SupportedOperationsRequest) (*meshes.SupportedOperationsResponse, error) {
	f.supportedOpsCalls.Add(1)
	if f.supportedOpsErr != nil {
		return nil, f.supportedOpsErr
	}
	return &meshes.SupportedOperationsResponse{
		Ops: []*meshes.SupportedOperation{
			{
				Key:      "install",
				Value:    "Install",
				Category: meshes.OpCategory_INSTALL,
			},
		},
	}, nil
}

func (f *fakeMeshService) ComponentInfo(context.Context, *meshes.ComponentInfoRequest) (*meshes.ComponentInfoResponse, error) {
	f.componentInfoCalls.Add(1)
	if f.componentInfoErr != nil {
		return nil, f.componentInfoErr
	}
	return &meshes.ComponentInfoResponse{
		Type:    "adapter",
		Name:    "meshery-test-adapter",
		Version: "v1.2.3",
		GitSha:  "abc123",
	}, nil
}

func TestMeshAdapterConfigHandler_PostEmptyMeshLocationURL(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	h.config.AdapterTracker = helpers.NewAdaptersTracker(nil)
	provider := newAdapterConfigSpyProvider()

	req := httptest.NewRequest(http.MethodPost, "/api/system/adapter/manage", nil)
	rec := httptest.NewRecorder()

	h.MeshAdapterConfigHandler(rec, req, models.NewDefaultPreference(), &models.User{UserId: "user-1"}, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d (body=%q)", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if provider.recorded.Load() != 0 {
		t.Fatalf("RecordPreferences called %d times, want 0", provider.recorded.Load())
	}
}

func TestMeshAdapterConfigHandler_PostRegistersReachableAdapter(t *testing.T) {
	fakeAdapter, location := startFakeMeshService(t)
	h := newTestHandler(t, map[string]models.Provider{}, "")
	h.config.AdapterTracker = helpers.NewAdaptersTracker(nil)
	pref := models.NewDefaultPreference()
	provider := newAdapterConfigSpyProvider()

	req := newAdapterConfigRequest(t, http.MethodPost, location)
	rec := httptest.NewRecorder()

	h.MeshAdapterConfigHandler(rec, req, pref, &models.User{UserId: "user-1"}, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body=%q)", http.StatusOK, rec.Code, rec.Body.String())
	}
	if fakeAdapter.supportedOpsCalls.Load() != 1 {
		t.Fatalf("SupportedOperations called %d times, want 1", fakeAdapter.supportedOpsCalls.Load())
	}
	if fakeAdapter.componentInfoCalls.Load() != 1 {
		t.Fatalf("ComponentInfo called %d times, want 1", fakeAdapter.componentInfoCalls.Load())
	}
	if provider.recorded.Load() != 1 {
		t.Fatalf("RecordPreferences called %d times, want 1", provider.recorded.Load())
	}

	got := decodeAdaptersResponse(t, rec)
	assertSingleAdapter(t, got, location)
	assertSingleAdapter(t, pref.MeshAdapters, location)

	tracked := h.config.AdapterTracker.GetAdapters(context.Background())
	if len(tracked) != 1 || tracked[0].Location != location {
		t.Fatalf("tracked adapters=%#v, want one adapter at %q", tracked, location)
	}
}

func TestMeshAdapterConfigHandler_PostDuplicateAdapterDoesNotDuplicate(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	location := "meshery-istio:10000"
	h.config.AdapterTracker = helpers.NewAdaptersTracker([]string{location})
	pref := models.NewDefaultPreference()
	pref.MeshAdapters = []*models.Adapter{{Location: location, Name: "meshery-istio"}}
	provider := newAdapterConfigSpyProvider()

	req := newAdapterConfigRequest(t, http.MethodPost, location)
	rec := httptest.NewRecorder()

	h.MeshAdapterConfigHandler(rec, req, pref, &models.User{UserId: "user-1"}, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body=%q)", http.StatusOK, rec.Code, rec.Body.String())
	}
	if provider.recorded.Load() != 1 {
		t.Fatalf("RecordPreferences called %d times, want 1", provider.recorded.Load())
	}
	got := decodeAdaptersResponse(t, rec)
	if len(got) != 1 {
		t.Fatalf("expected one adapter after duplicate add, got %#v", got)
	}
	if got[0].Location != location {
		t.Fatalf("adapter location=%q, want %q", got[0].Location, location)
	}
}

func TestMeshAdapterConfigHandler_PostUnavailableAdapterReturnsServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	h.config.AdapterTracker = helpers.NewAdaptersTracker(nil)
	pref := models.NewDefaultPreference()
	provider := newAdapterConfigSpyProvider()

	req := newAdapterConfigRequest(t, http.MethodPost, unusedLocalAddress(t))
	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.MeshAdapterConfigHandler(rec, req, pref, &models.User{UserId: "user-1"}, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d (body=%q)", http.StatusInternalServerError, rec.Code, rec.Body.String())
	}
	if provider.recorded.Load() != 0 {
		t.Fatalf("RecordPreferences called %d times, want 0", provider.recorded.Load())
	}
	if len(pref.MeshAdapters) != 0 {
		t.Fatalf("preference adapters=%#v, want empty", pref.MeshAdapters)
	}
}

func TestMeshAdapterConfigHandler_DeleteRemovesAdapter(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	h.config.AdapterTracker = helpers.NewAdaptersTracker(nil)
	location := "meshery-istio:10000"
	pref := models.NewDefaultPreference()
	pref.MeshAdapters = []*models.Adapter{
		{Location: location, Name: "meshery-istio"},
		{Location: "meshery-linkerd:10001", Name: "meshery-linkerd"},
	}
	provider := newAdapterConfigSpyProvider()

	req := httptest.NewRequest(http.MethodDelete, "/api/system/adapter/manage?adapter="+url.QueryEscape(location), nil)
	rec := httptest.NewRecorder()

	h.MeshAdapterConfigHandler(rec, req, pref, &models.User{UserId: "user-1"}, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body=%q)", http.StatusOK, rec.Code, rec.Body.String())
	}
	if provider.recorded.Load() != 1 {
		t.Fatalf("RecordPreferences called %d times, want 1", provider.recorded.Load())
	}
	got := decodeAdaptersResponse(t, rec)
	if len(got) != 1 || got[0].Location != "meshery-linkerd:10001" {
		t.Fatalf("adapters after delete=%#v, want only meshery-linkerd", got)
	}
	if len(pref.MeshAdapters) != 1 || pref.MeshAdapters[0].Location != "meshery-linkerd:10001" {
		t.Fatalf("preference adapters after delete=%#v, want only meshery-linkerd", pref.MeshAdapters)
	}
}

func TestMeshAdapterConfigHandler_DeleteMissingAdapterReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	h.config.AdapterTracker = helpers.NewAdaptersTracker(nil)
	pref := models.NewDefaultPreference()
	pref.MeshAdapters = []*models.Adapter{{Location: "meshery-istio:10000", Name: "meshery-istio"}}
	provider := newAdapterConfigSpyProvider()

	req := httptest.NewRequest(http.MethodDelete, "/api/system/adapter/manage?adapter=missing:9999", nil)
	rec := httptest.NewRecorder()

	h.MeshAdapterConfigHandler(rec, req, pref, &models.User{UserId: "user-1"}, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d (body=%q)", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if provider.recorded.Load() != 0 {
		t.Fatalf("RecordPreferences called %d times, want 0", provider.recorded.Load())
	}
	if len(pref.MeshAdapters) != 1 || pref.MeshAdapters[0].Location != "meshery-istio:10000" {
		t.Fatalf("preference adapters=%#v, want unchanged adapter list", pref.MeshAdapters)
	}
}

func startFakeMeshService(t *testing.T) (*fakeMeshService, string) {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on loopback: %v", err)
	}

	srv := grpc.NewServer()
	fake := &fakeMeshService{}
	meshes.RegisterMeshServiceServer(srv, fake)

	go func() {
		_ = srv.Serve(lis)
	}()
	t.Cleanup(func() {
		srv.Stop()
		_ = lis.Close()
	})

	return fake, lis.Addr().String()
}

func newAdapterConfigRequest(t *testing.T, method, meshLocationURL string) *http.Request {
	t.Helper()
	form := url.Values{}
	form.Set("meshLocationURL", meshLocationURL)
	req := httptest.NewRequest(method, "/api/system/adapter/manage", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func unusedLocalAddress(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve loopback address: %v", err)
	}
	addr := lis.Addr().String()
	if err := lis.Close(); err != nil {
		t.Fatalf("release loopback address %q: %v", addr, err)
	}
	return addr
}

func decodeAdaptersResponse(t *testing.T, rec *httptest.ResponseRecorder) []*models.Adapter {
	t.Helper()
	var adapters []*models.Adapter
	if err := json.Unmarshal(rec.Body.Bytes(), &adapters); err != nil {
		t.Fatalf("decode adapters response %q: %v", rec.Body.String(), err)
	}
	return adapters
}

func assertSingleAdapter(t *testing.T, adapters []*models.Adapter, location string) {
	t.Helper()
	if len(adapters) != 1 {
		t.Fatalf("expected one adapter, got %#v", adapters)
	}
	if adapters[0].Location != location {
		t.Fatalf("adapter location=%q, want %q", adapters[0].Location, location)
	}
	if adapters[0].Name != "meshery-test-adapter" {
		t.Fatalf("adapter name=%q, want meshery-test-adapter", adapters[0].Name)
	}
	if adapters[0].Version != "v1.2.3" {
		t.Fatalf("adapter version=%q, want v1.2.3", adapters[0].Version)
	}
	if adapters[0].GitCommitSHA != "abc123" {
		t.Fatalf("adapter git SHA=%q, want abc123", adapters[0].GitCommitSHA)
	}
	if len(adapters[0].Ops) != 1 || adapters[0].Ops[0].Key != "install" {
		t.Fatalf("adapter ops=%#v, want install op", adapters[0].Ops)
	}
}
