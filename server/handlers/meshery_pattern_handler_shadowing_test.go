package handlers

import (
	"context"
	"testing"
	"bytes"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	coreV1 "github.com/meshery/schemas/models/v1alpha1/core"
	"github.com/meshery/meshkit/logger"
	"database/sql"
)

// shadowingMockProvider implementation for testing VerifyAndConvertToDesign
type shadowingMockProvider struct {
	models.Provider // Embed the interface so we don't have to implement all methods
	
	fetchedSourceContent []byte
	fetchError           error
	
	savedPattern         *models.MesheryPattern
	saveError            error
	saveResp             []byte
}

func (m *shadowingMockProvider) GetDesignSourceContent(token, patternID string) ([]byte, error) {
	return m.fetchedSourceContent, m.fetchError
}

func (m *shadowingMockProvider) SaveMesheryPattern(token string, pattern *models.MesheryPattern) ([]byte, error) {
	m.savedPattern = pattern
	return m.saveResp, m.saveError
}

func TestVerifyAndConvertToDesign_VariableShadowing(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	
	h := &Handler{
		log: log,
		// registryManager can be nil if ConvertFileToDesign doesn't crash, 
		// but ConvertFileToDesign requires a valid filename with extension.
	}
	
	patternID, _ := uuid.NewV4()
	pattern := &models.MesheryPattern{
		ID:   &patternID,
		Name: "test-helm-chart.tgz", // needs extension for conversion to work without error
		Type: sql.NullString{String: string(coreV1.HelmChart), Valid: true},
		SourceContent: []byte{}, // empty initially
		PatternFile: "", // empty initially
	}
	
	expectedBytes := []byte("mocked helm chart content")
	
	provider := &shadowingMockProvider{
		fetchedSourceContent: expectedBytes,
		saveResp: []byte(`[]`), // mock successful json response from save
	}
	
	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "mock-token")
	
	// ConvertFileToDesign will likely fail without a valid registry or actual valid helm chart,
	// but the bug was in assigning `mesheryPattern.SourceContent = sourceContent` before `ConvertFileToDesign` is called.
	// So we can still assert that `mesheryPattern.SourceContent` gets populated correctly even if the function returns an error later.
	
	_ = h.VerifyAndConvertToDesign(ctx, pattern, provider)
	
	// Verify that the SourceContent was updated with the fetched bytes.
	// This would fail before the bug fix because of the variable shadowing bug.
	if !bytes.Equal(pattern.SourceContent, expectedBytes) {
		t.Errorf("Expected pattern.SourceContent to be %q, but got %q (variable shadowing bug still present)", expectedBytes, pattern.SourceContent)
	}
}
