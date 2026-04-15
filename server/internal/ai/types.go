package ai

// CompleteRequest carries everything a provider needs for a single completion.
type CompleteRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model,omitempty"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	Temperature *float64  `json:"temperature,omitempty"`
}

// CompleteResponse is the normalized output across all providers.
type CompleteResponse struct {
	Content      string     `json:"content"`
	Model        string     `json:"model"`
	FinishReason string     `json:"finish_reason"`
	Usage        TokenUsage `json:"usage"`
}

// Message represents a single turn in the conversation.
type Message struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

// TokenUsage tracks token consumption for a single request.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamChunk is a single piece of a streaming response.
type StreamChunk struct {
	Content      string `json:"content"`
	Done         bool   `json:"done"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// EmbedRequest asks a provider to produce a vector embedding.
type EmbedRequest struct {
	Text  string `json:"text"`
	Model string `json:"model,omitempty"`
}

// EmbedResponse carries the resulting embedding vector.
type EmbedResponse struct {
	Embedding []float64  `json:"embedding"`
	Usage     TokenUsage `json:"usage"`
}

// ProviderInfo describes a registered provider's capabilities.
type ProviderInfo struct {
	Name          string   `json:"name"`
	Models        []string `json:"models,omitempty"`
	SupportsEmbed bool     `json:"supports_embed"`
}
