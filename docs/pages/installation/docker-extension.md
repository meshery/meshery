---
layout: default
title: Docker Extension
permalink: installation/platforms/docker-extension
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/docker.svg
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

## Docker Extension for Meshery

The Docker Extension for Meshery extends Docker Desktop’s position as the cloud native developer’s go-to Kubernetes environment with easy access to the next layer of cloud native infrastructure: service meshes.

### Using Docker Desktop

1) Navigate to the Extensions Marketplace of Docker Desktop.

2) From the Dashboard, select Add Extensions in the menu bar or open the Extensions Marketplace from the menu options.

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension-marketplace-1.png">
  <img style="width:350px;" src="{{ site.baseurl }}/assets/img/platforms/docker-extension-marketplace-1.png">
</a>

3) Navigate to Meshery in the Marketplace and press install.

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension.png">
  <img style="width:90%" src="{{ site.baseurl }}/assets/img/platforms/docker-extension.png">
</a>

OR

You can visit the [Docker Hub](https://hub.docker.com/extensions/meshery/docker-extension-meshery) marketplace to directly install Meshery extension in your Docker Desktop.

### Using `Docker CLI`

Meshery runs as a set of containers inside your Docker Desktop virtual machine.

{% capture code_content %}docker extension install meshery/docker-extension-meshery{% endcapture %}
{% include code.html code=code_content %}

