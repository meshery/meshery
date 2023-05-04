package handlers

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/open-policy-agent/opa/rego"
	"github.com/sirupsen/logrus"
)

// Define a simple policy struct
type Policy struct {
	Rules string `json:"rules"`
}

func main() {
	// Load the policy
	policyFile, err := ioutil.ReadFile("../meshmodel/policies/network-policy.rego")
	if err != nil {
		panic(err)
	}

	// Initialize the policy
	policy := Policy{string(policyFile)}

	// Initialize the rego engine with the policy
	engine, err := rego.New(
		rego.Query("result = data"),
		rego.Module("policy.rego", policy.Rules),
	).PrepareForEval(context.Background())

	if err != nil {
		panic(err)
	}

	rs, err := engine.Eval(context.Background(), rego.EvalInput( // SAMPLE DATA, HAS TO BE REMOVED
		[]byte(`{
			"name": "service-nginx-deployment",
			"services": {
					"nginx-1-pod": {
							"apiVersion": "v1",
							"model": "kubernetes",
							"name": "nginx-1",
							"namespace": "default",
							"settings": {
									"spec": {
											"containers": [
													{
															"image": "nginx:stable",
															"name": "nginx",
															"ports": [
																	{
																			"container Port": 80,
																			"name": "http-web-svc"
																	}
															]
													}
											]
									}
							},
							"traits": {
									"meshmap": {
											"edges": [],
											"id": "19d2a03e-141d-41a9-bf7f-9ca12cf3b9be",
											"label": "nginx",
											"meshmodel-metadata": {
													"genealogy": "",
													"isCustomResource": false,
													"isNamespaced": true,
													"logoURL": "https://github.com/cncf/artwork/blob/master/projects/kubernetes/icon/white/kubernetes-icon-white.svg",
													"model": "kubernetes",
													"modelDisplayName": "Kubernetes",
													"primaryColor": "#326CE5",
													"published": true,
													"secondaryColor": "#7aa1f0",
													"shape": "round-rectangle",
													"styleOverrides": "",
													"subCategory": "Scheduling & Orchestration",
													"svgColor": "ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg",
													"svgComplete": "",
													"svgWhite": "ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg"
											},
											"position": {
													"posX": 113.0763599675155,
													"posY": -26.13721055079196
											}
									}
							},
							"type": "Pod"
					},
					"nginx-2-pod": {
							"apiVersion": "v1",
							"model": "kubernetes",
							"name": "nginx-2",
							"namespace": "default",
							"settings": {
									"spec": {
											"containers": [
													{
															"image": "nginx:stable",
															"name": "nginx",
															"ports": [
																	{
																			"container Port": 80,
																			"name": "http-web-svc"
																	}
															]
													}
											]
									}
							},
							"traits": {
									"meshmap": {
											"edges": [],
											"id": "19d2a03e-141d-41a9-bf7f-9ca12cf3b9be",
											"label": "nginx",
											"meshmodel-metadata": {
													"genealogy": "",
													"isCustomResource": false,
													"isNamespaced": true,
													"logoURL": "https://github.com/cncf/artwork/blob/master/projects/kubernetes/icon/white/kubernetes-icon-white.svg",
													"model": "kubernetes",
													"modelDisplayName": "Kubernetes",
													"primaryColor": "#326CE5",
													"published": true,
													"secondaryColor": "#7aa1f0",
													"shape": "round-rectangle",
													"styleOverrides": "",
													"subCategory": "Scheduling & Orchestration",
													"svgColor": "ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg",
													"svgComplete": "",
													"svgWhite": "ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg"
											},
											"position": {
													"posX": 113.0763599675155,
													"posY": -26.13721055079196
											}
									}
							},
							"type": "Pod"
					},
					"nginx-deployment": {
							"apiVersion": "apps/v1",
							"model": "kubernetes",
							"name": "nginx-deployment",
							"namespace": "default",
							"settings": {
									"spec": {
											"replicas": 1,
											"selector": {
													"match Labels": {
															"app": "nginx-deployment"
													}
											},
											"template": {
													"metadata": {
															"labels": {
																	"app": "nginx-deployment"
															}
													},
													"spec": {
															"containers": [
																	{
																			"image": "nginx:latest",
																			"image Pull Policy": "IfNotPresent",
																			"name": "nginx",
																			"ports": [
																					{
																							"container Port": 80,
																							"name": "nginx-port",
																							"protocol": "TCP"
																					}
																			]
																	}
															]
													}
											}
									}
							},
							"traits": {
									"meshmap": {
											"edges": [],
											"id": "87f36c76-cf0c-463e-950a-0a61948aaa41",
											"label": "nginx-deployment",
											"meshmodel-data": {
													"category": {
															"metadata": null,
															"name": "Orchestration & Management"
													},
													"displayName": "Kubernetes",
													"metadata": {
															"svgColor": "ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg",
															"svgWhite": "ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg"
													},
													"name": "kubernetes",
													"version": "v1.25.2"
											},
											"meshmodel-metadata": {
													"genealogy": "parent",
													"isCustomResource": false,
													"isNamespaced": true,
													"logoURL": "Created by Lee",
													"model": "kubernetes",
													"modelDisplayName": "Kubernetes",
													"primaryColor": "#326CE5",
													"published": true,
													"secondaryColor": "#7aa1f0",
													"shape": "round-pentagon",
													"styleOverrides": "{\"background-image\":\"none\",\"border-width\":2,\"background-opacity\":0.5}",
													"subCategory": "Scheduling & Orchestration",
													"svgColor": "ui/public/static/img/meshmodels/kubernetes/color/deployment-color.svg",
													"svgComplete": "ui/public/static/img/meshmodels/kubernetes/complete/deployment-complete.svg",
													"svgWhite": "ui/public/static/img/meshmodels/kubernetes/color/deployment-color.svg"
											},
											"position": {
													"posX": 264.94755032496664,
													"posY": 148.80052387476576
											}
									}
							},
							"type": "Deployment",
							"version": "v1.25.2"
					},
					"nginx-service-1-deployment": {
							"apiVersion": "v1",
							"model": "kubernetes",
							"name": "nginx-service",
							"namespace": "default",
							"settings": {
									"spec": {
											"ports": [
													{
															"port": 8080,
															"protocol": "TCP",
															"target Port": 80
													}
											]
									}
							},
							"traits": {
									"meshmap": {
											"edges": [],
											"id": "f69fad7e-6027-4265-8e1d-6fa5b0eee710",
											"label": "nginx-service",
											"meshmodel-data": {
													"category": {
															"metadata": null,
															"name": "Orchestration & Management"
													},
													"displayName": "Kubernetes",
													"metadata": {
															"svgColor": "ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg",
															"svgWhite": "ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg"
													},
													"name": "kubernetes",
													"version": "v1.25.2"
											},
											"meshmodel-metadata": {
													"genealogy": "",
													"isCustomResource": false,
													"isNamespaced": true,
													"logoURL": "https://github.com/cncf/artwork/blob/master/projects/kubernetes/icon/white/kubernetes-icon-white.svg",
													"model": "kubernetes",
													"modelDisplayName": "Kubernetes",
													"primaryColor": "#326CE5",
													"published": true,
													"secondaryColor": "#7aa1f0",
													"shape": "round-triangle",
													"styleOverrides": "{\"height\":16,\"width\":17,\"padding\":12,\"background-fit\":\"none\",\"background-position-y\":4.5}",
													"subCategory": "Scheduling & Orchestration",
													"svgColor": "ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg",
													"svgComplete": "",
													"svgWhite": "ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg"
											},
											"position": {
													"posX": 14.952012136066841,
													"posY": 133.77595704299856
											}
									}
							},
							"type": "Service",
							"version": "v1.25.2"
					},
					"nginx-service-2": {
							"apiVersion": "v1",
							"model": "kubernetes",
							"name": "nginx-service",
							"namespace": "default",
							"settings": {
									"spec": {
											"ports": [
													{
															"name": "name-of-service-port",
															"port": 80,
															"protocol": "TCP",
															"target Port": "http-web-svc"
													}
											],
											"selector": {
													"app.kubernetes.io/name": "proxy"
											}
									}
							},
							"traits": {
									"meshmap": {
											"edges": [],
											"id": "5b140a92-0f93-4500-b64a-ce703b49f166",
											"label": "nginx-service",
											"meshmodel-metadata": {
													"genealogy": "",
													"isCustomResource": false,
													"isNamespaced": true,
													"logoURL": "https://github.com/cncf/artwork/blob/master/projects/kubernetes/icon/white/kubernetes-icon-white.svg",
													"model": "kubernetes",
													"modelDisplayName": "Kubernetes",
													"primaryColor": "#326CE5",
													"published": true,
													"secondaryColor": "#7aa1f0",
													"shape": "round-triangle",
													"styleOverrides": "{\"height\":16,\"width\":17,\"padding\":12,\"background-fit\":\"none\",\"background-position-y\":4.5}",
													"subCategory": "Scheduling & Orchestration",
													"svgColor": "ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg",
													"svgComplete": "",
													"svgWhite": "ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg"
											},
											"position": {
													"posX": -32.24910429536563,
													"posY": 10.040920246138025
											}
									}
							},
							"type": "Service"
					}
			}
	}`),
	))

	// Check the result of the policy
	if len(rs) == 1 {
		fmt.Println("rego engine response:", rs)
	} else {
		logrus.Error("Failed to evaluate policy")
		return
	}
}

type ServiceStruct struct {
	Name       string                 `yaml:"name,omitempty" json:"name,omitempty"`
	Type       string                 `yaml:"type,omitempty" json:"type,omitempty"`
	APIVersion string                 `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	Version    string                 `yaml:"version,omitempty" json:"version,omitempty"`
	Model      string                 `yaml:"model,omitempty" json:"model,omitempty"`
	Settings   map[string]interface{} `yaml:"settings,omitempty" json:"settings,omitempty"`
	Traits     struct {
		Meshmap struct {
			Edges             []interface{}          `yaml:"edges,omitempty" json:"edges,omitempty"`
			ID                string                 `yaml:"id,omitempty" json:"id,omitempty"`
			Label             string                 `yaml:"label,omitempty" json:"label,omitempty"`
			MeshmodelData     map[string]interface{} `yaml:"meshmodel-data,omitempty" json:"meshmodel-data,omitempty"`
			MeshmodelMetadata map[string]interface{} `yaml:"meshmap,omitempty" json:"meshmap,omitempty"`
		} `yaml:"traits,omitempty" json:"traits,omitempty"`
	}
}

type MeshheryDesignFileFormat struct {
	Name    string          `yaml:"name"`
	Service []ServiceStruct `yaml:service`
}

// handlerelationship takes the design-file and resolves all the relationships supported by us and returns the updated design file with the updated relationships
func handleRelationship(designFile []byte) {

}

func handleHierarchyRelationship(designFile []byte) {

}

func handleNetworkRelationship(designFile []byte) {

}

func handleMountRelationship(designFile []byte) {

}
