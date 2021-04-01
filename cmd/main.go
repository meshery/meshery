package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/internal/graphql"
	"github.com/layer5io/meshery/internal/store"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/models/oam"
	"github.com/layer5io/meshery/router"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
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

	// Initialize Logger instance
	log, err := logger.New("meshery", logger.Options{
		Format: logger.SyslogLogFormat,
	})
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	operatingSystem, err := exec.Command("uname", "-s").Output()
	if err != nil {
		logrus.Error(err)
	}

	ctx := context.Background()

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("ADAPTER_URLS", "")
	viper.SetDefault("BUILD", version)
	viper.SetDefault("OS", strings.ToLower(string(operatingSystem)))
	viper.SetDefault("COMMITSHA", commitsha)
	viper.SetDefault("RELEASE_CHANNEL", releasechannel)

	store.Initialize()

	// Register local OAM traits and workloads
	if err := oam.RegisterMesheryOAMTraits(); err != nil {
		logrus.Error(err)
	}
	if err := oam.RegisterMesheryOAMWorkloads(); err != nil {
		logrus.Error(err)
	}
	logrus.Info("Registered Meshery local Capabilities")

	// Get the channel
	logrus.Info("Meshery server current channel: ", releasechannel)

	home, err := os.UserHomeDir()
	if viper.GetString("USER_DATA_FOLDER") == "" {
		if err != nil {
			logrus.Fatalf("unable to retrieve the user's home directory: %v", err)
		}
		viper.SetDefault("USER_DATA_FOLDER", path.Join(home, ".meshery", "config"))
	}
	logrus.Infof("Using '%s' to store user data", viper.GetString("USER_DATA_FOLDER"))

	if viper.GetString("KUBECONFIG_FOLDER") == "" {
		if err != nil {
			logrus.Fatalf("unable to retrieve the user's home directory: %v", err)
		}
		viper.SetDefault("KUBECONFIG_FOLDER", path.Join(home, ".kube"))
	}
	logrus.Infof("Using '%s' as the folder to look for kubeconfig file", viper.GetString("KUBECONFIG_FOLDER"))

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

	smiResultPersister, err := models.NewBitCaskSmiResultsPersister(viper.GetString("USER_DATA_FOLDER"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer smiResultPersister.CloseResultPersister()

	testConfigPersister, err := models.NewBitCaskTestProfilesPersister(viper.GetString("USER_DATA_FOLDER"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer testConfigPersister.CloseTestConfigsPersister()

	dbHandler, err := database.New(database.Options{
		Filename: fmt.Sprintf("%s/mesherydb.sql", viper.GetString("USER_DATA_FOLDER")),
		Engine:   database.SQLITE,
		Logger:   log,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	kubeclient := mesherykube.Client{}
	meshsyncCh := make(chan struct{})

	err = dbHandler.AutoMigrate(
		meshsyncmodel.KeyValue{},
		meshsyncmodel.Object{},
		models.PerformanceProfile{},
		models.MesheryResult{},
		models.MesheryPattern{},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	lProv := &models.DefaultLocalProvider{
		ProviderBaseURL:              DefaultProviderURL,
		MapPreferencePersister:       preferencePersister,
		ResultPersister:              &models.MesheryResultsPersister{DB: &dbHandler},
		SmiResultPersister:           smiResultPersister,
		TestProfilesPersister:        testConfigPersister,
		PerformanceProfilesPersister: &models.PerformanceProfilePersister{DB: &dbHandler},
		MesheryPatternPersister:      &models.MesheryPatternPersister{DB: &dbHandler},
		GenericPersister:             dbHandler,
		GraphqlHandler: graphql.New(graphql.Options{
			Logger:          log,
			DBHandler:       &dbHandler,
			KubeClient:      &kubeclient,
			MeshSyncChannel: meshsyncCh,
		}),
		GraphqlPlayground: graphql.NewPlayground(graphql.Options{
			URL: "/api/system/graphql/query",
		}),
	}
	lProv.Initialize()
	provs[lProv.Name()] = lProv

	cPreferencePersister, err := models.NewBitCaskPreferencePersister(viper.GetString("USER_DATA_FOLDER"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer preferencePersister.ClosePersister()

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
			BitCaskPreferencePersister: cPreferencePersister,
			ProviderVersion:            "v0.3.14",
			SmiResultPersister:         smiResultPersister,
			GenericPersister:           dbHandler,
			GraphqlHandler: graphql.New(graphql.Options{
				Logger:          log,
				DBHandler:       &dbHandler,
				KubeClient:      &kubeclient,
				MeshSyncChannel: meshsyncCh,
			}),
			GraphqlPlayground: graphql.NewPlayground(graphql.Options{
				URL: "/api/system/graphql/query",
			}),
		}

		cp.Initialize()

		cp.SyncPreferences()
		defer cp.StopSyncPreferences()
		provs[cp.Name()] = cp
	}

	h := handlers.NewHandlerInstance(&models.HandlerConfig{
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
	}, &kubeclient, meshsyncCh, log)

	port := viper.GetInt("PORT")
	r := router.NewRouter(ctx, h, port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		logrus.Infof("Starting Server listening on :%d", port)
		if err := r.Run(); err != nil {
			logrus.Fatalf("ListenAndServe Error: %v", err)
		}
	}()
	<-c
	logrus.Info("Shutting down Meshery")
}
