package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// DefaultLocalProvider - represents a local provider
type DefaultLocalProvider struct {
	*MapPreferencePersister
	SaaSBaseURL     string
	ResultPersister *BitCaskResultsPersister
}

// Name - Returns Provider's friendly name
func (l *DefaultLocalProvider) Name() string {
	return "None"
}

// Description - returns a short description of the provider for display in the Provider UI
func (l *DefaultLocalProvider) Description() string {
	return `Provider: None
	- ephemeral sessions
	- environment setup not saved
	- no performance test result history
	- free use`
}

// GetProviderType - Returns ProviderType
func (l *DefaultLocalProvider) GetProviderType() ProviderType {
	return LocalProviderType
}

// GetProviderProperties - Returns all the provider properties required
func (l *DefaultLocalProvider) GetProviderProperties() ProviderProperties {
	var result ProviderProperties
	result.ProviderType = l.GetProviderType()
	result.DisplayName = l.Name()
	result.Description = l.Description()
	result.Capabilities = make([]Capability, 0)
	return result
}

// InitiateLogin - initiates login flow and returns a true to indicate the handler to "return" or false to continue
func (l *DefaultLocalProvider) InitiateLogin(w http.ResponseWriter, r *http.Request, fromMiddleWare bool) {
	l.issueSession(w, r, fromMiddleWare)
	return
}

// issueSession issues a cookie session after successful login
func (l *DefaultLocalProvider) issueSession(w http.ResponseWriter, req *http.Request, fromMiddleWare bool) {
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

func (l *DefaultLocalProvider) fetchUserDetails() *User {
	return &User{
		UserID:    "meshery",
		FirstName: "Meshery",
		LastName:  "Meshery",
		AvatarURL: "",
	}
}

// GetUserDetails - returns the user details
func (l *DefaultLocalProvider) GetUserDetails(req *http.Request) (*User, error) {
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
func (l *DefaultLocalProvider) GetSession(req *http.Request) (*sessions.Session, error) {
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
func (l *DefaultLocalProvider) GetProviderToken(req *http.Request) (string, error) {
	return "", nil
}

// Logout - logout from provider backend
func (l *DefaultLocalProvider) Logout(w http.ResponseWriter, req *http.Request) {
	// sess, err := l.SessionStore.Get(req, l.SessionName)
	// if err == nil {
	// 	sess.Options.MaxAge = -1
	// 	_ = sess.Save(req, w)
	// }

	http.Redirect(w, req, "/login", http.StatusFound)
}

// FetchResults - fetches results from provider backend
func (l *DefaultLocalProvider) FetchResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
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

// GetResult - fetches result from provider backend for the given result id
func (l *DefaultLocalProvider) GetResult(req *http.Request, resultID uuid.UUID) (*MesheryResult, error) {
	// key := uuid.FromStringOrNil(resultID)
	if resultID == uuid.Nil {
		return nil, fmt.Errorf("given resultID is not valid")
	}
	return l.ResultPersister.GetResult(resultID)
}

// PublishResults - publishes results to the provider backend syncronously
func (l *DefaultLocalProvider) PublishResults(req *http.Request, result *MesheryResult) (string, error) {
	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery result for shipping"))
		return "", err
	}
	user, _ := l.GetUserDetails(req)
	pref, _ := l.ReadFromPersister(user.UserID)
	if !pref.AnonymousPerfResults {
		return "", nil
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	resultID, _ := l.shipResults(req, data)

	key := uuid.FromStringOrNil(resultID)
	logrus.Debugf("key: %s, is nil: %t", key.String(), (key == uuid.Nil))
	if key == uuid.Nil {
		key, _ = uuid.NewV4()
		result.ID = key
		data, err = json.Marshal(result)
		if err != nil {
			logrus.Error(errors.Wrap(err, "error - unable to marshal meshery result for persisting"))
			return "", err
		}
	}
	if err := l.ResultPersister.WriteResult(key, data); err != nil {
		return "", err
	}

	return key.String(), nil
}

func (l *DefaultLocalProvider) shipResults(req *http.Request, data []byte) (string, error) {
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
func (l *DefaultLocalProvider) PublishMetrics(_ string, result *MesheryResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return err
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
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
func (l *DefaultLocalProvider) RecordPreferences(req *http.Request, userID string, data *Preference) error {
	return l.MapPreferencePersister.WriteToPersister(userID, data)
}
