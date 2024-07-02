---
layout: default
title: Deploying PHP Guestbook application(Stateless) with Redis
abstract: |
  In this tutorial, you'll build and deploy a simple multi-tier web application using Kubernetes and Docker. It includes a Redis database for storing guestbook entries and multiple web frontend instances. You'll learn how to containerize these components, set up Kubernetes deployments and services, and establish communication between them. This hands-on guide is perfect for understanding basic application deployment in a clustered environment.
permalink: guides/tutorials/deploying-guestbook-application-with-redis
redirect_from: guides/tutorials/Deploying-Guestbook-Application-With-Redis
model: kubernetes
type: guides
category: tutorials
language: en
list: include
---
### Introduction

This tutorial shows you how to build and deploy a simple (not production ready), multi-tier web application using Kubernetes and Docker. This example consists of the following components:
- A single-instance Redis to store guestbook entries
- Multiple web frontend instance

### Prerequisites:
- Basic understanding of containerization , Kubernetes concepts and basic knowledge of reddis .
- Access to the _Me
- shery Playground_. If you don't have an account, sign up at [Meshery Playground](https://meshery.layer5.io/).

### lab scenario
In this lab, you will set up a Redis leader followed by two Redis followers to create a replicated Redis cluster. Next, you will deploy the guestbook frontend application that interacts with this Redis cluster. You will expose the frontend service to view and interact with the guestbook. Finally, after completing your tasks and testing, you will clean up by deleting the deployed resources to ensure no lingering components remain in your Kubernetes environment.

### Learning Objectives
In this tutorial, you will learn to deploy a Redis cluster with a leader and followers, understand Docker containerization for a guestbook frontend, configure Kubernetes deployments and services, expose services for external access, and practice resource cleanup with the help of meshery playground . These skills will provide a foundational understanding of managing distributed applications with Kubernetes and Docker in a clustered environment.

### Steps:

#### Access Meshery Playground
   - Log in to the [Meshery Playground](https://meshery.layer5.io/) using your credentials. On successful login, you should be at the dashboard. Press the **X** on the _Where do you want to start?_ popup to close it (if required).
   - Click **Explore** in the Cloud Native Playground tile to navigate to _MeshMap_.

> **_NOTE:_**  MeshMap is still in beta.

#### Step 1- Start up the Redis Database

The guestbook application uses Redis to store its data.

### Creating the Redis Deployment
The manifest file, included below, specifies a Deployment controller that runs a single replica Redis Pod.
```\yaml
# SOURCE: https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-leader
  labels:
    app: redis
    role: leader
    tier: backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
        role: leader
        tier: backend
    spec:
      containers:
      - name: leader
        image: "docker.io/redis:6.0.5"
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
        ports:
        - containerPort: 6379
```
#### Step2- Creating the Redis leader Service

The guestbook application needs to communicate to the Redis to write its data. You need to apply a Service to proxy the traffic to the Redis Pod. A Service defines a policy to access the Pods.
```yaml\
# SOURCE: https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook
apiVersion: v1
kind: Service
metadata:
  name: redis-leader
  labels:
    app: redis
    role: leader
    tier: backend
spec:
  ports:
  - port: 6379
    targetPort: 6379
  selector:
    app: redis
    role: leader
    tier: backend
```
## Note 
This manifest file creates a Service named redis-leader with a set of labels that match the labels previously defined, so the Service routes network traffic to the Redis Pod.

#### Step3 -Set up Redis followers 
Although the Redis leader is a single Pod, you can make it highly available and meet traffic demands by adding a few Redis followers, or replicas.
```\yaml
# SOURCE: https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-follower
  labels:
    app: redis
    role: follower
    tier: backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
        role: follower
        tier: backend
    spec:
      containers:
      - name: follower
        image: us-docker.pkg.dev/google-samples/containers/gke/gb-redis-follower:v2
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
        ports:
        - containerPort: 6379
```
#### Step4 - Creating the Redis follower service
The guestbook application needs to communicate with the Redis followers to read data. To make the Redis followers discoverable, you must set up another Service.
```\yaml
# SOURCE: https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook
apiVersion: v1
kind: Service
metadata:
  name: redis-follower
  labels:
    app: redis
    role: follower
    tier: backend
spec:
  ports:
    # the port that this service should serve on
  - port: 6379
  selector:
    app: redis
    role: follower
    tier: backend
```
#### Step5 -Set up and Expose the Guestbook Frontend
Now that you have the Redis storage of your guestbook up and running, start the guestbook web servers. Like the Redis followers, the frontend is deployed using a Kubernetes Deployment.

The guestbook app uses a PHP frontend. It is configured to communicate with either the Redis follower or leader Services, depending on whether the request is a read or a write. The frontend exposes a JSON interface, and serves a jQuery-Ajax-based UX.
### Creating the Guestbook Frontend Deployment
Here is deployment manifest:
```\yaml
# SOURCE: https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
spec:
  replicas: 3
  selector:
    matchLabels:
        app: guestbook
        tier: frontend
  template:
    metadata:
      labels:
        app: guestbook
        tier: frontend
    spec:
      containers:
      - name: php-redis
        image: us-docker.pkg.dev/google-samples/containers/gke/gb-frontend:v5
        env:
        - name: GET_HOSTS_FROM
          value: "dns"
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
        ports:
        - containerPort: 80
```
### Creating the Frontend Service
The Redis Services you applied is only accessible within the Kubernetes cluster because the default type for a Service is ClusterIP. ClusterIP provides a single IP address for the set of Pods the Service is pointing to. This IP address is accessible only within the cluster.

If you want guests to be able to access your guestbook, you must configure the frontend Service to be externally visible, so a client can request the Service from outside the Kubernetes cluster. However a Kubernetes user can use kubectl port-forward to access the service even though it uses a ClusterIP

## Note
Some cloud providers, like Google Compute Engine or Google Kubernetes Engine, support external load balancers. If your cloud provider supports load balancers and you want to use it, uncomment type: LoadBalancer.
Here Frontend Service manifest:
```\yaml
# SOURCE: https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook
apiVersion: v1
kind: Service
metadata:
  name: frontend
  labels:
    app: guestbook
    tier: frontend
spec:
  # if your cluster supports it, uncomment the following to automatically create
  # an external load-balanced IP for the frontend service.
  # type: LoadBalancer
  #type: LoadBalancer
  ports:
    # the port that this service should serve on
  - port: 80
  selector:
    app: guestbook
    tier: frontend
```
#### Step6 - Vizualizing  the Cassandra StatefulSet

To view the resources deployed we will use the **Visualize** section of the _MeshMap_. A view is created with necessary filters to show the relevant resources.

1.  Click **Visualize** to begin.
2.  Give the view a name (rename).
3.  Click the filter icon.
4.  Choose appropriate filters to limited displayed resources in the view. For example, here we want to display ,Service,deployments etc  
    Additionally, we will also add a label filter i.e. `tutorial=deployment` in this case. This should show a filtered view with only your resources something 
    similar to the screenshot below:

### Viewing deployment and Service information

Select the Statefullset from the _View_ to load the Statefullset details to the right. Ensure the _Details_ tab is selected.

Now, select one of the Service to display the service details. 



#### Step7 -Cleaning up
Click on `undeploy` on console to clean up resources
