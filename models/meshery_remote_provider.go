package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MesheryRemoteProvider - represents a local provider
type MesheryRemoteProvider struct {
	*BitCaskPreferencePersister

	SaaSTokenName string
	SaaSBaseURL   string

	SessionName   string
	RefCookieName string

	SessionStore        sessions.Store
	LoginCookieDuration time.Duration

	syncStopChan chan struct{}
	syncChan     chan *userSession
}

type userSession struct {
	token   string
	session *Preference
}

// UserPref - is just use to separate out the user info from preference
type UserPref struct {
	User
	Preferences *Preference `json:"preferences,omitempty"`
}

// Name - Returns Provider's friendly name
func (l *MesheryRemoteProvider) Name() string {
	return "Meshery"
}

// Description - returns a short description of the provider for display in the Provider UI
func (l *MesheryRemoteProvider) Description() string {
	return `Provider: Meshery (default)
	- persistent sessions 
	- save environment setup 
	- retrieve performance test results 
	- free use`
}

// GetProviderType - Returns ProviderType
func (l *MesheryRemoteProvider) GetProviderType() ProviderType {
	return RemoteProviderType
}

// GetProviderProperties - Returns all the provider properties required
func (l *MesheryRemoteProvider) GetProviderProperties() ProviderProperties {
	var result ProviderProperties
	result.ProviderType = l.GetProviderType()
	result.DisplayName = l.Name()
	result.Description = l.Description()
	result.Capabilities = make([]Capability, 0)
	return result
}

// SyncPreferences - used to sync preferences with the remote provider
func (l *MesheryRemoteProvider) SyncPreferences() {
	l.syncStopChan = make(chan struct{})
	l.syncChan = make(chan *userSession, 100)
	go func() {
		for {
			select {
			case uSess := <-l.syncChan:
				l.executePrefSync(uSess.token, uSess.session)
			case <-l.syncStopChan:
				return
			}
		}
	}()
}

// StopSyncPreferences - used to stop sync preferences
func (l *MesheryRemoteProvider) StopSyncPreferences() {
	l.syncStopChan <- struct{}{}
}

func (l *MesheryRemoteProvider) executePrefSync(tokenVal string, sess *Preference) {
	bd, err := json.Marshal(sess)
	if err != nil {
		logrus.Errorf("unable to marshal preference data: %v", err)
		return
	}
	saasURL, _ := url.Parse(l.SaaSBaseURL + "/user/preferences")
	req, _ := http.NewRequest(http.MethodPut, saasURL.String(), bytes.NewReader(bd))
	req.AddCookie(&http.Cookie{
		Name:     l.SaaSTokenName,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		logrus.Errorf("unable to upload user preference data: %v", err)
		return
	}
	if resp.StatusCode != http.StatusCreated {
		logrus.Errorf("unable to upload user preference data, status code received: %d", resp.StatusCode)
	}
}

// InitiateLogin - initiates login flow and returns a true to indicate the handler to "return" or false to continue
func (l *MesheryRemoteProvider) InitiateLogin(w http.ResponseWriter, r *http.Request, _ bool) {
	tu := "http://" + r.Host + r.RequestURI
	token := r.URL.Query().Get(l.SaaSTokenName)
	if token == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     l.RefCookieName,
			Value:    "/",
			Expires:  time.Now().Add(l.LoginCookieDuration),
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, l.SaaSBaseURL+"?source="+base64.URLEncoding.EncodeToString([]byte(tu)), http.StatusFound)
		return
	}
	l.issueSession(w, r)
	return
}

// issueSession issues a cookie session after successful login
func (l *MesheryRemoteProvider) issueSession(w http.ResponseWriter, req *http.Request) {
	var reffURL string
	reffCk, _ := req.Cookie(l.RefCookieName)
	if reffCk != nil {
		reffURL = reffCk.Value
	}
	logrus.Infof("preparing to issue session. retrieved reff url: %s", reffURL)
	if reffURL == "" {
		reffURL = "/"
	}
	// session, err := h.config.SessionStore.New(req, h.config.SessionName)
	session, _ := l.SessionStore.New(req, l.SessionName)
	// if err != nil {
	// 	logrus.Errorf("unable to create session: %v", err)
	// 	http.Error(w, "unable to create session", http.StatusInternalServerError)
	// 	return
	// }
	session.Options.Path = "/"
	token := ""
	for k, va := range req.URL.Query() {
		for _, v := range va {
			if k == l.SaaSTokenName {
				// logrus.Infof("setting user in session: %s", v)
				token = v
				break
			}
		}
	}
	if reffCk != nil && reffCk.Name != "" {
		reffCk.Expires = time.Now().Add(-2 * time.Second)
		http.SetCookie(w, reffCk)
	}
	session.Values[l.SaaSTokenName] = token
	user, err := l.fetchUserDetails(token)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)

	}
	session.Values["user"] = user
	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
	}
	http.Redirect(w, req, reffURL, http.StatusFound)
}

func (l *MesheryRemoteProvider) fetchUserDetails(tokenVal string) (*User, error) {
	saasURL, _ := url.Parse(l.SaaSBaseURL + "/user")
	req, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)
	req.AddCookie(&http.Cookie{
		Name:     l.SaaSTokenName,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		logrus.Errorf("unable to fetch user data: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read body: %v", err)
		return nil, err
	}

	up := &UserPref{}
	err = json.Unmarshal(bd, up)
	if err != nil {
		logrus.Errorf("unable to unmarshal user: %v", err)
		return nil, err
	}

	prefLocal, _ := l.ReadFromPersister(up.UserID)
	if prefLocal == nil || up.Preferences.UpdatedAt.After(prefLocal.UpdatedAt) {
		_ = l.WriteToPersister(up.UserID, up.Preferences)
	}

	logrus.Infof("retrieved user: %v", up.User)
	return &up.User, nil
}

// GetUserDetails - returns the user details
func (l *MesheryRemoteProvider) GetUserDetails(req *http.Request) (*User, error) {
	// ensuring session is intact before running load test
	session, err := l.GetSession(req)
	if err != nil {
		return nil, err
	}

	token, _ := l.GetProviderToken(req)

	user, _ := session.Values["user"].(*User)
	_, _ = l.fetchUserDetails(token)
	return user, nil
}

// GetSession - returns the session
func (l *MesheryRemoteProvider) GetSession(req *http.Request) (*sessions.Session, error) {
	session, err := l.SessionStore.Get(req, l.SessionName)
	if err != nil {
		err = errors.Wrap(err, "Error: unable to get session")
		logrus.Error(err)
		return nil, err
	}
	return session, nil
}

// GetProviderToken - returns provider token
func (l *MesheryRemoteProvider) GetProviderToken(req *http.Request) (string, error) {
	session, err := l.GetSession(req)
	if err != nil {
		return "", err
	}
	tokenVal, _ := session.Values[l.SaaSTokenName].(string)
	return tokenVal, nil
}

// Logout - logout from provider backend
func (l *MesheryRemoteProvider) Logout(w http.ResponseWriter, req *http.Request) {
	client := http.Client{}
	cReq, err := http.NewRequest(http.MethodGet, l.SaaSBaseURL+"/logout", req.Body)
	if err != nil {
		logrus.Errorf("Error creating a client to logout from tweet app: %v", err)
		http.Error(w, "unable to logout at the moment", http.StatusInternalServerError)
		return
	}
	_, _ = client.Do(cReq)
	// sessionStore.Destroy(w, sessionName)

	sess, err := l.SessionStore.Get(req, l.SessionName)
	if err == nil {
		sess.Options.MaxAge = -1
		_ = sess.Save(req, w)
	}

	http.Redirect(w, req, "/login", http.StatusFound)
}

// FetchResults - fetches results from provider backend
func (l *MesheryRemoteProvider) FetchResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	logrus.Infof("attempting to fetch results from cloud")
	session, _ := l.GetSession(req)

	tokenVal, _ := session.Values[l.SaaSTokenName].(string)

	saasURL, _ := url.Parse(l.SaaSBaseURL + "/results")
	q := saasURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("page_size", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	saasURL.RawQuery = q.Encode()
	logrus.Debugf("constructed results url: %s", saasURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)
	cReq.AddCookie(&http.Cookie{
		Name:     l.SaaSTokenName,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		logrus.Errorf("unable to get results: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("results successfully retrieved from SaaS")
		return bdr, nil
	}
	logrus.Errorf("error while fetching results: %s", bdr)
	return nil, fmt.Errorf("error while fetching results - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetResult - fetches result from provider backend for the given result id
func (l *MesheryRemoteProvider) GetResult(req *http.Request, resultID uuid.UUID) (*MesheryResult, error) {
	logrus.Infof("attempting to fetch result from cloud for id: %s", resultID)
	session, _ := l.GetSession(req)

	tokenVal, _ := session.Values[l.SaaSTokenName].(string)

	saasURL, _ := url.Parse(fmt.Sprintf("%s/result/%s", l.SaaSBaseURL, resultID.String()))
	logrus.Debugf("constructed result url: %s", saasURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)
	cReq.AddCookie(&http.Cookie{
		Name:     l.SaaSTokenName,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		logrus.Errorf("unable to get results: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("result successfully retrieved from SaaS")
		res := &MesheryResult{}
		err = json.Unmarshal(bdr, res)
		if err != nil {
			logrus.Errorf("unable to unmarshal meshery result: %v", err)
			return nil, err
		}
		return res, nil
	}
	logrus.Errorf("error while fetching result: %s", bdr)
	return nil, fmt.Errorf("error while getting result - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// PublishResults - publishes results to the provider backend syncronously
func (l *MesheryRemoteProvider) PublishResults(req *http.Request, result *MesheryResult) (string, error) {
	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return "", err
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	logrus.Infof("attempting to publish results to SaaS")
	bf := bytes.NewBuffer(data)
	session, _ := l.GetSession(req)

	tokenVal, _ := session.Values[l.SaaSTokenName].(string)

	saasURL, _ := url.Parse(l.SaaSBaseURL + "/result")
	cReq, _ := http.NewRequest(http.MethodPost, saasURL.String(), bf)
	cReq.AddCookie(&http.Cookie{
		Name:     l.SaaSTokenName,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		logrus.Errorf("unable to send results: %v", err)
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return "", err
	}
	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("results successfully published to SaaS")
		idMap := map[string]string{}
		if err = json.Unmarshal(bdr, &idMap); err != nil {
			logrus.Errorf("unable to unmarshal body: %v", err)
			return "", err
		}
		resultID, ok := idMap["id"]
		if ok {
			return resultID, nil
		}
		return "", nil
	}
	logrus.Errorf("error while sending results: %s", bdr)
	return "", fmt.Errorf("error while sending results - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// PublishMetrics - publishes metrics to the provider backend asyncronously
func (l *MesheryRemoteProvider) PublishMetrics(tokenVal string, result *MesheryResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return err
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	logrus.Infof("attempting to publish metrics to SaaS")
	bf := bytes.NewBuffer(data)

	saasURL, _ := url.Parse(l.SaaSBaseURL + "/result/metrics")
	cReq, _ := http.NewRequest(http.MethodPut, saasURL.String(), bf)
	cReq.AddCookie(&http.Cookie{
		Name:     l.SaaSTokenName,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		logrus.Errorf("unable to send metrics: %v", err)
		return err
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
		logrus.Errorf("unable to read response body: %v", err)
		return err
	}
	logrus.Errorf("error while sending metrics: %s", bdr)
	return fmt.Errorf("error while sending metrics - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// RecordPreferences - records the user preference
func (l *MesheryRemoteProvider) RecordPreferences(req *http.Request, userID string, data *Preference) error {
	if err := l.BitCaskPreferencePersister.WriteToPersister(userID, data); err != nil {
		return err
	}
	tokenVal, _ := l.GetProviderToken(req)
	l.syncChan <- &userSession{
		token:   tokenVal,
		session: data,
	}
	return nil
}
