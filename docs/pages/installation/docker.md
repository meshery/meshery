---
layout: default
title: Docker
permalink: installation/docker
type: installation
category: docker
redirect_from:
- installation/platforms/docker
display-title: "false"
language: en
list: include
image: /assets/img/platforms/docker.svg
abstract: Install Meshery on Docker
---

{% include installation/installation_prerequisites.html %}

## Deploying Meshery on Docker

Follow these installation steps to use Docker and Docker Compose to run Meshery. Users often choose this installation approach in order to run Meshery on their local machine. If you need to install *docker*, see [Getting Started with Docker](https://docs.docker.com/get-started/) and if you need to install *docker-compose*, see [Installing Docker Compose](https://docs.docker.com/compose/install/).

Start Meshery by executing:

<!-- <pre class="codeblock-pre" style="padding: 0; font-size:0px;"><div class="codeblock" style="display: block;">
 <div class="clipboardjs" style="padding: 0">
 <span style="font-size:0;">curl -L https://meshery.io/install | PLATFORM=docker bash -</span>  
 </div>
 <div class="window-buttons"></div>
 <div id="termynal1" style="width:100%; height:150px; max-width:100%;" data-termynal="">
            <span data-ty="input">curl -L https://meshery.io/install | PLATFORM=docker bash -</span>
            <span data-ty="progress"></span>
            <span data-ty="">Successfully installed Meshery</span>
  </div>
 </div>
 </pre>

 <script src="/assets/js/terminal.js" data-termynal-container="#termynal1"></script> -->

{% capture code_content %}mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment, using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{% include_cached installation/accessing-meshery-ui.md %}

{% include related-discussions.html tag="meshery" %}
