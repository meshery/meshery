---
layout: enhanced
title: Upgrading Meshery and all of its components
abstract: How to Meshery and all of its components
permalink: guides/upgrade
display-title: "false"
type: guides
category: operating
language: en
abstract: How to upgrade Meshery and all of its components
---

# Upgrade Guide

## Upgrading Meshery Server, Adapters, and UI

Various components of Meshery will need to be upgraded as new releases become available. Meshery is comprised of a number of components including a server, adapters, UI, and CLI. As an application, Meshery is a composition of different functional components.

<p style="text-align:center">
<a href="{{site.baseurl}}/assets/img/architecture/upgrading-meshery.svg">
    <img src="{{site.baseurl}}/assets/img/architecture/upgrading-meshery.svg" style="margin: 1rem;" />
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
        <td>Docker Deployment: Watchtower updates this component in accordance with the user’s release channel subscription.</td>
    </tr>
    <tr>
        <td rowspan="3" class="childcomponent">Meshery Server</td>
        <td>Meshery UI</td>
        <td rowspan="3">Manages lifecycle of Meshery Operator; Adapters, UI, Load Generators, Database.<br /><br />
Docker Deployment: Watchtower updates this component in accordance with the user’s release channel subscription.</td>
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
            <li><code>system start</code> calls system update by default, which updates server and existing adapters, but doesn’t update meshery.yaml. Unless the <code>skipUpdate</code> flag is used, operators are also updated here.</li>
            <li><code>system reset</code> retrieving docker-compose.yaml from GitHub (use git tag to reset to the right Meshery version).</li>
            <li><code>system restart</code> also updates operators, unless the <code>skipUpdate</code> flag is used.</li>
            <li><code>system update</code> updates operators in case of both docker and kubernetes deployments.</li>
            <li><code>system context</code> manages config.yaml, which manages meshery.yaml. </li>
            <li><code>mesheryctl</code> should generally be checking for latest release and informing user.</li>
        </ul>
        </td>
    </tr>
    <tr>
        <td rowspan="2" class="childcomponent"><a style="color:white;" ref="/extensibility/providers">Remote Providers</a></td>
        <td>Meshery Cloud</td>
        <td>Process Extension: Integrators manage the lifecycle of their Remote Providers. Process is unique per provider.</td>
    </tr>
    <tr>
        <td>Meshery Cloud</td>
        <td> Static Extension: Integrators manage the lifecycle of their Meshery Extensions. Process is unique per provider.</td>
    </tr>
</table>

Sub-components deploy as a unit, however, they do not share the same version number.

### Meshery Docker Deployments

In order to pull the latest images for Meshery Server, Adapters, and UI, execute the following command:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl system update</div></div>
 </pre>

If you wish to update a running Meshery deployment with the images you just pulled, you'll also have to execute:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl system restart</div></div>
 </pre>

### Meshery Kubernetes Deployments

Use `kubectl apply` or `helm` to upgrade the Meshery application manifests in your Kubernetes cluster.

## Upgrading Meshery CLI

The Meshery command line client, `mesheryctl`, is available in different package managers. Use the instructions relevant to your environment.

### Upgrading `mesheryctl` using Homebrew

<p>To upgrade `mesheryctl`, execute the following command:</p>

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">brew upgrade mesheryctl</div></div>
 </pre>

### Upgrading `mesheryctl` using Bash

Upgrade `mesheryctl` and run Meshery on Mac or Linux with this script:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -</div></div>
 </pre>

### Upgrading `mesheryctl` using Scoop

To upgrade `mesheryctl`, execute the following command:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">scoop update mesheryctl</div></div>
 </pre>

{% include related-discussions.html tag="meshery" %}
