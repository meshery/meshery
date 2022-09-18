---
layout: default
title: Application Management
permalink: functionality/application-management
type: functionality
language: en
list: include
---

Meshery's application management features allows you to store and onboard/offboard applications onto service meshes. Meshery also comes with some sample applications out-of-the-box for you to quickly try out service mesh features.

## What is a Meshery Application?

A single file (future: which could be split into multiple files upon user request) that comprises Kubernetes objects representative of a complete set of Kubernetes workload resources. 

### What is a Meshery Rollout?
	"v1" Application = Services + Deployments + ReplicaSets
	"v2" Application = Services + Deployments + ReplicaSets (Rollout Strategy) 


You can define the applications as _Kubernetes manifest files\*_ as well as [pattern files](./patterns.md).

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-applications-seeded.png" />

You can also import application files from the filesystem, a URL or from a GitHub repository.

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-applications.png" />

See [Guide: Application Management]() for a more detailed guide on working with your applications in Meshery.

**Known Caveats:**
- Rollouts do not support custom ServiceAccount names. User's ServiceAccount name must be the same as the Rollouts name.
