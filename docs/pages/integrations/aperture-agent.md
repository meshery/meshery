---
layout: enhanced
title: Aperture Agent
subtitle: Collaborative and visual infrastructure as code for Aperture Agent
image: /assets/img/integrations/aperture-agent/icons/color/aperture-agent-color.svg
permalink: extensibility/integrations/aperture-agent
docURL: https://docs.meshery.io/extensibility/integrations/aperture-agent
description: 
integrations-category: Orchestration & Management
integrations-subcategory: Scheduling & Orchestration
registrant: artifacthub
components: 
- name: agent
  colorIcon: assets/img/integrations/aperture-agent/components/agent/icons/color/agent-color.svg
  whiteIcon: assets/img/integrations/aperture-agent/components/agent/icons/white/agent-white.svg
  description: 
featureList: [
  "Drag-n-drop cloud native infrastructure designer to configure, model, and deploy your workloads.",
  "Invite anyone to review and make changes to your private designs.",
  "Ongoing synchronization of Kubernetes configuration and changes across any number of clusters."
]
howItWorks: "Collaborative Infrastructure as Code"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
display-title: "false"
---
<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>

<p>
The Aperture Agent is the decision executor of the Aperture system. In addition to gathering data, the Aperture Agent functions as a gatekeeper, acting on traffic based on decisions made by the Aperture Controller. Specifically, depending on feedback from the Controller, the Agent will effectively allow or drop incoming requests. Further supporting the Controller, the Agent works to inject information into traffic, including the specific traffic-shaping decisions made and classification labels which can later be used in policing. One Agent is deployed per node.

</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
