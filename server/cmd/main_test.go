package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	t.Log("Need to run main() skipping")
}
