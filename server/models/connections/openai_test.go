package connections

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestOpenAICred_SecretNeverMarshaled(t *testing.T) {
	cred := NewOpenAICred("my-openai-conn", "sk-supersecretvalue123")

	b, err := json.Marshal(cred)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if strings.Contains(string(b), "supersecretvalue123") {
		t.Fatalf("secret leaked into marshaled output: %s", b)
	}

	if cred.Secret() != "sk-supersecretvalue123" {
		t.Fatalf("Secret() accessor did not return the original value")
	}
}
