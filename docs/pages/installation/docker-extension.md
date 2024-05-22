---
layout: default
title: Docker Extension
permalink: installation/docker/docker-extension
type: installation
category: docker
redirect_from:
- installation/platforms/docker-extension
display-title: "false"
language: en
list: include
image: /assets/img/platforms/docker.svg
abstract: Install Docker Extension for Meshery
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

The Docker Extension for Meshery extends Docker Desktop’s position as the developer’s go-to Kubernetes environment with easy access to full the capabilities of Meshery's collaborative cloud native management features.

## Install the Docker Meshery Extension

Select one of the following three options to install the Docker Meshery Extension:

- [Install the Docker Meshery Extension](#install-the-docker-meshery-extension)
  - [Using Docker Desktop](#using-docker-desktop)
  - [Using Docker Hub](#using-docker-hub)
  - [Using Docker CLI](#using-docker-cli)
- [Remove Meshery as a Docker Extension](#remove-meshery-as-a-docker-extension)

### Using Docker Desktop

Navigate to the Extensions Marketplace of Docker Desktop. From the Dashboard, select Add Extensions in the menu bar or open the Extensions Marketplace from the menu options.

[![Docker Meshery Extension]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension.png)]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension.png)

### Using Docker Hub

You can find the [Docker Meshery Extension in Docker Hub](https://hub.docker.com/extensions/meshery/docker-extension-meshery) marketplace to install the Docker Meshery Extension.

### Using Docker CLI

Meshery runs as a set of one or more containers inside your Docker Desktop virtual machine.

<!--
{% capture code_content %}docker extension install meshery/docker-extension-meshery{% endcapture %} -->
<!-- {% include code.html code=code_content %} -->

<pre class="codeblock-pre" style="padding: 0; font-size: 0px;">
  <div class="codeblock" style="display: block;">
    <!-- Updated style for clipboardjs -->
    <div class="clipboardjs" style="padding: 0; height: 0.5rem; overflow: hidden;">
      <span style="font-size: 0;">docker extension install meshery/docker-extension-meshery</span> 
    </div>
    <div class="window-buttons"></div>
    <div id="termynal2" style="width: 100%; height: 200px; max-width: 100%;" data-termynal="">
      <span data-ty="input">docker extension install meshery/docker-extension-meshery</span>
      <span data-ty="progress"></span>
      <span data-ty="">Successfully installed Meshery</span>
      <span data-ty="input">mesheryctl system dashboard</span>
    </div>
  </div>
</pre>



## Remove Meshery as a Docker Extension

If you want to remove Meshery as a Docker extension from your system, follow these steps:

**Stop Meshery Container:**

- First, stop the running Meshery container (if it's currently running) using the following Docker command:
{% capture code_content %}$ docker stop meshery-container{% endcapture %}
{% include code.html code=code_content %}
<br />
    
**Remove Meshery Container:**

- After stopping the container, you can remove it using the following command:
{% capture code_content %}$ docker rm meshery-container{% endcapture %}
{% include code.html code=code_content %}
<br />

**Remove Meshery Images:**

- Meshery might have pulled Docker images for its components. You can remove these images using the `docker rmi` command. Replace the image names with the actual ones you want to remove:
{% capture code_content %}$ docker rmi meshery/meshery:latest{% endcapture %}
{% include code.html code=code_content %}
{% capture code_content %}$ docker rmi meshery/adapters:latest{% endcapture %}
{% include code.html code=code_content %}
...and so on for other Meshery-related images
<br />
<br />

**Remove Meshery Volumes (if necessary):**

- Meshery may have created Docker volumes to persist data. You can list and remove these volumes using the `docker volume ls` and `docker volume rm` commands. For example:
{% capture code_content %}$ docker volume ls{% endcapture %}
{% include code.html code=code_content %}
{% capture code_content %}$ docker volume rm meshery-data-volume{% endcapture %}
{% include code.html code=code_content %}
...remove other Meshery-related volumes if present
<br />
<br />

**Remove Docker Network (if necessary):**

- If Meshery created a custom Docker network, you can remove it using the `docker network rm` command. For example:
{% capture code_content %}$ docker network rm meshery-network{% endcapture %}
{% include code.html code=code_content %}
<br />

**Clean Up Configuration (optional):**
- If Meshery created configuration files or directories on your host machine, you can remove them manually if you no longer need them.

<script src="{{ site.baseurl }}/assets/js/terminal.js" data-termynal-container="#termynal2"></script>

{% include related-discussions.html tag="meshery" %}
