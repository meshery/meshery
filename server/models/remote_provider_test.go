package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
)

func TestAnonymousFlowResponse_UnmarshalJSON_StringUUID(t *testing.T) {
	// Test case 1: user_id as string UUID (expected format)
	jsonData := `{
		"access_token": "test_token_123",
		"user_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	var response AnonymousFlowResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal string UUID format: %v", err)
	}

	if response.AccessToken != "test_token_123" {
		t.Errorf("Expected access_token to be 'test_token_123', got '%s'", response.AccessToken)
	}

	expectedUUID, _ := uuid.FromString("550e8400-e29b-41d4-a716-446655440000")
	if response.UserID != expectedUUID {
		t.Errorf("Expected user_id to be %s, got %s", expectedUUID.String(), response.UserID.String())
	}
}

func TestAnonymousFlowResponse_UnmarshalJSON_ByteArray(t *testing.T) {
	// Test case 2: user_id as byte array [16]uint8 (Go array format from Layer5 Cloud)
	jsonData := `{
		"access_token": "test_token_456",
		"user_id": [85, 14, 132, 0, 226, 155, 65, 212, 167, 22, 68, 102, 85, 68, 0, 0]
	}`

	var response AnonymousFlowResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal byte array format: %v", err)
	}

	if response.AccessToken != "test_token_456" {
		t.Errorf("Expected access_token to be 'test_token_456', got '%s'", response.AccessToken)
	}

	expectedUUID, _ := uuid.FromString("550e8400-e29b-41d4-a716-446655440000")
	if response.UserID != expectedUUID {
		t.Errorf("Expected user_id to be %s, got %s", expectedUUID.String(), response.UserID.String())
	}
}

func TestAnonymousFlowResponse_UnmarshalJSON_NoUserID(t *testing.T) {
	// Test case 3: Missing user_id (should be allowed since it's omitempty)
	jsonData := `{
		"access_token": "test_token_789"
	}`

	var response AnonymousFlowResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON without user_id: %v", err)
	}

	if response.AccessToken != "test_token_789" {
		t.Errorf("Expected access_token to be 'test_token_789', got '%s'", response.AccessToken)
	}

	if response.UserID != uuid.Nil {
		t.Errorf("Expected user_id to be nil UUID, got %s", response.UserID.String())
	}
}

func TestAnonymousFlowResponse_UnmarshalJSON_InvalidFormat(t *testing.T) {
	// Test case 4: Invalid user_id format (not string or byte array)
	jsonData := `{
		"access_token": "test_token_invalid",
		"user_id": 12345
	}`

	var response AnonymousFlowResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err == nil {
		t.Fatal("Expected error when unmarshaling invalid user_id format, but got nil")
	}
}

func TestAnonymousFlowResponse_UnmarshalJSON_InvalidUUIDString(t *testing.T) {
	// Test case 5: Invalid UUID string format
	jsonData := `{
		"access_token": "test_token_invalid_uuid",
		"user_id": "not-a-valid-uuid"
	}`

	var response AnonymousFlowResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err == nil {
		t.Fatal("Expected error when unmarshaling invalid UUID string, but got nil")
	}
}
