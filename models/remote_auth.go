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
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

//DoRequest - executes a request and does refreshing automatically
func (l *RemoteProvider) DoRequest(req *http.Request, tokenString string) (*http.Response, error) {
	resp, err := l.doRequestHelper(req, tokenString)
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		logrus.Error("trying after refresh")
		newToken, err := l.refreshToken(tokenString)
		logrus.Info("token refresh successful")
		if err != nil {
			logrus.Errorf("error doing token : %v", err.Error())
			return nil, err
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
		logrus.Errorf("error refreshing token : %v", err.Error())
		return "", err
	}
	r, err := http.Post(l.RemoteProviderURL+"/refresh", "application/json; charset=utf-8", bytes.NewReader(jsonString))
	if err != nil {
		logrus.Errorf("error refreshing token : %v", err.Error())
		return "", err
	}
	defer SafeClose(r.Body)
	var target map[string]string
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		logrus.Errorf("error refreshing token : %v", err.Error())
		return "", err
	}
	l.TokenStore[tokenString] = target[tokenName]
	time.AfterFunc(1*time.Hour, func() {
		logrus.Infof("deleting old ts")
		delete(l.TokenStore, tokenString)
	})
	return target[tokenName], nil
}

func (l *RemoteProvider) doRequestHelper(req *http.Request, tokenString string) (*http.Response, error) {
	token, err := l.DecodeTokenData(tokenString)
	if err != nil {
		logrus.Errorf("error performing the request, %s", err.Error())
		return nil, err
	}
	c := &http.Client{}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token.AccessToken))
	resp, err := c.Do(req)
	if err != nil {
		logrus.Errorf("error performing the request, %s", err.Error())
		return nil, err
	}
	return resp, nil
}

// GetToken - extracts token form a request
func (l *RemoteProvider) GetToken(req *http.Request) (string, error) {
	ck, err := req.Cookie(tokenName)
	if err != nil {
		logrus.Errorf("error in getting the token, %s", err.Error())
		return "", err
	}
	return ck.Value, nil
}

// DecodeTokenData - Decodes a tokenString to a token
func (l *RemoteProvider) DecodeTokenData(tokenStringB64 string) (*oauth2.Token, error) {
	var token oauth2.Token
	// logrus.Debugf("Token string %s", tokenStringB64)
	tokenString, err := base64.RawStdEncoding.DecodeString(tokenStringB64)
	if err != nil {
		logrus.Errorf("token decode error : %s", err.Error())
		return nil, err
	}
	err = json.Unmarshal(tokenString, &token)
	if err != nil {
		logrus.Errorf("token extraction error : %s", err.Error())
		return nil, err
	}
	return &token, nil
}

// UpdateJWKs - Updates Keys to the JWKS
func (l *RemoteProvider) UpdateJWKs() error {
	resp, err := http.Get(l.RemoteProviderURL + "/keys")
	if err != nil {
		logrus.Errorf("error fetching keys from remote provider : %v", err.Error())
		return err
	}
	defer SafeClose(resp.Body)
	jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error reading keys response as JSON : %v", err.Error())
		return err
	}
	jwksJSON := map[string][]map[string]string{}
	if err := json.Unmarshal([]byte(jsonDataFromHTTP), &jwksJSON); err != nil {
		logrus.Errorf("error converting JSON response to map : %v", err.Error())
		return err
	}

	jwks := jwksJSON["keys"]

	if jwks == nil {
		err = fmt.Errorf("key response invalid")
		logrus.Errorf("Key response invalid : %v", err.Error())
		return err
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
	return nil, fmt.Errorf("key not found")
}

// GenerateKey - generate the actual key from the JWK
func (l *RemoteProvider) GenerateKey(jwk JWK) (*rsa.PublicKey, error) {
	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(jwk["n"])
	if err != nil {
		logrus.Fatalf("error decoding JWK, %v", err.Error())
		return nil, err
	}

	e := 0
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if jwk["e"] == "AQAB" || jwk["e"] == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		err = fmt.Errorf("error decoding JWK as big-endian int")
		logrus.Fatalf("%v", err.Error())
		return nil, err
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}

	der, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		logrus.Fatalf("error mashalling PKIX, : %v", err.Error())
		return nil, err
	}

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	var out bytes.Buffer
	if err := pem.Encode(&out, block); err != nil {
		logrus.Fatalf("error encoding jwk to pem, : %v", err.Error())
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(out.Bytes())
}

// VerifyToken - verifies the validity of a tokenstring
func (l *RemoteProvider) VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	dtoken, err := l.DecodeTokenData(tokenString)
	if err != nil {
		logrus.Fatalf("error decoding token : %v", err.Error())
		return nil, err
	}
	tokenString = dtoken.AccessToken
	tokenUP, x, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		logrus.Fatalf("error parsing token (unverified) : %v", err.Error())
		return nil, err
	}
	kid := tokenUP.Header["kid"].(string)

	var jtk map[string]interface{}
	t, _ := base64.RawStdEncoding.DecodeString(x[1])
	if err := json.Unmarshal(t, &jtk); err != nil {
		logrus.Fatalf("error parsing token (unverified) : %v", err.Error())
		return nil, err
	}

	// TODO: Once hydra fixes https://github.com/ory/hydra/issues/1542
	// we should rather configure hydra auth server to remove nbf field in the token
	exp := int64(jtk["exp"].(float64))
	if jwt.TimeFunc().Unix() > exp {
		err := fmt.Errorf("token has expired")
		logrus.Errorf("error validating token : %v", err.Error())
		return nil, err
	}

	keyJSON, err := l.GetJWK(kid)
	if err != nil {
		logrus.Errorf("error fetching JWK corresponding to token : %v", err.Error())
		return nil, err
	}
	key, err := l.GenerateKey(keyJSON)
	if err != nil {
		logrus.Errorf("error generating Key from JWK : %v", err.Error())
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
		logrus.Errorf("error validating token : %v", err.Error())
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error parsing claims")
	}
	return &claims, nil
}
