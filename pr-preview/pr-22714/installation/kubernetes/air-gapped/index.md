# Air-Gapped Installation

> Install Meshery on Kubernetes in an air-gapped (offline) environment using Helm, a private registry, or side-loaded images.

---

LLMS index: [llms.txt](/pr-preview/pr-22714/llms.txt)

---

# Install Meshery on Kubernetes in an Air-Gapped Environment

An air-gapped cluster cannot pull images or charts from the internet. Installing Meshery there takes three steps, in order:

1. Collect every container image the installation needs, on a machine with internet access.
2. Move the images and the Helm chart into the air-gapped environment, either by mirroring the images to a private registry the cluster can reach or by side-loading them onto the cluster nodes.
3. Install from the local chart with a values file that points the chart at your local copies.

This guide covers installation on Kubernetes with Helm. It assumes you have already chosen the Meshery release and Meshery Operator release to deploy, and that you pin both. Do not install a moving tag such as `stable-latest` in an air-gapped environment: record the exact tags (or digests) you mirrored, so the installation is reproducible.

## Images Required

The images fall into two groups. Images deployed directly by the Helm charts can be repointed in `values.yaml`. Images deployed by Meshery Operator at runtime cannot (see [Operator-Managed Images](#operator-managed-images)).

### Deployed by the Helm Charts

| Image | Component | Notes |
| --- | --- | --- |
| `meshery/meshery` | Meshery Server | Chart default tag is `stable-latest`; pin the release tag you mirrored. |
| `meshery/meshery-operator:1.0.6` | Meshery Operator | Only when the operator is enabled (`meshery-operator.enabled`). Tag is the operator release. |
| `alpine/k8s:1.35.6` | Operator CRD update job | Runs on install and upgrade when `meshery-operator.crds.updateJob.enabled` is `true` (the default). Easy to miss; a missing copy fails the job, not the operator pod. |
| `meshery/meshery-istio:stable-latest` | Meshery Istio adapter | Only when the adapter is enabled. The adapter charts carry the tag inside their image reference; pin the release tag you mirrored. |
| `meshery/meshery-linkerd:stable-latest` | Meshery Linkerd adapter | Only when enabled. |
| `meshery/meshery-consul:stable-latest` | Meshery Consul adapter | Only when enabled. |
| `meshery/meshery-kuma:stable-latest` | Meshery Kuma adapter | Only when enabled. |
| `meshery/meshery-nsm:stable-latest` | Meshery NSM adapter | Only when enabled. |
| `meshery/meshery-traefik-mesh:stable-latest` | Meshery Traefik Mesh adapter | Only when enabled. |
| `meshery/meshery-cilium:stable-latest` | Meshery Cilium adapter | Only when enabled. |

Adapters are disabled by default, and so is Meshery Operator. Mirror only the optional images for the components you intend to enable.

### Operator-Managed Images

When Meshery Operator runs, it creates the MeshSync and Broker workloads itself. Their images are chosen by the operator, not by the Helm charts:

| Image | Component |
| --- | --- |
| `meshery/meshsync:1.0.3` | MeshSync (tag is the operator's pinned default; see below) |
| `nats:2.14.2-alpine` | Broker (NATS server) |
| `natsio/nats-server-config-reloader:0.23.0` | Broker (NATS configuration reloader sidecar) |

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Verify this list against your own render</div>


Tags change with each operator release. Before mirroring, render the chart with the values you intend to install with and list the images it names, then check the operator release notes for the MeshSync and NATS versions that operator deploys:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm template meshery ./meshery-&lt;chart-version&gt;.tgz --namespace meshery -f air-gapped-values.yaml | grep -E &#39;^\s+image:&#39; | sort -u</code>
	</div>
</pre>

</div>


## Move the Images

### Option A: Mirror to a Private Registry

On a machine with internet access, pull each image, retag it for your registry, and push it:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">registry=registry.example.local
for image in meshery/meshery:&lt;release-tag&gt; meshery/meshery-operator:1.0.6 alpine/k8s:1.35.6 meshery/meshsync:1.0.3 nats:2.14.2-alpine natsio/nats-server-config-reloader:0.23.0; do
  docker pull &#34;$image&#34;
  docker tag &#34;$image&#34; &#34;$registry/$image&#34;
  docker push &#34;$registry/$image&#34;
done</code>
	</div>
</pre>


If your registry requires authentication, create the pull secret in the `meshery` namespace before installing, and name it in your values file (below):



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl create secret docker-registry private-registry --namespace meshery --docker-server=registry.example.local --docker-username=&lt;username&gt; --docker-password=&lt;password&gt;</code>
	</div>
</pre>


### Option B: Side-Load the Images onto Each Node

If you have no registry, archive the images, transfer the archive, and load it on every node:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">docker pull meshery/meshery:&lt;release-tag&gt; meshery/meshery-operator:1.0.6 alpine/k8s:1.35.6 meshery/meshsync:1.0.3 nats:2.14.2-alpine natsio/nats-server-config-reloader:0.23.0
docker save -o meshery-images.tar meshery/meshery:&lt;release-tag&gt; meshery/meshery-operator:1.0.6 alpine/k8s:1.35.6 meshery/meshsync:1.0.3 nats:2.14.2-alpine natsio/nats-server-config-reloader:0.23.0</code>
	</div>
</pre>


On each node, import the archive into the container runtime. For containerd (the runtime most Kubernetes distributions use today):



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">sudo ctr --namespace k8s.io images import meshery-images.tar</code>
	</div>
</pre>


Images loaded this way keep their original names. That matters for the operator-managed images, which the operator requests by their original names (see below).

## Get the Helm Chart

The usual `helm repo add` fetches the chart index from the internet, so transfer the chart itself instead. On the connected machine:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm repo add meshery https://meshery.io/charts/
helm repo update
helm pull meshery/meshery --version &lt;chart-version&gt;</code>
	</div>
</pre>


Transfer the resulting `meshery-<chart-version>.tgz` into the air-gapped environment. Installing from the packaged chart (or from a checkout of the Meshery repository) needs no repository access.

## Point values.yaml at Your Registry

Create `air-gapped-values.yaml`. Two defaults must change for every component: the image reference must name your registry, and the pull policy must be `IfNotPresent`. The chart defaults for Meshery Server and the adapters use the moving `stable-latest` tag with pull policy `Always`, which an offline cluster can neither resolve nor rely on.

```yaml
image:
  repository: registry.example.local/meshery/meshery
  tag: "<release-tag>"
  pullPolicy: IfNotPresent
imagePullSecrets:
  - name: private-registry

meshery-operator:
  enabled: true
  image:
    repository: registry.example.local/meshery/meshery-operator
    # tag remains the operator release, 1.0.6 in this example
    pullPolicy: IfNotPresent
  imagePullSecrets:
    - name: private-registry
  crds:
    updateJob:
      image:
        repository: registry.example.local/alpine/k8s

# Adapters, for each one you enable. Note that an adapter's image
# reference carries its tag, because that is how the adapter
# subcharts define it.
meshery-istio:
  enabled: true
  image:
    repository: registry.example.local/meshery/meshery-istio:<release-tag>
    pullPolicy: IfNotPresent
```

Notes on this file:

- The operator has its own `image`, `imagePullSecrets`, and CRD job image settings, nested under `meshery-operator`. There is no single global registry setting shared by the charts, so each component is repointed separately.
- If you side-loaded the images instead of mirroring them, keep the original image references, omit `imagePullSecrets`, and set only the pinned tags and `pullPolicy: IfNotPresent`.

### Operator-Managed Images

MeshSync and Broker images cannot be repointed in `values.yaml` or in their custom resources. Meshery Operator names the repositories (`meshery/meshsync` and `nats`) in its own code; the MeshSync and Broker resources accept only a version (`spec.version`), which selects the tag, and the NATS configuration reloader image is fixed in the operator's embedded Broker manifests. You have two ways to serve these images offline:

1. **Configure a registry mirror on every node.** Point your container runtime's mirror for Docker Hub at your private registry (for containerd, the `registry.mirrors` / registry configuration for `docker.io`), and mirror the three operator-managed images under their original names. The operator's requests then resolve to your registry without any Meshery configuration change.
2. **Side-load the three images onto every node** under their original names, as in Option B above. The operator pulls pinned tags with `IfNotPresent`, so pre-loaded images are used without contacting a registry.

## Install



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm install meshery ./meshery-&lt;chart-version&gt;.tgz --namespace meshery --create-namespace -f air-gapped-values.yaml
kubectl get pods --namespace meshery -w</code>
	</div>
</pre>


A pod stuck in `ImagePullBackOff` names an image that was not mirrored or loaded, or a reference your values file did not repoint. Check the exact image it wants:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl describe pod --namespace meshery &lt;pod-name&gt;</code>
	</div>
</pre>


## Providers in an Air-Gapped Environment

Use the Local Provider. The chart's default `env.PROVIDER_BASE_URLS` lists Remote Providers hosted outside your environment; signing in through a Remote Provider, and retrieving extensions and catalog content from one, requires internet access and will not work air-gapped. Meshery Server, Meshery Operator, MeshSync, and the adapters run entirely inside the cluster once their images are present.

Managing a service mesh through an adapter additionally requires that mesh's own images, which are outside the list above. Mirror those from the mesh project's own air-gapped instructions before installing the mesh.

<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-seeking-an-experienced-meshery-contributor-as-a-guide/8166" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate: Seeking an Experienced Meshery Contributor as a Guide
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-to-reactivate-my-slack-account/8168" target="_blank" rel="noopener noreferrer">
            Request to Reactivate My Slack Account
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/bug-new-catalog-publish-request-email-renders-with-oversized-social-icons-and-squashed-content-on-mobile/8138" target="_blank" rel="noopener noreferrer">
            [Bug] &#34;New Catalog publish request&#34; email renders with oversized social icons and squashed content on mobile
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/new-contributor-intro-vedit-agrawal-bughunterx2101-looking-to-pair-with-a-meshmate-for-meshery-ui-ai-adapter/8172" target="_blank" rel="noopener noreferrer">
            New Contributor Intro: Vedit Agrawal (@BugHunterX2101) | Looking to pair with a MeshMate for Meshery UI &amp; AI Adapter
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-frontend-backend-cloud-native-development/8161" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate — Frontend, Backend &amp; Cloud-Native Development
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-for-a-meshmate-active-meshery-contributor/8158" target="_blank" rel="noopener noreferrer">
            Request for a MeshMate — Active Meshery Contributor
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/how-to-present-an-issue-pr-in-a-development-meeting-feedback-wanted/8107" target="_blank" rel="noopener noreferrer">
            How to Present an Issue/PR in a Development Meeting: feedback wanted
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newbie-here-excited-to-join/8052" target="_blank" rel="noopener noreferrer">
            Newbie here -- Excited to join
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
