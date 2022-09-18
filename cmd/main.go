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
	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/internal/graphql"
	"github.com/layer5io/meshery/internal/store"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/models/pattern/core"
	"github.com/layer5io/meshery/router"
	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/broadcast"
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
)

func main() {
	if globalTokenForAnonymousResults != "" {
		models.GlobalTokenForAnonymousResults = globalTokenForAnonymousResults
	}

	instanceID, err := uuid.NewV4()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// Initialize Logger instance
	log, err := logger.New("meshery", logger.Options{
		Format: logger.SyslogLogFormat,
	})
	if err != nil {
		logrus.Error(err)
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

	viper.SetDefault("SKIP_DOWNLOAD_CONTENT", false)
	viper.SetDefault("SKIP_COMP_GEN", false)
	store.Initialize()

	// Register local OAM traits and workloads
	if err := core.RegisterMesheryOAMTraits(); err != nil {
		logrus.Error(err)
	}
	if err := core.RegisterMesheryOAMWorkloads(); err != nil {
		logrus.Error(err)
	}
	logrus.Info("Local Provider capabilities are: ", version)

	// Get the channel
	logrus.Info("Meshery Server release channel is: ", releasechannel)

	home, err := os.UserHomeDir()
	if viper.GetString("USER_DATA_FOLDER") == "" {
		if err != nil {
			logrus.Fatalf("unable to retrieve the user's home directory: %v", err)
		}
		viper.SetDefault("USER_DATA_FOLDER", path.Join(home, ".meshery", "config"))
	}

	errDir := os.MkdirAll(viper.GetString("USER_DATA_FOLDER"), 0755)
	if errDir != nil {
		logrus.Fatalf("unable to create the directory for storing user data at %v", viper.GetString("USER_DATA_FOLDER"))
	}

	logrus.Infof("Meshery Database is at: %s", viper.GetString("USER_DATA_FOLDER"))
	if viper.GetString("KUBECONFIG_FOLDER") == "" {
		if err != nil {
			logrus.Fatalf("unable to retrieve the user's home directory: %v", err)
		}
		viper.SetDefault("KUBECONFIG_FOLDER", path.Join(home, ".kube"))
	}
	logrus.Infof("Using kubeconfig at: %s", viper.GetString("KUBECONFIG_FOLDER"))

	if viper.GetBool("DEBUG") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Infof("Log level: %s", logrus.GetLevel())

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
		logrus.Fatal(err)
	}
	defer preferencePersister.ClosePersister()

	dbHandler := models.GetNewDBInstance()
	if err != nil {
		logrus.Fatal(err)
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
	)
	if err != nil {
		logrus.Fatal(err)
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
	lProv.SeedContent(log)
	provs[lProv.Name()] = lProv

	RemoteProviderURLs := viper.GetStringSlice("PROVIDER_BASE_URLS")
	for _, providerurl := range RemoteProviderURLs {
		parsedURL, err := url.Parse(providerurl)
		if err != nil {
			logrus.Error(providerurl, "is invalid url skipping provider")
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

	hc := &models.HandlerConfig{
		Providers:              provs,
		ProviderCookieName:     "meshery-provider",
		ProviderCookieDuration: 30 * 24 * time.Hour,

		AdapterTracker: adapterTracker,
		QueryTracker:   queryTracker,

		Queue: mainQueue,

		KubeConfigFolder: viper.GetString("KUBECONFIG_FOLDER"),

		GrafanaClient:         models.NewGrafanaClient(),
		GrafanaClientForQuery: models.NewGrafanaClientWithHTTPClient(&http.Client{Timeout: time.Second}),

		PrometheusClient:         models.NewPrometheusClient(),
		PrometheusClientForQuery: models.NewPrometheusClientWithHTTPClient(&http.Client{Timeout: time.Second}),

		ConfigurationChannel: models.NewConfigurationHelper(),

		K8scontextChannel: models.NewContextHelper(),
	}

	operatorDeploymentConfig := models.NewOperatorDeploymentConfig(adapterTracker)
	mctrlHelper := models.NewMesheryControllersHelper(log, operatorDeploymentConfig, dbHandler)
	k8sComponentsRegistrationHelper := models.NewComponentsRegistrationHelper(log)

	h := handlers.NewHandlerInstance(hc, meshsyncCh, log, brokerConn, k8sComponentsRegistrationHelper, mctrlHelper, dbHandler)

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
		logrus.Infof("Meshery Server listening on :%d", port)
		if err := r.Run(); err != nil {
			logrus.Fatalf("ListenAndServe Error: %v", err)
		}
	}()
	<-c
	logrus.Info("Doing seeded content cleanup...")
	err = lProv.Cleanup()
	if err != nil {
		log.Error(err)
	}

	logrus.Info("Closing database instance...")
	err = dbHandler.DBClose()
	if err != nil {
		log.Error(err)
	}

	logrus.Info("Shutting down Meshery Server...")
}
