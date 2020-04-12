package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (l *MesheryRemoteProvider) DoRequest(req *http.Request, tokenString string) (*http.Response, error) {
	resp, err := l.doRequestHelper(req, tokenString)
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		logrus.Errorf("Trying after refresh")
		newToken, err := l.refreshToken(tokenString)
		if err != nil {
			logrus.Errorf("error doing token : %v", err.Error())
			return nil, err
		}
		logrus.Debugf("new token %v", newToken)
		return l.doRequestHelper(req, newToken)
	}
	return resp, err
}

func (l *MesheryRemoteProvider) refreshToken(tokenString string) (string, error) {
	bd := map[string]string{
		tokenName: tokenString,
	}
	jsonString, err := json.Marshal(bd)
	if err != nil {
		logrus.Errorf("error refreshing token : %v", err.Error())
		return "", err
	}
	r, err := http.Post("http://localhost:9876/refresh", "application/json; charset=utf-8", bytes.NewReader(jsonString))
	defer r.Body.Close()
	var target map[string]string
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		logrus.Errorf("error refreshing token : %v", err.Error())
		return "", err
	}
	return target[tokenName], nil
}

func (l *MesheryRemoteProvider) doRequestHelper(req *http.Request, tokenString string) (*http.Response, error) {
	token, err := l.DecodeTokenData(tokenString)
	if err != nil {
		logrus.Errorf("Error performing the request, %s", err.Error())
		return nil, err
	}
	c := &http.Client{}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token.AccessToken))
	resp, err := c.Do(req)
	if err != nil {
		logrus.Errorf("Error performing the request, %s", err.Error())
		return nil, err
	}
	return resp, nil
}

func (l *MesheryRemoteProvider) GetToken(req *http.Request) (string, error) {
	ck, err := req.Cookie(tokenName)
	if err != nil {
		logrus.Errorf("Error in getting the token, %s", err.Error())
		return "", err
	}
	return ck.Value, nil
}

func (l *MesheryRemoteProvider) DecodeTokenData(tokenStringB64 string) (*oauth2.Token, error) {
	var token oauth2.Token
	// logrus.Debugf("Token string %s", tokenStringB64)
	tokenString, err := base64.RawStdEncoding.DecodeString(tokenStringB64)
	if err != nil {
		logrus.Errorf("token decode error : %s", err.Error())
		return nil, err
	}
	err = json.Unmarshal(tokenString, &token)
	if err != nil {
		logrus.Errorf("token decode error : %s", err.Error())
		return nil, err
	}
	return &token, nil
}
