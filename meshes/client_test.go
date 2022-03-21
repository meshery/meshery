package meshes

import (
	"testing"
)

m := MeshClient("nil", "nil")

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	err := m.Close()
	if err != nil {
		t.Errorf("Close() failed with error: %s", err)
	}
}
