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

- [mesheryctl](/installation/mesheryctl)
- [ORAS CLI](https://oras.land/docs/installation) (optional, for pushing model images as OCI artifacts)
- [Docker CLI](https://docs.docker.com/get-docker/) (optional, for Docker Hub authentication)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html) (optional, for AWS ECR authentication)
- Credentials for your target registry

## Pull a Model Image

Use `mesheryctl model import` to pull and import a model from any OCI registry:

{% highlight bash %}
mesheryctl model import -f [OCI]
{% endhighlight %}

**Example output (success):**

{% highlight text %}
Model imported successfully: [model-name] v[version]
{% endhighlight %}

**Example output (failure):**

{% highlight text %}
Error: failed to pull model image: unauthorized: authentication required
{% endhighlight %}

## Push a Model Image

Meshery model images are OCI artifacts. Use the [ORAS CLI](https://oras.land) to push them to any OCI-compatible registry:

{% highlight bash %}
oras push [registry-host]/[repository]:[tag] [model-file]
{% endhighlight %}

**Example output (success):**

{% highlight text %}
Pushed [registry-host]/[repository]:[tag]
Digest: sha256:abc123...
{% endhighlight %}

## Registry-Specific Examples

### Docker Hub

**Login:**

{% highlight bash %}
docker login
{% endhighlight %}

**Push with ORAS:**

{% highlight bash %}
oras push docker.io/[your-username]/[model-name]:[version] [model-file]
{% endhighlight %}

**Pull with mesheryctl:**

{% highlight bash %}
mesheryctl model import -f oci://docker.io/[your-username]/[model-name]:[version]
{% endhighlight %}

---

### Azure Container Registry (ACR)

**Login:**

{% highlight bash %}
az acr login --name [your-acr-name]
{% endhighlight %}

**Push with ORAS:**

{% highlight bash %}
oras push [your-acr-name].azurecr.io/[model-name]:[version] [model-file]
{% endhighlight %}

**Pull with mesheryctl:**

{% highlight bash %}
mesheryctl model import -f oci://[your-acr-name].azurecr.io/[model-name]:[version]
{% endhighlight %}

> **Note:** Ensure your Azure service principal has `AcrPush`/`AcrPull` role assigned.

---

### AWS Elastic Container Registry (ECR)

**Login:**

{% highlight bash %}
aws ecr get-login-password --region [region] | docker login --username AWS --password-stdin [account-id].dkr.ecr.[region].amazonaws.com
{% endhighlight %}

**Push with ORAS:**

{% highlight bash %}
oras push [account-id].dkr.ecr.[region].amazonaws.com/[model-name]:[version] [model-file]
{% endhighlight %}

**Pull with mesheryctl:**

{% highlight bash %}
mesheryctl model import -f oci://[account-id].dkr.ecr.[region].amazonaws.com/[model-name]:[version]
{% endhighlight %}

**Create the repository (if it doesn't exist):**

{% highlight bash %}
aws ecr create-repository --repository-name [model-name]
{% endhighlight %}

> **Note:** ECR requires the repository to exist before pushing.

---

### GitHub Container Registry (GHCR)

**Login:**

{% highlight bash %}
echo $GITHUB_TOKEN | oras login ghcr.io -u [github-username] --password-stdin
{% endhighlight %}

**Push with ORAS:**

{% highlight bash %}
oras push ghcr.io/[github-username]/[model-name]:[version] [model-file]
{% endhighlight %}

**Pull with mesheryctl:**

{% highlight bash %}
mesheryctl model import -f oci://ghcr.io/[github-username]/[model-name]:[version]
{% endhighlight %}

> **Note:** Your `GITHUB_TOKEN` needs `write:packages` scope for push and `read:packages` for pull.

---

## See Also

- [mesheryctl model import](/reference/mesheryctl/model/import)
- [ORAS CLI documentation](https://oras.land/docs/)
- [Managing Models](/guides/configuration-management/creating-models)
