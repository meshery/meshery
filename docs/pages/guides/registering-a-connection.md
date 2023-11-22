---
layout: default
title: Managing Events with Notification Center
description: A Comprehensive Guide to Managing Events with Meshery
permalink: guides/events-management
type: guides
language: en
abstract: "Meshery tracks operations that you perform on infrastructurees and their workloads. Meshery provides notification of environment issues, application conflicts with infrastructure configuration, and so on."
---

MeshSync is your tool for efficient management and synchronization of Kubernetes clusters! This user guide will walk you through how to establish a connection with any element that Meshery is able to connect to and manage. Use the MeshSync page to register connections, view connection metadata, and perform connectivity tests, as weel as to manage credentials and environments.

## MeshSync Discovery

MeshSync involves the discovery and registration of elements in your infrastructure. Follow these steps to seamlessly integrate and manage your connections:

1. **Discovered**: Fill out the registration form, specifying metrics to listen to for Prometheus.
2. **Register**: Your object will now appear on the connection page.

Tasks performed during registration:

- Enrich connections with more information about your infrastructure.
- Offload fingerprinting process to MeshSync.
- Connect with KubeAPI for connections management.
- Handle curated list of connections (port, protocol, etc.).
- Display Meshery Models icons/colors for every discovered object with connection definitions pointing to components.

### Meshery Server Database Updates

Meshery Server writes updates from MeshSync to the server database:

Meshery uses the component's Group, Version, Kind to perform the initial tier of fingerprinting. 

### MeshSync Endpoint

The Meshery Server endpoint performs a simple JOIN of GVKs to enrich responses with Component Metadata. Key points to note:

## Connection Management

### During Connection Registration

Upon registration, the MeshSync object should have a status of "REGISTERED." Follow these steps:

1. Create a new entry in the Connections table with either "REGISTERED" or "CONNECTED," based on the Meshery Server's ability to establish a connection.
2. Use the Registration Modal to assign the connection to an environment.
3. Confirm connection metadata details, including the assigned model and component.
4. Perform a connectivity test to determine the status: "REGISTERED" (administratively connected) or "CONNECTED" (both administratively and actually connected).

### Registration Modal

The Registration Modal allows users to assign connections to environments, view connection metadata, and perform connectivity tests. Key features include:

- Connectivity test to determine registration status.
- Identification of credentials to be used, with a selection based on CredentialDefinition for the ConnectionDefinition being registered.
- Future options to offer generic CredentialDefinitions for various authentication methods.

### Credential Registration

During credential registration, users can assign credentials to one or more connections and perform ad hoc connectivity tests. UI buttons include "Test" and "Register."

## Connection Configuration

The benefit of normalizing and extracting the status of a component as a direct property of the connection is to allow different systems to connect to the same component with different states. For example, multiple Meshery servers can connect to the same Kubernetes cluster, each having its own individual connection with a unique status.

Now that you are familiar with MeshSync and its powerful features, you are ready to streamline your Kubernetes cluster management. Happy syncing!