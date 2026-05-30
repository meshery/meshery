package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const RedactedValue = "[REDACTED]"

var (
	bearerTokenPattern = regexp.MustCompile(`(?i)(authorization\s*[:=]\s*bearer\s+)[^\s,;]+`)
	apiKeyPattern      = regexp.MustCompile(`(?i)((?:api[-_ ]?key|token|secret|password)\s*[:=]\s*)[^\s,;]+`)
	privateKeyPattern  = regexp.MustCompile(`(?is)-----BEGIN [^-]*PRIVATE KEY-----.*?-----END [^-]*PRIVATE KEY-----`)
)

var sensitiveFieldNames = map[string]struct{}{
	"apikey":                {},
	"authorization":         {},
	"bearertoken":           {},
	"clientcertificatedata": {},
	"clientkeydata":         {},
	"credential":            {},
	"credentialsecret":      {},
	"kubeconfig":            {},
	"password":              {},
	"secret":                {},
	"token":                 {},
	"xapikey":               {},
}

var promptFieldNames = map[string]struct{}{
	"content":       {},
	"prompt":        {},
	"system":        {},
	"systemprompt":  {},
	"sensitiveinput": {},
	"text":          {},
}

// RedactSensitiveValue returns a JSON-shaped copy of value with credentials,
// tokens, kubeconfigs, private keys, and sensitive prompt-like content replaced
// by a stable placeholder.
func RedactSensitiveValue(value any) any {
	return redactSensitiveValue(reflect.ValueOf(value), "")
}

// RedactSensitiveString removes common secret patterns from strings before
// they are written to logs, events, or client-facing error payloads.
func RedactSensitiveString(value string) string {
	value = privateKeyPattern.ReplaceAllString(value, RedactedValue)
	value = bearerTokenPattern.ReplaceAllString(value, "${1}"+RedactedValue)
	value = apiKeyPattern.ReplaceAllString(value, "${1}"+RedactedValue)
	return value
}

func redactSensitiveValue(value reflect.Value, key string) any {
	if !value.IsValid() {
		return nil
	}

	if value.CanInterface() {
		if err, ok := value.Interface().(error); ok {
			return RedactSensitiveString(err.Error())
		}
	}

	for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
		if value.CanInterface() {
			if err, ok := value.Interface().(error); ok {
				return RedactSensitiveString(err.Error())
			}
		}
	}

	if shouldRedactField(key) {
		return RedactedValue
	}

	switch value.Kind() {
	case reflect.String:
		return RedactSensitiveString(value.String())
	case reflect.Map:
		return redactSensitiveMap(value)
	case reflect.Slice, reflect.Array:
		items := make([]any, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			items = append(items, redactSensitiveValue(value.Index(i), key))
		}
		return items
	case reflect.Struct:
		return redactSensitiveStruct(value)
	case reflect.Bool:
		return value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	default:
		if value.CanInterface() {
			return value.Interface()
		}
		return fmt.Sprint(value)
	}
}

func redactSensitiveMap(value reflect.Value) map[string]any {
	redacted := make(map[string]any, value.Len())
	for _, mapKey := range value.MapKeys() {
		key := fmt.Sprint(mapKey.Interface())
		if shouldRedactField(key) {
			redacted[key] = RedactedValue
			continue
		}
		redacted[key] = redactSensitiveValue(value.MapIndex(mapKey), key)
	}
	return redacted
}

func redactSensitiveStruct(value reflect.Value) any {
	if !value.CanInterface() {
		return nil
	}

	encoded, err := json.Marshal(value.Interface())
	if err != nil {
		return fmt.Sprint(value.Interface())
	}

	var decoded any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return fmt.Sprint(value.Interface())
	}

	return RedactSensitiveValue(decoded)
}

func shouldRedactField(key string) bool {
	normalized := normalizeRedactionKey(key)
	if _, ok := sensitiveFieldNames[normalized]; ok {
		return true
	}
	_, ok := promptFieldNames[normalized]
	return ok
}

func normalizeRedactionKey(key string) string {
	replacer := strings.NewReplacer("_", "", "-", "", " ", "", ".", "")
	return strings.ToLower(replacer.Replace(key))
}
