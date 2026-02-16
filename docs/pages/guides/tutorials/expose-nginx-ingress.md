---
title: Expose a Kubernetes App with NGINX Ingress using Meshery Playground
author: Sapna Mohanta
description: Step-by-step tutorial to expose an NGINX app with Kubernetes Deployment, Service, and Ingress using Meshery Playground.
tags: 
  - kubernetes
  - ingress
  - nginx
  - tutorial
  - meshery
---


In this tutorial, we will learn how to expose a simple NGINX application using **Kubernetes Ingress** with the help of **Meshery Playground**.  
By the end of this tutorial, you will have:  
- A **Meshery Design** published to the Catalog (Deployment + Service + Ingress).  
- A working **Ingress rule** to route external traffic to your NGINX app.  
- A step-by-step guide with validation screenshots.  

---

## Step 1 — Open Meshery Playground
- Go to [Meshery Playground](https://playground.meshery.io).  
- Create a new design.  

---

## Step 2 — Add Deployment
- Drag a **Deployment** resource onto the Kanvas.  
- Configure it:
  - **Name:** `nginx-deployment`
  - **Image:** `nginx:latest`
  - **Replicas:** `1`
  - **Port:** `80`

📸 *Deployment created*

---

## Step 3 — Add Service
- Drag a **Service** resource onto the Kanvas.  
- Configure it:
  - **Name:** `nginx-service`
  - **Type:** `ClusterIP`
  - **Port:** `80`
  - **Target Port:** `80`  

- Connect the Service to the Deployment.  

📸 *Service connected to Deployment*

---

## Step 4 — Add Ingress
- Drag an **Ingress** resource onto the canvas.  
- Configure it:
  - **Name:** `nginx-ingress`
  - **Host:** `example.com`
  - **Path:** `/`
  - **Service Name:** `nginx-service`
  - **Service Port:** `80`  

- Connect the Ingress to the Service.  

📸 *Ingress configured and connected*

---

## Step 5 — Save the Design
- Save the design and name it something like:  
  **Expose NGINX with Ingress — Tutorial**

---

## Step 6 — Publish the Design
1. From **Meshery UI → Designs**, select your design.  
2. Click **Publish**.  
3. Fill in the details:
   - **Title:** Expose a Kubernetes App with NGINX Ingress  
   - **Description:** A Meshery Playground design that demonstrates exposing an NGINX app using Deployment, Service, and Ingress.  
   - **Tags:** `tutorial, kubernetes, nginx, ingress, meshery`  
4. Mark it as **Public** and click **Publish**.  

👉 [View this Design in Meshery Catalog](https://playground.meshery.io/extension/meshmap?mode=design&design=8b3a14d0-b822-4c0d-8051-571032558408)
📸 ![Sapna Mohanta – Published Catalog Design](/assets/img/tutorials/sapna-catalog.png)

---

## Step 7 — Validate
- Use the Meshery Playground simulation to validate that traffic flows from Ingress → Service → Deployment → Pod (NGINX).  
- You should see the app exposed successfully.  

📸 *Validation result*

---

## Conclusion
Congratulations 🎉  
You have successfully created and exposed an NGINX application using Kubernetes **Deployment**, **Service**, and **Ingress** — all through Meshery Playground.  

Your design is now **published to the Meshery Catalog** and can be reused by others.  
