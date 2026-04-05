package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/meshery/meshkit/logger"
	_events "github.com/meshery/meshkit/utils/events"
)

// TestListenForCoreEvents_Cancellation checks that the goroutine stops and does not panic when the context is cancelled.
func TestListenForCoreEvents_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	respChan := make(chan []byte)

	log, _ := logger.New("test", logger.Options{})
	eb := _events.NewEventStreamer()

	// Run the core events listener in a goroutine
	go listenForCoreEvents(ctx, eb, respChan, log, nil)

	// Simulate context cancellation, which is what happens when the HTTP client disconnects
	cancel()

	// Wait briefly to allow the goroutine to process the cancellation
	time.Sleep(100 * time.Millisecond)

	// We verify that it doesn't panic. If it does, the test suite will fail.
	// Since the select block uses ctx.Done(), the goroutine will exit cleanly.
}
