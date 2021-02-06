---
layout: default
title: Lifecycle Management
permalink: functionality/lifecycle-management
type: functionality
language: en
list: include
---
<a name="lifecycle-management"></a>

# Service Mesh Lifecycle Management

Meshery supports a number of [different service meshes]({{ site.baseurl }}/service-meshes), so that you can learn and compare different service mesh functionality. Once you have selected to operate a given service mesh, Meshery will assist you with ongoing activities.

## Deploying a service mesh

Meshery automates the provisioning of various service meshes, allowing users to select different configuration profiles in order to support a variety of deployment models.

Meshery adapters will dynamically retrieve the specific service mesh's release package upon initial deployment of a given service mesh. Meshery adapters cache (in `.meshery/bin`) the service mesh installation package.

#### [Installing sample applications]({{ site.baseurl }}/guides/sample-apps)

Commonly adjoining each service mesh is a sample application that facilitates demonstration of the value of the given service mesh. Meshery allows you to efficiently deploy the same sample application across different service meshes.

#### Use custom service mesh configuration

Meshery provides the ability for you as a service mesh manager to customize your service mesh deployment.

#### Deleting a service mesh

Just as Meshery automates the provisioning of various service meshes, so does it facilitate the deprovisioning of service meshes.

## Notifications 

Meshery tracks service mesh and application health. Meshery provides notification of environment issues, application conflicts with service mesh configuration.

### Acknowledging and dismissing notifications

Many notifications are informational and can readily be dismissed after reading. Some notifications include actionable alerts.

<strong>Taking action on alerts</strong>

Particular notifications that Meshery presents are immediately actionable. These are denoted by a red colored highlight. Be on the lookout for these alerts and take action promptly to ensure a smooth running service mesh.
