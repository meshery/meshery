---
layout: default
title: Meshery Operator CRD Reference
abstract: "Details for the Custom Resource Definitions included in Meshery Operator and used by it's custom controllers."
permalink: reference/meshery-operator-crds
redirect_from: reference/
type: Reference
language: en
---

Included in [Meshery Operator]({{site.baseurl}}/concepts/architecture/operator) are a couple of Kubernetes Custom Resource Definitions (CRDs) and a ConfigMap.

## Broker CRD

The CRD is used to configure [Broker]({{site.baseurl}}/concepts/architecture/broker) instances in a cluster.

## YAML synopsis

The following section shows a summary of the structure of the Custom Resource and the required fields.

```yaml
apiVersion: meshery.layer5.io/v1alpha1
kind: Broker
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
  size:
```

### Broker CRD Properties

The following section outlines the fields and their descriptions

- **apiVersion** – API version being used. Must be **v1alpha1** as its the only version supported at the moment.
- **kind** – Resource type. Must be set to **Broker**, also helps in quering for custom resources in the cluster using its plural form **brokers**
- **metadata** - The metadata section allows us to pass data that uniquely identifies a specific custom resource. For Broker, the following metadata is required
  - **name** : The name of this custom resource
  - **namespace**: The namespace that this custom resource will live in, usually **meshery** namespace
  - **labels**: labels are used to organize kubernetes objects and can be used to filter for objects either by kubectl or the Kubernetes API
    - **app**: The name of the application, in this case **meshery**
    - **component**: In the architecture diagram of meshery, the section that this application belongs to, for this case the **controller**
    - **version**: The current version of the meshery application as from its release
  - **annotations**: Annotations are used to provide non-indentifying attributes of a resource i.e cannot be used in filtration, but are informational attributes of an object
    - **meshery/component-type**: The component type pf this custom resource with respect to meshery design, for this case **management-plane**
- **spec**
  The specification section defines the desired state of our custom resource that Kubernetes can then use to take corrective measures to bring the cluster to.
  - **size**: The size is an integer value denoting the number of Broker instances that that should be in one cluster, currently it is adviced to have one Broker instance in a cluster but that can be scaled vertically up or down depending on load.

### Usage

The following section defines the usage for this Custom Resource.

### Example Use

```yaml
apiVersion: meshery.layer5.io/v1alpha1
kind: Broker
metadata:
  name: meshery-broker
  namespace: meshery
  labels:
    app: meshery
    component: controller
    version: v0.1.15
  annotations:
    meshery/component-type: management-plane
spec:
  size: 1
```
## MeshSync CRD

The MeshSync CRD is used as a configuration tool the [MeshSync](https://docs.meshery.io/concepts/architecture/MeshSync). The CRD is used to control the replica count for MeshSync instances, [Broker](https://docs.meshery.io/concepts/architecture/broker) configuration and the resources that MeshSync watches and listens to in a cluster.

### YAML Synopsis

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

### MeshSync CRD Properties

The following section outlines the fields and their descriptions

- **apiVersion** – API version being used. Must be **v1alpha1** as its the only version supported at the moment.

- **kind** – Resource type. Must be set to **MeshSync**, also helps in quering for custom resources in the cluster using its plural form **meshsyncs**

- **metadata** - The metadata section allows us to pass data that uniquely identifies a specific custom resource. For MeshSync, the following metadata is required

  - **name** : The name of this custom resource

  - **namespace**: The namespace that this custom resource will live in, usually **meshery** namespace
  - **labels**: labels are used to organize kubernetes objects and can be used to filter for objects either by kubectl or the Kubernetes API
    - **app**: The name of the application, in this case **meshery**
    - **component**: In the architecture diagram of meshery, the section that this application belongs to, for this case the **controller**
    - **version**: The current version of the meshery application as from its release
  - **annotations**: Annotations are used to provide non-indentifying attributes of a resource i.e cannot be used in filtration, but are informational attributes of an object
    - **meshery/component-type**: The component type pf this custom resource with respect to meshery design, for this case **management-plane**

- **spec**
  The specification section defines the desired state of our custom resource that Kubernetes can then use to take corrective measures to bring the cluster to.
  - **size**: The size is an integer value denoting the number of MeshSync instances that that should be in one cluster, currently only one instance of MeshSync is supported.
  - **broker**:
    This section contains the configuration for [Broker](https://docs.meshery.io/concepts/architecture/broker)
    - **native**:
      Configuration instruction for a broker instance in the cluster.
      - **name**:
        The name of the broker application to be deployed.
      - **namespace**:
        The namespace that the broker would run in.
    - **custom**:
      Alternatively, one can supply the url to an already running instance of the broker. \* **url**: The URL to the external broker instance.
  - **watch-list**:
    The watchlist is a configMap that is used to set the resources e.g. pods,deployments and event types e.g. ADDED, MODIFIED... that Meshync tracks or ignores in the Kubernetes cluster.
    - **whitelist**: Whitelisted resources are resources that MeshSync should track and is an array of the format
      ```yaml
      [
        #Name of the resource eg replicasets.v1.apps in resource.version.group format
        Resource string
        #Events to track, Either ADDED, MODIFIED or DELETED
        Events   []string
      ]
      ```
      The string should be json enconded for example
      ```yaml
      '[{"Resource":"namespaces.v1.","Events":["ADDED","MODIFIED","DELETED"]}]'
      ```
    - **blacklist**:
      Alternatively one can use the blacklist field to define resources that should be ingored by MeshSync. The resources and their respective events will not be tracked.
      The blacklist is an array of resources of the format **_resource.version.group_** format
      ```yaml
      '["namespaces.v1.","configmaps.v1."]'
      ```

### Usage

The following section defines the usage for this Custom Resource

The Custom Resource is used to configure the [Broker](https://docs.meshery.io/concepts/architecture/broker) and Events to be tracked by [MeshSync](https://docs.meshery.io/concepts/architecture/MeshSync)

#### Example Use

```yaml
apiVersion: meshery.layer5.io/v1alpha1
kind: MeshSync
metadata:
  name: meshery-MeshSync
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
      whitelist: '[{"Resource":"namespaces.v1.","Events":["ADDED","MODIFIED","DELETED"]}]'
```

