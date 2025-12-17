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

# Kubeconfig Permissions

When connecting Meshery to your Kubernetes cluster, the required permissions vary based on:
- Which service meshes you're deploying (Istio, Linkerd, Consul, etc.)
- Your deployment method (in-cluster vs out-of-cluster)
- Which Meshery features you're using

This guide shows you **how to discover and troubleshoot** the permissions Meshery needs for your specific use case.

## Understanding Permission Requirements

Meshery's permission requirements are **dynamic** and depend on:

1. **Service Mesh Adapters**: Each service mesh (Istio, Linkerd, Consul) has its own CRDs and resource requirements
2. **Deployment Mode**: In-cluster deployments may need different permissions than out-of-cluster
3. **Features Used**: Different Meshery features (performance testing, configuration management, observability) access different resources

This means there's no single "correct" permission set - your needs will vary.

## Minimal Starting Permissions

For initial Meshery deployment, start with these minimal permissions and add more as needed:

{% capture code_content %}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: meshery-minimal
rules:
# Minimal permissions to deploy Meshery
- apiGroups: [""]
  resources: ["namespaces", "pods", "services"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["get", "list", "watch"]
{% endcapture %}
{% include code.html code=code_content %}

**Note**: This minimal role allows Meshery to start and observe your cluster, but you'll need additional permissions to deploy and manage service meshes.

## Discovering Required Permissions

When Meshery needs permissions it doesn't have, you'll see errors like:
```
Error: services is forbidden: User "system:serviceaccount:meshery:meshery-sa" 
cannot create resource "services" in API group "" in the namespace "istio-system"
```

### How to Grant Missing Permissions

1. **Identify the missing permission** from the error message:
   - Resource: `services`
   - Verb: `create`
   - API Group: `""` (core)
   - Namespace: `istio-system`

2. **Update your ClusterRole** to add the permission:

{% capture code_content %}
# Add this rule to your ClusterRole
- apiGroups: [""]
  resources: ["services"]
  verbs: ["create", "get", "list", "watch", "update", "delete"]
{% endcapture %}
{% include code.html code=code_content %}

3. **Apply the updated role**:

{% capture code_content %}kubectl apply -f meshery-clusterrole.yaml{% endcapture %}
{% include code.html code=code_content %}

## Common Permission Patterns

Based on the service mesh you're deploying, you'll typically need permissions for:

### For Istio
- CRDs: `VirtualService`, `DestinationRule`, `Gateway`, etc.
- Resources: Services, Deployments, ConfigMaps in `istio-system` namespace

### For Linkerd
- CRDs: `ServiceProfile`, `TrafficSplit`
- Resources: Deployments, Services in `linkerd` namespace

### For Consul
- CRDs: `ServiceDefaults`, `ServiceRouter`, `ServiceSplitter`
- Resources: StatefulSets, Services in `consul` namespace

## Checking Your Current Permissions

To see what permissions your kubeconfig currently has:

{% capture code_content %}# Check if you can perform specific actions
kubectl auth can-i create deployments --all-namespaces
kubectl auth can-i create customresourcedefinitions
kubectl auth can-i get pods --all-namespaces

# List all your permissions
kubectl auth can-i --list{% endcapture %}
{% include code.html code=code_content %}

## Creating a ServiceAccount (Recommended)

Instead of using your personal kubeconfig, create a dedicated ServiceAccount for Meshery:

### Step 1: Create Namespace and ServiceAccount

{% capture code_content %}kubectl create namespace meshery
kubectl create serviceaccount meshery-sa -n meshery{% endcapture %}
{% include code.html code=code_content %}

### Step 2: Create ClusterRole with Minimal Permissions

Start with minimal permissions and expand as needed based on errors you encounter.

### Step 3: Bind the Role

{% capture code_content %}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: meshery-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: meshery-minimal
subjects:
- kind: ServiceAccount
  name: meshery-sa
  namespace: meshery
{% endcapture %}
{% include code.html code=code_content %}

### Step 4: Generate Token

{% capture code_content %}# Create a short-lived token (adjust duration as needed)
kubectl create token meshery-sa -n meshery --duration=8760h{% endcapture %}
{% include code.html code=code_content %}

**Security Best Practice**: Use the shortest token duration that works for your workflow. For production, consider integrating with your organization's identity provider using OIDC.

## Troubleshooting Permission Errors

### Common Error: "Forbidden"
```
Error: services is forbidden: User cannot create resource
```

**Solution**: Add the missing verb and resource to your ClusterRole.

### Common Error: "Unknown CustomResourceDefinition"
```
Error: the server doesn't have a resource type "virtualservices"
```

**Solution**: The CRD isn't installed, or you lack permissions to access it. First, verify the CRD exists:

{% capture code_content %}kubectl get crd virtualservices.networking.istio.io{% endcapture %}
{% include code.html code=code_content %}

### Common Error: "Cannot list resource in the namespace"

**Solution**: Your role may have cluster-wide permissions but not namespace-specific permissions, or vice versa.

## Security Considerations

When granting permissions to Meshery:

1. **Start Minimal**: Begin with read-only permissions and add write permissions only as needed
2. **Use Namespaces**: Limit permissions to specific namespaces where possible
3. **Audit Regularly**: Review and remove unused permissions periodically
4. **Use ServiceAccounts**: Don't use cluster-admin for Meshery in production
5. **Token Rotation**: Implement token rotation policies for long-running deployments

## Advanced: Environment-Specific Considerations

### GKE Autopilot
- Some permissions may be restricted by Google's security policies
- Workload Identity is recommended over service account keys

### OpenShift
- Use SecurityContextConstraints (SCC) instead of PSPs
- May need additional permissions for routes

### AKS
- Azure AD integration is recommended for authentication
- Use managed identities where possible

## Getting Help

If you continue to experience permission errors:

1. Check Meshery logs for specific permission denials
2. Post in [Meshery Slack](https://slack.layer5.io) with the error message
3. Review [GitHub Discussions](https://github.com/meshery/meshery/discussions) for similar issues

## Reference: Full Permissions Example

For users who need a comprehensive starting point (understanding it may grant more permissions than necessary):

See [Meshery's Helm chart RBAC configuration](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery/templates) for the permissions used in standard installations.

**Warning**: These permissions are designed for a full Meshery installation with all features enabled. For production use, audit and reduce these permissions based on your actual needs.

{% include related-discussions.html tag="meshery" %}
