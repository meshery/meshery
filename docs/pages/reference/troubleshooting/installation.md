---
layout: default
title: Troubleshooting Meshery Installations
description: Troubleshoot Meshery installation and deployment
permalink: guides/troubleshooting/installation
type: Guides
category: troubleshooting
---

## Meshery's Preflight Checks

Anytime a `mesheryctl system` command is executed, a series of preflight checks are run. An attempt will be made to connect to the Kubernetes cluster configured in the user's kubeconfig as their current-context .

1. Check whether `mesheryctl` can initialize a Kubernetes client.

Situation: `mesheryctl` fails to query for pods in the default namespace of the user's current Kubernetes context.
