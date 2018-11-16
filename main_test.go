// main is an example web app using Login with Twitter.
package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_aoDashRenderer(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "valid test",
			token: "ad6c84be90e7e16ef9150e0c0d809644956d5df6897b73d2340b3238fda40d9d",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(aoDashRenderer))
	defer server.Close()

	// server.Start()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := server.URL + "?token=" + tt.token

			resp, err := http.DefaultClient.Get(u)
			if err != nil {
				logrus.Errorf("error: %v", err)
				t.Fail()
				return
			}
			defer resp.Body.Close()
			d, _ := ioutil.ReadAll(resp.Body)
			logrus.Infof("Resp body: %s", d)
		})
	}
}
