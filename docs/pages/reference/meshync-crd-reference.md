---
layout: default
title: MeshSync CRD Reference
abstract: "Meshync CRD Documentation and Reference"
permalink: reference/meshync-crd
redirect_from: reference/
type: Reference
language: en
---

# MeshSync CRD
The MeshSync CRD is used as a configuration tool the the Meshsync service. The CRD is used to control the replica count for Meshsync instances, broker configuration and the resources that Meshsync watches and listens to in a cluster. 

## YAML synopsis

The following section shows a summary of the structure of the Custom Resource and the required fields. 
```yaml

apiVersion: meshery.layer5.io/v1alpha1
kind: MeshSync
metadata:
  name: 
  namespace: 
  labels:
    app: 
    component: 
    version: 
  annotations:
    meshery/component-type: 
spec:
  size: 1
  broker:
    native:
      name: 
      namespace: 
  watch-list:
    apiVersion: 
    data:
        blacklist: 
        whitelist: 

```

## FIELDS
The following section outlines the fields and their descriptions


* **apiVersion** – API version being used. Must be **v1alpha1** as its the only version supported at the moment.

* **kind** – Resource type. Must be set to **MeshSync**, also helps in quering for custom resources in the cluster using its plural form **meshsyncs**

* **metadata** - The metadata section allows us to pass data that uniquely identifies a specific custom resource. For MeshSync, the following metadata is required
    * **name** - The name of this custom resource 

    * **namespace**: The namespace that this custom resource will live in, usually **meshery** namespace
    * **labels**: labels are used to organize kubernetes objects and can be used to filter for objects either by kubectl or the Kubernetes API
      * **app**: The name of the application, in this case **meshery**
      * **component**: In the architecture diagram of meshery, the section that this application belongs to, for this case the  **controller**
      * **version**: The current version of the meshery application as from its release
    * **annotations**: Annotations are used to provide non-indentifying attributes of a resource i.e cannot be used in filtration, but are informational attributes of an object
      * **meshery/component-type**: The component type pf this custom resource with respect to meshery design, for this case **management-plane**

* **spec** 
The specification section defines the desired state of our custom resource that Kubernetes can then use to take corrective measures to bring the cluster to.
  * **size**: The size is an integer value denoting the number of Meshsync instances that that should be in one cluster, currently only one instance of Meshsync is supported.
  * **broker**:
  This section contains the configuration for [Meshery-Broker](https://docs.meshery.io/concepts/architecture/broker) 
    * **native**:
    Configuration instruction for a broker instance in the cluster.
      * **name**: 
      The name of the broker application to be deployed.
      * **namespace**: 
      The namespace that the broker would run in.        
    * **custom**:
    Alternatively, one can supply the url to an already running instance of the broker.
        * **url**:  The URL to the external broker instance.  
  * **watch-list**:
  The watchlist is a configMap that is used to set the resources e.g. pods,deployments and event types e.g. ADDED, MODIFIED... that Meshync tracks or ignores in the Kubernetes cluster. 
    * **whitelist**: Whitelisted resources are resources that Meshsync should track and is an array of the format
    ```yaml
    [
      #Name of the resource eg replicasets.v1.apps in resource.version.group format
      Resource string
      #Events to track, Either ADDED, MODIFIED or DELETED
      Events   []string   
    ]
    ```
    The string should be json enconded
# Usage
The following section defines the usage for this Custom Resource

The Custom Resource is used to configure the [broker](https://docs.meshery.io/concepts/architecture/broker) and Events to be tracked by [Meshysync](https://docs.meshery.io/concepts/architecture/meshsync) 
## Example Use
```yaml
apiVersion: meshery.layer5.io/v1alpha1
kind: MeshSync
metadata:
  name: meshery-meshsync
  namespace: meshery
  labels:
    app: meshery
    component: controller
    version: v0.1.15
  annotations:
    meshery/component-type: management-plane
spec:
  size: 1
  broker:
    native:
      name: meshery-broker
      namespace: meshery
  watch-list:
    apiVersion: v1
    data:
      whitelist: "[{\"Resource\":\"namespaces.v1.\",\"Events\":[\"ADDED\",\"MODIFIED\",\"DELETED\"]}]"
```