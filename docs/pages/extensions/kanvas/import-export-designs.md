---
layout: default
title: Import and Export Designs
abstract: Manage and transfer Meshery designs seamlessly through import and export using CLI and UI.
permalink: extensions/import-export-designs
redirect_from:
   - /guides/configuration-management/importing-designs
   - /extensions/importing-a-design
   - /extensions/importing-an-application
list: include
display-title: 'false'
category: kanvas
type: extensions
language: en
---

# Import and Export Designs

Meshery supports two primary operations for working with designs: **Import** and **Export**. A **[design](https://docs.meshery.io/concepts/logical/designs)** in Meshery is a structured model
describing how various components (e.g., Kubernetes manifests, Helm charts, or Docker Compose objects) should be managed and deployed. This document details the core concepts, the import/export
processes, and the available methods to perform them.

## Core Operations

Meshery treats each design as a combination of:

-  **Components**: Defined resources or services (e.g., Deployments, Services).
-  **Relationships**: Possible connections or dependencies among components (e.g., hierarchical, sibling, edge).

When **importing** a design, Meshery:

1. **Parses** and **validates** the design file.
2. Processes design structure, applying transformations as needed.
3. **Routes** validated components to the appropriate Meshery Adapters for further deployment or processing.

When **exporting** a design, Meshery:

1. **Fetches** the design data from the Meshery database using `/api/pattern/download/{id}`.
2. **Converts** the design format if necessary (e.g., from V1alpha2 to V1beta1).
3. **Applies** requested export format (e.g., `current`, `original`, `oci`).
4. **Sends** the exported design file for download as a YAML file or OCI artifact.

## Import Methods

You can import designs into Meshery through either **[Meshery CLI](https://docs.meshery.io/reference/mesheryctl)** or **Meshery UI**.

#### 1. Using Meshery CLI

```bash
mesheryctl design import --file <PATH/URL> --source-type <TYPE>
```

```bash
# Example: Import Kubernetes manifest
mesheryctl design import -f ./app-deployment.yaml -s manifest
```

#### 2. Using Meshery UI

1. Go to **Designs > Import Design**.
2. Choose your input method.
3. Select the source type from the dropdown (e.g., Kubernetes manifest, Helm chart).
4. Submit the design for processing.

### Import Processing Details

After you initiate an import, Meshery executes a **dataflow sequence** to properly handle the design:

1. **Client submits a design**  
   A POST request is sent to `/api/pattern/import`, containing a design file. If the design is sourced from Kubernetes Manifests, Docker Compose, or Helm Charts, it is first converted into a standard
   Kubernetes Manifest.

2. **Pattern Engine queries registry**  
   Meshery queries the component registry to match the design’s `kind`, `apiVersion`, and other identifiers. If the design originates from non-Meshery formats, Meshery attempts to transform it
   accordingly.

3. **System converts design format**

   -  Converts Helm Charts, Docker Compose, and Kubernetes Manifests into Meshery Design.
   -  Removes unnecessary fields and ensures compatibility.

4. **Validated components are stored in Meshery**  
   Once validated, the design is stored in Meshery. Users can later deploy it to a supported platform (e.g., Kubernetes, Consul, Istio) using Meshery UI or CLI.

## Export Methods

Designs can also be exported either via **Meshery CLI** or **Meshery UI**. During export, Meshery gathers the design’s components, relationships, and metadata into a specified format.

#### 1. Using Meshery CLI

```bash
mesheryctl design export --type <FORMAT> --output <DIRECTORY>
```

```bash
# Example: Export to OCI format
mesheryctl design export --type oci -o ./exports
```

#### 2. Using Meshery UI

1. Go to **Designs > My Designs** and select the design you want to export.
2. Click the **Export** action button.
3. Choose an export format (e.g., `current`, `original`, `oci`, `helm`).
4. Specify the download location for your exported file.

### Supported Export Formats

1. **Design File**: Meshery’s native YAML format. Preserves all Meshery-specific details for backup, sharing, or re-import.
2. **Kubernetes Manifest**: Standard Kubernetes YAML for direct cluster deployment.
3. **OCI Image**: Packages your design as a container image for sharing via registries.
4. **Helm Chart**: Packages your design as a Helm chart (`.tgz`) for Helm-based deployments.
5. **Embed Design**: Generates an embeddable version for docs or web pages.

**Behavioral Note:**  
Both **Kubernetes Manifest** and **Helm Chart** exports are **lossy**: Meshery-specific information (like visual layout, comments, and metadata) is removed. Only core Kubernetes resources are included
in these formats. Use the **Design File** export to retain all Meshery details.

## Summary

-  **Import**: Meshery parses and validates your design, converts non-Meshery formats (e.g., Kubernetes Manifests, Helm Charts, Docker Compose) into Meshery Designs, and stores the design in Meshery.
   Users can later deploy it to the correct adapters as needed.
-  **Export**: Meshery **retrieves** the design from storage, **converts** its format if needed, **packages** it in the selected format (`current`, `original`, `oci`,`helm`), and **delivers** it as a
   YAML file or OCI artifact via CLI or UI download.
