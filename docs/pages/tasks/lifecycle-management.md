---
layout: default
title: Lifecycle Management
permalink: tasks/lifecycle-management
type: tasks
language: en
list: include
---

<a name="lifecycle-management"></a>

Meshery manages hundreds of different types of cloud native infrastructure. See [full the set of integrations](https://meshery.io/integrations).

{% include alert.html type="info" title="Meshery's Capabilities Registry" content="See <a href='/concepts/architecture'>Architecture</a> for a description of Meshery's Capabilities Registry for an understanding of how Meshery Adapters deliver infrastructure-specific functionality." %}

## Cloud Native Infrastructure Lifecycle Management

![Lifecycle Management]({{ site.baseurl }}/assets/img/tasks/lifecycle-management.svg)

Meshery Adapters can optionally be deployed to provide deep support of different types of infrastructure including [service mesh]({{ site.baseurl }}/service-meshes), so that you can learn and compare different infrastructure functionality. Once you have selected to operate a given infrastructure, Meshery will assist you with ongoing operations.

Meshery automates the provisioning of various infrastructurees, allowing users to select different configuration profiles in order to support a variety of deployment models.

Meshery adapters will dynamically retrieve the specific infrastructure's release package upon initial deployment of a given infrastructure. Meshery adapters cache (in `.meshery/bin`) the infrastructure installation package.

#### Deprovisioning infrastructure

Just as Meshery automates the provisioning of various cloud native infrastructure, so too does it facilitate the deprovisioning of infrastructure. Find the "Undeploy" button in Meshery UI or the `--undeploy` flag in Meshery CLI.

## Workload Lifecycle Management

Users may bring their applications (perform workload onboarding) on the infrastructure using the Custom Configuration operation.

Commonly adjoining each infrastructure is a sample application that facilitates demonstration of the value of the given infrastructure. Meshery allows you to efficiently [install sample applications]({{ site.baseurl }}/guides/sample-apps) across different infrastructurees.

## Configuration Management

Meshery provides the ability for you as a infrastructure manager to customize your infrastructure deployment.

## Notifications

Meshery tracks operations that you perform on infrastructurees and their workloads. Meshery provides notification of environment issues, application conflicts with infrastructure configuration, and so on.

#### Acknowledging and dismissing notifications

Depending upon the severity of the notification, many notifications are informational and can readily be dismissed after reading. Some notifications include actionable alerts. Learn more in "[Managing Events with Notification Center](managing-notifications)".
