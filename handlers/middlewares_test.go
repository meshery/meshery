package handlers

import (
	"testing"
)

func TestAuthMiddleWare(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	_, err := AuthMiddleWare(nil, nil)
	if err != nil {
		t.Errorf("AuthMiddleWare() failed with error: %s", err)
	}
}
