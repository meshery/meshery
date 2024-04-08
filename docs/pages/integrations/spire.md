---
layout: enhanced
title: SPIRE
subtitle: Collaborative and visual infrastructure as code for SPIRE
image: /assets/img/integrations/spire/icons/color/spire-color.svg
permalink: extensibility/integrations/spire
docURL: https://docs.meshery.io/extensibility/integrations/spire
description: 
integrations-category: Security & Compliance
integrations-subcategory: Key Management
registrant: Artifact Hub
components: 
- name: cluster-federated-trust-domain
  colorIcon: assets/img/integrations/spire/components/cluster-federated-trust-domain/icons/color/cluster-federated-trust-domain-color.svg
  whiteIcon: assets/img/integrations/spire/components/cluster-federated-trust-domain/icons/white/cluster-federated-trust-domain-white.svg
  description: 
- name: cluster-spiffeid
  colorIcon: assets/img/integrations/spire/components/cluster-spiffeid/icons/color/cluster-spiffeid-color.svg
  whiteIcon: assets/img/integrations/spire/components/cluster-spiffeid/icons/white/cluster-spiffeid-white.svg
  description: 
- name: controller-manager-config
  colorIcon: assets/img/integrations/spire/components/controller-manager-config/icons/color/controller-manager-config-color.svg
  whiteIcon: assets/img/integrations/spire/components/controller-manager-config/icons/white/controller-manager-config-white.svg
  description: 
featureList: [
  "SPIRE design and deployment best practices",
  "Configure SPIRE to securely issue and renew SVIDs.",
  "Design workloads to establish trust between each other by establishing an mTLS connection or by signing and verifying a JWT token."
]
howItWorks: "Configure your signing framework"
howItWorksDetails: "Simplify access from identified services to secret stores, databases, services meshes and cloud provider services."
language: en
list: include
type: extensibility
category: integrations
display-title: "false"
---
<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>

<p>
The SPIFFE Runtime Environment
</p>
<p>
    Meshery deploys and manages SPIRE (the SPIFFE Runtime Environment). SPIRE is a toolchain of APIs for establishing trust between software systems across a wide variety of hosting platforms. SPIRE exposes the SPIFFE Workload API, which can attest running software systems and issue SPIFFE IDs and SVIDs to them. 
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
