---
layout: default
title: How Meshery Manages Deployments
abstract: Understand how Meshery modifies infrastructure during deployment and undeployment
permalink: guides/infrastructure-management/how-meshery-manages-deployments
category: infrastructure
type: guides
language: en
---

Meshery streamlines deployment and undeployment with model tracking, real-time sync, and automated workflows, ensuring efficient infrastructure management. 

---

## 1. Deployment Architecture

### 1.1 Meshery Registry: Single Source of Truth

The **[Meshery Registry](https://docs.meshery.io/concepts/logical/registry)** is the central repository for tracking **[models](https://docs.meshery.io/concepts/logical/models), [components](https://docs.meshery.io/concepts/logical/components), and [relationships](https://docs.meshery.io/concepts/logical/relationships)** within your infrastructure. Meshery Registry orchestrates deployment/undeployment by enforcing dependency order, ensuring efficient management of resources across cloud-native environments. To explain its structure, we use a university curriculum analogy.

[![Meshery Registry Analogy]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-architecture-registry.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-architecture-registry.svg)

### 1.2 Model Generation
Meshery provides two primary ways to generate infrastructure models, **Static Model Generation** and **Dynamic Model Generation**.
The key differences between these approaches are:

<table>
  <tr>
    <th>Feature</th>
    <th>Static Model Generation</th>
    <th>Dynamic Model Generation</th>
  </tr>
  <tr>
    <td><b>Data Source</b></td>
    <td>Helm, Git, Google Sheets</td>
    <td>Kubernetes API, User Uploads</td>
  </tr>
  <tr>
    <td><b>Processing Method</b></td>
    <td>Batch Processing (Predefined)</td>
    <td>Real-Time (Runtime)</td>
  </tr>
  <tr>
    <td><b>Update Frequency</b></td>
    <td>Periodic (Scheduled Jobs)</td>
    <td>On-Demand (Trigger-Based)</td>
  </tr>
  <tr>
    <td><b>Supported Inputs</b></td>
    <td>Helm Charts, Git Repositories, Google Sheet</td>
    <td>CRDs, Manifests, OCI files</td>
  </tr>
  <tr>
    <td><b>Primary Use Case</b></td>
    <td>Model Definitions & Updates</td>
    <td>Live Cluster State Syncing</td>
  </tr>
</table>

For CLI commands and UI steps, refer to:[Generating Models](https://docs.meshery.io/guides/configuration-management/generating-models)

<br>

#### **Static Generation (Predefined)**

The Static Generation of Models in Meshery allows users to **predefine infrastructure models** using external sources like Helm charts, Git repositories, and Google Sheets.

[![Dynamic Model Generation Flow]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-static-generation-models.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-static-generation-models.svg)

#### Understanding Static Model Generation

- **Batch Processing**: Runs on a **scheduled workflow** to generate models.
- **Input Sources**:
  - **Helm Charts**: Extracts metadata automatically.
  - **Git Repositories**: Defines model structures.
  - **Google Sheets**: Manually curated component lists.
- **Data Flow**:
  - Meshery **downloads a Google Sheet** → **Reads CSV data** → **Generates models**.
  - Converts data into **versioned model packages** stored in **Meshery Registry**.
  - Models are stored in a **versioned format** within the Meshery Registry, making them available for future deployments.
- **Version Control & Updates**:
  - New components are **added incrementally**, avoiding duplication.
  - Models can be **updated and re-synced** using `mesheryctl` commands.

#### **Dynamic Generation (Runtime)**

Meshery supports **runtime model generation**, enabling dynamic discovery and registration of models from **Kubernetes clusters, predefined sources, and user-imported manifests**.

[![Dynamic Model Generation Flow]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-dynamic-generation-models.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-dynamic-generation-models.svg)

#### Key Capabilities of Dynamic Generation

- **Auto-Discovery**: Connects to Kubernetes clusters and extracts CRDs.
- **User-Defined Models**: Supports uploading **unknown models** (e.g., manifests) for registration.
- **Predefined Model Imports**: Accepts OCI, tar.gz, and file-based models for offline deployments.
- **Automatic Validation**: Ensures model consistency before persisting in the registry.

#### Understanding Dynamic Model Generation

- **Server Discovery**: Meshery connects to Kubernetes and detects available CRDs.
- **User Imports Unknown Model**:  
  - Meshery processes an **uploaded manifest** and attempts to register the model.
- **User Imports Predefined Model**:  
  - If a valid format (**OCI, tar.gz, file**) is provided, the model is imported.
  - UI does not retrieve these models for air-gapped environments, but  
    `mesheryctl model import` can be used.
- **Adapter-Driven Model Generation**:  
  - Meshery Adapters generate **components, relationships, and policies** dynamically.
- **Validation & Registration**:  
  - Model is **validated and stored** in **Meshery Registry**.
  - **Future**: Cloud storage support for **multi-tenancy**.

---

## 2. Deployment Validation & State Management

### 2.1 MeshSync: Real-Time Synchronization

[MeshSync](https://docs.meshery.io/concepts/architecture/meshsync) is a Kubernetes-native synchronization mechanism, ensuring Meshery maintains an **accurate and up-to-date** understanding of the infrastructure it manages. It continuously tracks state changes, whether resources are created by Meshery itself (**greenfield**) or pre-existing (**brownfield**).

**Key Features:**
-  **GVK Metadata Tracking** – Identifies Kubernetes resources uniquely based on Group, Version, and Kind (GVK).
-  **Prometheus Health Checks** – Monitors and validates the state of infrastructure components.
-  **State Mapping (`REGISTERED` vs `CONNECTED`)** – Differentiates between known and actively monitored resources.

MeshSync operates on an **event-driven architecture**, leveraging NATS messaging for **high-speed and scalable** data synchronization. It can **detect infrastructure changes in real time**, making it a core component of Meshery’s **multi-cloud and cloud-native** management capabilities.

### 2.2 Versioned Entity Lifecycle

Meshery employs a structured versioning model to ensure that infrastructure components are well-defined, traceable, and adaptable to changes over time. The diagram above illustrates the different stages of versioning, helping users understand how Meshery manages and evolves its entities.

[![Meshery Entities]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-entities.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-entities.svg)

#### Understanding Meshery’s Versioning Model

Each Meshery entity follows a structured versioning lifecycle, ensuring consistency before deployment and management. This process enables users to track changes, apply upgrades, and maintain compatibility across environments.

- **Schema** (`v1alpha1`): Defines the fundamental structure of an entity, ensuring consistency and validation.  
- **Definition**: Establishes the expected behavior of an entity, detailing constraints and relationships.  
- **Declaration**: Provides a human-readable YAML-based configuration that Meshery uses for deployment.  
- **Instance**: Represents the actual deployed resource, tracking runtime changes and operational metadata.  

---

## 3. Core Operation Logic

### 3.1 Deploy Action Breakdown

Meshery follows a structured process when deploying infrastructure components, ensuring that dependencies are validated and operations are executed efficiently.

[![Deploy Action Flow]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-deploy-action.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-deploy-action.svg)

#### Deployment Flow
- **Parse Design File** → Meshery reads and interprets the provided `meshery-design.yaml` file.
- **Validate Dependencies** → Ensures that required CRDs, operators, and other components are available.
- **Deploy via Adapter** → Delegates to Meshery Adapters (e.g., Istio, Linkerd) via gRPC, enabling interaction with the target platform (e.g., Kubernetes).
- **Sync Status with MeshSync** → Meshery tracks deployment status through **MeshSync**, ensuring real-time updates.

#### Error Handling and Recovery
- If a dependency is missing, Meshery will either attempt to install it or notify the user.
- Failures during deployment are logged and surfaced through the UI and CLI.

### 3.2 Undeploy Action Breakdown & Dependency Cleanup

Meshery **relies on Adapters** to execute `undeploy`. Different adapters may use **Helm** or **direct Kubernetes resource deletion** to perform the undeployment process.

#### Undeploy Process
- **Remove Workloads** → Deletes `Pods`, `Deployments`, `StatefulSets`, etc.
- **Remove Services** → Deletes `Services`, `Ingress`, and `VirtualService`.
- **Remove CRDs & Operators (Optional)**:
  - If the deployment **was done via Helm**, Meshery will call Helm’s built-in undeploy mechanisms, ensuring that all dependencies (including CRDs) are removed.
  - If the deployment **was done via Git**, **Meshery does not track Helm charts** inside the repository, which can lead to **orphaned CRDs**.

#### Helm vs Git Handling
<table>
  <thead>
    <tr>
      <th>Source</th>
      <th>Dependency Tracking</th>
      <th>Undeploy Behavior</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><strong>Helm</strong></td>
      <td> <strong>Automatic</strong></td>
      <td> Dependencies are removed properly</td>
    </tr>
    <tr>
      <td><strong>Git</strong></td>
      <td> <strong>Partial / Manual</strong></td>
      <td> Dependencies may not be fully removed</td>
    </tr>
  </tbody>
</table>

#### Limitations of Git-based Dependencies
Meshery **does not track the original Git source reference** in its registry, which leads to the following limitations:
- **CRDs may persist** even after undeployment.
- **Manual cleanup is required** to avoid conflicts or inconsistencies in the cluster.

> ⚠️ **Note:**  
> If your model was deployed using **Git**, Meshery does **not** automatically remove CRDs.  
> To ensure a clean undeployment, manually delete any remaining CRDs and namespaces:

```bash
kubectl delete crd <your-crd>
kubectl delete namespace <your-namespace>
```

#### Recommended Workaround
To prevent issues with orphaned dependencies:
- **Track dependencies manually** when deploying from Git.
- **Ensure proper cleanup** of CRDs and operators after undeployment.
- Refer to [Helm Dependency Management](https://helm.sh/docs/helm/helm_dependency/){:target="_blank"} for best practices.

---

### 3.3 Sample Application Deployment 

Meshery provides a set of [Sample Applications](https://docs.meshery.io/guides/infrastructure-management/sample-apps), designed for testing and learning cloud-native infrastructure.

#### Available Sample Applications
- **BookInfo** → A sample application originally built by Istio to test service meshes.
- **Emojivoto** → A microservices demo app from Linkerd that allows users to vote for their favorite emoji.
- **ImageHub** → A WebAssembly-based sample application for experimenting with Envoy filters.
- **HTTPBin** → A simple HTTP request and response service for testing API calls.
- **Online Boutique** → A cloud-native e-commerce demo application originally built by Google.

---

