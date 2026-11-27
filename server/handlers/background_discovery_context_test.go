package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/meshery/meshery/server/models"
)

func TestBackgroundDiscoveryContextOutlivesRequest(t *testing.T) {
	requestCtx, cancelRequest := context.WithCancel(context.WithValue(context.Background(), models.TokenCtxKey, "test-token"))
	ctx, cancel := backgroundDiscoveryContext(requestCtx)
	defer cancel()

	cancelRequest()
	if err := ctx.Err(); err != nil {
		t.Fatalf("background discovery canceled with its request: %v", err)
	}
	if got := ctx.Value(models.TokenCtxKey); got != "test-token" {
		t.Fatalf("token = %v, want test-token", got)
	}
	deadline, ok := ctx.Deadline()
	if !ok || time.Until(deadline) <= 0 || time.Until(deadline) > 2*time.Minute {
		t.Fatalf("unexpected discovery deadline: %v, %v", deadline, ok)
	}
	cancel()
	if ctx.Err() != context.Canceled {
		t.Fatalf("cancel did not stop background discovery: %v", ctx.Err())
	}
}
