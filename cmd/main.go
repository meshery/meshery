package main

import (
	"context"

	"github.com/satori/go.uuid"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/meshes/istio"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/router"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)

	if viper.GetBool("DEBUG") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Infof("Log level: %s", logrus.GetLevel())

	byPassAuth := viper.GetBool("BYPASS_AUTH")

	saasBaseURL := viper.GetString("SAAS_BASE_URL")
	if saasBaseURL == "" && !byPassAuth {
		logrus.Fatalf("SAAS_BASE_URL environment variable not set.")
	}

	fortio := viper.GetString("FORTIO_URL")
	if fortio == "" {
		logrus.Fatalf("FORTIO_URL environment variable not set.")
	}

	loadTestURL := viper.GetString("LOAD_TEST_URL")
	if loadTestURL == "" {
		logrus.Fatalf("LOAD_TEST_URL environment variable not set.")
	}

	meshClient, err := istio.CreateIstioClient(ctx)
	if err != nil {
		logrus.Fatalf("Error creating an istio client: %v", err)
	}

	h := handlers.NewHandlerInstance(&models.HandlerConfig{
		ByPassAuth:  byPassAuth,
		FortioURL:   fortio,
		SaaSBaseURL: saasBaseURL,
		LoadTestURL: loadTestURL,

		RefCookieName: "meshery_ref",

		SessionName:  "meshery",
		SessionStore: sessions.NewCookieStore([]byte(uuid.NewV4().Bytes()), nil),

		SaaSTokenName: "meshery_saas",
	})

	port := viper.GetInt("PORT")
	r := router.NewRouter(ctx, h, meshClient, port)

	logrus.Infof("Starting Server listening on :%d", port)
	if err := r.Run(); err != nil {
		logrus.Fatalf("ListenAndServe Error: %v", err)
	}
}
