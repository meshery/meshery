---
title: "Upgrading Meshery"
description: "How to upgrade Meshery and all of its components"
weight: 50
aliases:
  - /installation/platforms/upgrades
display_title: "false"
---

# Upgrade Guide

## Upgrading Meshery Server, Adapters, and UI

Various components of Meshery will need to be upgraded as new releases become available. Meshery is comprised of a number of components including a server, adapters, UI, and CLI. As an application, Meshery is a composition of different functional components.

<p style="text-align:center">
<a href="/images/architecture/upgrading-meshery.svg">
    <img src="/images/architecture/upgrading-meshery.svg" style="margin: 1rem;" />
</a><br /><i><small>Figure: Meshery components</small></i>
</p>

Some of the components must be upgraded simultaneously, while others may be upgraded independently. The following table depicts components, their versions, and deployment units (deployment groups).

### Versioning of Meshery components

<table class="mesherycomponents">
    <tr>
        <th>Components</th>
        <th>Sub-component</th>
        <th>Considering or Updating</th>
    </tr>
    <tr>
        <td class="childcomponent">Meshery Adapters</td>
        <td>Any and All Adapters</td>
        <td>Docker Deployment: Watchtower updates this component in accordance with the user's release channel subscription.</td>
    </tr>
    <tr>
        <td rowspan="3" class="childcomponent">Meshery Server</td>
        <td>Meshery UI</td>
        <td rowspan="3">Manages lifecycle of Meshery Operator; Adapters, UI, Load Generators, Database.<br /><br />
Docker Deployment: Watchtower updates this component in accordance with the user's release channel subscription.</td>
    </tr>
    <tr>
        <td>Load Generators</td>
    </tr>
    <tr>
        <td>Database</td>
    </tr>
    <tr>
        <td rowspan="2" class="childcomponent">Meshery Operator</td>
        <td>MeshSync</td>
        <td>Meshery Operator manages the lifecycle of this component and its sub-components.</td>
    </tr>
    <tr>
        <td>Meshery Broker</td>
        <td>Meshery Operator manages the lifecycle of this event bus component.</td>
    </tr>
    <tr>
        <td class="childcomponent">`mesheryctl`</td>
        <td></td>
        <td><code>mesheryctl</code> manages the lifecycle of Meshery Server. <br /><br />
        <ul> 
            <li><code>system start</code> calls system update by default, which updates the server and existing adapters, but doesn't update <code>meshery.yaml</code>. Unless the <code>skipUpdate</code> flag is used, operators are also updated here.</li>
            <li><code>system reset</code> retrieves <code>docker-compose.yaml</code> from GitHub (use a Git tag to reset to the right Meshery version).</li>
            <li><code>system restart</code> also updates operators, unless the <code>skipUpdate</code> flag is used.</li>
            <li><code>system update</code> updates operators in case of both Docker and Kubernetes deployments.</li>
            <li><code>system context</code> manages <code>config.yaml</code>, which manages <code>meshery.yaml</code>.</li>
            <li><code>mesheryctl</code> should generally check for the latest release and inform the user.</li>
        </ul>
        </td>
    </tr>
    <tr>
        <td rowspan="2" class="childcomponent"><a style="color:white;" href="/extensibility/providers">Remote Providers</a></td>
        <td>Meshery Cloud</td>
        <td>Process Extension: Integrators manage the lifecycle of their Remote Providers. The process is unique per provider.</td>
    </tr>
    <tr>
        <td>Meshery Extensions</td>
        <td>Static Extension: Integrators manage the lifecycle of their Meshery Extensions. The process is unique per provider.</td>
    </tr>
</table>

Sub-components deploy as a unit; however, they do not share the same version number.

### Meshery Docker Deployments

In order to pull the latest images for Meshery Server, Adapters, and UI, execute the following command:

{{< code >}}
mesheryctl system update
{{< /code >}}

If you wish to update a running Meshery deployment with the images you just pulled, you'll also have to execute:

{{< code >}}
mesheryctl system restart
{{< /code >}}

### Meshery Kubernetes Deployments

Use `kubectl apply` or `helm` to upgrade the Meshery application manifests in your Kubernetes cluster.

## Upgrading Meshery CLI

The Meshery command-line client, `mesheryctl`, is available in different package managers. Use the instructions relevant to your environment.

### Upgrading `mesheryctl` using Homebrew

To upgrade `mesheryctl`, execute the following command:

{{< code >}}
brew upgrade mesheryctl
{{< /code >}}

### Upgrading `mesheryctl` using Bash

Upgrade `mesheryctl` and run Meshery on Mac or Linux with this script:

{{< code >}}
curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -
{{< /code >}}

### Upgrading `mesheryctl` using Scoop

To upgrade `mesheryctl`, execute the following command:

{{< code >}}
scoop update mesheryctl
{{< /code >}}

{{< related-discussions tag="meshery" >}}
