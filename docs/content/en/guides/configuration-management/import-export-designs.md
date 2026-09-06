---
title: Importing and Exporting Designs
description: Manage and transfer Meshery designs seamlessly through import and export using CLI and UI.
aliases:
- /guides/configuration-management/importing-designs
- /guides/configuration-management/import-export-a-design/import-export-designs
- /extensions/importing-a-design
- /extensions/importing-an-application
display_title: false
---

# Importing and Exporting Designs

Meshery supports two primary operations for working with designs: **Import** and **Export**.
A **[design]({{< ref "concepts/logical/designs.md" >}})** in Meshery is a structured model describing how various components (e.g., Kubernetes manifests, Helm charts, or Docker Compose objects) should be managed and deployed. This document details the core concepts, the import/export processes, and the available methods to perform them.

As a refresher, [designs]({{< ref "concepts/logical/designs.md" >}}) consist of:

- **Components**: Defined resources or services (e.g., Deployments, Services).
- **Relationships**: Possible connections or dependencies among components (e.g., hierarchical, sibling, edge).

When **importing** a design, Meshery performs the following steps:

1. **Parses** and **validates** the design file, processing design structure, applying transformations as needed.

<!-- 
Meshery does this when deploying, not when importing.
1. **Routes** validated components to the appropriate Meshery Adapters for further deployment or processing. -->

When **exporting** a design, Meshery:

1. **Fetches** the design data from the Meshery database using `/api/pattern/download/{id}`.
2. **Converts** the design format if necessary (e.g., from V1alpha2 to V1beta1).
3. **Applies** requested export format (e.g., `current`, `original`, `oci`, `helm-chart`).
4. **Sends** the exported design file for download as a YAML file or OCI artifact or Helm chart(.tgz).

## Import Methods

You can import designs into Meshery through either **[Meshery CLI]({{< ref "reference/references/mesheryctl/_index.md" >}})** or **Meshery UI**.

### 1. Using Meshery CLI

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
   A POST request is sent to `/api/pattern/import`, containing a design file. If the design is sourced from Kubernetes Manifests, Docker Compose, or Helm Charts, it is first converted into a standard Kubernetes Manifest.

2. **Pattern Engine queries registry**
   Meshery queries the component registry to match the design’s `kind`, `apiVersion`, and other identifiers. If the design originates from non-Meshery formats, Meshery attempts to transform it accordingly.

3. **System converts design format**
   - Converts Helm Charts, Docker Compose, and Kubernetes Manifests into Meshery Design.
   - Removes unnecessary fields and ensures compatibility.

4. **Validated components are stored in Meshery**
   Once validated, the design is stored in Meshery. Users can later deploy it to a supported platform (e.g., Kubernetes, Consul, Istio) using Meshery UI or CLI.

## Export Methods & Formats

Designs can be exported either via **Meshery UI** or **Meshery CLI (`mesheryctl`)**. Exporting allows you to back up designs, distribute them as artifacts across registries, publish them, or package them into standard tools like Helm.

### Supported Export Formats

Meshery supports multiple export targets tailored to different lifecycle stages:

| Format | Target Extension | Best For | Description |
| :--- | :--- | :--- | :--- |
| **Meshery Design (YAML)** | `.yaml` / `.yml` | GitOps & version control | Exports the full declarative Meshery Design manifest conforming to `designs.meshery.io/v1beta1`. |
| **OCI Artifact** | `.tar` / OCI layer | Container registries & Artifact Hub | Packages the design into an Open Container Initiative (OCI) image artifact for registry distribution. |
| **Helm Chart** | `.tgz` | Helm CLI & Kubernetes packaging | Converts and packages deployable components into a standard, deployable Helm Chart archive. |

---

### 1. Exporting via Meshery UI

To export a design from the web interface:

1. Navigate to **Designs** (or **Kanvas**) from the left navigation menu.
2. In the designs list or grid, locate the design you wish to export.
3. Click the **Action (three dots / menu)** or **Export** button on the design card.
4. Select your desired export format:
   - **Current**: Exports the active version of the design as a Meshery YAML manifest.
   - **Original**: Exports the original imported source manifest (if imported from Helm or Kubernetes).
   - **OCI**: Packages the design as an OCI artifact archive.
   - **Helm Chart**: Generates a standard Helm chart archive (`.tgz`).
5. Choose your local destination directory to download the exported file.

---

### 2. Exporting via Meshery CLI (`mesheryctl`)

Use the `mesheryctl design export` command to export designs programmatically:

```bash
# General syntax
mesheryctl design export <design-name-or-id> [--type <FORMAT>] [--output <DIRECTORY>]
```

Supported `--type` values:

| Value | Description |
| :--- | :--- |
| `current` | *(default)* Exports the active design version as a Meshery YAML manifest (`.yaml`). |
| `oci` | Packages and exports the design as an OCI artifact archive. |
| `original` | Exports the original imported source artifact (e.g., the original Helm chart or Kubernetes manifest, when available). |

#### Examples:

```bash
# Export design as a Meshery YAML manifest (default format)
mesheryctl design export "MyDesign" --type current -o ./exports

# Export design as an OCI artifact
mesheryctl design export "MyDesign" --type oci -o ./exports

# Export the original imported source (e.g., the original Helm chart or manifest)
mesheryctl design export "MyDesign" --type original -o ./exports
```

---

## Summary

- **Import**: Meshery parses and validates your design, converts non-Meshery formats (e.g., Kubernetes Manifests, Helm Charts, Docker Compose) into Meshery Designs, and stores the design in your workspace.
- **Export**: Meshery retrieves the design from storage, converts its format if needed, packages it in the selected format (`current` for YAML, `oci` for OCI artifact, or `original` for the original source artifact), and delivers it via CLI or UI download for GitOps, container registries, or package managers.

