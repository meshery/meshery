---
layout: page
title: Docker
permalink: installation/platforms/docker
---

## Quick Start with Docker
*Note: a minimum of 4GB RAM is needed for [Istio (and BookInfo sample app)](/docs/service-meshes/adapters/istio/istio) deployments.*

### **Steps**

Follow these installation steps to use Docker and Docker Compose to run Meshery. Users often choose this installation approach in order to run Meshery on their local machine. If you need to install *docker*, see [Getting Started with Docker](https://docs.docker.com/get-started/) and if you need to install *docker-compose*, see [Installing Docker Compose](https://docs.docker.com/compose/install/). 

Meshery repository includes a *docker-compose.yaml* file. We can use *docker-compose* to spin up all the Meshery services by running:

<code> curl -L https://git.io/meshery | bash - </code>

Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).
You will be redirected to a social login page where you can pick one of the available Social Login methods to login to Meshery.

Upon starting Meshery successfully, instructions to access Meshery will be printed on the sceen. Take a look at the [Meshery guides](/docs/guides) for advanced usage tips.
