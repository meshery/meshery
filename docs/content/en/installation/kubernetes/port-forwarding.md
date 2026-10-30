---
title: Port Forwarding
categories: [kubernetes]
aliases:
 - /installation/kubernetes/port-forwarding
display_title: false
description: Access Meshery UI using port forwarding when LoadBalancer is not available.
---

# Accessing Meshery UI Using Port Forwarding

Port forwarding is a universal method to access Meshery UI when a `LoadBalancer` service is not available or practical. This is common in local development environments (Minikube, KinD) or when you don't want to expose Meshery externally.

This guide is applicable to all cloud providers including AWS (EKS), Azure (AKS), GCP (GKE), and local cluster services (Minikube, KinD).

## When to Use Port Forwarding

Use port forwarding in the following scenarios:

- **Local development clusters** (Minikube, KinD) where LoadBalancer controllers are not installed
- **Testing environments** where you don't want to expose Meshery UI publicly
- **Temporary access** to Meshery UI without configuring ingress or load balancers
- **Debugging** Meshery deployments in restricted networks

## Prerequisites

- Meshery is installed and running in your Kubernetes cluster
- `kubectl` is configured to access your cluster
- Port `9081` is available on your local machine (or specify an alternative port)

## Using mesheryctl (Recommended)

The easiest way to access Meshery UI with port forwarding is using `mesheryctl`:

{{< code code="mesheryctl system dashboard --port-forward" >}}

This command:
- Automatically sets up port forwarding from port `9081` on your local machine to the Meshery service in the cluster
- Opens Meshery UI in your default browser
- Handles service discovery automatically

### Custom Port

If port `9081` is already in use, specify an alternative local port:

{{< code code="mesheryctl system dashboard --port-forward -p 9082" >}}

### Skip Browser Opening

To set up port forwarding without automatically opening the browser:

{{< code code="mesheryctl system dashboard --port-forward --skip-browser" >}}

## Using kubectl

For more control or when `mesheryctl` is not available, use `kubectl` directly:

{{< code code="kubectl port-forward svc/meshery 9081:9081 --namespace meshery" >}}

This command forwards:
- Local port `9081` → Service port `9081` in the `meshery` namespace

### Custom Namespace or Service Name

If you installed Meshery in a different namespace or with a different service name:

{{< code code="kubectl port-forward svc/meshery 9081:9081 --namespace <your-namespace>" >}}

### Port Forward to a Specific Pod

To port forward directly to a Meshery pod (useful for debugging):

{{< code code="kubectl port-forward deployment/meshery 9081:8080 --namespace meshery" >}}

## Verify the Connection

Once port forwarding is established:

1. Open your browser and navigate to `http://localhost:9081`
2. You should see the Meshery UI login screen
3. Log in with your Meshery credentials

## Troubleshooting

### Port Already in Use

If port `9081` is already in use on your machine:

- Use an alternative port: `mesheryctl system dashboard --port-forward -p 9082`
- Or stop the process using port `9081`

### Connection Refused

If you get a connection refused error:

1. Verify Meshery is running: `kubectl get pods --namespace meshery`
2. Check the service exists: `kubectl get svc --namespace meshery`
3. Ensure the service name is `meshery` (or adjust the port-forward command accordingly)

### Timeout Errors

If port forwarding times out:

1. Check your kubeconfig context: `kubectl config current-context`
2. Verify network connectivity to the cluster
3. Check if firewall rules are blocking the connection

## LoadBalancer vs Port Forwarding

| Feature | LoadBalancer | Port Forwarding |
|---------|--------------|-----------------|
| **External Access** | Yes (cloud provider assigns external IP) | No (local access only) |
| **Production Use** | Recommended | Not recommended |
| **Setup Complexity** | Requires cloud load balancer | Simple, one command |
| **Use Case** | Production, multi-user access | Development, testing, temporary access |

For production deployments, consider using:
- **LoadBalancer** service type (cloud providers provision external IP automatically)
- **Ingress controllers** (NGINX, Traefik) for advanced routing
- **Reverse proxies** for additional security layers

## Related Documentation

- [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) - Install Meshery using Helm
- [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}}) - General Kubernetes installation guide
- [mesheryctl system dashboard]({{< ref "reference/mesheryctl/system/dashboard.md" >}}) - CLI reference for dashboard command
- [Accessing Meshery UI]({{< ref "installation/accessing-meshery-ui.md" >}}) - General guide for accessing Meshery UI

{{< related-discussions tag="meshery" >}}
