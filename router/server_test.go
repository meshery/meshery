package router

import (
	"testing"
)

r := Router("nil", 0)

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	err := r.Close()
	if err != nil {
		t.Errorf("Close() failed with error: %s", err)
	}
} 
