---
title: "Ingesting Helm Charts from External Sources"
description: "Guide to import Helm charts from Artifact Hub and Bitnami into Meshery"
date: 2025-01-18
weight: 1
---

# Ingesting Helm Charts from Artifact Hub or Bitnami

This guide explains how to download Helm charts from external sources like [Artifact Hub](https://artifacthub.io) or [Bitnami](https://bitnami.com) and import them into Meshery as designs.

## Prerequisites

Before starting, ensure you have:

1. [Helm CLI](https://helm.sh/docs/intro/install/) installed
2. [mesheryctl](https://docs.meshery.io/installation) installed  
3. [Meshery](https://docs.meshery.io/installation/local) running locally

## Step 1: Download Helm Chart

### From Bitnami Repository

```bash
# Add Bitnami repository
helm repo add bitnami https://charts.bitnami.com/bitnami

# Update repositories
helm repo update

# Download a chart (example: nginx)
helm pull bitnami/nginx --untar

This creates a nginx directory containing the chart files.
From Artifact Hub

Visit artifacthub.io
Search for your desired chart (e.g., "nginx")
Click on the chart and scroll to find the Repository details
Note the repository URL and name

bash# Add the repository (example)
helm repo add stable https://charts.helm.sh/stable

# Update repositories
helm repo update

# Download the chart
helm pull stable/nginx-ingress --untar
Step 2: Import into Meshery
Method 1: Using mesheryctl (CLI)
Navigate to the downloaded chart directory and import:
bash# Change to chart directory
cd nginx

# Import as Meshery design
mesheryctl design import -f . -s helm-chart
Method 2: Using Meshery UI

Open Meshery UI (usually http://localhost:9081)
Navigate to Designs → Import
Drag & drop the chart directory or .tgz file
Select Helm Chart as the source type
Click Import

Step 3: Publish to Meshery Catalog

In Meshery UI, go to Designs
Find your newly imported design
Click Edit to add description, tags, and other metadata
Click Publish to submit for catalog approval

Once approved by Meshery maintainers, your design will be publicly available at meshery.io/catalog.
Complete Example Workflow
bash# 1. Setup and download
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm pull bitnami/nginx --untar

# 2. Start Meshery (if not running)
mesheryctl system start

# 3. Import to Meshery
cd nginx
mesheryctl design import -f . -s helm-chart

# 4. Open UI to publish
# Visit: http://localhost:9081






















IssueSolutionmesheryctl: command not found Install mesheryctl: `curl -sL https://git.io/mesheryImport failed: Meshery not runningStart Meshery: mesheryctl system startChart validation errorsRun helm lint ./chart-dir to validate
