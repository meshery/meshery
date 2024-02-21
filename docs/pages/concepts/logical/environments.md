---
layout: enhanced
title: Environments
permalink: concepts/logical/environments
type: concepts
abstract: "Environments are how you organize your deployment targets (whether on-premises servers or cloud services) into resource groups."
language: en
list: include
redirect_from:
  - concepts/environments
---

Meshery Environments allow you to logically group related [Connections](#connections) and their associated [Credentials](#credentials). Environments make it easier for you to manage, share, and work with a collection of resources as a group, instead of dealing with all your Connections and Credentials on an individual basis.

### Assigning Resources to an Environment

Assign any number of Connections to an environment whether that Connection is managed or unmanaged (see [MeshSync](/concepts/architecture/meshsync) to learn more about managed and unmanaged Connections). In-turn, assign any number of Environments to one or more [Workspaces](/concepts/logical/workspaces). Connections (and any associated Credentials) that are assigned to an Environment become immediately available for use in any associated Workspace.

### Sharing Resources between Environments

Environments can share resources. For example, you might create an environment named "production" and assign three connections: a GitHub connection, a Kubernetes connection, and a Prometheus connection. Subsequently, you also define a an environment named "dev/test "and assign three connections: a different Kubernetes connection, a different Prometheus connection, _and_ the same GitHub connection that is also assigned to the "production" environment.

### Deleting an Environment

Deleting an environment does not delete any resources (e.g. connections) currently contained with the environment. Resources that belong to others environments will continue to belong to those other environments. Learn more about the behavior of [lifecycle of connections]({{site.baseurl}}/concepts/logical/connections).

## Summary

Environments represent a collection of resources in the form of Connections - both of managed and unmanaged Connections. Environment resources are comprised of Connections (and implicitly any Credentials used by those assigned Connections). Create and use environments to organize your connections and credentials into groups, and then make these resources available to you and your teams by assigning environments to Workspaces.

## Key Features

- **Logical Grouping** Environments allow you to logically group related connections and their associated credentials. This makes it easier to manage, share, and work with a subset of resources instead of dealing with all your connections individually.

- **Resource Sharing** Environments can be seamlessly assigned to [Workspaces](/concepts/logical/workspaces), another essential concept in Meshery. When you assign an Environment to a Workspace, you enable resource sharing among team members. This collaborative approach simplifies the sharing of connections and resources, making it easier to work together in cloud-native environments.

## Key Components

### Connections <a id="connections"></a>

Connections are an integral part of Environment. These are cloud-native resources that can be both managed and unmanaged, and they're registered by the Meshery Server. Examples of connections include Kubernetes clusters, Prometheus instances, Jaeger tracers, and Nginx web servers.

See "[Connections](/concepts/logical/connections)" section for more information.

### Credentials <a id="credentials"></a>

Credentials in an Environment are the keys to securely authenticate and access managed connections. For example, valid Prometheus secrets or Kubernetes API tokens are essential credentials for securely interacting with these managed resources.

See "[Credentials](/concepts/logical/credentials)" section for more information.

