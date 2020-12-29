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
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// DefaultLocalProvider - represents a local provider
type DefaultLocalProvider struct {
	*MapPreferencePersister
	ProviderProperties
	SaaSBaseURL           string
	ResultPersister       *BitCaskResultsPersister
	SmiResultPersister    *BitCaskSmiResultsPersister
	TestProfilesPersister *BitCaskTestProfilesPersister
}

// Initialize will initialize the local provider
func (l *DefaultLocalProvider) Initialize() {
	l.ProviderName = "None"
	l.ProviderDescription = []string{
		"Ephemeral sessions",
		"Environment setup not saved",
		"No performance test result history",
		"Free Use",
	}
	l.ProviderType = LocalProviderType
	l.PackageVersion = "v0.0.1"
	l.PackageURL = ""
	l.Extensions = Extensions{}
	l.Capabilities = Capabilities{}
}

// Name - Returns Provider's friendly name
func (l *DefaultLocalProvider) Name() string {
	return l.ProviderName
}

// Description - returns a short description of the provider for display in the Provider UI
func (l *DefaultLocalProvider) Description() []string {
	return l.ProviderDescription
}

// GetProviderType - Returns ProviderType
func (l *DefaultLocalProvider) GetProviderType() ProviderType {
	return l.ProviderType
}

// GetProviderProperties - Returns all the provider properties required
func (l *DefaultLocalProvider) GetProviderProperties() ProviderProperties {
	return l.ProviderProperties
}

// PackageLocation returns an empty string as there is no extension package for
// the local provider
func (l *DefaultLocalProvider) PackageLocation() string {
	return ""
}

// GetProviderCapabilities returns all of the provider properties
func (l *DefaultLocalProvider) GetProviderCapabilities(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(l.ProviderProperties)
}

// InitiateLogin - initiates login flow and returns a true to indicate the handler to "return" or false to continue
func (l *DefaultLocalProvider) InitiateLogin(w http.ResponseWriter, r *http.Request, fromMiddleWare bool) {
	// l.issueSession(w, r, fromMiddleWare)
}

// issueSession issues a cookie session after successful login
func (l *DefaultLocalProvider) issueSession(w http.ResponseWriter, req *http.Request, fromMiddleWare bool) {
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
	return l.fetchUserDetails(), nil
}

// GetSession - returns the session
func (l *DefaultLocalProvider) GetSession(req *http.Request) error {
	return nil
}

// GetProviderToken - returns provider token
func (l *DefaultLocalProvider) GetProviderToken(req *http.Request) (string, error) {
	return "", nil
}

// Logout - logout from provider backend
func (l *DefaultLocalProvider) Logout(w http.ResponseWriter, req *http.Request) {
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

// PublishResults - publishes results to the provider backend synchronously
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

// FetchSmiResults - fetches results from provider backend
func (l *DefaultLocalProvider) FetchSmiResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
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

// PublishSmiResults - publishes results to the provider backend synchronously
func (l *DefaultLocalProvider) PublishSmiResults(result *SmiResult) (string, error) {
	key, _ := uuid.NewV4()
	result.ID = key
	data, err := json.Marshal(result)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal meshery result for persisting"))
		return "", err
	}

	if err := l.SmiResultPersister.WriteResult(key, data); err != nil {
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

// UpdateToken - specific to remote auth
func (l *DefaultLocalProvider) UpdateToken(http.ResponseWriter, *http.Request) {}

// TokenHandler - specific to remote auth
func (l *DefaultLocalProvider) TokenHandler(w http.ResponseWriter, r *http.Request, fromMiddleWare bool) {
}

// ExtractToken - Returns the auth token and the provider type
func (l *DefaultLocalProvider) ExtractToken(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"meshery-provider": l.Name(),
		tokenName:          "",
	}
	logrus.Debugf("encoded response : %v", resp)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.Errorf("Unable to extract auth details: %v", err)
		http.Error(w, "unable to extract auth details", http.StatusInternalServerError)
	}
}

// SMPTestConfigStore Stores the given PerformanceTestConfig into local datastore
func (l *DefaultLocalProvider) SMPTestConfigStore(req *http.Request, perfConfig *SMP.PerformanceTestConfig) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to generate new UUID"))
		return "", err
	}
	perfConfig.Id = uid.String()
	data, err := json.Marshal(perfConfig)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal test config for persisting"))
		return "", err
	}
	return uid.String(), l.TestProfilesPersister.WriteTestConfig(uid, data)
}

// SMPTestConfigGet gets the given PerformanceTestConfig from the local datastore
func (l *DefaultLocalProvider) SMPTestConfigGet(req *http.Request, testUUID string) (*SMP.PerformanceTestConfig, error) {
	uid, err := uuid.FromString(testUUID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to generate new UUID"))
		return nil, err
	}
	return l.TestProfilesPersister.GetTestConfig(uid)
}

// SMPTestConfigFetch gets all the PerformanceTestConfigs from the local datastore
func (l *DefaultLocalProvider) SMPTestConfigFetch(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
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
	return l.TestProfilesPersister.GetTestConfigs(pg, pgs)
}

// SMPTestConfigDelete deletes the given PerformanceTestConfig from the local datastore
func (l *DefaultLocalProvider) SMPTestConfigDelete(req *http.Request, testUUID string) error {
	uid, err := uuid.FromString(testUUID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to generate new UUID"))
		return err
	}
	return l.TestProfilesPersister.DeleteTestConfig(uid)
}
