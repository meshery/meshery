package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/layer5io/meshery/server/models"

	"github.com/ghodss/yaml"
	"github.com/open-policy-agent/opa/rego"
	"github.com/sirupsen/logrus"
)

// Define a simple policy struct
type Policy struct {
	Rules string `json:"rules"`
}

// TODO remove it
// func TestIt() {
// 	yamlData := "name: nginx-service.yml\nservices:\n  nginx:\n    name: nginx\n    type: Pod\n    apiVersion: v1\n    namespace: default\n    model: kubernetes\n    settings:\n      spec:\n        containers:\n        - image: nginx:stable\n          name: nginx\n          ports:\n          - container Port: 80\n            name: http-web-svc\n    traits:\n      meshmap:\n        edges: []\n        id: 19d2a03e-141d-41a9-bf7f-9ca12cf3b9be\n        label: nginx\n        meshmodel-metadata:\n          genealogy: \"\"\n          isCustomResource: false\n          isNamespaced: true\n          logoURL: https://github.com/cncf/artwork/blob/master/projects/kubernetes/icon/white/kubernetes-icon-white.svg\n          model: kubernetes\n          modelDisplayName: Kubernetes\n          primaryColor: '#326CE5'\n          published: true\n          secondaryColor: '#7aa1f0'\n          shape: round-rectangle\n          styleOverrides: \"\"\n          subCategory: Scheduling & Orchestration\n          svgColor: ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg\n          svgComplete: \"\"\n          svgWhite: ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg\n        position:\n          posX: 113.0763599675155\n          posY: -26.13721055079196\n  nginx-service:\n    name: nginx-service\n    type: Service\n    apiVersion: v1\n    namespace: default\n    model: kubernetes\n    settings:\n      spec:\n        ports:\n        - name: name-of-service-port\n          port: 80\n          protocol: TCP\n          target Port: http-web-svc\n        selector:\n          app.kubernetes.io/name: proxy\n    traits:\n      meshmap:\n        edges: []\n        id: 5b140a92-0f93-4500-b64a-ce703b49f166\n        label: nginx-service\n        meshmodel-metadata:\n          genealogy: \"\"\n          isCustomResource: false\n          isNamespaced: true\n          logoURL: https://github.com/cncf/artwork/blob/master/projects/kubernetes/icon/white/kubernetes-icon-white.svg\n          model: kubernetes\n          modelDisplayName: Kubernetes\n          primaryColor: '#326CE5'\n          published: true\n          secondaryColor: '#7aa1f0'\n          shape: round-triangle\n          styleOverrides: '{\"height\":16,\"width\":17,\"padding\":12,\"background-fit\":\"none\",\"background-position-y\":4.5}'\n          subCategory: Scheduling & Orchestration\n          svgColor: ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg\n          svgComplete: \"\"\n          svgWhite: ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg\n        position:\n          posX: -32.24910429536563\n          posY: 10.040920246138025\n"
// 	res := PolicyRelationshipRegoHandler([]byte(yamlData))
// 	fmt.Println("*******", *res)
// }

type RelationObject struct {
	DestinationId   string                 `json:"destination_id,omitempty"`
	DestinationName string                 `json:"destination_name,omitempty"`
	SourceId        string                 `json:"source_id,omitempty"`
	SourceName      string                 `json:"source_name,omitempty"`
	Port            map[string]interface{} `json:"port,omitempty"`
}

type NetworkPolicyRegoResponse struct {
	ServicePodRelationships        []RelationObject `json:"service_pod_relationships,omitempty"`
	ServiceDeploymentRelationships []RelationObject `json:"service_deployment_relationships,omitempty"`
}

func PolicyRelationshipRegoHandler(ctx context.Context, designFile []byte) *NetworkPolicyRegoResponse {
	// Load the policy
	policyFile, err := ioutil.ReadFile("../meshmodel/policies/network-policy.rego") // absolute path, needs to be changed
	if err != nil {
		logrus.Fatal("error reading rego file", err.Error())
	}

	// Initialize the policy
	policy := Policy{string(policyFile)}
	fmt.Println("policy relationship1")

	// Initialize the rego engine with the policy
	engine, err := rego.New(
		rego.Query("data.network_policy"),
		rego.Module("policy", policy.Rules),
	).PrepareForEval(ctx)

	if err != nil {
		panic(err)
	}

	var input map[string]interface{}
	err = yaml.Unmarshal(designFile, &input)
	if err != nil {
		logrus.Error("error unmarshalling design file format", err.Error())
	}

	eval_reponse, err := engine.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		fmt.Println("an error occured evaluating rego policy", err.Error())
	}

	// Check the result of the policy
	if len(eval_reponse) > 0 && len(eval_reponse[0].Expressions) > 0 {
		if eval_resp, ok := (eval_reponse[0].Expressions[0].Value).(map[string]interface{}); ok {
			return NewRelationPolicy(eval_resp)
		}
	} else {
		logrus.Error("Failed to evaluate policy")
	}
	return nil
}

func (h *Handler) HandleNetworkRelationship(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)

		rw.WriteHeader((http.StatusBadRequest))
		return
	}

	networkPolicy := PolicyRelationshipRegoHandler(context.Background(), body)
	ec := json.NewEncoder(rw)
	_ = ec.Encode(networkPolicy)
}

func NewRelationPolicy(networkResponse map[string]interface{}) *NetworkPolicyRegoResponse {
	var result NetworkPolicyRegoResponse

	b, err := json.Marshal(networkResponse)
	if err != nil {
		logrus.Error("Error marhsalling json", err.Error())
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		logrus.Error("Error unmarshalling json into network response", err.Error())
	}

	return &result
}
