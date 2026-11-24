---
title: Push or Pull a Model Image
categories: [configuration]
description: Push or pull a model image to or from an OCI-compatible image repository.
---

## Use mesheryctl to Push or Pull a Model Image

A Meshery [Model]({{< ref "concepts/logical/models/index.md" >}}) can be packaged as an
**OCI image** - a portable, versioned artifact that bundles the model definition together with
all of its components and relationships. Because the package is OCI-compliant, it can be moved
between Meshery Servers and stored in any OCI-compatible image repository (Docker Hub, GitHub
Container Registry, AWS ECR, and so on), exactly like a container image.

### Package (build) a model image

Package model files that follow the [scaffold layout]({{< ref "guides/configuration-management/creating-models/index.md" >}})
(`[model-name]/[model-version]`) into an OCI artifact:

```bash
mesheryctl model build [model-name]/[model-version]
```

`model build` writes an OCI artifact named `<model-name>-<version>.tar` to your current
directory (for example, `mesheryctl model build digitalocean-icons/v0.1.0` produces
`digitalocean-icons-v0-1-0.tar`). This command is local - it does not require a running Meshery
Server.

### Push a model image into a Meshery Server

"Pushing" a model into a Meshery Server means registering it into that server's
[Registry]({{< ref "concepts/logical/registry.md" >}}). Import the packaged artifact (or the
model directory) with:

```bash
mesheryctl model import -f digitalocean-icons-v0-1-0.tar
```

`mesheryctl model import` talks to the Meshery Server of your current context, so you must be
logged in (see [`mesheryctl system login`]({{< ref "reference/mesheryctl/system/login.md" >}})).
The file can be an OCI `.tar`, a compressed `tar.gz`, a model directory, or a remote URL - see
[Importing Models]({{< ref "guides/configuration-management/importing-models/index.md" >}}) for
the full list of supported sources.

### Pull a model image from a Meshery Server

"Pulling" a model retrieves an already-registered model from a Meshery Server as an OCI image
(or a compressed archive) so it can be inspected, versioned, or moved to another server:

```bash
mesheryctl model export [model-name] -o oci
```

Use `-o` to choose the output format (`oci` or `tar`), `-t` for the file type of the definitions
inside the archive (`yaml`, the default, or `json`), and `-l` to choose the output directory
(the directory must already exist). See
[Exporting Models]({{< ref "guides/configuration-management/exporting-models/index.md" >}}) for
details.

### Push to or pull from an external OCI registry

Because the artifact produced by `mesheryctl model build` or `model export` is a standard OCI
image, you can push it to - or pull it from - any OCI-compatible registry using standard container
and artifact tooling, such as [ORAS (OCI Registry As Storage)](https://oras.land/) or `docker`. Once
pulled, you can import the model into any Meshery Server using `mesheryctl model import`.

This makes model images easy to distribute, version, and share across teams and environments
alongside your other OCI artifacts.

#### Docker Hub

Authenticate using your Docker ID and a Personal Access Token (PAT):

```bash
oras login docker.io -u <username> -p <personal-access-token>
```

Push the model artifact:

```bash
oras push docker.io/<username>/<model-name>:<version> digitalocean-icons-v0-1-0.tar
```

Pull the model artifact:

```bash
oras pull docker.io/<username>/<model-name>:<version>
```

#### Azure Container Registry (ACR)

Authenticate to your registry using the Azure CLI or ORAS with credentials:

Using Azure CLI:

```bash
az acr login --name <registry-name>
```

Using ORAS with service principal or repository credentials:

```bash
oras login <registry-name>.azurecr.io -u <username-or-sp-app-id> -p <password-or-sp-secret>
```

Push the model artifact:

```bash
oras push <registry-name>.azurecr.io/<model-name>:<version> digitalocean-icons-v0-1-0.tar
```

Pull the model artifact:

```bash
oras pull <registry-name>.azurecr.io/<model-name>:<version>
```

#### GitHub Container Registry (GHCR)

Authenticate using your GitHub username and a Personal Access Token (classic) with `write:packages` and `read:packages` scopes:

```bash
echo $CR_PAT | oras login ghcr.io -u <username> --password-stdin
```

Push the model artifact:

```bash
oras push ghcr.io/<owner>/<model-name>:<version> digitalocean-icons-v0-1-0.tar
```

Pull the model artifact:

```bash
oras pull ghcr.io/<owner>/<model-name>:<version>
```

#### AWS Elastic Container Registry (ECR)

Authenticate by obtaining an authorization token via the AWS CLI and piping it to ORAS:

```bash
aws ecr get-login-password --region <region> | oras login --username AWS --password-stdin <aws-account-id>.dkr.ecr.<region>.amazonaws.com
```

> **Note on ECR Repository Creation:**
> AWS ECR requires the target repository (`<model-name>`) to exist before pushing. You can create it with:
> `aws ecr create-repository --repository-name <model-name> --region <region>`

Push the model artifact:

```bash
oras push <aws-account-id>.dkr.ecr.<region>.amazonaws.com/<model-name>:<version> digitalocean-icons-v0-1-0.tar
```

Pull the model artifact:

```bash
oras pull <aws-account-id>.dkr.ecr.<region>.amazonaws.com/<model-name>:<version>
```

#### Import the pulled model into Meshery Server

Once you have pulled the model artifact (`.tar` archive) from your OCI registry, register it with your active Meshery Server:

```bash
mesheryctl model import -f <model-name>-<version>.tar
```
