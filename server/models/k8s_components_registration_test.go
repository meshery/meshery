// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadK8sMeshModelMetadata_Success verifies that a valid model template JSON file
// is correctly read, parsed into ModelTemplate, and that the underlying file descriptor
// is closed immediately upon return (Fixes #22752).
func TestLoadK8sMeshModelMetadata_Success(t *testing.T) {
	tempDir := t.TempDir()
	validJSON := `{
		"model": "kubernetes",
		"primaryColor": "#326CE5",
		"secondaryColor": "#7aa1f0",
		"shape": "circle"
	}`

	filePath := filepath.Join(tempDir, "model_template_valid.json")
	err := os.WriteFile(filePath, []byte(validJSON), 0600)
	require.NoError(t, err, "failed to write test file")

	// Track open file descriptors around the call (effective on Linux /proc/self/fd)
	beforeFDs, _ := filepath.Glob("/proc/self/fd/*")
	metadata, err := loadK8sMeshModelMetadata(filePath)
	require.NoError(t, err, "loadK8sMeshModelMetadata should succeed for valid JSON")
	afterFDs, _ := filepath.Glob("/proc/self/fd/*")
	if len(beforeFDs) > 0 && len(afterFDs) > 0 {
		assert.Equal(t, len(beforeFDs), len(afterFDs), "file descriptor count should remain unchanged after successful load")
	}

	assert.Equal(t, "#326CE5", metadata.PrimaryColor)
	require.NotNil(t, metadata.SecondaryColor)
	assert.Equal(t, "#7aa1f0", *metadata.SecondaryColor)
	require.NotNil(t, metadata.Shape)
	assert.Equal(t, "circle", string(*metadata.Shape))

	// Verify the file descriptor was closed: on Windows, an open file descriptor
	// prevents removing or renaming the file. Successful deletion confirms closure.
	removeErr := os.Remove(filePath)
	assert.NoError(t, removeErr, "file descriptor should be closed, allowing removal")
}

// TestLoadK8sMeshModelMetadata_FileNotFound verifies that non-existent paths return
// an error gracefully without panic or resource leaks.
func TestLoadK8sMeshModelMetadata_FileNotFound(t *testing.T) {
	nonExistentPath := filepath.Join(t.TempDir(), "does_not_exist.json")

	metadata, err := loadK8sMeshModelMetadata(nonExistentPath)
	assert.Error(t, err, "expected error for non-existent file")
	assert.Empty(t, metadata.PrimaryColor)
	assert.Nil(t, metadata.SecondaryColor)
}

// TestLoadK8sMeshModelMetadata_CorruptedJSON verifies that malformed JSON files return
// an unmarshal error while still properly closing the open file descriptor via defer (Fixes #22752).
func TestLoadK8sMeshModelMetadata_CorruptedJSON(t *testing.T) {
	tempDir := t.TempDir()
	malformedJSON := `{"model": "kubernetes", "primaryColor": ` // incomplete JSON

	filePath := filepath.Join(tempDir, "model_template_corrupt.json")
	err := os.WriteFile(filePath, []byte(malformedJSON), 0600)
	require.NoError(t, err, "failed to write corrupt test file")

	// Track open file descriptors around the call (effective on Linux /proc/self/fd)
	beforeFDs, _ := filepath.Glob("/proc/self/fd/*")
	metadata, err := loadK8sMeshModelMetadata(filePath)
	assert.Error(t, err, "expected error for corrupted JSON content")
	afterFDs, _ := filepath.Glob("/proc/self/fd/*")
	if len(beforeFDs) > 0 && len(afterFDs) > 0 {
		assert.Equal(t, len(beforeFDs), len(afterFDs), "file descriptor should be closed even on unmarshal error")
	}
	assert.Empty(t, metadata.PrimaryColor)

	// Verify deferred close executed despite JSON unmarshal error (effective on Windows file locks):
	removeErr := os.Remove(filePath)
	assert.NoError(t, removeErr, "file descriptor should be closed even on unmarshal error")
}

// TestLoadK8sMeshModelMetadata_LiveTemplate verifies that the actual project-relative
// template file at k8sMeshModelPath can be read and parsed without coupling to volatile values.
func TestLoadK8sMeshModelMetadata_LiveTemplate(t *testing.T) {
	if _, err := os.Stat(k8sMeshModelPath); os.IsNotExist(err) {
		t.Skipf("skipping live template test: %s not found from current directory", k8sMeshModelPath)
	}

	metadata, err := loadK8sMeshModelMetadata(k8sMeshModelPath)
	require.NoError(t, err, "failed to load live model_template.json")
	assert.NotEmpty(t, metadata.PrimaryColor)
	require.NotNil(t, metadata.SecondaryColor)
	assert.NotEmpty(t, *metadata.SecondaryColor)
	require.NotNil(t, metadata.Shape)
	assert.NotEmpty(t, string(*metadata.Shape))
}

// TestK8sMeshModelMetadata_PackageVar verifies that the package-level K8sMeshModelMetadata
// is initialized by init() when the template file is available without coupling to volatile values.
func TestK8sMeshModelMetadata_PackageVar(t *testing.T) {
	if _, err := os.Stat(k8sMeshModelPath); os.IsNotExist(err) {
		t.Skipf("skipping package var test: %s not found from current directory", k8sMeshModelPath)
	}

	assert.NotEmpty(t, K8sMeshModelMetadata.PrimaryColor)
	require.NotNil(t, K8sMeshModelMetadata.SecondaryColor)
	assert.NotEmpty(t, *K8sMeshModelMetadata.SecondaryColor)
	require.NotNil(t, K8sMeshModelMetadata.Shape)
	assert.NotEmpty(t, string(*K8sMeshModelMetadata.Shape))
}

type mockCompRegProvider struct {
	Provider
	mu     sync.Mutex
	events []events.Event
}

func (m *mockCompRegProvider) PersistSystemEvent(event events.Event) (PersistedEvent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	return PersistedEvent{Event: &event}, nil
}

func (m *mockCompRegProvider) getEvents() []events.Event {
	m.mu.Lock()
	defer m.mu.Unlock()
	copied := make([]events.Event, len(m.events))
	copy(copied, m.events)
	return copied
}

func TestRegisterComponents_FailureResetsStatusAndAllowsRetry(t *testing.T) {
	log, err := logger.New("test", logger.Options{Format: logger.SyslogLogFormat, LogLevel: 0})
	require.NoError(t, err)

	cg := NewComponentsRegistrationHelper(log)
	instID := core.Uuid(uuid.Must(uuid.NewV4()))
	connID := uuid.Must(uuid.NewV4()).String()
	ctxID := "test-fail-ctx"
	k8sCtx := &K8sContext{
		ID:                ctxID,
		Name:              "test-fail-cluster",
		Server:            "https://k8s.example.com",
		ConnectionID:      connID,
		MesheryInstanceID: &instID,
		Cluster: map[string]interface{}{
			"name":   "test-cluster",
			"server": "https://k8s.example.com",
		},
		Auth: map[string]interface{}{
			"name": "test-user",
		},
	}

	cg.UpdateContexts([]*K8sContext{k8sCtx})
	mockProv := &mockCompRegProvider{}
	broadcaster := NewBroadcaster("test-broadcaster")
	testUserID := uuid.Must(uuid.NewV4()).String()

	failRan := make(chan struct{})
	failingFunc := func(provider *Provider, ctxt context.Context, config []byte, ctxID string, connectionID string, userID string, MesheryInstanceID core.Uuid, reg *meshmodel.RegistryManager, eb *Broadcast, log logger.Handler, ctxName string) error {
		close(failRan)
		return fmt.Errorf("simulated cluster unreachable error")
	}

	cg.RegisterComponents([]*K8sContext{k8sCtx}, []K8sRegistrationFunction{failingFunc}, nil, broadcaster, mockProv, testUserID, false)

	// Wait for failing registration function to execute
	select {
	case <-failRan:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for failing registration function to run")
	}

	// Wait for background goroutine defer to finish
	require.Eventually(t, func() bool {
		s, _ := cg.GetContextRegistrationStatus(ctxID)
		return s == NotRegistered && len(mockProv.getEvents()) >= 2
	}, 2*time.Second, 10*time.Millisecond)

	status, _ := cg.GetContextRegistrationStatus(ctxID)
	// Status must be reset to NotRegistered on failure to allow retry
	require.Equal(t, NotRegistered, status, "status must be reset to NotRegistered on failure to allow retry")

	// Verify that the failure event was persisted and published with Error severity
	evts := mockProv.getEvents()
	require.GreaterOrEqual(t, len(evts), 2, "must have emitted started and failure events")
	assert.Equal(t, events.Informational, evts[0].Severity)
	assert.Contains(t, evts[0].Description, "started")
	assert.Equal(t, events.Error, evts[1].Severity)
	assert.Contains(t, evts[1].Description, "Failed to register components")
	assert.Contains(t, evts[1].Metadata["error"], "simulated cluster unreachable error")

	// Retry with a successful registration function
	successRan := make(chan struct{})
	successFunc := func(provider *Provider, ctxt context.Context, config []byte, ctxID string, connectionID string, userID string, MesheryInstanceID core.Uuid, reg *meshmodel.RegistryManager, eb *Broadcast, log logger.Handler, ctxName string) error {
		close(successRan)
		return nil
	}

	cg.RegisterComponents([]*K8sContext{k8sCtx}, []K8sRegistrationFunction{successFunc}, nil, broadcaster, mockProv, testUserID, false)

	select {
	case <-successRan:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for retry registration function to run; retry was dropped because status was falsely complete")
	}

	require.Eventually(t, func() bool {
		s, _ := cg.GetContextRegistrationStatus(ctxID)
		return s == RegistrationComplete
	}, 2*time.Second, 10*time.Millisecond, "status must be RegistrationComplete after successful retry")
}

// TestRegisterComponents_NilContextAndInvalidUserID verifies that RegisterComponents
// safely handles nil contexts in the slice and rejects invalid user IDs without panic.
func TestRegisterComponents_NilContextAndInvalidUserID(t *testing.T) {
	log, err := logger.New("test", logger.Options{Format: logger.SyslogLogFormat, LogLevel: 0})
	require.NoError(t, err)

	cg := NewComponentsRegistrationHelper(log)
	mockProv := &mockCompRegProvider{}
	broadcaster := NewBroadcaster("test-broadcaster")

	// 1. Invalid user ID should log error and return without registering or panicking
	cg.RegisterComponents([]*K8sContext{{ID: "ctx-1"}}, nil, nil, broadcaster, mockProv, "invalid-uuid", false)
	assert.Empty(t, mockProv.getEvents(), "no events should be emitted for invalid user ID")

	// 2. Nil context inside slice should be safely skipped without panicking
	validUserID := uuid.Must(uuid.NewV4()).String()
	cg.RegisterComponents([]*K8sContext{nil}, nil, nil, broadcaster, mockProv, validUserID, false)
	assert.Empty(t, mockProv.getEvents(), "no events should be emitted for nil context")
}

// TestComponentsRegistrationHelper_RemoveContext verifies that RemoveContext
// drops the context from the tracking map, allowing a deleted and re-imported
// cluster to register its components afresh.
func TestComponentsRegistrationHelper_RemoveContext(t *testing.T) {
	log, err := logger.New("test", logger.Options{Format: logger.SyslogLogFormat, LogLevel: 0})
	require.NoError(t, err)

	cg := NewComponentsRegistrationHelper(log)
	ctxID := "ctx-remove-test"
	k8sCtx := &K8sContext{ID: ctxID, Name: "test-remove-cluster"}

	cg.UpdateContexts([]*K8sContext{k8sCtx})
	status, exists := cg.GetContextRegistrationStatus(ctxID)
	require.True(t, exists)
	require.Equal(t, NotRegistered, status)

	// Simulate completed registration
	cg.mx.Lock()
	cg.ctxRegStatusMap[ctxID] = RegistrationComplete
	cg.mx.Unlock()

	status, exists = cg.GetContextRegistrationStatus(ctxID)
	require.True(t, exists)
	require.Equal(t, RegistrationComplete, status)

	// Remove context upon deletion
	cg.RemoveContext(ctxID)
	_, exists = cg.GetContextRegistrationStatus(ctxID)
	require.False(t, exists, "context must be removed from map")

	// Re-importing the context re-initializes it as NotRegistered
	cg.UpdateContexts([]*K8sContext{k8sCtx})
	status, exists = cg.GetContextRegistrationStatus(ctxID)
	require.True(t, exists)
	require.Equal(t, NotRegistered, status, "re-imported context must be NotRegistered")
}

// TestRegisterComponents_RemovedContextNotResurrected verifies that if a context is
// removed via RemoveContext while registration is in-flight, the completing goroutine
// does not resurrect the context in the status map or overwrite newer registrations.
func TestRegisterComponents_RemovedContextNotResurrected(t *testing.T) {
	log, err := logger.New("test", logger.Options{Format: logger.SyslogLogFormat, LogLevel: 0})
	require.NoError(t, err)

	cg := NewComponentsRegistrationHelper(log)
	mockProv := &mockCompRegProvider{}
	broadcaster := NewBroadcaster("test-broadcaster")
	testUserID := uuid.Must(uuid.NewV4()).String()

	ctxID := "ctx-inflight-remove"
	k8sCtx := &K8sContext{ID: ctxID, Name: "test-inflight-cluster"}
	cg.UpdateContexts([]*K8sContext{k8sCtx})

	inFlight := make(chan struct{})
	continueReg := make(chan struct{})

	regFunc := func(provider *Provider, ctxt context.Context, config []byte, ctxID string, connectionID string, userID string, MesheryInstanceID core.Uuid, reg *meshmodel.RegistryManager, eb *Broadcast, log logger.Handler, ctxName string) error {
		close(inFlight)
		<-continueReg
		return nil
	}

	cg.RegisterComponents([]*K8sContext{k8sCtx}, []K8sRegistrationFunction{regFunc}, nil, broadcaster, mockProv, testUserID, false)

	// Wait until registration function is actively executing
	select {
	case <-inFlight:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for registration function to start")
	}

	// Context is deleted while registration is in flight
	cg.RemoveContext(ctxID)
	_, exists := cg.GetContextRegistrationStatus(ctxID)
	require.False(t, exists, "context must be removed from map")

	// Allow in-flight registration to finish
	close(continueReg)

	// Wait briefly and verify that the context was NOT resurrected into ctxRegStatusMap
	time.Sleep(100 * time.Millisecond)
	_, exists = cg.GetContextRegistrationStatus(ctxID)
	require.False(t, exists, "lingering registration must not resurrect removed context")
}
