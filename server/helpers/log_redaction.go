package helpers

import (
	"fmt"
	"regexp"
)

const LogRedactedValue = "[REDACTED]"

var logSecretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(authorization\s*[:=]\s*bearer\s+)[^\s,;]+`),
	regexp.MustCompile(`(?i)((?:api[-_ ]?key|token|secret|password)\s*[:=]\s*)[^\s,;]+`),
	regexp.MustCompile(`(?is)-----BEGIN [^-]*PRIVATE KEY-----.*?-----END [^-]*PRIVATE KEY-----`),
}

// RedactLogMessage normalizes arbitrary log values into a string with common
// secret patterns removed before callers pass the value to Meshery's logger.
func RedactLogMessage(value any) string {
	if value == nil {
		return ""
	}
	if err, ok := value.(error); ok {
		return redactLogString(err.Error())
	}
	return redactLogString(fmt.Sprint(value))
}

func redactLogString(value string) string {
	for _, pattern := range logSecretPatterns {
		value = pattern.ReplaceAllString(value, "${1}"+LogRedactedValue)
	}
	return value
}
