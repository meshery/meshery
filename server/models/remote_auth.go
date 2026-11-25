package models

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/tracing"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	// Stores meshery provider related info.
	ProviderCookieName = "meshery-provider"

	// Stores the JWT issued by the remote provider to provide secure access to its API
	TokenCookieName = "token"

	// Stores the remote provider session cookie (identity cookie) to facilitate logout from remote provider as user logs out of Meshery
	ProviderSessionCookieName = "session_cookie"

	// Tracks the number of consecutive auth redirect attempts to detect and break redirect loops
	AuthRetryCookieName = "meshery_auth_retries"

	// Maximum number of auth redirect attempts before breaking the loop
	MaxAuthRetries = 3
)

// JWK - a type respresting the JSON web Key
type JWK map[string]string

// SafeClose is a helper function help to close the io
func SafeClose(co io.Closer, log logger.Handler) {
	if cerr := co.Close(); cerr != nil {
		log.Error(ErrCloseIoReader(cerr))
	}
}

// DoRequest - executes a request and does refreshing automatically
func (l *RemoteProvider) DoRequest(req *http.Request, tokenString string) (*http.Response, error) {
	resp, err := l.doRequestHelper(req, tokenString)
	if err != nil {
		return nil, ErrDoRequest(err, req.Method, req.URL.String())
	}

	if resp.StatusCode == 401 {
		if !l.canRefreshToken(tokenString) {
			return resp, nil
		}

		// Read and close response body before reusing request
		// https://github.com/golang/go/issues/19653#issuecomment-341540384
		_, _ = io.ReadAll(resp.Body)
		if err := resp.Body.Close(); err != nil {
			l.Log.Error(err)
		}
		l.Log.Warn(ErrTokenRetry)
		newToken, err := l.refreshToken(tokenString)
		if err != nil {
			return nil, ErrTokenRefresh(err)
		}
		l.Log.Info("token refresh successful")
		retryReq, err := cloneRequestForRetry(req)
		if err != nil {
			return nil, ErrDoRequest(err, req.Method, req.URL.String())
		}

		resp, err := l.doRequestHelper(retryReq, newToken)
		if err != nil {
			return nil, ErrDoRequest(err, req.Method, req.URL.String())
		}
		return resp, nil
	}
	return resp, err
}

// refreshToken - takes a tokenString and returns a new tokenString
func (l *RemoteProvider) refreshToken(tokenString string) (string, error) {
	l.TokenStoreMut.Lock()
	defer l.TokenStoreMut.Unlock()
	newTokenString := strings.TrimSpace(l.TokenStore[tokenString])
	if newTokenString != "" {
		return newTokenString, nil
	}
	bd := map[string]string{
		TokenCookieName: tokenString,
	}
	jsonString, err := json.Marshal(bd)
	if err != nil {
		return "", ErrMarshal(err, "refreshing token")
	}
	r, err := l.postJSON(l.RemoteProviderURL+"/refresh", "application/json; charset=utf-8", jsonString)
	if err != nil {
		return "", err
	}

	defer SafeClose(r.Body, l.Log)
	if r.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			return "", fmt.Errorf("failed to refresh token. status code %d", r.StatusCode)
		}
		return "", fmt.Errorf("failed to refresh token. status code %d: %s", r.StatusCode, strings.TrimSpace(string(body)))
	}

	var target map[string]string
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return "", err
	}
	newTokenString = strings.TrimSpace(target[TokenCookieName])
	if newTokenString == "" {
		return "", fmt.Errorf("failed to refresh token. response missing %q", TokenCookieName)
	}
	l.TokenStore[tokenString] = newTokenString
	time.AfterFunc(1*time.Hour, func() {
		l.Log.Info("deleting old token string")
		delete(l.TokenStore, tokenString)
	})
	return newTokenString, nil
}

func (l *RemoteProvider) canRefreshToken(tokenString string) bool {
	if strings.TrimSpace(tokenString) == "" {
		return false
	}

	token, err := l.DecodeTokenData(tokenString)
	if err != nil {
		return false
	}

	return strings.TrimSpace(token.RefreshToken) != ""
}

func cloneRequestForRetry(req *http.Request) (*http.Request, error) {
	retryReq := req.Clone(req.Context())

	switch {
	case req.GetBody != nil:
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		retryReq.Body = body
	case req.Body == nil || req.Body == http.NoBody:
		retryReq.Body = http.NoBody
	default:
		return nil, fmt.Errorf("request body is not replayable")
	}

	return retryReq, nil
}

func (l *RemoteProvider) getHTTPClient() *http.Client {
	l.httpClientOnce.Do(func() {
		var transport *http.Transport
		if base, ok := http.DefaultTransport.(*http.Transport); ok {
			transport = base.Clone()
		} else {
			transport = &http.Transport{}
		}
		transport.MaxIdleConns = 100
		transport.MaxIdleConnsPerHost = 20
		transport.IdleConnTimeout = 90 * time.Second
		transport.TLSHandshakeTimeout = 10 * time.Second

		// Do not set a blanket http.Client.Timeout or ResponseHeaderTimeout here:
		// DoRequest shares this pooled transport for /upload and /download transfers,
		// while auth/profile paths enforce remoteAuthRequestTimeout via request context.
		l.httpClient = &http.Client{
			Transport: tracing.NewTransport(transport),
		}
	})
	return l.httpClient
}

type cancelOnCloseReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c *cancelOnCloseReadCloser) Close() error {
	err := c.ReadCloser.Close()
	c.cancel()
	return err
}

func tokenCacheKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func evictExpiredAndCap[V any](m map[string]V, now time.Time, expiryOf func(V) time.Time, maxEntries int) {
	if len(m) < maxEntries {
		return
	}
	for k, v := range m {
		if !now.Before(expiryOf(v)) {
			delete(m, k)
		}
	}
	for k := range m {
		if len(m) < maxEntries {
			break
		}
		delete(m, k)
	}
}

func (l *RemoteProvider) resolveEffectiveToken(token string) string {
	l.TokenStoreMut.Lock()
	defer l.TokenStoreMut.Unlock()
	if refreshed := strings.TrimSpace(l.TokenStore[token]); refreshed != "" {
		return refreshed
	}
	return token
}

func (l *RemoteProvider) tokenExpiry(tokenString string) time.Time {
	dtoken, err := l.DecodeTokenData(tokenString)
	if err != nil || dtoken == nil {
		return time.Time{}
	}
	_, parts, err := new(jwt.Parser).ParseUnverified(dtoken.AccessToken, jwt.MapClaims{})
	if err != nil || len(parts) < 2 {
		return time.Time{}
	}
	rawClaims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(rawClaims, &claims); err != nil {
		return time.Time{}
	}
	if expFloat, ok := claims["exp"].(float64); ok {
		return time.Unix(int64(expFloat), 0)
	}
	return time.Time{}
}

func (l *RemoteProvider) authCacheTTLForToken(tokenString string) time.Duration {
	return authCacheTTLFromExpiry(l.tokenExpiry(tokenString))
}

func authCacheTTLFromExpiry(expTime time.Time) time.Duration {
	ttl := defaultAuthCacheTTL
	if !expTime.IsZero() {
		remaining := time.Until(expTime)
		if remaining <= 0 {
			return 0
		}
		if remaining < ttl {
			ttl = remaining
		}
	}
	return ttl
}

func (l *RemoteProvider) isIntrospectCached(cacheKey string) bool {
	now := time.Now()
	l.authCacheMu.RLock()
	exp, ok := l.introspectCache[cacheKey]
	l.authCacheMu.RUnlock()
	return ok && now.Before(exp)
}

func (l *RemoteProvider) setIntrospectCache(cacheKey string, ttl time.Duration) {
	if ttl <= 0 {
		return
	}
	now := time.Now()
	l.authCacheMu.Lock()
	if l.introspectCache == nil {
		l.introspectCache = make(map[string]time.Time)
	}
	evictExpiredAndCap(l.introspectCache, now, func(exp time.Time) time.Time { return exp }, maxAuthCacheEntries)
	l.introspectCache[cacheKey] = now.Add(ttl)
	l.authCacheMu.Unlock()
}

func (l *RemoteProvider) introspectTokenCached(tokenString string, expTime time.Time) error {
	cacheKey := tokenCacheKey(tokenString)
	if l.isIntrospectCached(cacheKey) {
		return nil
	}

	_, err, _ := l.introspectFlight.Do(cacheKey, func() (interface{}, error) {
		if l.isIntrospectCached(cacheKey) {
			return nil, nil
		}
		if err := l.introspectToken(tokenString); err != nil {
			return nil, err
		}
		l.setIntrospectCache(cacheKey, authCacheTTLFromExpiry(expTime))
		return nil, nil
	})
	return err
}

func (l *RemoteProvider) getCachedUserProfile(cacheKey string) (*User, bool) {
	now := time.Now()
	l.authCacheMu.RLock()
	entry, ok := l.userProfileCache[cacheKey]
	l.authCacheMu.RUnlock()
	if !ok || !now.Before(entry.expiresAt) {
		return nil, false
	}
	userCopy := entry.user
	return &userCopy, true
}

func (l *RemoteProvider) setCachedUserProfile(cacheKey string, user *User, ttl time.Duration) {
	if user == nil || ttl <= 0 {
		return
	}
	now := time.Now()
	l.authCacheMu.Lock()
	if l.userProfileCache == nil {
		l.userProfileCache = make(map[string]cachedUserEntry)
	}
	evictExpiredAndCap(l.userProfileCache, now, func(e cachedUserEntry) time.Time { return e.expiresAt }, maxAuthCacheEntries)
	l.userProfileCache[cacheKey] = cachedUserEntry{
		user:      *user,
		expiresAt: now.Add(ttl),
	}
	l.authCacheMu.Unlock()
}

func (l *RemoteProvider) invalidateTokenCache(tokenString string) {
	if strings.TrimSpace(tokenString) == "" {
		return
	}
	keys := []string{tokenCacheKey(tokenString)}
	if effective := l.resolveEffectiveToken(tokenString); effective != "" && effective != tokenString {
		keys = append(keys, tokenCacheKey(effective))
	}
	l.authCacheMu.Lock()
	for _, k := range keys {
		delete(l.introspectCache, k)
		delete(l.userProfileCache, k)
	}
	l.authCacheMu.Unlock()
}

func (l *RemoteProvider) postJSON(urlStr, contentType string, body []byte) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), remoteAuthRequestTimeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bytes.NewReader(body))
	if err != nil {
		cancel()
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := l.getHTTPClient().Do(req)
	if err != nil {
		cancel()
		return nil, err
	}
	resp.Body = &cancelOnCloseReadCloser{ReadCloser: resp.Body, cancel: cancel}
	return resp, nil
}

func (l *RemoteProvider) doRequestHelper(req *http.Request, token string) (*http.Response, error) {
	c := l.getHTTPClient()
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	if token == GlobalTokenForAnonymousResults {
		req.Header.Set("X-API-Key", token)
	} else {
		// Defensive: drop any inbound X-API-Key copied through proxy paths
		// (e.g. ExtensionProxy). The cloud's static-token fallback rejects
		// anything other than GlobalTokenForAnonymousResults, so leaking a
		// caller-supplied value here would only confuse the auth path.
		req.Header.Del("X-API-Key")
	}
	req.Header.Set("SystemID", viper.GetString("INSTANCE_ID")) // Adds the system id to the header for event tracking
	resp, err := c.Do(req)
	if err != nil {
		return nil, ErrTokenClientCheck(err)
	}
	return resp, nil
}

// GetToken - extracts token form a request
func (l *RemoteProvider) GetToken(req *http.Request) (string, error) {
	ck, err := req.Cookie(TokenCookieName)
	if err != nil {
		return "", ErrGetToken(err)
	}
	return ck.Value, nil
}

// DecodeTokenData - Decodes a tokenString to a token
func (l *RemoteProvider) DecodeTokenData(tokenStringB64 string) (*oauth2.Token, error) {
	var token oauth2.Token
	// logrus.Debugf("Token string %s", tokenStringB64)
	tokenString, err := base64.RawStdEncoding.DecodeString(tokenStringB64)
	if err != nil {
		// Some providers set raw access tokens in cookies instead of base64-encoded payloads.
		return &oauth2.Token{AccessToken: tokenStringB64}, nil
	}
	err = json.Unmarshal(tokenString, &token)
	if err != nil {
		return nil, ErrUnmarshal(err, "Token String")
	}
	return &token, nil
}

// UpdateJWKs - Updates Keys to the JWKS
func (l *RemoteProvider) UpdateJWKs() error {
	ctx, cancel := context.WithTimeout(context.Background(), remoteAuthRequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, l.RemoteProviderURL+"/keys", nil)
	if err != nil {
		return ErrJWKsKeys(err)
	}
	resp, err := l.getHTTPClient().Do(req)
	if err != nil {
		return ErrJWKsKeys(err)
	}
	defer SafeClose(resp.Body, l.Log)
	jsonDataFromHTTP, err := io.ReadAll(resp.Body)
	if err != nil {
		return ErrDataRead(err, "Response Body")
	}
	jwksJSON := map[string][]map[string]string{}
	if err := json.Unmarshal([]byte(jsonDataFromHTTP), &jwksJSON); err != nil {
		return ErrUnmarshal(err, "JWKs Keys")
	}

	jwks := jwksJSON["keys"]

	if jwks == nil {
		return ErrNilJWKs
	}

	l.keysMu.Lock()
	l.Keys = jwks
	l.parsedKeys = make(map[string]*rsa.PublicKey)
	l.keysMu.Unlock()

	return nil
}

func (l *RemoteProvider) findJWK(kid string) (JWK, bool) {
	l.keysMu.RLock()
	defer l.keysMu.RUnlock()
	for _, key := range l.Keys {
		if key["kid"] == kid {
			return key, true
		}
	}
	return nil, false
}

// GetJWK - Takes a keyId and returns the JWK
func (l *RemoteProvider) GetJWK(kid string) (JWK, error) {
	if key, ok := l.findJWK(kid); ok {
		return key, nil
	}
	_, err, _ := l.keysFlight.Do("update-jwks", func() (interface{}, error) {
		if _, ok := l.findJWK(kid); ok {
			return nil, nil
		}
		return nil, l.UpdateJWKs()
	})
	if err != nil {
		return nil, err
	}
	if key, ok := l.findJWK(kid); ok {
		return key, nil
	}
	return nil, ErrNilKeys
}

func (l *RemoteProvider) getParsedKey(kid string) (*rsa.PublicKey, error) {
	l.keysMu.RLock()
	if key, ok := l.parsedKeys[kid]; ok && key != nil {
		l.keysMu.RUnlock()
		return key, nil
	}
	l.keysMu.RUnlock()

	keyJSON, err := l.GetJWK(kid)
	if err != nil {
		return nil, err
	}
	key, err := l.GenerateKey(keyJSON)
	if err != nil {
		return nil, err
	}

	l.storeParsedKeyIfCurrent(kid, keyJSON, key)

	return key, nil
}

func (l *RemoteProvider) storeParsedKeyIfCurrent(kid string, keyJSON JWK, key *rsa.PublicKey) {
	l.keysMu.Lock()
	defer l.keysMu.Unlock()
	if l.parsedKeys == nil {
		l.parsedKeys = make(map[string]*rsa.PublicKey)
	}
	for _, current := range l.Keys {
		if current["kid"] == kid && current["n"] == keyJSON["n"] && current["e"] == keyJSON["e"] {
			l.parsedKeys[kid] = key
			break
		}
	}
}

// GenerateKey - generate the actual key from the JWK
func (l *RemoteProvider) GenerateKey(jwk JWK) (*rsa.PublicKey, error) {
	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(jwk["n"])
	if err != nil {
		return nil, ErrDecodeBase64(err, "JWK")
	}

	e := 0
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if jwk["e"] == "AQAB" || jwk["e"] == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		return nil, ErrDecodeBase64(err, "JWK as big-endian int")
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}

	der, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return nil, ErrMarshalPKIX(err)
	}

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	var out bytes.Buffer
	if err := pem.Encode(&out, block); err != nil {
		return nil, ErrEncodingPEM(err)
	}
	return jwt.ParseRSAPublicKeyFromPEM(out.Bytes())
}

// VerifyToken - verifies the validity of a tokenstring
func (l *RemoteProvider) VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	dtoken, err := l.DecodeTokenData(tokenString)
	if err != nil {
		return nil, ErrTokenDecode(err)
	}
	tokenString = dtoken.AccessToken
	tokenUP, x, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, ErrPraseUnverified(err)
	}
	kid, ok := tokenUP.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, ErrNilKeys
	}

	var jtk map[string]interface{}
	t, _ := base64.RawURLEncoding.DecodeString(x[1])
	if err := json.Unmarshal(t, &jtk); err != nil {
		return nil, ErrPraseUnverified(err)
	}

	// TODO: Once hydra fixes https://github.com/ory/hydra/issues/1542
	// we should rather configure hydra auth server to remove nbf field in the token
	if expVal, ok := jtk["exp"].(float64); ok {
		exp := int64(expVal)
		if time.Now().Unix() > exp {
			return nil, ErrTokenExpired
		}
	}

	key, err := l.getParsedKey(kid)
	if err != nil {
		return nil, err
	}

	// Verifies the signature
	claims := jwt.MapClaims{}
	tokenParser := jwt.NewParser()
	token, err := tokenParser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, ErrTokenPrase(err)
	}
	if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok {
		return &tokenClaims, nil
	}
	if tokenClaims, ok := token.Claims.(*jwt.MapClaims); ok {
		return tokenClaims, nil
	}
	return nil, ErrTokenClaims
}

func (l *RemoteProvider) revokeToken(tokenString string) error {
	jsonData := make(map[string]string)
	token, err := l.DecodeTokenData(tokenString)

	if err != nil {
		return ErrTokenDecode(err)
	}

	jsonData["token"] = token.RefreshToken

	body, err := json.Marshal(jsonData)

	if err != nil {
		return ErrMarshal(err, "refresh token")
	}

	remoteURL, err := url.Parse(fmt.Sprintf("%s/revoke", l.RemoteProviderURL))
	if err != nil {
		l.Log.Error(ErrUrlParse(err))
		return err
	}
	r, err := l.postJSON(remoteURL.String(), "application/json", body)

	if err != nil {
		err = ErrTokenRevoke(err)
		l.Log.Error(err)
		return err
	}
	defer SafeClose(r.Body, l.Log)

	if r.StatusCode != http.StatusOK {
		return ErrTokenRevoke(fmt.Errorf("failed to revoke token: status %d", r.StatusCode))
	}
	l.invalidateTokenCache(tokenString)
	return nil
}

func (l *RemoteProvider) introspectToken(tokenString string) error {
	jsonData := make(map[string]string)
	token, err := l.DecodeTokenData(tokenString)

	if err != nil {
		return ErrTokenDecode(err)
	}

	jsonData["token"] = token.AccessToken
	body, err := json.Marshal(jsonData)

	if err != nil {
		return ErrMarshal(err, "refresh token")
	}

	remoteURL, err := url.Parse(fmt.Sprintf("%s/introspect", l.RemoteProviderURL))
	if err != nil {
		l.Log.Error(ErrUrlParse(err))
		return err
	}
	r, err := l.postJSON(remoteURL.String(), "application/json", body)
	if err != nil {
		err = ErrTokenIntrospect(err)
		l.Log.Error(err)
		return err
	}
	defer SafeClose(r.Body, l.Log)

	if r.StatusCode == http.StatusUnauthorized {
		return ErrTokenIntrospect(fmt.Errorf("unauthorized access: status %d", r.StatusCode))
	}

	if r.StatusCode != http.StatusOK {
		return ErrTokenIntrospect(fmt.Errorf("failed to introspect token: status %d", r.StatusCode))
	}

	return nil
}

func setCookie(w http.ResponseWriter, name, value string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(duration),
	})
}

func unsetCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func (l *RemoteProvider) SetProviderCookie(w http.ResponseWriter, provider string) {
	setCookie(w, ProviderCookieName, provider, l.CookieDuration)
}

func (l *RemoteProvider) UnSetProviderCookie(w http.ResponseWriter) {
	unsetCookie(w, ProviderCookieName)
}

func (l *RemoteProvider) SetJWTCookie(w http.ResponseWriter, token string) {
	setCookie(w, TokenCookieName, token, l.CookieDuration)
}

func (l *RemoteProvider) UnSetJWTCookie(w http.ResponseWriter) {
	unsetCookie(w, TokenCookieName)
}

func (l *RemoteProvider) SetProviderSessionCookie(w http.ResponseWriter, sessionCookie string) {
	setCookie(w, ProviderSessionCookieName, sessionCookie, l.CookieDuration)
}

func (l *RemoteProvider) UnSetProviderSessionCookie(w http.ResponseWriter) {
	unsetCookie(w, ProviderSessionCookieName)
}

// GetAuthRetryCount reads the auth retry counter from the request cookie.
// Returns 0 if the cookie is missing or malformed.
func GetAuthRetryCount(req *http.Request) int {
	ck, err := req.Cookie(AuthRetryCookieName)
	if err != nil {
		return 0
	}
	count := 0
	_, _ = fmt.Sscanf(ck.Value, "%d", &count)
	if count < 0 {
		count = 0
	}
	return count
}

// SetAuthRetryCookie sets the auth retry counter cookie with a short TTL.
// The short TTL ensures stale counters are automatically cleaned up.
func SetAuthRetryCookie(w http.ResponseWriter, count int) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthRetryCookieName,
		Value:    fmt.Sprintf("%d", count),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   120, // 2 minutes — enough to detect loops, auto-cleans after
	})
}

// ClearAuthRetryCookie removes the auth retry counter cookie.
// Called after a successful authentication to reset the counter.
func ClearAuthRetryCookie(w http.ResponseWriter) {
	unsetCookie(w, AuthRetryCookieName)
}
