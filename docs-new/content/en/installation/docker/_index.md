---
title: Docker
categories: [docker]
aliases:
- /installation/platforms/docker
display_title: false
image: /installation/docker/images/docker.svg
description: Install Meshery on Docker
---

{{< installation/installation_prerequisites >}}

## Deploying Meshery on Docker

Follow these installation steps to use Docker and Docker Compose to run Meshery. Users often choose this installation approach in order to run Meshery on their local machine. If you need to install *Docker*, see [Getting Started with Docker](https://docs.docker.com/get-started/), and if you need to install *Docker Compose*, see [Installing Docker Compose](https://docs.docker.com/compose/install/).


Start Meshery by executing the following command:

{{< code code="mesheryctl system start -p docker" >}}

## Advanced Configuration

### Customizing Kubernetes Configuration Location

By default, Meshery looks for Kubernetes configuration in the `$HOME/.kube` directory. You can customize this location by setting the `KUBECONFIG_FOLDER` environment variable in your Docker deployment.

To use a custom kubeconfig location with Docker Compose, modify your `docker-compose.yaml`:

{{< code code=`services:
  meshery:
    environment:
      - "KUBECONFIG_FOLDER=/custom/path/to/.kube"
    volumes:
      - /custom/path/to/.kube:/custom/path/to/.kube:ro` >}}

This is useful when:
- Providing a Meshery deployment with a predefined Kubernetes context
- Running Meshery in containerized environments with custom kubeconfig paths
- Managing multiple Kubernetes configurations

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}

### See Also
