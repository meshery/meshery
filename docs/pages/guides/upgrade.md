---
layout: default
title: Upgrade `mesheryctl` and Meshery
description: How to Meshery and all of its components
permalink: guides/upgrade
display-title: "false"
type: Guides
---
# Upgrading Meshery

| Components  | Sub-component | Considering or Updating | 	 
| ----------- | ----------- | ----------- |
| `Meshery Adapters` | - `Istio Adapter`      | `Docker Deployment`: Watchtower updates this component in accordance with the user’s release channel subscription.|
|                  | - `Linkerd Adapter`    |
| `Meshery Server` |- `Meshery UI` | `Docker Deployment`: Watchtower updates this component in accordance with the user’s release channel subscription. |
| | | Manages lifecycle of `Meshery Operator`; `Adapters`, `UI`, `Load Generators`, `Database` |
| `Meshery Operator` | `MeshSync` | `Meshery Operator` manages the lifecycle of this component and its sub-components. |
|  | `Meshery Broker` | `Meshery Operator` manages the lifecycle of this event bus component. | 
| `mesheryctl`       |                | `mesheryctl` manages the lifecycle of `Meshery Server`.|
|        |                | - system start calls system update by default, which updates server and existing adapters, but doesn’t update meshery.yaml.|
|        |                | - system reset retrieving docker-compose.yaml from `GitHub` (use git tag to reset to the right `Meshery` version)|
|        |                | - system context manages config.yaml, which manages meshery.yaml|
|        |                | `mesheryctl` should generally be checking for latest release and informing user.|
| `Remote Providers`  | `Meshery Cloud`  |  `Process Extension`: Integrators manage the lifecycle of their `Remote Providers`. Process is unique per provider. |    
|   | `Meshery Extensions`  | `Static Extension`: Integrators manage the lifecycle of their `Meshery Extensions`. Process is unique per provider.
  |

As an application, Meshery is a composition of different functional components. Some of the components must be upgraded simultaneously, while others may be upgraded independently.

## Upgrading Meshery Server, Adapters, and UI

Various components of Meshery will need to be upgraded as new releases become available. Meshery is comprised of a number of components including a server, adapters, UI, and CLI.

### Meshery Docker Deployments

In order to upgrade Meshery Server, Adapters, and UI, execute the following command:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system upgrade
 </div></div>
 </pre>

### Meshery Kubernetes Deployments

Use `kubectl apply` or `helm` to upgrade the Meshery application manifests in your Kubernetes cluster.

## Upgrading `mesheryctl`

The Meshery command line client is available in different package managers. Use the instructions relevant to your environment.

### Upgrading `mesheryctl` using Homebrew

<p>To upgrade `mesheryctl`, execute the following command:</p>

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew upgrade mesheryctl
 </div></div>
 </pre>

### Upgrading `mesheryctl` using Bash

Upgrade `mesheryctl` and run Meshery on Mac or Linux with this script:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

### Upgrading `mesheryctl` using Scoop

To upgrade `mesheryctl`, execute the following command:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 scoop update mesheryctl
 </div></div>
 </pre>
