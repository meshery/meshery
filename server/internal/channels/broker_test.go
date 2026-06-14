package channels

import (
	"testing"
)

func TestNewBrokerSubscribeChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	ch := NewBrokerSubscribeChannel()
	if ch == nil {
		t.Error("NewBrokerSubscribeChannel() returned nil.")
	}
}
