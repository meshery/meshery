---
layout: default
title: Octarine
name: Meshery Adapter for Octarine
mesh_name: Octarine
version: v1.0
port: 10003/tcp
project_status: stable
github_link: https://github.com/layer5io/meshery-octarine
image: /assets/img/service-meshes/octarine.svg
permalink: service-meshes/adapters/octarine
---
{% include adapter-status.html %}

## Lifecycle Management

The {{page.name}} can install **{{page.version}}** of the {{page.mesh_name}} service mesh. A number of sample applications for {{page.mesh_name}} can also be installed using Meshery.

### Install {{ page.mesh_name }}

##### **Choose the Meshery Adapter for {{ page.mesh_name }}**

<a href="{{ site.baseurl }}/assets/img/adapters/octarine/octarine-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/octarine/octarine-adapter.png" />
</a>

##### **Click on (+) and choose the {{page.version}} of the {{page.mesh_name}} service mesh.**

<a href="{{ site.baseurl }}/assets/img/adapters/octarine/octarine-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/octarine/octarine-install.png" />
</a>

### Features

1. Policy-based validation that k8s workloads spec are secure.
1. Visibility of layer 4-7 traffic between workloads, as well as ingress and egress.
1. Encryption and authentication based on mTLS.
1. Automation and enforcement of access control policy based on observed traffic.
1. Threat detection based on signatures and anomalies.

### Configuration
In order to connect to the Octarine Control Plane, the adapter requires the following environment variables to be set:

* **OCTARINE_DOCKER_USERNAME** : The docker username needed to pull Octarine's images to the target cluster, supplied by Octarine. Do not use your own docker credentials.
* **OCTARINE_DOCKER_EMAIL** : The docker email, supplied by Octarine.
* **OCTARINE_DOCKER_PASSWORD** : The docker password, supplied by Octarine.
* **OCTARINE_ACC_MGR_PASSWD** : The password that will be assigned to the user 'meshery' in the new account.
* **OCTARINE_CREATOR_PASSWD** : The password needed to create an account in Octarine.
* **OCTARINE_DELETER_PASSWD** : The password needed to delete the account in Octarine.
* **OCTARINE_CP** : The address of the Octarine Control Plane. `Example: meshery-cp.octarinesec.com`
* **OCTARINE_DOMAIN** : The name that will be assigned to the target cluster in Octarine. `Example: meshery:domain`


### Usage

Once the Octarine's data plane services are deployed, the adapter can be used to deploy Bookinfo:

* Enable the target namespace for automatic sidecar injection.
* Deploy Bookinfo to the target namespace.

### Architecture

#### Control Plane

[![Octarine Control Plane](./octarine_cparch.jpg?raw=true)](./octarine_cparch.jpg?raw=true)

#### Data Plane

[![Octarine Data Plane](./octarine_dparch.jpg?raw=true)](./octarine_dparch.jpg?raw=true)

### Sample Applications 

The {{ page.name }} includes the below sample application operation. Meshery can be use to deploy this sample application.

- [Bookinfo]({{ site.baseurl }}/guides/sample-apps#bookinfo) 
    - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
