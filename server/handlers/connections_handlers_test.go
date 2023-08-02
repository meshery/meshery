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
	Type:    "platform",
	MetaData: map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	},
	Status: models.CONNECTED,
	CredentialSecret: map[string]interface{}{
		"token": "eyJhY2Nlc3NfdG9rZW4iOiJleUpoYkdjaU9pSlNVekkxTmlJc0ltdHBaQ0k2SW5CMVlteHBZenBvZVdSeVlTNXFkM1F1WVdOalpYTnpMWFJ2YTJWdUlpd2lkSGx3SWpvaVNsZFVJbjAuZXlKaGRXUWlPbHRkTENKamJHbGxiblJmYVdRaU9pSnRaWE5vWlhKNUxXTnNiM1ZrSWl3aVpYaHdJam94TmpnNU56TXlNRE16TENKbGVIUWlPbnQ5TENKcFlYUWlPakUyT0RrMk5EVTJNeklzSW1semN5STZJbWgwZEhBNkx5OXNiMk5oYkdodmMzUTZPVEF3TUM4aUxDSnFkR2tpT2lJek9UVmpOMlV3TnkwME1XTXhMVFF6T1RrdE9ERXhaQzFqWmpKak56ZGlabVZtTVRraUxDSnVZbVlpT2pFMk9EazJORFUyTXpJc0luTmpjQ0k2V3lKdmNHVnVhV1FpTENKdlptWnNhVzVsSWl3aWIyWm1iR2x1WlY5aFkyTmxjM01pWFN3aWMzVmlJam9pWWxkV2VtRkhNV2hqUlVKeldWaHNiR05xVlhWaFZ6ZzlPbUpYVm5waFIxWjVaVk14YW1KSE9URmFRVDA5SW4wLmpVelhmMkRrTi15NVpiWThSVkFuTjJZeXFTS05BT2tCZ0dDNHMxb2dnUFoxQS1CakdNbkFyMTZjRVBLcmVPVENXSzNfaU1DdWctcmltZG9QeGlqWFllX1d5bGpqeU9Ta1RyR0xpY0xubE1WNmhXNndJbXptT0E1VWh2UDVWZ09oeXRDZ1VQU3U0OFRVRFFTUE1NU1VkRnRFTzY2c0tvM3pXeURiV01lbm1rTHp6bGJwaFhFbXJpTDFhcnZkRHBFZV9ndHlWa0xjbktCS3p6U2Q2WWF3eDUzS0NhWGhGa2Q0Y0FVaUFVX2lGZjZnYW9ib0VRNUVfV2ZtRDJJaG9BaEtTWHJNbjdmTkxJTERwOFVXRUk3VGdwUUVFZ1QyY1VXdWY4eXlhV2NxOWpmRW9ia3ZQLXpTa29oXzJYZWlLRmtFRjAydTJ2VGp6Z3VpSmJRRl9DWmtQZGlVcGt1VlU0VEJHZVk5U2h6SWxUeGxYdWNySnp1UDI0YVJOMjl6cnNxN1hUMTktTVd2WVRkb3RCek1Dd3hhOTV2Q1BkdUFOTXhPTzVoYWJ3X2hfVUoyTGwycFFfT1dkeHJUSjFNV2pQR25hdzBFcGtqcVlDV1lSdzVJMnBwZGctWGRJUzM3aEpBbTk0MXpLNFhVanhIZzluTVJxSE53Y1RjWnFYRGU0UUstejBVWTJ0Nm0tNjBTZjZJRDg1anpkeFlxMmVfQ3BFMno3MFhaSmJZT1A1OGthU0Zfekxac1FUSkwyUFBDRmpGOGl2bnFFRzlzLUZqemRSSUxmUi1Tc2VsUjZzWk4tb3JteXlLbWF4MnlIYmljZkxTRjhveWQzWW5QSHk5ZWpYSkY1bVpwd3lPWUpmS1VRNEtOVTZ3QXg2cjljRnVvQ0NSRjJFRnJDdnJTVWE0IiwidG9rZW5fdHlwZSI6ImJlYXJlciIsInJlZnJlc2hfdG9rZW1iOiIteXI0Z2VDeEExWlJPTFJyeTlsa3doUXl5VGR4SzV4clFzNzY1ejBkYUxnLmEzSVdLVGZCVHQ5VzRkZ0Mwdl93U0dZX1ozYTZPSndWYl9jU1h4c1hUalUiLCJleHBpcnkiOiIyMDIzLTA3LTE5VDA3OjMwOjMyLjkxNjozNpswNToqMCA1",
	},
}
var provide models.Provider

//Save test
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
