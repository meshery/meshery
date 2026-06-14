package meshes

import (
	"testing"
)

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	t.Log("Need to run Close() skipping")
	//err := m.Close()
	//if err != nil {
	//	t.Errorf("Close() failed with error: %s", err)
	//}
}
