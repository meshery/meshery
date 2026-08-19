---
title: Creating a Meshery Design
description: Learn how to create a Meshery Design using the built-in Design Configurator in Meshery UI or the mesheryctl CLI.
categories: [configuration]
aliases:
- /guides/configuration-management/working-with-designs
- /tasks/patterns
---

A Meshery Design is the primary unit of configuration management in Meshery. It is a declarative document that describes the desired state of your infrastructure, cloud architectures, and applications — the components you want, their configuration, and their relationships. Designs can be deployed, shared, versioned, exported, and imported.

See [Meshery Designs Concept]({{< ref "concepts/logical/designs.md" >}}) for detailed information on design features and classifications.

---

## Design Types Overview

Before creating a design, decide on the appropriate design type for your use case:

1. **Annotation-Only Designs**: Intended for visual architecture diagrams, conceptual system maps, boundary groups, and documentation. Nodes carry `isAnnotation: true` and are non-deployable.
2. **Configurable and Deployable Designs**: Composed of model-backed infrastructure components (such as Kubernetes workloads, Helm charts, or cloud CRDs) with configurable properties that can be deployed to live environments.

---

## Ways to Create a Design

{{< tabs id="creating-a-meshery-design-tabs" >}}
Meshery UI (Visual Designer) | fa fa-desktop

### Option 1 — Creating an Annotation-Only Design (Visual Diagrams)

Annotation-only designs allow you to draft system architecture diagrams, document data flows, and illustrate cloud topologies.

```
┌────────────────────────────────────────────────────────┐
│ Visual Designer / Kanvas                               │
│                                                        │
│  ┌──────────────┐      HTTP / TLS      ┌────────────┐  │
│  │ API Gateway  ├─────────────────────►│ Microservice│ │
│  │ (Annotation) │                      │(Annotation)│  │
│  └──────┬───────┘                      └─────┬──────┘  │
│         │                                    │         │
│         ▼                                    ▼         │
│  ┌──────────────────────────────────────────────────┐  │
│  │               VPC Private Subnet                 │  │
│  │                  (Boundary)                      │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
```

#### Step 1: Open the Visual Designer
1. Log in to Meshery UI and navigate to **Designs** (or **Kanvas**) from the left sidebar.
2. Click **+ New Design** to open an empty canvas.

#### Step 2: Name and Tag the Design
1. Provide a descriptive name in the **Design Name** field (e.g., `3-Tier Microservices Architecture`).
2. Add relevant metadata, such as architecture tags (`Architecture`, `Diagram`, `AWS`, `Kubernetes`).

#### Step 3: Add Annotation Elements & Shapes
1. In the component palette on the left, switch to the **Annotations** tab or category.
2. Select shapes, boundaries, text blocks, or cloud service icons (e.g., Cloud Provider logos, firewall boundaries, databases).
3. Drag and drop the annotation shapes onto the canvas.

#### Step 4: Define Visual Relationships and Flow
1. Connect components by dragging connection lines between node anchor points.
2. Select the relationship line to style it (e.g., dashed for logical dependencies, solid for network traffic).
3. Group related components inside boundary containers (e.g., *VPC Boundary*, *Kubernetes Cluster Boundary*).

#### Step 5: Save the Annotation Design
1. Click the **Save** icon in the top toolbar.
2. The design is saved to your workspace and is ready for export or publishing to the [Meshery Catalog]({{< ref "concepts/architecture/catalog/index.md" >}}).

---

### Option 2 — Creating a Configurable & Deployable Design

Deployable designs define operational infrastructure workloads backed by registered Meshery Models.

#### Step 1: Open the Design Configurator
1. Navigate to **Designs** and click **+ New Design**.
2. Enter a unique name for your deployable design.

#### Step 2: Add Components from Registered Models
1. In the **Category** dropdown, select the target category (e.g., `Orchestration & Management`, `Cloud Native Network`, `Database`).
2. In the **Model** dropdown, choose the model (e.g., `Kubernetes`, `Istio`, `AWS`).
3. Click or drag components (e.g., `Deployment`, `Service`, `ConfigMap`, `VirtualService`) onto the canvas.

#### Step 3: Configure Component Properties
1. Click any component on the canvas to open its configuration panel.
2. The form displays dynamic fields generated from the component's JSON Schema:
   - **Required properties** (highlighted)
   - **Environment-specific values** (e.g., `replicas`, `image`, `ports`, `environment variables`)
   - **Selectors and Labels**
3. Configure each component according to your target deployment specification.

#### Step 4: Declare Relationships
1. Connect dependent components (e.g., link a `Service` to a `Deployment`, or a `ConfigMap` to a `Pod` template).
2. Meshery automatically evaluates relationship rules (hierarchical, edge, or sibling relationships) based on defined policies.

#### Step 5: Validate and Save
1. Click the **Validate** button (or check icon) to verify the design against schema definitions.
2. Click **Save** to persist the design.

---

### Option 3 — Direct YAML Authoring (Code Editor)

The Design Configurator includes a bidirectional **Code Editor**:

1. In the design toolbar, switch to the **YAML / Code** view.
2. Author or paste your declarative design manifest following the [Meshery Schemas](https://github.com/meshery/schemas) specification.
3. Switch back to the visual canvas at any time to see the visual layout updated automatically.
4. Click **Save** to persist changes.

<!-- tab -->
mesheryctl (CLI) | fa fa-terminal

## Creating and Managing Designs via CLI

You can create, import, validate, and deploy designs directly using `mesheryctl`.

### 1. Import a Design from a Local File or URL
```bash
# Import a design from a local YAML file
mesheryctl design import -f ./my-design.yaml

# Import a design from a remote URL
mesheryctl design import -f https://raw.githubusercontent.com/service-mesh-patterns/service-mesh-patterns/master/samples/istio-bookinfo.yaml

# Import a Kubernetes manifest or Helm chart
mesheryctl design import -f ./deployment.yaml -s manifest
mesheryctl design import -f ./chart.tgz -s helm
```

### 2. View and List Saved Designs
```bash
# List all designs in your workspace
mesheryctl design list

# View details of a specific design
mesheryctl design view <Design-Name-or-ID>
```

### 3. Deploy (Apply) a Design
```bash
# Apply a design file directly to the connected cluster
mesheryctl design apply -f ./my-design.yaml

# Apply an already-saved design by name
mesheryctl design apply "MyDesignName"
```

See the [`mesheryctl design` reference]({{< ref "reference/references/mesheryctl/design/_index.md" >}}) for full CLI options.
{{< /tabs >}}

---

## Working Design Examples

### Example 1: Annotation-Only Architecture Design (YAML)

The following example demonstrates an **Annotation-only** architecture design for a distributed service topology. Note the `isAnnotation: true` metadata flag marking components as non-deployable visual nodes.

```yaml
name: Sample-Annotation-Architecture
version: 0.0.1
schemaVersion: designs.meshery.io/v1beta1
components:
  - name: api-gateway-node
    type: Gateway
    model:
      name: annotations
      version: 0.0.1
    metadata:
      isAnnotation: true
      annotations:
        meshery.io/component-type: "annotation"
        meshery.io/shape: "round-rectangle"
        meshery.io/label: "API Gateway (Public Ingress)"
        meshery.io/position-x: 120
        meshery.io/position-y: 150
    configuration:
      description: "Public edge entry point for API clients"
      protocol: "HTTPS / TLS 1.3"
  - name: auth-service-node
    type: Service
    model:
      name: annotations
      version: 0.0.1
    metadata:
      isAnnotation: true
      annotations:
        meshery.io/component-type: "annotation"
        meshery.io/shape: "rectangle"
        meshery.io/label: "Auth Service"
        meshery.io/position-x: 350
        meshery.io/position-y: 150
    configuration:
      description: "JWT authentication and authorization service"
relationships:
  - kind: edge
    type: visual
    metadata:
      isAnnotation: true
    selectors:
      - allow:
          from:
            id: api-gateway-node
          to:
            id: auth-service-node
```

---

### Example 2: Deployable Kubernetes Service Stack (YAML)

The following example demonstrates a functional, **Deployable** design defining an NGINX web application with a Kubernetes `Deployment` and `Service`.

```yaml
name: nginx-web-stack
version: 0.0.1
schemaVersion: designs.meshery.io/v1beta1
components:
  - name: web-deployment
    type: Deployment
    apiVersion: apps/v1
    model:
      name: kubernetes
      version: v1.30.0
    configuration:
      spec:
        replicas: 2
        selector:
          matchLabels:
            app: web
        template:
          metadata:
            labels:
              app: web
          spec:
            containers:
              - name: nginx
                image: nginx:1.25-alpine
                ports:
                  - containerPort: 80
                    name: http
  - name: web-service
    type: Service
    apiVersion: v1
    model:
      name: kubernetes
      version: v1.30.0
    configuration:
      spec:
        type: ClusterIP
        selector:
          app: web
        ports:
          - port: 80
            targetPort: 80
            name: http
relationships:
  - kind: edge
    type: network
    selectors:
      - allow:
          from:
            id: web-service
          to:
            id: web-deployment
```

---

## Next Steps

- **[Exporting Designs]({{< ref "guides/configuration-management/import-export-designs.md" >}})**: Export your design as YAML, OCI image, or Helm chart.
- **[Publishing to Meshery Catalog]({{< ref "concepts/architecture/catalog/index.md" >}})**: Share your designs with the global community.
- **[Deploying a Design]({{< ref "guides/configuration-management/working-with-designs/index.md" >}})**: Deploy and manage designs across multi-cluster environments.

