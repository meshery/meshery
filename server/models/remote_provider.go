package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/meshery/schemas/models/v1beta1"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

// RemoteProvider - represents a local provider
type RemoteProvider struct {
	ProviderProperties
	*SessionPreferencePersister
	*EventsPersister
	*UserCapabilitiesPersister

	SaaSTokenName     string
	RemoteProviderURL string

	SessionName   string
	RefCookieName string

	TokenStore    map[string]string
	TokenStoreMut sync.Mutex
	Keys          []map[string]string

	LoginCookieDuration time.Duration

	// provider and token cookie expiry bound
	CookieDuration time.Duration

	syncStopChan chan struct{}
	syncChan     chan *userSession

	ProviderVersion    string
	SmiResultPersister *SMIResultsPersister
	GenericPersister   *database.Handler
	KubeClient         *mesherykube.Client
	Log                logger.Handler
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
	refURLCookie      = "meshery_ref"
)

// Initialize function will initialize the RemoteProvider instance with the metadata
// fetched from the remote providers capabilities endpoint
func (l *RemoteProvider) Initialize() {
	// Get the capabilities with no token
	// assuming that this will help get basic info
	// of the provider
	providerProperties := l.loadCapabilities("")
	l.ProviderProperties = providerProperties
}

func (l *RemoteProvider) SetProviderProperties(providerProperties ProviderProperties) {
	l.ProviderProperties = providerProperties
}

// loadCapabilities loads the capabilities of the remote provider
//
// It takes in "token" string of the user for loading the capbilities
// if an empty string is provided then it will try to make a request
// with no token, however a remote provider is free to refuse to
// serve requests with no token
func (l *RemoteProvider) loadCapabilities(token string) ProviderProperties {
	var resp *http.Response
	var err error

	providerProperties := ProviderProperties{
		ProviderURL: l.RemoteProviderURL,
	}

	version := viper.GetString("BUILD")
	os := viper.GetString("OS")
	playground := viper.GetString("PLAYGROUND")
	finalURL := fmt.Sprintf("%s/%s/capabilities?os=%s&playground=%s", l.RemoteProviderURL, version, os, playground)
	finalURL = strings.TrimSuffix(finalURL, "\n")
	remoteProviderURL, err := url.Parse(finalURL)
	if err != nil {
		l.Log.Error(ErrUrlParse(err))
		return providerProperties
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
	if err != nil && resp == nil {
		l.Log.Error(ErrUnreachableRemoteProvider(err))
		return providerProperties
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		l.Log.Error(ErrFetch(err, "Capabilities", http.StatusInternalServerError))
		return providerProperties
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			err := ErrCloseIoReader(err)
			l.Log.Error(err)
		}
	}()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&providerProperties); err != nil {
		err = ErrUnmarshal(err, "provider properties")
		l.Log.Error(err)
	}
	return providerProperties
}

// downloadProviderExtensionPackage will download the remote provider extensions
// package
func (l *ProviderProperties) DownloadProviderExtensionPackage(log logger.Handler) {
	// Location for the package to be stored
	loc := l.PackageLocation()

	log.Infof("Package location %s", loc)

	// Skip download if the file is already present
	if _, err := os.Stat(loc); err == nil {
		log.Debug(fmt.Sprintf("[Initialize]: Package found at %s skipping download", loc))
		return
	}

	log.Info(fmt.Sprintf("[Initialize]: Package not found at %s proceeding to download", loc))
	// logrus the provider package
	if err := TarXZF(l.PackageURL, loc, log); err != nil {
		log.Error(ErrDownloadPackage(err, "provider package"))
	}
}

// PackageLocation returns the location of where the package for the current
// provider is located
func (l *ProviderProperties) PackageLocation() string {
	return path.Join(homedir.HomeDir(), ".meshery", "provider", l.ProviderName, l.PackageVersion)
}

// Name - Returns Provider's friendly name
func (l *RemoteProvider) Name() string {
	return l.ProviderName
}

func (l *RemoteProvider) GetProviderURL() string {
	return l.ProviderURL
}

// Description - returns a short description of the provider for display in the Provider UI
func (l *RemoteProvider) Description() []string {
	return l.ProviderDescription
}

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
func (l *RemoteProvider) GetProviderCapabilities(w http.ResponseWriter, req *http.Request, userID string) {
	tokenString := req.Context().Value(TokenCtxKey).(string)

	providerProperties := l.loadCapabilities(tokenString)
	providerProperties.ProviderURL = l.RemoteProviderURL
	if err := l.WriteCapabilitiesForUser(userID, &providerProperties); err != nil {
		l.Log.Error(ErrDBPut(errors.Join(err, fmt.Errorf("failed to write capabilities for the user %s to the server store", userID))))
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(providerProperties); err != nil {
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
		l.Log.Warn(ErrInvalidCapability("SyncPrefs", l.ProviderName))
		return
	}

	bd, err := json.Marshal(sess)
	if err != nil {
		l.Log.Error(ErrMarshal(err, "preference data"))
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
		if resp == nil {
			l.Log.Error(ErrUnreachableRemoteProvider(err))
			return
		}
		l.Log.Error(ErrPost(err, "user preference data", 0))
		return
	}
	if resp.StatusCode != http.StatusCreated {
		err = ErrPost(fmt.Errorf("status code: %d", resp.StatusCode), "user preference data", resp.StatusCode)
		l.Log.Error(err)
	}
}

// InitiateLogin - initiates login flow and returns a true to indicate the handler to "return" or false to continue
//
// Every Remote Provider must offer this function
func (l *RemoteProvider) InitiateLogin(w http.ResponseWriter, r *http.Request, _ bool) {
	baseCallbackURL := r.Context().Value(MesheryServerCallbackURL).(string)

	// Support for deep-link and redirection to land user on their originally requested page post authentication instead of dropping user on the root (home) page.
	refURLqueryParam := r.URL.Query().Get("ref")

	mesheryVersion := viper.GetString("BUILD")

	callbackURL, _ := url.Parse(baseCallbackURL)
	callbackURL = callbackURL.JoinPath(r.URL.EscapedPath())
	callbackURL.RawQuery = r.URL.RawQuery

	ck, err := r.Cookie(TokenCookieName)
	if err != nil || ck.Value == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     l.RefCookieName,
			Value:    "/",
			Expires:  time.Now().Add(l.LoginCookieDuration),
			Path:     "/",
			HttpOnly: true,
		})

		var refURL []string
		// If refURL is empty, generate the refURL based on the current request path and query param.
		if refURLqueryParam == "" {
			refURL = []string{base64.RawURLEncoding.EncodeToString([]byte(strings.TrimPrefix(callbackURL.String(), baseCallbackURL)))}
		} else {
			refURL = append(refURL, refURLqueryParam)
		}

		queryParams := url.Values{
			"source":           []string{base64.RawURLEncoding.EncodeToString([]byte(baseCallbackURL))},
			"provider_version": []string{l.ProviderVersion},
			"meshery_version":  []string{mesheryVersion},
			"ref":              refURL,
		}

		releaseChannel := NewReleaseChannelInterceptor(viper.GetString("RELEASE_CHANNEL"), l, l.Log)
		if releaseChannel != nil {
			releaseChannel.Intercept(r, w)
			return
		}

		http.Redirect(w, r, l.RemoteProviderURL+"/login?"+queryParams.Encode(), http.StatusFound)
		return
	}

	// TODO: go to ref cookie
	http.Redirect(w, r, "/", http.StatusFound)
}

// GetUserDetails - returns the user details
func (l *RemoteProvider) GetUserDetails(req *http.Request) (*User, error) {
	if !l.Capabilities.IsSupported(UsersProfile) {
		l.Log.Warn(ErrOperationNotAvaibale)
		return &User{}, ErrInvalidCapability("UserProfile", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(UsersProfile)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "User Data", http.StatusUnauthorized)
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

	if prefLocal == nil || up.Preferences.UpdatedAt.After(prefLocal.UpdatedAt) || !reflect.DeepEqual(up.Preferences.RemoteProviderPreferences, prefLocal.RemoteProviderPreferences) {
		_ = l.WriteToPersister(up.UserID, up.Preferences)
	}

	// Uncomment when Debug verbosity is figured out project wide. | @leecalcote
	// l.Log.Debug("retrieved user: %v", up.User)
	return &up.User, nil
}

func (l *RemoteProvider) GetUserByID(req *http.Request, userID string) ([]byte, error) {
	systemID := viper.GetString("INSTANCE_ID")
	if userID == systemID {
		return nil, nil
	}
	if !l.Capabilities.IsSupported(UsersProfile) {
		err := ErrInvalidCapability("UsersProfile", l.ProviderName)
		l.Log.Warn(err)
		return []byte{}, err
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(UsersProfile)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, userID))
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "User Profile", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "response body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("User profile retrieved from remote provider.")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("Error retrieving user with ID: %s", userID), "User Profile", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

func (l *RemoteProvider) GetUsers(token, page, pageSize, search, order, filter string) ([]byte, error) {
	if !l.Capabilities.IsSupported(UsersIdentity) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("UsersIdentity", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(UsersIdentity)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Users Data", http.StatusUnauthorized)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Users Data")
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("user data successfully retrieved from remote provider")
		return bd, nil
	}
	err = ErrFetch(err, "Users Data", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// Returns Keys from a user /api/identity/users/keys
func (l *RemoteProvider) GetUsersKeys(token, page, pageSize, search, order, filter string, orgID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(UsersKeys) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("UsersKeys", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistOrganizations)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + orgID + "/users/keys")
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	remoteProviderURL.RawQuery = q.Encode()
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(fmt.Errorf("unable to fetch keys for the org id %s", orgID), "Users keys", http.StatusUnauthorized)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Users Keys")
	}
	if resp.StatusCode == http.StatusOK {
		l.Log.Info("user keys successfully retrieved from remote provider")
		return bd, nil
	}
	err = ErrFetch(fmt.Errorf("unable to fetch keys for the org id %s", orgID), "Users Keys", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetSession - validates the current request, attempts for a refresh of token, and then return its validity
//
// It is assumed that each remote provider will support this feature
func (l *RemoteProvider) GetSession(req *http.Request) error {
	ts, err := l.GetToken(req)
	if err != nil || ts == "" {
		l.Log.Info("session not found")
		return ErrEmptySession
	}
	jwtClaims, err := l.VerifyToken(ts)
	if err != nil {
		err = ErrTokenClaims
		l.Log.Error(err)
		return err
	}
	if jwtClaims == nil {
		err = ErrNilJWKs
		l.Log.Error(err)
		return err
	}
	// we verify the signature of the token and check if it has exp claim,
	// if not present it's an infinite JWT, hence skip the introspect step
	//
	if (*jwtClaims)["exp"] != nil {
		err = l.introspectToken(ts)
		if err != nil {
			return err
		}
	}

	if err != nil {
		l.Log.Info("Token validation error : ", err.Error())
		newts, err := l.refreshToken(ts)
		if err != nil {
			err = ErrTokenRefresh(err)
			l.Log.Error(err)
			return err
		}
		_, err = l.VerifyToken(newts)
		if err != nil {
			err = ErrTokenVerify(err)
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
func (l *RemoteProvider) Logout(w http.ResponseWriter, req *http.Request) error {
	// construct remote provider logout url
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, "/logout"))
	l.Log.Debug("constructed url: ", remoteProviderURL.String())

	// make http.Request type variable with the constructed URL
	cReq, _ := http.NewRequest(req.Method, remoteProviderURL.String(), req.Body)
	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrLogout(err)
		l.Log.Error(err)
		return err
	}

	// gets session cookie from the request headers
	sessionCookie, err := req.Cookie("session_cookie")
	if err != nil {
		err = ErrGetSessionCookie(err)
		l.Log.Error(err)
		return err
	}

	// adds session cookie to the new request headers
	// necessary to run logout flow on the remote provider
	cReq.AddCookie(&http.Cookie{
		Name:  "session_cookie",
		Value: sessionCookie.Value,
	})

	// adds return_to cookie to the new request headers
	// necessary to inform remote provider to return back to Meshery UI
	cReq.AddCookie(&http.Cookie{Name: "return_to", Value: "provider_ui"})

	// make request to remote provider with contructed URL and updated headers (like session_cookie, return_to cookies)
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		err = ErrUnreachableRemoteProvider(err)
		l.Log.Error(err)
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrLogout(err))
		return err
	}
	l.Log.Info("response successfully retrieved from remote provider")
	// if request succeeds then redirect to Provider UI
	// And empties the token and session cookies
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusOK {
		// gets the token from the request headers
		ck, err := req.Cookie(TokenCookieName)
		if err == nil {
			err = l.revokeToken(ck.Value)
		}
		if err != nil {
			l.Log.Error(ErrLogout(fmt.Errorf("token cannot be revoked:%v", err)))
		}
		l.UnSetJWTCookie(w)

		l.UnSetProviderSessionCookie(w)
		return nil
	}
	l.Log.Error(ErrLogout(fmt.Errorf("error performing logout: %s", bd)))
	return errors.New(string(bd))
}

// HandleUnAuthenticated
//
// Redirects to alert user of expired sesion
func (l *RemoteProvider) HandleUnAuthenticated(w http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("meshery-provider")
	if err == nil {
		// remove the cookie from the browser and redirect to inform about expired session.
		l.UnSetJWTCookie(w)
		http.Redirect(w, req, "/auth/login?"+req.URL.RawQuery, http.StatusFound)
		return
	}
	http.Redirect(w, req, "/provider", http.StatusFound)
}

func (l *RemoteProvider) SaveK8sContext(token string, k8sContext K8sContext) (connections.Connection, error) {

	k8sServerID := *k8sContext.KubernetesServerID

	_metadata := map[string]string{
		"id":                   k8sContext.ID,
		"server":               k8sContext.Server,
		"meshery_instance_id":  k8sContext.MesheryInstanceID.String(),
		"deployment_type":      k8sContext.DeploymentType,
		"version":              k8sContext.Version,
		"name":                 k8sContext.Name,
		"kubernetes_server_id": k8sServerID.String(),
	}
	metadata := make(map[string]interface{}, len(_metadata))
	for k, v := range _metadata {
		metadata[k] = v
	}

	cred := map[string]interface{}{
		"auth":    k8sContext.Auth,
		"cluster": k8sContext.Cluster,
	}

	conn := &connections.ConnectionPayload{
		Kind:    "kubernetes",
		Type:    "platform",
		SubType: "orchestrator",
		// Eventually the status would depend on other factors like, whether user administratively processed it or not
		// Is clsuter reachable and other reasons.
		Status:           connections.DISCOVERED,
		MetaData:         metadata,
		CredentialSecret: cred,
	}

	connection, err := l.SaveConnection(conn, token, true)
	if err != nil {
		l.Log.Error(ErrPersistConnection(err))
		return connections.Connection{}, err
	}

	return *connection, nil
}
func (l *RemoteProvider) GetK8sContexts(token, page, pageSize, search, order string, withStatus string, withCredentials bool) ([]byte, error) {
	MesheryInstanceID, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		return nil, ErrMesheryInstanceID
	}
	mi := MesheryInstanceID.String()
	l.Log.Info("attempting to fetch kubernetes contexts from cloud for Meshery instance: ", mi)
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/kubernetes", l.RemoteProviderURL, ep))
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if withStatus != "" {
		q.Set("status", string(connections.CONNECTED))
	}
	if !withCredentials {
		q.Set("with_credentials", "false")
	}
	q.Set("meshery_instance_id", mi)
	remoteProviderURL.RawQuery = q.Encode()
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "Kubernetes Contexts", resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("kubernetes contexts successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("unable to fetch kubernetes contexts"), "Kubernetes Contexts", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

func (l *RemoteProvider) LoadAllK8sContext(token string) ([]*K8sContext, error) {
	page := 0
	pageSize := 25
	results := []*K8sContext{}

	for {
		res, err := l.GetK8sContexts(token, strconv.Itoa(page), strconv.Itoa(pageSize), "", "", string(connections.CONNECTED), true)
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
	l.Log.Info("attempting to delete kubernetes context from cloud for id: ", id)
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return K8sContext{}, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, id))
	l.Log.Debug("constructed kubernetes contexts url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return K8sContext{}, ErrUnreachableRemoteProvider(err)
		}
		err = ErrDelete(err, "Kubernetes Context", resp.StatusCode)
		l.Log.Error(err)
		return K8sContext{}, err
	}

	if resp.StatusCode == http.StatusOK {
		defer func() {
			_ = resp.Body.Close()
		}()

		deletedContext := K8sContext{}
		_ = json.NewDecoder(resp.Body).Decode(&deletedContext)
		l.Log.Info("kubernetes successfully deleted from remote provider")
		return deletedContext, nil
	}
	err = ErrDelete(fmt.Errorf("unable to delete kubernetes context with id: %s", id), "Kubernetes Context", resp.StatusCode)
	l.Log.Error(err)
	return K8sContext{}, err
}

func (l *RemoteProvider) GetK8sContext(token, connectionID string) (K8sContext, error) {
	l.Log.Info("attempting to fetch kubernetes contexts from cloud for connection id: ", connectionID)

	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return K8sContext{}, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/kubernetes/%s", l.RemoteProviderURL, ep, connectionID))

	l.Log.Debug("constructed kubernetes contexts url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return K8sContext{}, ErrUnreachableRemoteProvider(err)
		}
		return K8sContext{}, ErrFetch(err, "Kubernetes Context", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		var kc MesheryK8sContextPage
		if err := json.NewDecoder(resp.Body).Decode(&kc); err != nil {
			return K8sContext{}, ErrUnmarshal(err, "Kubernetes context")
		}

		if len(kc.Contexts) == 0 {
			return K8sContext{}, fmt.Errorf("no Kubernetes contexts available")
		}

		l.Log.Info("Retrieved Kubernetes context from remote provider.")
		// Response will contain single context
		return *kc.Contexts[0], nil
	}

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return K8sContext{}, ErrDataRead(err, "Kubernetes context")
	}
	err = ErrFetch(fmt.Errorf("Failed to retrieve Kubernetes context."), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return K8sContext{}, err
}

// FetchResults - fetches results for profile id from provider backend
func (l *RemoteProvider) FetchResults(tokenVal string, page, pageSize, search, order, profileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	l.Log.Info("attempting to fetch results from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL)
	remoteProviderURL.Path = path.Join(ep, profileID, "results")
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed results url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	// tokenString, err := l.GetToken(req)
	// if err != nil {
	// 	return nil, err
	// }
	resp, err := l.DoRequest(cReq, tokenVal)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("Retrieved results from remote provider")
		return bdr, nil
	}
	err = ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// FetchAllResults - fetches results from provider backend
func (l *RemoteProvider) FetchAllResults(tokenString string, page, pageSize, search, order, from, to string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistResults) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("Persist Results", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistResults)

	l.Log.Info("Fetching results from remote provider.")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
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
	l.Log.Debug("constructed results url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("Retrieved results from remote provider.")
		return bdr, nil
	}
	err = ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// FetchSmiResults - fetches results from provider backend
func (l *RemoteProvider) FetchSmiResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSMIResults) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("PersistSMIResults", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMIResults)

	l.Log.Info("attempting to fetch SMI conformance results from remote provider")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed smi results url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("results successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(err, "SMI Result", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// FetchSmiResult - fetches single result from provider backend with given id
func (l *RemoteProvider) FetchSmiResult(req *http.Request, page, pageSize, search, order string, resultID uuid.UUID) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSMIResults) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("PersistSMIResults", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMIResults)

	l.Log.Info("attempting to fetch smi result from cloud for id: ", resultID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resultID.String()))
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed smi result url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("Retrieved result from remote provider")
		return bdr, nil
	}
	err = ErrFetch(err, "SMI Result", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetResult - fetches result from provider backend for the given result id
func (l *RemoteProvider) GetResult(tokenVal string, resultID uuid.UUID) (*MesheryResult, error) {
	if !l.Capabilities.IsSupported(PersistResult) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistResult", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistResult)

	l.Log.Info("attempting to fetch result from cloud for id: ", resultID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resultID.String()))
	l.Log.Debug("constructed result url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	// tokenString, err := l.GetToken(req)
	// if err != nil {
	// 	return nil, err
	// }
	resp, err := l.DoRequest(cReq, tokenVal)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "Perf Result "+resultID.String(), resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Result "+resultID.String())
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("result successfully retrieved from remote provider")
		res := &MesheryResult{}
		err = json.Unmarshal(bdr, res)
		if err != nil {
			err = ErrUnmarshal(err, "Perf Result "+resultID.String())
			l.Log.Error(err)
			return nil, err
		}
		return res, nil
	}
	err = ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
}

// PublishResults - publishes results to the provider backend synchronously
func (l *RemoteProvider) PublishResults(req *http.Request, result *MesheryResult, profileID string) (string, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		l.Log.Error(ErrOperationNotAvaibale)
		return "", ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	data, err := json.Marshal(result)
	if err != nil {
		return "", ErrMarshal(err, "meshery metrics for shipping")
	}

	l.Log.Debug(fmt.Sprintf("Result: %s, size: %d", data, len(data)))
	l.Log.Info("attempting to publish results to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL)
	remoteProviderURL.Path = path.Join(ep, profileID, "results")
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return "", err
	}
	l.Log.Info("request: ", cReq)
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return "", ErrUnreachableRemoteProvider(err)
		}
		err = ErrPost(err, "Perf Results", resp.StatusCode)
		l.Log.Error(err)
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return "", ErrDataRead(err, "Perf Result")
	}
	if resp.StatusCode == http.StatusCreated {
		l.Log.Info("Published results to remote provider.")
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
	err = ErrPost(err, fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return "", err
}

// PublishSmiResults - publishes results to the provider backend synchronously
func (l *RemoteProvider) PublishSmiResults(result *SmiResult) (string, error) {
	if !l.Capabilities.IsSupported(PersistSMIResults) {
		l.Log.Error(ErrOperationNotAvaibale)
		return "", ErrInvalidCapability("PersistSMIResults", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMIResults)

	data, err := json.Marshal(result)
	if err != nil {
		return "", ErrMarshal(err, "meshery metrics for shipping")
	}

	l.Log.Debug(fmt.Sprintf("Result: %s, size: %d", data, len(data)))
	l.Log.Info("attempting to publish results to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString := viper.GetString("opt-token")
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return "", ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("results successfully published to remote provider")
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

func (l *RemoteProvider) PublishEventToProvider(tokenString string, event events.Event) error {
	if !l.Capabilities.IsSupported(PersistMesheryPatternResources) {
		l.Log.Error(ErrInvalidCapability("PersistEvents", l.ProviderName))
		return ErrInvalidCapability("PersistEvents", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEvents)

	data, err := json.Marshal(event)
	if err != nil {
		return ErrMarshal(err, "meshery event")
	}

	l.Log.Info("attempting to publish event to remote provider, size: %d", len(data))
	bf := bytes.NewBuffer(data)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return ErrUnreachableRemoteProvider(err)
	}

	if resp.StatusCode != http.StatusOK {
		l.Log.Error(ErrPost(fmt.Errorf("error persisting event with the remote provider"), "event", resp.StatusCode))
		return ErrPost(fmt.Errorf("error persisting event with the remote provider"), "event", resp.StatusCode)
	}
	return nil
}

// PublishMetrics - publishes metrics to the provider backend asyncronously
func (l *RemoteProvider) PublishMetrics(tokenString string, result *MesheryResult) error {
	if !l.Capabilities.IsSupported(PersistMetrics) {
		l.Log.Error(ErrInvalidCapability("PersistMetrics", l.ProviderName))
		return ErrInvalidCapability("PersistMetrics", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMetrics)

	data, err := json.Marshal(result)
	if err != nil {
		return ErrMarshal(err, "meshery metrics for shipping")
	}

	l.Log.Debug("Result: %s, size: %d", data, len(data))
	l.Log.Info("attempting to publish metrics to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return ErrUnreachableRemoteProvider(err)
		}
		l.Log.Error(ErrPost(err, "metrics", resp.StatusCode))
		return ErrPost(err, "metrics", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusOK {
		l.Log.Info("metrics successfully published to remote provider")
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return ErrDataRead(err, "metrics Data")
	}
	l.Log.Error(ErrPost(err, fmt.Sprint(bdr), resp.StatusCode))
	return ErrPost(err, fmt.Sprint(bdr), resp.StatusCode)
}

func (l *RemoteProvider) SaveMesheryPatternResource(token string, resource *PatternResource) (*PatternResource, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatternResources) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryDesignResources", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatternResources)

	data, err := json.Marshal(resource)
	if err != nil {
		return nil, ErrMarshal(err, "meshery design resource")
	}

	l.Log.Info("attempting to save design resource to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, err := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrPost(err, "design", cReq.Response.StatusCode)
		l.Log.Error(err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		var pr PatternResource
		if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
			return nil, ErrUnmarshal(err, "Design Resource")
		}

		// l.Log.Info("design successfully sent to remote provider: %+v", pr)
		return &pr, nil
	}

	return nil, ErrPost(err, fmt.Sprint(resp.Body), resp.StatusCode)
}

func (l *RemoteProvider) GetMesheryPatternResource(token, resourceID string) (*PatternResource, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatternResources) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryPatternResources", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatternResources)

	l.Log.Info("attempting to fetch design resource from cloud for id: ", resourceID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resourceID))
	l.Log.Debug("constructed design url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Design resource", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		var pr PatternResource
		if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
			return nil, ErrUnmarshal(err, "Design resource")
		}

		l.Log.Info("Retrieved design from remote provider.")
		return &pr, nil
	}

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Design resource")
	}
	err = ErrFetch(err, fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
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
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatternResources, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatternResources)
	l.Log.Debug("Fetching designs resource from remote provider.")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
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
	l.Log.Debug("constructed design resource url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "design Page Resource", resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		var pr PatternResourcePage
		if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
			l.Log.Error(ErrDataRead(err, "respone body"))
			return nil, ErrUnmarshal(err, "design Page Resource")
		}

		l.Log.Debug("Retrieved design from remote provider")
		return &pr, nil
	}

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "design Page Resource")
	}
	err = ErrFetch(fmt.Errorf("Failed to fetch design: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

func (l *RemoteProvider) DeleteMesheryPatternResource(token, resourceID string) error {
	if !l.Capabilities.IsSupported(PersistMesheryPatternResources) {
		return ErrInvalidCapability("PersistMesheryPatternResources", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatternResources)

	l.Log.Info("Fetching design from remote provider for ID: ", resourceID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, resourceID))
	l.Log.Debug("constructed design url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return ErrUnreachableRemoteProvider(err)
		}
		return ErrDelete(err, "design: "+resourceID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("Deleted design from remote provider.")
		return nil
	}
	err = ErrDelete(fmt.Errorf("Error while deleting design."), "design: "+resourceID, resp.StatusCode)
	l.Log.Error(err)
	return err
}

func (l *RemoteProvider) SaveMesheryPatternSourceContent(tokenString string, patternID string, sourceContent []byte) error {
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	l.Log.Debug("Pattern Content size ", len(sourceContent))
	bf := bytes.NewBuffer(sourceContent)

	uploadURL := fmt.Sprintf("%s%s%s/%s", l.RemoteProviderURL, ep, remoteUploadURL, patternID)
	remoteProviderURL, _ := url.Parse(uploadURL)

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return ErrUnreachableRemoteProvider(err)
		}
		return ErrPost(err, "Pattern Source Content", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("pattern source successfully uploaded to remote provider")
		return nil
	}

	return ErrPost(fmt.Errorf("failed to upload pattern source to remote provider"), "", resp.StatusCode)
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

	l.Log.Debug(fmt.Sprintf("design: %s, size: %d", pattern.Name, len(data)))
	l.Log.Info("attempting to save design to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, err := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		err = ErrDoRequest(err, http.MethodPost, remoteProviderURL.String())
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrDoRequest(err, http.MethodPost, remoteProviderURL.String())
		l.Log.Error(err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info(fmt.Sprintf("design %s successfully sent to remote provider", pattern.Name))
		return bdr, nil
	}
	err = ErrPost(fmt.Errorf("failed to send design %s to remote provider", pattern.Name), "", resp.StatusCode)
	l.Log.Error(err)
	return bdr, err
}

// GetMesheryPatterns gives the patterns stored with the provider
func (l *RemoteProvider) GetMesheryPatterns(tokenString string, page, pageSize, search, order, updatedAfter string, visibility []string, includeMetrics string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	l.Log.Info("attempting to fetch designs from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if updatedAfter != "" {
		q.Set("updated_after", updatedAfter)
	}

	q.Set("metrics", includeMetrics)

	if len(visibility) > 0 {
		for _, v := range visibility {
			l.Log.Debug("visibility: ", v)
			q.Add("visibility", v)
		}
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed design url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "designs", resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("designs successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("unable to fetch designs"), "designs", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetCatalogMesheryPatterns gives the catalog patterns stored with the provider
func (l *RemoteProvider) GetCatalogMesheryPatterns(tokenString string, page, pageSize, search, order, includeMetrics string) ([]byte, error) {
	if !l.Capabilities.IsSupported(MesheryPatternsCatalog) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("MesheryPatternsCatalog", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(MesheryPatternsCatalog)

	l.Log.Info("attempting to fetch catalog designs from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	q.Set("metrics", includeMetrics)
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed catalog design url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "design Page - Catalog", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "design Page - Catalog")
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("catalog design successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("error while fetching catalog design: %s", bdr), "design page - Catalog", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetMesheryPattern gets pattern for the given patternID
func (l *RemoteProvider) GetMesheryPattern(req *http.Request, patternID string, includeMetrics string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryPatterns", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)
	l.Log.Info("attempting to fetch design from cloud for id: ", patternID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, patternID))
	q := remoteProviderURL.Query()
	q.Set("metrics", includeMetrics)

	l.Log.Debug("constructed design url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrFetch(err, "design:", http.StatusUnauthorized)
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "design:"+patternID, resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, ErrDataRead(err, "design:"+patternID)
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("design successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("could not retrieve design from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// DeleteMesheryPattern deletes a meshery pattern with the given id
func (l *RemoteProvider) DeleteMesheryPattern(req *http.Request, patternID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	l.Log.Info("attempting to fetch design from cloud for id: ", patternID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, patternID))
	l.Log.Debug("constructed design url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrDelete(err, "design:", http.StatusUnauthorized)
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrDelete(err, "design:"+patternID, resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("design successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrDelete(fmt.Errorf("could not retrieve design from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// CloneMesheryPattern clones a meshery pattern with the given id
func (l *RemoteProvider) CloneMesheryPattern(req *http.Request, patternID string, clonePatternRequest *MesheryClonePatternRequestBody) ([]byte, error) {
	if !l.Capabilities.IsSupported(CloneMesheryPatterns) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", CloneMesheryPatterns, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(CloneMesheryPatterns)

	l.Log.Info("attempting to clone design from cloud for id: ", patternID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, patternID))
	l.Log.Debug("constructed design url: ", remoteProviderURL.String())

	data, err := json.Marshal(clonePatternRequest)
	if err != nil {
		err = ErrMarshal(err, "design request to clone")
		l.Log.Error(err)
		return nil, err
	}

	bf := bytes.NewBuffer(data)

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrClone(err, "design")
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrClone(err, "design")
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("design successfully cloned from remote provider")
		return bdr, nil
	}
	err = ErrClone(fmt.Errorf("error while cloning design file with id %s: %s", patternID, bdr), "design")
	l.Log.Error(err)
	return nil, err
}

// PublishMesheryPattern publishes a meshery pattern with the given id to catalog
func (l *RemoteProvider) PublishCatalogPattern(req *http.Request, publishPatternRequest *MesheryCatalogPatternRequestBody) ([]byte, error) {
	if !l.Capabilities.IsSupported(MesheryPatternsCatalog) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", MesheryPatternsCatalog, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(MesheryPatternsCatalog)

	l.Log.Info("attempting to pubish design with id: ", publishPatternRequest.ID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))
	l.Log.Debug("constructed design url: ", remoteProviderURL.String())

	data, err := json.Marshal(publishPatternRequest)
	if err != nil {
		return nil, ErrMarshal(err, "design request to publish to catalog")
	}
	bf := bytes.NewBuffer(data)

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrPublish(err, "design")
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrPublish(err, "design")
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("design successfully published to catalog")
		return bdr, nil
	}
	err = ErrPublish(fmt.Errorf("error while publishing design file to catalog with id %s: %s", publishPatternRequest.ID, bdr), "design")
	l.Log.Error(err)
	return nil, err
}

// UnPublishMesheryPattern unpublishes a meshery pattern with the given id to catalog
func (l *RemoteProvider) UnPublishCatalogPattern(req *http.Request, publishPatternRequest *MesheryCatalogPatternRequestBody) ([]byte, error) {
	if !l.Capabilities.IsSupported(MesheryPatternsCatalog) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", MesheryPatternsCatalog, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(MesheryPatternsCatalog)

	l.Log.Info("attempting to unpubish design with id: ", publishPatternRequest.ID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))
	l.Log.Debug("constructed design url: ", remoteProviderURL.String())

	data, err := json.Marshal(publishPatternRequest)
	if err != nil {
		return nil, ErrMarshal(err, "design request to unpublish from catalog")
	}
	bf := bytes.NewBuffer(data)

	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), bf)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrPublish(err, "design")
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrPublish(err, "design")
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("design successfully unpublished from catalog")
		return bdr, nil
	}
	err = ErrPublish(fmt.Errorf("error while unpublishing design file from catalog with id %s: %s", publishPatternRequest.ID, bdr), "design")
	l.Log.Error(err)
	return nil, err
}

// DeleteMesheryPatterns deletes meshery patterns with the given ids and names
func (l *RemoteProvider) DeleteMesheryPatterns(req *http.Request, patterns MesheryPatternDeleteRequestBody) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", PersistMesheryPatterns, l.ProviderName)
	}

	var reqBodyBuffer bytes.Buffer
	if err := json.NewEncoder(&reqBodyBuffer).Encode(patterns); err != nil {
		err = ErrEncoding(err, "pattern delete request")
		l.Log.Error(err)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)

	// Create remote provider-url
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))
	l.Log.Debug("constructed design url: ", remoteProviderURL.String())

	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), &reqBodyBuffer)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrFetch(err, "designs", http.StatusUnauthorized)
		l.Log.Error(err)
		return nil, err
	}

	// make request
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrDoRequest(err, http.MethodDelete, remoteProviderURL.String())
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("design successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("could not retrieve design from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

func (l *RemoteProvider) RemotePatternFile(req *http.Request, resourceURL, path string, save bool) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		l.Log.Error(ErrOperationNotAvaibale)
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

	l.Log.Debug(fmt.Sprintf("design: %s, size: %d", data, len(data)))
	l.Log.Info("attempting to save design to remote provider")
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
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrPost(err, "design File", resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "design File")
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("design successfully sent to remote provider: ", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("could not send design to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// SaveMesheryFilter saves given filter with the provider
func (l *RemoteProvider) SaveMesheryFilter(tokenString string, filter *MesheryFilter) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	data, err := json.Marshal(map[string]interface{}{
		"filter_data": filter,
		"save":        true,
	})

	if err != nil {
		return nil, ErrMarshal(err, "Meshery Filters")
	}

	l.Log.Debug("size of filter: ", len(data))
	l.Log.Info("attempting to save filter to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrPost(err, "Filters", resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filters")
	}

	if resp.StatusCode == http.StatusOK {
		// l.Log.Info("filter successfully sent to remote provider: %s", string(bdr)) stop logging filter data
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("could not send filter to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// GetMesheryFilters gives the filters stored with the provider
func (l *RemoteProvider) GetMesheryFilters(tokenString string, page, pageSize, search, order string, visibility []string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	l.Log.Info("attempting to fetch filters from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}

	if len(visibility) > 0 {
		for _, v := range visibility {
			q.Add("visibility", v)
		}
	}

	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed filters url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("filters successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("error while fetching filters: %s", bdr), "Filters page", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetCatalogMesheryFilters gives the catalog filters stored with the provider
func (l *RemoteProvider) GetCatalogMesheryFilters(tokenString string, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(MesheryFiltersCatalog) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("MesheryFiltersCatalog", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(MesheryFiltersCatalog)

	l.Log.Info("attempting to fetch catalog filters from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed catalog filters url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("catalog filters successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("error while fetching catalog filters: %s", bdr), "Filters page - Catalog", resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetMesheryFilterFile gets filter for the given filterID without the metadata
func (l *RemoteProvider) GetMesheryFilterFile(req *http.Request, filterID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	l.Log.Info("attempting to fetch filter from cloud for id: ", filterID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/download/%s", l.RemoteProviderURL, ep, filterID))
	l.Log.Debug("constructed filter url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "Filter File: "+filterID, resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filter File: "+filterID)
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("filter successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("could not retrieve filter from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// GetMesheryFilter gets filter for the given filterID
func (l *RemoteProvider) GetMesheryFilter(req *http.Request, filterID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	l.Log.Info("attempting to fetch filter from cloud for id: ", filterID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, filterID))
	l.Log.Debug("constructed filter url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("filter successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("could not retrieve filter from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeleteMesheryFilter deletes a meshery filter with the given id
func (l *RemoteProvider) DeleteMesheryFilter(req *http.Request, filterID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	l.Log.Info("attempting to fetch filter from cloud for id: ", filterID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, filterID))
	l.Log.Debug("constructed filter url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrDelete(err, "Filter: "+filterID, resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Filter: "+filterID)
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("filter successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrDelete(fmt.Errorf("error while fetching filter: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// CloneMesheryFilter clones a meshery filter with the given id
func (l *RemoteProvider) CloneMesheryFilter(req *http.Request, filterID string, cloneFilterRequest *MesheryCloneFilterRequestBody) ([]byte, error) {
	if !l.Capabilities.IsSupported(CloneMesheryFilters) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", CloneMesheryFilters, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(CloneMesheryFilters)

	l.Log.Info("attempting to clone filter from cloud for id: ", filterID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, filterID))
	l.Log.Debug("constructed filter url: ", remoteProviderURL.String())
	data, err := json.Marshal(cloneFilterRequest)
	if err != nil {
		err = ErrMarshal(err, "filter request to clone")
		l.Log.Error(err)
		return nil, err
	}

	bf := bytes.NewBuffer(data)

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrClone(err, "filter")
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrClone(err, "filter")
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("filter successfully cloned from remote provider")
		return bdr, nil
	}
	err = ErrClone(fmt.Errorf("error while cloning filter with id %s: %s", filterID, bdr), "filter")
	l.Log.Error(err)
	return nil, err
}

// CloneMesheryFilter publishes a meshery filter with the given id to catalog
func (l *RemoteProvider) PublishCatalogFilter(req *http.Request, publishFilterRequest *MesheryCatalogFilterRequestBody) ([]byte, error) {
	if !l.Capabilities.IsSupported(MesheryFiltersCatalog) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", MesheryFiltersCatalog, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(MesheryFiltersCatalog)

	l.Log.Info("attempting to pubish filter with id: ", publishFilterRequest.ID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))
	l.Log.Debug("constructed filter url: ", remoteProviderURL.String())

	data, err := json.Marshal(publishFilterRequest)
	if err != nil {
		return nil, ErrMarshal(err, "filter request to publish to catalog")
	}
	bf := bytes.NewBuffer(data)

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrPublish(err, "filter")
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrPublish(err, "filter")
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("filter successfully published to catalog")
		return bdr, nil
	}
	err = ErrPublish(fmt.Errorf("error while publishing filter file to catalog with id %s: %s", publishFilterRequest.ID, bdr), "filter")
	l.Log.Error(err)
	return nil, err
}

// UnPublishMesheryFilter publishes a meshery filter with the given id to catalog
func (l *RemoteProvider) UnPublishCatalogFilter(req *http.Request, publishFilterRequest *MesheryCatalogFilterRequestBody) ([]byte, error) {
	if !l.Capabilities.IsSupported(MesheryFiltersCatalog) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, fmt.Errorf("%s is not suppported by provider: %s", MesheryFiltersCatalog, l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(MesheryFiltersCatalog)

	l.Log.Info("attempting to unpubish filter with id: ", publishFilterRequest.ID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))
	l.Log.Debug("constructed filter url: ", remoteProviderURL.String())

	data, err := json.Marshal(publishFilterRequest)
	if err != nil {
		return nil, ErrMarshal(err, "filter request to unpublish from catalog")
	}
	bf := bytes.NewBuffer(data)

	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), bf)

	tokenString, err := l.GetToken(req)
	if err != nil {
		err = ErrUnpPublish(err, "filter")
		l.Log.Error(err)
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrUnpPublish(err, "filter")
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("filter successfully unpublished from catalog")
		return bdr, nil
	}
	err = ErrUnpPublish(fmt.Errorf("error while unpublishing filter file from catalog with id %s: %s", publishFilterRequest.ID, bdr), "filter")
	l.Log.Error(err)
	return nil, err
}

func (l *RemoteProvider) RemoteFilterFile(req *http.Request, resourceURL, path string, save bool, resource string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryFilters) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryFilters)

	data, err := json.Marshal(map[string]interface{}{
		"url":  resourceURL,
		"save": save,
		"path": path,
		"filter_data": MesheryFilter{
			FilterResource: resource,
		},
	})

	if err != nil {
		err = ErrMarshal(err, "meshery metrics for shipping")
		return nil, err
	}

	l.Log.Debug(fmt.Sprintf("Filter: %s, size: %d", data, len(data)))
	l.Log.Info("attempting to save filter to remote provider")
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
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		// l.Log.Info("filter successfully sent to remote provider: %s", string(bdr)) stop logging filter data
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("could not send filter to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// SaveMesheryApplication saves given application with the provider
func (l *RemoteProvider) SaveMesheryApplication(tokenString string, application *MesheryApplication) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryApplications) {
		l.Log.Error(ErrOperationNotAvaibale)
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

	l.Log.Debug(fmt.Sprintf("Application size: %d", len(data)))
	l.Log.Info("attempting to save application to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("application successfully sent to remote provider: ", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("failed to send application to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// SaveApplicationSourceContent saves given application source content with the provider after successful save of Application with the provider
func (l *RemoteProvider) SaveApplicationSourceContent(tokenString string, applicationID string, sourceContent []byte) error {
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	l.Log.Debug("Application Content size ", len(sourceContent))
	bf := bytes.NewBuffer(sourceContent)

	uploadURL := fmt.Sprintf("%s%s%s/%s", l.RemoteProviderURL, ep, remoteUploadURL, applicationID)
	remoteProviderURL, _ := url.Parse(uploadURL)

	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return ErrUnreachableRemoteProvider(err)
		}
		return ErrPost(err, "Application Source Content", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("application source successfully uploaded to remote provider")
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

	l.Log.Info("attempting to fetch application source content from cloud for id: ", applicationID)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "Application source content", resp.StatusCode)
		l.Log.Error((err))
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, ErrDataRead(err, "Application")
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("applications successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("error while fetching applications: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetDesignSourceContent returns design source-content from provider
func (l *RemoteProvider) GetDesignSourceContent(token, designID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryPatterns) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryPatterns", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryPatterns)
	downloadURL := fmt.Sprintf("%s%s%s/%s", l.RemoteProviderURL, ep, remoteDownloadURL, designID)
	remoteProviderURL, _ := url.Parse(downloadURL)
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	l.Log.Info("attempting to fetch design source content from cloud for id: ", designID)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "Design source content", resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, ErrDataRead(err, "Pattern")
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("design source content successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("error while fetching designs: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)

	return nil, err
}

// GetMesheryApplications gives the applications stored with the provider
func (l *RemoteProvider) GetMesheryApplications(tokenString string, page, pageSize, search, order string, updaterAfter string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryApplications) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("PersistMesheryApplications", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	l.Log.Info("attempting to fetch applications from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if updaterAfter != "" {
		q.Set("updated_after", updaterAfter)
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed applications url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "Application", resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, ErrDataRead(err, "Application")
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("applications successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("error while fetching applications: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetMesheryApplication gets application for the given applicationID
func (l *RemoteProvider) GetMesheryApplication(req *http.Request, applicationID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryApplications) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryApplications", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	l.Log.Info("attempting to fetch application from cloud for id: ", applicationID)
	urls := fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, applicationID)
	remoteProviderURL, _ := url.Parse(urls)
	l.Log.Debug("constructed application url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("application successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to retrieve application from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeleteMesheryApplication deletes a meshery application with the given id
func (l *RemoteProvider) DeleteMesheryApplication(req *http.Request, applicationID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistMesheryApplications) {
		l.Log.Error(ErrOperationNotAvaibale)
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistMesheryApplications", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistMesheryApplications)

	l.Log.Info("attempting to fetch application from cloud for id: ", applicationID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, applicationID))
	l.Log.Debug("constructed application url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("application successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrDelete(fmt.Errorf("could not retrieve application from remote provider"), "Application :"+applicationID, resp.StatusCode)
}

func (l *RemoteProvider) ShareDesign(req *http.Request) (int, error) {
	if !l.Capabilities.IsSupported(ShareDesigns) {
		l.Log.Error(ErrOperationNotAvaibale)
		return http.StatusForbidden, ErrInvalidCapability("ShareDesigns", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(ShareDesigns)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))
	bd, _ := io.ReadAll(req.Body)
	defer func() {
		_ = req.Body.Close()
	}()
	bf := bytes.NewBuffer(bd)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return http.StatusInternalServerError, ErrUnreachableRemoteProvider(err)
		}
		return http.StatusInternalServerError, ErrShareDesign(err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, ErrShareDesign(fmt.Errorf("unable to share design"))
	}
	return resp.StatusCode, nil
}

// SavePerformanceProfile saves a performance profile into the remote provider
func (l *RemoteProvider) SavePerformanceProfile(tokenString string, pp *PerformanceProfile) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	data, err := json.Marshal(pp)
	if err != nil {
		err = ErrMarshal(err, "meshery metrics for shipping")
		return nil, err
	}

	l.Log.Debug(fmt.Sprintf("performance profile: %s, size: %d", data, len(data)))
	l.Log.Info("attempting to save performance profile to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("performance profile successfully sent to remote provider: ", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("failed to send performance profile to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// GetPerformanceProfiles gives the performance profiles stored with the provider
func (l *RemoteProvider) GetPerformanceProfiles(tokenString string, page, pageSize, search, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	l.Log.Info("attempting to fetch performance profiles from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed performance profiles url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("performance profiles successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrPost(fmt.Errorf("failed to retrieve performance profile from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// GetPerformanceProfile gets performance profile for the given the performanceProfileID
func (l *RemoteProvider) GetPerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	l.Log.Info("attempting to fetch performance profile from cloud for id: ", performanceProfileID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, performanceProfileID))
	l.Log.Debug("constructed performance profile url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrFetch(err, "Perf Profile :"+performanceProfileID, resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Profile :"+performanceProfileID)
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("performance profile successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to retrieve performance profile from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeletePerformanceProfile deletes a performance profile with the given performanceProfileID
func (l *RemoteProvider) DeletePerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistPerformanceProfiles) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistPerformanceProfiles", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistPerformanceProfiles)

	l.Log.Info("attempting to fetch performance profile from cloud for id: ", performanceProfileID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, performanceProfileID))
	l.Log.Debug("constructed performance profile url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("performance profile successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrDelete(fmt.Errorf("failed to retrieve performance profile from remote provider"), "Perf Profile :"+performanceProfileID, resp.StatusCode)
}

// SaveSchedule saves a SaveSchedule into the remote provider
func (l *RemoteProvider) SaveSchedule(tokenString string, s *Schedule) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistSchedules", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	data, err := json.Marshal(s)
	if err != nil {
		return nil, ErrMarshal(err, "schedule for shipping")
	}

	l.Log.Debug(fmt.Sprintf("schedule: %s, size: %d", data, len(data)))
	l.Log.Info("attempting to save schedule to remote provider")
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrPost(err, "Perf Schedule", resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Schedule")
	}

	if resp.StatusCode == http.StatusCreated {
		l.Log.Info("schedule successfully sent to remote provider: ", string(bdr))
		return bdr, nil
	}

	return bdr, ErrPost(fmt.Errorf("failed to send schedule to remote provider: %s", string(bdr)), fmt.Sprint(bdr), resp.StatusCode)
}

// GetSchedules gives the schedules stored with the provider
func (l *RemoteProvider) GetSchedules(req *http.Request, page, pageSize, order string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("PersistSchedules", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	l.Log.Info("attempting to fetch schedules from cloud")

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if order != "" {
		q.Set("order", order)
	}
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed schedules url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("schedules successfully retrieved from remote provider")
		return bdr, nil
	}
	err = ErrFetch(fmt.Errorf("error while fetching schedules: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

// GetSchedule gets schedule for the given the scheduleID
func (l *RemoteProvider) GetSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistSchedules", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	l.Log.Info("attempting to fetch schedule from cloud for id: ", scheduleID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, scheduleID))
	l.Log.Debug("constructed schedule url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Info("schedule successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("could not retrieve schedule from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeleteSchedule deletes a schedule with the given scheduleID
func (l *RemoteProvider) DeleteSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistSchedules) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistSchedules", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSchedules)

	l.Log.Info("attempting to fetch schedule from cloud for id: ", scheduleID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, scheduleID))
	l.Log.Debug("constructed schedule url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		err = ErrDelete(err, "Perf Schedule :"+scheduleID, resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Schedule :"+scheduleID)
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("schedule successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrDelete(fmt.Errorf("could not retrieve schedule from remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

// RecordPreferences - records the user preference
func (l *RemoteProvider) RecordPreferences(req *http.Request, userID string, data *Preference) error {
	if !l.Capabilities.IsSupported(SyncPrefs) {
		l.Log.Error(ErrOperationNotAvaibale)
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
func (l *RemoteProvider) TokenHandler(w http.ResponseWriter, r *http.Request, _ bool) {
	tokenString := r.URL.Query().Get(TokenCookieName)
	// gets the session cookie from remote provider
	sessionCookie := r.URL.Query().Get("session_cookie")

	l.SetJWTCookie(w, tokenString)
	// sets the session cookie for Meshery Session
	l.SetProviderSessionCookie(w, sessionCookie)

	// Get new capabilities
	// Doing this here is important so that
	providerProperties := l.loadCapabilities(tokenString)
	l.ProviderProperties = providerProperties

	// Download the package for the user only if they have extension capability
	if len(l.GetProviderProperties().Extensions.Navigator) > 0 {
		l.DownloadProviderExtensionPackage(l.Log)
	}

	// Proceed to redirect once the capabilities has loaded
	// and the package has been downloaded
	redirectURL := "/"
	isPlayGround, _ := strconv.ParseBool(viper.GetString("PLAYGROUND"))
	if isPlayGround {
		redirectURL = GetRedirectURLForNavigatorExtension(&providerProperties)
	}

	refQueryParam := r.URL.Query().Get("ref")
	if refQueryParam != "" {
		redirectURL = refQueryParam
	}

	go func() {
		credential := make(map[string]interface{}, 0)
		var temp *uuid.UUID
		credential["token"] = temp

		connectionPayload := connections.BuildMesheryConnectionPayload(r.Context().Value(MesheryServerURL).(string), credential)

		_, err := l.SaveConnection(connectionPayload, tokenString, true)
		if err != nil {
			l.Log.Error(ErrSaveConnection(err))
		}
	}()
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// UpdateToken - in case the token was refreshed, this routine updates the response with the new token
func (l *RemoteProvider) UpdateToken(w http.ResponseWriter, r *http.Request) string {
	l.TokenStoreMut.Lock()
	defer l.TokenStoreMut.Unlock()

	tokenString, _ := l.GetToken(r)
	newts := l.TokenStore[tokenString]
	if newts != "" {
		l.Log.Debug("set updated token: ", newts)
		l.SetJWTCookie(w, newts)
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
		l.Log.Error(ErrGetToken(err))
		return
	}
	newts := l.TokenStore[tokenString]
	if newts != "" {
		tokenString = newts
	}

	resp := map[string]interface{}{
		"meshery-provider": l.Name(),
		TokenCookieName:    tokenString,
	}
	l.Log.Debug("token sent for meshery-provider ", l.Name())
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		err = ErrEncoding(err, "Auth Details")
		l.Log.Error(ErrEncoding(err, "Auth Details"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SMPTestConfigStore - persist test profile details to provider
func (l *RemoteProvider) SMPTestConfigStore(req *http.Request, perfConfig *SMP.PerformanceTestConfig) (string, error) {
	if !l.Capabilities.IsSupported(PersistSMPTestProfile) {
		l.Log.Error(ErrOperationNotAvaibale)
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
		if resp == nil {
			return "", ErrUnreachableRemoteProvider(err)
		}
		err = ErrPost(err, "Perf Test Config", resp.StatusCode)
		l.Log.Error(err)

		return "", err
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
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistSMPTestProfile", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMPTestProfile)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	q.Add("test_uuid", testUUID)
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("Making request to : ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Perf Test Config: "+testUUID, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Perf Test Config :"+testUUID)
	}
	l.Log.Debug(string(bdr))
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
		l.Log.Error(ErrOperationNotAvaibale)
		return []byte{}, ErrInvalidCapability("PersistSMPTestProfile", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistSMPTestProfile)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	q.Add("page", page)
	q.Add("pagesize", pageSize)
	q.Add("search", search)
	q.Add("order", order)

	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("Making request to : ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
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
		l.Log.Error(ErrOperationNotAvaibale)
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
		if resp == nil {
			return ErrUnreachableRemoteProvider(err)
		}
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

// ExtensionProxy - proxy requests to the remote provider which are specific to user_account extension
func (l *RemoteProvider) ExtensionProxy(req *http.Request) (*ExtensionProxyResponse, error) {
	l.Log.Info("attempting to request remote provider")
	// gets the requested path from user_account extension UI in Meshery UI
	// splits the requested path into '/api/extensions' and '/<remote-provider-endpoint>'
	p := req.URL.Path
	split := strings.Split(p, "/api/extensions")
	path := split[1]
	// gets the available query parameters
	q := req.URL.Query().Encode()
	if len(q) > 0 {
		// if available, then add it to <remote-provider-endpoint>
		// eg: /<remote-provider-endpoint>?<query-parameters>
		path = fmt.Sprintf("%s?%s", path, q)
	}
	// then attach the final path to the remote provider URL
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, path))
	l.Log.Debug("constructed url: ", remoteProviderURL.String())

	// make http.Request type variable with the constructed URL
	cReq, err := http.NewRequest(req.Method, remoteProviderURL.String(), req.Body)
	if err != nil {
		return nil, err
	}

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	// update the request headers with the headers from the original request
	// original headers includes cookies likes session_cookie necessary to run user management flows
	cReq.Header = req.Header

	// make request to remote provider with contructed URL and updated headers (like session cookie)
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &ExtensionProxyResponse{
		Body:       bdr,
		StatusCode: resp.StatusCode,
	}

	// check for all success status codes
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if statusOK {
		l.Log.Info("response successfully retrieved from remote provider")
		return response, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to request to remote provider"), fmt.Sprint(bdr), resp.StatusCode)
}

func (l *RemoteProvider) SaveConnection(conn *connections.ConnectionPayload, token string, skipTokenCheck bool) (*connections.Connection, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	_conn, err := json.Marshal(conn)
	if err != nil {
		return nil, err
	}
	bf := bytes.NewBuffer(_conn)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	l.Log.Info(fmt.Sprintf("attempting to save %s connection %s to remote provider with status %s", conn.Name, conn.Kind, conn.Status))
	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Save Connection", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Save Connection")
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		connectionPage := &connections.ConnectionPage{}
		err = json.Unmarshal(bdr, connectionPage)
		if err != nil {
			return nil, ErrUnmarshal(err, "Connection \"%s\" of type \"%s\" with status \"%s\" from the remote provider")
		}
		l.Log.Debug("connections, ", connectionPage)
		// On POST request to Remote Provider API, the response always contains single entry/connection.
		if len(connectionPage.Connections) > 0 {
			return connectionPage.Connections[0], nil
		}
		return nil, ErrPost(fmt.Errorf("failed to save the connection"), fmt.Sprint(bdr), resp.StatusCode)
	}

	return nil, ErrPost(fmt.Errorf("failed to save the connection \"%s\" of type \"%s\" with status \"%s\"", conn.Name, conn.Kind, conn.Status), fmt.Sprint(bdr), resp.StatusCode)
}

func (l *RemoteProvider) GetConnections(req *http.Request, userID string, page, pageSize int, search, order string, filter string, status []string, kind []string) (*connections.ConnectionPage, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))

	q := remoteProviderURL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("pagesize", strconv.Itoa(pageSize))
	q.Add("search", search)
	q.Add("order", order)
	if filter != "" {
		q.Set("filter", filter)
	}

	if len(status) > 0 {
		for _, v := range status {
			q.Add("status", v)
		}
	}
	if len(kind) > 0 {
		for _, v := range kind {
			q.Add("kind", v)
		}
	}

	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("Making request to : ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Connections Page", resp.StatusCode)
	}
	defer resp.Body.Close()

	bdr, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, ErrFetch(fmt.Errorf("could not retrieve list of connections: %d", resp.StatusCode), fmt.Sprint(bdr), resp.StatusCode)
	}

	var cp connections.ConnectionPage
	if err = json.Unmarshal(bdr, &cp); err != nil {
		return nil, err
	}

	return &cp, nil
}

// GetConnectionsByKind - to get saved credentials
func (l *RemoteProvider) GetConnectionsByKind(req *http.Request, _ string, page, pageSize int, search, order, connectionKind string) (*map[string]interface{}, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, connectionKind))

	q := remoteProviderURL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("pagesize", strconv.Itoa(pageSize))
	q.Add("search", search)
	q.Add("order", order)

	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("Making request to : ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Connections Page", resp.StatusCode)
	}
	defer resp.Body.Close()

	bdr, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, ErrFetch(fmt.Errorf("%s", string(bdr)), "connections", resp.StatusCode)
	}

	var res map[string]interface{}
	if err = json.Unmarshal(bdr, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (l *RemoteProvider) GetConnectionByIDAndKind(token string, connectionID uuid.UUID, kind string) (*connections.Connection, int, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, http.StatusForbidden, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s/%s", l.RemoteProviderURL, ep, kind, connectionID))

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		l.Log.Error(err)
		statusCode := http.StatusInternalServerError
		if resp != nil {
			statusCode = resp.StatusCode
		}
		return nil, http.StatusInternalServerError, ErrFetch(err, "connection", statusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, ErrFetch(err, "connection", http.StatusInternalServerError)
	}
	if resp.StatusCode == http.StatusOK {
		var conn connections.Connection
		if err = json.Unmarshal(bdr, &conn); err != nil {
			return nil, http.StatusInternalServerError, ErrUnmarshal(err, "connection")
		}
		return &conn, resp.StatusCode, nil
	}

	l.Log.Debug(string(bdr))
	return nil, resp.StatusCode, ErrFetch(fmt.Errorf("unable to retrieve connection with id %s", connectionID), "connection", resp.StatusCode)
}

func (l *RemoteProvider) GetConnectionByID(token string, connectionID uuid.UUID) (*connections.Connection, int, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrInvalidCapability("PersistConnection", l.ProviderName))
		return nil, http.StatusForbidden, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, connectionID))

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		l.Log.Error(err)
		statusCode := http.StatusInternalServerError
		return nil, statusCode, ErrFetch(err, "connection", statusCode)
	}

	if resp.StatusCode == http.StatusOK {
		defer func() {
			_ = resp.Body.Close()
		}()

		bdr, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, ErrFetch(fmt.Errorf("unable to retrieve connection with id %s", connectionID), "connection", resp.StatusCode)
		}
		var conn connections.Connection
		if err = json.Unmarshal(bdr, &conn); err != nil {
			l.Log.Error(ErrUnmarshal(err, "connection"))
			return nil, http.StatusInternalServerError, ErrFetch(fmt.Errorf("unable to retrieve connection with id %s", connectionID), "connection", resp.StatusCode)
		}
		return &conn, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, ErrFetch(fmt.Errorf("unable to retrieve connection with id %s", connectionID), "connection", resp.StatusCode)
}

func (l *RemoteProvider) GetConnectionsStatus(req *http.Request, userID string) (*connections.ConnectionsStatusPage, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/status", l.RemoteProviderURL, ep))

	l.Log.Debug("Making request to : ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if resp != nil {
			statusCode = resp.StatusCode
		}
		return nil, ErrFetch(err, "Connections Status Page", statusCode)
	}
	defer resp.Body.Close()

	bdr, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, ErrFetch(fmt.Errorf("could not retrieve list of connections status: %d", resp.StatusCode), fmt.Sprint(bdr), resp.StatusCode)
	}

	var cp connections.ConnectionsStatusPage
	if err = json.Unmarshal(bdr, &cp); err != nil {
		return nil, err
	}
	return &cp, nil
}

func (l *RemoteProvider) UpdateConnectionStatusByID(token string, connectionID uuid.UUID, connectionStatus connections.ConnectionStatus) (*connections.Connection, int, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, http.StatusForbidden, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	bf := bytes.NewBuffer([]byte(connectionStatus))
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/status/%s", l.RemoteProviderURL, ep, connectionID))

	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		l.Log.Error(err)
		return nil, http.StatusInternalServerError, ErrUpdateConnectionStatus(err, http.StatusInternalServerError)
	}
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, ErrDataRead(err, "Update Connection")
	}
	if resp.StatusCode == http.StatusOK {
		var conn connections.Connection
		if err = json.Unmarshal(bdr, &conn); err != nil {
			return nil, http.StatusInternalServerError, ErrUnmarshal(err, "connection")
		}
		return &conn, resp.StatusCode, nil
	}

	l.Log.Debug(string(bdr))
	return nil, resp.StatusCode, ErrUpdateConnectionStatus(fmt.Errorf("unable to update connection with id %s", connectionID), resp.StatusCode)
}

// UpdateConnection - to update an existing connection
func (l *RemoteProvider) UpdateConnection(req *http.Request, connection *connections.Connection) (*connections.Connection, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	_creds, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}
	bf := bytes.NewBuffer(_creds)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, connection.Kind))
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)
	tokenString, _ := l.GetToken(req)
	if err != nil {
		l.Log.Error(ErrGetToken(err))
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Update Connection", resp.StatusCode)
	}
	defer resp.Body.Close()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Update Connection")
	}

	if resp.StatusCode == http.StatusOK {
		var conn connections.Connection
		if err = json.Unmarshal(bdr, &conn); err != nil {
			return nil, err
		}
		return &conn, nil
	}

	return nil, ErrFetch(fmt.Errorf("failed to update the connection"), fmt.Sprint(bdr), resp.StatusCode)
}

// UpdateConnectionById - to update an existing connection using the connection id
func (l *RemoteProvider) UpdateConnectionById(req *http.Request, connection *connections.ConnectionPayload, connId string) (*connections.Connection, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	_conn, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}
	bf := bytes.NewBuffer(_conn)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, connId))
	l.Log.Debug("Making request to : ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)
	tokenString, err := l.GetToken(req)
	if err != nil {
		l.Log.Error(ErrGetToken(err))
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Update Connection", resp.StatusCode)
	}
	defer resp.Body.Close()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Update Connection")
	}

	if resp.StatusCode == http.StatusOK {
		var conn connections.Connection
		if err = json.Unmarshal(bdr, &conn); err != nil {
			return nil, ErrUnmarshal(err, "connection")
		}
		return &conn, nil
	}

	return nil, ErrFetch(fmt.Errorf("failed to update the connection"), string(bdr), resp.StatusCode)
}

// DeleteConnection - to delete a saved connection
func (l *RemoteProvider) DeleteConnection(req *http.Request, connectionID uuid.UUID) (*connections.Connection, error) {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistConnection", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)

	l.Log.Info("attempting to delete connection from cloud for id: ", connectionID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, connectionID))
	l.Log.Debug("constructed connection url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		err = ErrDelete(err, "Connection: "+connectionID.String(), resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Connection: "+connectionID.String())
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("connection successfully deleted from remote provider")
		var conn connections.Connection
		if err = json.Unmarshal(bdr, &conn); err != nil {
			return nil, err
		}
		return &conn, nil
	}
	err = ErrDelete(fmt.Errorf("error while deleting connection: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

func (l *RemoteProvider) DeleteMesheryConnection() error {
	if !l.Capabilities.IsSupported(PersistConnection) {
		l.Log.Error(ErrOperationNotAvaibale)
		return ErrInvalidCapability("PersistConnection", l.ProviderName)
	}

	mesheryServerID := viper.GetString("INSTANCE_ID")
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistConnection)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/meshery/%s", l.RemoteProviderURL, ep, mesheryServerID))
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)
	cReq.Header.Set("X-API-Key", GlobalTokenForAnonymousResults)
	cReq.Header.Set("SystemID", viper.GetString("INSTANCE_ID")) // Adds the system id to the header for event tracking
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		if resp == nil {
			return ErrUnreachableRemoteProvider(err)
		}
		return ErrDelete(err, "Meshery Connection", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return ErrDelete(fmt.Errorf("could not delete meshery connection"), " Meshery Connection", resp.StatusCode)
}

// TarXZF takes in a source url downloads the tar.gz file
// uncompresses and then save the file to the destination
func TarXZF(srcURL, destination string, log logger.Handler) error {
	filename := filepath.Base(srcURL)

	// Check if filename ends with tar.gz or tgz extension
	if !strings.HasSuffix(filename, ".tar.gz") && !strings.HasSuffix(filename, ".tgz") {
		return fmt.Errorf("file is not of type tar.gz or tgz")
	}

	resp, err := http.Get(srcURL)
	if err != nil {
		if resp == nil {
			return fmt.Errorf("could not reach %v: %w", srcURL, err)
		}
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(ErrCloseIoReader(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return ErrFetch(fmt.Errorf("failed GET request"), "TarTZF file :"+srcURL, resp.StatusCode)
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

// SaveCredential - to save a creadential for an integration
func (l *RemoteProvider) SaveUserCredential(token string, credential *Credential) (*Credential, error) {
	var createdCredential Credential
	if !l.Capabilities.IsSupported(PersistCredentials) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistCredentials", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistCredentials)
	_creds, err := json.Marshal(credential)
	if err != nil {
		return nil, err
	}
	bf := bytes.NewBuffer(_creds)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		return nil, ErrFetch(err, "Save Credential", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Save Credential")
	}

	if resp.StatusCode == http.StatusCreated {
		err := json.Unmarshal(bdr, &createdCredential)
		if err != nil {
			return nil, ErrUnmarshal(err, "created credential")
		}
		return &createdCredential, nil
	}

	return nil, ErrFetch(fmt.Errorf("failed to save the credential"), fmt.Sprint(bdr), resp.StatusCode)
}

// GetCredentials - to get saved credentials
func (l *RemoteProvider) GetUserCredentials(req *http.Request, _ string, page, pageSize int, search, order string) (*CredentialsPage, error) {
	if !l.Capabilities.IsSupported(PersistCredentials) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistCredentials", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistCredentials)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("pagesize", strconv.Itoa(pageSize))
	q.Add("search", search)
	q.Add("order", order)

	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("Making request to : ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Credentials Page", resp.StatusCode)
	}
	defer resp.Body.Close()

	bdr, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, ErrFetch(fmt.Errorf("could not retrieve list of test profiles: %d", resp.StatusCode), fmt.Sprint(bdr), resp.StatusCode)
	}

	var cp CredentialsPage
	if err = json.Unmarshal(bdr, &cp); err != nil {
		return nil, ErrFetch(err, "Unmarshal Credentials Page", resp.StatusCode)
	}
	return &cp, nil
}

func (l *RemoteProvider) GetCredentialByID(token string, credentialID uuid.UUID) (*Credential, int, error) {
	if !l.Capabilities.IsSupported(PersistCredentials) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, http.StatusForbidden, ErrInvalidCapability("PersistCredentials", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistCredentials)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s/%s", l.RemoteProviderURL, ep, credentialID))
	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if resp != nil {
			statusCode = resp.StatusCode
		}
		return nil, statusCode, ErrFetch(err, "Credentials Page", resp.StatusCode)
	}
	defer resp.Body.Close()

	bdr, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, ErrFetch(fmt.Errorf("could not retrieve credential with id %s: %d", credentialID, resp.StatusCode), fmt.Sprint(bdr), resp.StatusCode)
	}

	var cp Credential
	if err = json.Unmarshal(bdr, &cp); err != nil {
		return nil, http.StatusInternalServerError, ErrFetch(err, "Unmarshal Credentials Page", resp.StatusCode)
	}
	return &cp, resp.StatusCode, nil
}

// UpdateUserCredential - to update an existing credential
func (l *RemoteProvider) UpdateUserCredential(req *http.Request, credential *Credential) (*Credential, error) {
	if !l.Capabilities.IsSupported(PersistCredentials) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistCredentials", l.ProviderName)
	}
	ep, _ := l.Capabilities.GetEndpointForFeature(PersistCredentials)
	_creds, err := json.Marshal(credential)
	if err != nil {
		return nil, err
	}
	bf := bytes.NewBuffer(_creds)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)
	tokenString, _ := l.GetToken(req)
	if err != nil {
		l.Log.Error(ErrGetToken(err))
		return nil, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return nil, ErrFetch(err, "Update Credential", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Update Credential")
	}

	if resp.StatusCode == http.StatusOK {
		var cred Credential
		if err = json.Unmarshal(bdr, &cred); err != nil {
			return nil, err
		}
		return &cred, nil
	}

	return nil, ErrFetch(fmt.Errorf("failed to update the credential"), fmt.Sprint(bdr), resp.StatusCode)
}

// DeleteUserCredential - to delete a saved credential
func (l *RemoteProvider) DeleteUserCredential(req *http.Request, credentialID uuid.UUID) (*Credential, error) {
	if !l.Capabilities.IsSupported(PersistCredentials) {
		l.Log.Error(ErrOperationNotAvaibale)
		return nil, ErrInvalidCapability("PersistCredentials", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistCredentials)

	l.Log.Info("attempting to delete credential from cloud for id: ", credentialID)

	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))
	q := remoteProviderURL.Query()
	q.Add("credential_id", credentialID.String())
	remoteProviderURL.RawQuery = q.Encode()
	l.Log.Debug("constructed credential url: ", remoteProviderURL.String())
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)

	tokenString, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		err = ErrDelete(err, "Credential: "+credentialID.String(), resp.StatusCode)
		l.Log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Credential: "+credentialID.String())
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("credential successfully deleted from remote provider")
		return nil, nil
	}
	err = ErrDelete(fmt.Errorf("error while deleting credential: %s", bdr), fmt.Sprint(bdr), resp.StatusCode)
	l.Log.Error(err)
	return nil, err
}

func (l *RemoteProvider) ShareFilter(req *http.Request) (int, error) {
	if !l.Capabilities.IsSupported(ShareFilters) {
		l.Log.Error(ErrOperationNotAvaibale)
		return http.StatusForbidden, ErrInvalidCapability("ShareFilters", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(ShareFilters)
	remoteProviderURL, _ := url.Parse(fmt.Sprintf("%s%s", l.RemoteProviderURL, ep))
	bd, _ := io.ReadAll(req.Body)
	defer func() {
		_ = req.Body.Close()
	}()
	bf := bytes.NewBuffer(bd)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString, err := l.GetToken(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		return http.StatusInternalServerError, ErrShareFilter(err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, ErrShareFilter(fmt.Errorf("unable to share filter"))
	}
	return resp.StatusCode, nil
}

func (l *RemoteProvider) GetEnvironments(token, page, pageSize, search, order, filter, orgID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistEnvironments) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Environment", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEnvironments)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	if orgID != "" {
		q.Set("orgID", orgID)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Users Data", http.StatusUnauthorized)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Environments")
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Environments data successfully retrieved from remote provider")
		return bd, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get environments"), "Environments", resp.StatusCode)
}

func (l *RemoteProvider) GetEnvironmentByID(req *http.Request, environmentID, orgID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistEnvironments) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Environment", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEnvironments)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + environmentID)
	q := remoteProviderURL.Query()
	if orgID != "" {
		q.Set("orgID", orgID)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Environment", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("Environment successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get environment by ID"), "Environment", resp.StatusCode)
}

func (l *RemoteProvider) SaveEnvironment(req *http.Request, env *v1beta1.EnvironmentPayload, token string, skipTokenCheck bool) ([]byte, error) {

	if !l.Capabilities.IsSupported(PersistEnvironments) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Environment", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEnvironments)
	_env, err := json.Marshal(env)
	if err != nil {
		return []byte{}, err
	}
	bf := bytes.NewBuffer(_env)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString := token
	if !skipTokenCheck {
		tokenString, err = l.GetToken(req)
		if err != nil {
			l.Log.Error(ErrGetToken(err))
			return []byte{}, err
		}
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return []byte{}, ErrUnreachableRemoteProvider(err)
		}
		return []byte{}, ErrFetch(err, "Save Environment", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBf, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, ErrDataRead(err, "Save Environment")
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return respBf, nil
	}

	return []byte{}, ErrPost(fmt.Errorf("failed to save the environment"), "Environment", resp.StatusCode)
}

func (l *RemoteProvider) DeleteEnvironment(req *http.Request, environmentID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistEnvironments) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Environment", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEnvironments)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + environmentID)
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Environment", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Environment successfully deleted from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to delete environment"), "Environment", resp.StatusCode)
}

func (l *RemoteProvider) UpdateEnvironment(req *http.Request, env *v1beta1.EnvironmentPayload, environmentID string) (*v1beta1.Environment, error) {
	if !l.Capabilities.IsSupported(PersistEnvironments) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return &v1beta1.Environment{}, ErrInvalidCapability("Environment", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEnvironments)
	_env, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}
	bf := bytes.NewBuffer(_env)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + environmentID)
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Environment", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Update Environment")
	}

	if resp.StatusCode == http.StatusOK {
		var environment v1beta1.Environment
		if err = json.Unmarshal(bdr, &environment); err != nil {
			return nil, err
		}
		return &environment, nil
	}

	return nil, ErrFetch(fmt.Errorf("failed to update the environment"), "Environment", resp.StatusCode)
}

func (l *RemoteProvider) AddConnectionToEnvironment(req *http.Request, environmentID string, connectionID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistEnvironments) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Environment", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEnvironments)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + environmentID + "/connections/" + connectionID)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Environment", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Connection successfully added to environment")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get environments"), "Environment", resp.StatusCode)
}

func (l *RemoteProvider) RemoveConnectionFromEnvironment(req *http.Request, environmentID string, connectionID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistEnvironments) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Environment", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEnvironments)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + environmentID + "/connections/" + connectionID)
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Environment", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Connection successfully removed from environment")
		return bdr, nil
	}

	return nil, ErrFetch(fmt.Errorf("failed to unassign connection from environment"), "Environment", resp.StatusCode)
}

func (l *RemoteProvider) GetConnectionsOfEnvironment(req *http.Request, environmentID, page, pageSize, search, order, filter string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistEnvironments) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Environment", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistEnvironments)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + environmentID + "/connections")

	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Environment", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Connections successfully retrieved from environment")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get environments"), "Environment", resp.StatusCode)
}

func (l *RemoteProvider) GetOrganizations(token, page, pageSize, search, order, filter string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistOrganizations) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Organization", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistOrganizations)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Organization", http.StatusUnauthorized)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Organization")
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("user data successfully retrieved from remote provider")
		return bd, nil
	}

	return nil, ErrFetch(fmt.Errorf("failed to get organizations"), "Organization", resp.StatusCode)
}

func (l *RemoteProvider) GetWorkspaces(token, page, pageSize, search, order, filter, orgID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	if orgID != "" {
		q.Set("orgID", orgID)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Users Data", http.StatusUnauthorized)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Workspaces")
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Workspaces data successfully retrieved from remote provider")
		return bd, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get workspaces"), "Workspaces", resp.StatusCode)
}

func (l *RemoteProvider) GetWorkspaceByID(req *http.Request, workspaceID, orgID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID)
	q := remoteProviderURL.Query()
	if orgID != "" {
		q.Set("orgID", orgID)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		l.Log.Info("Workspace successfully retrieved from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get workspace by ID"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) SaveWorkspace(req *http.Request, env *v1beta1.WorkspacePayload, token string, skipTokenCheck bool) ([]byte, error) {

	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte(""), ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	_env, err := json.Marshal(env)
	if err != nil {
		return []byte(""), err
	}
	bf := bytes.NewBuffer(_env)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	tokenString := token
	if !skipTokenCheck {
		tokenString, err = l.GetToken(req)
		if err != nil {
			l.Log.Error(ErrGetToken(err))
			return []byte(""), err
		}
	}
	resp, err := l.DoRequest(cReq, tokenString)
	if err != nil {
		if resp == nil {
			return []byte(""), ErrUnreachableRemoteProvider(err)
		}
		return []byte(""), ErrFetch(err, "Save Workspace", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bfResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), ErrDataRead(err, "Save Workspace")
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return bfResp, nil
	}

	return []byte(""), ErrPost(fmt.Errorf("failed to save the workspace"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) DeleteWorkspace(req *http.Request, workspaceID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID)
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Workspace successfully deleted from remote provider")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to delete workspace"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) UpdateWorkspace(req *http.Request, env *v1beta1.WorkspacePayload, workspaceID string) (*v1beta1.Workspace, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return &v1beta1.Workspace{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	_env, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}
	bf := bytes.NewBuffer(_env)

	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID)
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrDataRead(err, "Update Workspace")
	}

	if resp.StatusCode == http.StatusOK {
		var workspace v1beta1.Workspace
		if err = json.Unmarshal(bdr, &workspace); err != nil {
			return nil, err
		}
		return &workspace, nil
	}

	return nil, ErrFetch(fmt.Errorf("failed to update the workspace"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) AddEnvironmentToWorkspace(req *http.Request, workspaceID string, environmentID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID + "/environments/" + environmentID)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Environment successfully added to workspace")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get workspaces"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) RemoveEnvironmentFromWorkspace(req *http.Request, workspaceID string, environmentID string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID + "/environments/" + environmentID)
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		if resp == nil {
			return nil, ErrUnreachableRemoteProvider(err)
		}
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Environment successfully removed from workspace")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to unassign environment from workspace"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) GetEnvironmentsOfWorkspace(req *http.Request, workspaceID, page, pageSize, search, order, filter string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID + "/environments")

	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Environments successfully retrieved from workspace")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get environments of workspace"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) AddDesignToWorkspace(req *http.Request, workspaceID string, designId string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID + "/designs/" + designId)
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Design successfully added to workspace")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to add design to workspace"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) RemoveDesignFromWorkspace(req *http.Request, workspaceID string, designId string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID + "/designs/" + designId)
	cReq, _ := http.NewRequest(http.MethodDelete, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Design successfully removed from workspace")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to remove design from workspace"), "Workspace", resp.StatusCode)
}

func (l *RemoteProvider) GetDesignsOfWorkspace(req *http.Request, workspaceID, page, pageSize, search, order, filter string) ([]byte, error) {
	if !l.Capabilities.IsSupported(PersistWorkspaces) {
		l.Log.Warn(ErrOperationNotAvaibale)

		return []byte{}, ErrInvalidCapability("Workspace", l.ProviderName)
	}

	ep, _ := l.Capabilities.GetEndpointForFeature(PersistWorkspaces)
	remoteProviderURL, _ := url.Parse(l.RemoteProviderURL + ep + "/" + workspaceID + "/designs")

	q := remoteProviderURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("pagesize", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	remoteProviderURL.RawQuery = q.Encode()

	cReq, _ := http.NewRequest(http.MethodGet, remoteProviderURL.String(), nil)
	token, err := l.GetToken(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.DoRequest(cReq, token)
	if err != nil {
		return nil, ErrFetch(err, "Workspace", resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Error(ErrDataRead(err, "respone body"))
		return nil, err
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		l.Log.Info("Designs successfully retrieved from workspace")
		return bdr, nil
	}
	return nil, ErrFetch(fmt.Errorf("failed to get designs of workspace"), "Workspace", resp.StatusCode)
}
