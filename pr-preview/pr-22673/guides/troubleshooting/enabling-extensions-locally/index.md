# Troubleshooting Extensions for Local Development

> Find out why a Meshery extension doesn't load in your local development environment, and what you can fix yourself.

---

LLMS index: [llms.txt](/pr-preview/pr-22673/llms.txt)

---

Use this guide when an extension doesn't show up in your local Meshery development environment. Meshery extensions integrate in two different ways, and each fails differently, so first find out which kind you have.

## Step 1: Find the extension's repository

Most Meshery extensions are open source, with repositories in the [meshery-extensions GitHub organization](https://github.com/orgs/meshery-extensions/repositories?type=all). Find the repository for the extension you are troubleshooting.

If the repository is public, follow the steps for [Loosely-Coupled Extensions](#loosely-coupled-extensions).

If you can't find it, or it is private, the extension is probably tightly coupled, and the problem is probably a dependency mismatch. Follow the steps for [Tightly-Coupled Extensions](#tightly-coupled-extensions).

## Step 2: Troubleshoot by extension type

### Loosely-Coupled Extensions

These extensions run as separate processes and talk to Meshery Server through its APIs, so you can inspect them directly.

- Example: [Meshery Adapters](/pr-preview/pr-22673/concepts/architecture/adapters/) like the [meshery-istio](https://github.com/meshery-extensions/meshery-istio) adapter.
- Common issues: networking problems, port conflicts, or errors in the extension itself.
- Debugging steps: run `docker ps` to check that the container is running, and `docker logs <container-id>` to read its logs. If Meshery runs in Kubernetes, use `kubectl get pods -n meshery` and `kubectl logs`.

### Tightly-Coupled Extensions

These extensions load into Meshery UI as plugins, so they need package versions that exactly match your Meshery build.

- Example: [Kanvas](/pr-preview/pr-22673/extensions/extensions/kanvas/), the collaborative visual designer.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Resolving Dependency Mismatches</div>


Problems with tightly-coupled extensions are usually not code bugs. They are version conflicts between the extension and your local Meshery build.
</div>


#### Why you can't fix this locally

The extension's source repository is private, so:

- You can't build the extension from source.
- You can only use a [pre-built package](https://github.com/meshery-extensions/meshery-extensions-packages).

#### The dependency mismatch problem

A pre-built package works only if its dependencies happen to match your local Meshery build:

- Updating your local Meshery can break a package that worked before.
- Failures often show no clear error message.
- Without source access, you can't rebuild the package to match.

#### Next steps

You can try other package versions from the [meshery-extensions-packages](https://github.com/meshery-extensions/meshery-extensions-packages) repository, but trial and error isn't guaranteed to work. Instead:

- **Ask for help** in the [Meshery community](/pr-preview/pr-22673/project/community/).
- **Check access requirements** or find maintainers in the [community handbook](https://meshery.io/community).
