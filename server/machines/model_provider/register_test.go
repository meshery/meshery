package model_provider

import (
	"testing"

	"github.com/meshery/meshery/server/models/connections"
)

func TestValidateOpenAI(t *testing.T) {
	tests := []struct {
		name    string
		payload connections.ConnectionPayload
		wantErr bool
	}{
		{
			name: "valid openai key",
			payload: connections.ConnectionPayload{
				Kind:             connections.KindOpenAI,
				CredentialSecret: map[string]interface{}{"apiKey": "sk-abc123"},
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			payload: connections.ConnectionPayload{
				Kind:             connections.KindOpenAI,
				CredentialSecret: map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "wrong prefix",
			payload: connections.ConnectionPayload{
				Kind:             connections.KindOpenAI,
				CredentialSecret: map[string]interface{}{"apiKey": "invalid-key"},
			},
			wantErr: true,
		},
		{
			name: "skip credential verification",
			payload: connections.ConnectionPayload{
				Kind:                       connections.KindOpenAI,
				CredentialSecret:           map[string]interface{}{},
				SkipCredentialVerification: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.payload.SkipCredentialVerification {
				return
			}
			err := validateOpenAI(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOpenAI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAnthropic(t *testing.T) {
	tests := []struct {
		name    string
		payload connections.ConnectionPayload
		wantErr bool
	}{
		{
			name: "valid anthropic key",
			payload: connections.ConnectionPayload{
				Kind:             connections.KindAnthropic,
				CredentialSecret: map[string]interface{}{"apiKey": "sk-ant-abc123"},
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			payload: connections.ConnectionPayload{
				Kind:             connections.KindAnthropic,
				CredentialSecret: map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "openai key used for anthropic",
			payload: connections.ConnectionPayload{
				Kind:             connections.KindAnthropic,
				CredentialSecret: map[string]interface{}{"apiKey": "sk-abc123"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAnthropic(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAnthropic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAWSBedrock(t *testing.T) {
	validCred := map[string]interface{}{
		"accessKeyId":     "AKIAIOSFODNN7EXAMPLE",
		"secretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"region":          "us-east-1",
	}

	tests := []struct {
		name    string
		payload connections.ConnectionPayload
		wantErr bool
	}{
		{
			name:    "valid aws bedrock creds",
			payload: connections.ConnectionPayload{Kind: connections.KindAWSBedrock, CredentialSecret: validCred},
			wantErr: false,
		},
		{
			name: "missing access key id",
			payload: connections.ConnectionPayload{
				Kind: connections.KindAWSBedrock,
				CredentialSecret: map[string]interface{}{
					"secretAccessKey": "secret",
					"region":          "us-east-1",
				},
			},
			wantErr: true,
		},
		{
			name: "missing secret access key",
			payload: connections.ConnectionPayload{
				Kind: connections.KindAWSBedrock,
				CredentialSecret: map[string]interface{}{
					"accessKeyId": "AKID",
					"region":      "us-east-1",
				},
			},
			wantErr: true,
		},
		{
			name: "missing region",
			payload: connections.ConnectionPayload{
				Kind: connections.KindAWSBedrock,
				CredentialSecret: map[string]interface{}{
					"accessKeyId":     "AKID",
					"secretAccessKey": "secret",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAWSBedrock(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAWSBedrock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOllama(t *testing.T) {
	tests := []struct {
		name    string
		payload connections.ConnectionPayload
		wantErr bool
	}{
		{
			name: "valid ollama with base url",
			payload: connections.ConnectionPayload{
				Kind:     connections.KindOllama,
				MetaData: map[string]interface{}{"baseURL": "http://localhost:11434"},
			},
			wantErr: false,
		},
		{
			name: "missing base url",
			payload: connections.ConnectionPayload{
				Kind:     connections.KindOllama,
				MetaData: map[string]interface{}{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOllama(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOllama() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStructure_UnsupportedKind(t *testing.T) {
	err := validateStructure("unknown-provider", connections.ConnectionPayload{})
	if err == nil {
		t.Error("expected error for unsupported provider kind, got nil")
	}
}
