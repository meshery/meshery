package telemetry

import "strings"

// NormalizeBaseURL trims surrounding space, any query or fragment, and a
// trailing slash from a telemetry instance's base URL, so that the API path
// appended to it on every request is always the effective request path.
//
// It lives here rather than in each client package so the Prometheus and
// Grafana clients cannot drift apart on what a stored base URL means.
func NormalizeBaseURL(raw string) string {
	s := strings.TrimSpace(raw)
	if i := strings.IndexAny(s, "?#"); i >= 0 {
		s = s[:i]
	}
	return strings.TrimRight(s, "/")
}
