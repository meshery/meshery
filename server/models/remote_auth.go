package models

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// JWK - a type respresting the JSON web Key
type JWK map[string]string

// SafeClose is a helper function help to close the io
func SafeClose(co io.Closer) {
	if cerr := co.Close(); cerr != nil {
		logrus.Error(cerr)
	}
}

// DoRequest - executes a request and does refreshing automatically
func (l *RemoteProvider) DoRequest(req *http.Request, tokenString string) (*http.Response, error) {
	resp, err := l.doRequestHelper(req, tokenString)
	if err != nil {
		return nil, ErrTokenRefresh(err)
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		// Read and close response body before reusing request
		// https://github.com/golang/go/issues/19653#issuecomment-341540384
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		logrus.Warn("trying after refresh")
		newToken, err := l.refreshToken(tokenString)
		logrus.Info("token refresh successful")
		if err != nil {
			return nil, ErrTokenRefresh(err)
		}
		return l.doRequestHelper(req, newToken)
	}
	return resp, err
}

// refreshToken - takes a tokenString and returns a new tokenString
func (l *RemoteProvider) refreshToken(tokenString string) (string, error) {
	l.TokenStoreMut.Lock()
	defer l.TokenStoreMut.Unlock()
	newTokenString := l.TokenStore[tokenString]
	if newTokenString != "" {
		return newTokenString, nil
	}
	bd := map[string]string{
		tokenName: tokenString,
	}
	jsonString, err := json.Marshal(bd)
	if err != nil {
		return "", ErrMarshal(err, "refreshing token")
	}
	r, err := http.Post(l.RemoteProviderURL+"/refresh", "application/json; charset=utf-8", bytes.NewReader(jsonString))
	if err != nil {
		return "", err
	}
	if r.StatusCode == http.StatusInternalServerError {
		return "", ErrTokenRefresh(fmt.Errorf("failed to refresh token: status code 500"))
	}

	defer SafeClose(r.Body)
	var target map[string]string
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return "", err
	}
	l.TokenStore[tokenString] = target[tokenName]
	time.AfterFunc(1*time.Hour, func() {
		logrus.Infof("deleting old token string")
		delete(l.TokenStore, tokenString)
	})
	return target[tokenName], nil
}

func (l *RemoteProvider) doRequestHelper(req *http.Request, tokenString string) (*http.Response, error) {
	token, err := l.DecodeTokenData(tokenString)
	if err != nil {
		return nil, ErrTokenDecode(err)
	}
	c := &http.Client{}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token.AccessToken))
	resp, err := c.Do(req)
	if err != nil {
		return nil, ErrTokenClientCheck(err)
	}
	return resp, nil
}

// GetToken - extracts token form a request
func (l *RemoteProvider) GetToken(req *http.Request) (string, error) {
	ck, err := req.Cookie(tokenName)
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
		return nil, err
	}
	err = json.Unmarshal(tokenString, &token)
	if err != nil {
		return nil, ErrUnmarshal(err, "Token String")
	}
	return &token, nil
}

// UpdateJWKs - Updates Keys to the JWKS
func (l *RemoteProvider) UpdateJWKs() error {
	resp, err := http.Get(l.RemoteProviderURL + "/keys")
	if err != nil {
		return ErrJWKsKeys(err)
	}
	defer SafeClose(resp.Body)
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

	l.Keys = jwks

	return nil
}

// GetJWK - Takes a keyId and returns the JWK
func (l *RemoteProvider) GetJWK(kid string) (JWK, error) {
	for _, key := range l.Keys {
		if key["kid"] == kid {
			return key, nil
		}
	}
	err := l.UpdateJWKs()
	if err != nil {
		return nil, err
	}
	for _, key := range l.Keys {
		if key["kid"] == kid {
			return key, nil
		}
	}
	return nil, ErrNilKeys
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
	kid := tokenUP.Header["kid"].(string)

	var jtk map[string]interface{}
	t, _ := base64.RawStdEncoding.DecodeString(x[1])
	if err := json.Unmarshal(t, &jtk); err != nil {
		return nil, ErrPraseUnverified(err)
	}

	// TODO: Once hydra fixes https://github.com/ory/hydra/issues/1542
	// we should rather configure hydra auth server to remove nbf field in the token
	_, ok := jtk["exp"]
	if ok {
		exp := int64(jtk["exp"].(float64))
		if jwt.TimeFunc().Unix() > exp {
			return nil, ErrTokenExpired
		}
	}

	keyJSON, err := l.GetJWK(kid)
	if err != nil {
		return nil, err
	}
	key, err := l.GenerateKey(keyJSON)
	if err != nil {
		return nil, err
	}

	// Verifies the signature
	tokenParser := jwt.Parser{
		SkipClaimsValidation: true,
	}
	token, err := tokenParser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, ErrTokenPrase(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenClaims
	}
	return &claims, nil
}
