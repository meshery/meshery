---
title: Install on Docker
categories: [docker]
aliases:
- /installation/platforms/docker
display_title: false
image: installation/docker/images/docker.svg
description: Install Meshery on Docker
weight: 10
---

{{< installation/installation_prerequisites >}}

## Deploying Meshery on Docker

Follow these steps to run Meshery with Docker and Docker Compose, a common way to run Meshery on your local machine. To install Docker, see [Getting Started with Docker](https://docs.docker.com/get-started/). To install Docker Compose, see [Installing Docker Compose](https://docs.docker.com/compose/install/).

Start Meshery:

{{< code code="mesheryctl system start -p docker" >}}

## Advanced Configuration

### Customizing Kubernetes Configuration Location

By default, Meshery Server looks for Kubernetes configuration in the `$HOME/.kube` directory of the user it runs as, and `mesheryctl` mounts your `$HOME/.kube` there read-only. To use a different location, set the `KUBECONFIG_FOLDER` environment variable and mount the folder into the `meshery` container.

`mesheryctl system start` rewrites the `meshery` service in `~/.meshery/meshery.yaml` each time it starts, so edits to that file do not persist. To use a custom kubeconfig location, run Meshery with Docker Compose directly, using a copy of [`install/docker/docker-compose.yaml`](https://github.com/meshery/meshery/blob/master/install/docker/docker-compose.yaml), and change the `meshery` service:

{{< code code=`services:
  meshery:
    environment:
      - "KUBECONFIG_FOLDER=/custom/path/to/.kube"
    volumes:
      - /custom/path/to/.kube:/custom/path/to/.kube:ro` >}}

Then start it with `docker compose up -d`. A custom location is useful when you are:

- Providing a Meshery deployment with a predefined Kubernetes context
- Running Meshery in containerized environments with custom kubeconfig paths
- Managing several Kubernetes configurations

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with [mesheryctl system check]({{< ref "reference/references/mesheryctl/system/check.md" >}}).

You're ready to use Meshery. Open Meshery UI in your browser.

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
