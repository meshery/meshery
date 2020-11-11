# Meshery

Meshery is a the multi-service mesh management plane offering lifecycle, configuration and performance management of service meshes and applications running atop them.

## Introduction

The chart bootstraps a single nodes Meshery deployment on Kubernetes cluster using the Helm package manager.

## Prerequisites

* Kubernetes v1.15.0
* K3s v1.14.4-k3s.1
* Helm v3.0.2

## Installing the Chart

> If you in Helm2 please kindly use `helm template`, the chart only was passed tested under Helm3.

```
$ kubectl create namespace meshery
$ helm repo add meshery https://meshery.io/charts/
$ helm install meshery --namespace meshery meshery/meshery
```


## Check the resource

```
# Check the resource
$ kubectl get pods -n meshery
```

```
NAME                                   READY   STATUS    RESTARTS   AGE
pod/meshery-664668d77c-ztsgp           1/1     Running   0          98m
pod/meshery-istio-5b77685db9-k2mz6     1/1     Running   0          98m
pod/meshery-linkerd-569877fcff-ljtfs   1/1     Running   0          98m

NAME                      TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
service/meshery           ClusterIP   10.43.84.165    <none>        8080/TCP    98m
service/meshery-istio     ClusterIP   10.43.3.85      <none>        10000/TCP   98m
service/meshery-linkerd   ClusterIP   10.43.196.211   <none>        10001/TCP   98m

NAME                              READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/meshery           1/1     1            1           98m
deployment.apps/meshery-istio     1/1     1            1           98m
deployment.apps/meshery-linkerd   1/1     1            1           98m

NAME                                         DESIRED   CURRENT   READY   AGE
replicaset.apps/meshery-664668d77c           1         1         1       98m
replicaset.apps/meshery-istio-5b77685db9     1         1         1       98m
replicaset.apps/meshery-linkerd-569877fcff   1         1         1       98m
```

```
# Check the Helm chart release
$ helm ls -n meshery
```

```
NAME   	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART        	APP VERSION
meshery	meshery  	1       	2020-01-22 11:53:20.407751 +0800 CST	deployed	meshery-1.0.0	latest
```



## Uninstall the Chart

```
$ helm uninstall meshery -n meshery
release "meshery" uninstalled

$ kubectl delete namespace meshery
namespace "meshery" deleted
```

# License
This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).

# About Layer5

[Layer5.io](https://layer5.io/) is the service mesh community, serving as a repository for information pertaining to the surrounding technology ecosystem (service meshes, api gateways, edge proxies, ingress and egress controllers) of microservice management in cloud native environments.