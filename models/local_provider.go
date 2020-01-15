package models

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LocalProvider - represents a local provider
type LocalProvider struct {
	*MapPreferencePersister
	SaaSBaseURL     string
	ResultPersister *BitCaskResultsPersister
}

// GetProviderType - Returns ProviderType
func (l *LocalProvider) GetProviderType() ProviderType {
	return LocalProviderType
}

// InitiateLogin - initiates login flow and returns a true to indicate the handler to "return" or false to continue
func (l *LocalProvider) InitiateLogin(w http.ResponseWriter, r *http.Request, fromMiddleWare bool) {
	l.issueSession(w, r, fromMiddleWare)
	return
}

// issueSession issues a cookie session after successful login
func (l *LocalProvider) issueSession(w http.ResponseWriter, req *http.Request, fromMiddleWare bool) {
	// session, _ := l.SessionStore.New(req, l.SessionName)
	// session.Options.Path = "/"
	// user := l.fetchUserDetails()
	// session.Values["user"] = user
	// if err := session.Save(req, w); err != nil {
	// 	logrus.Errorf("unable to save session: %v", err)
	// }
	if !fromMiddleWare {
		returnURL := "/"
		if req.RequestURI != "" {
			returnURL = req.RequestURI
		}
		http.Redirect(w, req, returnURL, http.StatusFound)
	}
}

func (l *LocalProvider) fetchUserDetails() *User {
	return &User{
		UserID:    "meshery",
		FirstName: "Meshery",
		LastName:  "Meshery",
		AvatarURL: "",
	}
}

// GetUserDetails - returns the user details
func (l *LocalProvider) GetUserDetails(req *http.Request) (*User, error) {
	// ensuring session is intact before running load test
	// session, err := l.GetSession(req)
	// if err != nil {
	// 	return nil, err
	// }

	// user, _ := session.Values["user"].(*User)
	// return user, nil

	return l.fetchUserDetails(), nil
}

// GetSession - returns the session
func (l *LocalProvider) GetSession(req *http.Request) (*sessions.Session, error) {
	// session, err := l.SessionStore.Get(req, l.SessionName)
	// if err != nil {
	// 	err = errors.Wrap(err, "Error: unable to get session")
	// 	logrus.Error(err)
	// 	return nil, err
	// }
	// return session, nil
	return &sessions.Session{}, nil
}

// GetProviderToken - returns provider token
func (l *LocalProvider) GetProviderToken(req *http.Request) (string, error) {
	return "", nil
}

// Logout - logout from provider backend
func (l *LocalProvider) Logout(w http.ResponseWriter, req *http.Request) {
	// sess, err := l.SessionStore.Get(req, l.SessionName)
	// if err == nil {
	// 	sess.Options.MaxAge = -1
	// 	_ = sess.Save(req, w)
	// }

	http.Redirect(w, req, "/login", http.StatusFound)
}

// FetchResults - fetches results from provider backend
func (l *LocalProvider) FetchResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		err = errors.Wrapf(err, "unable to parse page number")
		logrus.Error(err)
		return nil, err
	}
	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		err = errors.Wrapf(err, "unable to parse page size")
		logrus.Error(err)
		return nil, err
	}
	return l.ResultPersister.GetResults(pg, pgs)
}

// PublishResults - publishes results to the provider backend syncronously
func (l *LocalProvider) PublishResults(req *http.Request, data []byte) (string, error) {
	if err := l.ResultPersister.WriteResult(data); err != nil {
		return "", err
	}
	user, _ := l.GetUserDetails(req)
	pref, _ := l.ReadFromPersister(user.UserID)
	if !pref.AnonymousPerfResults {
		return "", nil
	}

	bf := bytes.NewBuffer(data)
	saasURL, _ := url.Parse(l.SaaSBaseURL + "/result")
	cReq, _ := http.NewRequest(http.MethodPost, saasURL.String(), bf)
	cReq.Header.Set("X-API-Key", GlobalTokenForAnonymousResults)
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		logrus.Warnf("unable to send results: %v", err)
		return "", nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Warnf("unable to read response body: %v", err)
		return "", nil
	}
	if resp.StatusCode == http.StatusCreated {
		// logrus.Infof("results successfully published to SaaS")
		idMap := map[string]string{}
		if err = json.Unmarshal(bdr, &idMap); err != nil {
			logrus.Warnf("unable to unmarshal body: %v", err)
			return "", nil
		}
		resultID, ok := idMap["id"]
		if ok {
			return resultID, nil
		}
		return "", nil
	}
	logrus.Warnf("error while sending results: %s", bdr)
	return "", nil
}

// PublishMetrics - publishes metrics to the provider backend asyncronously
func (l *LocalProvider) PublishMetrics(_ string, data []byte) error {
	bf := bytes.NewBuffer(data)

	saasURL, _ := url.Parse(l.SaaSBaseURL + "/result/metrics")
	cReq, _ := http.NewRequest(http.MethodPut, saasURL.String(), bf)
	cReq.Header.Set("X-API-Key", GlobalTokenForAnonymousResults)
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		logrus.Warnf("unable to send metrics: %v", err)
		return nil
	}
	if resp.StatusCode == http.StatusOK {
		logrus.Infof("metrics successfully published to SaaS")
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Warnf("unable to read response body: %v", err)
		return nil
	}
	logrus.Warnf("error while sending metrics: %s", bdr)
	return nil
}

// RecordPreferences - records the user preference
func (l *LocalProvider) RecordPreferences(req *http.Request, userID string, data *Preference) error {
	return l.MapPreferencePersister.WriteToPersister(userID, data)
}
