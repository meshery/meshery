---
layout: integration
title: Azure Data Protection
subtitle: Collaborative and visual infrastructure as design for Azure Data Protection
image: /assets/img/integrations/azure-data-protection/icons/color/azure-data-protection-color.svg
permalink: extensibility/integrations/azure-data-protection
docURL: https://docs.meshery.io/extensibility/integrations/azure-data-protection
description: 
integrations-category: Provisioning
integrations-subcategory: Automation & Configuration
registrant: GitHub
components: 
- name: backup-vault
  colorIcon: assets/img/integrations/azure-data-protection/components/backup-vault/icons/color/backup-vault-color.svg
  whiteIcon: assets/img/integrations/azure-data-protection/components/backup-vault/icons/white/backup-vault-white.svg
  description: 
- name: backup-vaults-backup-instance
  colorIcon: assets/img/integrations/azure-data-protection/components/backup-vaults-backup-instance/icons/color/backup-vaults-backup-instance-color.svg
  whiteIcon: assets/img/integrations/azure-data-protection/components/backup-vaults-backup-instance/icons/white/backup-vaults-backup-instance-white.svg
  description: 
- name: backup-vaults-backup-policy
  colorIcon: assets/img/integrations/azure-data-protection/components/backup-vaults-backup-policy/icons/color/backup-vaults-backup-policy-color.svg
  whiteIcon: assets/img/integrations/azure-data-protection/components/backup-vaults-backup-policy/icons/white/backup-vaults-backup-policy-white.svg
  description: 
components-count: 3
relationships: 
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between BackupVaultsBackupInstance and ManagedCluster(azure-container-service)"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between BackupVaultsBackupInstance and BackupVaultsBackupPolicy(azure-data-protection)"
- type: "parent"
  kind: "hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of (parent component) is patched with the configuration of (child component). "
- type: "parent"
  kind: "hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of (parent component) is patched with the configuration of (child component). "
relationship-count: 4
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
