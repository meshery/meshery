package environments

import (
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/google/uuid"
)

var update = flag.Bool("update", false, "update golden files")

var testConstants = map[string]string{
	"orgID":           "2d2c0b60-076a-4f0a-8a63-de538570a553",
	"environmentName": "test-environment",
}

// isValidUUID validates if a string is a valid UUID
func isValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}

// validateEnvironmentFixture validates fixture files for proper UUID format
func validateEnvironmentFixture(t *testing.T, fixturePath string) {
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("Failed to read fixture: %v", err)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("Failed to parse fixture JSON: %v", err)
	}
	
	// Validate UUIDs in the response
	if environments, ok := response["environments"].([]interface{}); ok {
		for _, env := range environments {
			if envMap, ok := env.(map[string]interface{}); ok {
				if id, ok := envMap["id"].(string); ok && !isValidUUID(id) {
					t.Errorf("Invalid environment ID in fixture: %s", id)
				}
				if orgID, ok := envMap["organization_id"].(string); ok && !isValidUUID(orgID) {
					t.Errorf("Invalid organization ID in fixture: %s", orgID)
				}
				// Check if owner field exists and is not a valid UUID
				if owner, ok := envMap["owner"].(string); ok && !isValidUUID(owner) {
					t.Errorf("Invalid owner UUID in fixture: %s", owner)
				}
			}
		}
	}
}