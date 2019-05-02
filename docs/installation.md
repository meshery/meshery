---
layout: page
title: Installation
parent: Meshery
permalink: /installation
nav_order: 2
---
## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---


### Quick install guide (using docker on local machine)
This approach outlines quick install leveraging docker on local machine. For installing docker please follow these instructions [Getting Started with Docker](https://docs.docker.com/get-started/)
Step 1 
### Using Docker, install Meshery on your local machine by running the following:
```
$ sudo curl -L https://git.io/meshery -o /usr/local/bin/meshery
$ sudo chmod a+x /usr/local/bin/meshery
$ meshery start
```
### Custom install guide
##### Kubernetes deployments (in-cluster and out of cluster)
Using Kubernetes, install Meshery on your cluster by cloning the Meshery repo:
```
$ git clone https://github.com/layer5io/meshery.git; cd meshery      
```
Install Meshery on your cluster by running the following:
```
$ kubectl create ns meshery
$ kubectl -n meshery apply -f deployment_yamls/k8s
```

### Running Meshery
To run Meshery point your browser to  [http://localhost:9081](http://localhost:9081)
Login using twitter or github account. 
### Configuration
This is where you configure your settings on the adaptor(Istio etc) and other things 

##### Connecting grafana

##### Connecting adapters

### What is meshery CLI and how to use meshery CLI
