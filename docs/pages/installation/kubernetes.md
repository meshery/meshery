---
layout: default
title: Kubernetes
permalink: installation/kubernetes
type: installation
category: kubernetes
redirect_from:
- installation/platforms/kubernetes
display-title: "false"
language: en
list: include
image: /assets/img/platforms/kubernetes.svg
abstract: Install Meshery on Kubernetes. Deploy Meshery in Kubernetes in-cluster or outside of Kubernetes out-of-cluster.
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Manage your kubernetes clusters with Meshery. Deploy Meshery in kubernetes [in-cluster](#in-cluster-installation) or outside of kubernetes [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your kubernetes clusters_**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Access to an active kubernetes cluster.</li>
  </ol>
</div>

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Set up Ingress on Minikube with the NGINX Ingress Controller](#set-up-ingress-on-minikube-with-the-nginx-ingress-controller)
  - [Installing cert-manager with kubectl](#installing-cert-manager-with-kubectl)

# In-cluster Installation

Follow the steps below to install Meshery in your kubernetes cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on kubernetes.

### Preflight: Cluster Connectivity

Verify your kubeconfig's current context is set the kubernetes cluster you want to deploy Meshery.
{% capture code_content %}kubectl config current-context{% endcapture %}
{% include code.html code=code_content %}

## Installation: Using `mesheryctl`

Once configured, execute the following command to start Meshery.

Before executing the below command, go to ~/.meshery/config.yaml and ensure that current platform is set to kubernetes.
{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment, using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{% include_cached installation/accessing-meshery-ui.md display-title="true" %}

# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your Kubernetes cluster.

<!-- ## Installation: Upload Config File in Meshery Web UI

- Run the below command to generate the _"config_minikube.yaml"_ file for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl config view --minify --flatten > config_minikube.yaml</div></div>
 </pre>

- Upload the generated config file by navigating to _Settings > Environment > Out of Cluster Deployment_ in the Web UI and using the _"Upload kubeconfig"_ option. -->

## Set up Ingress on Minikube with the NGINX Ingress Controller
- Run the below command to enable the NGINX Ingress controller for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">minikube addons enable ingress</div></div>
 </pre>

- To check if NGINX Ingress controller is running
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl get pods -n ingress-nginx</div></div>
 </pre>

## Installing cert-manager with kubectl
- Run the below command to install cert-manager for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml</div></div>
 </pre>

# Kubeconfig Permissions Required

When connecting Meshery to your Kubernetes cluster, your kubeconfig must have sufficient RBAC permissions to allow Meshery to deploy and manage service meshes and applications.

## Required Permissions

Meshery requires permissions to manage the following Kubernetes resources. Each permission is necessary for specific Meshery operations:

### Core Resources
- **Namespaces, Pods, Services, ConfigMaps, Secrets, ServiceAccounts**
  - *Why needed*: Meshery deploys service mesh components and applications into namespaces, manages their lifecycle, and stores configuration data
- **PersistentVolumes, PersistentVolumeClaims**
  - *Why needed*: Service mesh control planes may require persistent storage for state and metrics

### Workload Resources
- **Deployments, StatefulSets, DaemonSets, ReplicaSets**
  - *Why needed*: Service mesh components are deployed as these workload types; Meshery manages their deployment, scaling, and updates
- **Jobs, CronJobs**
  - *Why needed*: Meshery may run maintenance tasks, certificate rotation, or scheduled operations for service meshes

### Network Resources
- **Ingresses, NetworkPolicies**
  - *Why needed*: Meshery configures external access to service mesh UIs and enforces network security policies between services

### RBAC Resources
- **Roles, RoleBindings, ClusterRoles, ClusterRoleBindings**
  - *Why needed*: Service mesh components require specific RBAC permissions; Meshery creates these roles following the principle of least privilege
  - **Note on `escalate` verb**: The `escalate` permission allows Meshery to grant service mesh components the specific permissions they need to function. While this is a powerful permission, it's necessary because different service meshes require different RBAC configurations. In production environments, consider using a dedicated cluster or namespace for Meshery with strict access controls.

### Custom Resources
- **CustomResourceDefinitions (CRDs)**
  - *Why needed*: Service meshes extend Kubernetes with custom resources (e.g., VirtualServices, DestinationRules); Meshery installs and manages these CRDs
- **All Custom Resources (wildcard `*` permissions)**
  - *Why needed*: Different service meshes use different CRDs. The wildcard permission allows Meshery to manage any service mesh without requiring per-mesh permission updates. This is a trade-off between operational convenience and the principle of least privilege. For enhanced security in production, consider creating specific roles for each service mesh you plan to use.

## Minimum Required ClusterRole

Below is a ClusterRole definition with the minimum permissions required for Meshery:

{% capture code_content %}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: meshery-role
rules:
# Core resources
- apiGroups: [""]
  resources: ["namespaces", "pods", "services", "configmaps", "secrets", "serviceaccounts", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["create", "get", "list", "watch", "update", "delete", "patch"]
# Workloads
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets", "daemonsets", "replicasets"]
  verbs: ["create", "get", "list", "watch", "update", "delete", "patch"]
# Batch resources
- apiGroups: ["batch"]
  resources: ["jobs", "cronjobs"]
  verbs: ["create", "get", "list", "watch", "update", "delete", "patch"]
# Network resources
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses", "networkpolicies"]
  verbs: ["create", "get", "list", "watch", "update", "delete", "patch"]
# RBAC
- apiGroups: ["rbac.authorization.k8s.io"]
  resources: ["roles", "rolebindings", "clusterroles", "clusterrolebindings"]
  verbs: ["create", "get", "list", "watch", "update", "delete", "bind", "escalate"]
# CRDs
- apiGroups: ["apiextensions.k8s.io"]
  resources: ["customresourcedefinitions"]
  verbs: ["create", "get", "list", "watch", "update", "delete", "patch"]
# All custom resources
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["create", "get", "list", "watch", "update", "delete", "patch"]
{% endcapture %}
{% include code.html code=code_content %}

**Security Note**: This ClusterRole grants broad permissions necessary for Meshery to manage multiple service meshes. In production environments:
- Consider using a dedicated Kubernetes cluster for Meshery
- Implement additional security controls like admission controllers
- Regularly audit the permissions granted to Meshery
- For specific service mesh deployments, you can create more restrictive roles targeting only the CRDs needed

## Creating a ServiceAccount with Required Permissions

### Step 1: Create the ClusterRole

Save the ClusterRole definition above to a file (e.g., `meshery-clusterrole.yaml`) and apply it:

{% capture code_content %}kubectl apply -f meshery-clusterrole.yaml{% endcapture %}
{% include code.html code=code_content %}

### Step 2: Create a ServiceAccount

{% capture code_content %}kubectl create namespace meshery
kubectl create serviceaccount meshery-sa -n meshery{% endcapture %}
{% include code.html code=code_content %}

### Step 3: Create ClusterRoleBinding

Create a file named `meshery-clusterrolebinding.yaml`:

{% capture code_content %}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: meshery-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: meshery-role
subjects:
- kind: ServiceAccount
  name: meshery-sa
  namespace: meshery
{% endcapture %}
{% include code.html code=code_content %}

Apply the binding:

{% capture code_content %}kubectl apply -f meshery-clusterrolebinding.yaml{% endcapture %}
{% include code.html code=code_content %}

### Step 4: Generate ServiceAccount Token

{% capture code_content %}kubectl create token meshery-sa -n meshery --duration=8760h{% endcapture %}
{% include code.html code=code_content %}

**Security Warning**: The example uses a 1-year token duration (`8760h`). For production use:
- Use shorter token durations aligned with your organization's security policies
- Implement token rotation procedures
- Consider using workload identity or OIDC authentication where available
- Long-lived tokens pose security risks if compromised
- Some organizations may require tokens with durations of days or weeks, not years

## Verifying Permissions

To verify your kubeconfig has the required permissions:

{% capture code_content %}# Check namespace permissions
kubectl auth can-i create namespaces

# Check deployment permissions
kubectl auth can-i create deployments --all-namespaces

# Check CRD permissions
kubectl auth can-i create customresourcedefinitions{% endcapture %}
{% include code.html code=code_content %}

All commands should return `yes` for Meshery to function properly.

## Troubleshooting Permission Issues

If you encounter permission denied errors when using Meshery:

1. **Verify your current permissions:**
   {% capture code_content %}kubectl auth can-i --list{% endcapture %}
   {% include code.html code=code_content %}

2. **Check the ClusterRoleBinding:**
   {% capture code_content %}kubectl get clusterrolebinding meshery-binding -o yaml{% endcapture %}
   {% include code.html code=code_content %}

3. **For restricted environments** (GKE Autopilot, OpenShift), consult your cluster administrator for appropriate permissions or SecurityContextConstraints.

{% include related-discussions.html tag="meshery" %}
