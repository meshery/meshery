# Authentication Flow Security Analysis

## Executive Summary

This document provides a detailed analysis of the authentication flow in Meshery, identifying security vulnerabilities, bugs, and potential issues in session validation, invalidation, and redirect management.

## Authentication Architecture Overview

Meshery supports two provider types:
1. **DefaultLocalProvider** - Local authentication with ephemeral sessions
2. **RemoteProvider** - OAuth2-based authentication with remote provider integration

### Key Components

- **Middlewares**: `ProviderMiddleware`, `AuthMiddleware`, `SessionInjectorMiddleware`
- **Session Management**: Token-based (JWT) for RemoteProvider, no auth for LocalProvider
- **Token Storage**: Cookie-based with `TokenCookieName` and `ProviderSessionCookieName`
- **Routes**: `/user/login`, `/user/logout`, `/api/user/token`, `/auth/login`, `/auth/redirect`

---

## Critical Security Issues

### 1. **Cookie Security Flags Missing - HIGH SEVERITY**

**Location**: `server/models/remote_auth.go:359-376`

**Issue**: Cookies are created without `Secure` and `SameSite` flags.

```go
func setCookie(w http.ResponseWriter, name, value string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,  // ✓ Good
		Expires:  time.Now().Add(duration),
		// Missing: Secure flag
		// Missing: SameSite flag
	})
}
```

**Impact**: 
- **Missing `Secure` flag**: Cookies can be transmitted over unencrypted HTTP connections, exposing JWT tokens to man-in-the-middle attacks
- **Missing `SameSite` flag**: Application is vulnerable to Cross-Site Request Forgery (CSRF) attacks

**Attack Scenario**:
1. Attacker creates malicious website with form POST to `/api/user/logout`
2. Authenticated user visits malicious site
3. Browser automatically includes cookies with cross-site request
4. User is logged out without consent (CSRF attack)

**Recommendation**:
```go
func setCookie(w http.ResponseWriter, name, value string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,  // Add this
		SameSite: http.SameSiteStrictMode,  // Add this
		Expires:  time.Now().Add(duration),
	})
}
```

---

### 2. **Token Expiration Not Validated in All Cases - HIGH SEVERITY**

**Location**: `server/models/remote_auth.go:232-284`

**Issue**: Token expiration validation has inconsistent error handling and allows infinite JWTs without proper controls.

```go
func (l *RemoteProvider) VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	// ... token parsing ...
	
	_, ok := jtk["exp"]
	if ok {
		exp := int64(jtk["exp"].(float64))
		if time.Now().Unix() > exp {  // ← Only checks if "exp" exists
			return nil, ErrTokenExpired
		}
	}
	// If no "exp" claim, token is considered valid indefinitely
	// ...
}
```

**Impact**:
- Tokens without expiration claims are accepted as valid indefinitely
- Compromised tokens remain valid forever unless explicitly revoked
- No max lifetime enforcement for sessions

**Recommendation**:
- Enforce maximum token lifetime even for tokens without `exp` claim
- Require `exp` claim for all tokens
- Add `nbf` (not before) claim validation

---

### 3. **Session Validation Race Condition - MEDIUM SEVERITY**

**Location**: `server/models/remote_auth.go:737-779`

**Issue**: `GetSession()` has unreachable error handling code and potential race conditions.

```go
func (l *RemoteProvider) GetSession(req *http.Request) error {
	// ... verification logic ...
	
	if err != nil {  // ← This condition is NEVER true at this point
		l.Log.Info("Token validation error : ", err.Error())
		newts, err := l.refreshToken(ts)
		// ... refresh logic ...
	}
	return nil
}
```

The code at line 764-777 is unreachable because:
1. If `introspectToken()` fails (line 758), the function returns immediately (line 760)
2. The `err` variable is always `nil` at line 764

**Impact**:
- Dead code indicates incomplete refactoring or logic error
- Token refresh mechanism may not work as intended
- Potential for stale tokens to be accepted

---

### 4. **Redirect Validation Missing - MEDIUM SEVERITY**

**Location**: Multiple locations including `server/models/remote_provider.go:455-468` and `server/models/remote_provider.go:4021-4038`

**Issue**: No validation of redirect URLs, allowing open redirect vulnerabilities.

```go
func (l *RemoteProvider) TokenHandler(w http.ResponseWriter, r *http.Request, _ bool) {
	// ...
	redirectURL := "/"
	// ...
	refQueryParam := r.URL.Query().Get("ref")
	if refQueryParam != "" {
		redirectURL = refQueryParam  // ← No validation!
	}
	// ...
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
```

**Attack Scenario**:
1. Attacker crafts URL: `/api/user/token?ref=https://evil.com/phishing`
2. After authentication, user is redirected to attacker's site
3. Attacker can capture credentials or perform phishing

**Recommendation**:
```go
// Validate redirect URL is relative or same-origin
func isValidRedirect(redirectURL string) bool {
	if redirectURL == "" || redirectURL == "/" {
		return true
	}
	// Must be relative URL (starts with /)
	if !strings.HasPrefix(redirectURL, "/") {
		return false
	}
	// Must not be protocol-relative (starts with //)
	if strings.HasPrefix(redirectURL, "//") {
		return false
	}
	return true
}
```

---

### 5. **Inconsistent Session Invalidation - MEDIUM SEVERITY**

**Location**: `server/handlers/common_handlers.go:42-63` and `server/models/remote_provider.go:793-861`

**Issue**: Cookie deletion uses `Expires` instead of `MaxAge` for immediate invalidation.

```go
func (h *Handler) LogoutHandler(w http.ResponseWriter, req *http.Request, user *models.User, p models.Provider) {
	// ...
	http.SetCookie(w, &http.Cookie{
		Name:     h.config.ProviderCookieName,
		Value:    p.Name(),
		Expires:  time.Now().Add(-time.Hour),  // ← Should use MaxAge: -1
		Path:     "/",
		HttpOnly: true,
	})
	// ...
}
```

**Impact**:
- Browser compatibility issues - some browsers handle `Expires` inconsistently
- Cookies may not be deleted immediately
- User sessions may persist after logout

**Recommendation**:
```go
http.SetCookie(w, &http.Cookie{
	Name:     h.config.ProviderCookieName,
	Value:    "",
	Path:     "/",
	MaxAge:   -1,  // More reliable for deletion
	HttpOnly: true,
})
```

---

### 6. **Logout Flow Error Handling Issues - MEDIUM SEVERITY**

**Location**: `server/models/remote_provider.go:793-861`

**Issue**: Multiple error handling problems in logout flow:

1. **Session cookie error handling** (line 809-813):
```go
sessionCookie, err := req.Cookie("session_cookie")
if err != nil {
	err = ErrGetSessionCookie(err)
	l.Log.Error(err)
	return err  // ← Logout fails if session cookie missing
}
```
This prevents logout if the session cookie is already missing/expired.

2. **Token revocation happens only on success** (lines 847-853):
```go
if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusOK {
	ck, err := req.Cookie(TokenCookieName)
	if err == nil {
		err = l.revokeToken(ck.Value)
	}
	// ...
}
```
If remote provider returns error, token is NOT revoked locally.

**Impact**:
- Users cannot logout if session cookie is missing
- Failed remote logout prevents local cleanup
- Tokens remain in TokenStore map indefinitely

**Recommendation**:
- Make session cookie optional (log warning but continue)
- Always attempt local cleanup regardless of remote provider response
- Add timeout for TokenStore entries

---

### 7. **Token Store Memory Leak - MEDIUM SEVERITY**

**Location**: `server/models/remote_auth.go:72-106`

**Issue**: `TokenStore` map grows unbounded with only 1-hour delayed cleanup.

```go
func (l *RemoteProvider) refreshToken(tokenString string) (string, error) {
	// ...
	l.TokenStore[tokenString] = target[TokenCookieName]
	time.AfterFunc(1*time.Hour, func() {
		l.Log.Info("deleting old token string")
		delete(l.TokenStore, tokenString)
	})
	return target[TokenCookieName], nil
}
```

**Impact**:
- Memory grows with each token refresh
- If users never logout and refresh frequently, memory consumption increases
- No upper bound on map size
- Goroutines accumulate for cleanup

**Recommendation**:
- Implement LRU cache with size limit
- Use time-based eviction strategy
- Consider using a proper cache library

---

### 8. **Middleware Session Validation Commented Out - LOW SEVERITY**

**Location**: `server/handlers/middlewares.go:166-178`

**Issue**: Critical session validation code is commented out in `SessionInjectorMiddleware`.

```go
// ensuring session is intact
// err := provider.GetSession(req)
// if err != nil {
// 	err1 := provider.Logout(w, req)
// 	if err1 != nil {
// 		logrus.Errorf("Error performing logout: %v", err1.Error())
// 		provider.HandleUnAuthenticated(w, req)
// 		return
// 	}
// 	logrus.Errorf("Error: unable to get session: %v", err)
// 	http.Error(w, "unable to get session", http.StatusUnauthorized)
// 	return
// }
```

**Impact**:
- Sessions are not re-validated in SessionInjectorMiddleware
- Relies solely on earlier AuthMiddleware validation
- Potential for stale session data to be used

---

### 9. **Insufficient Error Information Leakage Prevention - LOW SEVERITY**

**Location**: Multiple error handlers

**Issue**: Error responses may leak sensitive information about authentication state.

**Recommendation**:
- Use generic error messages for authentication failures
- Log detailed errors server-side only
- Avoid exposing JWT validation errors to clients

---

## Additional Findings

### 10. **No Rate Limiting on Authentication Endpoints**

**Issue**: No rate limiting on `/user/login`, `/api/user/token` endpoints.

**Impact**: 
- Vulnerable to brute force attacks
- Vulnerable to token enumeration
- DoS potential

---

### 11. **Inconsistent Redirect URL Encoding**

**Location**: `server/core/redirects.go:9-33`

**Issue**: Base64 encoding is used but validation is minimal.

```go
func DecodeRefURL(refURLB64 string) (string, error) {
	refURLBytes, err := base64.RawURLEncoding.DecodeString(refURLB64)
	if err != nil {
		return "", err
	}
	return string(refURLBytes), nil  // No validation of decoded content
}
```

---

### 12. **Auth Redirect Endpoint Security**

**Location**: `server/router/server.go:413-421`

**Issue**: `/auth/redirect` endpoint sets cookies from query parameters without validation.

```go
gMux.HandleFunc("/auth/redirect", func(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	http.SetCookie(w, &http.Cookie{
		Name:     models.TokenCookieName,
		Value:    token,  // ← No validation!
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
```

**Impact**: Potential for token injection if this endpoint is reachable.

---

## Recommended Security Improvements

### Immediate (Critical)

1. **Add Secure and SameSite flags to all cookies**
2. **Validate and restrict redirect URLs**
3. **Fix logout flow to always perform local cleanup**
4. **Remove dead code from GetSession()**

### Short-term (High Priority)

5. **Implement token lifetime limits**
6. **Fix token store memory management**
7. **Add rate limiting to auth endpoints**
8. **Validate tokens on /auth/redirect endpoint**

### Long-term (Best Practices)

9. **Add CSRF token mechanism**
10. **Implement session fingerprinting**
11. **Add security headers (CSP, HSTS, etc.)**
12. **Add audit logging for authentication events**
13. **Consider implementing refresh token rotation**

---

## Testing Recommendations

1. **Add tests for cookie security flags**
2. **Add tests for redirect validation**
3. **Add tests for logout with missing cookies**
4. **Add tests for token expiration edge cases**
5. **Add CSRF tests**
6. **Add open redirect tests**

---

## Conclusion

The authentication flow has several security vulnerabilities ranging from HIGH to LOW severity. The most critical issues are:

1. Missing cookie security flags (CSRF and MitM vulnerabilities)
2. Open redirect vulnerability
3. Incomplete logout flow
4. Token store memory leak

These issues should be addressed prioritizing by severity to improve the security posture of Meshery's authentication system.
