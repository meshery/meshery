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

{% include alert.html type="info" title="Prerequisites" content="You need Docker Desktop version of 4.10 or higher for this." %}


## Install the Docker Meshery Extension

Select one of the following three options to install the Docker Meshery Extension:

- [Install the Docker Meshery Extension](#install-the-docker-meshery-extension)
  - [Using Docker Desktop](#using-docker-desktop)
  - [Using Docker Hub](#using-docker-hub)
  - [Using Docker CLI](#using-docker-cli)
- [Remove Meshery as a Docker Extension](#remove-meshery-as-a-docker-extension)

### Using Docker Desktop

Navigate to the **Extensions** marketplace of Docker Desktop. Search for Meshery and click the Install button to install the extension.

[![Docker Meshery Extension Install]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-install.png)]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-install.png)

Click **Open** when installation is done or click **Meshery** on the left under **Extensions**.

[![Docker Meshery Extension Open]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-open.png)]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-open.png)

Click **Login** to open the _Layer5 Cloud_ login page. Login or Sign up and you will be redirected back to Docker Desktop.

[![Docker Meshery Extension Login]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-login.png)]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-login.png)

Finally, click **Launch Meshery** to load Meshery Dashboard on a browser window. It runs at http://localhost:9081/ by default.

[![Docker Meshery Extension Launch]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-launch.png)]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-launch.png)

You can also open http://localhost:9081/ directly on a browser on the local machine after installing the Docker extension and complete the _Layer5 Cloud_ login process to achieve the same result.

### Using Docker Hub

Another way to install the Meshery Docker Extension is from the Docker Hub. Navigate to the [Meshery Docker Extension](https://hub.docker.com/extensions/meshery/docker-extension-meshery) page and click Open in Docker Desktop to get started. Once installed, the rest of the process is same as above.

[![Docker Hub Extension]({{site.baseurl}}/assets/img/platforms/docker-hub-meshery-extension.png)]({{site.baseurl}}/assets/img/platforms/docker-hub-meshery-extension.png)

### Using Docker CLI

Finally, you can also install the Meshery Docker Extension using the Docker CLI. Follow the commands in the clipboard below. 

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

It runs as a set of one or more containers inside your Docker Desktop virtual machine.

## Remove Meshery Docker Extension

You can remove the Docker Extension from Docker Desktop interface or from the CLI. 

### Removing from Docker Desktop

Navigate to **Manage** under Extensions, click the ellipsis button (three vertical dots) and select **Uninstall**.

[![Remove Meshery Docker Desktop Extension]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-remove.png)]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension-remove.png)

### Removing using Docker CLI

To remove the extension from the command line, use the `docker extension rm` command.

{% capture code_content %}$ docker extension rm meshery/docker-extension-meshery{% endcapture %}
{% include code.html code=code_content %}

### Additional Cleanup

There could be redisual Images and Networks to remove after removing/uninstalling the extension. Follow the steps below to do so. 

**Remove Meshery Images (if necessary)**

Meshery pulls Docker images for deploying the extension and there could be additional Meshery/Layer5 images based on how it was configured. You can remove these images using the `docker rmi` command. Start by listing all the images and then running the command for each image you want to remove. For example:

{% capture code_content %}$ docker rmi layer5/meshery:stable-latest{% endcapture %}
{% include code.html code=code_content %}


**Remove Meshery Docker Networks (if necessary)**

Meshery creates custom Docker networks, and they could still be left after the extension uninstall. These can be removed using the `docker network rm` command. For example:

{% capture code_content %}$ docker network rm meshery_default{% endcapture %}
{% include code.html code=code_content %}
<br />


<script src="{{ site.baseurl }}/assets/js/terminal.js" data-termynal-container="#termynal2"></script>

{% include related-discussions.html tag="meshery" %}
