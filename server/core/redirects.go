package core

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

func EncodeRefUrl(url url.URL) string {
	refURL := url.String()
	// If the source is "/", and doesn't include any path or param, set refURL as empty string.
	// Even if this isn't handle, it doesn't lead to issues but adds an extra /? after login in the URL.
	if refURL == "" || refURL == "/" {
		return ""
	}
	refURLB64 := base64.RawURLEncoding.EncodeToString([]byte(refURL))
	return refURLB64
}

func DecodeRefURL(refURLB64 string) (string, error) {
	refURLBytes, err := base64.RawURLEncoding.DecodeString(refURLB64)
	if err != nil {
		return "", err
	}
	return string(refURLBytes), nil
}

// IsValidRedirectURL validates that a redirect URL is safe.
// It only allows relative URLs that don't start with // to prevent open redirects.
func IsValidRedirectURL(redirectURL string) bool {
	if redirectURL == "" || redirectURL == "/" {
		return true
	}
	
	// Must be a relative URL (starts with /)
	if !strings.HasPrefix(redirectURL, "/") {
		return false
	}
	
	// Must not be a protocol-relative URL (starts with //)
	if strings.HasPrefix(redirectURL, "//") {
		return false
	}
	
	// Additional validation: check for suspicious patterns
	// Prevent URLs like /%0d%0aLocation:%20http://evil.com (CRLF injection)
	if strings.Contains(redirectURL, "\r") || strings.Contains(redirectURL, "\n") {
		return false
	}
	
	return true
}

// GetRefURLFromRequest retrieves the ref URL from the request query parameters.
// It returns the decoded ref URL and an error if decoding fails.
func GetRefURLFromRequest(req *http.Request) (string, error) {
	encodedRef := req.URL.Query().Get("ref")
	return DecodeRefURL(encodedRef)
}
