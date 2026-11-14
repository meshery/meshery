---
title: Kubernetes resource examples
description: Downloadable Kubernetes manifests and a custom resource definition for exploring Meshery components.
categories: [contributing]
weight: 15
---

Use these examples to explore how Kubernetes resources map to [Meshery
Components]({{< ref "concepts/logical/components.md" >}}).
The files are Kubernetes manifests, rather than Meshery Model definitions.
Importing a manifest as a Design configures existing components; a CRD describes
a new resource type.

## Example files

| File | Resources | What to explore |
| --- | --- | --- |
| [Deployment and Service](./deployment-service.yaml) | Deployment, Service | Pod labels, Service selectors, and named ports |
| [Custom resource definition](./widget-crd.yaml) | CustomResourceDefinition | A namespaced resource type with a validation schema |
| [Custom resource](./widget.yaml) | Widget | A declaration that conforms to the CRD |

Save the files in the same directory. Run the commands below from that
directory.

## Prerequisites

- A running Meshery instance for importing the Deployment and Service as a
  Design.
- A disposable Kubernetes cluster and `kubectl` configured for that cluster for
  the deployment and CRD exercises.
- Permission to create a namespace and, for the CRD exercise, cluster-scoped
  CustomResourceDefinitions.

Check your current context before applying the examples:

```bash
kubectl config current-context
kubectl create namespace meshery-examples
```

If the namespace already exists, use a fresh cluster or change the namespace in
the commands. The examples do not require access to a production cluster.

## Deployment and Service

The Deployment runs one NGINX Pod. Its Pod template has the label `app:
meshery-example-web`, which the Service uses as its selector.
The Service forwards port `80` to the container's named `http` port. It uses
`ClusterIP`, so it does not request an external load balancer.

In Meshery, open **Designs > Import Design**, choose **File Upload**,
select `deployment-service.yaml`, and click **Import**.
Meshery detects the Kubernetes manifest format automatically.
Inspect the imported Deployment and Service and compare their configuration with
the manifest.
Importing a Design does not deploy it. See [Importing and Exporting Designs]({{<
ref "guides/configuration-management/import-export-designs.md" >}}) for the full
workflow.

To check the example directly against Kubernetes:

```bash
kubectl apply --dry-run=server -n meshery-examples -f deployment-service.yaml
kubectl apply -n meshery-examples -f deployment-service.yaml
kubectl rollout status -n meshery-examples deployment/meshery-example-web --timeout=120s
kubectl get pods -n meshery-examples -l app=meshery-example-web
kubectl port-forward -n meshery-examples service/meshery-example-web 8080:80
```

With port forwarding running, open `http://localhost:8080` in a browser. The
NGINX welcome page confirms that the Service reaches the Pod.
Press `Ctrl+C` to stop port forwarding. If the Pod does not become ready, check
its events with `kubectl describe pod -n meshery-examples -l
app=meshery-example-web`.

## CRD and custom resource

`widget-crd.yaml` defines a namespaced `Widget` resource in the
`examples.meshery.io` API group.
The `v1alpha1` version is both served and used for storage. Its schema requires
a `spec.message` string and restricts optional `spec.replicas` to an integer of
at least one.
`widget.yaml` supplies both fields.

This example defines an API and stores a resource. It has no controller, so
creating a Widget does not create Pods or reconcile `spec.replicas` into a
running workload.
For background, see [Kubernetes custom resource
definitions](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/).

Apply the CRD first and wait until Kubernetes establishes the new API before
creating the Widget:

```bash
kubectl apply --dry-run=server -f widget-crd.yaml
kubectl apply -f widget-crd.yaml
kubectl wait --for=condition=Established --timeout=60s crd/widgets.examples.meshery.io
kubectl apply --dry-run=server -n meshery-examples -f widget.yaml
kubectl apply -n meshery-examples -f widget.yaml
kubectl get widgets.examples.meshery.io meshery-example-widget \
  -n meshery-examples -o yaml
```

The returned resource should contain `spec.message: Hello from Meshery` and
`spec.replicas: 1`.
To explore CRD-based Meshery components, use this CRD as a schema reference
alongside [Contributing to Models]({{< ref
"project/contributing/models/models.md" >}}) and [Creating Models]({{< ref
"guides/configuration-management/creating-models/index.md" >}}).
The Widget must have a corresponding component registered in Meshery before it
can be imported as a Design.

## Cleanup

Delete the example resources before deleting their namespace. Deleting the CRD
removes all Widgets in its API group across the cluster; use these commands only
in the disposable cluster used for this exercise.

```bash
kubectl delete -n meshery-examples -f widget.yaml
kubectl delete -f widget-crd.yaml
kubectl delete -n meshery-examples -f deployment-service.yaml
kubectl delete namespace meshery-examples
```

If you completed only the Deployment and Service exercise, run only the last two
commands.
