---
layout: default
title: Hello Kubernetes with Dapr
abstract: |
  Dive into this step-by-step tutorial on deploying Dapr in Kubernetes using Meshery Playground. Learn to deploy a Python app for generating messages and a Node app for consuming and persisting them, all orchestrated with Dapr for enhanced application integration.
permalink: guides/tutorials/Hello-Kubernetes-with-Dapr
redirect_from: /dapr-quickstart/Hello-Kubernetes-with-Dapr
type: guide
category: tutorial
language: en
---

### Introduction :
This tutorial shows how to deploy Dapr in a Kubernetes cluster using Meshery Playground. You'll deploy a Python app for message generation and a Node app for consumption and persistence, following the classic Hello World example with Dapr. Let's dive in and explore the architecture diagram that illustrates this setup.

### Prerequisites:
- Basic understanding of containerization , Kubernetes concepts and basic knowledge of Dapr .
- Access to the _Meshery Playground_. If you don't have an account, sign up at [Meshery Playground](https://meshery.layer5.io/).

### Lab Scenario:
You will be deploying `Hello Kuberentes` applicationswith the help of `Meshery Playground` which has  Python App that generates messages and the Node app consumes , persists them.

### Learning Objective:
Upon completing this tutorial, you will:
- Deploy Hello Kubernetes with  Dapr using Meshery Playground.
- Implement a Python app for message generation and a Node app for consumption and persistence with Dapr.
- Explore Dapr's features for service invocation, state management, and application observability in Kubernetes.

### Steps:

#### Access Meshery Playground
   - Log in to the [Meshery Playground](https://meshery.layer5.io/) using your credentials. On successful login, you should be at the dashboard. Press the **X** on the _Where do you want to start?_ popup to close it (if required).
   - Click **Explore** in the Cloud Native Playground tile to navigate to _MeshMap_.

> **_NOTE:_**  MeshMap is still in beta.




