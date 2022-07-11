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

{% include installation_prerequisites.html %}

## Docker Extension for Meshery

The Docker Extension for Meshery extends Docker Desktop’s position as the cloud native developer’s go-to Kubernetes environment with easy access to the next layer of cloud native infrastructure: service meshes.

### Using Docker Desktop

1) Navigate to the Extensions area of Docker Desktop.

2) From the Dashboard, select Add Extensions in the menu bar or open the Extensions Marketplace from the menu options.

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension-marketplace-1.png">
  <img style="width:350px;" src="{{ site.baseurl }}/assets/img/platforms/docker-extension-marketplace-1.png">
</a>

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension-marketplace-2.png">
  <img style="width:350px;" src="{{ site.baseurl }}/assets/img/platforms/docker-extension-marketplace-2.png">
</a>

3) Navigate to Meshery in the Marketplace and press install.

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension.png">
  <img style="width:550px;height:80px" src="{{ site.baseurl }}/assets/img/platforms/docker-extension.png">
</a>

### Using `Docker CLI`

Meshery runs as a set of containers inside your Docker Desktop virtual machine.

<pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 docker extension install meshery/docker-extension-meshery
 </div></div>
 </pre>

