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

## Remove Meshery as a Docker Extension

If you want to remove Meshery as a Docker extension from your system, follow these steps:

1. **Stop Meshery Container:**

   - First, stop the running Meshery container (if it's currently running) using the following Docker command:

   <pre class="codeblock-pre" style="padding: 0; margin-top: 2px; font-size:0px;"><div class="codeblock" style="display: block;">
    <div class="clipboardjs" style="padding: 0">
    <span style="font-size:0;">docker stop meshery-container</span> 
    </div>
    <div class="window-buttons"></div>
    <div id="termynal2" style="width:100%; height:200px; max-width:100%;" data-termynal="">
      <span data-ty="input">docker stop meshery-container</span>
    </div>
    </div>
   </pre>
    
2. **Remove Meshery Container:**

   - After stopping the container, you can remove it using the following command:

   <pre class="codeblock-pre" style="padding: 0; margin-top: 2px; font-size:0px;"><div class="codeblock" style="display: block;">
    <div class="clipboardjs" style="padding: 0">
    <span style="font-size:0;">docker rm meshery-container</span> 
    </div>
    <div class="window-buttons"></div>
    <div id="termynal2" style="width:100%; height:200px; max-width:100%;" data-termynal="">
      <span data-ty="input">docker rm meshery-container</span>
    </div>
    </div>
   </pre>

3. **Remove Meshery Images:**

   - Meshery might have pulled Docker images for its components. You can remove these images using the `docker rmi` command. Replace the image names with the actual ones you want to remove:

   <pre class="codeblock-pre" style="padding: 0; margin-top: 2px; font-size:0px;"><div class="codeblock" style="display: block;">
    <div class="clipboardjs" style="padding: 0">
    <span style="font-size:0;">docker rmi meshery/meshery:latest</span> 
    </div>
    <div class="window-buttons"></div>
    <div id="termynal2" style="width:100%; height:200px; max-width:100%;" data-termynal="">
      <span data-ty="input">docker rmi meshery/meshery:latest</span>
    </div>
    </div>
   </pre>

   <pre class="codeblock-pre" style="padding: 0; margin-top: 2px; font-size:0px;"><div class="codeblock" style="display: block;">
    <div class="clipboardjs" style="padding: 0">
    <span style="font-size:0;">docker rmi meshery/adapters:latest</span> 
    </div>
    <div class="window-buttons"></div>
    <div id="termynal2" style="width:100%; height:200px; max-width:100%;" data-termynal="">
      <span data-ty="input">docker rmi meshery/adapters:latest</span>
    </div>
    </div>
   </pre>

      ...and so on for other Meshery-related images

4. **Remove Meshery Volumes (if necessary):**

   - Meshery may have created Docker volumes to persist data. You can list and remove these volumes using the `docker volume ls` and `docker volume rm` commands. For example:

   <pre class="codeblock-pre" style="padding: 0; margin-top: 2px; font-size:0px;"><div class="codeblock" style="display: block;">
    <div class="clipboardjs" style="padding: 0">
    <span style="font-size:0;">docker volume ls</span> 
    </div>
    <div class="window-buttons"></div>
    <div id="termynal2" style="width:100%; height:200px; max-width:100%;" data-termynal="">
      <span data-ty="input">docker volume ls</span>
    </div>
    </div>
   </pre>

   <pre class="codeblock-pre" style="padding: 0; margin-top: 2px; font-size:0px;"><div class="codeblock" style="display: block;">
    <div class="clipboardjs" style="padding: 0">
    <span style="font-size:0;">docker volume rm meshery-data-volume</span> 
    </div>
    <div class="window-buttons"></div>
    <div id="termynal2" style="width:100%; height:200px; max-width:100%;" data-termynal="">
      <span data-ty="input">docker volume rm meshery-data-volume</span>
    </div>
    </div>
   </pre>

    ...remove other Meshery-related volumes if present

5. **Remove Docker Network (if necessary):**

   - If Meshery created a custom Docker network, you can remove it using the `docker network rm` command. For example:

   <pre class="codeblock-pre" style="padding: 0; margin-top: 2px; font-size:0px;"><div class="codeblock" style="display: block;">
    <div class="clipboardjs" style="padding: 0">
    <span style="font-size:0;">docker network rm meshery-network</span> 
    </div>
    <div class="window-buttons"></div>
    <div id="termynal2" style="width:100%; height:200px; max-width:100%;" data-termynal="">
      <span data-ty="input">docker network rm meshery-network</span>
    </div>
    </div>
   </pre>

6. **Clean Up Configuration (optional):**
   - If Meshery created configuration files or directories on your host machine, you can remove them manually if you no longer need them.

<script src="{{ site.baseurl }}/assets/js/terminal.js" data-termynal-container="#termynal2"></script>