package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/models/walker"
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/model"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DefaultLocalProvider - represents a local provider
type DefaultLocalProvider struct {
	*MapPreferencePersister
	ProviderProperties
	ProviderBaseURL              string
	ResultPersister              *MesheryResultsPersister
	SmiResultPersister           *BitCaskSmiResultsPersister
	TestProfilesPersister        *BitCaskTestProfilesPersister
	PerformanceProfilesPersister *PerformanceProfilePersister
	MesheryPatternPersister      *MesheryPatternPersister
	MesheryFilterPersister       *MesheryFilterPersister
	GenericPersister             database.Handler
	GraphqlHandler               http.Handler
	GraphqlPlayground            http.Handler
	KubeClient                   *mesherykube.Client
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
	l.PackageVersion = viper.GetString("BUILD")
	l.PackageURL = ""
	l.Extensions = Extensions{}
	l.Capabilities = Capabilities{
		// {Feature: PersistMesheryPatterns},
	}
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
	if err := encoder.Encode(l.ProviderProperties); err != nil {
		http.Error(w, "failed to encode provider capabilities", http.StatusInternalServerError)
	}
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
func (l *DefaultLocalProvider) FetchResults(req *http.Request, page, pageSize, search, order, profileID string) ([]byte, error) {
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
	return l.ResultPersister.GetResults(pg, pgs, profileID)
}

// FetchResults - fetches results from provider backend
func (l *DefaultLocalProvider) FetchAllResults(req *http.Request, page, pageSize, search, order, from, to string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}
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

	return l.ResultPersister.GetAllResults(pg, pgs)
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
func (l *DefaultLocalProvider) PublishResults(req *http.Request, result *MesheryResult, profileID string) (string, error) {
	profileUUID, err := uuid.FromString(profileID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - invalid performance profile id"))
		return "", err
	}

	result.PerformanceProfile = &profileUUID
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
	remoteProviderURL, _ := url.Parse(l.ProviderBaseURL + "/result")
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
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
		// logrus.Infof("results successfully published to reomote provider")
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

	remoteProviderURL, _ := url.Parse(l.ProviderBaseURL + "/result/metrics")
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)
	cReq.Header.Set("X-API-Key", GlobalTokenForAnonymousResults)
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		logrus.Warnf("unable to send metrics: %v", err)
		return nil
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

// SaveMesheryPattern saves given pattern with the provider
func (l *DefaultLocalProvider) SaveMesheryPattern(tokenString string, pattern *MesheryPattern) ([]byte, error) {
	return l.MesheryPatternPersister.SaveMesheryPattern(pattern)
}

// GetMesheryPatterns gives the patterns stored with the provider
func (l *DefaultLocalProvider) GetMesheryPatterns(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

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

	return l.MesheryPatternPersister.GetMesheryPatterns(search, order, pg, pgs)
}

// GetMesheryPattern gets pattern for the given patternID
func (l *DefaultLocalProvider) GetMesheryPattern(req *http.Request, patternID string) ([]byte, error) {
	id := uuid.FromStringOrNil(patternID)
	return l.MesheryPatternPersister.GetMesheryPattern(id)
}

// DeleteMesheryPattern deletes a meshery pattern with the given id
func (l *DefaultLocalProvider) DeleteMesheryPattern(req *http.Request, patternID string) ([]byte, error) {
	id := uuid.FromStringOrNil(patternID)
	return l.MesheryPatternPersister.DeleteMesheryPattern(id)
}

// RemotePatternFile takes in the
func (l *DefaultLocalProvider) RemotePatternFile(req *http.Request, resourceURL, path string, save bool) ([]byte, error) {
	parsedURL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, err
	}

	// Check if hostname is github
	if parsedURL.Host == "github.com" {
		parsedPath := strings.Split(parsedURL.Path, "/")
		if len(parsedPath) < 3 {
			return nil, fmt.Errorf("malformed URL: url should be of type github.com/<owner>/<repo>/[branch]")
		}

		owner := parsedPath[1]
		repo := parsedPath[2]
		branch := "main"

		if len(parsedPath) == 4 {
			branch = parsedPath[3]
		}

		pfs, err := githubRepoPatternScan(owner, repo, path, branch)
		if err != nil {
			return nil, err
		}

		if save {
			return l.MesheryPatternPersister.SaveMesheryPatterns(pfs)
		}

		return json.Marshal(pfs)
	}

	// Fallback to generic HTTP import
	pfs, err := genericHTTPPatternFile(resourceURL)
	if err != nil {
		return nil, err
	}
	if save {
		return l.MesheryPatternPersister.SaveMesheryPatterns(pfs)
	}

	return json.Marshal(pfs)
}

// SaveMesheryFilter saves given filter with the provider
func (l *DefaultLocalProvider) SaveMesheryFilter(tokenString string, filter *MesheryFilter) ([]byte, error) {
	return l.MesheryFilterPersister.SaveMesheryFilter(filter)
}

// GetMesheryFilters gives the filters stored with the provider
func (l *DefaultLocalProvider) GetMesheryFilters(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

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

	return l.MesheryFilterPersister.GetMesheryFilters(search, order, pg, pgs)
}

// GetMesheryFilter gets filter for the given filterID
func (l *DefaultLocalProvider) GetMesheryFilter(req *http.Request, filterID string) ([]byte, error) {
	id := uuid.FromStringOrNil(filterID)
	return l.MesheryFilterPersister.GetMesheryFilter(id)
}

// DeleteMesheryFilter deletes a meshery filter with the given id
func (l *DefaultLocalProvider) DeleteMesheryFilter(req *http.Request, filterID string) ([]byte, error) {
	id := uuid.FromStringOrNil(filterID)
	return l.MesheryFilterPersister.DeleteMesheryFilter(id)
}

// RemoteFilterFile takes in the
func (l *DefaultLocalProvider) RemoteFilterFile(req *http.Request, resourceURL, path string, save bool) ([]byte, error) {
	parsedURL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, err
	}

	// Check if hostname is github
	if parsedURL.Host == "github.com" {
		parsedPath := strings.Split(parsedURL.Path, "/")
		if len(parsedPath) < 3 {
			return nil, fmt.Errorf("malformed URL: url should be of type github.com/<owner>/<repo>/[branch]")
		}

		owner := parsedPath[1]
		repo := parsedPath[2]
		branch := "main"

		if len(parsedPath) == 4 {
			branch = parsedPath[3]
		}

		ffs, err := githubRepoFilterScan(owner, repo, path, branch)
		if err != nil {
			return nil, err
		}

		if save {
			return l.MesheryFilterPersister.SaveMesheryFilters(ffs)
		}

		return json.Marshal(ffs)
	}

	// Fallback to generic HTTP import
	ffs, err := genericHTTPFilterFile(resourceURL)
	if err != nil {
		return nil, err
	}
	if save {
		return l.MesheryFilterPersister.SaveMesheryFilters(ffs)
	}

	return json.Marshal(ffs)
}

// SavePerformanceProfile saves given performance profile with the provider
func (l *DefaultLocalProvider) SavePerformanceProfile(tokenString string, performanceProfile *PerformanceProfile) ([]byte, error) {
	var uid uuid.UUID
	if performanceProfile.ID != nil {
		uid = *performanceProfile.ID
	} else {
		var err error
		uid, err = uuid.NewV4()
		if err != nil {
			logrus.Error(errors.Wrap(err, "error - unable to generate new UUID"))
			return nil, err
		}
		performanceProfile.ID = &uid
	}

	data, err := json.Marshal(performanceProfile)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal performance profile for persisting"))
		return nil, err
	}

	return data, l.PerformanceProfilesPersister.SavePerformanceProfile(uid, performanceProfile)
}

// GetPerformanceProfiles gives the performance profiles stored with the provider
func (l *DefaultLocalProvider) GetPerformanceProfiles(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

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

	return l.PerformanceProfilesPersister.GetPerformanceProfiles("", "", "", pg, pgs)
}

// GetPerformanceProfile gets performance profile for the given performance profileID
func (l *DefaultLocalProvider) GetPerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error) {
	uid, err := uuid.FromString(performanceProfileID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to generate new UUID"))
		return nil, err
	}

	profile, err := l.PerformanceProfilesPersister.GetPerformanceProfile(uid)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to get performance profile"))
		return nil, err
	}

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to marshal performance profile"))
		return nil, err
	}

	return profileJSON, nil
}

// DeletePerformanceProfile deletes a meshery performance profile with the given id
func (l *DefaultLocalProvider) DeletePerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error) {
	uid, err := uuid.FromString(performanceProfileID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error - unable to generate new UUID"))
		return nil, err
	}

	return l.PerformanceProfilesPersister.DeletePerformanceProfile(uid)
}

// SaveSchedule saves a schedule
func (l *DefaultLocalProvider) SaveSchedule(tokenString string, schedule *Schedule) ([]byte, error) {
	return []byte{}, fmt.Errorf("function not supported by local provider")
}

// GetSchedules gets the schedules stored by the current user
func (l *DefaultLocalProvider) GetSchedules(req *http.Request, page, pageSize, order string) ([]byte, error) {
	return []byte{}, fmt.Errorf("function not supported by local provider")
}

// GetSchedule gets a schedule with the given id
func (l *DefaultLocalProvider) GetSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	return []byte{}, fmt.Errorf("function not supported by local provider")
}

// DeleteSchedule deletes a schedule with the given id
func (l *DefaultLocalProvider) DeleteSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	return []byte{}, fmt.Errorf("function not supported by local provider")
}

// RecordMeshSyncData records the mesh sync data
func (l *DefaultLocalProvider) RecordMeshSyncData(obj model.Object) error {
	result := l.GenericPersister.Create(&obj)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ReadMeshSyncData reads the mesh sync data
func (l *DefaultLocalProvider) ReadMeshSyncData() ([]model.Object, error) {
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

// GetGenericPersister - to return persister
func (l *DefaultLocalProvider) GetGenericPersister() *database.Handler {
	return &l.GenericPersister
}

// GetGraphqlHandler - to return graphql handler instance
func (l *DefaultLocalProvider) GetGraphqlHandler() http.Handler {
	return l.GraphqlHandler
}

// GetGraphqlPlayground - to return graphql playground instance
func (l *DefaultLocalProvider) GetGraphqlPlayground() http.Handler {
	return l.GraphqlPlayground
}

// SetKubeClient - to set meshery kubernetes client
func (l *DefaultLocalProvider) SetKubeClient(client *mesherykube.Client) {
	l.KubeClient = client
}

// GetKubeClient - to get meshery kubernetes client
func (l *DefaultLocalProvider) GetKubeClient() *mesherykube.Client {
	return l.KubeClient
}

// githubRepoPatternScan & githubRepoFilterScan takes in github repo owner, repo name, path from where the file/files are needed
// to be imported
//
// If the path name is like dir1/dir2 then only the given path will be scanned for the file
// but if the path name is like dir1/dir2/** then it will recursively scan all of the directories
// inside the path. If path is empty then only the root is scanned for patterns/ filters
func githubRepoPatternScan(
	owner,
	repo,
	path,
	branch string,
) ([]MesheryPattern, error) {
	var mu sync.Mutex
	ghWalker := walker.NewGithub()
	result := make([]MesheryPattern, 0)

	err := ghWalker.
		Owner(owner).
		Repo(repo).
		Branch(branch).
		Root(path).
		RegisterFileInterceptor(func(data walker.GithubContentAPI) error {
			ext := filepath.Ext(data.Name)
			if ext == ".yml" || ext == ".yaml" {
				decodedContent, err := base64.StdEncoding.DecodeString(data.Content)
				if err != nil {
					return err
				}

				name, err := GetPatternName(string(decodedContent))
				if err != nil {
					return err
				}

				pf := MesheryPattern{
					Name:        name,
					PatternFile: string(decodedContent),
					Location: map[string]interface{}{
						"type":   "github",
						"host":   fmt.Sprintf("github.com/%s/%s", owner, repo),
						"path":   data.Path,
						"branch": branch,
					},
				}

				mu.Lock()
				result = append(result, pf)
				mu.Unlock()
			}

			return nil
		}).
		Walk()

	return result, err
}

func githubRepoFilterScan(
	owner,
	repo,
	path,
	branch string,
) ([]MesheryFilter, error) {
	var mu sync.Mutex
	ghWalker := walker.NewGithub()
	result := make([]MesheryFilter, 0)

	err := ghWalker.
		Owner(owner).
		Repo(repo).
		Branch(branch).
		Root(path).
		RegisterFileInterceptor(func(data walker.GithubContentAPI) error {
			ext := filepath.Ext(data.Name)
			if ext == ".yml" || ext == ".yaml" {
				decodedContent, err := base64.StdEncoding.DecodeString(data.Content)
				if err != nil {
					return err
				}

				name, err := GetFilterName(string(decodedContent))
				if err != nil {
					return err
				}

				ff := MesheryFilter{
					Name:       name,
					FilterFile: string(decodedContent),
					Location: map[string]interface{}{
						"type":   "github",
						"host":   fmt.Sprintf("github.com/%s/%s", owner, repo),
						"path":   data.Path,
						"branch": branch,
					},
				}

				mu.Lock()
				result = append(result, ff)
				mu.Unlock()
			}

			return nil
		}).
		Walk()

	return result, err
}

func genericHTTPPatternFile(fileURL string) ([]MesheryPattern, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file not found")
	}

	defer SafeClose(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := string(body)

	name, err := GetPatternName(result)
	if err != nil {
		return nil, err
	}

	pf := MesheryPattern{
		Name:        name,
		PatternFile: result,
		Location: map[string]interface{}{
			"type":   "http",
			"host":   fileURL,
			"path":   "",
			"branch": "",
		},
	}

	return []MesheryPattern{pf}, nil
}

func genericHTTPFilterFile(fileURL string) ([]MesheryFilter, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file not found")
	}

	defer SafeClose(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := string(body)

	name, err := GetFilterName(result)
	if err != nil {
		return nil, err
	}

	ff := MesheryFilter{
		Name:       name,
		FilterFile: result,
		Location: map[string]interface{}{
			"type":   "http",
			"host":   fileURL,
			"path":   "",
			"branch": "",
		},
	}

	return []MesheryFilter{ff}, nil
}
