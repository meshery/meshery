package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (l *MesheryRemoteProvider) DoRequest(req *http.Request, tokenString string) (*http.Response, error) {
	token, err := l.DecodeTokenData(tokenString)
	if err != nil {
		logrus.Errorf("Error performing the request, %s", err.Error())
		return nil, err
	}
	c := &http.Client{}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token.AccessToken))
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
