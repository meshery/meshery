package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/internal/graphql"
	"github.com/layer5io/meshery/server/internal/store"
	meshmodelhelper "github.com/layer5io/meshery/server/meshmodel"
	"github.com/layer5io/meshery/server/models"
	mesherymeshmodel "github.com/layer5io/meshery/server/models/meshmodel"
	"github.com/layer5io/meshery/server/router"
	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel/core/policies"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/utils/broadcast"
	"github.com/layer5io/meshkit/utils/events"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
)

var (
	globalTokenForAnonymousResults string
	version                        = "Not Set"
	commitsha                      = "Not Set"
	releasechannel                 = "Not Set"
)

const (
	// DefaultProviderURL is the provider url for the "none" provider
	DefaultProviderURL = "https://meshery.layer5.io"
	PoliciesPath       = "../meshmodel/kubernetes/policies"
	RelationshipsPath  = "../meshmodel/kubernetes/relationships"
)

func main() {
	if globalTokenForAnonymousResults != "" {
		models.GlobalTokenForAnonymousResults = globalTokenForAnonymousResults
	}

	// Initialize Logger instance
	log, err := logger.New("meshery", logger.Options{
		Format: logger.SyslogLogFormat,
	})
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

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
	viper.SetDefault("PROVIDER", "")
	viper.SetDefault("REGISTER_STATIC_K8S", true)
	viper.SetDefault("SKIP_DOWNLOAD_CONTENT", false)
	viper.SetDefault("SKIP_COMP_GEN", false)
	viper.SetDefault("PLAYGROUND", false)
	store.Initialize()

	log.Info("Local Provider capabilities are: ", version)

	// Get the channel
	log.Info("Meshery Server release channel is: ", releasechannel)

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

	log.Info("Meshery Database is at: ", viper.GetString("USER_DATA_FOLDER"))
	if viper.GetString("KUBECONFIG_FOLDER") == "" {
		if err != nil {
			log.Error(ErrRetrievingUserHomeDirectory(err))
			os.Exit(1)
		}
		viper.SetDefault("KUBECONFIG_FOLDER", path.Join(home, ".kube"))
	}
	log.Info("Using kubeconfig at: ", viper.GetString("KUBECONFIG_FOLDER"))

	if viper.GetBool("DEBUG") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	log.Info("Log level: ", logrus.GetLevel())

	adapterURLs := viper.GetStringSlice("ADAPTER_URLS")

	adapterTracker := helpers.NewAdaptersTracker(adapterURLs)
	queryTracker := helpers.NewUUIDQueryTracker()

	// Uncomment line below to generate a new UUID and force the user to login every time Meshery is started.
	// fileSessionStore := sessions.NewFilesystemStore("", []byte(uuid.NewV4().Bytes()))
	// fileSessionStore := sessions.NewFilesystemStore("", []byte("Meshery"))
	// fileSessionStore.MaxLength(0)

	QueueFactory := memqueue.NewFactory()
	mainQueue := QueueFactory.RegisterQueue(&taskq.QueueOptions{
		Name: "loadTestReporterQueue",
	})

	provs := map[string]models.Provider{}

	preferencePersister, err := models.NewMapPreferencePersister()
	if err != nil {
		log.Error(ErrCreatingMapPreferencePersisterInstance(err))
		os.Exit(1)
	}
	defer preferencePersister.ClosePersister()

	// eventsPersister, err := models.
	dbHandler := models.GetNewDBInstance()
	regManager, err := meshmodel.NewRegistryManager(dbHandler)
	if err != nil {
		log.Error(ErrInitializingRegistryManager(err))
		os.Exit(1)
	}
	meshsyncCh := make(chan struct{}, 10)
	brokerConn := nats.NewEmptyConnection

	err = dbHandler.AutoMigrate(
		&meshsyncmodel.KeyValue{},
		&meshsyncmodel.Object{},
		&meshsyncmodel.ResourceSpec{},
		&meshsyncmodel.ResourceStatus{},
		&meshsyncmodel.ResourceObjectMeta{},
		&models.PerformanceProfile{},
		&models.MesheryResult{},
		&models.MesheryPattern{},
		&models.MesheryFilter{},
		&models.PatternResource{},
		&models.MesheryApplication{},
		&models.UserPreference{},
		&models.PerformanceTestConfig{},
		&models.SmiResultWithID{},
		models.K8sContext{},
		models.Event{},
	)
	if err != nil {
		log.Error(ErrDatabaseAutoMigration(err))
		os.Exit(1)
	}

	lProv := &models.DefaultLocalProvider{
		ProviderBaseURL:                 DefaultProviderURL,
		MapPreferencePersister:          preferencePersister,
		ResultPersister:                 &models.MesheryResultsPersister{DB: dbHandler},
		SmiResultPersister:              &models.SMIResultsPersister{DB: dbHandler},
		TestProfilesPersister:           &models.TestProfilesPersister{DB: dbHandler},
		PerformanceProfilesPersister:    &models.PerformanceProfilePersister{DB: dbHandler},
		MesheryPatternPersister:         &models.MesheryPatternPersister{DB: dbHandler},
		MesheryFilterPersister:          &models.MesheryFilterPersister{DB: dbHandler},
		MesheryApplicationPersister:     &models.MesheryApplicationPersister{DB: dbHandler},
		MesheryPatternResourcePersister: &models.PatternResourcePersister{DB: dbHandler},
		MesheryK8sContextPersister:      &models.MesheryK8sContextPersister{DB: dbHandler},
		GenericPersister:                dbHandler,
	}
	lProv.Initialize()

	hc := &models.HandlerConfig{
		Providers:              provs,
		ProviderCookieName:     "meshery-provider",
		ProviderCookieDuration: 30 * 24 * time.Hour,
		PlaygroundBuild:        viper.GetBool("PLAYGROUND"),
		AdapterTracker:         adapterTracker,
		QueryTracker:           queryTracker,

		Queue: mainQueue,

		KubeConfigFolder: viper.GetString("KUBECONFIG_FOLDER"),

		GrafanaClient:         models.NewGrafanaClient(),
		GrafanaClientForQuery: models.NewGrafanaClientWithHTTPClient(&http.Client{Timeout: time.Second}),

		PrometheusClient:         models.NewPrometheusClient(),
		PrometheusClientForQuery: models.NewPrometheusClientWithHTTPClient(&http.Client{Timeout: time.Second}),

		ConfigurationChannel: models.NewConfigurationHelper(),

		DashboardK8sResourcesChan: models.NewDashboardK8sResourcesHelper(),
		MeshModelSummaryChannel:   mesherymeshmodel.NewSummaryHelper(),

		K8scontextChannel: models.NewContextHelper(),
		OperatorTracker:   models.NewOperatorTracker(viper.GetBool("DISABLE_OPERATOR")),
	}

	//seed the local meshmodel components
	ch := meshmodelhelper.NewEntityRegistrationHelper(hc, regManager, log)
	go func() {
		ch.SeedComponents()
		go hc.MeshModelSummaryChannel.Publish()
	}()

	lProv.SeedContent(log)
	provs[lProv.Name()] = lProv

	RemoteProviderURLs := viper.GetStringSlice("PROVIDER_BASE_URLS")
	for _, providerurl := range RemoteProviderURLs {
		parsedURL, err := url.Parse(providerurl)
		if err != nil {
			log.Error(ErrInvalidURLSkippingProvider(providerurl))
			continue
		}
		cp := &models.RemoteProvider{
			RemoteProviderURL:          parsedURL.String(),
			RefCookieName:              parsedURL.Host + "_ref",
			SessionName:                parsedURL.Host,
			TokenStore:                 make(map[string]string),
			LoginCookieDuration:        1 * time.Hour,
			SessionPreferencePersister: &models.SessionPreferencePersister{DB: dbHandler},
			ProviderVersion:            version,
			SmiResultPersister:         &models.SMIResultsPersister{DB: dbHandler},
			GenericPersister:           dbHandler,
		}

		cp.Initialize()

		cp.SyncPreferences()
		defer cp.StopSyncPreferences()
		provs[cp.Name()] = cp
	}

	operatorDeploymentConfig := models.NewOperatorDeploymentConfig(adapterTracker)
	mctrlHelper := models.NewMesheryControllersHelper(log, operatorDeploymentConfig, dbHandler)
	k8sComponentsRegistrationHelper := models.NewComponentsRegistrationHelper(log)
	rego, err := policies.NewRegoInstance(PoliciesPath, RelationshipsPath)
	if err != nil {
		logrus.Warn("error creating rego instance, policies will not be evaluated")
	}
	h := handlers.NewHandlerInstance(hc, meshsyncCh, log, brokerConn, k8sComponentsRegistrationHelper, mctrlHelper, dbHandler, events.NewEventStreamer(), regManager, viper.GetString("PROVIDER"), rego)

	b := broadcast.NewBroadcaster(100)
	defer b.Close()

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
		// skipping none provider for now
		// so it doesn't throw error each server is stopped. Reason: support for none provider is not yet implemented
		if p.Name() != "None" {
			log.Info("De-registering Meshery server.")
			err = p.DeleteMesheryConnection()
			if err != nil {
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
