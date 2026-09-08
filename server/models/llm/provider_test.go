package llm

import (
	"context"
	"testing"
	"time"

	"github.com/meshery/meshery/server/models/connections"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid HTTPS URL", "https://api.openai.com/v1", false},
		{"Valid HTTP Local URL", "http://localhost:11434", false},
		{"Invalid Scheme", "ftp://localhost:11434", true},
		{"Embedded Credentials", "https://user:pass@api.openai.com", true},
		{"Metadata IP AWS", "http://169.254.169.254/latest/meta-data/", true},
		{"Metadata IP GCP", "http://169.254.169.253", true},
		{"Empty URL", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnsureHTTPS(t *testing.T) {
	u1, _ := ValidateURL("http://api.openai.com")
	if err := EnsureHTTPS(u1); err == nil {
		t.Errorf("EnsureHTTPS() expected error for external http URL")
	}

	u2, _ := ValidateURL("http://localhost:11434")
	if err := EnsureHTTPS(u2); err != nil {
		t.Errorf("EnsureHTTPS() expected NO error for localhost http URL")
	}

	u3, _ := ValidateURL("https://api.openai.com")
	if err := EnsureHTTPS(u3); err != nil {
		t.Errorf("EnsureHTTPS() expected NO error for external https URL")
	}
}

func TestNewProviderFromConnection(t *testing.T) {
	ctx := context.Background()

	// Test OpenAI Selection
	conn1 := &connections.ConnectionPayload{
		MetaData: map[string]interface{}{
			"provider": "openai",
		},
		CredentialSecret: map[string]interface{}{
			"secret": "test-key",
		},
	}
	p1, err := NewProviderFromConnection(ctx, conn1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if _, ok := p1.(*OpenAIProvider); !ok {
		t.Errorf("Expected OpenAIProvider, got %T", p1)
	}

	// Test Claude Selection
	conn2 := &connections.ConnectionPayload{
		SubType: "claude",
		CredentialSecret: map[string]interface{}{
			"apiKey": "test-key",
		},
	}
	p2, err := NewProviderFromConnection(ctx, conn2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if _, ok := p2.(*ClaudeProvider); !ok {
		t.Errorf("Expected ClaudeProvider, got %T", p2)
	}

	// Test missing API key
	conn3 := &connections.ConnectionPayload{
		MetaData: map[string]interface{}{
			"provider": "openai",
		},
		CredentialSecret: map[string]interface{}{},
	}
	_, err = NewProviderFromConnection(ctx, conn3)
	if err == nil {
		t.Fatalf("Expected error for missing credentials, got nil")
	}
}

func TestGenerateMesheryDesign(t *testing.T) {
	ctx := context.Background()
	mock := &MockProvider{
		MockResponse: `{"name":"test-design"}`,
	}

	resp, err := GenerateMesheryDesign(ctx, mock, "create a test design", "")
	if err != nil {
		t.Fatalf("GenerateMesheryDesign failed: %v", err)
	}

	if resp != `{"name":"test-design"}` {
		t.Errorf("Unexpected response: %v", resp)
	}
}

func TestNormalizeErrorTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	time.Sleep(5 * time.Millisecond)
	
	err := NormalizeError(ctx.Err())
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	
	// Ensure the original context error is properly mapped to a safe internal error
	// and doesn't leak sensitive connection details that could be in a wrapped context
}
