package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/walker"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/meshery/schemas/models/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

// DefaultLocalProvider - represents a local provider
type DefaultLocalProvider struct {
	*MapPreferencePersister
	*UserCapabilitiesPersister
	*EventsPersister
	ProviderProperties
	ProviderBaseURL                 string
	ResultPersister                 *MesheryResultsPersister
	SmiResultPersister              *SMIResultsPersister
	TestProfilesPersister           *TestProfilesPersister
	PerformanceProfilesPersister    *PerformanceProfilePersister
	MesheryPatternPersister         *MesheryPatternPersister
	MesheryPatternResourcePersister *PatternResourcePersister
	MesheryApplicationPersister     *MesheryApplicationPersister
	MesheryFilterPersister          *MesheryFilterPersister
	MesheryK8sContextPersister      *MesheryK8sContextPersister
	OrganizationPersister           *OrganizationPersister
	KeyPersister                    *KeyPersister
	ConnectionPersister             *ConnectionPersister
	EnvironmentPersister            *EnvironmentPersister
	WorkspacePersister              *WorkspacePersister

	GenericPersister *database.Handler
	KubeClient       *mesherykube.Client
	Log              logger.Handler
}

// Initialize will initialize the local provider
func (l *DefaultLocalProvider) Initialize() {
	l.ProviderName = "None"
	l.ProviderDescription = []string{
		"Ephemeral sessions",
		"Environment setup not saved",
		"No performance or conformance test result history",
		"Free Use",
	}
	l.ProviderType = LocalProviderType
	l.PackageVersion = viper.GetString("BUILD")
	l.PackageURL = ""
	l.Extensions = Extensions{}
	l.Capabilities = Capabilities{
		{Feature: PersistMesheryPatterns},
		{Feature: PersistMesheryApplications},
		{Feature: PersistMesheryFilters},
		{Feature: PersistCredentials},
	}
}

// Name - Returns Provider's friendly name
func (l *DefaultLocalProvider) Name() string {
	return l.ProviderName
}

func (l *DefaultLocalProvider) GetProviderURL() string {
	return l.ProviderBaseURL
}

// Description - returns a short description of the provider for display in the Provider UI
func (l *DefaultLocalProvider) Description() []string {
	return l.ProviderDescription
}

// GetProviderType - Returns ProviderType
func (l *DefaultLocalProvider) GetProviderType() ProviderType {
	return l.ProviderType
}

func (l *DefaultLocalProvider) DownloadProviderExtensionPackage() {
}

func (l *DefaultLocalProvider) SetProviderProperties(providerProperties ProviderProperties) {
	l.ProviderProperties = providerProperties
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

func (l *DefaultLocalProvider) SetJWTCookie(_ http.ResponseWriter, _ string) {
}

func (l *DefaultLocalProvider) UnSetJWTCookie(_ http.ResponseWriter) {
}

// GetProviderCapabilities returns all of the provider properties
func (l *DefaultLocalProvider) GetProviderCapabilities(w http.ResponseWriter, _ *http.Request, _ string) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(l.ProviderProperties); err != nil {
		http.Error(w, "failed to encode provider capabilities", http.StatusInternalServerError)
	}
}

// InitiateLogin - initiates login flow and returns a true to indicate the handler to "return" or false to continue
func (l *DefaultLocalProvider) InitiateLogin(_ http.ResponseWriter, _ *http.Request, _ bool) {
	// l.issueSession(w, r, fromMiddleWare)
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
func (l *DefaultLocalProvider) GetUserDetails(_ *http.Request) (*User, error) {
	return l.fetchUserDetails(), nil
}

func (l *DefaultLocalProvider) GetUserByID(_ *http.Request, _ string) ([]byte, error) {
	return nil, nil
}

func (l *DefaultLocalProvider) GetUsers(_, _, _, _, _, _ string) ([]byte, error) {
	return []byte(""), ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) GetEnvironments(_, page, pageSize, search, order, filter string, orgID string) ([]byte, error) {
	return l.EnvironmentPersister.GetEnvironments(orgID, search, order, page, pageSize, filter)
}

func (l *DefaultLocalProvider) GetEnvironmentByID(_ *http.Request, environmentID string, _ string) ([]byte, error) {
	id := uuid.FromStringOrNil(environmentID)
	return l.EnvironmentPersister.GetEnvironmentByID(id)
}

func (l *DefaultLocalProvider) DeleteEnvironment(_ *http.Request, environmentID string) ([]byte, error) {
	id := uuid.FromStringOrNil(environmentID)
	return l.EnvironmentPersister.DeleteEnvironmentByID(id)
}

func (l *DefaultLocalProvider) SaveEnvironment(_ *http.Request, environmentPayload *v1beta1.EnvironmentPayload, _ string, _ bool) ([]byte, error) {
	orgId, _ := uuid.FromString(environmentPayload.OrgId)
	environment := &v1beta1.Environment{
		CreatedAt:      time.Now(),
		Description:    environmentPayload.Description,
		Name:           environmentPayload.Name,
		OrganizationId: orgId,
		Owner:          "Meshery",
		UpdatedAt:      time.Now(),
	}
	return l.EnvironmentPersister.SaveEnvironment(environment)
}

func (l *DefaultLocalProvider) UpdateEnvironment(_ *http.Request, environmentPayload *v1beta1.EnvironmentPayload, environmentID string) (*v1beta1.Environment, error) {
	id, _ := uuid.FromString(environmentID)
	orgId, _ := uuid.FromString(environmentPayload.OrgId)
	environment := &v1beta1.Environment{
		ID:             id,
		CreatedAt:      time.Now(),
		Description:    environmentPayload.Description,
		Name:           environmentPayload.Name,
		OrganizationId: orgId,
		Owner:          "Meshery",
		UpdatedAt:      time.Now(),
	}
	return l.EnvironmentPersister.UpdateEnvironmentByID(environment)
}

func (l *DefaultLocalProvider) AddConnectionToEnvironment(_ *http.Request, environmentID string, connectionID string) ([]byte, error) {
	envId, _ := uuid.FromString(environmentID)
	conId, _ := uuid.FromString(connectionID)
	return l.EnvironmentPersister.AddConnectionToEnvironment(envId, conId)
}

func (l *DefaultLocalProvider) RemoveConnectionFromEnvironment(_ *http.Request, environmentID string, connectionID string) ([]byte, error) {
	envId, _ := uuid.FromString(environmentID)
	conId, _ := uuid.FromString(connectionID)
	return l.EnvironmentPersister.DeleteConnectionFromEnvironment(envId, conId)
}

func (l *DefaultLocalProvider) GetConnectionsOfEnvironment(_ *http.Request, environmentID, page, pageSize, search, order, filter string) ([]byte, error) {
	envId, _ := uuid.FromString(environmentID)
	return l.EnvironmentPersister.GetEnvironmentConnections(envId, search, order, page, pageSize, filter)
}

// GetSession - returns the session
func (l *DefaultLocalProvider) GetSession(_ *http.Request) error {
	return nil
}

// GetProviderToken - returns provider token
func (l *DefaultLocalProvider) GetProviderToken(_ *http.Request) (string, error) {
	return "", nil
}

// Logout - logout from provider backend
func (l *DefaultLocalProvider) Logout(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// HandleUnAuthenticated - logout from provider backend
func (l *DefaultLocalProvider) HandleUnAuthenticated(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/user/login", http.StatusFound)
}

func (l *DefaultLocalProvider) SaveK8sContext(_ string, k8sContext K8sContext) (connections.Connection, error) {

	k8sServerID := *k8sContext.KubernetesServerID

	var connID uuid.UUID
	id, err := K8sContextGenerateID(k8sContext)
	if err == nil {
		connID, _ = uuid.FromString(id)
	}

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
		ID:      connID,
		Kind:    "kubernetes",
		Name:    k8sContext.Name,
		Type:    "platform",
		SubType: "orchestrator",
		// Eventually the status would depend on other factors like, whether user administratively processed it or not
		// Is clsuter reachable and other reasons.
		Status:           connections.DISCOVERED,
		MetaData:         metadata,
		CredentialSecret: cred,
	}
	connectionCreated, err := l.SaveConnection(conn, "", true)
	if err != nil {
		return connections.Connection{}, fmt.Errorf("error in saving k8s context %v", err)
	}
	_, _ = l.MesheryK8sContextPersister.SaveMesheryK8sContext(k8sContext)
	return *connectionCreated, nil
}

func (l *DefaultLocalProvider) GetK8sContexts(_, page, pageSize, search, order string, withStatus string, withCredentials bool) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}

	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}

	return l.MesheryK8sContextPersister.GetMesheryK8sContexts(search, order, pg, pgs)
}

func (l *DefaultLocalProvider) DeleteK8sContext(_, id string) (K8sContext, error) {
	return l.MesheryK8sContextPersister.DeleteMesheryK8sContext(id)
}

func (l *DefaultLocalProvider) GetK8sContext(_, id string) (K8sContext, error) {
	idStrWithoutDashes := strings.ReplaceAll(id, "-", "")
	return l.MesheryK8sContextPersister.GetMesheryK8sContext(idStrWithoutDashes)
}

func (l *DefaultLocalProvider) LoadAllK8sContext(token string) ([]*K8sContext, error) {
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
			obj := "k8s context"
			return nil, ErrUnmarshal(err, obj)
		}
		results = append(results, k8scontext.Contexts...)

		if page*pageSize >= k8scontext.TotalCount {
			break
		}

		page++
	}

	return results, nil
}

// func (l *DefaultLocalProvider) SetCurrentContext(token, id string) (K8sContext, error) {
// 	if err := l.MesheryK8sContextPersister.SetMesheryK8sCurrentContext(id); err != nil {
// 		return K8sContext{}, err
// 	}

// 	return l.MesheryK8sContextPersister.GetMesheryK8sContext(id)
// }

// func (l *DefaultLocalProvider) GetCurrentContext(token string) (K8sContext, error) {
// 	return l.MesheryK8sContextPersister.GetMesheryK8sCurrentContext()
// }

// FetchResults - fetches results from provider backend
func (l *DefaultLocalProvider) FetchResults(_, page, pageSize, _, _, profileID string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}
	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}
	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}
	return l.ResultPersister.GetResults(pg, pgs, profileID, l.Log)
}

// FetchResults - fetches results from provider backend
func (l *DefaultLocalProvider) FetchAllResults(_, page, pageSize, _, _, _, _ string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}
	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}
	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}

	return l.ResultPersister.GetAllResults(pg, pgs, l.Log)
}

// GetResult - fetches result from provider backend for the given result id
func (l *DefaultLocalProvider) GetResult(_ string, resultID uuid.UUID) (*MesheryResult, error) {
	// key := uuid.FromStringOrNil(resultID)
	if resultID == uuid.Nil {
		return nil, ErrResultID
	}
	return l.ResultPersister.GetResult(resultID, l.Log)
}

// PublishResults - publishes results to the provider backend synchronously
func (l *DefaultLocalProvider) PublishResults(req *http.Request, result *MesheryResult, profileID string) (string, error) {
	profileUUID, err := uuid.FromString(profileID)
	if err != nil {
		return "", ErrPerfID(err)
	}

	result.PerformanceProfile = &profileUUID
	data, err := json.Marshal(result)
	if err != nil {
		return "", ErrMarshal(err, "meshery result for shipping")
	}
	user, _ := l.GetUserDetails(req)
	pref, _ := l.ReadFromPersister(user.UserID)
	if !pref.AnonymousPerfResults {
		return "", nil
	}

	l.Log.Debug(fmt.Sprintf("Result: %s, size: %d", data, len(data)))
	resultID, _ := l.shipResults(req, data)

	key := uuid.FromStringOrNil(resultID)
	l.Log.Debug(fmt.Sprintf("key: %s, is nil: %t", key.String(), (key == uuid.Nil)))
	if key == uuid.Nil {
		key, _ = uuid.NewV4()
		result.ID = key
		data, err = json.Marshal(result)
		if err != nil {
			return "", ErrMarshal(err, "Meshery Result for Persisting")
		}
	}
	if err := l.ResultPersister.WriteResult(key, data); err != nil {
		return "", err
	}

	return key.String(), nil
}

// FetchSmiResults - fetches results from provider backend
func (l *DefaultLocalProvider) FetchSmiResults(_ *http.Request, page, pageSize, _, _ string) ([]byte, error) {
	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}
	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}
	return l.SmiResultPersister.GetResults(pg, pgs)
}

// FetchSmiResults - fetches results from provider backend
func (l *DefaultLocalProvider) FetchSmiResult(_ *http.Request, _, _, _, _ string, resultID uuid.UUID) ([]byte, error) {
	return l.SmiResultPersister.GetResult(resultID)
}

// PublishSmiResults - publishes results to the provider backend synchronously
func (l *DefaultLocalProvider) PublishSmiResults(result *SmiResult) (string, error) {
	key, _ := uuid.NewV4()
	result.ID = key
	data, err := json.Marshal(result)
	if err != nil {
		return "", ErrMarshal(err, "Meshery Results for Persisting")
	}

	if err := l.SmiResultPersister.WriteResult(key, data); err != nil {
		return "", err
	}

	return key.String(), nil
}

func (l *DefaultLocalProvider) shipResults(_ *http.Request, data []byte) (string, error) {
	bf := bytes.NewBuffer(data)
	remoteProviderURL, _ := url.Parse(l.ProviderBaseURL + "/result")
	cReq, _ := http.NewRequest(http.MethodPost, remoteProviderURL.String(), bf)
	cReq.Header.Set("X-API-Key", GlobalTokenForAnonymousResults)
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		l.Log.Warn(ErrDoRequest(err, cReq.Method, remoteProviderURL.String()))
		return "", nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bdr, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Warn(ErrDataRead(err, "body"))
		return "", nil
	}
	if resp.StatusCode == http.StatusCreated {
		// logrus.Infof("results successfully published to reomote provider")
		idMap := map[string]string{}
		if err = json.Unmarshal(bdr, &idMap); err != nil {
			l.Log.Warn(ErrUnmarshal(err, "body"))
			return "", nil
		}
		resultID, ok := idMap["id"]
		if ok {
			return resultID, nil
		}
		return "", nil
	}
	l.Log.Warn(ErrDoRequest(err, resp.Request.Method, remoteProviderURL.String()))
	return "", nil
}

func (l *DefaultLocalProvider) PublishEventToProvider(_ string, _ events.Event) error {
	return nil
}

// PublishMetrics - publishes metrics to the provider backend asyncronously
func (l *DefaultLocalProvider) PublishMetrics(_ string, result *MesheryResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return ErrMarshal(err, "Meshery Matrics for shipping")
	}

	l.Log.Debug(fmt.Sprintf("Result: %s, size: %d", data, len(data)))
	bf := bytes.NewBuffer(data)

	remoteProviderURL, _ := url.Parse(l.ProviderBaseURL + "/result/metrics")
	cReq, _ := http.NewRequest(http.MethodPut, remoteProviderURL.String(), bf)
	cReq.Header.Set("X-API-Key", GlobalTokenForAnonymousResults)
	c := &http.Client{}
	resp, err := c.Do(cReq)
	if err != nil {
		l.Log.Warn(ErrDoRequest(err, cReq.Method, remoteProviderURL.String()))
		return nil
	}
	if resp.StatusCode == http.StatusOK {
		l.Log.Info("metrics successfully published to remote provider")
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		l.Log.Warn(ErrDataRead(err, "body"))
		return nil
	}
	l.Log.Warn(ErrDoRequest(err, resp.Request.Method, remoteProviderURL.String()))
	return nil
}

// RecordPreferences - records the user preference
func (l *DefaultLocalProvider) RecordPreferences(_ *http.Request, userID string, data *Preference) error {
	return l.MapPreferencePersister.WriteToPersister(userID, data)
}

// UpdateToken - specific to remote auth
func (l *DefaultLocalProvider) UpdateToken(http.ResponseWriter, *http.Request) string {
	return ""
}

// TokenHandler - specific to remote auth
func (l *DefaultLocalProvider) TokenHandler(_ http.ResponseWriter, _ *http.Request, _ bool) {
}

// ExtractToken - Returns the auth token and the provider type
func (l *DefaultLocalProvider) ExtractToken(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]interface{}{
		"meshery-provider": l.Name(),
		TokenCookieName:    "",
	}
	l.Log.Debug(fmt.Sprintf("token sent for meshery-provider %v", l.Name()))
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		l.Log.Error(ErrEncoding(err, "token"))
		http.Error(w, ErrEncoding(err, "token").Error(), http.StatusInternalServerError)
	}
}

// SMPTestConfigStore Stores the given PerformanceTestConfig into local datastore
func (l *DefaultLocalProvider) SMPTestConfigStore(_ *http.Request, perfConfig *SMP.PerformanceTestConfig) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", ErrGenerateUUID(err)
	}
	perfConfig.Id = uid.String()
	data, err := json.Marshal(perfConfig)
	if err != nil {
		return "", ErrMarshal(err, "test config for persisting")
	}
	return uid.String(), l.TestProfilesPersister.WriteTestConfig(uid, data)
}

// SMPTestConfigGet gets the given PerformanceTestConfig from the local datastore
func (l *DefaultLocalProvider) SMPTestConfigGet(_ *http.Request, testUUID string) (*SMP.PerformanceTestConfig, error) {
	uid, err := uuid.FromString(testUUID)
	if err != nil {
		return nil, ErrGenerateUUID(err)
	}
	return l.TestProfilesPersister.GetTestConfig(uid)
}

// SMPTestConfigFetch gets all the PerformanceTestConfigs from the local datastore
func (l *DefaultLocalProvider) SMPTestConfigFetch(_ *http.Request, page, pageSize, _, _ string) ([]byte, error) {
	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}
	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}
	return l.TestProfilesPersister.GetTestConfigs(pg, pgs)
}

// SMPTestConfigDelete deletes the given PerformanceTestConfig from the local datastore
func (l *DefaultLocalProvider) SMPTestConfigDelete(_ *http.Request, testUUID string) error {
	uid, err := uuid.FromString(testUUID)
	if err != nil {
		return ErrGenerateUUID(err)
	}
	return l.TestProfilesPersister.DeleteTestConfig(uid)
}

// SaveMesheryPatternSourceContent nothing needs to be done as pattern is saved with source content for local provider
func (l *DefaultLocalProvider) SaveMesheryPatternSourceContent(_, _ string, _ []byte) error {
	return nil
}

func (l *DefaultLocalProvider) SaveMesheryPatternResource(_ string, resource *PatternResource) (*PatternResource, error) {
	return l.MesheryPatternResourcePersister.SavePatternResource(resource)
}

func (l *DefaultLocalProvider) GetMesheryPatternResource(_, resourceID string) (*PatternResource, error) {
	id := uuid.FromStringOrNil(resourceID)
	return l.MesheryPatternResourcePersister.GetPatternResource(id)
}

func (l *DefaultLocalProvider) GetMesheryPatternResources(
	_,
	page,
	pageSize,
	search,
	order,
	name,
	namespace,
	typ,
	oamType string,
) (*PatternResourcePage, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}

	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}

	return l.MesheryPatternResourcePersister.GetPatternResources(
		search, order, name, namespace, typ, oamType, pg, pgs,
	)
}

func (l *DefaultLocalProvider) DeleteMesheryPatternResource(_, resourceID string) error {
	id := uuid.FromStringOrNil(resourceID)
	return l.MesheryPatternResourcePersister.DeletePatternResource(id)
}

// SaveMesheryPattern saves given pattern with the provider
func (l *DefaultLocalProvider) SaveMesheryPattern(_ string, pattern *MesheryPattern) ([]byte, error) {
	return l.MesheryPatternPersister.SaveMesheryPattern(pattern)
}

// GetMesheryPatterns gives the patterns stored with the provider
func (l *DefaultLocalProvider) GetMesheryPatterns(_, page, pageSize, search, order, updatedAfter string, visibility []string, _ string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}

	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}
	return l.MesheryPatternPersister.GetMesheryPatterns(search, order, pg, pgs, updatedAfter, visibility)
}

// GetCatalogMesheryPatterns gives the catalog patterns stored with the provider
func (l *DefaultLocalProvider) GetCatalogMesheryPatterns(_, page, pageSize, search, order, _ string) ([]byte, error) {
	return l.MesheryPatternPersister.GetMesheryCatalogPatterns(page, pageSize, search, order)
}

// PublishCatalogPattern publishes pattern to catalog
// Not supported by local provider
func (l *DefaultLocalProvider) PublishCatalogPattern(_ *http.Request, _ *MesheryCatalogPatternRequestBody) ([]byte, error) {
	return []byte(""), ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) UnPublishCatalogPattern(_ *http.Request, _ *MesheryCatalogPatternRequestBody) ([]byte, error) {
	return []byte(""), ErrLocalProviderSupport
}

// GetMesheryPattern gets pattern for the given patternID
func (l *DefaultLocalProvider) GetMesheryPattern(_ *http.Request, patternID, _ string) ([]byte, error) {
	id := uuid.FromStringOrNil(patternID)
	return l.MesheryPatternPersister.GetMesheryPattern(id)
}

// DeleteMesheryPattern deletes a meshery pattern with the given id
func (l *DefaultLocalProvider) DeleteMesheryPattern(_ *http.Request, patternID string) ([]byte, error) {
	id := uuid.FromStringOrNil(patternID)
	return l.MesheryPatternPersister.DeleteMesheryPattern(id)
}

// DeleteMesheryPattern deletes a meshery pattern with the given id
func (l *DefaultLocalProvider) DeleteMesheryPatterns(_ *http.Request, patterns MesheryPatternDeleteRequestBody) ([]byte, error) {
	return l.MesheryPatternPersister.DeleteMesheryPatterns(patterns)
}

// CloneMesheryPattern clones a meshery pattern with the given id
func (l *DefaultLocalProvider) CloneMesheryPattern(_ *http.Request, patternID string, clonePatternRequest *MesheryClonePatternRequestBody) ([]byte, error) {
	return l.MesheryPatternPersister.CloneMesheryPattern(patternID, clonePatternRequest)
}

// GetDesignSourceContent returns design source-content from provider
func (l *DefaultLocalProvider) GetDesignSourceContent(_, designID string) ([]byte, error) {
	id := uuid.FromStringOrNil(designID)
	return l.MesheryPatternPersister.GetMesheryPatternSource(id)
}

// RemotePatternFile takes in the
func (l *DefaultLocalProvider) RemotePatternFile(_ *http.Request, resourceURL, path string, save bool) ([]byte, error) {
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
		if len(parsedPath) >= 4 && parsedPath[3] == "tree" {
			parsedPath = append(parsedPath[0:3], parsedPath[4:]...)
		}

		owner := parsedPath[1]
		repo := parsedPath[2]
		branch := "master"

		if len(parsedPath) == 4 {
			branch = parsedPath[3]
		}
		if path == "" && len(parsedPath) > 4 {
			path = strings.Join(parsedPath[4:], "/")
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
	pfs, err := genericHTTPPatternFile(resourceURL, l.Log)
	if err != nil {
		return nil, err
	}
	if save {
		return l.MesheryPatternPersister.SaveMesheryPatterns(pfs)
	}

	return json.Marshal(pfs)
}

// SaveMesheryFilter saves given filter with the provider
func (l *DefaultLocalProvider) SaveMesheryFilter(_ string, filter *MesheryFilter) ([]byte, error) {
	return l.MesheryFilterPersister.SaveMesheryFilter(filter)
}

// GetMesheryFilters gives the filters stored with the provider
func (l *DefaultLocalProvider) GetMesheryFilters(_, page, pageSize, search, order string, visibility []string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}

	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}

	return l.MesheryFilterPersister.GetMesheryFilters(search, order, pg, pgs, visibility)
}

// GetCatalogMesheryFilters gives the catalog filters stored with the provider
func (l *DefaultLocalProvider) GetCatalogMesheryFilters(_ string, page, pageSize, search, order string) ([]byte, error) {
	return l.MesheryFilterPersister.GetMesheryCatalogFilters(page, pageSize, search, order)
}

// PublishCatalogFilter publishes filter to catalog
// Not supported by local provider
func (l *DefaultLocalProvider) PublishCatalogFilter(_ *http.Request, _ *MesheryCatalogFilterRequestBody) ([]byte, error) {
	return []byte(""), nil
}

func (l *DefaultLocalProvider) UnPublishCatalogFilter(_ *http.Request, _ *MesheryCatalogFilterRequestBody) ([]byte, error) {
	return []byte(""), ErrLocalProviderSupport
}

// GetMesheryFilterFile gets filter for the given filterID without the metadata
func (l *DefaultLocalProvider) GetMesheryFilterFile(_ *http.Request, filterID string) ([]byte, error) {
	id := uuid.FromStringOrNil(filterID)
	return l.MesheryFilterPersister.GetMesheryFilterFile(id)
}

// GetMesheryFilter gets filter for the given filterID
func (l *DefaultLocalProvider) GetMesheryFilter(_ *http.Request, filterID string) ([]byte, error) {
	id := uuid.FromStringOrNil(filterID)
	return l.MesheryFilterPersister.GetMesheryFilter(id)
}

// DeleteMesheryFilter deletes a meshery filter with the given id
func (l *DefaultLocalProvider) DeleteMesheryFilter(_ *http.Request, filterID string) ([]byte, error) {
	id := uuid.FromStringOrNil(filterID)
	return l.MesheryFilterPersister.DeleteMesheryFilter(id)
}

// CloneMesheryFilter clones a meshery filter with the given id
func (l *DefaultLocalProvider) CloneMesheryFilter(_ *http.Request, filterID string, cloneFilterRequest *MesheryCloneFilterRequestBody) ([]byte, error) {
	return l.MesheryFilterPersister.CloneMesheryFilter(filterID, cloneFilterRequest)
}

// RemoteFilterFile takes in the
func (l *DefaultLocalProvider) RemoteFilterFile(_ *http.Request, resourceURL, path string, save bool, resource string) ([]byte, error) {
	parsedURL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, err
	}

	// Check if hostname is github
	if parsedURL.Host == "github.com" {
		parsedPath := strings.Split(parsedURL.Path, "/")
		if parsedPath[3] == "tree" {
			parsedPath = append(parsedPath[0:3], parsedPath[4:]...)
		}
		if len(parsedPath) < 3 {
			return nil, fmt.Errorf("malformed URL: url should be of type github.com/<owner>/<repo>/[branch]")
		}

		owner := parsedPath[1]
		repo := parsedPath[2]
		branch := "master"

		if len(parsedPath) == 4 {
			branch = parsedPath[3]
		}
		if path == "" && len(parsedPath) > 4 {
			path = strings.Join(parsedPath[4:], "/")
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
	ffs, err := genericHTTPFilterFile(resourceURL, l.Log)
	if err != nil {
		return nil, err
	}
	if save {
		return l.MesheryFilterPersister.SaveMesheryFilters(ffs)
	}

	return json.Marshal(ffs)
}

// SaveMesheryApplication saves given application with the provider
func (l *DefaultLocalProvider) SaveMesheryApplication(_ string, application *MesheryApplication) ([]byte, error) {
	return l.MesheryApplicationPersister.SaveMesheryApplication(application)
}

// SaveApplicationSourceContent nothing needs to be done as application is saved with source content for local provider
func (l *DefaultLocalProvider) SaveApplicationSourceContent(_, _ string, _ []byte) error {
	return nil
}

// GetApplicationSourceContent returns application source-content from provider
func (l *DefaultLocalProvider) GetApplicationSourceContent(_ *http.Request, applicationID string) ([]byte, error) {
	id := uuid.FromStringOrNil(applicationID)
	return l.MesheryApplicationPersister.GetMesheryApplicationSource(id)
}

// GetMesheryApplications gives the applications stored with the provider
func (l *DefaultLocalProvider) GetMesheryApplications(_, page, pageSize, search, order string, updatedAfter string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}

	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}

	return l.MesheryApplicationPersister.GetMesheryApplications(search, order, pg, pgs, updatedAfter)
}

// GetMesheryApplication gets application for the given applicationID
func (l *DefaultLocalProvider) GetMesheryApplication(_ *http.Request, applicationID string) ([]byte, error) {
	id := uuid.FromStringOrNil(applicationID)
	return l.MesheryApplicationPersister.GetMesheryApplication(id)
}

// DeleteMesheryApplication deletes a meshery application with the given id
func (l *DefaultLocalProvider) DeleteMesheryApplication(_ *http.Request, applicationID string) ([]byte, error) {
	id := uuid.FromStringOrNil(applicationID)
	return l.MesheryApplicationPersister.DeleteMesheryApplication(id)
}

func (l *DefaultLocalProvider) ShareDesign(_ *http.Request) (int, error) {
	return http.StatusForbidden, ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) ShareFilter(_ *http.Request) (int, error) {
	return http.StatusForbidden, ErrLocalProviderSupport
}

// SavePerformanceProfile saves given performance profile with the provider
func (l *DefaultLocalProvider) SavePerformanceProfile(_ string, performanceProfile *PerformanceProfile) ([]byte, error) {
	var uid uuid.UUID
	if performanceProfile.ID != nil {
		uid = *performanceProfile.ID
	} else {
		var err error
		uid, err = uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}
		performanceProfile.ID = &uid
	}

	data, err := json.Marshal(performanceProfile)
	if err != nil {
		return nil, ErrMarshal(err, "Perf Profile for persisting")
	}

	return data, l.PerformanceProfilesPersister.SavePerformanceProfile(uid, performanceProfile)
}

// GetPerformanceProfiles gives the performance profiles stored with the provider
func (l *DefaultLocalProvider) GetPerformanceProfiles(_, page, pageSize, _, _ string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}

	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}

	return l.PerformanceProfilesPersister.GetPerformanceProfiles("", "", "", pg, pgs)
}

// GetPerformanceProfile gets performance profile for the given performance profileID
func (l *DefaultLocalProvider) GetPerformanceProfile(_ *http.Request, performanceProfileID string) ([]byte, error) {
	uid, err := uuid.FromString(performanceProfileID)
	if err != nil {
		return nil, ErrPerfID(err)
	}

	profile, err := l.PerformanceProfilesPersister.GetPerformanceProfile(uid)
	if err != nil {
		return nil, err
	}

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return nil, ErrMarshal(err, "Perf Profile")
	}

	return profileJSON, nil
}

// DeletePerformanceProfile deletes a meshery performance profile with the given id
func (l *DefaultLocalProvider) DeletePerformanceProfile(_ *http.Request, performanceProfileID string) ([]byte, error) {
	uid, err := uuid.FromString(performanceProfileID)
	if err != nil {
		return nil, ErrPerfID(err)
	}

	return l.PerformanceProfilesPersister.DeletePerformanceProfile(uid)
}

// SaveSchedule saves a schedule
func (l *DefaultLocalProvider) SaveSchedule(_ string, _ *Schedule) ([]byte, error) {
	return []byte{}, ErrLocalProviderSupport
}

// GetSchedules gets the schedules stored by the current user
func (l *DefaultLocalProvider) GetSchedules(_ *http.Request, _, _, _ string) ([]byte, error) {
	return []byte{}, ErrLocalProviderSupport
}

// GetSchedule gets a schedule with the given id
func (l *DefaultLocalProvider) GetSchedule(_ *http.Request, _ string) ([]byte, error) {
	return []byte{}, ErrLocalProviderSupport
}

// DeleteSchedule deletes a schedule with the given id
func (l *DefaultLocalProvider) DeleteSchedule(_ *http.Request, _ string) ([]byte, error) {
	return []byte{}, ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) ExtensionProxy(_ *http.Request) (*ExtensionProxyResponse, error) {
	return nil, ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) SaveConnection(conn *connections.ConnectionPayload, _ string, _ bool) (*connections.Connection, error) {
	id := uuid.Nil
	if conn.ID != uuid.Nil {
		id = conn.ID
	}
	connection := &connections.Connection{
		ID:           id,
		Name:         conn.Name,
		CredentialID: uuid.Nil,
		Type:         conn.Type,
		SubType:      conn.SubType,
		Kind:         conn.Kind,
		Metadata:     conn.MetaData,
		Status:       conn.Status,
		UserID:       &uuid.Nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	connectionCreated, err := l.ConnectionPersister.SaveConnection(connection)
	if err != nil {
		return nil, err
	}
	return connectionCreated, nil
}

func (l *DefaultLocalProvider) GetConnections(_ *http.Request, userID string, page, pageSize int, search, order string, filter string, status []string, kind []string) (*connections.ConnectionPage, error) {
	connectionsPage, err := l.ConnectionPersister.GetConnections(search, order, page, pageSize, filter, status, kind)
	if err != nil {
		return nil, err
	}
	return connectionsPage, nil
}

func (l *DefaultLocalProvider) GetConnectionByID(token string, connectionID uuid.UUID) (*connections.Connection, int, error) {
	return nil, http.StatusForbidden, ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) GetConnectionByIDAndKind(token string, connectionID uuid.UUID, kind string) (*connections.Connection, int, error) {
	connection, err := l.ConnectionPersister.GetConnection(connectionID, kind)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}
	return connection, http.StatusOK, err
}

func (l *DefaultLocalProvider) GetConnectionsByKind(_ *http.Request, _ string, _, _ int, _, _, _ string) (*map[string]interface{}, error) {
	return nil, ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) GetConnectionsStatus(_ *http.Request, _ string) (*connections.ConnectionsStatusPage, error) {
	return nil, ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) UpdateConnection(_ *http.Request, connection *connections.Connection) (*connections.Connection, error) {
	return l.ConnectionPersister.UpdateConnectionByID(connection)
}

func (l *DefaultLocalProvider) UpdateConnectionStatusByID(token string, connectionID uuid.UUID, connectionStatus connections.ConnectionStatus) (*connections.Connection, int, error) {
	updatedConnection, err := l.ConnectionPersister.UpdateConnectionStatusByID(connectionID, connectionStatus)
	if err != nil {
		fmt.Println("on update connection error", err)
		return nil, http.StatusInternalServerError, err
	}
	return updatedConnection, http.StatusOK, nil
}

func (l *DefaultLocalProvider) UpdateConnectionById(req *http.Request, conn *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	connection := connections.Connection{
		ID:           conn.ID,
		Name:         conn.Name,
		Type:         conn.Type,
		SubType:      conn.SubType,
		Kind:         conn.Kind,
		Metadata:     conn.MetaData,
		Status:       conn.Status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CredentialID: *conn.CredentialID,
	}
	return l.UpdateConnection(req, &connection)
}

func (l *DefaultLocalProvider) DeleteConnection(_ *http.Request, connectionID uuid.UUID) (*connections.Connection, error) {
	connection := connections.Connection{ID: connectionID}
	return l.ConnectionPersister.DeleteConnection(&connection)
}

func (l *DefaultLocalProvider) DeleteMesheryConnection() error {
	mesheryConnectionID := uuid.FromStringOrNil(viper.GetString("INSTANCE_ID"))
	_, err := l.DeleteConnection(nil, mesheryConnectionID)
	return err
}

// GetGenericPersister - to return persister
func (l *DefaultLocalProvider) GetGenericPersister() *database.Handler {
	return l.GenericPersister
}

// SetKubeClient - to set meshery kubernetes client
func (l *DefaultLocalProvider) SetKubeClient(client *mesherykube.Client) {
	l.KubeClient = client
}

// GetKubeClient - to get meshery kubernetes client
func (l *DefaultLocalProvider) GetKubeClient() *mesherykube.Client {
	return l.KubeClient
}

// SeedContent- to seed various contents with local provider-like patterns, filters, and applications
func (l *DefaultLocalProvider) SeedContent(log logger.Handler) {
	seededUUIDs := make([]uuid.UUID, 0)
	seedContents := []string{"Pattern", "Filter"}
	nilUserID := ""
	for _, seedContent := range seedContents {
		go func(comp string, log logger.Handler, seededUUIDs *[]uuid.UUID) {
			names, content, err := getSeededComponents(comp, log)
			if err != nil {
				log.Error(ErrGettingSeededComponents(err, comp))
			} else {
				log.Info("seeding sample ", comp, "s")
				switch comp {
				case "Pattern":
					for i, name := range names {
						id, _ := uuid.NewV4()

						var pattern = &MesheryPattern{
							PatternFile: content[i],
							Name:        name,
							ID:          &id,
							UserID:      &nilUserID,
							Visibility:  Published,
							Location: map[string]interface{}{
								"host":   "",
								"path":   "",
								"type":   "local",
								"branch": "",
							},
						}
						log.Debug("seeding "+comp+": ", name)
						_, err := l.MesheryPatternPersister.SaveMesheryPattern(pattern)
						if err != nil {
							log.Error(ErrGettingSeededComponents(err, comp+"s"))
						}
						*seededUUIDs = append(*seededUUIDs, id)
					}
				case "Filter":
					for i, name := range names {
						id, _ := uuid.NewV4()
						var filter = &MesheryFilter{
							FilterFile: []byte(content[i]),
							Name:       name,
							ID:         &id,
							UserID:     &nilUserID,
							Visibility: Published,
							Location: map[string]interface{}{
								"host":   "",
								"path":   "",
								"type":   "local",
								"branch": "",
							},
						}
						log.Debug("seeding "+comp+": ", name)
						_, err := l.MesheryFilterPersister.SaveMesheryFilter(filter)
						if err != nil {
							log.Error(ErrGettingSeededComponents(err, comp+"s"))
						}
						*seededUUIDs = append(*seededUUIDs, id)

					}
				}
			}
		}(seedContent, log, &seededUUIDs)
	}

	// seed default organization
	go func() {
		id, _ := uuid.NewV4()
		org := &Organization{
			ID:          &id,
			Name:        "My Org",
			Country:     "",
			Region:      "",
			Description: "This is default organization",
			Owner:       uuid.Nil,
		}
		count, _ := l.OrganizationPersister.GetOrganizationsCount()
		if count == 0 {
			_, err := l.OrganizationPersister.SaveOrganization(org)
			if err != nil {
				log.Error(ErrGettingSeededComponents(err, "organization"))
			}
		}
	}()
}
func (l *DefaultLocalProvider) Cleanup() error {
	if err := l.MesheryK8sContextPersister.DB.Migrator().DropTable(&K8sContext{}); err != nil {
		return err
	}
	if err := l.MesheryK8sContextPersister.DB.Migrator().DropTable(&connections.Connection{}); err != nil {
		return err
	}
	if err := l.MesheryK8sContextPersister.DB.Migrator().DropTable(&MesheryPattern{}); err != nil {
		return err
	}
	if err := l.MesheryK8sContextPersister.DB.Migrator().DropTable(&MesheryApplication{}); err != nil {
		return err
	}

	if err := l.KeyPersister.DB.Migrator().DropTable(&Key{}); err != nil {
		return err
	}

	return l.MesheryK8sContextPersister.DB.Migrator().DropTable(&MesheryFilter{})
}

func (l *DefaultLocalProvider) SaveUserCredential(token string, credential *Credential) (*Credential, error) {
	result := l.GetGenericPersister().Table("credentials").Create(&credential)
	if result.Error != nil {
		return nil, fmt.Errorf("error saving user credentials: %v", result.Error)
	}
	return nil, nil
}

func (l *DefaultLocalProvider) GetCredentialByID(token string, credentialID uuid.UUID) (*Credential, int, error) {
	return nil, http.StatusForbidden, ErrLocalProviderSupport
}

func (l *DefaultLocalProvider) GetUserCredentials(_ *http.Request, userID string, page, pageSize int, search, order string) (*CredentialsPage, error) {
	result := l.GetGenericPersister().Select("*").Where("user_id=? and deleted_at is NULL", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		result = result.Where("(lower(name) like ?)", like)
	}

	result = result.Order(order)

	var count int64
	if err := result.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("error retrieving count of credentials for user id: %s - %v", userID, err)
	}

	var credentialsList []*Credential
	if count > 0 {
		if err := result.Offset(page * pageSize).Limit(pageSize).Find(&credentialsList).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("error retrieving credentials for user id: %s - %v", userID, err)
			}
		}
	}

	credentialsPage := &CredentialsPage{
		Credentials: credentialsList,
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  int(count),
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error getting user credentials: %v", result.Error)
	}
	return credentialsPage, nil
}

func (l *DefaultLocalProvider) UpdateUserCredential(_ *http.Request, credential *Credential) (*Credential, error) {
	updatedCredential := &Credential{}
	if err := l.GetGenericPersister().Model(*updatedCredential).Where("user_id = ? AND id = ? AND deleted_at is NULL", credential.UserID, credential.ID).Updates(credential); err != nil {
		return nil, fmt.Errorf("error updating user credential: %v", err)
	}

	if err := l.GetGenericPersister().Where("user_id = ? AND id = ?", credential.UserID, credential.ID).First(updatedCredential).Error; err != nil {
		return nil, fmt.Errorf("error getting updated user credential: %v", err)
	}
	return updatedCredential, nil
}

func (l *DefaultLocalProvider) DeleteUserCredential(_ *http.Request, credentialID uuid.UUID) (*Credential, error) {
	delCredential := &Credential{}
	if err := l.GetGenericPersister().Model(&Credential{}).Where("id = ?", credentialID).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, err
	}
	if err := l.GetGenericPersister().Where("id = ?", credentialID).First(delCredential).Error; err != nil {
		return nil, err
	}
	return delCredential, nil
}

// GetOrganizations returns the list of organizations
func (l *DefaultLocalProvider) GetOrganizations(_, page, pageSize, search, order, updatedAfter string) ([]byte, error) {
	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pg, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, ErrPageNumber(err)
	}

	pgs, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, ErrPageSize(err)
	}

	return l.OrganizationPersister.GetOrganizations(search, order, pg, pgs, updatedAfter)
}

// GetKeys returns the list of keys
func (l *DefaultLocalProvider) GetUsersKeys(_, _, _, search, order, updatedAfter string, _ string) ([]byte, error) {
	keys, err := l.KeyPersister.GetUsersKeys(search, order, updatedAfter)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// GetKey returns the key for the given keyID
func (l *DefaultLocalProvider) GetUsersKey(_ *http.Request, keyID string) ([]byte, error) {
	id := uuid.FromStringOrNil(keyID)
	return l.KeyPersister.GetUsersKey(id)
}

// SaveKey saves the given key with the provider
func (l *DefaultLocalProvider) SaveUsersKey(_ string, keys []*Key) ([]*Key, error) {
	// fmt.Printf(keys)
	return l.KeyPersister.SaveUsersKeys(keys)
}

func (l *DefaultLocalProvider) GetWorkspaces(_, page, pageSize, search, order, filter string, orgID string) ([]byte, error) {
	return l.WorkspacePersister.GetWorkspaces(orgID, search, order, page, pageSize, filter)
}

func (l *DefaultLocalProvider) GetWorkspaceByID(_ *http.Request, workspaceID string, _ string) ([]byte, error) {
	id := uuid.FromStringOrNil(workspaceID)
	return l.WorkspacePersister.GetWorkspaceByID(id)
}

func (l *DefaultLocalProvider) DeleteWorkspace(_ *http.Request, workspaceID string) ([]byte, error) {
	id := uuid.FromStringOrNil(workspaceID)
	return l.WorkspacePersister.DeleteWorkspaceByID(id)
}

func (l *DefaultLocalProvider) SaveWorkspace(_ *http.Request, workspacePayload *v1beta1.WorkspacePayload, _ string, _ bool) ([]byte, error) {
	orgId, _ := uuid.FromString(workspacePayload.OrganizationID)
	workspace := &v1beta1.Workspace{
		CreatedAt:      time.Now(),
		Description:    workspacePayload.Description,
		Name:           workspacePayload.Name,
		OrganizationId: orgId,
		Owner:          "Meshery",
		UpdatedAt:      time.Now(),
	}
	return l.WorkspacePersister.SaveWorkspace(workspace)
}

func (l *DefaultLocalProvider) UpdateWorkspace(_ *http.Request, workspacePayload *v1beta1.WorkspacePayload, workspaceID string) (*v1beta1.Workspace, error) {
	id, _ := uuid.FromString(workspaceID)
	orgId, _ := uuid.FromString(workspacePayload.OrganizationID)
	workspace := &v1beta1.Workspace{
		ID:             id,
		CreatedAt:      time.Now(),
		Description:    workspacePayload.Description,
		Name:           workspacePayload.Name,
		OrganizationId: orgId,
		Owner:          "Meshery",
		UpdatedAt:      time.Now(),
	}
	return l.WorkspacePersister.UpdateWorkspaceByID(workspace)
}

func (l *DefaultLocalProvider) AddEnvironmentToWorkspace(_ *http.Request, workspaceID string, environmentID string) ([]byte, error) {
	workspaceId, _ := uuid.FromString(workspaceID)
	envId, _ := uuid.FromString(environmentID)
	return l.WorkspacePersister.AddEnvironmentToWorkspace(workspaceId, envId)
}

func (l *DefaultLocalProvider) RemoveEnvironmentFromWorkspace(_ *http.Request, workspaceID string, environmentID string) ([]byte, error) {
	workspaceId, _ := uuid.FromString(workspaceID)
	envId, _ := uuid.FromString(environmentID)
	return l.WorkspacePersister.DeleteEnvironmentFromWorkspace(workspaceId, envId)
}

func (l *DefaultLocalProvider) GetEnvironmentsOfWorkspace(_ *http.Request, workspaceID, page, pageSize, search, order, filter string) ([]byte, error) {
	workspaceId, _ := uuid.FromString(workspaceID)
	return l.WorkspacePersister.GetWorkspaceEnvironments(workspaceId, search, order, page, pageSize, filter)
}

func (l *DefaultLocalProvider) AddDesignToWorkspace(_ *http.Request, workspaceID string, designID string) ([]byte, error) {
	workspaceId, _ := uuid.FromString(workspaceID)
	designId, _ := uuid.FromString(designID)
	return l.WorkspacePersister.AddDesignToWorkspace(workspaceId, designId)
}

func (l *DefaultLocalProvider) RemoveDesignFromWorkspace(_ *http.Request, workspaceID string, designID string) ([]byte, error) {
	workspaceId, _ := uuid.FromString(workspaceID)
	designId, _ := uuid.FromString(designID)
	return l.WorkspacePersister.DeleteDesignFromWorkspace(workspaceId, designId)
}

func (l *DefaultLocalProvider) GetDesignsOfWorkspace(_ *http.Request, workspaceID, page, pageSize, search, order, filter string) ([]byte, error) {
	workspaceId, _ := uuid.FromString(workspaceID)
	return l.WorkspacePersister.GetWorkspaceDesigns(workspaceId, search, order, page, pageSize, filter)
}

// GetOrganization returns the organization for the given organizationID
func (l *DefaultLocalProvider) GetOrganization(_ *http.Request, organizationId string) ([]byte, error) {
	id := uuid.FromStringOrNil(organizationId)
	return l.OrganizationPersister.GetOrganzation(id)
}

// SaveOrganization saves given organization with the provider
func (l *DefaultLocalProvider) SaveOrganization(_ string, organization *Organization) ([]byte, error) {
	return l.OrganizationPersister.SaveOrganization(organization)
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
	ghWalker := walker.NewGit()
	result := make([]MesheryPattern, 0)
	err := ghWalker.Owner(owner).
		Repo(repo).
		Branch(branch).
		Root(path).
		RegisterFileInterceptor(func(f walker.File) error {
			ext := filepath.Ext(f.Name)
			if ext == ".yml" || ext == ".yaml" {
				name, err := GetPatternName(string(f.Content))
				if err != nil {
					return err
				}

				pf := MesheryPattern{
					Name:        name,
					PatternFile: f.Content,
					Location: map[string]interface{}{
						"type":   "github",
						"host":   fmt.Sprintf("github.com/%s/%s", owner, repo),
						"path":   f.Path,
						"branch": branch,
					},
				}

				mu.Lock()
				result = append(result, pf)
				mu.Unlock()
			}

			return nil
		}).Walk()
	return result, err
}

func githubRepoFilterScan(
	owner,
	repo,
	path,
	branch string,
) ([]MesheryFilter, error) {
	var mu sync.Mutex
	ghWalker := walker.NewGit()
	result := make([]MesheryFilter, 0)

	err := ghWalker.
		Owner(owner).
		Repo(repo).
		Branch(branch).
		Root(path).
		RegisterFileInterceptor(func(f walker.File) error {
			ext := filepath.Ext(f.Name)
			if ext == ".yml" || ext == ".yaml" {
				name, err := GetFilterName(string(f.Content))
				if err != nil {
					return err
				}

				ff := MesheryFilter{
					Name:       name,
					FilterFile: []byte(f.Content),
					Location: map[string]interface{}{
						"type":   "github",
						"host":   fmt.Sprintf("github.com/%s/%s", owner, repo),
						"path":   f.Path,
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

func genericHTTPPatternFile(fileURL string, log logger.Handler) ([]MesheryPattern, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file not found")
	}

	defer SafeClose(resp.Body, log)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var patternFile pattern.PatternFile
	err = yaml.Unmarshal(body, &patternFile)
	if err != nil {
		err = errors.Wrapf(err, "error decoding design file")
		return nil, utils.ErrDecodeYaml(err)
	}

	pf := MesheryPattern{
		Name:        patternFile.Name,
		PatternFile: string(body),
		Location: map[string]interface{}{
			"type":   "http",
			"host":   fileURL,
			"path":   "",
			"branch": "",
		},
	}

	return []MesheryPattern{pf}, nil
}

func genericHTTPFilterFile(fileURL string, log logger.Handler) ([]MesheryFilter, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file not found")
	}

	defer SafeClose(resp.Body, log)

	body, err := io.ReadAll(resp.Body)
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
		FilterFile: []byte(result),
		Location: map[string]interface{}{
			"type":   "http",
			"host":   fileURL,
			"path":   "",
			"branch": "",
		},
	}

	return []MesheryFilter{ff}, nil
}

// getSeededComponents reads the directory recursively looking for seed content
// Note- This function does not throw meshkit errors because the only method that calls it,"SeedContent" wraps the errors in meshkit errors.
// If this function is reused somewhere else, make sure to wrap its errors in appropriate meshkit errors, otherwise it can cause can a panic.
func getSeededComponents(comp string, log logger.Handler) ([]string, []string, error) {
	wd := utils.GetHome()
	switch comp {
	case "Pattern":
		wd = filepath.Join(wd, ".meshery", "content", "patterns")
	case "Filter":
		wd = filepath.Join(wd, ".meshery", "content", "filters", "binaries")

	}
	_, err := os.Stat(wd)
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, err
	} else if os.IsNotExist(err) {
		log.Info("creating directories for seeding... ", wd)
		er := os.MkdirAll(wd, 0777)
		if er != nil {
			return nil, nil, er
		}
	}
	if !viper.GetBool("SKIP_DOWNLOAD_CONTENT") {
		err = downloadContent(comp, wd)
		if err != nil {
			log.Error(ErrDownloadingSeededComponents(err, comp))
		}
	}
	log.Info("extracting "+comp+"s from ", wd)
	var names []string
	var contents []string
	err = filepath.WalkDir(wd,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				file, err := os.OpenFile(path, os.O_RDONLY, 0444)
				if err != nil {
					return err
				}
				content, err := io.ReadAll(file)
				if err != nil {
					return err
				}
				names = append(names, d.Name())
				contents = append(contents, string(content))
			}
			return nil
		})
	if err != nil {
		return nil, nil, err
	}
	return names, contents, nil
}

// Below helper functions are for downloading seed content and can be re used in future as the way of extracting them changes, like endpoint changes or so. That is why, they are encapsulated into general functions
func downloadContent(comp string, downloadpath string) error {
	switch comp {
	case "Pattern":
		walk := walker.NewGithub()
		return walk.Owner("service-mesh-patterns").Repo("service-mesh-patterns").Root("samples/").Branch("master").RegisterFileInterceptor(func(gca walker.GithubContentAPI) error {
			path := filepath.Join(downloadpath, gca.Name)
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()
			content, err := base64.StdEncoding.DecodeString(gca.Content)
			if err != nil {
				return err
			}
			fmt.Fprintf(file, "%s", content)
			return nil
		}).Walk()
	case "Filter":
		return getFiltersFromWasmFiltersRepo(downloadpath)

	}
	return nil
}

func getFiltersFromWasmFiltersRepo(downloadPath string) error {
	// releaseName, err := getLatestStableReleaseTag()
	// if err != nil {
	// 	return err
	// }
	//Temporary hardcoding until https://github.com/layer5io/wasm-filters/issues/38 is resolved
	downloadURL := "https://github.com/layer5io/wasm-filters/releases/download/v0.1.0/wasm-filters-v0.1.0.tar.gz"
	res, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	gzipStream := res.Body
	return extractTarGz(gzipStream, downloadPath)
}
func extractTarGz(gzipStream io.Reader, downloadPath string) error {
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
		switch header.Typeflag {
		case tar.TypeReg:
			if strings.HasSuffix(header.Name, ".wasm") {
				outFile, err := os.Create(filepath.Join(downloadPath, header.Name))
				if err != nil {
					return err
				}
				if _, err := io.Copy(outFile, tarReader); err != nil {
					return err
				}
				outFile.Close()
			}
		}
	}
	return nil
}

// // GetLatestStableReleaseTag fetches and returns the latest release tag from GitHub
// func getLatestStableReleaseTag() (string, error) {
// 	url := "https://github.com/layer5io/wasm-filters/releases/latest"
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", errors.New("failed to get latest stable release tag")
// 	}
// 	defer SafeClose(resp.Body)

// 	if resp.StatusCode != http.StatusOK {
// 		return "", errors.New("failed to get latest stable release tag")
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", errors.New("failed to get latest stable release tag")
// 	}
// 	re := regexp.MustCompile("/releases/tag/(.*?)\"")
// 	releases := re.FindAllString(string(body), -1)
// 	latest := strings.ReplaceAll(releases[0], "/releases/tag/", "")
// 	latest = strings.ReplaceAll(latest, "\"", "")
// 	return latest, nil
// }
