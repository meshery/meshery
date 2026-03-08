---
layout: default
title: Deploying, Scaling, and Updating NGINX with Meshery
abstract: Learn how to observe and manage a workload's lifecycle using Meshery and kubectl on a local Kubernetes cluster.
permalink: guides/tutorials/deploying-nginx
type: guide
category: tutorial
language: en
tags:
- kubernetes
- nginx
- workloads
- minikube
---

> **Note for Tutorial Users:** Reference screenshots are provided throughout this guide to verify your progress. If you're unsure whether your output matches expected results, refer to the [Screenshots Reference](#screenshots-reference) section at the bottom of this page.

## Introduction

This tutorial demonstrates how to deploy, scale, update, and rollback an NGINX web server using Meshery with a local Kubernetes cluster (minikube). 

### Why a Local Cluster?

You cannot complete this tutorial using the public Meshery Playground. The Playground sandbox restricts resource creation—you will encounter a **403 Forbidden** error for the `playground-user` service account when attempting to create new deployments or namespaces. A local Minikube or Kind cluster provides the necessary administrative access to deploy and manage resources.

### Why Use Terminal for State Changes?

Out of the box, the basic self-hosted Meshery installation serves as a powerful observation deck. Visual design tools like Kanvas (MeshMap) are premium extensions. Therefore, this tutorial teaches the standard real-world workflow: mutating state via the Kubernetes CLI (`kubectl`) and using the Meshery dashboard to monitor those changes in real time.

## Prerequisites

- A local Kubernetes cluster running (e.g., [Minikube](https://minikube.sigs.k8s.io/docs/start/))
- `mesheryctl` installed and configured
- `kubectl` configured to communicate with your local cluster
- Meshery running locally (via Docker)

## Lab Scenario: NGINX Deployment Lifecycle

**Objective:**  
Learn to deploy NGINX, scale it to handle load, perform rolling updates, and execute rollbacks using `kubectl` while observing changes through Meshery's dashboard.

---

### Step 1: Verify Cluster Connectivity

First, start your cluster and launch Meshery.

1. **Start your Minikube cluster:**

   ```bash
   minikube start
   ```

2. **Start Meshery using Docker:**

   ```bash
   docker run -d --name meshery \
     --network host \
     -e PROVIDER_BASE_URLS=https://meshery.layer5.io \
     -p 8080:8080 \
     layer5/meshery:stable-latest
   ```

3. **Access Meshery UI:**
   
   Open your browser and navigate to `http://localhost:8080`

4. **Verify Connection:**
   
   Navigate to **Connections** from the sidebar. You should see your `minikube` cluster listed with a green **Connected** badge.

   ![Minikube Connected](/assets/img/tutorials/deploying-nginx/1.png)

---

### Step 2: The Dashboard Overview

Navigate to the **Dashboard** via the left sidebar. This provides a high-level overview of everything running in your cluster. 

![Dashboard Overview](/assets/img/tutorials/deploying-nginx/2.png)

---

### Step 3: Deploy NGINX

Deploy an NGINX web server in the `meshery` namespace.

```bash
kubectl create deployment nginx-deployment --image=nginx:latest -n meshery
```

Navigate to **Dashboard** → **Workload** tab → **Deployment** sub-tab. You will see `nginx-deployment` with **1/1** replicas ready.

![Initial Deployment 1/1](/assets/img/tutorials/deploying-nginx/3.png)

---

### Step 4: Scale to 3 Replicas

Scale the deployment to handle more traffic:

```bash
kubectl scale deployment nginx-deployment --replicas=3 -n meshery
```

Refresh the **Deployment** tab. Meshery will now show **3/3** replicas ready.

![Scaled Deployment 3/3](/assets/img/tutorials/deploying-nginx/4.png)

---

### Step 5: Inspect Individual Pods

To see the individual instances backing your deployment:

1. Click the **Pod** sub-tab (next to Deployment).
2. Monitor the 3 NGINX pods running with `nginx:latest` image.

   ![Pods with Latest Image](/assets/img/tutorials/deploying-nginx/5.png)

---

### Step 6: Perform a Rolling Update

Update the NGINX image to the lighter `alpine` version with zero downtime.

1. **Keep the Pod tab open** in Meshery to observe the transition.

2. **Trigger the image update:**

   ```bash
   kubectl set image deployment/nginx-deployment nginx=nginx:alpine -n meshery
   ```

3. **Observe in Meshery:**
   
   Watch Kubernetes cycle the pods. New pods will spin up with `nginx:alpine` while old pods terminate. We can see the updated time saying 2 minutes while in step 5 it said 1 hour (which means they were updated correctly)
   ![Pods with Alpine Image](/assets/img/tutorials/deploying-nginx/6.png)

---

### Step 7: Verify the Update

Confirm the update was successful using Meshery's resource inspector.

1. Go back to the **Deployment** tab.
2. Click on `nginx-deployment` name.
3. Scroll down to the **Images** section. It should show `nginx:alpine`.

   ![Verify NGINX Alpine Image](/assets/img/tutorials/deploying-nginx/7.png)

---

### Step 8: Execute a Rollback

If an update causes issues, revert to the previous version.

1. **Run the rollback command:**

   ```bash
   kubectl rollout undo deployment/nginx-deployment -n meshery
   ```

2. **Verify in Meshery:**
   
   Check the Deployment Details page. The Images row will revert to `nginx:latest`.

   ![Rollback to NGINX Latest](/assets/img/tutorials/deploying-nginx/8.png)

---

### Step 9: Clean Up

Delete the resources when finished.

```bash
kubectl delete deployment nginx-deployment -n meshery
```

Check the Meshery Dashboard—the deployment will disappear from the workloads list.

---

## Conclusion

Congratulations! You have successfully walked through the core lifecycle of a Kubernetes workload:

- ✅ Deployed an NGINX web server
- ✅ Scaled it horizontally to 3 replicas
- ✅ Performed a zero-downtime rolling update
- ✅ Executed a rollback to previous version
- ✅ Cleaned up resources

By combining the mutation power of `kubectl` with the observability of the Meshery dashboard, you can confidently deploy, scale, and update applications while monitoring your infrastructure in real time.

