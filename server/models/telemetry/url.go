package telemetry

import (
	"net/url"
	"strings"
)

// NormalizeBaseURL strips user/query/fragment components and trailing slashes
// from a baseURL string so API path concatenation remains safe and deterministic.
func NormalizeBaseURL(baseURL string) string {
	raw := strings.TrimSpace(baseURL)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return strings.TrimRight(raw, "/")
	}
	parsed.RawQuery = ""
	parsed.ForceQuery = false
	parsed.Fragment = ""
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	if parsed.RawPath != "" {
		parsed.RawPath = strings.TrimRight(parsed.RawPath, "/")
	}
	return strings.TrimRight(parsed.String(), "/")
}
