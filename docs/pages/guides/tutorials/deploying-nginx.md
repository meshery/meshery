---
layout: default
title: Deploying NGINX with Meshery Playground
permalink: guides/tutorials/deploying-nginx
model: kubernetes
kind: deployment
type: guides
category: tutorials
language: en
list: include
abstract: "In this tutorial, we will explore how to use Meshery Playground, an interactive live cluster environment, to perform hands-on labs for deploying and managing an NGINX web server on Kubernetes."
---

Introduction:
In this tutorial, we will explore how to deploy and manage an NGINX web server using Meshery Playground. We will use the MeshMap Designer to design, validate, deploy, and undeploy an NGINX Deployment and Service, following the full Kubernetes workload lifecycle entirely within the Meshery Playground interface.

Prerequisites:
- Basic understanding of Kubernetes concepts.
- Meshery Playground access. If you don't have an account, you can sign up at [Meshery Playground](https://play.meshery.io).

Lab Scenario: Deploying and Managing an NGINX Web Server

Objective:
Learn how to use Meshery Playground to design an NGINX deployment, validate it, deploy it to the cluster, and cleanly undeploy it using the platform.

### Steps:

#### 1. **Accessing Meshery Playground:**
   - Log in to the [Meshery Playground](https://play.meshery.io) using your credentials.
   - Once logged in, navigate to the Meshery Playground dashboard.
   [![Dashboard]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/Dashboard.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/Dashboard.png)

#### 2. **Opening Kanvas Designer:**
   - Click the **Open Kanvas** button at the bottom right of the dashboard.
   - You will be taken to the Kanvas Designer canvas.
   [![Kanvas Empty]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/kanvas-empty.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/kanvas-empty.png)

#### 3. **Creating the NGINX Design:**
   - Click on **"Untitled Design"** at the top and rename it to `nginx-playground`.
   [![New Design]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/new-design.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/new-design.png)

#### 4. **Adding a Deployment Component:**
   1. Click the **Kubernetes wheel icon** in the bottom toolbar.
   2. Search for **"Deployment"** and click to place it on the canvas.
   [![Add Deployment]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/add-deployment.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/add-deployment.png)
   3. Click the Deployment component and configure it:
      - **Metadata → Name:** `nginx-deployment`
      - **Spec → Replicas:** `3`
      - **Spec → Template → Spec → Containers → Image:** `nginx:latest`
      - **Spec → Template → Spec → Containers → Name:** `nginx`
   [![Configure Deployment]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/deployment.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/deployment.png)

#### 5. **Adding a Service Component:**
   1. Search for **"Service"** in the toolbar and place it on the canvas.
   2. Configure the Service:
      - **Spec → Type:** `ClusterIP`
      - **Spec → Ports → Port:** `80`
   [![Add Service]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/service-add.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/service-add.png)
   [![Service Ports]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/ports-service.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/ports-service.png)

#### 6. **Saving the Design:**
   - Click the **cloud/save icon** at the top bar next to the design name.
   - Confirm the name `nginx-playground` and save.
   [![Save Design]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/save-design.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/save-design.png)
   [![Canvas Overview]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/canvas-overview.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/canvas-overview.png)

#### 7. **Deploying the Design:**
   1. Click the **Actions** dropdown at the top right.
   2. Select **Deploy**.
   [![Click Deploy]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/click-deploy.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/click-deploy.png)
   3. The deployment pipeline will validate your design and deploy it to the cluster.
   4. You will receive a confirmation notification at the bottom of the screen.
   [![Deployment Completed]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/deployment-completed.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/deployment-completed.png)

#### 8. **Verifying the Deployment:**
   1. Click the **Operate** tab at the top center of Kanvas.
   2. In the right panel, click **WORKLOADS** → **Deployments** to view running deployments in the cluster.
   [![Operator View]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/operator-view.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/operator-view.png)

#### 9. **Undeploying the Design:**
   1. Return to the **Design** tab.
   2. Click **Actions** → **Undeploy**.
   3. Confirm the undeploy action in the popup window to remove all NGINX resources from the cluster.
   [![Undeploy]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/undeploy.png)]({{site.baseurl}}/assets/img/tutorials/deploying-nginx/undeploy.png)

### Conclusion
Congratulations! You have successfully completed the lab on deploying NGINX using Meshery Playground. You have learned how to:

- ✅ Access and navigate Meshery Playground and Kanvas Designer
- ✅ Design an NGINX Deployment and Service using MeshMap Designer
- ✅ Deploy NGINX to the Kubernetes cluster via the Meshery Playground platform
- ✅ Verify the deployment using the Operate tab
- ✅ Cleanly undeploy all resources

Explore more scenarios in the Meshery Playground to enhance your skills in cloud-native technologies.
