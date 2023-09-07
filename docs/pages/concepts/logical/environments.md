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

Environments represent a collection of resources in the form of Connections - both of managed and unmanaged Connections. Environment resources are comprised of Connections (and implicitly those Connections' Credentials), which can be assigned and unassigned by Workspace Managers.

### Environment Relationships and Restrictions

- Environments can be assigned to zero or more Workspaces.
- Connections can be assigned to zero or more Environments.