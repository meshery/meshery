package sql

import (
	"testing"
)

var Sample map[string]interface{}

func TestScan(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	err := Map.Interface(Sample)
	if err != nil {
		t.Errorf("Scan() failed with error: %s", err)
	}
}
