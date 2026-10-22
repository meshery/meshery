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

Import a design to bring infrastructure configuration you already have into Meshery. Export a design to back it up, move it to another Meshery Server, or use it outside Meshery.

A [design]({{< ref "concepts/logical/designs.md" >}}) holds **components**, such as Deployments and Services, and the **relationships** between them.

## What happens when you import

You can import any of these formats:

- Meshery Design (YAML, or an OCI image exported from Meshery)
- Kubernetes manifest (YAML)
- Helm chart (`.tgz` archive)
- Docker Compose file
- Kustomize package (archive)

Meshery Server identifies the format from the file itself. It converts Helm charts, Docker Compose files, and Kustomize packages to a Kubernetes manifest, matches each resource's `apiVersion` and `kind` to a component in the registry, builds a design from the matched components, and saves the design. Importing does not deploy anything. Deploy the design afterward from Meshery UI or with `mesheryctl design deploy`.

## Import a design

{{< tabs id="import-design-tabs" >}}
Meshery UI | fa fa-desktop

1. Go to **Configuration** > **Designs** and click **Import Design**.
2. Enter a **Design Name**.
3. Choose an **Upload method**:
    - **File Upload**, then pick the file from your computer.
    - **URL Import**, then enter a direct URL to the file, for example `https://raw.githubusercontent.com/<org>/<repo>/main/app.yaml`. The URL must return the file itself, not a web page.
4. Click **Import**.

<!-- tab -->
mesheryctl | fa fa-terminal

```bash
mesheryctl design import -f <path-or-URL> [-n <design-name>] [-s <source-type>]
```

- `-f` is a local path or a direct URL to the file.
- `-n` names the design. Without it, `mesheryctl` uses the file name.
- `-s` is optional because Meshery Server detects the format. If you pass it, use one of `"Kubernetes Manifest"`, `"Helm Chart"`, `"Docker Compose"`, or `"Design"` (case-insensitive).

For example, to import a Kubernetes manifest:

```bash
mesheryctl design import -f ./app-deployment.yaml -n my-app
```

`mesheryctl` prints the new design's name and ID.
{{< /tabs >}}

### Import a Kustomize package

To import a Kustomize package, upload a `.zip`, `.tar`, `.tar.gz`, or `.tgz` archive with a file named `kustomization.yaml` at the root of the archive. Include the files it references in the same archive. For example:

```text
kustomization.yaml
deployment.yaml
service.yaml
```

You can also import the archive from a direct URL. Meshery downloads the file and takes its name from the response's `Content-Disposition` header or the last segment of the URL, so make sure that name ends in `.zip`, `.tar`, `.tar.gz`, or `.tgz`.

## Export a design

{{< tabs id="export-design-tabs" >}}
Meshery UI | fa fa-desktop

1. Go to **Configuration** > **Designs**.
2. Click **Download** on the design you want to export.
3. Click the download icon next to the format you want:

| Format | What you get |
| --- | --- |
| **Meshery Design (yaml)** | A complete Meshery design, including embedded images. Use it for backups or to move a design to another Meshery Server without losing anything. |
| **Meshery Design (OCI image)** | The design as an OCI image that you can push to and pull from a container registry. |
| **Kubernetes Manifest (yaml)** | The design's Kubernetes resources as a plain manifest, for use with `kubectl` or other tools. |
| **Helm Chart (tar.gz)** | The design packaged as a Helm chart. |

<!-- tab -->
mesheryctl | fa fa-terminal

```bash
mesheryctl design export <design-name-or-ID> [--type <type>] [-o <directory>]
```

- `--type` is one of:
    - `current` (default): the design as Meshery stores it now.
    - `original`: the file you originally imported, such as the Helm chart or manifest.
    - `oci`: the design as an OCI image.
- `-o` is the directory to save the file in. The default is the current directory.

For example, to export a design as an OCI image into `./exports`:

```bash
mesheryctl design export my-app --type oci -o ./exports
```

To export a design as a Helm chart or a Kubernetes manifest, use Meshery UI.
{{< /tabs >}}

## See also

- [Creating a Meshery Design]({{< ref "guides/configuration-management/creating-a-meshery-design.md" >}})
- [`mesheryctl design` reference]({{< ref "reference/references/mesheryctl/design/_index.md" >}})
