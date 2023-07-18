package handlers

import (
	"bytes"
	"encoding/json"
	
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/layer5io/meshery/server/models"
)

var testcon = &models.ConnectionPayload{
	Kind:    "meshery",
	SubType: "management",
	Type:    "type-value",
	MetaData: map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	},
	Status: models.CONNECTED,
	CredentialSecret: map[string]interface{}{
		"secretKey1": "secretValue1",
		"secretKey2": true,
	},
}
var provide models.Provider
func TestSaveConnection(t *testing.T) {
	//create handler instance
	handlers := &Handler{}

	// Convert the connection payload to JSON
	payloadJSON, err := json.Marshal(testcon)
	if err != nil {
		t.Fatalf("Failed to marshal connection payload to JSON: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Log("Error in making request with Payload", err)
	}
	///create a repsonse recorder
	rr := httptest.NewRecorder()
	handlers.SaveConnection(rr, req, nil, nil,provide)
	// Check the response status code
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
	t.Log("Saved Connection Successfully")
}
