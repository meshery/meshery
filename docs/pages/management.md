---
layout: page
title: Mesh Management
permalink: mesh-management
---

# Service mesh lifecycle management
Meshery supports a number of [different service meshes](installation/adapters), so that you can learn and compare different service mesh functionality. Once you have selected to operate a given service mesh, Meshery will assist you with ongoing activities.

## Deploying a service mesh
Meshery automates the provisioning of various service meshes, allowing users to select different configuration profiles in order to support a variety of deployment models.

## Installing sample applications
Commonly adjoining each service mesh is a sample application that facilitates demonstration of the value of the given service mesh. Meshery allows you to quickly deploy the same sample application across different service meshes.

## Use custom service mesh configuration
Meshery provides the ability for you as a service mesh manager to customize your service mesh deployment.

## Deleting a service mesh
Just as Meshery automates the provisioning of various service meshes, so does it facilitate the deprovisioning of service meshes.

# Management notifications 
Meshery tracks service mesh and application health. Meshery provides notification of environment issues, application conflicts with service mesh configuration.

**Acknowledging and dismissing notifications**
Many notifications are informational and can readily be dismissed after reading. Some notifications include actionable alerts.

<strong>Taking action on alerts</strong>

Particular notifications that Meshery presents are immediately actionable. These are denoted by a red colored highlight. Be on the lookout for these alerts and take action promptly to ensure a smooth running service mesh.

# Node and Service Mesh Metrics

## Grafana and Meshery

Connect Meshery to your existing Grafana instance and Meshery will import the boards of your choosing. 

### Connecting to Grafana
If you have an API key configured to restrict access to your Grafana boards, you will need to enter the API key when establishing Meshery's connection to Grafana.

* Importing Grafana boards
    - Importing existing Grafana boards via API
    - Importing custom Grafana board via yaml
* Configuring graph panel preferences

## Prometheus and Meshery
Meshery allows users to connect to one or more Prometheus instances in order to gather telemetric data (in the form of metrics). These metrics may pertain to service meshes, Kubernetes, applications on the mesh or really... any metric that Prometheus has collected.

Once you have connected Meshery to your Prometheus deployment(s), you may perform ad-hoc connectivity tests to verify communication between Meshery and Prometheus.
