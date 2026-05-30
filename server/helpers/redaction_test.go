package helpers

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestRedactSensitiveValueRedactsSecretsAndPromptContent(t *testing.T) {
	input := map[string]any{
		"provider": "openai",
		"apiKey":   "sk-live-secret",
		"headers": map[string]any{
			"Authorization": "Bearer raw-token",
		},
		"credentialSecret": map[string]any{
			"client-key-data":         "private-key-data",
			"client-certificate-data": "certificate-data",
		},
		"kubeconfig": "apiVersion: v1\nusers:\n- token: raw",
		"messages": []any{
			map[string]any{"role": "user", "content": "deploy my private app"},
		},
	}

	got := RedactSensitiveValue(input)
	redacted, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("RedactSensitiveValue() type = %T, want map[string]any", got)
	}

	if redacted["apiKey"] != RedactedValue {
		t.Fatalf("apiKey = %v, want %q", redacted["apiKey"], RedactedValue)
	}
	if redacted["kubeconfig"] != RedactedValue {
		t.Fatalf("kubeconfig = %v, want %q", redacted["kubeconfig"], RedactedValue)
	}

	headers := redacted["headers"].(map[string]any)
	if headers["Authorization"] != RedactedValue {
		t.Fatalf("Authorization = %v, want %q", headers["Authorization"], RedactedValue)
	}

	credentials := redacted["credentialSecret"].(map[string]any)
	if credentials["client-key-data"] != RedactedValue {
		t.Fatalf("client-key-data = %v, want %q", credentials["client-key-data"], RedactedValue)
	}

	messages := redacted["messages"].([]any)
	message := messages[0].(map[string]any)
	if message["content"] != RedactedValue {
		t.Fatalf("message content = %v, want %q", message["content"], RedactedValue)
	}
	if message["role"] != "user" {
		t.Fatalf("message role = %v, want %q", message["role"], "user")
	}
}

func TestRedactSensitiveValueDoesNotMutateInput(t *testing.T) {
	input := map[string]any{
		"apiKey": "sk-live-secret",
		"nested": map[string]any{
			"prompt": "private prompt",
		},
	}
	want := map[string]any{
		"apiKey": "sk-live-secret",
		"nested": map[string]any{
			"prompt": "private prompt",
		},
	}

	_ = RedactSensitiveValue(input)

	if !reflect.DeepEqual(input, want) {
		t.Fatalf("input mutated = %#v, want %#v", input, want)
	}
}

func TestRedactSensitiveStringRedactsCommonSecretPatterns(t *testing.T) {
	input := strings.Join([]string{
		"Authorization: Bearer raw-token",
		"api_key=sk-live-secret",
		"token: raw-token",
		"-----BEGIN PRIVATE KEY-----raw-----END PRIVATE KEY-----",
	}, "\n")

	got := RedactSensitiveString(input)

	for _, leaked := range []string{"raw-token", "sk-live-secret", "BEGIN PRIVATE KEY"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("RedactSensitiveString() leaked %q in %q", leaked, got)
		}
	}
	if strings.Count(got, RedactedValue) < 4 {
		t.Fatalf("RedactSensitiveString() = %q, want at least 4 redactions", got)
	}
}

func TestRedactSensitiveValueRedactsErrors(t *testing.T) {
	got := RedactSensitiveValue(errors.New("provider failed with token=raw-token"))
	if got != "provider failed with token="+RedactedValue {
		t.Fatalf("RedactSensitiveValue(error) = %v", got)
	}
}
