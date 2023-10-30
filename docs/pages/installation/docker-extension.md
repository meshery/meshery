---
layout: default
title: Docker Extension
permalink: installation/docker/docker-extension
type: installation
category: docker
display-title: "false"
language: en
list: include
image: /assets/img/platforms/docker.svg
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

## Docker Extension for Meshery

The Docker Extension for Meshery extends Docker Desktop’s position as the developer’s go-to Kubernetes environment with easy access to full the capabilities of Meshery's collaborative cloud native management features.

### Using Docker Desktop

1. Navigate to the Extensions Marketplace of Docker Desktop.

2. From the Dashboard, select Add Extensions in the menu bar or open the Extensions Marketplace from the menu options.

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension-marketplace-1.png">
  <img style="width:350px;" src="{{ site.baseurl }}/assets/img/platforms/docker-extension-marketplace-1.png">
</a>

3. Navigate to Meshery in the Marketplace and press install.

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension.png">
  <img style="width:90%" src="{{ site.baseurl }}/assets/img/platforms/docker-extension.png">
</a>

OR

You can visit the [Docker Hub](https://hub.docker.com/extensions/meshery/docker-extension-meshery) marketplace to directly install Meshery extension in your Docker Desktop.

### Using Docker CLI

Meshery runs as a set of one or more containers inside your Docker Desktop virtual machine.

<!--
{% capture code_content %}docker extension install meshery/docker-extension-meshery{% endcapture %} -->
<!-- {% include code.html code=code_content %} -->

<pre class="codeblock-pre" style="padding: 0; font-size:0px;"><div class="codeblock" style="display: block;">
 <div class="clipboardjs" style="padding: 0">
 <span style="font-size:0;">docker extension install meshery/docker-extension-meshery</span> 
 </div>
 <div class="window-buttons"></div>
 <div id="termynal2" style="width:100%; height:200px; max-width:100%;" data-termynal="">
            <span data-ty="input">docker extension install meshery/docker-extension-meshery</span>
            <span data-ty="progress"></span>
            <span data-ty="">Successfully installed Meshery</span>
            <span data-ty="input">mesheryctl system dashboard</span>
  </div>
 </div>
</pre>
<script src="{{ site.baseurl }}/assets/js/terminal.js" data-termynal-container="#termynal2"></script>

{% include suggested-reading.html language="en" %}

{% include related-discussions.html tag="meshery" %}
