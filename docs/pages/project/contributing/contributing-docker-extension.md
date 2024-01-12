---
layout: page
title: Contributing to Meshery Docker Extension
permalink: project/contributing/contributing-docker-extension
# redirect_from: project/contributing/contributing-docker-extension/
abstract: How to contribute to Meshery Docker Extension
language: en
type: project
category: contributing
list: include
---

## Prerequisites
To start contributing to Meshery Docker Extension, make sure you have [Docker](https://docs.docker.com/get-docker/) installed on your system.
### Docker Extension for Meshery

The Docker Extension for Meshery extends Docker Desktop’s position as the cloud native developer’s go-to Kubernetes environment with easy access to the next layer of cloud native infrastructure: service meshes.

#### Using Docker Desktop

1) Navigate to the Extensions Marketplace of Docker Desktop.

2) From the Dashboard, select Add Extensions in the menu bar or open the Extensions Marketplace from the menu options.

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension-browse.png">
  <img style="width:350px;" src="{{ site.baseurl }}/assets/img/platforms/docker-extension-browse.png">
</a>

3) Navigate to Meshery in the Marketplace and press install.

<a href="{{ site.baseurl }}/assets/img/platforms/docker-extension.png">
  <img style="width:90%" src="{{ site.baseurl }}/assets/img/platforms/docker-extension.png">
</a>

OR

You can visit the [Docker Hub](https://hub.docker.com/extensions/meshery/docker-extension-meshery) marketplace to directly install Meshery extension in your Docker Desktop.

#### Using `Docker CLI`

Meshery runs as a set of containers inside your Docker Desktop virtual machine.

{% capture code_content %}docker extension install meshery/docker-extension-meshery{% endcapture %}
{% include code.html code=code_content %}


## Set up the server

In the root directory of meshery, run the following command:

### To install/update the UI dependencies:
{% capture code_content %}make ui-setup{% endcapture %}
{% include code.html code=code_content %}


### Start the server locally
{% capture code_content %}make server{% endcapture %}
{% include code.html code=code_content %}

This will ensure that the server is up and running at port 9081

## Set up docker extension Locally

Open another terminal while the server is running,
Go inside the docker-extension directory {% capture code_content %}cd install/docker-extension{% endcapture %}
{% include code.html code=code_content %}

### Build and export UI

 {% capture code_content %}make ui-build{% endcapture %}
{% include code.html code=code_content %}


### UI Development Server

If you want to work on the Docker UI, it will be a good idea to use the included UI development server. You can run the UI development server by running the following command:
 {% capture code_content %}make ui{% endcapture %}
{% include code.html code=code_content %}

Now the meshery docker-extension is up and running.

### Linking the docker extension locally
To see the changes reflected in the docker extension locally and open the devTools window, we can run the command:
 {% capture code_content %}make link{% endcapture %}
{% include code.html code=code_content %}

Now that our local development environment is connected with the meshery docker extension, we can start contributing to it.
