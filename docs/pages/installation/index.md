---
layout: page
title: Installation Guide
permalink: /installation
---

<a name="getting-started"></a>

# Quick Start

<h6>Getting Meshery up and running locally on a Docker-enabled system is easy with Meshery's command line interface, <a href="/docs/guides/mesheryctl">mesheryctl</a>:</h6>

###### 1. Configure your local environment to support Meshery

Meshery works with both, Docker and Kubernetes enabled systems. Setup the prerequisites on any of Meshery's [supported platforms](/docs/installation/platforms)

###### 2. Install Meshery

Use the Meshery command line interface, [mesheryctl](/docs/guides/mesheryctl), to install and start Meshery

###### 3. Access Meshery

When Meshery is up and running, instructions to access Meshery on the host:`9081` will be printed on your screen. Your default browser should be auto-directed to the Meshery login screen

###### 4. Choose your [provider](/docs/extensibility#providers)

<a href="#meshery-login-page">
  <img style="width:300px;" src="/docs/assets/img/meshery-server-page.png" />
</a>
<a href="#" class="lightbox" id="meshery-login-page">
  <span style="background-image: url('/docs/assets/img/meshery-server-page.png')"></span>
</a>

###### 5. Sign in with your preferred method

<a href="#meshery-sign-up">
<img style="width:300px;height=auto;" src="/docs/assets/img/meshery-login-page.png" />
</a>
<a href="#" class="lightbox" id="meshery-sign-up">
  <span style="background-image: url('/docs/assets/img/meshery-login-page.png')"></span>
</a>

###### 6. You will now be directed to the Meshery UI

<a href="#meshery-ui">
<img style="width:450px;height=auto;" src="/docs/assets/img/adapters/meshery-ui.png" />
</a>
<a href="#" class="lightbox" id="meshery-ui">
  <span style="background-image: url('/docs/assets/img/adapters/meshery-ui.png')"></span>
</a>

###### 7. Ensure that your kubernetes cluster is connected to Meshery. Go to <i class="fas fa-cog"></i> Settings:

- Meshery attempts to auto detect your kubernetes config if it is stored in the default path (`$HOME/.kube` directory) on your system. If your configuration has been auto-detected, you will be able to see your configuration details listed

  <a href="#meshery-settings">
  <img style="width:600px;" src="/docs/assets/img/adapters/meshery-settings.png" />
  </a>
  <a href="#" class="lightbox" id="meshery-settings">
  <span style="background-image: url('/docs/assets/img/adapters/meshery-settings.png')"></span>
  </a>

- If your config has not been auto-detected, you may manually locate and upload your **kube config** file and select the `context name` (docker-desktop, kind-clsuter, minikube etc.)

###### 8. Check up on your connection

Check up on the successful configuration of your Kubernetes settings by clicking on your configuration `context` name. You will be notified of your connection status. You can also verify a successful connection between Meshery and its adapters by clicking on any of the available [Adapter Ports](/docs/architecture#adapter-ports).

###### 9. You may now proceed to install and work with any [service mesh](/docs/service-meshes) supported by Meshery

<video class="videoTest" width="750" height="auto" autoplay muted loop>
  <source src="/docs/assets/img/adapters/meshery-ui-setup.mp4" type="video/mp4">
 Your browser does not support the video tag
</video>
