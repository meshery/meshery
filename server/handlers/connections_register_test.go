package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
)

// The Create Connection wizard sends {"status":"register"} for every
// non-Kubernetes kind. Kinds without a verification step (github,
// artifacthub) must advance to REGISTERED instead of failing the transition.
//
// The connection row is saved up front: ProcessConnectionRegistration does not
// persist it yet (#21813), and the status write after the transition needs it.
func TestProcessConnectionRegistrationRegistersKindsWithoutVerifyAction(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.Credential{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
		ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
		},
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister: &models.ConnectionPersister{DB: &db},
		EventsPersister:     &models.EventsPersister{DB: &db},
		GenericPersister:    &db,
	}

	for _, kind := range []string{"github", "artifacthub"} {
		t.Run(kind, func(t *testing.T) {
			id := core.Uuid(uuid.Must(uuid.NewV4()))
			if _, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
				ID: id, Name: "my-" + kind, Kind: kind, Status: connections.DISCOVERED,
			}); err != nil {
				t.Fatalf("save connection: %v", err)
			}

			body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
				ID:                         &id,
				Kind:                       kind,
				Name:                       "my-" + kind,
				Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
				SkipCredentialVerification: true,
			})
			if err != nil {
				t.Fatalf("marshal registration event: %v", err)
			}
			user := &models.User{ID: uuid.Must(uuid.NewV4())}
			ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
			ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
			req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(body)).WithContext(ctx)
			rec := httptest.NewRecorder()

			h.ProcessConnectionRegistration(rec, req, nil, user, provider)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
			}
			conn, err := provider.ConnectionPersister.GetConnection(id, kind)
			if err != nil {
				t.Fatalf("get connection: %v", err)
			}
			if conn.Status != connections.REGISTERED {
				t.Fatalf("persisted status = %q, want %q", conn.Status, connections.REGISTERED)
			}
		})
	}
}

func TestProcessConnectionRegistrationReturnsErrorWhenRegistryLookupFails(t *testing.T) {
	rm, db := newTestRegistryManager(t)

	// Simulate a registry database failure.
	if err := db.DBClose(); err != nil {
		t.Fatalf("close registry database: %v", err)
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		log:             newTestLogger(t),
		SystemID:        &systemID,
		registryManager: rm,
	}

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	provider := &models.DefaultLocalProvider{}

	body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
		Kind:   "github",
		Status: connectionv1beta3.ConnectionRegistrationEventStatusInitialize,
	})
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/integrations/connections/register",
		bytes.NewReader(body),
	)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"status = %d, want %d; body: %s",
			rec.Code,
			http.StatusInternalServerError,
			rec.Body.String(),
		)
	}
}

func TestProcessConnectionRegistration_KubernetesKindWithoutContextReturnsGetContextError(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.Credential{}, models.K8sContext{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
		ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
		},
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister:        &models.ConnectionPersister{DB: &db},
		EventsPersister:            &models.EventsPersister{DB: &db},
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &db},
		GenericPersister:           &db,
	}

	id := core.Uuid(uuid.Must(uuid.NewV4()))
	body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         &id,
		Kind:                       "kubernetes",
		Name:                       "my-k8s",
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	})
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	code, _ := resp["code"].(string)
	if code != ErrGetK8sContextsCode {
		t.Fatalf("expected error code %s (ErrGetK8sContexts), got %s; body = %s", ErrGetK8sContextsCode, code, rec.Body.String())
	}
}

type trackerReuseAction struct {
	executed bool
}

func (a *trackerReuseAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *trackerReuseAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	a.executed = true
	return machines.NoOp, nil, nil
}

func (a *trackerReuseAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func TestProcessConnectionRegistration_KubernetesKindPreExistingTrackerInstance(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.Credential{}, models.K8sContext{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister:        &models.ConnectionPersister{DB: &db},
		EventsPersister:            &models.EventsPersister{DB: &db},
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &db},
		GenericPersister:           &db,
	}

	id := core.Uuid(uuid.Must(uuid.NewV4()))
	userUUID := uuid.Must(uuid.NewV4())
	// Pre-populate an existing machine in the tracker
	sm, err := kubernetes.New(id.String(), userUUID, h.log)
	if err != nil {
		t.Fatalf("new kubernetes machine: %v", err)
	}
	sm.CurrentState = machines.DISCOVERED
	// Attach a valid machine context
	sm.Context = &kubernetes.MachineCtx{
		K8sContext: models.K8sContext{
			Name:         "test-k8s",
			Server:       "https://127.0.0.1:1",
			ConnectionID: id.String(),
		},
	}
	mockAction := &trackerReuseAction{}
	sm.States[machines.REGISTERED] = machines.State{
		Action: mockAction,
	}
	tracker.Add(id, sm)

	if _, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID: id, Name: "my-k8s", Kind: "kubernetes", Status: connections.DISCOVERED,
	}); err != nil {
		t.Fatalf("save connection: %v", err)
	}

	body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         &id,
		Kind:                       "kubernetes",
		Name:                       "my-k8s",
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	})
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}
	user := &models.User{ID: userUUID}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	// Since the machine is already tracked, the handler reuses it without looking up
	// K8sContext from the provider (no context row exists in DB).
	// The tracked machine is reused and no 500 is returned.
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !mockAction.executed {
		t.Fatal("expected pre-existing tracked machine to be reused and executed")
	}
}

// TestProcessConnectionRegistration_KubernetesKindFreshRegistrationWithPersistedContext verifies
// that a fresh registration without an existing tracker entry resolves the persisted K8sContext
// (via contextId in metadata or connection ID) and successfully initializes the machine context.
func TestProcessConnectionRegistration_KubernetesKindFreshRegistrationWithPersistedContext(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.Credential{}, models.K8sContext{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister:        &models.ConnectionPersister{DB: &db},
		EventsPersister:            &models.EventsPersister{DB: &db},
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &db},
		GenericPersister:           &db,
	}

	// Persist the K8sContext beforehand (e.g. from kubeconfig discovery/upload)
	ctxID := "persisted-context-id"
	if _, err := provider.MesheryK8sContextPersister.SaveMesheryK8sContext(models.K8sContext{
		ID:     ctxID,
		Name:   "persisted-k8s",
		Server: "https://127.0.0.1:1",
		Cluster: map[string]interface{}{
			"name": "persisted-k8s",
			"cluster": map[string]interface{}{
				"server": "https://127.0.0.1:1",
			},
		},
		Auth: map[string]interface{}{
			"name": "persisted-k8s",
			"user": map[string]interface{}{
				"token": "dummy-token",
			},
		},
	}); err != nil {
		t.Fatalf("save k8s context: %v", err)
	}

	// Fresh registration tracker ID differs from persisted context ID
	trackerID := core.Uuid(uuid.Must(uuid.NewV4()))
	userUUID := uuid.Must(uuid.NewV4())

	if _, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID: trackerID, Name: "my-k8s", Kind: "kubernetes", Status: connections.DISCOVERED,
	}); err != nil {
		t.Fatalf("save connection: %v", err)
	}

	body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         &trackerID,
		Kind:                       "kubernetes",
		Name:                       "my-k8s",
		Metadata:                   core.Map{"contextId": ctxID},
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	})
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}
	user := &models.User{ID: userUUID}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	// Context lookup succeeds using Metadata["contextId"], bypassing ErrGetK8sContexts.
	// Machine initializes and executes RegisterAction, reaching unreachable ping test (ErrSendMachineEventCode).
	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	code, _ := resp["code"].(string)
	if code == ErrGetK8sContextsCode {
		t.Fatalf("unexpected ErrGetK8sContexts error: %s", rec.Body.String())
	}
	if code == machines.ErrAssetMachineCtxCode {
		t.Fatalf("unexpected ErrAssetMachineCtxCode error: %s", rec.Body.String())
	}
	if code != ErrSendMachineEventCode {
		t.Fatalf("expected error code %s (ErrSendMachineEventCode), got %s; body = %s", ErrSendMachineEventCode, code, rec.Body.String())
	}
	// Verify that the machine was added to the tracker
	if _, ok := tracker.Get(trackerID); !ok {
		t.Fatal("expected newly initialized machine to be registered in tracker")
	}
}

// TestProcessConnectionRegistration_KubernetesKindResolvesContextByConnectionIDLink verifies
// that if Metadata does not contain contextId, the context is looked up via its ConnectionID link.
func TestProcessConnectionRegistration_KubernetesKindResolvesContextByConnectionIDLink(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.Credential{}, models.K8sContext{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister:        &models.ConnectionPersister{DB: &db},
		EventsPersister:            &models.EventsPersister{DB: &db},
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &db},
		GenericPersister:           &db,
	}

	trackerID := core.Uuid(uuid.Must(uuid.NewV4()))
	userUUID := uuid.Must(uuid.NewV4())

	// Persist the K8sContext with an MD5-like ID, linked by ConnectionID = trackerID
	if _, err := provider.MesheryK8sContextPersister.SaveMesheryK8sContext(models.K8sContext{
		ID:           "e10d308c2bb4418199bb446f189b418e",
		ConnectionID: trackerID.String(),
		Name:         "linked-k8s",
		Server:       "https://127.0.0.1:1",
		Cluster: map[string]interface{}{
			"name": "linked-k8s",
			"cluster": map[string]interface{}{
				"server": "https://127.0.0.1:1",
			},
		},
		Auth: map[string]interface{}{
			"name": "linked-k8s",
			"user": map[string]interface{}{
				"token": "dummy-token",
			},
		},
	}); err != nil {
		t.Fatalf("save k8s context: %v", err)
	}

	if _, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID: trackerID, Name: "linked-k8s", Kind: "kubernetes", Status: connections.DISCOVERED,
	}); err != nil {
		t.Fatalf("save connection: %v", err)
	}

	// Payload contains NO contextId in metadata; must resolve by trackerID (ConnectionID link)
	body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         &trackerID,
		Kind:                       "kubernetes",
		Name:                       "linked-k8s",
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	})
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}
	user := &models.User{ID: userUUID}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	code, _ := resp["code"].(string)
	if code == ErrGetK8sContextsCode {
		t.Fatalf("unexpected ErrGetK8sContexts error: %s", rec.Body.String())
	}
	if code != ErrSendMachineEventCode {
		t.Fatalf("expected error code %s (ErrSendMachineEventCode), got %s; body = %s", ErrSendMachineEventCode, code, rec.Body.String())
	}
	if _, ok := tracker.Get(trackerID); !ok {
		t.Fatal("expected newly initialized machine to be registered in tracker")
	}
}

// TestDefaultLocalProvider_GetK8sContextDisambiguatesExactFromNormalizedID verifies that
// GetK8sContext prioritizes exact ID match so that an ID with dashes like "foo-bar" does not
// inadvertently select "foobar".
func TestDefaultLocalProvider_GetK8sContextDisambiguatesExactFromNormalizedID(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(models.K8sContext{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	persister := &models.MesheryK8sContextPersister{DB: &db}
	provider := &models.DefaultLocalProvider{
		MesheryK8sContextPersister: persister,
	}

	ctxDashed := models.K8sContext{
		ID:     "foo-bar",
		Name:   "dashed-context",
		Server: "https://10.0.0.1",
	}
	ctxNormalized := models.K8sContext{
		ID:     "foobar",
		Name:   "normalized-context",
		Server: "https://10.0.0.2",
	}
	if err := db.Create(&ctxDashed).Error; err != nil {
		t.Fatalf("create dashed: %v", err)
	}
	if err := db.Create(&ctxNormalized).Error; err != nil {
		t.Fatalf("create normalized: %v", err)
	}

	res, err := provider.GetK8sContext("", "foo-bar")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if res.ID != "foo-bar" || res.Server != "https://10.0.0.1" {
		t.Fatalf("expected exact context foo-bar (https://10.0.0.1), got ID %s (server %s)", res.ID, res.Server)
	}
}

type mockRemoteProviderForRegistration struct {
	*models.DefaultLocalProvider
	receivedLookupID string
}

func (m *mockRemoteProviderForRegistration) GetProviderType() models.ProviderType {
	return models.RemoteProviderType
}

func (m *mockRemoteProviderForRegistration) GetK8sContext(token, connectionID string) (models.K8sContext, error) {
	m.receivedLookupID = connectionID
	return m.DefaultLocalProvider.GetK8sContext(token, connectionID)
}

// TestProcessConnectionRegistration_RemoteProviderUsesConnectionUUIDEvenWithMetadataContextID verifies
// that when the provider is a RemoteProvider, the registration handler uses the connection UUID rather
// than passing a non-UUID contextId from metadata.
func TestProcessConnectionRegistration_RemoteProviderUsesConnectionUUIDEvenWithMetadataContextID(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.Credential{}, models.K8sContext{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}
	localProvider := &models.DefaultLocalProvider{
		ConnectionPersister:        &models.ConnectionPersister{DB: &db},
		EventsPersister:            &models.EventsPersister{DB: &db},
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &db},
		GenericPersister:           &db,
	}
	remoteProvider := &mockRemoteProviderForRegistration{
		DefaultLocalProvider: localProvider,
	}

	trackerID := core.Uuid(uuid.Must(uuid.NewV4()))
	userUUID := uuid.Must(uuid.NewV4())

	// Persist the K8sContext linked by ConnectionID = trackerID
	if _, err := localProvider.MesheryK8sContextPersister.SaveMesheryK8sContext(models.K8sContext{
		ID:           trackerID.String(),
		ConnectionID: trackerID.String(),
		Name:         "remote-k8s",
		Server:       "https://127.0.0.1:1",
		Cluster: map[string]interface{}{
			"name": "remote-k8s",
			"cluster": map[string]interface{}{
				"server": "https://127.0.0.1:1",
			},
		},
		Auth: map[string]interface{}{
			"name": "remote-k8s",
			"user": map[string]interface{}{
				"token": "dummy-token",
			},
		},
	}); err != nil {
		t.Fatalf("save k8s context: %v", err)
	}

	if _, err := localProvider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID: trackerID, Name: "remote-k8s", Kind: "kubernetes", Status: connections.DISCOVERED,
	}); err != nil {
		t.Fatalf("save connection: %v", err)
	}

	// Payload contains a non-UUID contextId in metadata. RemoteProvider requires connection UUID.
	body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         &trackerID,
		Kind:                       "kubernetes",
		Name:                       "remote-k8s",
		Metadata:                   core.Map{"contextId": "non-uuid-custom-context-id"},
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	})
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}
	user := &models.User{ID: userUUID}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, remoteProvider)

	if remoteProvider.receivedLookupID != trackerID.String() {
		t.Fatalf("expected remote provider to receive connection UUID %s, got %s", trackerID.String(), remoteProvider.receivedLookupID)
	}
}

// TestProcessConnectionRegistration_KubernetesKindFullLifecycleFromSaveK8sContext verifies the
// end-to-end fresh registration flow:
// 1. K8s context discovery / upload calls provider.SaveK8sContext, which mints a connection and links k8sContext.ConnectionID.
// 2. The client sends a registration event with the saved connection's ID.
// 3. ProcessConnectionRegistration resolves the context via its ConnectionID link without needing metadata contextId.
func TestProcessConnectionRegistration_KubernetesKindFullLifecycleFromSaveK8sContext(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.Credential{}, models.K8sContext{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister:        &models.ConnectionPersister{DB: &db},
		EventsPersister:            &models.EventsPersister{DB: &db},
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &db},
		GenericPersister:           &db,
	}

	// 1. Fresh context discovered/uploaded via SaveK8sContext
	mesheryID := core.Uuid(uuid.Must(uuid.NewV4()))
	k8sCtx := models.K8sContext{
		Name:              "discovered-cluster",
		Server:            "https://127.0.0.1:1",
		MesheryInstanceID: &mesheryID,
		Cluster: map[string]interface{}{
			"name": "discovered-cluster",
			"cluster": map[string]interface{}{
				"server": "https://127.0.0.1:1",
			},
		},
		Auth: map[string]interface{}{
			"name": "discovered-cluster",
			"user": map[string]interface{}{
				"token": "test-token",
			},
		},
	}
	createdConn, err := provider.SaveK8sContext("", k8sCtx, nil)
	if err != nil {
		t.Fatalf("save k8s context via provider: %v", err)
	}

	// 2. Client initiates registration for the newly saved connection ID
	connID := createdConn.ID
	body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         &connID,
		Kind:                       "kubernetes",
		Name:                       createdConn.Name,
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	})
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	// 3. ProcessConnectionRegistration resolves the context via its ConnectionID link
	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	code, _ := resp["code"].(string)
	if code == ErrGetK8sContextsCode {
		t.Fatalf("unexpected ErrGetK8sContexts: context resolution failed for fresh connection: %s", rec.Body.String())
	}
	if code == machines.ErrAssetMachineCtxCode {
		t.Fatalf("unexpected ErrAssetMachineCtxCode: machine context was not initialized as *kubernetes.MachineCtx: %s", rec.Body.String())
	}
	if code != ErrSendMachineEventCode {
		t.Fatalf("expected ErrSendMachineEventCode (unreachable ping), got %s; body = %s", code, rec.Body.String())
	}
	if _, ok := tracker.Get(connID); !ok {
		t.Fatal("expected newly initialized machine to be registered in tracker")
	}
}


