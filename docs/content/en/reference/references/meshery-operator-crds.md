---
title: Meshery Operator CRDs
description: Details of the Custom Resource Definitions included in Meshery Operator and used by its custom controllers.
aliases:
- /reference/meshery-operator-crds/
---

Included in [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}) are a couple of Kubernetes Custom Resource Definitions (CRDs) and a ConfigMap.

## Broker CRD

The CRD is used to configure [Broker]({{< ref "concepts/architecture/broker/index.md" >}}) instances in a cluster.

## YAML synopsis

The following section shows a summary of the structure of the Custom Resource and the required fields.

```yaml
apiVersion: meshery.io/v1alpha1
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
#### 2. `### Broker CRD Properties` should NOT have two spaces

It should be:

```markdown
### Broker CRD Properties

The following section outlines the fields and their descriptions

- **apiVersion** – API version being used. Must be **v1alpha1** as its the only version supported at the moment.
- **kind** – Resource type. Must be set to **Broker**, also helps in querying for custom resources in the cluster using its plural form **brokers**
- **metadata** - The metadata section allows us to pass data that uniquely identifies a specific custom resource. For Broker, the following metadata is required
  - **name** : The name of this custom resource
  - **namespace**: The namespace that this custom resource will live in, usually **meshery** namespace
  - **labels**: labels are used to organize kubernetes objects and can be used to filter for objects either by kubectl or the Kubernetes API
    - **app**: The name of the application, in this case **meshery**
    - **component**: In the architecture diagram of meshery, the section that this application belongs to, for this case the **controller**
    - **version**: The current version of the meshery application as from its release
  - **annotations**: Annotations are used to provide non-identifying attributes of a resource i.e cannot be used in filtration, but are informational attributes of an object
    - **meshery/component-type**: The component type of this custom resource with respect to meshery design, for this case **management-plane**
- **spec**

  The specification section defines the desired state of our custom resource that Kubernetes can then use to take corrective measures to bring the cluster to.

  - **size**: The size is an integer value denoting the number of Broker instances that should be in one cluster, currently it is advised to have one Broker instance in a cluster but that can be scaled vertically up or down depending on load.
