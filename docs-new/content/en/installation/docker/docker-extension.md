---
title: "Docker Extension"
description: "Install Docker Extension for Meshery"
weight: 10
aliases:
  - /installation/platforms/docker-extension
image: /images/platforms/docker.svg
display_title: "false"
---

<h1>Quick Start with Docker Extension <img src="/images/platforms/docker.svg" style="width:35px;height:35px;" /></h1>

The Docker Extension for Meshery extends Docker Desktop's position as the developer's go-to Kubernetes environment with easy access to the full capabilities of Meshery's collaborative cloud native management features.

## Prerequisites

- You need Docker Desktop version of 4.10 or higher for this.
- This document applies only when Docker Desktop uses kubeadm with Kubernetes enabled. If you are using kind, please refer to the [Kind](/installation/kubernetes/kind) section.

## Install the Docker Meshery Extension

Select one of the following three options to install the Docker Meshery Extension:

- [Install the Docker Meshery Extension](#install-the-docker-meshery-extension)
  - [Using Docker Desktop](#using-docker-desktop)
  - [Using Docker Hub](#using-docker-hub)
  - [Using Docker CLI](#using-docker-cli)
- [Remove Meshery as a Docker Extension](#remove-meshery-as-a-docker-extension)

### Using Docker Desktop

Navigate to the **Extensions** marketplace of Docker Desktop. Search for Meshery and click the Install button to install the extension.

[![Docker Meshery Extension Install](/images/platforms/docker-desktop-meshery-extension-install.png)](/images/platforms/docker-desktop-meshery-extension-install.png)

Click **Open** when installation is done or click **Meshery** on the left under **Extensions**.

[![Docker Meshery Extension Open](/images/platforms/docker-desktop-meshery-extension-open.png)](/images/platforms/docker-desktop-meshery-extension-open.png)

Click **Login** to open the _Layer5 Cloud_ login page. Log in or sign up and you will be redirected back to Docker Desktop.

[![Docker Meshery Extension Login](/images/platforms/docker-desktop-meshery-extension-login.png)](/images/platforms/docker-desktop-meshery-extension-login.png)

Finally, click **Launch Meshery** to load Meshery Dashboard in a browser window. It runs at http://localhost:9081/ by default.

[![Docker Meshery Extension Launch](/images/platforms/docker-desktop-meshery-extension-launch.png)](/images/platforms/docker-desktop-meshery-extension-launch.png)

You can also open http://localhost:9081/ directly on a browser on the local machine after installing the Docker extension and complete the _Layer5 Cloud_ login process to achieve the same result.

### Using Docker Hub

Another way to install the Meshery Docker Extension is from the Docker Hub. Navigate to the [Meshery Docker Extension](https://hub.docker.com/extensions/meshery/docker-extension-meshery) page and click Open in Docker Desktop to get started. Once installed, the rest of the process is same as above.

[![Docker Hub Extension](/images/platforms/docker-hub-meshery-extension.png)](/images/platforms/docker-hub-meshery-extension.png)

### Using Docker CLI

Finally, you can also install the Meshery Docker Extension using the Docker CLI:

{{< code >}}
docker extension install meshery/docker-extension-meshery
{{< /code >}}

It runs as a set of one or more containers inside your Docker Desktop virtual machine.

Finally, you can now fully utilize Meshery to manage and monitor your cloud-native infrastructure.

## Remove Meshery as a Docker Extension

You can remove the Docker Extension from the Docker Desktop interface or from the CLI.

### Removing from Docker Desktop

Navigate to **Manage** under Extensions, click the ellipsis button (three vertical dots) and select **Uninstall**.

[![Remove Meshery Docker Desktop Extension](/images/platforms/docker-desktop-meshery-extension-remove.png)](/images/platforms/docker-desktop-meshery-extension-remove.png)

### Removing using Docker CLI

To remove the extension from the command line, use the `docker extension rm` command.

{{< code >}}
docker extension rm meshery/docker-extension-meshery
{{< /code >}}

### Additional Cleanup

There could be residual images and networks to remove after removing/uninstalling the extension. Follow the steps below to do so.

**Remove Meshery Images (if necessary)**

Meshery pulls Docker images for deploying the extension and there could be additional Meshery/Layer5 images based on how it was configured. You can remove these images using the `docker rmi` command. Start by listing all the images and then running the command for each image you want to remove. For example:

{{< code >}}
docker rmi meshery/meshery:stable-latest
{{< /code >}}

**Remove Meshery Docker Networks (if necessary)**

Meshery creates custom Docker networks, and they could still be left after the extension uninstall. These can be removed using the `docker network rm` command. For example:

{{< code >}}
docker network rm meshery_default
{{< /code >}}

{{< related-discussions tag="meshery" >}}