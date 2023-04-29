package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/internal/graphql"
	meshmodelhelper "github.com/layer5io/meshery/server/meshmodel"
	"github.com/layer5io/meshery/server/models"
	mesherymeshmodel "github.com/layer5io/meshery/server/models/meshmodel"
	"github.com/layer5io/meshery/server/router"
	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/utils/broadcast"
	"github.com/layer5io/meshkit/utils/events"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/spf13/viper"

	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
	"golang.org/x/exp/slog"
)

var (
	globalTokenForAnonymousResults string
	version                        = "Not Set"
	commitsha                      = "Not Set"
	releasechannel                 = "Not Set"
	serviceName                    = "meshery"
)

const (
	// DefaultProviderURL is the provider url for the "none" provider
	DefaultProviderURL           = "https://meshery.layer5.io"
	ArtifactHubComponentsHandler = "kubernetes" //The components generated in output directory will be handled by kubernetes
)

func main() {
	if globalTokenForAnonymousResults != "" {
		models.GlobalTokenForAnonymousResults = globalTokenForAnonymousResults
	}

	ctx := context.Background()

	log, err := logger.New(serviceName, logger.Options{
		Format:     logger.JSONLogFormat,
		DebugLevel: true,
		Output:     os.Stderr,
	})
	if err != nil {
		log.Error("Failed to create logger", err)
	}

	// print stack gracefully
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("A panic occurred:", "error", fmt.Sprintf("%v", r))
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	instanceID, err := uuid.NewV4()
	if err != nil {
		// log.Error(ErrCreatingUUIDInstance(err))
		os.Exit(1)
	}

	viper.AutomaticEnv()
	defaults := map[string]interface{}{
		"PORT":                  8080,
		"ADAPTER_URLS":          "",
		"BUILD":                 version,
		"OS":                    "meshery",
		"COMMITSHA":             commitsha,
		"RELEASE_CHANNEL":       releasechannel,
		"ISNTANCE_ID":           &instanceID,
		"PROVIDER":              "",
		"REGISTER_STATIC_K8S":   true,
		"SKIP_DOWNLOAD_CONTENT": false,
		"SKIP_COMP_GEN":         false,
		"PLAYGROUND":            false,
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	// Provide version
	log.Info("Local Provider capabilities", slog.String("version", version))

	// Get current channel - beta or stable
	log.Info("Meshery Server release channel is", slog.String("releaseChannel", releasechannel))

	// create user data directory
	userDataFolder := viper.GetString("USER_DATA_FOLDER")
	if _, err := os.Stat(userDataFolder); os.IsNotExist(err) {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Errorf("Failed to get user home directory", "error", err.Error())
			os.Exit(1)
		}

		userDataFolder := viper.GetString("USER_DATA_FOLDER")
		if userDataFolder == "" {
			userDataFolder = filepath.Join(home, ".meshery", "config")
			viper.SetDefault("USER_DATA_FOLDER", userDataFolder)
		}

		log.Debug("User data directory created", slog.String("path", userDataFolder))

		errDir := os.MkdirAll(userDataFolder, os.ModePerm)
		if errDir != nil {
			log.Errorf("Failed to create user data directory", "error", errDir.Error())
			os.Exit(1)
		}
	} else if err != nil {
		log.Errorf("Failed to create user data directory", "error", err.Error())
		os.Exit(1)
	}
	log.Debugf("Meshery Database is at", "path", viper.GetString("USER_DATA_FOLDER"))

	// create kubeconfig data directory
	if _, err := os.Stat(viper.GetString("KUBECONFIG_FOLDER")); os.IsNotExist(err) {
		kubeDataFolder := viper.GetString("KUBECONFIG_FOLDER")
		if kubeDataFolder == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				log.Errorf("Failed to get user home directory", "error", err.Error())
				os.Exit(1)
			}
			kubeDataFolder = path.Join(home, ".kube")
			viper.SetDefault("KUBECONFIG_FOLDER", kubeDataFolder)
		}

		errDir := os.MkdirAll(kubeDataFolder, 0755)
		if errDir != nil {
			log.Errorf("Failed to create kubeconfig data directory", "error", errDir.Error())
			os.Exit(1)
		}
		// logger.Debug("Kubeconfig data directory created", slog.String("path", kubeDataFolder))
	} else if err != nil {
		log.Errorf("Failed to create kubeconfig data directory", "error", err.Error())
		os.Exit(1)
	}

	// Adapters Configuration
	adapterURLS := viper.GetStringSlice("ADAPTER_URLS")
	adapterTracker := helpers.NewAdaptersTracker(adapterURLS)
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

	// Persist data from preferences
	preferencePersister, err := models.NewMapPreferencePersister(log)
	if err != nil {
		//log.Error(ErrCreatingMapPreferencePersisterInstance(err))
		log.Errorf("Failed to create preference persister", "error", err.Error())
		os.Exit(1)
	}
	defer preferencePersister.ClosePersister()

	// Database Configuration
	dbHandler := models.GetNewDBInstance()

	// Registry Manager configuration
	regManager, err := meshmodel.NewRegistryManager(dbHandler)
	if err != nil {
		log.Errorf("Failed to create registry manager", "error", err.Error())
		// log.Error(ErrInitializingRegistryManager(err))
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
	)
	if err != nil {
		log.Errorf("Failed to auto migrate database", "error", err.Error())
		// log.Error(ErrDatabaseAutoMigration(err))
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

	//--- Handler Configuration
	hc := &models.HandlerConfig{
		Providers:                 provs,
		ProviderCookieName:        "meshery-provider",
		ProviderCookieDuration:    30 * 24 * time.Hour,
		PlaygroundBuild:           viper.GetBool("PLAYGROUND"),
		AdapterTracker:            adapterTracker,
		QueryTracker:              queryTracker,
		Queue:                     mainQueue,
		KubeConfigFolder:          viper.GetString("KUBECONFIG_FOLDER"),
		GrafanaClient:             models.NewGrafanaClient(),
		GrafanaClientForQuery:     models.NewGrafanaClientWithHTTPClient(&http.Client{Timeout: time.Second}),
		PrometheusClient:          models.NewPrometheusClient(),
		PrometheusClientForQuery:  models.NewPrometheusClientWithHTTPClient(&http.Client{Timeout: time.Second}),
		ConfigurationChannel:      models.NewConfigurationHelper(),
		DashboardK8sResourcesChan: models.NewDashboardK8sResourcesHelper(),
		MeshModelSummaryChannel:   mesherymeshmodel.NewSummaryHelper(),
		K8scontextChannel:         models.NewContextHelper(),
		OperatorTracker:           models.NewOperatorTracker(viper.GetBool("DISABLE_OPERATOR")),
	}

	//seed the local meshmodel components
	ch := meshmodelhelper.NewEntityRegistrationHelper(hc, regManager, log)

	// @TODO add these to two different goroutine to avoid race conditions
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Error occurred while publishing meshmodel summary: %v", r)
			}
		}()
		log.Debug("Publishing meshmodel summary...")
		hc.MeshModelSummaryChannel.Publish()
		log.Debug("Done publishing meshmodel summary")
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Error occurred while seeding components: %v", r)
			}
		}()
		log.Debug("Seeding components now...")
		ch.SeedComponents()
		log.Debug("Done seeding components")
	}()

	log.Debug("Waiting for goroutines to finish...")
	wg.Wait()

	log.Debug("All goroutines have finished")

	buf := make([]byte, 1<<20)
	runtime.Stack(buf, true)
	log.Debugf("Running goroutines:\n%s", buf)

	select {
	case <-time.After(10 * time.Second):
		log.Error("Timeout waiting for goroutines to finish", err)
		buf := make([]byte, 1<<20)
		runtime.Stack(buf, true)
		log.Debugf("Running goroutines:\n%s", buf)
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Debug("Context deadline exceeded")
			return
		}
	}

	// logger.Debug("Server has finished")

	lProv.SeedContent(log)
	provs[lProv.Name()] = lProv

	//--- Remote Provider
	RemoteProviderURLs := viper.GetStringSlice("REMOTE_PROVIDER_URLS")
	for _, providerurl := range RemoteProviderURLs {
		parsedURL, err := url.Parse(providerurl)
		if err != nil {
			//log.Error(ErrInvalidURLSkippingProvider(providerurl))
			log.Errorf("Failed to parse provider url", "error", err.Error())
			os.Exit(1)
		}
		log.Debugf("Initializing remote provider with URL: %s", parsedURL.String())
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

		log.Debugf("Syncing preferences for remote provider: %s", cp.Name())
		cp.SyncPreferences()
		defer func(providerName string) {
			log.Debugf("Stopping preference sync for remote provider: %s", providerName)
			cp.StopSyncPreferences()
		}(cp.Name())

		provs[cp.Name()] = cp
	}

	operatorDeploymentConfig := models.NewOperatorDeploymentConfig(adapterTracker)
	mctrlHelper := models.NewMesheryControllersHelper(log, operatorDeploymentConfig, dbHandler)
	k8sComponentsRegistrationHelper := models.NewComponentsRegistrationHelper(log)

	h := handlers.NewHandlerInstance(
		hc,
		meshsyncCh,
		log,
		brokerConn,
		k8sComponentsRegistrationHelper,
		mctrlHelper,
		dbHandler,
		events.NewEventStreamer(),
		regManager,
		viper.GetString("PROVIDER"),
	)

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
		log.Info("Meshery server listening on", slog.Int("port", port))
		if err := r.Run(); err != nil {
			// log.Error(ErrListenAndServe(err))
			log.Error("error with starting the server", err)
			os.Exit(1)
		}
	}()
	<-c

	// Clean-up registry manager
	regManager.Cleanup()
	log.Infof("Doing seeded content cleanup...")

	for _, p := range hc.Providers {
		// skipping none provider for now
		// so it doesn't throw error each server is stopped. Reason: support for none provider is not yet implemented
		if p.Name() != "None" {
			log.Infof("De-registering Meshery server.")
			err = p.DeleteMesheryConnection()
			if err != nil {
				log.Errorf("error cleaning", err)
			}
		}
	}
	err = lProv.Cleanup()
	if err != nil {
		// log.Error(ErrCleaningUpLocalProvider(err))
		log.Errorf("error with cleaning", ErrCleaningUpLocalProvider(err))
	}
	utils.DeleteSVGsFromFileSystem()
	log.Infof("Closing database instance...")
	err = dbHandler.DBClose()
	if err != nil {
		// log.Error(ErrClosingDatabaseInstance(err))
		log.Error("error closing database", ErrClosingDatabaseInstance(err))
	}

	log.Infof("Shutting down Meshery Server...")

}
