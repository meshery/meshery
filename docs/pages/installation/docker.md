---
layout: default
title: Docker
permalink: installation/platforms/docker
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/docker.svg
---

{% include installation_prerequisites.html %}

## Deploying Meshery on Docker

Follow these installation steps to use Docker and Docker Compose to run Meshery. Users often choose this installation approach in order to run Meshery on their local machine. If you need to install *docker*, see [Getting Started with Docker](https://docs.docker.com/get-started/) and if you need to install *docker-compose*, see [Installing Docker Compose](https://docs.docker.com/compose/install/). 

Start Meshery by executing:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 mesheryctl system start
 </div></div>
 </pre>
 - Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start
 </div></div>
 </pre>

Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081). Log into the [Provider](/extensibility/providers) of your choice.

Upon starting Meshery successfully, instructions to access Meshery will be printed on the screen. See these [guides]({{ site.baseurl }}/guides) for usage tips.
