package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	_, err := main()
	if err != nil {
		t.Errorf("main() failed with error: %s", err)
	}
}
