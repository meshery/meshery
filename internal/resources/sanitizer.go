package resources

import (
	"encoding/json"
	"strings"
)

var sensitiveKeys = []string{
	"kubeconfig", "token", "password", "secret", "credential",
	"api_key", "apikey", "auth_token", "private_key",
}

// SanitizeJSON parses raw JSON and removes any key matching sensitive patterns
func SanitizeJSON(input []byte) ([]byte, error) {
	var raw interface{}
	if err := json.Unmarshal(input, &raw); err != nil {
		return nil, err
	}

	sanitized := sanitizeValue(raw)
	return json.Marshal(sanitized)
}

func sanitizeValue(val interface{}) interface{} {
	switch v := val.(type) {
	case map[string]interface{}:
		cleaned := make(map[string]interface{})
		for key, child := range v {
			if isSensitiveKey(key) {
				continue // Omit sensitive fields entirely
			}
			cleaned[key] = sanitizeValue(child)
		}
		return cleaned
	case []interface{}:
		cleaned := make([]interface{}, len(v))
		for i, child := range v {
			cleaned[i] = sanitizeValue(child)
		}
		return cleaned
	default:
		return v
	}
}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(lower, sensitive) {
			return true
		}
	}
	return false
}
