---
layout: default
title: Push or Pull a Model Image
permalink: guides/configuration-management/push-pull-model-image
type: guides
category: configuration
language: en
list: include
abstract: "Push or pull a model image to or from an OCI-compatible image repository."
---

## Use mesheryctl to Push or Pull a Model Image

You can push or pull Meshery model images to or from any OCI-compatible image repository.

## Prerequisites

- [mesheryctl installed](/installation/mesheryctl)
- [ORAS CLI installed](https://oras.land/docs/installation) (for pushing model images)
- Credentials for your target registry

## Pull a Model Image

Use `mesheryctl model import` to pull and import a model from any OCI registry:
```bash
mesheryctl model import -f [OCI]
```

## Push a Model Image

Meshery model images are OCI artifacts. Use the [ORAS CLI](https://oras.land) to push them to any OCI-compatible registry:
```bash
oras push [registry-host]/[repository]:[tag] [model-file]
```

## Registry-Specific Examples

### Docker Hub

**Login:**
```bash
docker login
```

**Push with ORAS:**
```bash
oras push docker.io/[your-username]/[model-name]:[version] [model-file]
```

**Pull with mesheryctl:**
```bash
mesheryctl model import -f oci://docker.io/[your-username]/[model-name]:[version]
```

---

### Azure Container Registry (ACR)

**Login:**
```bash
az acr login --name [your-acr-name]
```

**Push with ORAS:**
```bash
oras push [your-acr-name].azurecr.io/[model-name]:[version] [model-file]
```

**Pull with mesheryctl:**
```bash
mesheryctl model import -f oci://[your-acr-name].azurecr.io/[model-name]:[version]
```

> **Note:** Ensure your Azure service principal has `AcrPush`/`AcrPull` role assigned.

---

### AWS Elastic Container Registry (ECR)

**Login:**
```bash
aws ecr get-login-password --region [region] | docker login \
  --username AWS \
  --password-stdin [account-id].dkr.ecr.[region].amazonaws.com
```

**Push with ORAS:**
```bash
oras push [account-id].dkr.ecr.[region].amazonaws.com/[model-name]:[version] [model-file]
```

**Pull with mesheryctl:**
```bash
mesheryctl model import -f oci://[account-id].dkr.ecr.[region].amazonaws.com/[model-name]:[version]
```

> **Note:** ECR requires the repository to exist before pushing. Create it first:
> ```bash
> aws ecr create-repository --repository-name [model-name]
> ```

---

### GitHub Container Registry (GHCR)

**Login:**
```bash
echo $GITHUB_TOKEN | oras login ghcr.io -u [github-username] --password-stdin
```

**Push with ORAS:**
```bash
oras push ghcr.io/[github-username]/[model-name]:[version] [model-file]
```

**Pull with mesheryctl:**
```bash
mesheryctl model import -f oci://ghcr.io/[github-username]/[model-name]:[version]
```

> **Note:** Your `GITHUB_TOKEN` needs `write:packages` scope for push and `read:packages` for pull.

---

## See Also

- [mesheryctl model import](/reference/mesheryctl/model/import)
- [ORAS CLI documentation](https://oras.land/docs/)
- [Managing Models](/guides/configuration-management/creating-models)
