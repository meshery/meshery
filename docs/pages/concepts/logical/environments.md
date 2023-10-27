---
layout: default
title: Environments
permalink: concepts/environments
type: concepts
abstract: "Meshery offers support for Kubernetes cluster and cloud state synchronization with the help of MeshSync."
language: en
list: include
---


## Environments

Environments, in the context of Meshery, serve as logical groupings of connections and credentials. Connections are the crucial aspect of your cloud-native applications, ranging from Kubernetes clusters to Prometheus instances, Jaeger services, Nginx deployments, and more. These connections could be either managed or unmanaged, depending on the tiers of discovery performed by [MeshSync](concepts/architecture/meshsync), Meshery's dynamic, and sophisticated resource discovery component. The ability to organize these resources into logical groups brings a new level of convenience and efficiency to managing your cloud-native infrastructure.


## Key Features

- **Logical Grouping** Environments allow you to logically group related connections and their associated credentials. This makes it easier to manage, share, and work with a subset of resources instead of dealing with all your connections individually.

- **Resource Sharing** Environments can be seamlessly assigned to [Workspaces](/concepts/workspaces), another essential concept in Meshery. When you assign an Environment to a Workspace, you enable resource sharing among team members. This collaborative approach simplifies the sharing of connections and resources, making it easier to work together in cloud-native environments.


## Key Components

### Connections
Connections are an integral part of Environment. These are cloud-native resources that can be both managed and unmanaged, and they're registered by the Meshery Server. Examples of connections include Kubernetes clusters, Prometheus instances, Jaeger tracers, and Nginx web servers.
See "[Connections](/concepts/connections)" section for more information.

### Credentials
Credentials in an Environment are the keys to securely authenticate and access managed connections. For example, valid Prometheus secrets or Kubernetes API tokens are essential credentials for securely interacting with these managed resources.
See "[Credentials](/concepts/credentials)" section for more information.

