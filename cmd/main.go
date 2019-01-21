package main

import (
	"context"
	"net/http"
	"os"

	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/meshes/istio"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	twitterHost := os.Getenv("TWITTER_APP_HOST")
	if twitterHost == "" {
		logrus.Fatalf("TWITTER_APP_HOST environment variable not set.")
	}

	fortio := os.Getenv("FORTIO_URL")
	if fortio == "" {
		logrus.Fatalf("FORTIO_URL environment variable not set.")
	}

	productPageURL := os.Getenv("PRODUCT_PAGE_URL")
	if productPageURL == "" {
		logrus.Fatalf("PRODUCT_PAGE_URL environment variable not set.")
	}

	meshClient, err := istio.CreateIstioClient(ctx)
	if err != nil {
		logrus.Fatalf("Error creating an istio client: %v", err)
	}
	config := &handlers.ServerConfig{
		MeshClient: meshClient,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("Starting Server listening on %s", (":" + port))
	err = http.ListenAndServe(":"+port, handlers.New(ctx, config))
	if err != nil {
		logrus.Fatalf("ListenAndServe Error: %v", err)
	}
}
