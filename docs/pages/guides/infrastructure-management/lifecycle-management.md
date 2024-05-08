---
layout: default
title: Managing Connections
permalink: guides/infrastructure-management/lifecycle-management
redirect_from: 
- tasks/lifecycle-management
- guides/infrastructure-management/managing-connections
type: guides
category: infrastructure
language: en
list: include
abstract: 'Manage the lifecycle of your infrastructure by registering each infrastructure element with Meshery.'
---

<a name="lifecycle-management"></a>

Meshery manages hundreds of different types of cloud native infrastructure. See the [full set of integrations](/integrations).

## Cloud Native Infrastructure Lifecycle Management

[![Lifecycle Management]({{ site.baseurl }}/assets/img/lifecycle-management/states-for-kubernetes-cluster-connections.svg)](/assets/img/lifecycle-management/states-for-kubernetes-cluster-connections.svg)

Meshery Adapters can optionally be deployed to provide deep support of different types of infrastructure, so that you can learn and compare different infrastructure functionality. Once you have selected to operate a given infrastructure, Meshery will assist you with ongoing operations.

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

Depending upon the severity of the notification, many notifications are informational and can readily be dismissed after reading. Some notifications include actionable alerts. Learn more in "[Managing Events with Notification Center]({{site.baseurl}}/guides/events-management)".

## Lifecycle FAQs

<details>
<summary>
<strong>Question:</strong> “Meshery is a collaboration platform. When I bring my infrastructure under Meshery's management, will my Kubernetes clusters be available to all other users? Can other people access my cluster?
</summary><strong>Answer:</strong> <p>Yes, they can, <i>IF</i> you explicitly allow them to do so. It's important to understand the following controls and system behavior:</p>
<p><b>1. Ownership:</b> Every connection to a Kubernetes cluster is created by and owned by the individual that provided the Kubernetes context. That individual may elect to share the connection with others on their team (if you have invited anyone to your team(s)).</p>
<p><b>2. Permission:</b> If you do grant other team members access, you do so by creating an environment, assigning that Kubernetes connection to the environment, then creating a workspace and assigning that environment to the workspace. Users of any of your teams to which you have shared access to the workspace will then have permission to access the cluster.</p>
<p><b>3. Connectivity:</b> those individuals will have to be afforded network connectivity to that cluster (in whatever fashion you deem appropriate). Grossly, there are two ways in which this can occur:</p>
<p>3.a) The other user runs their own copy of Meshery (or signs into a shared instance like the Playground) and will see the connection as being available when they sign in. Their Meshery Server will need to be able to reach your Kube API over the network. How that is done can be any number of ways and is left unto your own devices.</p>
<p>3.b.) The other user signs into your Meshery Server instance, which has network access to your Kubernetes cluster. In order for the other user to sign into your Meshery Server, you would have to expose it to the Internet or VPN or… one of the many other ways to all them access to your Meshery Server.</p>
<p>So, in short, yes, you can share access to your Kubernetes cluster with other users, but you have to explicitly grant them access to do so. For more information please visit <a href="/extensibility/authorization">Remote Provider Permissions</a>.</p>
</details>

{% include discuss.html %}