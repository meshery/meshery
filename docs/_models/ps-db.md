---
layout: integration
title: Percona Server for MySQL
subtitle: Meshery provides performance reports, including performance test results, node resource metrics etc. so that operators may easily understand the overhead of their service mesh’s control plane and data plane in context of the overhead incurred on nodes running within the cluster. In order to generate performance test reports of service meshes and their workloads, Meshery uses Grafana and/or Prometheus as visualization and metrics systems, respectively. This guide outlines the requirements necessary for Meshery to connect to these systems. The steps may vary depending upon the service mesh and its configuration.
image: /assets/img/integrations/ps-db/icons/color/ps-db-color.svg
permalink: extensibility/integrations/ps-db
docURL: https://docs.meshery.io/extensibility/integrations/ps-db
description: 
integrations-category: Database
integrations-subcategory: App Definition and Development
registrant: Github
components: 
- name: percona-server-my-sql-backup
  colorIcon: assets/img/integrations/ps-db/components/percona-server-my-sql-backup/icons/color/percona-server-my-sql-backup-color.svg
  whiteIcon: assets/img/integrations/ps-db/components/percona-server-my-sql-backup/icons/white/percona-server-my-sql-backup-white.svg
  description: 
- name: percona-server-my-sql-restore
  colorIcon: assets/img/integrations/ps-db/components/percona-server-my-sql-restore/icons/color/percona-server-my-sql-restore-color.svg
  whiteIcon: assets/img/integrations/ps-db/components/percona-server-my-sql-restore/icons/white/percona-server-my-sql-restore-white.svg
  description: 
- name: percona-server-my-sql
  colorIcon: assets/img/integrations/ps-db/components/percona-server-my-sql/icons/color/percona-server-my-sql-color.svg
  whiteIcon: assets/img/integrations/ps-db/components/percona-server-my-sql/icons/white/percona-server-my-sql-white.svg
  description: 
featureList: [
  "Native support for PromQL",
  "Create custom charts with your own Prometheus queries",
  "Keep charts in-sync with Mesherys panel viewer"
]
howItWorks: "Meshery provides performance reports, including performance test results, node resource metrics etc. so that operators may easily understand the overhead of their service mesh’s control plane and data plane in context of the overhead incurred on nodes running within the cluster. In order to generate performance test reports of service meshes and their workloads, Meshery uses Grafana and/or Prometheus as visualization and metrics systems, respectively. This guide outlines the requirements necessary for Meshery to connect to these systems. The steps may vary depending upon the service mesh and its configuration."
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
