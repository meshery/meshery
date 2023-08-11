package models

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/model"
	SMP "github.com/layer5io/service-mesh-performance/spec"
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
	Navigator NavigatorExtensions `json:"navigator,omitempty"`
	UserPrefs UserPrefsExtensions `json:"user_prefs,omitempty"`
	GraphQL   GraphQLExtensions   `json:"graphql,omitempty"`
	Acccount  AccountExtensions   `json:"account,omitempty"`
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

type ConnectionPayload struct {
	Kind             string                 `json:"kind,omitempty"`
	SubType          string                 `json:"sub_type,omitempty"`
	Type             string                 `json:"type,omitempty"`
	MetaData         map[string]interface{} `json:"metadata,omitempty"`
	Status           ConnectionStatus       `json:"status,omitempty"`
	CredentialSecret map[string]interface{} `json:"credential_secret,omitempty"`
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

	PersistConnection Feature = "persist-connection"

	PersistCredentials Feature = "persist-credentials"

	UsersProfile Feature = "users-profile"

	UsersIdentity Feature = "users-identity"
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

	// UserPrefsCtxKey is the context key for persisting user preferences to context
	PerfObjCtxKey ContextKey = "perf_obj"

	KubeClustersKey   ContextKey = "kubeclusters"
	AllKubeClusterKey ContextKey = "allkubeclusters"

	MesheryControllerHandlersKey ContextKey = "mesherycontrollerhandlerskey"
	MeshSyncDataHandlersKey      ContextKey = "meshsyncdatahandlerskey"

	RegistryManagerKey ContextKey = "registrymanagerkey"

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

// Provider - interface for providers
type Provider interface {
	PreferencePersister

	// Initialize will initialize a provider instance
	// by loading its capabilities and other metadata in the memory
	Initialize()

	Name() string

	// Returns ProviderType
	GetProviderType() ProviderType

	PackageLocation() string

	GetProviderCapabilities(http.ResponseWriter, *http.Request)

	GetProviderProperties() ProviderProperties
	// InitiateLogin - does the needed check, returns a true to indicate "return" or false to continue
	InitiateLogin(http.ResponseWriter, *http.Request, bool)
	TokenHandler(http.ResponseWriter, *http.Request, bool)
	ExtractToken(http.ResponseWriter, *http.Request)
	GetSession(req *http.Request) error
	GetUserDetails(*http.Request) (*User, error)
	GetUserByID(req *http.Request, userID string) ([]byte, error)
	GetUsers(token, page, pageSize, search, order, filter string) ([]byte, error)
	GetProviderToken(req *http.Request) (string, error)
	UpdateToken(http.ResponseWriter, *http.Request) string
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

	SaveK8sContext(token string, k8sContext K8sContext) (K8sContext, error)
	GetK8sContexts(token, page, pageSize, search, order string) ([]byte, error)
	DeleteK8sContext(token, id string) (K8sContext, error)
	GetK8sContext(token, id string) (K8sContext, error)
	LoadAllK8sContext(token string) ([]*K8sContext, error)
	// SetCurrentContext(token, id string) (K8sContext, error)
	// GetCurrentContext(token string) (K8sContext, error)

	SMPTestConfigStore(req *http.Request, perfConfig *SMP.PerformanceTestConfig) (string, error)
	SMPTestConfigGet(req *http.Request, testUUID string) (*SMP.PerformanceTestConfig, error)
	SMPTestConfigFetch(req *http.Request, page, pageSize, search, order string) ([]byte, error)
	SMPTestConfigDelete(req *http.Request, testUUID string) error

	RecordMeshSyncData(model.Object) error
	ReadMeshSyncData() ([]model.Object, error)
	GetGenericPersister() *database.Handler

	SetKubeClient(client *mesherykube.Client)
	GetKubeClient() *mesherykube.Client

	SaveMesheryPattern(tokenString string, pattern *MesheryPattern) ([]byte, error)
	GetMesheryPatterns(tokenString, page, pageSize, search, order string, updatedAfter string) ([]byte, error)
	GetCatalogMesheryPatterns(tokenString string, page, pageSize, search, order string) ([]byte, error)
	PublishCatalogPattern(req *http.Request, publishPatternRequest *MesheryCatalogPatternRequestBody) ([]byte, error)
	UnPublishCatalogPattern(req *http.Request, publishPatternRequest *MesheryCatalogPatternRequestBody) ([]byte, error)
	DeleteMesheryPattern(req *http.Request, patternID string) ([]byte, error)
	DeleteMesheryPatterns(req *http.Request, patterns MesheryPatternDeleteRequestBody) ([]byte, error)
	CloneMesheryPattern(req *http.Request, patternID string, clonePatternRequest *MesheryClonePatternRequestBody) ([]byte, error)
	GetMesheryPattern(req *http.Request, patternID string) ([]byte, error)
	RemotePatternFile(req *http.Request, resourceURL, path string, save bool) ([]byte, error)
	SaveMesheryPatternResource(token string, resource *PatternResource) (*PatternResource, error)
	GetMesheryPatternResource(token, resourceID string) (*PatternResource, error)
	GetMesheryPatternResources(token, page, pageSize, search, order, name, namespace, typ, oamType string) (*PatternResourcePage, error)
	DeleteMesheryPatternResource(token, resourceID string) error

	SaveMesheryFilter(tokenString string, filter *MesheryFilter) ([]byte, error)
	GetMesheryFilters(tokenString, page, pageSize, search, order string) ([]byte, error)
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

	SavePerformanceProfile(tokenString string, performanceProfile *PerformanceProfile) ([]byte, error)
	GetPerformanceProfiles(tokenString string, page, pageSize, search, order string) ([]byte, error)
	GetPerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error)
	DeletePerformanceProfile(req *http.Request, performanceProfileID string) ([]byte, error)

	SaveSchedule(tokenString string, s *Schedule) ([]byte, error)
	GetSchedules(req *http.Request, page, pageSize, order string) ([]byte, error)
	GetSchedule(req *http.Request, scheduleID string) ([]byte, error)
	DeleteSchedule(req *http.Request, scheduleID string) ([]byte, error)

	ExtensionProxy(req *http.Request) (*ExtensionProxyResponse, error)

	SaveConnection(req *http.Request, conn *ConnectionPayload, token string, skipTokenCheck bool) error
	GetConnections(req *http.Request, userID string, page, pageSize int, search, order string) (*ConnectionPage, error)
	GetConnectionsByKind(req *http.Request, userID string, page, pageSize int, search, order, connectionKind string) (*map[string]interface{}, error)
	GetConnectionsStatus(req *http.Request, userID string) (*ConnectionsStatusPage, error)
	UpdateConnection(req *http.Request, conn *Connection) (*Connection, error)
	UpdateConnectionById(req *http.Request, conn *ConnectionPayload, connId string) (*Connection, error)
	DeleteConnection(req *http.Request, connID uuid.UUID) (*Connection, error)
	DeleteMesheryConnection() error

	SaveUserCredential(req *http.Request, credential *Credential) error
	GetUserCredentials(req *http.Request, userID string, page, pageSize int, search, order string) (*CredentialsPage, error)
	UpdateUserCredential(req *http.Request, credential *Credential) (*Credential, error)
	DeleteUserCredential(req *http.Request, credentialID uuid.UUID) (*Credential, error)
}
