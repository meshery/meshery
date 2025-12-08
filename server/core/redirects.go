package core

import (
	"encoding/base64"
	"net/http"
	"net/url"
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

// GetRefURLFromRequest retrieves the ref URL from the request query parameters.
// It returns the decoded ref URL and an error if decoding fails.
func GetRefURLFromRequest(req *http.Request) (string, error) {
	encodedRef := req.URL.Query().Get("ref")
	return DecodeRefURL(encodedRef)
}
