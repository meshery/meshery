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
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/model"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

// RemoteProvider - represents a local provider
type RemoteProvider struct {
	ProviderProperties
	*BitCaskPreferencePersister

	SaaSTokenName     string
	RemoteProviderURL string

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
	GenericPersister   database.Handler
	GraphqlHandler     http.Handler
	GraphqlPlayground  http.Handler
	KubeClient         *mesherykube.Client
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
	// Get the capabilities with no token
	// assuming that this will help get basic info
	// of the provider
	l.loadCapabilities("")
}

// loadCapabilities loads the capabilities of the remote provider
//
// It takes in "token" string of the user for loading the capbilities
// if an empty string is provided then it will try to make a request
// with no token, however a remote provider is free to refuse to
// serve requests with no token
func (l *RemoteProvider) loadCapabilities(token string) {
	var resp *http.Response
	var err error

	version := viper.GetString("BUILD")
	os := viper.GetString("OS")
	finalURL := fmt.Sprintf("%s/%s/capabilities?os=%s", l.RemoteProviderURL, version, os)
	finalURL = strings.TrimSuffix(finalURL, "\n")
	remoteProviderURL, err := url.Parse(finalURL)
	if err != nil {
		logrus.Errorf("Error while constructing url: %s", err)
		return
	}

	req, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	// If not token is provided then make a simple GET request
	if token == "" {
		c := &http.Client{}
		resp, err = c.Do(req)
	} else {
		// Proceed to make a request with the token
		resp, err = l.DoRequest(req, token)
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorf("[Initialize Provider]: Failed to get capabilities %s", err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logrus.Errorf("[Initialize]: Failed to close response body %s", err)
		}
	}()

	// Clear the previous capabilities before writing new one
	l.ProviderProperties = ProviderProperties{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&l.ProviderProperties); err != nil {
		logrus.Errorf("[Initialize]: Failed to decode provider properties %s", err)
	}
}

// downloadProviderExtensionPackage will download the remote provider extensions
// package
func (l *RemoteProvider) downloadProviderExtensionPackage() {
	// Location for the package to be stored
	loc := l.PackageLocation()

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

// PackageLocation returns the location of where the package for the current
// provider is located
func (l *RemoteProvider) PackageLocation() string {
	return path.Join(homedir.HomeDir(), ".meshery", "provider", l.ProviderName, l.PackageVersion)
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
	if err := encoder.Encode(l.ProviderProperties); err != nil {
		http.Error(w, "failed to encode provider capabilities", http.StatusInternalServerError)
	}
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

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	req, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bytes.NewReader(bd))

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
		http.Redirect(w, r, l.RemoteProviderURL+"?source="+base64.RawURLEncoding.EncodeToString([]byte(tu))+"&provider_version="+l.ProviderVersion, http.StatusFound)
		return
	}

	// TODO: go to ref cookie
	http.Redirect(w, r, "/", http.StatusFound)
}

func (l *RemoteProvider) fetchUserDetails(tokenString string) (*User, error) {
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + "/user")
	req, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

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

	up := &UserPref{
		Preferences: &Preference{
			AnonymousUsageStats:  true,
			AnonymousPerfResults: true,
		},
	}
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

// FetchResults - fetches results for profile id from provider backend
func (l *RemoteProvider) FetchResults(req *http.Request, page, pageSize, search, order, profileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return []byte{}, fmt.Errorf("%s is not suppported by provider: %s", PersistPerformanceProfiles, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	logrus.Infof("attempting to fetch results from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL)
	remoteProviderURL.Path = path.Join(ep, profileID, "results")
	q := remoteProviderURL.Query()
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
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed results url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

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
		logrus.Infof("results successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching results: %s", bdr)
	return nil, fmt.Errorf("error while fetching results - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// FetchAllResults - fetches results from provider backend
func (l *RemoteProvider) FetchAllResults(req *http.Request, page, pageSize, search, order, from, to string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistResults) {
		logrus.Error("operation not available")
		return []byte{}, fmt.Errorf("%s is not suppported by provider: %s", PersistResults, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistResults)

	logrus.Infof("attempting to fetch results from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
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
	if from != "" {
		q.Set("from", from)
	}
	if to != "" {
		q.Set("to", to)
	}

	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed results url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

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
		logrus.Infof("results successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching results: %s", bdr)
	return nil, fmt.Errorf("error while fetching results - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// FetchSmiResults - fetches results from provider backend
func (l *RemoteProvider) FetchSmiResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSMIResults) {
		logrus.Error("operation not available")
		return []byte{}, fmt.Errorf("%s is not suppported by provider: %s", PersistSMIResults, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMIResults)

	logrus.Infof("attempting to fetch SMI results from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
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
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed smi results url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get token: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get smi results: %v", err)
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
		logrus.Infof("results successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching smi results: %s", bdr)
	return nil, fmt.Errorf("error while fetching smi results - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetResult - fetches result from provider backend for the given result id
func (l *RemoteProvider) GetResult(req *http.Request, resultID uuid.UUID) (*MesheryResult, error) {
	if !l.Capabilities.IsSupported(PersistResult) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistResult, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistResult)

	logrus.Infof("attempting to fetch result from cloud for id: %s", resultID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resultID.String()))
	logrus.Debugf("constructed result url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

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
		logrus.Infof("result successfully retrieved from remote provider")
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
func (l *RemoteProvider) PublishResults(req *http.Request, result *MesheryResult, profileID string) (string, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return "", fmt.Errorf("%s is not supported by provider: %s", PersistPerformanceProfiles, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return "", err
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	logrus.Infof("attempting to publish results to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL)
	remoteProviderURL.Path = path.Join(ep, profileID, "results")
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
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
		logrus.Infof("results successfully published to remote provider")
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
	logrus.Infof("attempting to publish results to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString := viper.GetString("opt-token")
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
		logrus.Infof("results successfully published to remote provider")
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
	logrus.Infof("attempting to publish metrics to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)

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
		logrus.Infof("metrics successfully published to remote provider")
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

// SaveMesheryPattern saves given pattern with the provider
func (l *RemoteProvider) SaveMesheryPattern(tokenString string, pattern *MesheryPattern) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not supported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	data, err := json.Marshal(pattern)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return nil, err
	}

	logrus.Debugf("Pattern: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save pattern to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		logrus.Errorf("unable to get pattern: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send pattern: %v", err)
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

	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("pattern successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	logrus.Errorf("error while sending pattern: %s", bdr)
	return bdr, fmt.Errorf("error while sending pattern - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetMesheryPatterns gives the patterns stored with the provider
func (l *RemoteProvider) GetMesheryPatterns(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		logrus.Error("operation not available")
		return []byte{}, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	logrus.Infof("attempting to fetch patterns from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
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
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed patterns url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get patterns: %v", err)
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get patterns: %v", err)
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
		logrus.Infof("patterns successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching patterns: %s", bdr)
	return nil, fmt.Errorf("error while fetching patterns - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetMesheryPattern gets pattern for the given patternID
func (l *RemoteProvider) GetMesheryPattern(req *http.Request, patternID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	logrus.Infof("attempting to fetch pattern from cloud for id: %s", patternID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, patternID))
	logrus.Debugf("constructed pattern url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get patterns: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get patterns: %v", err)
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
		logrus.Infof("pattern successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching pattern: %s", bdr)
	return nil, fmt.Errorf("error while getting pattern - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// DeleteMesheryPattern deletes a meshery pattern with the given id
func (l *RemoteProvider) DeleteMesheryPattern(req *http.Request, patternID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	logrus.Infof("attempting to fetch pattern from cloud for id: %s", patternID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, patternID))
	logrus.Debugf("constructed pattern url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get patterns: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get patterns: %v", err)
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
		logrus.Infof("pattern successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching pattern: %s", bdr)
	return nil, fmt.Errorf("error while getting pattern - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// SavePerformanceProfile saves a performance profile into the remote provider
func (l *RemoteProvider) SavePerformanceProfile(tokenString string, pp *PerformanceProfile) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not supported by provider: %s", PersistPerformanceProfiles, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	data, err := json.Marshal(pp)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery metrics for shipping"))
		return nil, err
	}

	logrus.Debugf("performance profile: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save performance profile to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		logrus.Errorf("unable to get performance profile: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send performance profile: %v", err)
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

	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("performance profile successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	logrus.Errorf("error while sending performance profile: %s", bdr)
	return bdr, fmt.Errorf("error while sending performance profile - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetPerformanceProfiles gives the performance profiles stored with the provider
func (l *RemoteProvider) GetPerformanceProfiles(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return []byte{}, fmt.Errorf("%s is not suppported by provider: %s", PersistPerformanceProfiles, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	logrus.Infof("attempting to fetch performance profiles from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
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
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed performance profiles url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get performance profiles: %v", err)
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get performance profiles: %v", err)
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
		logrus.Infof("performance profiles successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching performance profiles: %s", bdr)
	return nil, fmt.Errorf("error while fetching performance profiles - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetPerformanceProfile gets performance profile for the given the performanceProfileID
func (l *RemoteProvider) GetPerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistPerformanceProfiles, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	logrus.Infof("attempting to fetch performance profile from cloud for id: %s", performanceProfileID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, performanceProfileID))
	logrus.Debugf("constructed performance profile url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get performance profiles: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get performance profiles: %v", err)
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
		logrus.Infof("performance profile successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching performance profile: %s", bdr)
	return nil, fmt.Errorf("error while getting performance profile - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// DeletePerformanceProfile deletes a performance profile with the given performanceProfileID
func (l *RemoteProvider) DeletePerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistPerformanceProfiles, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	logrus.Infof("attempting to fetch performance profile from cloud for id: %s", performanceProfileID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, performanceProfileID))
	logrus.Debugf("constructed performance profile url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to delete performance profiles: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to delete performance profiles: %v", err)
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
		logrus.Infof("performance profile successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching performance profile: %s", bdr)
	return nil, fmt.Errorf("error while getting performance profile - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// SaveSchedule saves a SaveSchedule into the remote provider
func (l *RemoteProvider) SaveSchedule(tokenString string, s *Schedule) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not supported by provider: %s", PersistSchedules, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	data, err := json.Marshal(s)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal schedule for shipping"))
		return nil, err
	}

	logrus.Debugf("schedule: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save schedule to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		logrus.Errorf("unable to get schedule: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send schedule: %v", err)
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

	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("schedule successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	logrus.Errorf("error while sending schedule: %s", bdr)
	return bdr, fmt.Errorf("error while sending schedule - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetSchedules gives the schedules stored with the provider
func (l *RemoteProvider) GetSchedules(req *http.Request, page, pageSize, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		logrus.Error("operation not available")
		return []byte{}, fmt.Errorf("%s is not suppported by provider: %s", PersistSchedules, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	logrus.Infof("attempting to fetch schedules from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("page_size", pageSize)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed schedules url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get schedules: %v", err)
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get schedules: %v", err)
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
		logrus.Infof("schedules successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching schedules: %s", bdr)
	return nil, fmt.Errorf("error while fetching schedules - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetSchedule gets schedule for the given the scheduleID
func (l *RemoteProvider) GetSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistSchedules, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	logrus.Infof("attempting to fetch schedule from cloud for id: %s", scheduleID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, scheduleID))
	logrus.Debugf("constructed schedule url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get schedules: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get schedules: %v", err)
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
		logrus.Infof("schedule successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching schedule: %s", bdr)
	return nil, fmt.Errorf("error while getting schedule - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// DeleteSchedule deletes a schedule with the given scheduleID
func (l *RemoteProvider) DeleteSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistSchedules, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	logrus.Infof("attempting to fetch schedule from cloud for id: %s", scheduleID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, scheduleID))
	logrus.Debugf("constructed schedule url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to delete schedules: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to delete schedules: %v", err)
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
		logrus.Infof("schedule successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching schedule: %s", bdr)
	return nil, fmt.Errorf("error while getting schedule - Status code: %d, Body: %s", resp.StatusCode, bdr)
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

	// Get new capabilities
	// Doing this here is important so that
	l.loadCapabilities(tokenString)

	// Download the package for the user
	l.downloadProviderExtensionPackage()

	// Proceed to redirect once the capabilities has loaded
	// and the package has been downloaded
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

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
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

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	q.Add("test_uuid", testUUID)
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("Making request to : %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
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

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	q.Add("page", page)
	q.Add("pageSize", pageSize)
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("Making request to : %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
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

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	q.Add("test_uuid", testUUID)
	remoteProviderURL.RawQuery = q.Encode()
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)
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

// RecordMeshSyncData records the mesh sync data
func (l *RemoteProvider) RecordMeshSyncData(obj model.Object) error {
	result := l.GenericPersister.Create(&obj)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ReadMeshSyncData records the mesh sync data
func (l *RemoteProvider) ReadMeshSyncData() ([]model.Object, error) {
	objects := make([]model.Object, 0)
	result := l.GenericPersister.
		Preload("TypeMeta").
		Preload("ObjectMeta").
		Preload("ObjectMeta.Labels").
		Preload("ObjectMeta.Annotations").
		Preload("Spec").
		Preload("Status").
		Find(&objects)

	if result.Error != nil {
		return nil, result.Error
	}
	return objects, nil
}

// TarXZF takes in a source url downloads the tar.gz file
// uncompresses and then save the file to the destination
func TarXZF(srcURL, destination string) error {
	filename := filepath.Base(srcURL)

	// Check if filename ends with tar.gz or tgz extension
	if !strings.HasSuffix(filename, ".tar.gz") && !strings.HasSuffix(filename, ".tgz") {
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

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		loc := path.Join(destination, header.Name)

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

// GetGenericPersister - to return persister
func (l *RemoteProvider) GetGenericPersister() *database.Handler {
	return &l.GenericPersister
}

// GetGraphqlHandler - to return graphql handler instance
func (l *RemoteProvider) GetGraphqlHandler() http.Handler {
	return l.GraphqlHandler
}

// GetGraphqlPlayground - to return graphql playground instance
func (l *RemoteProvider) GetGraphqlPlayground() http.Handler {
	return l.GraphqlPlayground
}

// SetKubeClient - to set meshery kubernetes client
func (l *RemoteProvider) SetKubeClient(client *mesherykube.Client) {
	l.KubeClient = client
}

// GetKubeClient - to get meshery kubernetes client
func (l *RemoteProvider) GetKubeClient() *mesherykube.Client {
	return l.KubeClient
}
