---
layout: page
title: Frequently Asked Questions
permalink: project/faq
description: General commonly asked questions and answers about Meshery.
language: en
type: project
---

## General FAQs

#### Question: What is Meshery?

**Answer:** _As a self-service engineering platform, Meshery enables collaborative design and operation of cloud native infrastructure._

#### Question: Why was Meshery created?

**Answer:** _As an open source, vendor neutral project, Meshery was created out of the necessity to enable platform engineers, site reliability engineers, devops engineers... engineers to collaborate in the management of their infrastucture and workloads. Meshery was created to enable you to expect more from your infrastructure and to do so with confidence._

#### Question: What does Meshery do?

**Answer:** _Infrastructure as design. Meshery enables you to design and operate cloud native infrastructure visually, collaboratively, with confidence and in partnership with your teammates. See all of Meshery's features [here]({{site.baseurl}}/features)._

<!-- - _offers a catalog of operational best practices._
- _offersompare apples-to-apples performance across different infrastructure configurations._
- _Understand behavioral differences between service deployments._
- _Track your application performance from version to version._ -->

#### Is Meshery open source project?

**Answer:** _Yes, Meshery is a Cloud Native Computing Foundation (CNCF) project and is licensed under Apache v2. As the cloud native management plane, Meshery is an extensible platform, offering multiple extension points within which users and partners can customize and extend Meshery's functionality._

#### Question: Why should I use Meshery?

**Answer:** _Meshery is a powerful tool for managing ​Kubernetes infrastructure. It seamlessly integrates with different hundreds of tools and offers extensibility through many different [extension points]({{site.baseurl}}/extensibility/#extension-points). With Meshery, you can easily discover your environment, collaboratively manage multiple Kubernetes clusters, connect your Git and Helm repos, and analyze app and infra performance._

## User FAQs

#### Question: What is mesheryctl?

**Answer:** _A command line interface to manage Meshery._

#### Question: How do I install Meshery?

**Answer:** _Meshery runs on a [number of platforms]({{site.baseurl}}/installation)._
_You are encouraged to use `mesheryctl` to configure and control Meshery deployments. Install `mesheryctl` using any of these options:_

- _[Bash user](/installation/linux-mac/bash)_
- _[Brew user](/installation/linux-mac/brew)_
- _[Scoop user](/installation/windows/scoop)_

#### Question: What architecture does Meshery have?

**Answer:** _.There are several compotents, languages and they have different purposes. See Meshery's [Architecture](https://docs.meshery.io/concepts/architecture)._

#### Question: What is the difference between `make server` and `mesheryctl system start`? Do they both run Meshery on my local machine?

**Answer:** _Yes, both of them do run Meshery on your local machine. `make server` builds Meshery from source and runs it on your local OS, while `mesheryctl system start` runs Meshery as a set of containers in Docker or in Kubernetes on your local machine._

#### Question: What systems can I deploy Meshery onto?

**Answer:** _Many. See Meshery's [Compatibility Matrix]({{site.baseurl}}/installation)._

#### Question: What systems does Meshery manage?

**Answer:** _Many. See Meshery's [Integrations](https://meshery.io/integrations)._

#### Question: Why is Meshery Server only receiving MeshSync updates from one of my Kubernetes Clusters?

**Answer:** _In order to receive MeshSync updates, Meshery Server subscribes for updates Meshery Broker. In other words, Meshery Server connects to the `meshery-broker` service port in order to subscribe for streaming MeshSync updates. By default, the Meshery Broker service is deployed as type Kubernetes Service type `LoadBalancer`, which requires that your Kubernetes cluster provides an external IP address to the Meshery Broker service, exposing it external to the Kubernetes cluster. _
_If you're running Kubernetes in Docker Desktop, an external IP address of `localhost` is assigned. If you're running Minikube, and execute `minikube tunnel` to gain access to Meshery Broker's service, you will find that both Meshery Broker service endpoints (from two different clusters) are sharing the same `localhost:4222` address and port number. This port sharing causes conflict and Meshery Server is only able to connect to one of the Meshery Brokers._

_Few ways to solve this problem:_

- _Use an external cloud provider which provides you with the LoadBalancer having an external IP address other than localhost_
- _Use [Kind](https://kind.sigs.k8s.io) cluster with [MetalLB](https://metallb.universe.tf) configuration_

#### Question: Why does the dashboard not show the infrastructure provisioned or discovered by Meshery?

**Answer:** _This issue is typically caused by either lack of connectivity between Meshery Server and Meshery Broker or by database corruption. Use the following troubleshooting steps to resolve this issue:_

**Lack of Connectivity**

1. Confirm that the Meshery Broker service is exposed from your cluster using `kubectl get svc -n meshery` and that an hostname or IP address is displayed in the External Address column. Meshery Server should be able to reach this address.
1. It is possible that MeshSync is not healthy and not sending cluster updates, check for MeshSync status by navigating to Settings in Meshery UI and clicking on the MeshSync connection.
1. If MeshSync is healthy, check the status of Meshery Broker by clicking on the NATS connection.

If either is the case, Meshery Operator will make sure MeshSync and Meshery Broker deployments are again healthy, wait for some time, otherwise try redeploying Meshery Operator.

**Database Corruption**
If MeshSync, Meshery Broker and Meshery Operator are healthy, then perhaps, there is corruption in the Meshery Database. Use the following troubleshooting steps to resolve this issue:\_

    1. Try clearing the database by clicking on the `Flush MeshSync` button associated with the corresponding cluster.
    1. If still `Service Mesh` is not visible in UI, move on to `Hard Reset` of Database. This option is in the `Reset System` Tab in `Settings` page.

Note: _You can also verify health of your system using [mesheryctl system check]({{site.baseurl}}/reference/mesheryctl/system/check)_

## Contributing FAQs

#### Question: Getting error while running `make server`?

**Answer:** _On Windows, set up the project on Ubuntu WSL2 and you will be able to run the Meshery UI and the server. For more information please visit [Supported Platforms](https://docs.meshery.io/installation)._

{% include discuss.html %}

<!--Add other questions-->

