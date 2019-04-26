package main

import (
	"context"

	"github.com/layer5io/meshery/helpers"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/router"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("ADAPTER_URLS", "")

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

	// fileSessionStore := sessions.NewFilesystemStore("", []byte(uuid.NewV4().Bytes())) // this is making us re-initiate login after every restart
	fileSessionStore := sessions.NewFilesystemStore("", []byte("Meshery2019"))
	fileSessionStore.MaxLength(0)

	h := handlers.NewHandlerInstance(&models.HandlerConfig{
		SaaSBaseURL: saasBaseURL,

		RefCookieName: "meshery_ref",

		SessionName:  "meshery",
		SessionStore: fileSessionStore,

		SaaSTokenName: "meshery_saas",

		AdapterTracker: adapterTracker,
	})

	port := viper.GetInt("PORT")
	r := router.NewRouter(ctx, h, port)

	logrus.Infof("Starting Server listening on :%d", port)
	if err := r.Run(); err != nil {
		logrus.Fatalf("ListenAndServe Error: %v", err)
	}
}
