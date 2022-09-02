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
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/model"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

// RemoteProvider - represents a local provider
type RemoteProvider struct {
	ProviderProperties
	*SessionPreferencePersister

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
	SmiResultPersister *SMIResultsPersister
	GenericPersister   *database.Handler
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

const (
	remoteUploadURL   = "/upload"
	remoteDownloadURL = "/download"
)

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
	l.ProviderProperties = ProviderProperties{
		ProviderURL: l.RemoteProviderURL,
	}
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
		http.Error(w, ErrEncoding(err, "Provider Capablity").Error(), http.StatusInternalServerError)
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
		logrus.Error(ErrMarshal(err, "preference data"))
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
	tu := viper.GetString("MESHERY_SERVER_CALLBACK_URL")
	if tu == "" {
		tu = "http://" + r.Host + "/api/user/token" // Hard coding the path because this is what meshery expects
	}

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
		return nil, ErrFetch(err, "User Data", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "User Data")
	}

	up := &UserPref{
		Preferences: &Preference{
			AnonymousUsageStats:  true,
			AnonymousPerfResults: true,
		},
	}
	err = json.Unmarshal(bd, up)
	if err != nil {
		return nil, ErrUnmarshal(err, "User Pref")
	}

	prefLocal, _ := l.ReadFromPersister(up.UserID)
	if prefLocal == nil || up.Preferences.UpdatedAt.After(prefLocal.UpdatedAt) {
		_ = l.WriteToPersister(up.UserID, up.Preferences)
	}

	// Uncomment when Debug verbosity is figured out project wide. | @leecalcote
	// logrus.Debugf("retrieved user: %v", up.User)
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
	if err != nil || ts == "" {
		logrus.Infof("session not found")
		return err
	}

	_, err = l.VerifyToken(ts)
	if err != nil {
		logrus.Infof("Token validation error : %v", err.Error())
		newts, err := l.refreshToken(ts)
		if err != nil {
			return ErrTokenRefresh(err)
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
		ck.Path = "/"
		http.SetCookie(w, ck)
	}
	http.Redirect(w, req, "/provider", http.StatusFound)
}

// HandleUnAuthenticated
//
// Redirects to alert user of expired sesion
func (l *RemoteProvider) HandleUnAuthenticated(w http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("meshery-provider")
	if err == nil {
		ck, err := req.Cookie(tokenName)
		if err == nil {
			ck.MaxAge = -1
			ck.Path = "/"
			http.SetCookie(w, ck)
		}

		http.Redirect(w, req, "/auth/login", http.StatusFound)
		return
	}
	http.Redirect(w, req, "/provider", http.StatusFound)
}

func (l *RemoteProvider) SaveK8sContext(token string, k8sContext K8sContext) (K8sContext, error) {
	data, err := json.Marshal(k8sContext)
	if err != nil {
		return k8sContext, ErrMarshal(err, "kubernetes context error")
	}

	logrus.Infof("attempting to save %s context to remote provider with ID %s", k8sContext.Name, k8sContext.ID)
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + "/user/contexts")
	cReq, err := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	if err != nil {
		return k8sContext, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		logrus.Errorf("unable to send kubernetes context: %v", err)
		return k8sContext, ErrPost(err, "kubernetes context", cReq.Response.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		var kcr K8sContextPersistResponse
		if err := json.NewDecoder(resp.Body).Decode(&kcr); err != nil {
			return k8sContext, ErrUnmarshal(err, "kubernetes context")
		}

		// Sensitive data. Commenting until better debug controls are put into place. - @leecalcote
		// logrus.Infof("kubernetes context successfully sent to remote provider: %+v", kc)

		// If the context already existed, return that as error
		if !kcr.Inserted {
			return kcr.K8sContext, ErrContextAlreadyPersisted
		}

		return kcr.K8sContext, nil
	}
	return k8sContext, ErrPost(fmt.Errorf("failed to save kubernetes context"), fmt.Sprint(resp.Body), resp.StatusCode)
}
func (l *RemoteProvider) GetK8sContexts(token, page, pageSize, search, order string) ([]byte, error) {
	MesheryInstanceID, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		return nil, ErrMesheryInstanceID
	}
	mi := MesheryInstanceID.String()
	logrus.Infof("attempting to fetch kubernetes contexts from cloud for Meshery instance: %s", mi)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + "/user/contexts/" + mi)
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
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		logrus.Errorf("unable to get kubernetes contexts: %v", err)
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
		logrus.Infof("kubernetes contexts successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching kubernetes contexts: %s", bdr)
	return nil, fmt.Errorf("error while fetching kubernetes contexts - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

func (l *RemoteProvider) LoadAllK8sContext(token string) ([]*K8sContext, error) {
	page := 0
	pageSize := 25
	results := []*K8sContext{}

	for {
		res, err := l.GetK8sContexts(token, strconv.Itoa(page), strconv.Itoa(pageSize), "", "")
		if err != nil {
			return results, err
		}
		var k8scontext MesheryK8sContextPage
		err = json.Unmarshal(res, &k8scontext)
		if err != nil {
			return results, ErrMarshal(err, "kubernetes context")
		}
		results = append(results, k8scontext.Contexts...)

		if (page+1)*pageSize >= k8scontext.TotalCount {
			break
		}

		page++
	}

	return results, nil
}

func (l *RemoteProvider) DeleteK8sContext(token, id string) (K8sContext, error) {
	logrus.Infof("attempting to delete kubernetes context from cloud for id: %s", id)
	mesheryInstanceID, _ := viper.Get("INSTANCE_ID").(*uuid.UUID)
	ep := "/user/contexts/" + mesheryInstanceID.String()

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/context/%s", l.RemoteProviderURL, ep, id))
	logrus.Debugf("constructed kubernetes contexts url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		logrus.Errorf("unable to delete kubernetes contexts: %v", err)
		return K8sContext{}, err
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("kubernetes successfully deleted from remote provider")
		return K8sContext{}, nil
	}
	logrus.Errorf("error while deleting kubernetes contexts")
	return K8sContext{}, fmt.Errorf("error while deleting kubernetes context - Status code: %d", resp.StatusCode)
}

func (l *RemoteProvider) GetK8sContext(token, id string) (K8sContext, error) {
	mesheryInstanceID, _ := viper.Get("INSTANCE_ID").(*uuid.UUID)
	ep := "/user/contexts/" + mesheryInstanceID.String()
	logrus.Infof("attempting to fetch kubernetes contexts from cloud for context id: %s", id)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/context/%s", l.RemoteProviderURL, ep, id))
	logrus.Debugf("constructed kubernetes contexts url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		return K8sContext{}, ErrFetch(err, "Kubernetes Context", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		var kc K8sContext
		if err := json.NewDecoder(resp.Body).Decode(&kc); err != nil {
			return kc, ErrUnmarshal(err, "Kubernetes context")
		}

		logrus.Infof("kubernetes context successfully retrieved from remote provider")
		return kc, nil
	}

	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return K8sContext{}, ErrDataRead(err, "Kubernetes context")
	}

	logrus.Errorf("error while fetching kubernetes context: %s", bdr)
	return K8sContext{}, ErrFetch(fmt.Errorf("failed to get kubernetes context"), fmt.Sprint(bdr), resp.StatusCode)
}

// func (l *RemoteProvider) SetCurrentContext(token, id string) (K8sContext, error) {
// 	if id == "" {
// 		return K8sContext{}, ErrContextID
// 	}

// 	mesheryInstanceID, _ := viper.Get("INSTANCE_ID").(*uuid.UUID)
// 	ep := "/user/contexts/" + mesheryInstanceID.String()
// 	logrus.Infof("attempting to set kubernetes contexts from cloud to id: %s", id)

// 	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/current/%s", l.RemoteProviderURL, ep, id))
// 	logrus.Debugf("constructed kubernetes contexts url: %s", remoteProviderURL.String())
// 	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), nil)

// 	resp, err := l.DoRequest(cReq, token)
// 	if err != nil {
// 		return K8sContext{}, ErrFetch(err, "Kubernetes Context", resp.StatusCode)
// 	}
// 	defer func() {
// 		_ = resp.Body.Close()
// 	}()

// 	if resp.StatusCode == http.StatusOK {
// 		var kc K8sContext

// 		logrus.Infof("kubernetes context successfully retrieved from remote provider")
// 		return kc, nil
// 	}

// 	bdr, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return K8sContext{}, ErrDataRead(err, "Kubernetes context")
// 	}

// 	logrus.Errorf("error while setting kubernetes context: %s", bdr)
// 	return K8sContext{}, ErrPost(fmt.Errorf("failed to set current context"), fmt.Sprint(bdr), resp.StatusCode)
// }

// func (l *RemoteProvider) GetCurrentContext(token string) (K8sContext, error) {
// 	mesheryInstanceID, _ := viper.Get("INSTANCE_ID").(*uuid.UUID)
// 	ep := "/user/contexts/" + mesheryInstanceID.String()
// 	logrus.Infof("attempting to fetch current kubernetes contexts from cloud")

// 	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, "current"))
// 	logrus.Debugf("constructed kubernetes contexts url: %s", remoteProviderURL.String())
// 	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

// 	resp, err := l.DoRequest(cReq, token)
// 	if err != nil {
// 		return K8sContext{}, ErrFetch(err, "Kubernetes Context", resp.StatusCode)
// 	}
// 	defer func() {
// 		_ = resp.Body.Close()
// 	}()

// 	if resp.StatusCode == http.StatusOK {
// 		var kc K8sContext
// 		if err := json.NewDecoder(resp.Body).Decode(&kc); err != nil {
// 			return kc, ErrUnmarshal(err, "Kubernetes context")
// 		}

// 		logrus.Infof("kubernetes context successfully retrieved from remote provider")
// 		return kc, nil
// 	}

// 	bdr, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return K8sContext{}, ErrDataRead(err, "Kubernetes context")
// 	}

// 	logrus.Errorf("error while fetching kubernetes current context: %s", bdr)
// 	return K8sContext{}, ErrFetch(fmt.Errorf("failed to retrieve kubernetes current contexts"), fmt.Sprint(bdr), resp.StatusCode)
// }

// FetchResults - fetches results for profile id from provider backend
func (l *RemoteProvider) FetchResults(tokenVal string, page, pageSize, search, order, profileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Warn("operation not available")
		return []byte{}, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
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

	// tokenString, err := l.GetToken(req)
	// if err != nil {
	// 	return nil, err
	// }
	resp, err := l.DoRequest(cReq, tokenVal)
	if err != nil {
		return nil, ErrFetch(err, "Perf Results", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Result")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("results successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching results: %s", bdr)
	return nil, ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
}

// FetchAllResults - fetches results from provider backend
func (l *RemoteProvider) FetchAllResults(tokenString string, page, pageSize, search, order, from, to string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistResults) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("Persist Results", l.ProviderName)
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

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "All Perf results", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "All Perf results")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("results successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching results: %s", bdr)
	return nil, ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
}

// FetchSmiResults - fetches results from provider backend
func (l *RemoteProvider) FetchSmiResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSMIResults) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("PersistSMIResults", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMIResults)

	logrus.Infof("attempting to fetch SMI conformance results from remote provider")

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
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "SMI Result", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "SMI Result")
	}
	if resp.StatusCode == http.StatusOK {
		logrus.Infof("results successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching smi results: %s", bdr)
	return nil, ErrFetch(err, "SMI Result", resp.StatusCode)
}

// FetchSmiResult - fetches single result from provider backend with given id
func (l *RemoteProvider) FetchSmiResult(req *http.Request, page, pageSize, search, order string, resultID uuid.UUID) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSMIResult) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("PersistSMIResult", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMIResult)

	logrus.Infof("attempting to fetch smi result from cloud for id: %s", resultID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resultID.String()))
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
	logrus.Debugf("constructed smi result url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "SMI Result", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "SMI Result")
	}
	if resp.StatusCode == http.StatusOK {
		logrus.Infof("result successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching smi result: %s", bdr)
	return nil, ErrFetch(err, "SMI Result", resp.StatusCode)
}

// GetResult - fetches result from provider backend for the given result id
func (l *RemoteProvider) GetResult(tokenVal string, resultID uuid.UUID) (*MesheryResult, error) {
	if !l.Capabilities.IsSupported(PersistResult) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistResult", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistResult)

	logrus.Infof("attempting to fetch result from cloud for id: %s", resultID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resultID.String()))
	logrus.Debugf("constructed result url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	// tokenString, err := l.GetToken(req)
	// if err != nil {
	// 	return nil, err
	// }
	resp, err := l.DoRequest(cReq, tokenVal)
	if err != nil {
		logrus.Errorf("unable to get results: %v", err)
		return nil, ErrFetch(err, "Perf Result "+resultID.String(), resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Result "+resultID.String())
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("result successfully retrieved from remote provider")
		res := &MesheryResult{}
		err = json.Unmarshal(bdr, res)
		if err != nil {
			logrus.Errorf("unable to unmarshal meshery result: %v", err)
			return nil, ErrUnmarshal(err, "Perf Result "+resultID.String())
		}
		return res, nil
	}
	logrus.Errorf("error while fetching result: %s", bdr)
	return nil, ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
}

// PublishResults - publishes results to the provider backend synchronously
func (l *RemoteProvider) PublishResults(req *http.Request, result *MesheryResult, profileID string) (string, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return "", ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	data, err := json.Marshal(result)
	if err != nil {
		return "", ErrMarshal(err, "meshery metrics for shipping")
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	logrus.Infof("attempting to publish results to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL)
	remoteProviderURL.Path = path.Join(ep, profileID, "results")
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return "", err
	}
	logrus.Info("request: ", cReq)
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send results: %v", err)
		return "", ErrPost(err, "Perf Results", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return "", ErrDataRead(err, "Perf Result")
	}
	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("results successfully published to remote provider")
		idMap := map[string]string{}
		if err = json.Unmarshal(bdr, &idMap); err != nil {
			return "", ErrUnmarshal(err, "Perf Result")
		}
		resultID, ok := idMap["id"]
		if ok {
			return resultID, nil
		}
		return "", nil
	}
	logrus.Errorf("error while sending results: %s", bdr)
	return "", ErrPost(err, fmt.Sprint(bdr), resp.StatusCode)
}

// PublishSmiResults - publishes results to the provider backend synchronously
func (l *RemoteProvider) PublishSmiResults(result *SmiResult) (string, error) {
	if !l.Capabilities.IsSupported(PersistSMIResult) {
		logrus.Error("operation not available")
		return "", ErrInvalidCapability("PersistSMIResult", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMIResult)

	data, err := json.Marshal(result)
	if err != nil {
		return "", ErrMarshal(err, "meshery metrics for shipping")
	}

	logrus.Debugf("Result: %s, size: %d", data, len(data))
	logrus.Infof("attempting to publish results to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString := viper.GetString("opt-token")
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return "", ErrPost(err, "SMI Result", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrDataRead(err, "SMI Result")
	}
	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("results successfully published to remote provider")
		idMap := map[string]string{}
		if err = json.Unmarshal(bdr, &idMap); err != nil {
			return "", ErrUnmarshal(err, "idMap")
		}
		resultID, ok := idMap["id"]
		if ok {
			return resultID, nil
		}
		return "", nil
	}
	return "", ErrPost(err, fmt.Sprint(bdr), resp.StatusCode)
}

// PublishMetrics - publishes metrics to the provider backend asyncronously
func (l *RemoteProvider) PublishMetrics(tokenString string, result *MesheryResult) error {
	if !l.Capabilities.IsSupported(PersistMetrics) {
		logrus.Error("operation not available")
		return ErrInvalidCapability("PersistMetrics", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMetrics)

	data, err := json.Marshal(result)
	if err != nil {
		return ErrMarshal(err, "meshery metrics for shipping")
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
		return ErrPost(err, "metrics", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusOK {
		logrus.Infof("metrics successfully published to remote provider")
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return ErrDataRead(err, "metrics Data")
	}
	logrus.Errorf("error while sending metrics: %s", bdr)
	return ErrPost(err, fmt.Sprint(bdr), resp.StatusCode)
}

func (l *RemoteProvider) SaveMesheryPatternResource(token string, resource *PatternResource) (*PatternResource, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatternResources) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryPatternResources", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatternResources)

	data, err := json.Marshal(resource)
	if err != nil {
		return nil, ErrMarshal(err, "meshery pattern resource")
	}

	logrus.Debugf("Pattern Resource: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save pattern resource to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	// logrus.Debugf("saving pattern to remote provider - constructed URL: %s", remoteProviderURL.String())
	cReq, err := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		logrus.Errorf("unable to send pattern resource: %v", err)
		return nil, ErrPost(err, "pattern", cReq.Response.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		var pr PatternResource
		if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
			return nil, ErrUnmarshal(err, "Pattern Resource")
		}

		logrus.Infof("pattern successfully sent to remote provider: %+v", pr)
		return &pr, nil
	}

	return nil, ErrPost(err, fmt.Sprint(resp.Body), resp.StatusCode)
}

func (l *RemoteProvider) GetMesheryPatternResource(token, resourceID string) (*PatternResource, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatternResources) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryPatternResources", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatternResources)

	logrus.Infof("attempting to fetch pattern resource from cloud for id: %s", resourceID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resourceID))
	logrus.Debugf("constructed pattern url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		return nil, ErrFetch(err, "Pattern Resource", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		var pr PatternResource
		if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
			return nil, ErrUnmarshal(err, "Pattern resource")
		}

		logrus.Infof("pattern resource successfully retrieved from remote provider")
		return &pr, nil
	}

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Pattern Resource")
	}

	logrus.Errorf("error while fetching pattern resource: %s", bdr)
	return nil, ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
}

func (l *RemoteProvider) GetMesheryPatternResources(
	token,
	page,
	pageSize,
	search,
	order,
	name,
	namespace,
	typ,
	oamType string,
) (*PatternResourcePage, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatternResources) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatternResources, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatternResources)
	logrus.Infof("attempting to fetch patterns resource from cloud")

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
	if name != "" {
		q.Set("name", name)
	}
	if namespace != "" {
		q.Set("namespace", namespace)
	}
	if typ != "" {
		q.Set("type", typ)
	}
	if oamType != "" {
		q.Set("oam_type", oamType)
	}

	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed pattern resource url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		logrus.Errorf("unable to get pattern resource: %v", err)
		return nil, ErrFetch(err, "Patterns Page Resource", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		var pr PatternResourcePage
		if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
			logrus.Errorf("unable to read response body: %v", err)
			return nil, ErrUnmarshal(err, "Patterns Page Resource")
		}

		logrus.Infof("pattern resources successfully retrieved from remote provider")
		return &pr, nil
	}

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Patterns Page Resource")
	}

	logrus.Errorf("error while fetching pattern resource: %s", bdr)
	return nil, ErrFetch(fmt.Errorf("error while fetching pattern resource: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
}

func (l *RemoteProvider) DeleteMesheryPatternResource(token, resourceID string) error {
	if !l.Capabilities.IsSupported(PersistMesheryPatternResources) {
		return ErrInvalidCapability("PersistMesheryPatternResources", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatternResources)

	logrus.Infof("attempting to fetch pattern from cloud for id: %s", resourceID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resourceID))
	logrus.Debugf("constructed pattern url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		return ErrDelete(err, "pattern: "+resourceID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("pattern resource successfully deleted from remote provider")
		return nil
	}

	logrus.Errorf("error while deleting pattern resource")
	return ErrDelete(fmt.Errorf("error while deleting pattern resource"), "pattern: "+resourceID, resp.StatusCode)
}

// SaveMesheryPattern saves given pattern with the provider
func (l *RemoteProvider) SaveMesheryPattern(tokenString string, pattern *MesheryPattern) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		return nil, fmt.Errorf("%s is not supported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	data, err := json.Marshal(map[string]interface{}{
		"pattern_data": pattern,
		"save":         true,
	})

	if err != nil {
		err = ErrMarshal(err, "meshery metrics for shipping")
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
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("pattern successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	logrus.Errorf("error while sending pattern: %s", bdr)
	return bdr, fmt.Errorf("error while sending pattern - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// GetMesheryPatterns gives the patterns stored with the provider
func (l *RemoteProvider) GetMesheryPatterns(tokenString string, page, pageSize, search, order string) ([]byte, error) {
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

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get patterns: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
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

// GetCatalogMesheryPatterns gives the catalog patterns stored with the provider
func (l *RemoteProvider) GetCatalogMesheryPatterns(tokenString string, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(MesheryPatternsCatalog) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("MesheryPatternsCatalog", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(MesheryPatternsCatalog)

	logrus.Infof("attempting to fetch catalog patterns from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed catalog patterns url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Pattern Page - Catalog", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Pattern Page - Catalog")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("catalog patterns successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching catalog patterns: %s", bdr)
	return nil, ErrFetch(fmt.Errorf("error while fetching catalog patterns: %s", bdr), "Patterns page - Catalog", resp.StatusCode)
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
	bdr, err := io.ReadAll(resp.Body)
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
	bdr, err := io.ReadAll(resp.Body)
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

// CloneMesheryPattern clones a meshery pattern with the given id
func (l *RemoteProvider) CloneMesheryPattern(req *http.Request, patternID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(CloneMesheryPatterns) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", CloneMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(CloneMesheryPatterns)

	logrus.Infof("attempting to clone pattern from cloud for id: %s", patternID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, patternID))
	logrus.Debugf("constructed pattern url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to clone patterns: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to clone patterns: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("pattern successfully cloned from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while cloning pattern file with id %s: %s", patternID, bdr)
	return nil, fmt.Errorf("error while cloning pattern - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

// DeleteMesheryPatterns deletes meshery patterns with the given ids and names
func (l *RemoteProvider) DeleteMesheryPatterns(req *http.Request, patterns MesheryPatternDeleteRequestBody) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	var reqBodyBuffer bytes.Buffer
	if err := json.NewEncoder(&reqBodyBuffer).Encode(patterns); err != nil {
		logrus.Error("unable to encode json: ", err)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	// Create remote provider-url
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, "delete"))
	logrus.Debugf("constructed pattern url: %s", remoteProviderURL.String())

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), &reqBodyBuffer)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to get patterns: %v", err)
		return nil, err
	}

	// make request
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
	return nil, fmt.Errorf("error while getting pattern - Status code: %d, Body: %s", 200, bdr)
}

func (l *RemoteProvider) RemotePatternFile(req *http.Request, resourceURL, path string, save bool) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryPatterns", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	data, err := json.Marshal(map[string]interface{}{
		"url":  resourceURL,
		"save": save,
		"path": path,
	})

	if err != nil {
		err = ErrMarshal(err, "meshery metrics for shipping")
		return nil, ErrMarshal(err, "meshery metrics for shipping")
	}

	logrus.Debugf("Pattern: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save pattern to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send pattern: %v", err)
		return nil, ErrPost(err, "Pattern File", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Pattern File")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("pattern successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("could not send pattern to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// SaveMesheryFilter saves given filter with the provider
func (l *RemoteProvider) SaveMesheryFilter(tokenString string, filter *MesheryFilter) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	data, err := json.Marshal(map[string]interface{}{
		"filter_data": filter,
		"save":        true,
	})

	if err != nil {
		return nil, ErrMarshal(err, "meshery metrics for shipping")
	}

	logrus.Debugf("Filter: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save filter to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send filter: %v", err)
		return nil, ErrPost(err, "Filters", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filters")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("filter successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("could not send filter to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// GetMesheryFilters gives the filters stored with the provider
func (l *RemoteProvider) GetMesheryFilters(tokenString string, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	logrus.Infof("attempting to fetch filters from cloud")

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
	logrus.Debugf("constructed filters url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Filter Page", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filter Page")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("filters successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching filters: %s", bdr)
	return nil, ErrFetch(fmt.Errorf("error while fetching filters: %s", bdr), "Filters page", resp.StatusCode)
}

// GetCatalogMesheryFilters gives the catalog filters stored with the provider
func (l *RemoteProvider) GetCatalogMesheryFilters(tokenString string, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(MesheryFiltersCatalog) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("MesheryFiltersCatalog", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(MesheryFiltersCatalog)

	logrus.Infof("attempting to fetch catalog filters from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	logrus.Debugf("constructed catalog filters url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Filter Page - Catalog", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filter Page - Catalog")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("catalog filters successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching catalog filters: %s", bdr)
	return nil, ErrFetch(fmt.Errorf("error while fetching catalog filters: %s", bdr), "Filters page - Catalog", resp.StatusCode)
}

// GetMesheryFilterFile gets filter for the given filterID without the metadata
func (l *RemoteProvider) GetMesheryFilterFile(req *http.Request, filterID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	logrus.Infof("attempting to fetch filter from cloud for id: %s", filterID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/file/%s", l.RemoteProviderURL, ep, filterID))
	logrus.Debugf("constructed filter url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get filters: %v", err)
		return nil, ErrFetch(err, "Filter File: "+filterID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filter File: "+filterID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("filter successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("could not retrieve filter from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// GetMesheryFilter gets filter for the given filterID
func (l *RemoteProvider) GetMesheryFilter(req *http.Request, filterID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	logrus.Infof("attempting to fetch filter from cloud for id: %s", filterID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, filterID))
	logrus.Debugf("constructed filter url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Filter:"+filterID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filter:"+filterID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("filter successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("could not retrieve filter from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeleteMesheryFilter deletes a meshery filter with the given id
func (l *RemoteProvider) DeleteMesheryFilter(req *http.Request, filterID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	logrus.Infof("attempting to fetch filter from cloud for id: %s", filterID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, filterID))
	logrus.Debugf("constructed filter url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get filters: %v", err)
		return nil, ErrDelete(err, "Filter: "+filterID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filter: "+filterID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("filter successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching filter: %s", bdr)
	return nil, ErrDelete(fmt.Errorf("error while fetching filter: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
}

// CloneMesheryFilter clones a meshery filter with the given id
func (l *RemoteProvider) CloneMesheryFilter(req *http.Request, filterID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(CloneMesheryFilters) {
		logrus.Error("operation not available")
		return nil, fmt.Errorf("%s is not suppported by provider: %s", CloneMesheryFilters, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(CloneMesheryFilters)

	logrus.Infof("attempting to clone filter from cloud for id: %s", filterID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, filterID))
	logrus.Debugf("constructed filter url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		logrus.Errorf("unable to clone filters: %v", err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to clone filters: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("filter successfully cloned from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while cloning filter with id %s: %s", filterID, bdr)
	return nil, fmt.Errorf("error while cloning filter - Status code: %d, Body: %s", resp.StatusCode, bdr)
}

func (l *RemoteProvider) RemoteFilterFile(req *http.Request, resourceURL, path string, save bool) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	data, err := json.Marshal(map[string]interface{}{
		"url":  resourceURL,
		"save": save,
		"path": path,
	})

	if err != nil {
		err = ErrMarshal(err, "meshery metrics for shipping")
		return nil, err
	}

	logrus.Debugf("Filter: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save filter to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrPost(err, "Filter File", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filter file")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("filter successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("could not send filter to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// SaveMesheryApplication saves given application with the provider
func (l *RemoteProvider) SaveMesheryApplication(tokenString string, application *MesheryApplication) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryApplications) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryApplications", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	data, err := json.Marshal(map[string]interface{}{
		"application_data": application,
		"save":             true,
	})

	if err != nil {
		err = ErrMarshal(err, "meshery metrics for shipping")
		return nil, err
	}

	logrus.Debugf("Application size: %d", len(data))
	logrus.Infof("attempting to save application to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrPost(err, "Application", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Application")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("application successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("failed to send application to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// SaveApplicationSourceContent saves given application source content with the provider after successful save of Application with the provider
func (l *RemoteProvider) SaveApplicationSourceContent(tokenString string, applicationID string, sourceContent []byte) error {
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	logrus.Debugf("Application Content size %d", len(sourceContent))
	bf := bytes.NewBuffer(sourceContent)

	uploadURL := fmt.Sprintf("%s%s%s/%s", l.RemoteProviderURL, ep, remoteUploadURL, applicationID)
	remoteProviderURL, _ := url.Parse(uploadURL)

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return ErrPost(err, "Application Source Content", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("application source successfully uploaded to remote provider")
		return nil
	}

	return ErrPost(fmt.Errorf("failed to upload application source to remote provider"), "", resp.StatusCode)
}

// GetApplicationSourceContent returns application source-content from provider
func (l *RemoteProvider) GetApplicationSourceContent(req *http.Request, applicationID string) ([]byte, error) {
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)
	downloadURL := fmt.Sprintf("%s%s%s/%s", l.RemoteProviderURL, ep, remoteDownloadURL, applicationID)
	remoteProviderURL, _ := url.Parse(downloadURL)
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	logrus.Infof("attempting to fetch application source content from cloud for id: %s", applicationID)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get application source content: %v", err)
		return nil, ErrFetch(err, "Application source content", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return nil, ErrDataRead(err, "Application")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("applications successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching source content: %s", bdr)
	return nil, ErrFetch(fmt.Errorf("error while fetching applications: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
}

// GetMesheryApplications gives the applications stored with the provider
func (l *RemoteProvider) GetMesheryApplications(tokenString string, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryApplications) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("PersistMesheryApplications", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	logrus.Infof("attempting to fetch applications from cloud")

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
	logrus.Debugf("constructed applications url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get applications: %v", err)
		return nil, ErrFetch(err, "Application", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return nil, ErrDataRead(err, "Application")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("applications successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching applications: %s", bdr)
	return nil, ErrFetch(fmt.Errorf("error while fetching applications: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
}

// GetMesheryApplication gets application for the given applicationID
func (l *RemoteProvider) GetMesheryApplication(req *http.Request, applicationID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryApplications) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryApplications", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	logrus.Infof("attempting to fetch application from cloud for id: %s", applicationID)
	urls := fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, applicationID)
	remoteProviderURL, _ := url.Parse(urls)
	logrus.Debugf("constructed application url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Application: "+applicationID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Application: "+applicationID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("application successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to retrieve application from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeleteMesheryApplication deletes a meshery application with the given id
func (l *RemoteProvider) DeleteMesheryApplication(req *http.Request, applicationID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryApplications) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistMesheryApplications", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	logrus.Infof("attempting to fetch application from cloud for id: %s", applicationID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, applicationID))
	logrus.Debugf("constructed application url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrDelete(err, "Application :"+applicationID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Application :"+applicationID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("application successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrDelete(fmt.Errorf("could not retrieve application from remote provider"), "Application :"+applicationID, resp.StatusCode)
}

// SavePerformanceProfile saves a performance profile into the remote provider
func (l *RemoteProvider) SavePerformanceProfile(tokenString string, pp *PerformanceProfile) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	data, err := json.Marshal(pp)
	if err != nil {
		err = ErrMarshal(err, "meshery metrics for shipping")
		return nil, err
	}

	logrus.Debugf("performance profile: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save performance profile to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrPost(err, "Perf Profile", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Profile")
	}

	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("performance profile successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("failed to send performance profile to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// GetPerformanceProfiles gives the performance profiles stored with the provider
func (l *RemoteProvider) GetPerformanceProfiles(tokenString string, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
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

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Perf Profile Page", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Profile Page")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("performance profiles successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrPost(fmt.Errorf("failed to retrieve performance profile from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// GetPerformanceProfile gets performance profile for the given the performanceProfileID
func (l *RemoteProvider) GetPerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	logrus.Infof("attempting to fetch performance profile from cloud for id: %s", performanceProfileID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, performanceProfileID))
	logrus.Debugf("constructed performance profile url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to get performance profiles: %v", err)
		return nil, ErrFetch(err, "Perf Profile :"+performanceProfileID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Profile :"+performanceProfileID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("performance profile successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to retrieve performance profile from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeletePerformanceProfile deletes a performance profile with the given performanceProfileID
func (l *RemoteProvider) DeletePerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	logrus.Infof("attempting to fetch performance profile from cloud for id: %s", performanceProfileID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, performanceProfileID))
	logrus.Debugf("constructed performance profile url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrDelete(err, "Perf Profile :"+performanceProfileID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Profile :"+performanceProfileID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("performance profile successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrDelete(fmt.Errorf("failed to retrieve performance profile from remote provider"), "Perf Profile :"+performanceProfileID, resp.StatusCode)
}

// SaveSchedule saves a SaveSchedule into the remote provider
func (l *RemoteProvider) SaveSchedule(tokenString string, s *Schedule) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistSchedules", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	data, err := json.Marshal(s)
	if err != nil {
		return nil, ErrMarshal(err, "schedule for shipping")
	}

	logrus.Debugf("schedule: %s, size: %d", data, len(data))
	logrus.Infof("attempting to save schedule to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send schedule: %v", err)
		return nil, ErrPost(err, "Perf Schedule", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Schedule")
	}

	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("schedule successfully sent to remote provider: %s", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("failed to send schedule to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// GetSchedules gives the schedules stored with the provider
func (l *RemoteProvider) GetSchedules(req *http.Request, page, pageSize, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("PersistSchedules", l.ProviderName)
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
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Perf Schedule Page", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Schedule Page")
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("schedules successfully retrieved from remote provider")
		return bdr, nil
	}
	logrus.Errorf("error while fetching schedules: %s", bdr)
	return nil, ErrFetch(fmt.Errorf("error while fetching schedules: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
}

// GetSchedule gets schedule for the given the scheduleID
func (l *RemoteProvider) GetSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistSchedules", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	logrus.Infof("attempting to fetch schedule from cloud for id: %s", scheduleID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, scheduleID))
	logrus.Debugf("constructed schedule url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Perf Schedule :"+scheduleID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Schedule :"+scheduleID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("schedule successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("could not retrieve schedule from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeleteSchedule deletes a schedule with the given scheduleID
func (l *RemoteProvider) DeleteSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistSchedules", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	logrus.Infof("attempting to fetch schedule from cloud for id: %s", scheduleID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, scheduleID))
	logrus.Debugf("constructed schedule url: %s", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to delete schedules: %v", err)
		return nil, ErrDelete(err, "Perf Schedule :"+scheduleID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Schedule :"+scheduleID)
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("schedule successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrDelete(fmt.Errorf("could not retrieve schedule from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// RecordPreferences - records the user preference
func (l *RemoteProvider) RecordPreferences(req *http.Request, userID string, data *Preference) error {
	if !l.Capabilities.IsSupported(SyncPrefs) {
		logrus.Error("operation not available")
		return ErrInvalidCapability("SyncPrefs", l.ProviderName)
	}
	if err := l.SessionPreferencePersister.WriteToPersister(userID, data); err != nil {
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
func (l *RemoteProvider) UpdateToken(w http.ResponseWriter, r *http.Request) string {
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
		return newts
	}

	return tokenString
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
		err = ErrEncoding(err, "Auth Details")
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SMPTestConfigStore - persist test profile details to provider
func (l *RemoteProvider) SMPTestConfigStore(req *http.Request, perfConfig *SMP.PerformanceTestConfig) (string, error) {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		logrus.Error("operation not available")
		return "", ErrInvalidCapability("PersistSMPTestProfile", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMPTestProfile)

	data, err := json.Marshal(perfConfig)
	if err != nil {
		return "", ErrMarshal(err, "testConfig for shipping")
	}

	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return "", err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		logrus.Errorf("unable to send testConfig: %v", err)
		return "", ErrPost(err, "Perf Test Config", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated || err != nil {
		return string(bdr), err
	}
	return "", ErrPost(fmt.Errorf("could not send test profile details to remote provider: %d", resp.StatusCode), fmt.Sprint(bdr), resp.StatusCode)
}

// SMPTestConfigGet - retrieve a single test profile details
func (l *RemoteProvider) SMPTestConfigGet(req *http.Request, testUUID string) (*SMP.PerformanceTestConfig, error) {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		logrus.Error("operation not available")
		return nil, ErrInvalidCapability("PersistSMPTestProfile", l.ProviderName)
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
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Perf Test Config: "+testUUID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Test Config :"+testUUID)
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
	return nil, ErrFetch(fmt.Errorf("could not retrieve test profile details: %d", resp.StatusCode), fmt.Sprint(bdr), resp.StatusCode)
}

// SMPTestConfigFetch - retrieve list of test profiles
func (l *RemoteProvider) SMPTestConfigFetch(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		logrus.Error("operation not available")
		return []byte{}, ErrInvalidCapability("PersistSMPTestProfile", l.ProviderName)
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
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Perf Test Config Page", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK || err != nil {
		return bdr, err
	}

	return nil, ErrFetch(fmt.Errorf("could not retrieve list of test profiles: %d", resp.StatusCode), fmt.Sprint(bdr), resp.StatusCode)
}

// SMPTestConfigDelete - tombstone a given test profile
func (l *RemoteProvider) SMPTestConfigDelete(req *http.Request, testUUID string) error {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		logrus.Error("operation not available")
		return ErrInvalidCapability("PersistSMPTestProfile", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMPTestProfile)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	q.Add("test_uuid", testUUID)
	remoteProviderURL.RawQuery = q.Encode()
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return ErrDelete(err, "Perf Test Config :"+testUUID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return ErrDelete(fmt.Errorf("could not delete the test profile: %d", resp.StatusCode), "Perf Test Config :"+testUUID, resp.StatusCode)
}

func (l *RemoteProvider) ExtensionProxy(req *http.Request) ([]byte, error) {
	logrus.Infof("attempting to request remote provider")
	p := req.URL.Path
	split := strings.Split(p, "/api/extensions")
	path := split[1]
	q := req.URL.Query().Encode()
	if len(q) > 0 {
		path = fmt.Sprintf("%s?%s", path, q)
	}
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, path))
	logrus.Debugf("constructed url: %s", remoteProviderURL.String())

	cReq, _ := http.NewRequest(req.Method, remoteProviderURL.String(), req.Body)
	tokenString, err := l.GetToken(req)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		logrus.Infof("response successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to request to remote provider"), fmt.Sprint(bdr), resp.StatusCode)
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
		return ErrFetch(err, "TarTZF file :"+srcURL, resp.StatusCode)
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

		// Prevent Arbitrary file write during zip extraction ("zip slip")
		// This checks that the zip doesn't contain filenames like
		// `../../../tmp/some.sh`
		// see https://snyk.io/research/zip-slip-vulnerability for more
		err = validateExtractPath(header.Name, destination)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path.Join(destination, header.Name), 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// When we encounter files that are in nested dirs this takes care of
			// creating parent dirs.
			// #nosec
			if _, err := os.Stat(path.Join(destination, path.Dir(header.Name))); err != nil {
				if err := os.MkdirAll(path.Join(destination, path.Dir(header.Name)), 0750); err != nil {
					return err
				}
			}

			outFile, err := os.Create(path.Join(destination, header.Name))
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.CopyN(outFile, tarReader, header.Size); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown type: %s", string(header.Typeflag))
		}
	}
	return nil
}

func validateExtractPath(filePath string, destination string) error {
	destpath := filepath.Join(destination, filePath)
	if !strings.HasPrefix(destpath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", filePath)
	}
	return nil
}

// GetGenericPersister - to return persister
func (l *RemoteProvider) GetGenericPersister() *database.Handler {
	return l.GenericPersister
}

// SetKubeClient - to set meshery kubernetes client
func (l *RemoteProvider) SetKubeClient(client *mesherykube.Client) {
	l.KubeClient = client
}

// GetKubeClient - to get meshery kubernetes client
func (l *RemoteProvider) GetKubeClient() *mesherykube.Client {
	return l.KubeClient
}
