package ai

import (
	"context"
	"io"
)

// StreamReader yields streamed completion chunks from a provider.
type StreamReader interface {
	Recv() (StreamChunk, error)
}

// LLMProvider defines the minimal provider contract for AI backends.
type LLMProvider interface {
	Complete(ctx context.Context, req CompleteRequest) (CompleteResponse, error)
	Stream(ctx context.Context, req StreamRequest) (StreamReader, error)
	Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error)
}

// DrainStream consumes a stream until completion and returns the collected chunks.
func DrainStream(reader StreamReader) ([]StreamChunk, error) {
	if reader == nil {
		return nil, nil
	}

	chunks := make([]StreamChunk, 0)
	for {
		chunk, err := reader.Recv()
		if err != nil {
			if err == io.EOF {
				return chunks, nil
			}

			return nil, err
		}

		chunks = append(chunks, chunk)
	}
}
