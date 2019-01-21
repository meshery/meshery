package handlers

import (
	"context"
	"net/http"

	"github.com/layer5io/meshery/appoptics"
	"github.com/layer5io/meshery/meshes"
	"github.com/sirupsen/logrus"
)

func aoDashRenderer(ctx context.Context, meshClient meshes.MeshClient, w http.ResponseWriter, req *http.Request) {
	token := req.FormValue("token")
	if token == "" {
		// http.Error(w, "Token not found", http.StatusNotFound)
		token = "ad6c84be90e7e16ef9150e0c0d809644956d5df6897b73d2340b3238fda40d9d"
	}
	spaceName := req.FormValue("dashboard")
	if spaceName == "" {
		spaceName = "istio"
	}
	logrus.Infof("retrieved token from query: %s", token)
	ao, err := appoptics.NewAOClient(token, spaceName)
	if err != nil {
		logrus.Errorf("error getting AO data: %v", err)
		http.Error(w, "unable to get data for the token", http.StatusNotFound)
		return
	}
	ad := ao.GenerateDataForTemplate()
	logrus.Infof("Retrieved AO data: %+#v", ad)

	ops, err := meshClient.Operations(ctx)
	if err != nil {
		logrus.Errorf("error getting operations for the mesh: %v", err)
		http.Error(w, "unable to render the page", http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"Ops":  ops,
		"AO":   ad,
		"Name": meshClient.MeshName(),
	}

	err = dashTempl.Execute(w, result)
	if err != nil {
		logrus.Errorf("error rendering the template for the dashboard: %v", err)
		http.Error(w, "unable to render the page", http.StatusInternalServerError)
		return
	}
}
