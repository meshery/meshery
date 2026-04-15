package ai

import (
	"context"
	"io"
)

// LLMProvider defines the contract every inference backend must satisfy.
// Implementations are stateless — no conversation history, no session tracking.
// Each call receives its full context through the request struct.
type LLMProvider interface {
	// Complete sends messages to the model and returns the full response.
	Complete(ctx context.Context, req CompleteRequest) (CompleteResponse, error)

	// Stream sends messages and returns a reader that emits newline-delimited
	// JSON chunks. Callers must close the reader when done.
	Stream(ctx context.Context, req CompleteRequest) (io.ReadCloser, error)

	// Embed returns a vector embedding for the input text.
	Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error)

	// Name returns the provider identifier used for registry lookups.
	Name() string
}
