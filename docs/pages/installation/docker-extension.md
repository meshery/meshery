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

Navigate to the Extensions area of Docker Desktop.

### Using `docker`

Meshery runs as a set of containers inside your Docker Desktop virtual machine.

<pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 docker extension install meshery/docker-extension-meshery
 </div></div>
 </pre>

