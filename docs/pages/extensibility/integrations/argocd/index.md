---
layout: integration
title: ArgoCD
subtitle: GitOps Integration with Meshery
permalink: /extensibility/integrations/argocd
docURL: https://docs.meshery.io/extensibility/integrations/argocd
description: Learn how Meshery integrates with ArgoCD, a CNCF continuous delivery tool, to implement GitOps workflows for managing Kubernetes applications and Meshery Designs.
integrations-category: Provisioning
integrations-subcategory: GitOps & Automation
registrant: meshery
---

# Overview

ArgoCD is a declarative, GitOps continuous delivery tool for Kubernetes. As a CNCF project, it provides automated deployment of applications by tracking the desired state in Git repositories. Meshery integrates with ArgoCD to enable GitOps workflows for managing Kubernetes applications and Meshery Designs.

# Setup and Requirements

- Kubernetes cluster with Meshery installed.
- ArgoCD deployed in the target cluster.
- Meshery CLI or Meshery UI access.
- Proper RBAC permissions for GitOps operations.

# Using Meshery with ArgoCD

Meshery’s integration with ArgoCD allows you to:

- Synchronize Meshery Designs with Git repositories.
- Deploy and manage applications using GitOps workflows.
- Visualize deployment states and monitor cluster configurations.

# Meshery Catalog and Design Patterns

Meshery Catalog contains reusable design patterns for ArgoCD, such as:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: sample-app
spec:
  project: default
  source:
    repoURL: https://github.com/example/repo
    targetRevision: HEAD
    path: manifests
  destination:
    server: https://kubernetes.default.svc
    namespace: default
```

You can import, modify, and deploy these designs directly from Meshery.

# Tutorials and Resources

- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [Meshery Docs - ArgoCD Integration](https://docs.meshery.io/extensibility/integrations/argocd)
- [Meshery Playground](https://play.meshery.io/)

# Related Integrations

- [Flux Integration](/extensibility/integrations/flux)
- [Helm Integration](/extensibility/integrations/helm)

# What’s Next

- Contribute new ArgoCD patterns to the [Meshery Catalog](https://meshery.io/catalog).
- Extend GitOps Learning Paths with Meshery and ArgoCD.
