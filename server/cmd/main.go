package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/meshery/schemas/models/core"

	"github.com/fsnotify/fsnotify"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/mesheryctl/pkg/constants"
	"github.com/meshery/meshery/server/handlers"
	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/internal/graphql"
	"github.com/meshery/meshery/server/internal/store"
	"github.com/meshery/meshery/server/machines"
	mhelpers "github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	mesherymeshmodel "github.com/meshery/meshery/server/models/meshmodel"
	"github.com/meshery/meshery/server/router"
	"github.com/meshery/meshkit/broker/nats"
	"github.com/meshery/meshkit/logger"
	_events "github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/models/meshmodel/core/policies"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/tracing"
	"github.com/meshery/meshkit/utils/broadcast"
	"github.com/meshery/meshkit/utils/events"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/meshery/schemas/models/v1beta1/workspace"
	schemasOrganization "github.com/meshery/schemas/models/v1beta2/organization"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	globalTokenForAnonymousResults string
	version                        = "Not Set"
	commitsha                      = "Not Set"
	releasechannel                 = "Not Set"
)

const (
	// DefaultProviderURL is the ProviderBaseURL stamped on the built-in local
	// provider. It is used for capability/package paths only; the local
	// provider does not authenticate against this URL. Sourced from the canonical
	// primary provider host (install/providers.env) so it cannot drift from the
	// PROVIDER_BASE_URLS default seeded into viper below.
	DefaultProviderURL = models.PrimaryProviderURL
	RelationshipsPath  = "../../models/kubernetes/"
)

func main() {
	if globalTokenForAnonymousResults != "" {
		models.GlobalTokenForAnonymousResults = globalTokenForAnonymousResults
	}

	viper.AutomaticEnv()

	// Meshery Server configuration
	viper.SetConfigFile("./server-config.env")
	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Errorf("error reading config %v", err)
	}

	logLevel := viper.GetInt("LOG_LEVEL")
	if viper.GetBool("DEBUG") {
		logLevel = int(logrus.DebugLevel)
	}
	logOption := logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
		// if debug, output caller
		EnableCallerInfo: logLevel == int(logrus.DebugLevel),
	}
	// Initialize Logger instance
	log, err := logger.New("meshery", logOption)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	viper.OnConfigChange(func(event fsnotify.Event) {
		log.Info("received change for", event.Name)
		log.SetLevel(logrus.Level(viper.GetInt("LOG_LEVEL")))
	})

	instanceID, err := uuid.NewV4()
	if err != nil {
		log.Error(ErrCreatingUUIDInstance(err))
		os.Exit(1)
	}

	// operatingSystem, err := exec.Command("uname", "-s").Output()
	// if err != nil {
	// 	logrus.Error(err)
	// }

	ctx := context.Background()

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("ADAPTER_URLS", "")
	viper.SetDefault("BUILD", version)
	viper.SetDefault("OS", "meshery")
	viper.SetDefault("COMMITSHA", commitsha)
	viper.SetDefault("RELEASE_CHANNEL", releasechannel)
	viper.SetDefault("INSTANCE_ID", &instanceID)
	viper.SetDefault(constants.ProviderENV, "")
	// Seed the canonical active remote-provider list (install/providers.env) so a
	// server started without PROVIDER_BASE_URLS still registers the default providers.
	viper.SetDefault(constants.ProviderURLsENV, models.DefaultRemoteProviderURLs)
	viper.SetDefault("REGISTER_STATIC_K8S", true)
	viper.SetDefault("SKIP_DOWNLOAD_CONTENT", false)
	viper.SetDefault("SKIP_DOWNLOAD_EXTENSIONS", false)
	viper.SetDefault("SKIP_COMP_GEN", false)
	viper.SetDefault("PLAYGROUND", false)
	viper.SetDefault("MESHSYNC_DEFAULT_DEPLOYMENT_MODE", connections.MeshsyncDeploymentModeDefault)
	store.Initialize()

	// initialize tracing. Skip entirely when OTEL_CONFIG is unset so local dev
	// doesn't pay for a failing OTLP gRPC exporter that logs
	// "traces export: ... connection refused" every ~10s.
	var tracingProvider *sdktrace.TracerProvider
	otelConfigString := strings.TrimSpace(viper.GetString("OTEL_CONFIG"))
	if otelConfigString == "" {
		log.Info("OpenTelemetry config not set; tracing disabled")
	} else {
		log.Info("Initializing OpenTelemetry tracing with config:", otelConfigString)
		provider, err := tracing.InitTracerFromYamlConfig(context.Background(), otelConfigString)
		if err != nil {
			log.Error(fmt.Errorf("failed to initialize OpenTelemetry tracing: %v", err))
		} else {
			tracingProvider = provider
			log.Info("OpenTelemetry tracing initialized with config:" + otelConfigString)
		}
	}
	// Defer shutdown of tracer provider
	defer func() {
		if tracingProvider != nil {
			if err := tracingProvider.Shutdown(context.Background()); err != nil {
				log.Error(fmt.Errorf("failed to shutdown OpenTelemetry tracer provider: %v", err))
			}
		}
	}()

	log.Info("Local Provider capabilities are: ", version)

	// Get the channel
	log.Info("Meshery Server release channel is: ", releasechannel)

	log.Infof("MESHSYNC_DEFAULT_DEPLOYMENT_MODE is %s", viper.GetString("MESHSYNC_DEFAULT_DEPLOYMENT_MODE"))

	home, err := os.UserHomeDir()
	if viper.GetString("USER_DATA_FOLDER") == "" {
		if err != nil {
			log.Error(ErrRetrievingUserHomeDirectory(err))
			os.Exit(1)
		}
		viper.SetDefault("USER_DATA_FOLDER", path.Join(home, ".meshery", "config"))
	}

	errDir := os.MkdirAll(viper.GetString("USER_DATA_FOLDER"), 0755)
	if errDir != nil {
		log.Error(ErrCreatingUserDataDirectory(viper.GetString("USER_DATA_FOLDER")))
		os.Exit(1)
	}
	logDir := path.Join(home, ".meshery", "logs", "registry")
	errDir = os.MkdirAll(logDir, 0755)
	if errDir != nil {
		logrus.Fatalf("Error creating user data directory: %v", err)
	}

	// Create or open the log file
	logFilePath := path.Join(logDir, "registry-logs.log")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		logrus.Fatalf("Could not create log file: %v", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Error(err)
		}
	}()
	viper.Set("REGISTRY_LOG_FILE", logFilePath)

	log.Info("Meshery Database is at: ", viper.GetString("USER_DATA_FOLDER"))
	if viper.GetString("KUBECONFIG_FOLDER") == "" {
		if err != nil {
			log.Error(ErrRetrievingUserHomeDirectory(err))
			os.Exit(1)
		}
		viper.SetDefault("KUBECONFIG_FOLDER", path.Join(home, ".kube"))
	}
	log.Info("Using kubeconfig at: ", viper.GetString("KUBECONFIG_FOLDER"))
	log.Info("Log level: ", log.GetLevel())

	adapterURLs := utils.SplitAndTrim(viper.GetString("ADAPTER_URLS"), ", \t\n\r")

	adapterTracker := helpers.NewAdaptersTracker(adapterURLs)

	// Uncomment line below to generate a new UUID and force the user to login every time Meshery is started.
	// fileSessionStore := sessions.NewFilesystemStore("", []byte(uuid.NewV4().Bytes()))
	// fileSessionStore := sessions.NewFilesystemStore("", []byte("Meshery"))
	// fileSessionStore.MaxLength(0)

	provs := map[string]models.Provider{}

	preferencePersister, err := models.NewMapPreferencePersister()
	if err != nil {
		log.Error(ErrCreatingMapPreferencePersisterInstance(err))
		os.Exit(1)
	}
	defer preferencePersister.ClosePersister()

	dbHandler := models.GetNewDBInstance()
	regManager, err := meshmodel.NewRegistryManager(dbHandler)
	if err != nil {
		log.Error(ErrInitializingRegistryManager(err))
		os.Exit(1)
	}
	meshsyncCh := make(chan struct{}, 10)
	brokerConn := nats.NewEmptyConnection

	err = dbHandler.AutoMigrate(
		&meshsyncmodel.KubernetesKeyValue{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
		&models.PerformanceProfile{},
		&models.MesheryResult{},
		&models.MesheryPattern{},
		&models.MesheryFilter{},
		&models.PatternResource{},
		&models.MesheryApplication{},
		&models.UserPreference{},
		&models.UserCapabilities{},
		&models.PerformanceTestConfig{},
		&models.SmiResultWithID{},
		models.K8sContext{},
		schemasOrganization.Organization{},
		models.Key{},
		&models.Credential{},
		connections.Connection{},
		environment.Environment{},
		environment.EnvironmentConnectionMapping{},
		workspace.Workspace{},
		workspace.WorkspacesEnvironmentsMapping{},
		workspace.WorkspacesDesignsMapping{},
		workspace.WorkspacesTeamsMapping{},
		workspace.WorkspacesViewsMapping{},
		_events.Event{},
		&models.SystemSetting{},
	)
	if err != nil {
		log.Error(ErrDatabaseAutoMigration(err))
		os.Exit(1)
	}

	meshsyncDefaultDeploymentMode := connections.MeshsyncDeploymentModeFromString(
		viper.GetString("MESHSYNC_DEFAULT_DEPLOYMENT_MODE"),
	)

	if meshsyncDefaultDeploymentMode == connections.MeshsyncDeploymentModeUndefined {
		meshsyncDefaultDeploymentMode = connections.MeshsyncDeploymentModeDefault
	}

	lProv := &models.DefaultLocalProvider{
		ProviderBaseURL:                 DefaultProviderURL,
		MapPreferencePersister:          preferencePersister,
		UserCapabilitiesPersister:       &models.UserCapabilitiesPersister{DB: dbHandler},
		ResultPersister:                 &models.MesheryResultsPersister{DB: dbHandler},
		SmiResultPersister:              &models.SMIResultsPersister{DB: dbHandler},
		TestProfilesPersister:           &models.TestProfilesPersister{DB: dbHandler},
		PerformanceProfilesPersister:    &models.PerformanceProfilePersister{DB: dbHandler},
		MesheryPatternPersister:         &models.MesheryPatternPersister{DB: dbHandler},
		MesheryFilterPersister:          &models.MesheryFilterPersister{DB: dbHandler},
		MesheryApplicationPersister:     &models.MesheryApplicationPersister{DB: dbHandler},
		MesheryPatternResourcePersister: &models.PatternResourcePersister{DB: dbHandler},
		MesheryK8sContextPersister:      &models.MesheryK8sContextPersister{DB: dbHandler},
		OrganizationPersister:           &models.OrganizationPersister{DB: dbHandler},
		ConnectionPersister:             &models.ConnectionPersister{DB: dbHandler},
		EnvironmentPersister:            &models.EnvironmentPersister{DB: dbHandler},
		WorkspacePersister:              &models.WorkspacePersister{DB: dbHandler},
		KeyPersister:                    &models.KeyPersister{DB: dbHandler},
		EventsPersister:                 &models.EventsPersister{DB: dbHandler},
		GenericPersister:                dbHandler,
		Log:                             log,
		MeshsyncDefaultDeploymentMode:   meshsyncDefaultDeploymentMode,
	}

	// Local remote provider is initalized here.
	lProv.Initialize()

	hc := &models.HandlerConfig{
		Providers:              provs,
		ProviderCookieName:     "meshery-provider",
		ProviderCookieDuration: 30 * 24 * time.Hour,
		// ProviderTracker is set below, once provider registration is
		// complete and the tracker has been built around the final
		// `provs` map. The handlers must not access it before that
		// assignment, but every route is registered after.
		PlaygroundBuild: viper.GetBool("PLAYGROUND"),
		AdapterTracker:  adapterTracker,

		KubeConfigFolder: viper.GetString("KUBECONFIG_FOLDER"),

		PatternChannel:            models.NewBroadcaster("Patterns"),
		FilterChannel:             models.NewBroadcaster("Filters"),
		EventBroadcaster:          models.NewBroadcaster("Events"),
		DashboardK8sResourcesChan: models.NewDashboardK8sResourcesHelper(),
		MeshModelSummaryChannel:   mesherymeshmodel.NewSummaryHelper(),

		K8scontextChannel: models.NewContextHelper(),
		OperatorTracker:   models.NewOperatorTracker(viper.GetBool("DISABLE_OPERATOR")),
	}
	krh, err := models.NewKeysRegistrationHelper(dbHandler, log)
	if err != nil {
		log.Error(ErrInitializingKeysRegistration(err))
		os.Exit(1)
	}
	//seed the local meshmodel components
	rego := policies.Rego{}

	go func() {
		// This is where models are seeded from meshmodel directory to registry
		models.SeedComponents(log, hc, regManager)
		// Rego is intialized for passing of policy if the policies are made to be per model base this needs to be removed.
		r, err := policies.NewRegoInstance(models.PoliciesPath, regManager)
		if err != nil {
			log.Warn(handlers.ErrCreatingOPAInstance(err))
		} else {
			rego = *r
		}
		krh.SeedKeys(viper.GetString("KEYS_PATH"))
		hc.MeshModelSummaryChannel.Publish()
	}()

	lProv.SeedContent(log)
	provs[lProv.Name()] = lProv

	providerEnvVar := viper.GetString(constants.ProviderENV)
	RemoteProviderURLs := utils.SplitAndTrim(viper.GetString("PROVIDER_BASE_URLS"), ", \t\n\r")
	for _, providerurl := range RemoteProviderURLs {
		parsedURL, err := url.Parse(providerurl)
		if err != nil {
			log.Error(ErrInvalidURLSkippingProvider(providerurl))
			continue
		}
		cp := &models.RemoteProvider{
			RemoteProviderURL:             parsedURL.String(),
			RefCookieName:                 parsedURL.Host + "_ref",
			SessionName:                   parsedURL.Host,
			TokenStore:                    make(map[string]string),
			LoginCookieDuration:           1 * time.Hour,
			SessionPreferencePersister:    &models.SessionPreferencePersister{DB: dbHandler},
			UserCapabilitiesPersister:     &models.UserCapabilitiesPersister{DB: dbHandler},
			ProviderVersion:               version,
			SmiResultPersister:            &models.SMIResultsPersister{DB: dbHandler},
			GenericPersister:              dbHandler,
			EventsPersister:               &models.EventsPersister{DB: dbHandler},
			Log:                           log,
			CookieDuration:                24 * time.Hour,
			MeshsyncDefaultDeploymentMode: meshsyncDefaultDeploymentMode,
		}

		// Initialize stamps the minimum addressable metadata
		// (ProviderType, ProviderURL, fallback ProviderName=host) so the
		// provider is renderable in /api/providers immediately. The
		// capabilities HTTP probe runs later, concurrently, via
		// ProviderTracker.VerifyAll - a single unreachable remote can no
		// longer block server startup or the provider chooser.
		cp.Initialize()

		// Pick a registration key that survives two failure modes:
		//   1. /capabilities has not been probed yet (or will fail at
		//      probe time), so cp.Name() is the URL host fallback that
		//      Initialize stamps. Two remotes sharing the same host
		//      would collide - the parsed URL keeps each addressable.
		//   2. Two configured URLs report the SAME canonical name from
		//      /capabilities once the probe lands (for example, both
		//      cloud.meshery.io and cloud.acme.io currently return
		//      "Meshery"). Since we now key by the host-fallback name
		//      before the probe, the name collision shows up later, in
		//      VerifyAll; but at registration time we already disambiguate
		//      by URL host so both entries reach the tracker.
		// On collision (same host or same fallback name) we re-key the
		// earlier entry by its URL host and let the new entry claim the
		// canonical name. This preserves the historical "last URL wins
		// the canonical name" routing that integrations rely on (e.g.
		// tokens minted against cloud.acme.io expect the cookie value
		// "Meshery" to resolve to cloud.acme.io regardless of the
		// iteration order of PROVIDER_BASE_URLS), while still giving the
		// re-keyed provider an addressable home in /api/providers.
		key := cp.Name()
		if key == "" {
			key = parsedURL.Host
			log.Warnf("remote provider at %q produced an empty Name() after Initialize; registering it under host %q so it remains addressable.", providerurl, key)
		}
		if existing, ok := provs[key]; ok {
			existingHost := existing.GetProviderURL()
			if u, err := url.Parse(existingHost); err == nil && u.Host != "" {
				existingHost = u.Host
			}
			log.Warnf("provider name collision: %q is already registered for %q; re-keying the earlier entry as %q so both remain addressable, and giving the canonical name %q to %q. Ensure each remote provider's /capabilities returns a unique providerName.", key, existing.GetProviderURL(), existingHost, key, providerurl)
			provs[existingHost] = existing
			delete(provs, key)
		}
		provs[key] = cp
	}

	// All providers are registered. Build the tracker now and kick off
	// the boot-time availability probe + post-probe SyncPreferences
	// activation in a background goroutine, so a slow remote cannot
	// delay server startup. The probe publishes status events as each
	// remote settles, so the chooser shows the local provider (and any
	// already-reachable remotes) immediately and updates each remote
	// entry independently as its probe completes.
	providerTracker := models.NewProviderTracker(provs, log)
	hc.ProviderTracker = providerTracker
	go func() {
		providerTracker.VerifyAll(ctx)
		// SyncPreferences requires capabilities to have been loaded
		// (the SyncPrefs entry is what tells the goroutine the remote
		// accepts preference sync). Activate it AFTER the boot probe
		// so the in-loop call doesn't no-op on still-empty caps.
		for _, p := range provs {
			rp, ok := p.(*models.RemoteProvider)
			if !ok {
				continue
			}
			rp.SyncPreferences()
		}
	}()

	// Defer StopSyncPreferences for every remote, regardless of whether
	// SyncPreferences ever started. The Stop call is internally guarded
	// so it is a safe no-op when sync was never activated (probe failed,
	// shutdown ran before activation, or the remote does not advertise
	// SyncPrefs).
	for _, p := range provs {
		rp, ok := p.(*models.RemoteProvider)
		if !ok {
			continue
		}
		defer rp.StopSyncPreferences()
	}

	// Resolve the configured PROVIDER to Meshery's internal registration key.
	// Remote providers are registered under a stable key (typically URL host)
	// before their async /capabilities probe reveals the canonical
	// providerName, so a pre-selected remote such as PROVIDER=Meshery may need
	// one bounded probe here to avoid falling back to the chooser.
	resolveStart := time.Now()
	resolvedProviderKey, providerResolved := models.ResolveProviderKeyWithProbe(ctx, providerEnvVar, provs)
	if probeElapsed := time.Since(resolveStart); providerEnvVar != "" && probeElapsed > time.Second {
		// Surface the boot-time remote /capabilities probe so operators can
		// see why startup paused (each parallel probe waits up to 15s on an
		// unreachable configured remote before timing out).
		log.Infof("resolving configured PROVIDER %q required a remote capability probe that took %s", providerEnvVar, probeElapsed.Round(time.Millisecond))
	}
	if providerResolved {
		providerEnvVar = resolvedProviderKey
	} else {
		if providerEnvVar != "" {
			// Informational, not an error: a configured PROVIDER that
			// matches no registered provider is a valid fallback to the
			// chooser, but operators should be able to see why auto-select
			// did not engage.
			log.Infof("configured PROVIDER %q could not be resolved to any registered provider; falling back to the provider chooser", providerEnvVar)
		}
		providerEnvVar = ""
	}

	operatorDeploymentConfig := models.NewOperatorDeploymentConfig(adapterTracker)
	// this mctrlHelper is not used it is being recreated per connection entity
	mctrlHelper := models.NewMesheryControllersHelper(
		log,
		operatorDeploymentConfig,
		dbHandler,
		hc.EventBroadcaster,
		nil,
		&instanceID,
	)
	connToInstanceTracker := machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine, 0),
	}

	k8sComponentsRegistrationHelper := models.NewComponentsRegistrationHelper(log)

	models.InitMeshSyncRegistrationQueue()
	mhelpers.InitRegistrationHelperSingleton(dbHandler, log, &connToInstanceTracker, hc.EventBroadcaster)
	policies.SyncRelationship.Lock()
	h := handlers.NewHandlerInstance(hc, meshsyncCh, log, brokerConn, k8sComponentsRegistrationHelper, mctrlHelper, dbHandler, events.NewEventStreamer(), regManager, providerEnvVar, &rego, &connToInstanceTracker, meshsyncDefaultDeploymentMode)
	policies.SyncRelationship.Unlock()

	b := broadcast.NewBroadcaster(100)
	defer func() {
		if err := b.Close(); err != nil {
			log.Error(err)
		}
	}()

	g := graphql.New(graphql.Options{
		Config:      hc,
		Logger:      log,
		BrokerConn:  brokerConn,
		Broadcaster: b,
	})

	gp := graphql.NewPlayground(graphql.Options{
		URL: "/api/system/graphql/query",
	})

	port := viper.GetInt("PORT")
	r := router.NewRouter(ctx, h, port, g, gp)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		log.Info("Meshery Server listening on: ", port)
		if err := r.Run(); err != nil {
			log.Error(ErrListenAndServe(err))
			os.Exit(1)
		}
	}()
	<-c
	regManager.Cleanup()
	log.Info("Doing seeded content cleanup...")

	for _, p := range hc.Providers {
		// Skip the local provider: it does not register a remote session, so
		// there is nothing to de-register or log out on shutdown.
		if p.Name() != models.LocalProviderName {
			log.Info("De-registering Meshery server.")
			if err = p.DeleteMesheryConnection(); err != nil {
				log.Error(err)
				continue
			}
			// Logout follows deregistration so the session that authorized
			// the delete is revoked only after the connection is removed.
			log.Info("Logging out Meshery server session.")
			if err = p.LogoutMesheryServer(); err != nil {
				log.Error(err)
			}
		}
	}

	err = lProv.Cleanup()
	if err != nil {
		log.Error(ErrCleaningUpLocalProvider(err))
	}
	utils.DeleteSVGsFromFileSystem()
	log.Info("Closing database instance...")
	err = dbHandler.DBClose()
	if err != nil {
		log.Error(ErrClosingDatabaseInstance(err))
	}

	log.Info("Shutting down Meshery Server...")
}
