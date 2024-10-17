package models

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/meshery/schemas/models/v1beta1"
)

// ContextKey is a custom type for setting context key
type ContextKey string

// ExtensionInput - input for a plugin
type ExtensionInput struct {
	DBHandler       *database.Handler
	MeshSyncChannel chan struct{}
	Logger          logger.Handler
	BrokerConn      broker.Handler
}

// Router
type Router struct {
	HTTPHandler http.Handler
	Path        string
}

// ExtensionOutput - output for a plugin
type ExtensionOutput struct {
	Router *Router
}

// ProviderType - for representing provider types
type ProviderType string

// ProviderProperties represents the structure of properties that a provider has
type ProviderProperties struct {
	ProviderType        ProviderType     `json:"provider_type,omitempty"`
	PackageVersion      string           `json:"package_version,omitempty"`
	PackageURL          string           `json:"package_url,omitempty"`
	ProviderName        string           `json:"provider_name,omitempty"`
	ProviderDescription []string         `json:"provider_description,omitempty"`
	ProviderURL         string           `json:"provider_url,omitempty"`
	Extensions          Extensions       `json:"extensions,omitempty"`
	Capabilities        Capabilities     `json:"capabilities,omitempty"`
	RestrictedAccess    RestrictedAccess `json:"restrictedAccess,omitempty"`
}

type Adapters struct {
	Istio   bool `json:"istio,omitempty"`
	Citrix  bool `json:"citrix,omitempty"`
	Consul  bool `json:"consul,omitempty"`
	Cilium  bool `json:"cilium,omitempty"`
	AppMesh bool `json:"appMesh,omitempty"`
	Kuma    bool `json:"kuma,omitempty"`
	Linkerd bool `json:"linkerd,omitempty"`
	Nginx   bool `json:"nginx,omitempty"`
	NSM     bool `json:"nsm,omitempty"`
}

type Configuration struct {
	Designs      bool `json:"designs,omitempty"`
	Applications bool `json:"applications,omitempty"`
	Filters      bool `json:"filters,omitempty"`
}

type NavigatorComponents struct {
	Dashboard     bool          `json:"dashboard,omitempty"`
	Performance   bool          `json:"performance,omitempty"`
	Conformance   bool          `json:"conformance,omitempty"`
	Extensions    bool          `json:"extensions,omitempty"`
	Toggler       bool          `json:"toggler,omitempty"`
	Help          bool          `json:"help,omitempty"`
	Lifecycle     Adapters      `json:"lifecycle,omitempty"`
	Configuration Configuration `json:"configuration,omitempty"`
}

type HeaderComponents struct {
	ContextSwitcher bool `json:"contextSwitcher,omitempty"`
	Settings        bool `json:"settings,omitempty"`
	Notifications   bool `json:"notifications,omitempty"`
	Profile         bool `json:"profile,omitempty"` // todo: account can have other structs, if needed needs to expand
}

type MesheryUICapabilities struct {
	Navigator NavigatorComponents `json:"navigator,omitempty"`
	Header    HeaderComponents    `json:"header,omitempty"`
}

type RestrictedAccess struct {
	IsMesheryUIRestricted bool                  `json:"isMesheryUiRestricted"`
	AllowedComponents     MesheryUICapabilities `json:"allowedComponents,omitempty"`
}

// Extensions defines the UI extension points
type Extensions struct {
	Navigator    NavigatorExtensions    `json:"navigator,omitempty"`
	UserPrefs    UserPrefsExtensions    `json:"user_prefs,omitempty"`
	GraphQL      GraphQLExtensions      `json:"graphql,omitempty"`
	Acccount     AccountExtensions      `json:"account,omitempty"`
	Collaborator CollaboratorExtensions `json:"collaborator,omitempty"`
}

// NavigatorExtensions is a collection of NavigatorExtension
type NavigatorExtensions []NavigatorExtension

// UserPrefsExtensions is a collection of UserPrefsExtension
type UserPrefsExtensions []UserPrefsExtension

// GraphQLExtensions is a collection of GraphQLExtension endpoints
type GraphQLExtensions []GraphQLExtension

// NavigatorExtensions is a collection of AccountExtension
type AccountExtensions []AccountExtension

// CollaboratorExtension describes the Collaborator extension point in the UI
type CollaboratorExtensions []CollaboratorExtension

// GraphQLExtension describes the graphql server extension point in the backend
type GraphQLExtension struct {
	Component string `json:"component,omitempty"`
	Path      string `json:"path,omitempty"`
	Type      string `json:"type,omitempty"`
}

type DesignerComponents struct {
	Design      bool `json:"design,omitempty"`
	Application bool `json:"application,omitempty"`
	Filter      bool `json:"filter,omitempty"`
	Save        bool `json:"save,omitempty"`
	New         bool `json:"new,omitempty"`
	SaveAs      bool `json:"saveAs,omitempty"`
	Validate    bool `json:"validate,omitempty"`
	Deploy      bool `json:"deploy,omitempty"`
	Undeploy    bool `json:"unDeploy,omitempty"`
}

type MeshMapComponentSet struct {
	Designer   DesignerComponents `json:"designer,omitempty"`
	Visualizer bool               `json:"visualizer,omitempty"` // todo: create a component set for visualizer
}

// NavigatorExtension describes the Navigator extension point in the UI
type NavigatorExtension struct {
	Title           string              `json:"title,omitempty"`
	OnClickCallback int                 `json:"on_click_callback,omitempty"`
	Href            Href                `json:"href,omitempty"`
	Component       string              `json:"component,omitempty"`
	Icon            string              `json:"icon,omitempty"`
	Link            *bool               `json:"link,omitempty"`
	Show            *bool               `json:"show,omitempty"`
	Children        NavigatorExtensions `json:"children,omitempty"`
	Type            string              `json:"type,omitempty"`
	AllowedTo       MeshMapComponentSet `json:"allowedTo,omitempty"`
	IsBeta          *bool               `json:"isBeta,omitempty"`
}

// AccountExtension describes the Account extension point in the UI
type AccountExtension struct {
	Title           string            `json:"title,omitempty"`
	OnClickCallback int               `json:"on_click_callback,omitempty"`
	Href            Href              `json:"href,omitempty"`
	Component       string            `json:"component,omitempty"`
	Link            *bool             `json:"link,omitempty"`
	Show            *bool             `json:"show,omitempty"`
	Children        AccountExtensions `json:"children,omitempty"`
	Type            string            `json:"type,omitempty"`
}

// UserPrefsExtension describes the user preference extension point in the UI
type UserPrefsExtension struct {
	Component string `json:"component,omitempty"`
	Type      string `json:"type,omitempty"`
}

// CollaboratorsExtension is the struct for collaborators extension
type CollaboratorExtension struct {
	Component string `json:"component,omitempty"`
	Type      string `json:"type,omitempty"`
}

// Href describes a link along with its type
type Href struct {
	URI      string `json:"uri,omitempty"`
	External *bool  `json:"external,omitempty"`
}

// Capabilities is the collection of capability
type Capabilities []Capability

// Capability is a capability of Provider indicating whether a feature is present
type Capability struct {
	Feature  Feature `json:"feature,omitempty"`
	Endpoint string  `json:"endpoint,omitempty"`
}

// K8sContextResponse - struct of response sent by provider when requested to persist k8s config
type K8sContextPersistResponse struct {
	K8sContext K8sContext `json:"k8s_context,omitempty"`
	Inserted   bool       `json:"inserted,omitempty"`
}

type ExtensionProxyResponse struct {
	Body       []byte `json:"body,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
}

// Feature is a type to store the features of the provider
type Feature string

const (
	// SyncPrefs indicates the Preference Synchronization feature
	SyncPrefs Feature = "sync-prefs" // /user/preferences

	PersistResults Feature = "persist-results" // /results

	PersistResult Feature = "persist-result" // /result

	// PersistSMIResults Feature = "persist-smi-results" // /smi/results

	PersistSMIResults Feature = "persist-smi-results" // /smi/results

	PersistMetrics Feature = "persist-metrics" // /result/metrics

	PersistEvents Feature = "persist-events"

	PersistSMPTestProfile Feature = "persist-smp-test-profile" // /user/test-config

	PersistMesheryPatterns Feature = "persist-meshery-patterns" // /patterns

	PersistMesheryPatternResources Feature = "persist-meshery-pattern-resources" // /patterns/resources

	PersistMesheryFilters Feature = "persist-meshery-filters" // /filter

	PersistMesheryApplications Feature = "persist-meshery-applications" // /applications

	PersistPerformanceProfiles Feature = "persist-performance-profiles" // /user/performance/profile

	PersistSchedules Feature = "persist-schedules" // /user/schedules

	MesheryPatternsCatalog Feature = "meshery-patterns-catalog" // /patterns/catalog

	MesheryFiltersCatalog Feature = "meshery-filters-catalog" // /filters/catalog

	CloneMesheryPatterns Feature = "clone-meshery-patterns" // /patterns/clone

	CloneMesheryFilters Feature = "clone-meshery-filters" // /filters/clone

	ShareDesigns Feature = "share-designs"

	ShareFilters Feature = "share-filters"

	PersistConnection Feature = "persist-connection"

	PersistCredentials Feature = "persist-credentials"

	UsersProfile Feature = "users-profile"

	UsersIdentity Feature = "users-identity"

	UsersKeys Feature = "users-keys"

	PersistOrganizations Feature = "organizations"

	PersistEnvironments Feature = "environments"

	PersistWorkspaces Feature = "workspaces"

	PersistAnonymousUser Feature = "persist-anonymous-user"
)

const (
	// LocalProviderType - represents local providers
	LocalProviderType ProviderType = "local"

	// RemoteProviderType - represents cloud providers
	RemoteProviderType ProviderType = "remote"

	// ProviderCtxKey is the context key for persisting provider to context
	ProviderCtxKey ContextKey = "provider"

	// TokenCtxKey is the context key for persisting token to context
	TokenCtxKey ContextKey = "token"

	// UserCtxKey is the context key for persisting user to context
	UserCtxKey ContextKey = "user"

	// UserIDCtxKey is the context key for persisting userID to context
	UserIDCtxKey ContextKey = "user_id"

	// UserPrefsCtxKey is the context key for persisting user preferences to context
	PerfObjCtxKey ContextKey = "perf_obj"

	KubeClustersKey   ContextKey = "kubeclusters"
	AllKubeClusterKey ContextKey = "allkubeclusters"

	MesheryControllerHandlersKey ContextKey = "mesherycontrollerhandlerskey"
	MeshSyncDataHandlersKey      ContextKey = "meshsyncdatahandlerskey"

	RegistryManagerKey ContextKey = "registrymanagerkey"

	HandlerKey               ContextKey = "handlerkey"
	SystemIDKey              ContextKey = "systemidKey"
	MesheryServerURL         ContextKey = "mesheryserverurl"
	MesheryServerCallbackURL ContextKey = "mesheryservercallbackurl"
)

// IsSupported returns true if the given feature is listed as one of
// the capabilities of the provider
func (caps Capabilities) IsSupported(feature Feature) bool {
	for _, cap := range caps {
		if feature == cap.Feature {
			return true
		}
	}

	return false
}

// GetEndpointForFeature returns the endpoint for the given feature
//
// Existence of a feature DOES NOT guarantee that the endpoint would be a not empty
// string as some of the features may not require an endpoint
func (caps Capabilities) GetEndpointForFeature(feature Feature) (string, bool) {
	for _, cap := range caps {
		if feature == cap.Feature {
			return cap.Endpoint, true
		}
	}

	return "", false
}

func VerifyMesheryProvider(provider string, supportedProviders map[string]Provider) bool {
	for prov := range supportedProviders {
		if prov == provider {
			return true
		}
	}
	return false
}

// Provider - interface for providers
type Provider interface {
	PreferencePersister
	CapabilitiesPersister
	MesheryEvents

	// Initialize will initialize a provider instance
	// by loading its capabilities and other metadata in the memory
	Initialize()

	Name() string
	GetProviderURL() string

	// Returns ProviderType
	GetProviderType() ProviderType

	PackageLocation() string

	GetProviderCapabilities(http.ResponseWriter, *http.Request, string)

	SetProviderProperties(providerProperties ProviderProperties)
	GetProviderProperties() ProviderProperties
	// InitiateLogin - does the needed check, returns a true to indicate "return" or false to continue
	InitiateLogin(http.ResponseWriter, *http.Request, bool)
	TokenHandler(http.ResponseWriter, *http.Request, bool)
	ExtractToken(http.ResponseWriter, *http.Request)
	GetSession(req *http.Request) error
	GetUserDetails(*http.Request) (*User, error)
	GetUserByID(req *http.Request, userID string) ([]byte, error)
	GetUsers(token, page, pageSize, search, order, filter string) ([]byte, error)
	GetUsersKeys(token, page, pageSize, search, order, filter string, orgID string) ([]byte, error)
	GetProviderToken(req *http.Request) (string, error)
	UpdateToken(http.ResponseWriter, *http.Request) string
	SetJWTCookie(w http.ResponseWriter, token string)
	UnSetJWTCookie(w http.ResponseWriter)
	Logout(http.ResponseWriter, *http.Request) error
	HandleUnAuthenticated(w http.ResponseWriter, req *http.Request)
	FetchResults(tokenVal string, page, pageSize, search, order, profileID string) ([]byte, error)
	FetchAllResults(tokenVal string, page, pageSize, search, order, from, to string) ([]byte, error)
	PublishResults(req *http.Request, result *MesheryResult, profileID string) (string, error)
	FetchSmiResults(req *http.Request, page, pageSize, search, order string) ([]byte, error)
	FetchSmiResult(req *http.Request, page, pageSize, search, order string, resultID uuid.UUID) ([]byte, error)
	PublishSmiResults(result *SmiResult) (string, error)
	PublishMetrics(tokenVal string, data *MesheryResult) error
	GetResult(tokenVal string, resultID uuid.UUID) (*MesheryResult, error)
	RecordPreferences(req *http.Request, userID string, data *Preference) error

	// We persist the events generated by Meshery Server in the local database (PersistEvent func in MesheryEvents interface) and use it as a cache to prevent increasing load on remote provider.
	// As events are persisted in the local database below func is a no-op for local provider, for the case of remote provider we want to persist the events to the provider's database.
	PublishEventToProvider(tokenVal string, event events.Event) error

	SaveK8sContext(token string, k8sContext K8sContext) (connections.Connection, error)
	GetK8sContexts(token, page, pageSize, search, order string, withStatus string, withCredentials bool) ([]byte, error)
	DeleteK8sContext(token, id string) (K8sContext, error)
	GetK8sContext(token, connectionID string) (K8sContext, error)
	LoadAllK8sContext(token string) ([]*K8sContext, error)
	// SetCurrentContext(token, id string) (K8sContext, error)
	// GetCurrentContext(token string) (K8sContext, error)

	SMPTestConfigStore(req *http.Request, perfConfig *SMP.PerformanceTestConfig) (string, error)
	SMPTestConfigGet(req *http.Request, testUUID string) (*SMP.PerformanceTestConfig, error)
	SMPTestConfigFetch(req *http.Request, page, pageSize, search, order string) ([]byte, error)
	SMPTestConfigDelete(req *http.Request, testUUID string) error

	GetGenericPersister() *database.Handler

	SetKubeClient(client *mesherykube.Client)
	GetKubeClient() *mesherykube.Client

	SaveMesheryPattern(tokenString string, pattern *MesheryPattern) ([]byte, error)
	GetMesheryPatterns(tokenString, page, pageSize, search, order string, updatedAfter string, visbility []string, includeMetrics string) ([]byte, error)
	GetCatalogMesheryPatterns(tokenString string, page, pageSize, search, order string, includeMetrics string) ([]byte, error)
	PublishCatalogPattern(req *http.Request, publishPatternRequest *MesheryCatalogPatternRequestBody) ([]byte, error)
	UnPublishCatalogPattern(req *http.Request, publishPatternRequest *MesheryCatalogPatternRequestBody) ([]byte, error)
	DeleteMesheryPattern(req *http.Request, patternID string) ([]byte, error)
	DeleteMesheryPatterns(req *http.Request, patterns MesheryPatternDeleteRequestBody) ([]byte, error)
	CloneMesheryPattern(req *http.Request, patternID string, clonePatternRequest *MesheryClonePatternRequestBody) ([]byte, error)
	GetMesheryPattern(req *http.Request, patternID string, includeMetrics string) ([]byte, error)
	RemotePatternFile(req *http.Request, resourceURL, path string, save bool) ([]byte, error)
	SaveMesheryPatternResource(token string, resource *PatternResource) (*PatternResource, error)
	GetMesheryPatternResource(token, resourceID string) (*PatternResource, error)
	GetMesheryPatternResources(token, page, pageSize, search, order, name, namespace, typ, oamType string) (*PatternResourcePage, error)
	DeleteMesheryPatternResource(token, resourceID string) error
	SaveMesheryPatternSourceContent(token string, applicationID string, sourceContent []byte) error
	GetDesignSourceContent(token, patternID string) ([]byte, error)

	SaveMesheryFilter(tokenString string, filter *MesheryFilter) ([]byte, error)
	GetMesheryFilters(tokenString, page, pageSize, search, order string, visibility []string) ([]byte, error)
	GetCatalogMesheryFilters(tokenString string, page, pageSize, search, order string) ([]byte, error)
	PublishCatalogFilter(req *http.Request, publishFilterRequest *MesheryCatalogFilterRequestBody) ([]byte, error)
	UnPublishCatalogFilter(req *http.Request, publishFilterRequest *MesheryCatalogFilterRequestBody) ([]byte, error)
	DeleteMesheryFilter(req *http.Request, filterID string) ([]byte, error)
	CloneMesheryFilter(req *http.Request, filterID string, cloneFilterRequest *MesheryCloneFilterRequestBody) ([]byte, error)
	GetMesheryFilter(req *http.Request, filterID string) ([]byte, error)
	GetMesheryFilterFile(req *http.Request, filterID string) ([]byte, error)
	RemoteFilterFile(req *http.Request, resourceURL, path string, save bool, resource string) ([]byte, error)

	SaveMesheryApplication(tokenString string, application *MesheryApplication) ([]byte, error)
	SaveApplicationSourceContent(token string, applicationID string, sourceContent []byte) error
	GetApplicationSourceContent(req *http.Request, applicationID string) ([]byte, error)
	GetMesheryApplications(tokenString, page, pageSize, search, order string, updatedAfter string) ([]byte, error)
	DeleteMesheryApplication(req *http.Request, applicationID string) ([]byte, error)
	GetMesheryApplication(req *http.Request, applicationID string) ([]byte, error)
	ShareDesign(req *http.Request) (int, error)
	ShareFilter(req *http.Request) (int, error)

	SavePerformanceProfile(tokenString string, performanceProfile *PerformanceProfile) ([]byte, error)
	GetPerformanceProfiles(tokenString string, page, pageSize, search, order string) ([]byte, error)
	GetPerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error)
	DeletePerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error)

	SaveSchedule(tokenString string, s *Schedule) ([]byte, error)
	GetSchedules(req *http.Request, page, pageSize, order string) ([]byte, error)
	GetSchedule(req *http.Request, scheduleID string) ([]byte, error)
	DeleteSchedule(req *http.Request, scheduleID string) ([]byte, error)

	ExtensionProxy(req *http.Request) (*ExtensionProxyResponse, error)

	SaveConnection(conn *connections.ConnectionPayload, token string, skipTokenCheck bool) (*connections.Connection, error)
	GetConnections(req *http.Request, userID string, page, pageSize int, search, order string, filter string, status []string, kind []string) (*connections.ConnectionPage, error)
	GetConnectionByIDAndKind(token string, connectionID uuid.UUID, kind string) (*connections.Connection, int, error)
	GetConnectionByID(token string, connectionID uuid.UUID) (*connections.Connection, int, error)
	GetConnectionsByKind(req *http.Request, userID string, page, pageSize int, search, order, connectionKind string) (*map[string]interface{}, error)
	GetConnectionsStatus(req *http.Request, userID string) (*connections.ConnectionsStatusPage, error)
	UpdateConnection(req *http.Request, conn *connections.Connection) (*connections.Connection, error)
	UpdateConnectionById(req *http.Request, conn *connections.ConnectionPayload, connId string) (*connections.Connection, error)
	UpdateConnectionStatusByID(token string, connectionID uuid.UUID, connectionStatus connections.ConnectionStatus) (*connections.Connection, int, error)
	DeleteConnection(req *http.Request, connID uuid.UUID) (*connections.Connection, error)
	DeleteMesheryConnection() error

	SaveUserCredential(token string, credential *Credential) (*Credential, error)
	GetUserCredentials(req *http.Request, userID string, page, pageSize int, search, order string) (*CredentialsPage, error)
	GetCredentialByID(token string, credentialID uuid.UUID) (*Credential, int, error)
	UpdateUserCredential(req *http.Request, credential *Credential) (*Credential, error)
	DeleteUserCredential(req *http.Request, credentialID uuid.UUID) (*Credential, error)

	GetEnvironments(token, page, pageSize, search, order, filter, orgID string) ([]byte, error)
	GetEnvironmentByID(req *http.Request, environmentID, orgID string) ([]byte, error)
	SaveEnvironment(req *http.Request, env *v1beta1.EnvironmentPayload, token string, skipTokenCheck bool) ([]byte, error)
	DeleteEnvironment(req *http.Request, environmentID string) ([]byte, error)
	UpdateEnvironment(req *http.Request, env *v1beta1.EnvironmentPayload, environmentID string) (*v1beta1.Environment, error)
	AddConnectionToEnvironment(req *http.Request, environmentID string, connectionID string) ([]byte, error)
	RemoveConnectionFromEnvironment(req *http.Request, environmentID string, connectionID string) ([]byte, error)
	GetConnectionsOfEnvironment(req *http.Request, environmentID, page, pagesize, search, order, filter string) ([]byte, error)

	GetOrganizations(token, page, pageSize, search, order, filter string) ([]byte, error)

	GetWorkspaces(token, page, pagesize, search, order, filter, orgID string) ([]byte, error)
	GetWorkspaceByID(req *http.Request, workspaceID, orgID string) ([]byte, error)
	SaveWorkspace(req *http.Request, workspace *v1beta1.WorkspacePayload, token string, skipTokenCheck bool) ([]byte, error)
	DeleteWorkspace(req *http.Request, workspaceID string) ([]byte, error)
	UpdateWorkspace(req *http.Request, workspace *v1beta1.WorkspacePayload, workspaceID string) (*v1beta1.Workspace, error)
	GetEnvironmentsOfWorkspace(req *http.Request, workspaceID, page, pagesize, search, order, filter string) ([]byte, error)
	AddEnvironmentToWorkspace(req *http.Request, workspaceID string, environmentID string) ([]byte, error)
	RemoveEnvironmentFromWorkspace(req *http.Request, workspaceID string, environmentID string) ([]byte, error)
	GetDesignsOfWorkspace(req *http.Request, workspaceID, page, pagesize, search, order, filter string) ([]byte, error)
	AddDesignToWorkspace(req *http.Request, workspaceID string, designID string) ([]byte, error)
	RemoveDesignFromWorkspace(req *http.Request, workspaceID string, designID string) ([]byte, error)
}
