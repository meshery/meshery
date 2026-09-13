package llm

import (
	"context"
	"fmt"
)

// GenerateMesheryDesign takes a natural language request, interacts with the configured LLM,
// and returns a raw response which should represent a Meshery design.
func GenerateMesheryDesign(ctx context.Context, provider Provider, userPrompt string, schemaContext string) (string, error) {
	sysPrompt := `You are an expert infrastructure as code assistant for Meshery.
Your goal is to output valid Meshery design representations (JSON/YAML) based on the user's request.
Always output structured data, avoid hallucinated components, and adhere to the provided schema context.
Do not wrap your output in markdown code blocks. Output only the valid JSON/YAML structure.`

	req := &GenerateRequest{
		SystemPrompt:  sysPrompt,
		UserPrompt:    userPrompt,
		SchemaContext: schemaContext,
	}

	resp, err := provider.Generate(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.RawResponse) == 0 {
		return "", ErrMalformedResponse(fmt.Errorf("provider returned an empty response"))
	}

	// Future: here we would validate resp.RawResponse against Meshery's schema registry.
	// For now, we return the raw response, which the caller can try to parse via core Meshery parsers.

	return resp.RawResponse, nil
}
