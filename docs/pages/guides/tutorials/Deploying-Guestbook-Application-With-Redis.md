---
layout: default
title: Deploying PHP Guestbook application(Stateless) with Redis
abstract: |
  In this tutorial, you'll build and deploy a simple multi-tier web application using Kubernetes and Docker. It includes a Redis database for storing guestbook entries and multiple web frontend instances. You'll learn how to containerize these components, set up Kubernetes deployments and services, and establish communication between them. This hands-on guide is perfect for understanding basic application deployment in a clustered environment.
permalink: guides/tutorials/Deploying-Guestbook-Application-With-Redis
redirect_from: guides/tutorials/Deploying-Guestbook-Application-With-Redis
type: guide
category: tutorial
language: en
---
### Introduction

This tutorial shows you how to build and deploy a simple (not production ready), multi-tier web application using Kubernetes and Docker. This example consists of the following components:
- A single-instance Redis to store guestbook entries
- Multiple web frontend instance

### Prerequisites:
- Basic understanding of containerization , Kubernetes concepts and basic knowledge of reddis .
- Access to the _Me
- shery Playground_. If you don't have an account, sign up at [Meshery Playground](https://meshery.layer5.io/).

### lab scenario
In this lab, you will set up a Redis leader followed by two Redis followers to create a replicated Redis cluster. Next, you will deploy the guestbook frontend application that interacts with this Redis cluster. You will expose the frontend service to view and interact with the guestbook. Finally, after completing your tasks and testing, you will clean up by deleting the deployed resources to ensure no lingering components remain in your Kubernetes environment.

### Learning Objectives
In this tutorial, you will learn to deploy a Redis cluster with a leader and followers, understand Docker containerization for a guestbook frontend, configure Kubernetes deployments and services, expose services for external access, and practice resource cleanup with the help of meshery playground . These skills will provide a foundational understanding of managing distributed applications with Kubernetes and Docker in a clustered environment.

### Steps:

#### Access Meshery Playground
   - Log in to the [Meshery Playground](https://meshery.layer5.io/) using your credentials. On successful login, you should be at the dashboard. Press the **X** on the _Where do you want to start?_ popup to close it (if required).
   - Click **Explore** in the Cloud Native Playground tile to navigate to _MeshMap_.

> **_NOTE:_**  MeshMap is still in beta.

#### Step 1- Start up the Redis Database

