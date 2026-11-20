package ai

// Role identifies the author of a message sent to a provider.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message represents a single chat message sent to or received from a provider.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// ToolDefinition captures the provider-agnostic tool metadata for a request.
type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"input_schema,omitempty"`
}

// RequestContext carries request-scoped execution settings for a provider call.
type RequestContext struct {
	Provider     string            `json:"provider,omitempty"`
	Model        string            `json:"model,omitempty"`
	SystemPrompt string            `json:"system_prompt,omitempty"`
	Tools        []ToolDefinition  `json:"tools,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// GenerationOptions contains backend-agnostic generation controls.
type GenerationOptions struct {
	Temperature *float64 `json:"temperature,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
}

// CompleteRequest describes a single-turn completion request.
type CompleteRequest struct {
	Context  RequestContext    `json:"context"`
	Messages []Message         `json:"messages"`
	Options  GenerationOptions `json:"options,omitempty"`
}

// StreamRequest describes a streaming completion request.
type StreamRequest struct {
	Context  RequestContext    `json:"context"`
	Messages []Message         `json:"messages"`
	Options  GenerationOptions `json:"options,omitempty"`
}

// EmbedRequest describes an embeddings request.
type EmbedRequest struct {
	Context RequestContext `json:"context"`
	Text    string         `json:"text"`
}

// TokenUsage tracks provider-reported token counts where available.
type TokenUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
	TotalTokens  int `json:"total_tokens,omitempty"`
}

// CompleteResponse contains the normalized output from a completion request.
type CompleteResponse struct {
	Content      string     `json:"content"`
	FinishReason string     `json:"finish_reason,omitempty"`
	Usage        TokenUsage `json:"usage,omitempty"`
}

// StreamChunk contains one unit of streamed completion output.
type StreamChunk struct {
	Content      string `json:"content,omitempty"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// EmbedResponse contains the embedding vector returned by a provider.
type EmbedResponse struct {
	Vector []float32  `json:"vector"`
	Usage  TokenUsage `json:"usage,omitempty"`
}
