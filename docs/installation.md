---
layout: page
title: Installation
parent: Meshery
permalink: /installation
nav_order: 2
---
## Table of contents
{: .no_toc }

1. TOC
{:toc}

---


### Quick install guide (using docker on local machine)
This approach outlines quick install leveraging `docker` and `docker-compose` on your local machine. For installing `docker` please follow these instructions [Getting Started with Docker](https://docs.docker.com/get-started/) and for `docker-compose` follow the instructions [here](https://docs.docker.com/compose/install/). 

Install Meshery on your local machine by running the following:

```
sudo curl -L https://git.io/meshery -o /usr/local/bin/meshery
```

```
sudo chmod a+x /usr/local/bin/meshery
```

```
meshery start
```
Upon starting Meshery successfully, instructions to access Meshery will be printed on the sceen.

#### Custom install guide
First clone the Meshery source code:
```
git clone https://github.com/layer5io/meshery.git; cd meshery      
```

##### Docker Compose
Meshery repository includes a `docker-compose.yaml` file. We can use `docker-compose` to spin up all the Meshery services by running:
```
docker-compose up -d
```

Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).
You will be redirected to a social login page where you can pick one of the available Social Login methods to login to Meshery.

##### Kubernetes deployment
Meshery can also be deployed on an existing Kubernetes cluster.

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
* Ingress
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



### Configuration
This is where you configure your settings on the adaptor(Istio etc) and other things 

##### Connecting grafana

##### Connecting adapters

### What is meshery CLI and how to use meshery CLI
