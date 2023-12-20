---
layout: default
title: KubeSphere
permalink: installation/kubesphere
type: installation
category: kubernetes
redirect_from:
- installation/platforms/kubershphere
display-title: "false"
language: en
list: include
image: /assets/img/platforms/kubesphere.png
abstract: Install Meshery on KubeSphere
---

{% include installation_prerequisites.html %}

[Meshery](https://meshery.io/) is the open source, cloud native management plane that enables the adoption, operation, and management of Kubernetes, any service mesh, and their workloads.

This tutorial walks you through an example of deploying Meshery from the App Store of KubeSphere.


## Prerequisites

- Please make sure you enable the OpenPitrix system.
- You need to create a workspace, a project, and a user account (`project-regular`) for this tutorial. The account needs to be a platform regular user and to be invited as the project operator with the `operator` role. In this tutorial, you log in as `project-regular` and work in the project `demo-project` in the workspace `demo-workspace`. For more information, see Create Workspaces, Projects, Users and Roles.


## Hands-on Lab

Perform the following steps in order:

### 1. <b>Deploy Meshery from the App Store</b>


1. On the **Overview** page of the project `demo-project`, click **App Store** in the upper-left corner.
2. Search for **Meshery** in the App Store, and click on the search result to enter the app.

    ![meshery-app]({{ site.baseurl }}/assets/img/platforms/meshery-app.png)
3. In the **App Information** page, click **Install** on the upper right corner.

    ![meshery-install]({{ site.baseurl }}/assets/img/platforms/Meshery-install.png)

4. In the App Settings page, set the application **Name**, **Location** (as your Namespace), and App Version, and then click Next on the upper right corner.

    ![meshery-info]({{ site.baseurl }}/assets/img/platforms/Meshery-info.png)

5. Configure the **values.yaml** file as needed, or click **Install** to use the default configuration.

    ![meshery-yaml]({{ site.baseurl }}/assets/img/platforms/Meshery-yaml.png)

6. Wait for the deployment to be finished. Upon completion, **Meshery** will be shown as **Running** in KubeSphere.

    ![meshery-app-running]({{ site.baseurl }}/assets/img/platforms/Meshery-app-running.png)



### 2. <b>Access the Meshery Dashboard</b>


1. Go to **Services** and click the service name of Meshery.
2. In the **Resource Status** page, copy the **NodePort** of Meshery.

    ![meshery-service]({{ site.baseurl }}/assets/img/platforms/Meshery-service.png)

3. Access the Meshery Dashboard by entering **${NodeIP}:${NODEPORT}** in your browser.

    ![meshery-dashboard]({{ site.baseurl }}/assets/img/platforms/meshery-dashboard.png)

{% include suggested-reading.html language="en" %}

{% include related-discussions.html tag="meshery" %}
