package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/homedir"
)

// RemoteProvider - represents a local provider
type RemoteProvider struct {
	ProviderProperties
	*BitCaskPreferencePersister

	SaaSTokenName string
	SaaSBaseURL   string

	SessionName   string
	RefCookieName string

	TokenStore    map[string]string
	TokenStoreMut sync.Mutex
	Keys          []map[string]string

	LoginCookieDuration time.Duration

	syncStopChan chan struct{}
	syncChan     chan *userSession

	ProviderVersion    string
	SmiResultPersister *BitCaskSmiResultsPersister
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

// Initialize function will initialize the RemoteProvider instance with the metadata
// fetched from the remote providers capabilities endpoint
func (l *RemoteProvider) Initialize() {
	// Get the capabilities
	saasURL, _ := url.Parse(l.SaaSBaseURL + "/capabilities")
	resp, err := http.Get(saasURL.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorf("[Initialize]: Failed to get capabilities %s", err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logrus.Errorf("[Initialize]: Failed to close response body %s", err)
		}
	}()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&l.ProviderProperties); err != nil {
		logrus.Errorf("[Initialize]: Failed to decode provider properties %s", err)
	}

	// Location for the package to be stored
	loc := path.Join(homedir.HomeDir(), ".meshery", "provider", l.ProviderName, l.PackageVersion)

	// Skip download if the file is already present
	if _, err := os.Stat(loc); err == nil {
		logrus.Debugf("[Initialize]: Package found at %s skipping download", loc)
		return
	}

	logrus.Debugf("[Initialize]: Package not found at %s proceeding to download", loc)
	// Download the provider package
	if err := TarXZF(l.PackageURL, loc); err != nil {
		logrus.Errorf("[Initialize]: Failed to download provider package %s", err)
	}
}

// Name - Returns Provider's friendly name
func (l *RemoteProvider) Name() string {
	return l.ProviderName
}

// Description - returns a short description of the provider for display in the Provider UI
func (l *RemoteProvider) Description() []string {
	return l.ProviderDescription
}

const tokenName = "token"

// GetProviderType - Returns ProviderType
func (l *RemoteProvider) GetProviderType() ProviderType {
	return l.ProviderType
}

// GetProviderProperties - Returns all the provider properties required
func (l *RemoteProvider) GetProviderProperties() ProviderProperties {
	return l.ProviderProperties
}

// SyncPreferences - used to sync preferences with the remote provider
func (l *RemoteProvider) SyncPreferences() {
	if !l.Capabilities.IsSupported(SyncPrefs) {
		return
	}

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

// GetProviderCapabilities returns all of the provider properties
func (l *RemoteProvider) GetProviderCapabilities(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(l.ProviderProperties)
}

// StopSyncPreferences - used to stop sync preferences
func (l *RemoteProvider) StopSyncPreferences() {
	if !l.Capabilities.IsSupported(SyncPrefs) {
		return
	}

	l.syncStopChan <- struct{}{}
}

func (l *RemoteProvider) executePrefSync(tokenString string, sess *Preference) {
	ep, exists := l.Capabilities.GetEndpointForFeature(SyncPrefs)
	if !exists {
		logrus.Warn("SyncPrefs is not a supported capability by provider:", l.ProviderName)
		return
	}

	bd, err := json.Marshal(sess)
	if err != nil {
		logrus.Errorf("unable to marshal preference data: %v", err)
		return
	}

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
	req, _ := http.NewRequest(http.MethodPut, saasURL.String(), bytes.NewReader(bd))

	// tokenString, err := l.GetToken(req)
	// if err != nil {
	// 	logrus.Errorf("unable to get results: %v", err)
	// 	return nil, err
	// }
	resp, err := l.DoRequest(req, tokenString)
	if err != nil {
		logrus.Errorf("unable to upload user preference data: %v", err)
		return
	}
	if resp.StatusCode != http.StatusCreated {
		logrus.Errorf("unable to upload user preference data, status code received: %d", resp.StatusCode)
	}
}

// InitiateLogin - initiates login flow and returns a true to indicate the handler to "return" or false to continue
//
// It is assumed that every remote provider will offer this feature
func (l *RemoteProvider) InitiateLogin(w http.ResponseWriter, r *http.Request, _ bool) {
	tu := "http://" + r.Host + r.RequestURI

	_, err := r.Cookie(tokenName)
	// logrus.Debugf("url token: %v %v", token, err)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     l.RefCookieName,
			Value:    "/",
			Expires:  time.Now().Add(l.LoginCookieDuration),
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, l.SaaSBaseURL+"?source="+base64.RawURLEncoding.EncodeToString([]byte(tu))+"&provider_version="+l.ProviderVersion, http.StatusFound)
		return
	}

	// TODO: go to ref cookie
	http.Redirect(w, r, "/", http.StatusFound)
}

func (l *RemoteProvider) fetchUserDetails(tokenString string) (*User, error) {
	saasURL, _ := url.Parse(l.SaaSBaseURL + "/user")
	req, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)

	resp, err := l.DoRequest(req, tokenString)
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
//
// It is assumed that every remote provider will support this feature
func (l *RemoteProvider) GetUserDetails(req *http.Request) (*User, error) {
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	user, err := l.fetchUserDetails(token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetSession - validates the current request, attempts for a refresh of token, and then return its validity
//
// It is assumed that each remote provider will support this feature
func (l *RemoteProvider) GetSession(req *http.Request) error {
	ts, err := l.GetToken(req)
	if err != nil {
		err = fmt.Errorf("session not found")
		logrus.Infof(err.Error())
		return err
	}
	_, err = l.VerifyToken(ts)
	if err != nil {
		logrus.Infof("Token validation error : %v", err.Error())
		newts, err := l.refreshToken(ts)
		if err != nil {
			logrus.Errorf("Token Refresh failed : %v", err.Error())
			return err
		}
		_, err = l.VerifyToken(newts)
		if err != nil {
			logrus.Errorf("Validation of refreshed token failed : %v", err.Error())
			return err
		}
	}
	return nil
}

// GetProviderToken - returns provider token
func (l *RemoteProvider) GetProviderToken(req *http.Request) (string, error) {
	tokenVal, err := l.GetToken(req)
	if err != nil {
		return "", err
	}
	return tokenVal, nil
}

// Logout - logout from provider backend
//
// It is assumed that every remote provider will support this feature
func (l *RemoteProvider) Logout(w http.ResponseWriter, req *http.Request) {
	ck, err := req.Cookie(tokenName)
	if err == nil {
		ck.MaxAge = -1
		http.SetCookie(w, ck)
	}
	http.Redirect(w, req, "/login", http.StatusFound)
}

// FetchResults - fetches results from provider backend
func (l *RemoteProvider) FetchResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistResults) {
		logrus.Error("operation not available")
		return []byte{}, fmt.Errorf("%s is not suppported by provider: %s", PersistResults, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistResults)

	logrus.Infof("attempting to fetch results from cloud")

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
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

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get results: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
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

// FetchSmiResults - fetches results from provider backend
//
// Assuming that every remote provider will support this feature
func (l *RemoteProvider) FetchSmiResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
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
	return l.SmiResultPersister.GetResults(pg, pgs)
}

// GetResult - fetches result from provider backend for the given result id
func (l *RemoteProvider) GetResult(req *http.Request, resultID uuid.UUID) (*MesheryResult, error) {
	if !l.Capabilities.IsSupported(PersistResult) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistResult, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistResult)

	logrus.Infof("attempting to fetch result from cloud for id: %s", resultID)

	saasURL, _ := url.Parse(fmt.Sprintf("%s/%s/%s", l.SaaSBaseURL, ep, resultID.String()))
	logrus.Debugf("constructed result url: %s", saasURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get results: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
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

// PublishResults - publishes results to the provider backend synchronously
func (l *RemoteProvider) PublishResults(req *http.Request, result *MesheryResult) (string, error) {
	if !l.Capabilities.IsSupported(PersistResult) {
		logrus.Error("operation not available")
		return "", fmt.Errorf("%s is not supported by provider: %s", PersistResult, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistResult)

	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return "", err
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	logrus.Infof("attempting to publish results to SaaS")
	bf := bytes.NewBuffer(data)

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, saasURL.String(), bf)
	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get results: %v", err)
		return "", err
	}
	resp, err := l.DoRequest(cReq, tokenString)
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

// PublishSmiResults - publishes results to the provider backend synchronously
func (l *RemoteProvider) PublishSmiResults(result *SmiResult) (string, error) {
	if !l.Capabilities.IsSupported(PersistSMIResult) {
		logrus.Error("operation not available")
		return "", fmt.Errorf("%s is not supported by provider: %s", PersistSMIResult, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMIResult)

	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return "", err
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	logrus.Infof("attempting to publish results to SaaS")
	bf := bytes.NewBuffer(data)

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, saasURL.String(), bf)
	tokenString, err := l.GetToken(nil)
	if err != nil {
		logrus.Errorf("unable to get results: %v", err)
		return "", err
	}
	resp, err := l.DoRequest(cReq, tokenString)
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
func (l *RemoteProvider) PublishMetrics(tokenString string, result *MesheryResult) error {
	if !l.Capabilities.IsSupported(PersistMetrics) {
		logrus.Error("operation not available")
		return fmt.Errorf("%s is not supported by provider: %s", PersistMetrics, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMetrics)

	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return err
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	logrus.Infof("attempting to publish metrics to SaaS")
	bf := bytes.NewBuffer(data)

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
	cReq, _ := http.NewRequest(http.MethodPut, saasURL.String(), bf)

	// tokenString, err := l.GetToken(req)
	// if err != nil {
	// 	logrus.Errorf("unable to get results: %v", err)
	// 	return nil, err
	// }
	resp, err := l.DoRequest(cReq, tokenString)
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
func (l *RemoteProvider) RecordPreferences(req *http.Request, userID string, data *Preference) error {
	if !l.Capabilities.IsSupported(SyncPrefs) {
		logrus.Error("operation not available")
		return fmt.Errorf("%s is not supported by provider: %s", SyncPrefs, l.ProviderName)
	}
	if err := l.BitCaskPreferencePersister.WriteToPersister(userID, data); err != nil {
		return err
	}
	tokenVal, _ := l.GetToken(req)
	l.syncChan <- &userSession{
		token:   tokenVal,
		session: data,
	}
	return nil
}

// TokenHandler - specific to remote auth
func (l *RemoteProvider) TokenHandler(w http.ResponseWriter, r *http.Request, fromMiddleWare bool) {
	tokenString := r.URL.Query().Get(tokenName)
	logrus.Debugf("token : %v", tokenString)
	ck := &http.Cookie{
		Name:     tokenName,
		Value:    string(tokenString),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, ck)
	http.Redirect(w, r, "/", http.StatusFound)
}

// UpdateToken - in case the token was refreshed, this routine updates the response with the new token
func (l *RemoteProvider) UpdateToken(w http.ResponseWriter, r *http.Request) {
	l.TokenStoreMut.Lock()
	defer l.TokenStoreMut.Unlock()

	tokenString, _ := l.GetToken(r)
	newts := l.TokenStore[tokenString]
	if newts != "" {
		logrus.Debugf("set updated token: %v", newts)
		http.SetCookie(w, &http.Cookie{
			Name:     tokenName,
			Value:    newts,
			Path:     "/",
			HttpOnly: true,
		})
	}
}

// ExtractToken - Returns the auth token and the provider type
func (l *RemoteProvider) ExtractToken(w http.ResponseWriter, r *http.Request) {
	l.TokenStoreMut.Lock()
	defer l.TokenStoreMut.Unlock()

	tokenString, err := l.GetToken(r)
	if err != nil {
		logrus.Errorf("Token not found:  %s", err.Error())
		return
	}
	newts := l.TokenStore[tokenString]
	if newts != "" {
		tokenString = newts
	}

	resp := map[string]interface{}{
		"meshery-provider": l.Name(),
		tokenName:          tokenString,
	}
	logrus.Debugf("encoded response : %v", resp)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.Errorf("Unable to extract auth details: %v", err)
		http.Error(w, "unable to extract auth details", http.StatusInternalServerError)
	}
}

// SMPTestConfigStore - persist test profile details to provider
func (l *RemoteProvider) SMPTestConfigStore(req *http.Request, perfConfig *SMP.PerformanceTestConfig) (string, error) {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		logrus.Error("operation not available")
		return "", fmt.Errorf("%s is not supported by provider: %s", PersistSMPTestProfile, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMPTestProfile)

	data, err := json.Marshal(perfConfig)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal testConfig for shipping"))
		return "", err
	}

	bf := bytes.NewBuffer(data)

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, saasURL.String(), bf)
	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get token: %v", err)
		return "", err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send testConfig: %v", err)
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated || err != nil {
		return string(bdr), err
	}
	logrus.Errorf("error while sending testConfig: %s", bdr)
	return "", fmt.Errorf("error while sending testConfig - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// SMPTestConfigGet - retrieve a single test profile details
func (l *RemoteProvider) SMPTestConfigGet(req *http.Request, testUUID string) (*SMP.PerformanceTestConfig, error) {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not supported by provider: %s", PersistSMPTestProfile, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMPTestProfile)

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
	q := saasURL.Query()
	q.Add("test_uuid", testUUID)
	saasURL.RawQuery = q.Encode()
	logrus.Debugf("Making request to : %s", saasURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get token: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get testConfig: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error sending the request: %v", err)
		return nil, err
	}
	logrus.Debugf("%v", string(bdr))
	if resp.StatusCode == http.StatusOK {
		testConfig := SMP.PerformanceTestConfig{}
		err := json.Unmarshal(bdr, &testConfig)
		if err != nil {
			return nil, err
		}
		return &testConfig, nil
	}
	logrus.Errorf("error while getting testConfig: %s", bdr)
	return nil, fmt.Errorf("error while getting testConfig - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// SMPTestConfigFetch - retrieve list of test profiles
func (l *RemoteProvider) SMPTestConfigFetch(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		logrus.Error("operation not available")
		return []byte{}, fmt.Errorf("%s is not supported by provider: %s", PersistSMPTestProfile, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMPTestProfile)

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
	q := saasURL.Query()
	q.Add("page", page)
	q.Add("pageSize", pageSize)
	saasURL.RawQuery = q.Encode()
	logrus.Debugf("Making request to : %s", saasURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get token: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get testConfigs: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK || err != nil {
		return bdr, err
	}
	logrus.Errorf("error while getting testConfigs: %s", bdr)
	return nil, fmt.Errorf("error while getting testConfigs - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// SMPTestConfigDelete - tombstone a given test profile
func (l *RemoteProvider) SMPTestConfigDelete(req *http.Request, testUUID string) error {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		logrus.Error("operation not available")
		return fmt.Errorf("%s is not supported by provider: %s", PersistSMPTestProfile, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMPTestProfile)

	saasURL, _ := url.Parse(l.SaaSBaseURL + ep)
	q := saasURL.Query()
	q.Add("test_uuid", testUUID)
	saasURL.RawQuery = q.Encode()
	cReq, _ := http.NewRequest(http.MethodDelete, saasURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get token: %v", err)
		return err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to delete testConfig: %v", err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	logrus.Errorf("error while deleting testConfig: %s", testUUID)
	return fmt.Errorf("error while deleting testConfig - Status code: %d, Body: %s", resp.StatusCode, testUUID)
}

// TarXZF takes in a source url downloads the tar.gz file
// uncompresses and then save the file to the destination
func TarXZF(srcURL, destination string) error {
	filename := filepath.Base(srcURL)

	// Check if filename ends with tar.gz or tgz extension
	if !strings.HasSuffix(filename, "tar.gz") && !strings.HasSuffix(filename, "tgz") {
		return fmt.Errorf("file is not of type tar.gz or tgz")
	}

	resp, err := http.Get(srcURL)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is not 200")
	}

	return TarXZ(resp.Body, destination)
}

// TarXZ takes in a gzip stream and untars and uncompresses it to
// the destination directory
//
// If the destination doesn't exists, it will create it
func TarXZ(gzipStream io.Reader, destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		loc := createDestinationPath(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(loc, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(loc)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown type: %s", string(header.Typeflag))
		}
	}

	return nil
}

func createDestinationPath(dest, filename string) string {
	return path.Join(dest, filename)
}
