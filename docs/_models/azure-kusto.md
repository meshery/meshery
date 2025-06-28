---
layout: integration
title: Azure Kusto
subtitle: Collaborative and visual infrastructure as design for Azure Kusto
image: /assets/img/integrations/azure-kusto/icons/color/azure-kusto-color.svg
permalink: extensibility/integrations/azure-kusto
docURL: https://docs.meshery.io/extensibility/integrations/azure-kusto
description: 
integrations-category: Orchestration & Management
integrations-subcategory: Kubernetes
registrant: GitHub
components: 
- name: cluster
  colorIcon: assets/img/integrations/azure-kusto/components/cluster/icons/color/cluster-color.svg
  whiteIcon: assets/img/integrations/azure-kusto/components/cluster/icons/white/cluster-white.svg
  description: 
- name: database
  colorIcon: assets/img/integrations/azure-kusto/components/database/icons/color/database-color.svg
  whiteIcon: assets/img/integrations/azure-kusto/components/database/icons/white/database-white.svg
  description: 
- name: data-connection
  colorIcon: assets/img/integrations/azure-kusto/components/data-connection/icons/color/data-connection-color.svg
  whiteIcon: assets/img/integrations/azure-kusto/components/data-connection/icons/white/data-connection-white.svg
  description: 
components-count: 3
relationships: 
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of (parent component) is patched with the configuration of (child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of (parent component) is patched with the configuration of (child component). "
relationship-count: 2
featureList: [
  "Drag-n-drop cloud native infrastructure designer to configure, model, and deploy your workloads.",
  "Invite anyone to review and make changes to your private designs.",
  "Ongoing synchronization of Kubernetes configuration and changes across any number of clusters."
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
