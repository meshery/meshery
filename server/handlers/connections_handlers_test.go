package handlers

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
)

var saveconn = []models.ConnectionPayload{
	{
		Kind:    "meshery",
		SubType: "management",
		Type:    "platform",
		MetaData: map[string]interface{}{
			"server_id":        "d058accd-290d-4123-9c19-3ee7c0b1cbdb",
			"server_version":   "v0.6.91",
			"server_location":  "https://playground.meshery.io",
			"server_build_sha": "4bf9a717",
		},
		Status: models.CONNECTED,
		CredentialSecret: map[string]interface{}{
			"token": "eyJhY2Nlc3NfdG9rZW4iOiJleUpoYkdjaU9pSlNVekkxTmlJc0ltdHBaQ0k2SW5CMVlteHBZenBvZVdSeVlTNXFkM1F1WVdOalpYTnpMWFJ2YTJWdUlpd2lkSGx3SWpvaVNsZFVJbjAuZXlKaGRXUWlPbHRkTENKamJHbGxiblJmYVdRaU9pSnRaWE5vWlhKNUxXTnNiM1ZrSWl3aVpYaHdJam94TmpnNU56TXlNRE16TENKbGVIUWlPbnQ5TENKcFlYUWlPakUyT0RrMk5EVTJNeklzSW1semN5STZJbWgwZEhBNkx5OXNiMk5oYkdodmMzUTZPVEF3TUM4aUxDSnFkR2tpT2lJek9UVmpOMlV3TnkwME1XTXhMVFF6T1RrdE9ERXhaQzFqWmpKak56ZGlabVZtTVRraUxDSnVZbVlpT2pFMk9EazJORFUyTXpJc0luTmpjQ0k2V3lKdmNHVnVhV1FpTENKdlptWnNhVzVsSWl3aWIyWm1iR2x1WlY5aFkyTmxjM01pWFN3aWMzVmlJam9pWWxkV2VtRkhNV2hqUlVKeldWaHNiR05xVlhWaFZ6ZzlPbUpYVm5waFIxWjVaVk14YW1KSE9URmFRVDA5SW4wLmpVelhmMkRrTi15NVpiWThSVkFuTjJZeXFTS05BT2tCZ0dDNHMxb2dnUFoxQS1CakdNbkFyMTZjRVBLcmVPVENXSzNfaU1DdWctcmltZG9QeGlqWFllX1d5bGpqeU9Ta1RyR0xpY0xubE1WNmhXNndJbXptT0E1VWh2UDVWZ09oeXRDZ1VQU3U0OFRVRFFTUE1NU1VkRnRFTzY2c0tvM3pXeURiV01lbm1rTHp6bGJwaFhFbXJpTDFhcnZkRHBFZV9ndHlWa0xjbktCS3p6U2Q2WWF3eDUzS0NhWGhGa2Q0Y0FVaUFVX2lGZjZnYW9ib0VRNUVfV2ZtRDJJaG9BaEtTWHJNbjdmTkxJTERwOFVXRUk3VGdwUUVFZ1QyY1VXdWY4eXlhV2NxOWpmRW9ia3ZQLXpTa29oXzJYZWlLRmtFRjAydTJ2VGp6Z3VpSmJRRl9DWmtQZGlVcGt1VlU0VEJHZVk5U2h6SWxUeGxYdWNySnp1UDI0YVJOMjl6cnNxN1hUMTktTVd2WVRkb3RCek1Dd3hhOTV2Q1BkdUFOTXhPTzVoYWJ3X2hfVUoyTGwycFFfT1dkeHJUSjFNV2pQR25hdzBFcGtqcVlDV1lSdzVJMnBwZGctWGRJUzM3aEpBbTk0MXpLNFhVanhIZzluTVJxSE53Y1RjWnFYRGU0UUstejBVWTJ0Nm0tNjBTZjZJRDg1anpkeFlxMmVfQ3BFMno3MFhaSmJZT1A1OGthU0Zfekxac1FUSkwyUFBDRmpGOGl2bnFFRzlzLUZqemRSSUxmUi1Tc2VsUjZzWk4tb3JteXlLbWF4MnlIYmljZkxTRjhveWQzWW5QSHk5ZWpYSkY1bVpwd3lPWUpmS1VRNEtOVTZ3QXg2cjljRnVvQ0NSRjJFRnJDdnJTVWE0IiwidG9rZW5fdHlwZSI6ImJlYXJlciIsInJlZnJlc2hfdG9rZW1iOiIteXI0Z2VDeEExWlJPTFJyeTlsa3doUXl5VGR4SzV4clFzNzY1ejBkYUxnLmEzSVdLVGZCVHQ5VzRkZ0Mwdl93U0dZX1ozYTZPSndWYl9jU1h4c1hUalUiLCJleHBpcnkiOiIyMDIzLTA3LTE5VDA3OjMwOjMyLjkxNjozNpswNToqMCA1",
		},
	},
}
var user = models.User{
	ID:        "9cb4356o",
	UserID:    "wer5653o",
	FirstName: "Lee",
	LastName:  "Caltcole",
	AvatarURL: "",
	Provider:  "",
	Email:     "lee@gmail.com",
	Status:    "",
	Bio:       "",
	RoleNames: []string{},
}
var provide models.Provider
var preference *models.Preference

func TestSaveConnections(t *testing.T) {
	t.Run("TestSaveConnection", func(t *testing.T) {
		handlers := Handler{}

		// Convert the connection payload to JSON
		payloadJSON, err := json.Marshal(saveconn)
		if err != nil {
			t.Fatalf("Failed to marshal connection payload to JSON: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatalf("Error in making request with Payload: %v", err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		handlers.SaveConnection(rr, req, preference, &user, provide)
		// Check the response status code
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
		t.Log("Saved Connection Successfully")
	})
}
func TestGetConnections(t *testing.T) {
	t.Run("TestGetConnection", func(t *testing.T) {
		h := Handler{}
		req, err := http.NewRequest(http.MethodGet, "/api/integrations/connections/{connectionKind}", nil)
		if err != nil {
			t.Fatalf("Error in making request with payload: %v", err)
		}
		req.URL.RawQuery = "page=0&pagesize=10&order=created_at+desc"
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		r.HandleFunc("/api/integrations/connections/{connectionKind}", func(w http.ResponseWriter, req *http.Request) {
			h.GetConnections(w, req, preference, &models.User{ID: "user123"}, provide)
		})
		// Serve the request
		r.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}
