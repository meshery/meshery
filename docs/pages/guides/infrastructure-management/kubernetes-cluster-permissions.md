---
layout: default
title: Kubernetes Cluster Permissions
permalink: guides/infrastructure-management/kubernetes-cluster-permissions
type: guides
category: infrastructure
language: en
abstract: Minimum kubeconfig permissions required for Meshery to connect to and manage Kubernetes clusters.
list: include
---

Meshery requires specific permissions to connect to and manage your Kubernetes clusters. This reference guide documents the minimum permissions needed in your kubeconfig for Meshery to function properly.

## Overview

When you connect a Kubernetes cluster to Meshery by uploading a kubeconfig file, Meshery performs several operations to:

1. Verify cluster connectivity and identify the cluster
2. Deploy Meshery components (Operator, MeshSync, Broker) to the cluster
3. Discover and synchronize cluster resources
4. Manage workloads and configurations

Each of these operations requires specific Kubernetes RBAC permissions.

## Minimum Required Permissions

### Core Connectivity Permissions

At a minimum, your kubeconfig must have the following permissions for Meshery to establish a connection to your cluster:

#### 1. Read Access to `kube-system` Namespace

Meshery identifies each Kubernetes cluster by querying the UID of the `kube-system` namespace. This serves as a unique identifier for the cluster.

**Required Permission:**
- `get` on `namespaces` resource for `kube-system` namespace

**RBAC Example:**
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: meshery-minimum-access
rules:
- apiGroups: [""]
  resources: ["namespaces"]
  resourceNames: ["kube-system"]
  verbs: ["get"]
```

{% include alert.html type="warning" title="Critical Requirement" content="Without read access to the <code>kube-system</code> namespace, Meshery cannot identify the cluster and the connection will fail." %}

### Full Meshery Management Permissions

For Meshery to fully manage your cluster, including deploying Meshery Operator, MeshSync, and discovering resources, the following permissions are required:

#### 2. Meshery Namespace Management

Meshery deploys its components to a dedicated namespace (by default, `meshery`). The kubeconfig needs permissions to create and manage this namespace.

**Required Permissions:**
- `create`, `get`, `list`, `watch`, `update`, `patch`, `delete` on `namespaces`

#### 3. Operator and Component Deployment

Meshery Operator, MeshSync, and Broker are deployed as Deployments, Services, and CustomResourceDefinitions.

**Required Permissions:**
- Full access to `deployments`, `replicasets`, `pods` in the `meshery` namespace
- Full access to `services`, `serviceaccounts` in the `meshery` namespace
- Full access to `configmaps`, `secrets` in the `meshery` namespace
- Cluster-wide access to `customresourcedefinitions`
- Cluster-wide `create`, `get`, `list`, `watch`, `update`, `patch`, `delete` on CRDs: `brokers.meshery.io`, `meshsyncs.meshery.io`

#### 4. RBAC for Meshery Components

Meshery Operator and MeshSync require their own ServiceAccounts, ClusterRoles, and ClusterRoleBindings.

**Required Permissions:**
- Cluster-wide `create`, `get`, `list`, `watch`, `update`, `patch`, `delete` on `clusterroles`, `clusterrolebindings`, `roles`, `rolebindings`

#### 5. Resource Discovery (MeshSync)

MeshSync discovers and continuously synchronizes the state of resources in your cluster. It requires broad read access across the cluster.

**Required Permissions:**
- Cluster-wide `get`, `list`, `watch` on most core and custom resources including:
  - Core resources: `pods`, `services`, `deployments`, `replicasets`, `statefulsets`, `daemonsets`, `jobs`, `cronjobs`, `configmaps`, `secrets`, `persistentvolumes`, `persistentvolumeclaims`, `namespaces`, `nodes`, `events`
  - Networking: `ingresses`, `networkpolicies`, `services`
  - Custom resources defined by installed operators and applications

{% include alert.html type="info" title="Configurable Discovery" content="You can configure which resources MeshSync discovers by editing the <code>meshsyncs.meshery.io</code> CRD. See the <a href='/concepts/architecture/meshsync'>MeshSync documentation</a> for details on blacklisting specific resources." %}

## Complete RBAC Example

Below is a complete example of a ClusterRole that grants Meshery all necessary permissions to manage a Kubernetes cluster:

{% capture code_content %}apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: meshery-full-access
rules:
# Cluster identification
- apiGroups: [""]
  resources: ["namespaces"]
  verbs: ["get", "list", "create", "delete", "watch"]

# Deployments and workloads
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets", "statefulsets", "daemonsets"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Core resources
- apiGroups: [""]
  resources: ["pods", "services", "configmaps", "secrets", "serviceaccounts", "persistentvolumes", "persistentvolumeclaims", "nodes", "events"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Batch resources
- apiGroups: ["batch"]
  resources: ["jobs", "cronjobs"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Networking
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses", "networkpolicies"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# RBAC
- apiGroups: ["rbac.authorization.k8s.io"]
  resources: ["roles", "rolebindings", "clusterroles", "clusterrolebindings"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Custom Resource Definitions
- apiGroups: ["apiextensions.k8s.io"]
  resources: ["customresourcedefinitions"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Meshery CRDs
- apiGroups: ["meshery.io"]
  resources: ["brokers", "meshsyncs"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Discovery of all resources (for MeshSync)
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["get", "list", "watch"]
{% endcapture %}
{% include code.html code=code_content %}

To use this ClusterRole, create a ServiceAccount and bind it:

{% capture code_content %}---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: meshery
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: meshery-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: meshery-full-access
subjects:
- kind: ServiceAccount
  name: meshery
  namespace: default
{% endcapture %}
{% include code.html code=code_content %}

Then generate a kubeconfig for this ServiceAccount to use with Meshery.

## Permission Levels and Use Cases

### Read-Only Monitoring

If you only want Meshery to monitor your cluster without deploying components or making changes:

**Required Permissions:**
- `get` on `namespaces/kube-system`
- Cluster-wide `get`, `list`, `watch` on all resources you want to monitor

**Limitations:**
- Cannot deploy Meshery Operator, MeshSync, or Broker
- Cannot create or modify resources
- Cannot perform lifecycle management operations

### Standard Management

For typical Meshery operations including deploying applications and managing configurations:

**Required Permissions:**
- All permissions in the [Complete RBAC Example](#complete-rbac-example)

**Capabilities:**
- Full cluster management
- Deploy Meshery components
- Create, update, and delete resources
- Real-time resource synchronization

### Namespace-Scoped Access

If you want to restrict Meshery to manage only specific namespaces:

1. Create a Role (instead of ClusterRole) with appropriate permissions in each target namespace
2. Grant cluster-wide read access to `namespaces` for cluster identification
3. Grant cluster-wide read access for resource discovery (optional)

{% include alert.html type="info" title="Operator Limitations" content="Meshery Operator and MeshSync typically require cluster-wide permissions. Namespace-scoped access may limit functionality." %}

## Security Considerations

### Principle of Least Privilege

While Meshery can function with cluster-admin privileges, consider granting only the minimum permissions needed for your use case:

- For **monitoring only**: Grant read-only access
- For **standard operations**: Use the complete RBAC example provided
- For **namespace-specific management**: Use namespace-scoped roles

### Credential Management

- Store kubeconfig files securely
- Use short-lived tokens when possible
- Rotate credentials regularly
- Consider using service account tokens instead of user credentials

### In-Cluster vs. Out-of-Cluster Deployments

**In-Cluster Deployment:**
When Meshery is deployed inside a Kubernetes cluster, it can use in-cluster authentication with a ServiceAccount. This is more secure as credentials don't leave the cluster.

**Out-of-Cluster Deployment:**
When Meshery runs outside Kubernetes (e.g., on Docker), you must provide a kubeconfig file. Ensure this file is protected and contains only the necessary permissions.

## Troubleshooting Permission Issues

### Connection Fails with "Unreachable Kubernetes API"

**Cause:** Meshery cannot read the `kube-system` namespace.

**Solution:** Verify your kubeconfig has `get` permission on `namespaces/kube-system`:

{% capture code_content %}kubectl auth can-i get namespace/kube-system --as=system:serviceaccount:default:meshery{% endcapture %}
{% include code.html code=code_content %}

### Operator Fails to Deploy

**Cause:** Insufficient permissions to create resources in the `meshery` namespace or cluster-wide CRDs.

**Solution:** Ensure your kubeconfig can create namespaces, deployments, and CustomResourceDefinitions:

{% capture code_content %}kubectl auth can-i create namespace
kubectl auth can-i create deployment -n meshery
kubectl auth can-i create customresourcedefinition{% endcapture %}
{% include code.html code=code_content %}

### MeshSync Not Discovering Resources

**Cause:** Missing `list` or `watch` permissions on specific resource types.

**Solution:** Check MeshSync logs and verify permissions:

{% capture code_content %}kubectl logs -n meshery -l app=meshsync
kubectl auth can-i list pods --all-namespaces{% endcapture %}
{% include code.html code=code_content %}

## Related Documentation

- [Meshery Operator](/concepts/architecture/operator) - Learn about Meshery Operator deployment and management
- [MeshSync](/concepts/architecture/meshsync) - Understand MeshSync's resource discovery and synchronization
- [Managing Connections](/guides/infrastructure-management/lifecycle-management) - Guide for managing Kubernetes cluster connections
- [Troubleshooting Operator & MeshSync](/guides/troubleshooting/meshery-operator-meshsync) - Debug common issues with Meshery components
- [Kubernetes Installation](/installation/kubernetes) - Deploy Meshery on Kubernetes

## FAQs

<details>
<summary><strong>Question:</strong> Why does Meshery need access to the <code>kube-system</code> namespace?</summary>
<strong>Answer:</strong> Meshery uses the UID of the <code>kube-system</code> namespace as a unique identifier for each Kubernetes cluster. This allows Meshery to distinguish between multiple clusters and track their individual states. Without this access, Meshery cannot establish a connection to the cluster.
</details>

<details>
<summary><strong>Question:</strong> Can I use Meshery with a kubeconfig that doesn't have cluster-admin privileges?</summary>
<strong>Answer:</strong> Yes, Meshery does not require cluster-admin privileges. You can create a custom ClusterRole with only the permissions documented in this guide. However, you must at minimum have read access to the <code>kube-system</code> namespace.
</details>

<details>
<summary><strong>Question:</strong> What happens if I don't grant permissions to deploy the Operator?</summary>
<strong>Answer:</strong> Without permissions to deploy the Meshery Operator, you can still connect to the cluster and view basic information, but you won't be able to:
<ul>
  <li>Deploy MeshSync for real-time resource synchronization</li>
  <li>Deploy the Broker for event streaming</li>
  <li>Use embedded mode for MeshSync as an alternative</li>
</ul>
You may be able to use MeshSync in embedded mode, which runs within the Meshery Server and requires fewer cluster permissions.
</details>

<details>
<summary><strong>Question:</strong> How do I generate a kubeconfig for a specific ServiceAccount?</summary>
<strong>Answer:</strong> Create a ServiceAccount with appropriate permissions, then generate a kubeconfig:
<pre><code>
# Create ServiceAccount and bindings (see RBAC example above)
kubectl apply -f meshery-rbac.yaml

# Get the ServiceAccount token
TOKEN=$(kubectl create token meshery -n default)

# Get cluster info
CLUSTER_NAME=$(kubectl config view --minify -o jsonpath='{.clusters[0].name}')
CLUSTER_SERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
CLUSTER_CA=$(kubectl config view --minify --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')

# Create kubeconfig
kubectl config set-cluster $CLUSTER_NAME --server=$CLUSTER_SERVER --certificate-authority-data=$CLUSTER_CA --kubeconfig=meshery-kubeconfig
kubectl config set-credentials meshery --token=$TOKEN --kubeconfig=meshery-kubeconfig
kubectl config set-context meshery --cluster=$CLUSTER_NAME --user=meshery --kubeconfig=meshery-kubeconfig
kubectl config use-context meshery --kubeconfig=meshery-kubeconfig
</code></pre>
</details>

<details>
<summary><strong>Question:</strong> Does MeshSync need write access to cluster resources?</summary>
<strong>Answer:</strong> No, MeshSync only needs read access (<code>get</code>, <code>list</code>, <code>watch</code>) to discover and synchronize cluster resources. However, Meshery Server itself may need write access if you want to deploy workloads or modify configurations through Meshery.
</details>

{% include discuss.html %}
