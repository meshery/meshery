---
layout: page
title: Installation
parent: Overview
permalink: /installation
has_children: true
nav_order: 2
---
# Installation Guide
{: .no_toc }

1. TOC
{:toc}

---
## Supported Platforms

Meshery's compatibility has been confirmed with the following platforms:
| Platform      | Version       |
| ------------- |:-------------:|   
| [Docker Engine](#docker) | 19.x and above |
| [Docker Desktop](#docker) | 2.0.x and above |
| [EKS](#eks) | 1.12.x and above |
| [GKE](#gke) | 1.14.x and above |
| [Kubernetes](#kubernetes) | 1.12.x and above |
| [Minikube](#minikube) | 1.2.x and above |

### Quick Start 
Download `mesheryctl`. Install Meshery on your local machine by running the following:

```
sudo curl -L https://git.io/meshery -o /usr/local/bin/meshery
sudo chmod a+x /usr/local/bin/meshery
meshery start
```
Upon starting Meshery successfully, instructions to access Meshery will be printed on the sceen.

### Platform: Docker <a name="docker"></a>
Follow these installation steps to use Docker and Docker Compose to run Meshery. Users often choose this installation approach in order to run Meshery on their local machine. If you need to install `docker`, see [Getting Started with Docker](https://docs.docker.com/get-started/) and if you need to install `docker-compose`, see [Installing Docker Compose](https://docs.docker.com/compose/install/). 

Meshery repository includes a `docker-compose.yaml` file. We can use `docker-compose` to spin up all the Meshery services by running:
```
sudo curl -L https://git.io/meshery -o /usr/local/bin/meshery
sudo chmod a+x /usr/local/bin/meshery
meshery start
```

Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).
You will be redirected to a social login page where you can pick one of the available Social Login methods to login to Meshery.

Upon starting Meshery successfully, instructions to access Meshery will be printed on the sceen.

### Platform: Kubernetes <a name="kubernetes"></a>
Meshery can also be deployed on an existing Kubernetes cluster running version 1.14.1 and above.

To install Meshery on your cluster, let us first create a namespace to host the Meshery components:
```
kubectl create ns meshery
```

All the needed deployment yamls for deploying Meshery are included in the `deployment_yamls/k8s` folder inside the cloned Meshery folder. To deploy the yamls on the cluster please run the following command:
```
kubectl -n meshery apply -f deployment_yamls/k8s
```
Once the yaml files are deployed, we need to expose the `meshery` service to be able to access the service from outside the cluster. 

There are several ways a service can be exposed on Kubernetes. 

Here we will describe 3 common ways we can expose a service:
**Ingress**
  * If your Kubernetes cluster has a functional Ingress Controller, then you can configure an ingress to 
    ```
    apiVersion: extensions/v1beta1
    kind: Ingress
    metadata:
    name: meshery-ingress
    annotations:
        kubernetes.io/ingress.class: "nginx"
    spec:
    rules:
    - host: *
        http:
        paths:
        - path: /
            backend:
            serviceName: meshery-service
            servicePort: 8080

    ```
* LoadBalancer
    If your Kubernetes cluster has an external load balancer, this might be a logical route.

* NodePort
    If your cluster does not have an Ingress Controller or a load balancer, then NodePort is probably the last resort.

### Platform: Minikube <a name="minikube"></a>
See the Meshery [Minikube installation guide](minikube).
## Configuration
This is where you configure your settings on the adaptor(Istio etc) and other things 

## Managed Kubernetes <a name="managedk8s"></a>
In order to run Meshery in a managed Kubernetes environment, you will need to assign an existing `ServiceAccount` or create a new `ServiceAccount`:

1. Create a `ServiceAccount` with `cluster-admin` role.
1. Get secret name from `ServiceAccount`.
1. Extract CA certificate and user token from the secret.
1. Generate new kubeconfig yaml file to use as input to Meshery.

### Platform: GKE <a name="gke"></a>
You may perform the steps outlined under [Managed Kubernetes](#managedk8s)following by hand or run the [generate_kubeconfig_gke.txt](https://github.com/layer5io/meshery/files/3166324/generate_kubeconfig_gke.txt) shell script using the desired ServiceAccount name and Namespace arguments, like so:

`./generate_kubeconfig_gke.sh cluster-admin-sa-gke default`

Having run this script, supply the generated file `config-cluster-admin-sa-gke-default.yaml` in your Meshery settings page.

### Platform: Kubernetes <a name="eks"></a>

### Connecting Grafana and Prometheus

### Connecting Meshery adapters

## What is `mesheryctl`?
`mesheryctl` is a command line interface to manage a Meshery deployment. `mesheryctl` allows you to control Meshery's lifecycle with commands like `start`, `stop`, `status`, `cleanup`. Running `cleanup` will remove all active container instanaces, prune pulled images and remove any local volumes crated by starting Meshery.