---
layout: default
title: How Meshery Manages Deployments
abstract: Learn how Meshery manages and modifies infrastructure components during deployment and teardown processes.
permalink: guides/infrastructure-management/how-meshery-manages-deployments
category: infrastructure
type: guides
language: en
---

This guide helps you understand how Meshery deploys and undeploys infrastructure components in cloud-native environments, and what roles the Meshery [Registry]({{ site.baseurl }}/concepts/logical/registry), [MeshSync]({{ site.baseurl }}/concepts/architecture/meshsync), and model generation play in the process.

By the end of this guide, you'll be able to:

- Understand the difference between static and dynamic model generation
- Know how Meshery handles deployments, dependencies, and undeployments
- Avoid common pitfalls when undeploying Git-based models

## 1. Deployment Architecture

To deploy infrastructure using Meshery, the first thing you need is a **[model]({{ site.baseurl }}/concepts/logical/models)** — a structured description of the resources you want to manage.

Meshery generates models in two ways: **from your own definitions**, or **from what's already running in your Kubernetes cluster**. These two modes are known as:

- **Static Model Generation** – You define the model
- **Dynamic Model Generation** – Meshery discovers the model for you

### 1.1 Meshery Registry: Single Source of Truth

The **[Meshery Registry](https://docs.meshery.io/concepts/logical/registry)** is the central hub that stores all models, components, and their relationships. It helps Meshery know:

- What needs to be deployed
- In what order
- How everything is connected

Think of it like your infrastructure's course catalog — all deployments go through it.

[![Meshery Registry Analogy]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-architecture-registry.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-architecture-registry.svg)

### 1.2 How to Provide a Model

Depending on what you already have, choose one of the following approaches:

#### Option 1: **Static Model Generation**

Use this if you already have a model or want to build one from scratch.

There are two ways to do this:
- **[Create Model]({{ site.baseurl }}/guides/configuration-management/creating-models)** : Recommended for first-time users. Define your model step by step in the UI.
- **[Import Model]({{ site.baseurl }}/guides/configuration-management/importing-models)**: Upload a model file (.json, .csv, or .tar.gz) you’ve prepared elsewhere (e.g., Git, Google Sheets, Helm).

These models are versioned and saved in the Meshery Registry.

| Item | Static Model Generation |
|------|--------------------------|
| **Best for** | Defining infrastructure before deployment |
| **Data sources** | Helm charts, Git repositories, spreadsheets |
| **Processing** | Batch-based (scheduled or on command) |

<details>
<summary>Step-by-Step: How Static Generation Works</summary>

<p>Static generation is like preparing a blueprint before building. You define the structure, and Meshery follows that to manage your infrastructure.</p>

<ul>
  <li><strong>Batch Processing:</strong> Runs on a schedule or via command, ideal for repeatable deployments.</li>
  <li><strong>Input Sources:</strong>
    <ul>
      <li><strong>Helm Charts:</strong> Meshery automatically extracts component metadata.</li>
      <li><strong>Git Repositories:</strong> Model structure is defined as code.</li>
      <li><strong>Google Sheets:</strong> Useful for manually curated lists.</li>
    </ul>
  </li>
  <li><strong>Data Flow:</strong>
    <ul>
      <li>Meshery downloads a Google Sheet → reads CSV data → generates models.</li>
      <li>Models are packaged as versioned bundles and stored in the Meshery Registry.</li>
    </ul>
  </li>
  <li><strong>Version Control & Updates:</strong>
    <ul>
      <li>New components are added incrementally to avoid duplication.</li>
      <li>You can re-sync models using <code>mesheryctl</code> commands.</li>
    </ul>
  </li>
</ul>

</details>


[![Static Model Generation Flow]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-static-generation-models.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-static-generation-models.svg)

#### Option 2: **Dynamic Model Generation**

Use this if your infrastructure is already running and you want Meshery to generate a model from it.

You can:
- Connect Meshery to your Kubernetes cluster to auto-discover CRDs and resources
- Upload a manifest (YAML file) that Meshery analyzes and registers
- Let a Meshery Adapter dynamically create models based on what it manages

| Item | Dynamic Model Generation |
|------|---------------------------|
| **Best for** | Working with existing Kubernetes environments |
| **Data sources** | Cluster APIs, uploaded YAMLs, Adapters |
| **Processing** | Real-time (on demand) |

<details>
<summary>Step-by-Step: How Dynamic Generation Works</summary>

<p>
Dynamic model generation enables Meshery to discover and register models from your existing infrastructure — no need to define anything manually.
</p>

<h4>Key Capabilities of Dynamic Generation</h4>

<ul>
  <li><strong>Auto-Discovery:</strong> Connects to Kubernetes clusters and automatically identifies CRDs.</li>
  <li><strong>User-Defined Models:</strong> Upload manifests and Meshery will attempt to register them.</li>
  <li><strong>Predefined Model Imports:</strong> Accepts OCI images, <code>.tar.gz</code>, or file-based models for offline deployments.</li>
  <li><strong>Automatic Validation:</strong> All models are checked before being stored.</li>
</ul>

<h4>Understanding Dynamic Model Generation</h4>

<p>Dynamic generation is like Meshery walking into your cluster and documenting everything it sees.</p>

<ul>
  <li><strong>Server Discovery:</strong> Connects to a live Kubernetes cluster and detects CRDs.</li>
  <li><strong>User Uploads:</strong>
    <ul>
      <li>Upload unknown manifests → Meshery registers them as models.</li>
      <li>Upload <code>.tar.gz</code> or OCI-formatted files → Meshery imports the model directly.</li>
      <li>For air-gapped environments, use <code>mesheryctl model import</code>.</li>
    </ul>
  </li>
  <li><strong>Adapter-Driven Generation:</strong> Some Meshery Adapters (e.g., for Istio, Linkerd) generate models automatically as part of their operation.</li>
  <li><strong>Validation & Registration:</strong>
    <ul>
      <li>All models are validated before being saved in the Meshery Registry.</li>
      <li>Future support will include cloud-based multi-tenant storage.</li>
    </ul>
  </li>
</ul>

</details>

[![Dynamic Model Generation Flow]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-dynamic-generation-models.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-dynamic-generation-models.svg)

> Example: You already deployed an app in Kubernetes. Meshery connects to the cluster and generates models based on what it finds.

### 1.3 Create or Import? What’s the Difference?

- Use **Create Model** when you want to build a model step by step with help from the UI (ideal for beginners).
- Use **Import Model** if you already have model files ready to go (ideal for advanced users or automation).
- Use **Dynamic Generation** when you don’t want to define anything manually — just let Meshery observe and build models from your cluster or files.

## 2. Deployment Validation & State Management

Once infrastructure is deployed, Meshery continuously monitors its state to ensure everything remains consistent and healthy. It does this through two core mechanisms:

- **MeshSync** – keeps Meshery’s internal model in sync with the actual cluster state
- **Versioned Entities** – ensures all changes are tracked in a structured way

### 2.1 MeshSync: Real-Time Synchronization

[MeshSync](https://docs.meshery.io/concepts/architecture/meshsync) is Meshery’s way of **watching your cluster in real time**.

It keeps Meshery updated with the actual status of Kubernetes resources, whether:
- They were deployed through Meshery (**greenfield**)
- Or already existed before Meshery connected to the cluster (**brownfield**)

#### Why it matters:
- You always see the **current state** in Meshery UI
- Meshery can **react to changes** made outside its own workflows
- You can track what’s registered vs what’s actively running

#### What MeshSync does:
- **GVK Tracking**  
  Meshery uses GVK (Group-Version-Kind) to uniquely identify Kubernetes objects across the cluster.
- **Prometheus Integration**  
  MeshSync can pull health signals from Prometheus to understand if components are behaving as expected.
- **State Awareness (`REGISTERED` vs `CONNECTED`)**  
  Meshery tracks both what it knows (registered models) and what it currently sees running (connected resources).

MeshSync is built on an **event-driven architecture**, using **NATS** for scalable, high-speed communication between components. This enables it to detect infrastructure changes almost instantly — a key capability for multi-cloud or dynamic environments.

### 2.2 Versioned Entity Lifecycle

In Meshery, every component — from a service to a CRD — follows a **versioned lifecycle**.  
This helps ensure changes are:

- Traceable
- Repeatable
- Compatible across environments

[![Meshery Entities]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-entities.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-entities.svg)

#### Understanding Meshery’s Versioning Model

Each Meshery entity follows a structured versioning lifecycle, ensuring consistency before deployment and management. This process enables users to track changes, apply upgrades, and maintain compatibility across environments.

- **Schema** (`v1alpha1`): Defines the fundamental structure of an entity, ensuring consistency and validation.  
- **Definition**: Establishes the expected behavior of an entity, detailing constraints and relationships.  
- **Declaration**: Provides a human-readable YAML-based configuration that Meshery uses for deployment.  
- **Instance**: Represents the actual deployed resource, tracking runtime changes and operational metadata.  

## 3. Core Operation Logic

Once a model is ready, Meshery takes care of deploying it into your Kubernetes cluster and tracking its lifecycle.

This section breaks down what happens during a **Deploy** or **Undeploy** action — how Meshery applies changes, handles dependencies, and updates its internal state.

### 3.1 Deploy Action Breakdown

Meshery follows a consistent workflow when you trigger a deployment. The goal is to ensure the process is safe, dependency-aware, and fully observable.

[![Deploy Action Flow]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-deploy-action.svg)]({{ site.baseurl }}/assets/img/infrastructure-management/meshery-deploy-action.svg)

#### What happens during deployment:

1. **Parse Design File**  
   Meshery reads the `meshery-design.yaml` file, which describes what should be deployed and how it’s structured.
2. **Validate Dependencies**  
   Meshery checks for required CRDs, operators, and supporting components before starting. If something’s missing, it may install it or alert you.
3. **Deploy via Adapter**  
   Each deployment request is passed to the appropriate Meshery Adapter (e.g., Istio, Linkerd) via gRPC. The Adapter translates the model into actual Kubernetes resources.
4. **Sync with MeshSync**  
   Meshery continuously tracks the deployed resources via **MeshSync** to monitor their status in real time.

#### Error Handling and Recovery
- If a dependency is missing, Meshery will either attempt to install it or notify the user.
- Failures during deployment are logged and surfaced through the UI and CLI.

### 3.2 Undeploy Action Breakdown & Dependency Cleanup

Just like deployment, Meshery handles undeployment in a structured way — but with some **differences depending on how the resources were originally deployed**.

#### What happens during undeployment:

1. **Remove Workloads**  
   Deletes core resources like `Pods`, `Deployments`, `StatefulSets`, etc.
2. **Remove Services**  
   Removes network-facing resources like `Services`, `Ingress`, `VirtualService`.
3. **Remove CRDs & Operators (optional)**  
   - If deployed using **Helm**, Meshery calls Helm’s uninstall mechanism, which usually cleans up all dependencies.
   - If deployed from **Git**, Meshery cannot track Helm charts embedded in Git repos, which can leave **CRDs orphaned**.

#### Helm vs Git: Undeploy Behavior

The behavior of undeployment in Meshery depends on **how the resources were originally deployed**:

| Source | Dependency Tracking | What Happens During Undeploy |
|--------|----------------------|-------------------------------|
| **Helm** | Automatic | Meshery calls Helm's uninstall process. All related resources, including CRDs, are typically cleaned up. |
| **Git** | Partial / Manual | Meshery cannot track Helm charts stored in Git. Some resources (especially CRDs) may be left behind. |

#### ⚠️ Limitation: Git-based Deployments

When you deploy using Git repositories, Meshery does **not** record the original Helm chart metadata. As a result:

- CRDs and other dependencies may **not be automatically deleted**
- You may end up with **orphaned resources** in your cluster
- This can cause **conflicts** in future deployments

> Tip: Use Git-based deployments when you want to version your design files, but be aware that Meshery doesn’t track internal Helm structure from those files.

#### How to Clean Up Manually

If you deployed using Git and want to ensure a clean undeployment, remove leftover CRDs and namespaces manually:

```bash
kubectl delete crd <your-crd>
kubectl delete namespace <your-namespace>
```
