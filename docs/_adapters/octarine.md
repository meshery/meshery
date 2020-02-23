---
layout: page
title: Octarine Adapter
name: Octarine
version: v1.0
port: 10003/tcp
project_status: stable
---
| Adapter Status |
| :------------: |
| [{{ page.project_status }}]({{ page.github_link }})|

# Configuration
In order to connect to the Octarine Control Plane the adapter requires the follwing environment variables to be set:
* `OCTARINE_DOCKER_USERNAME` : The docker username needed to pull Octarine's images to the target cluster. Do not use your own docker credentials. Use the ones supplies by Octarine.
* `OCTARINE_DOCKER_EMAIL` : The docker username needed to pull Octarine's images to the target cluster.
* `OCTARINE_DOCKER_PASSWORD` : The docker username needed to pull Octarine's images to the target cluster.
* `OCTARINE_ACC_MGR_PASSWD` : The password that will be assigned to the user 'meshery' in the new account.
* `OCTARINE_CREATOR_PASSWD` : The password needed to create an account in Octarine.
* `OCTARINE_DELETER_PASSWD` : The password needed to delete the account in Octarine.
* `OCTARINE_CP` : The address of the Octarine Control Plane. Example: meshery-cp.octarinesec.com
* `OCTARINE_DOMAIN` : The name that will be assigned to the target cluster in Octarine. Example: meshery:domain

# Features
Octarine provides:
- Policy-based validation that k8s workloads spec are secure.
- Visibility of layer 4-7 traffic between workloads, as well as ingress and egress.
- Encryption and authentication based on mTLS.
- Automation and enforcement of access control policy based on observed traffic.
- Threat detection based on signatures and anomalies.

# Usage
Once the Octarine's data plane services are deployed, the adapter can be used to deploy Bookinfo. The steps here are:
* Enable the target namespace for automatic sidecar injection.
* Deploy Bookinfo to the target namespace.

# Architecture
## Control Plane
![Alt text](../../assets/img/octarine_cparch.jpg?raw=true "Octarine Control Plane")

## Data Plane
![Alt text](../../assets/img/octarine_dparch.jpg?raw=true "Octarine Data Plane")
