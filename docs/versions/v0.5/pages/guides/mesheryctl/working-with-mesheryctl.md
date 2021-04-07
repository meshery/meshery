---
layout: default
title: Using mesheryctl
description: How to use mesheryctl
permalink: guides/mesheryctl/using-mesheryctl 
language: en
type: Guides
category: mesheryctl
---

Meshery's command line interface is `mesheryctl`. Use `mesheryctl` to both manage the lifecyle of Meshery itself and to access and invoke any of Meshery's application and service mesh management functions. `mesheryctl` commands can be categorized as follows:

- `mesheryctl` - Global flags
- `mesheryctl system` - Meshery Lifecycle and Troubleshooting
- `mesheryctl mesh` - Service Mesh Lifecycle & Configuration Management
- `mesheryctl perf` -  Service Mesh Performance Management
- `mesheryctl pattern` - Service Mesh Pattern Configuration & Management


## Configuring Meshery Deployments with meshconfig

*Meshconfig* refers to a configuration file found at `~/.meshery/config.yaml`. Your meshconfig file must contain one or more `contexts` in order for any `mesheryctl system` command to work. Each context represents a Meshery deployment.


Each of the `system` commands are used to control Meshery's lifecycle like `system start`, `stop`, `status`, `reset` and so on. 


## Meshery CLI FAQ
#### Question: What are differences between contexts ?
Each context is a block that defines certain parameters that are specific to meshery itself. 
This helps the user to save multiple configurations of the meshery instance, switching between them with ease.

#### Question: Why are contexts necessary ?
Having configured multiple contexts and switched between environments enabled automation during the deployment process.

Approach: 
Context allows us to maining context data in the configuration file similar to kube contexts. The properties could be under the context-name object and the current context can be stored using a current- context key.

#### Question: What is `current-context`?
`current-context` identifies the Meshery deployment that when any `mesheryctl` command is invoked will use the environment described in the `current-context`. You can switch between contexts. Only one context can be the `current-context`.
#### Question: What's the difference between contexts and environments?
Contexts configure Meshery deployments (server, adapters, operator and so on), while environments define a collection of Kubernetes clusters and service meshes under management in Meshery.

#### Question: What does the default meshconfig look like? 
The following template is used to create a config file from scratch. Not all of the following variables are required to be included. Some of the variables may have a null value or may be excluded (e.g. “adapters”).

```
contexts:
   <context1-name>:
	endpoint: <url to meshery server rest api>
	token: <name of token variable in this config file>
platform: <type of platform: ”docker” or “kubernetes”> 
# Future: specify type of kubernetes (e.g. eks)
adapters: <collection of names of service mesh adapters: “istio”,“linkerd”,”consul”,”nginx-sm”,”octarine”,”tanzu-sm”,”citrix-sm”,”kuma”,”osm”,”nsm”> 
# Future: ”app-mesh”,”traefik-mesh”
   <context2-name>:
	endpoint: <url to meshery server rest api>
	token: <name of token variable in this config file>
platform: <type of platform: ”docker” or “kubernetes”>
current-context: <context name>
tokens:
- name: <token1-name>
  location: <token-location>
- name: <token2-name>
  value: <token-value> 
  # Future: allow embedding of token certificate
```

Try it out and see for yourself. Run `mesheryctl system context create test` and `mesheryctl system context view test`.

## Advanced Installation

Users can control the specific container image and tag (version) of Meshery that they would like to run by editing their local *~/.meshery/meshery.yaml* (a docker compose file).
Aligned with the Meshery container image, instead of leaving the implicit :stable-latest tag behind image: layer5/meshery, users will instead identify a specific image tag like so:

```
bash
version: '3'
services:
  meshery:
    image: layer5/meshery:v0.5.0
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
```

## Related Guides

<div class="wrapper" style="text-align: left;">
  <div>
  <a href="{{ site.baseurl }}/reference/mesheryctl">
    <div class="overview">Command Reference</div>
  </a>
  <p>Find an exhaustive list of commands and their syntax.</p>
</div>

<div>
  <a href="{{ site.baseurl }}/guides/upgrade">
    <div class="overview">Upgrade Guide</div>
  </a>
  <p>To upgrade <code>mesheryctl</code>, refer to the Upgrade Guide.</p>
</div>
</div>


<!--
## Installing `mesheryctl`

### Mac or Linux

Use your choice of homebrew or bash to install `mesheryctl`. You only need to use one.
### Homebrew

Install `mesheryctl` and run Meshery on Mac with Homebrew.

#### Installing with Homebrew

To install `mesheryctl`, execute the following commands:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew tap layer5io/tap
 brew install mesheryctl
 mesheryctl system start
 </div></div>
 </pre>

**Upgrading with Homebrew**

To upgrade `mesheryctl`, execute the following command:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew upgrade mesheryctl
 </div></div>
 </pre>

#### Bash

**Installing with Bash**

Install `mesheryctl` and run Meshery on Mac or Linux with this script:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

**Upgrading with Bash**

Upgrade `mesheryctl` and run Meshery on Mac or Linux with this script:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

## Windows

### Installing the `mesheryctl` binary

Download and unzip `mesheryctl` from the [Meshery releases](https://github.com/layer5io/meshery/releases/latest) page. Add `mesheryctl` to your PATH for ease of use. Then, execute:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 ./mesheryctl system start
 </div></div>
 </pre>

### Scoop

Use [Scoop](https://scoop.sh) to install Meshery on your Windows machine.

**Installing with Scoop**

Add the Meshery Scoop Bucket and install:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
 scoop install mesheryctl
 </div></div>
 </pre>

**Upgrading with Scoop**

To upgrade `mesheryctl`, execute the following command:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 scoop update mesheryctl
 </div></div>
 </pre>

--> 
