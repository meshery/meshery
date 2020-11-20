# Meshery

Meshery is a the multi-service mesh management plane offering lifecycle, configuration and performance management of service meshes and applications running atop them.

## Introduction

The chart bootstraps a single nodes Meshery deployment on Kubernetes cluster using the Helm package manager.

## Prerequisites

* Kubernetes >= v1.15.0
* K3s v1.14.4-k3s.1
* Helm v3.x.x

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
NAME                                  READY   STATUS    RESTARTS   AGE
pod/meshery-78659f4d65-qfn5q          1/1     Running   0          7m57s
pod/meshery-consul-5578fb5d54-tlgb9   1/1     Running   0          7m3s
pod/meshery-istio-5f9cf5f678-mchnc    1/1     Running   0          7m3s
pod/meshery-kuma-65bc58d67c-t7l52     1/1     Running   0          7m57s
pod/meshery-linkerd-db46fbc88-ccnpb   1/1     Running   0          7m57s
pod/meshery-nsm-d4dbf8f9b-hqc84       1/1     Running   0          7m57s
pod/meshery-osm-654f884b76-ffx84      1/1     Running   0          7m57s

NAME                      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
service/meshery           ClusterIP   10.97.119.203    <none>        8080/TCP    7m57s
service/meshery-consul    ClusterIP   10.97.141.111    <none>        10002/TCP   7m3s
service/meshery-istio     ClusterIP   10.106.102.178   <none>        10000/TCP   7m3s
service/meshery-kuma      ClusterIP   10.110.169.103   <none>        10007/TCP   7m57s
service/meshery-linkerd   ClusterIP   10.104.36.196    <none>        10001/TCP   7m57s
service/meshery-nsm       ClusterIP   10.97.129.64     <none>        10004/TCP   7m57s
service/meshery-osm       ClusterIP   10.103.243.140   <none>        10009/TCP   7m57s

NAME                              READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/meshery           1/1     1            1           7m57s
deployment.apps/meshery-consul    1/1     1            1           7m3s
deployment.apps/meshery-istio     1/1     1            1           7m3s
deployment.apps/meshery-kuma      1/1     1            1           7m57s
deployment.apps/meshery-linkerd   1/1     1            1           7m57s
deployment.apps/meshery-nsm       1/1     1            1           7m57s
deployment.apps/meshery-osm       1/1     1            1           7m57s

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/meshery-78659f4d65          1         1         1       7m57s
replicaset.apps/meshery-consul-5578fb5d54   1         1         1       7m3s
replicaset.apps/meshery-istio-5f9cf5f678    1         1         1       7m3s
replicaset.apps/meshery-kuma-65bc58d67c     1         1         1       7m57s
replicaset.apps/meshery-linkerd-db46fbc88   1         1         1       7m57s
replicaset.apps/meshery-nsm-d4dbf8f9b       1         1         1       7m57s
replicaset.apps/meshery-osm-654f884b76      1         1         1       7m57s
```

```
# Check the Helm chart release
$ helm ls -n meshery
```

```
NAME   	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART        	APP VERSION
meshery	meshery  	1       	2020-01-22 11:53:20.407751 +0800 CST	deployed	meshery-2.0.0	stable-latest
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