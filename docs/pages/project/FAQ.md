---
layout: page
title: FAQ
permalink: project/faq
description: General commonly asked questions and answers about Meshery.
language: en
type: project
---

#### Question: What is service mesh?
**Answer:** _A service mesh is a way to control how different parts of an application share data with one another. Unlike other systems for managing this communication, a service mesh is a dedicated infrastructure layer built right into an app. This visible infrastructure layer can document how well (or not) different parts of an app interact, so it becomes easier to optimize communication and avoid downtime as an app grows._ 
#### Question: What is Meshery?
**Answer:** _Meshery is the open source, service mesh management plane that enables the adoption, operation, and management of any service mesh and their workloads._

#### Question: Why was Meshery created?

 **Answer:** _As an open source, vendor neutral project, Meshery was created out of the necessity to enable service mesh adopters to overcome the challenge of complex virtual networking; to come to expect more from their infrastructure; to enable the world to understand and operate any service mesh with confidence._

#### Question: What does Meshery do?
**Answer:** 
 - _Operates service meshes with confidence,_
   - _with operational best practices._ 
 - _Compare apples-to-apples performance across service meshes._
 - _Understand behavioral differences between service meshes._
 - _Track your application performance from version to version._

#### Question: What is mesheryctl?

**Answer:** _A command line interface to manage Meshery._

#### Question: How to install Meshery?

**Answer:** _Meshery runs on a number of platforms._
_Popular Installers:_
- _Bash user → Meshery Installation [Quick Start](https://docs.meshery.io/installation/quick-start)._
- _[Brew user](https://github.com/meshery/homebrew-tap)._
- _Kubernetes user._
- _[Scoop user](https://github.com/meshery/scoop-bucket)._ 


#### Question: What is the difference between `make server` and `mesheryctl system start`? Do they both run Meshery on my local machine?

  **Answer:** _Yes, both of them do run Meshery on your local machine. `make server` builds Meshery from source and runs it on your local OS, while `mesheryctl system start` runs Meshery as a set of containers in Docker or in Kubernetes on your local machine._ 

#### Question: What systems is Meshery compatible with?

  **Answer:** _Many. See Meshery's [Compatibility Matrix]({{site.baseurl}}/project/compatibility-matrix)._

#### Question: Why is Meshery Server only receiving MeshSync updates from one of my Kubernetes Clusters?

  **Answer:** _In order to receive MeshSync updates, Meshery Server subscribes for updates Meshery Broker. In other words, Meshery Server connects to the `meshery-broker` service port in order to subscribe for streaming MeshSync updates. By default, the Meshery Broker service is deployed as type Kubernetes Service type `LoadBalancer`, which requires that your Kubernetes cluster provides an external IP address to the Meshery Broker service, exposing it external to the Kubernetes cluster. _
_If you're running Kubernetes in Docker Desktop, an external IP address of `localhost` is assigned. If you're running Minikube, and execute `minikube tunnel` to gain access to Meshery Broker's service, you will find that both Meshery Broker service endpoints (from two different clusters) are sharing the same `localhost:4222` address and port number. This port sharing causes conflict and Meshery Server is only able to connect to one of the Meshery Brokers._
  
  _Few ways to solve this problem:_
  - _Use an external cloud provider which provides you with the LoadBalancer having an external IP address other than localhost_
  - _Use [Kind](https://kind.sigs.k8s.io) cluster with [MetalLB](https://metallb.universe.tf) configuration_

#### Question: Why Service Mesh deployed is not visible in UI despite successful installation in my cluster?
  **Answer:** _Few ways to solve this problem:_
  - Confirm that the Meshery Broker service is exposed from your cluster using `kubectl get svc -n meshery`
  - It is possible that MeshSync is not healthy and not sending cluster updates, check for MeshSync status by going in Settings from Meshery UI and click on MeshSync chip.
  - If MeshSync is healthy, next step is to check for status of Meshery Broker by clicking on the NATS chip. 
  
  If either is the case, Meshery Operator will make sure MeshSync and Meshery Broker deployments are again healthy, wait for some time, otherwise try redeploying Meshery Operator.
  
  _You can also verify health of you system using [mesheryctl system check]({{site.baseurl}}/reference/mesheryctl/system/check)_
  
  If MeshSync, Meshery Broker and Meshery Operator are healthy.

  - Try clearing the database by clicking on the `Flush MeshSync` button associated with the corresponding cluster.
  - If still `Service Mesh` is not visible in UI, move on to `Hard Reset` of Database. This option is in the `Reset System` Tab in `Settings` page.

{% include discuss.html %}

<!--Add other questions-->

