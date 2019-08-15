package main

import (
	"context"
	"os"
	"path"

	"github.com/layer5io/meshery/helpers"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/router"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	"github.com/vmihailenco/taskq"
	"github.com/vmihailenco/taskq/memqueue"
)

func main() {
	ctx := context.Background()

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("ADAPTER_URLS", "")

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

	saasBaseURL := viper.GetString("SAAS_BASE_URL")
	if saasBaseURL == "" {
		logrus.Fatalf("SAAS_BASE_URL environment variable not set.")
	}

	adapterURLs := viper.GetStringSlice("ADAPTER_URLS")

	adapterTracker := helpers.NewAdaptersTracker(adapterURLs)
	queryTracker := helpers.NewUUIDQueryTracker()

	// fileSessionStore := sessions.NewFilesystemStore("", []byte(uuid.NewV4().Bytes())) // this is making us re-initiate login after every restart
	fileSessionStore := sessions.NewFilesystemStore("", []byte("Meshery2019"))
	fileSessionStore.MaxLength(0)

	queueFactory := memqueue.NewFactory()
	mainQueue := queueFactory.NewQueue(&taskq.QueueOptions{
		Name: "loadTestReporterQueue",
	})

	// sessionPersister := helpers.NewFileSessionPersister(viper.GetString("USER_DATA_FOLDER"))
	sessionPersister, err := helpers.NewBadgerSessionPersister(viper.GetString("USER_DATA_FOLDER"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer sessionPersister.Close()

	h := handlers.NewHandlerInstance(&models.HandlerConfig{
		SaaSBaseURL: saasBaseURL,

		RefCookieName: "meshery_ref",

		SessionName:  "meshery",
		SessionStore: fileSessionStore,

		SaaSTokenName: "meshery_saas",

		AdapterTracker: adapterTracker,
		QueryTracker:   queryTracker,

		Queue: mainQueue,

		SessionPersister: sessionPersister,

		KubeConfigFolder: viper.GetString("KUBECONFIG_FOLDER"),
	})

	port := viper.GetInt("PORT")
	r := router.NewRouter(ctx, h, port)

	// go func() {
	// 	err := mainQueue.Consumer().Start()
	// 	if err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// }()

	logrus.Infof("Starting Server listening on :%d", port)
	if err := r.Run(); err != nil {
		logrus.Fatalf("ListenAndServe Error: %v", err)
	}
}
