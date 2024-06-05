---
layout: tutorials
title: Importing Helm Charts into Meshery
abstract: In this tutorial, we will explore how to import  Helm charts using Meshery. Helm charts simplify the deployment and management of Kubernetes applications.
permalink: guides/tutorials/Importing_HelmCharts_into_Meshery/
redirect_from: guides/tutorial/Importing_HelmCharts_into _Meshery/
model: kubernetes
kind: helm-charts
type: guides
category: tutorials
language: en
list: include
published: true
abstract: >
  In this tutorial, we will guide you through the process of importing a Helm chart into Meshery. You will learn how to leverage Meshery's capabilities to manage your Kubernetes applications effectively using Helm charts.
---
### Introduction

Meshery is a powerful tool for managing and operating service meshes. One of its key features is the ability to import Helm charts, which simplifies the deployment and management of Kubernetes applications. This tutorial will guide you through the process of importing your first Helm chart into Meshery.

### Prerequisites

- Basic understanding of kubernetes and Helm.
- Access to the _Meshery Playground_. If you don't have an account, sign up at [Meshery Playground](https://play.meshery.io/).

### Lab Scenario: Importing Dapr Helm Chart into Meshery Playground

In this lab, we will walk through the process of importing the Dapr Helm chart into Meshery Playground. Dapr (Distributed Application Runtime) is a portable, event-driven runtime that makes it easy for developers to build resilient, stateless, and stateful microservices. By leveraging Meshery Playground, you will experience a hands-on approach to deploying and managing the Dapr Helm chart within a live Kubernetes environment.
For more information on Dapr, visit the [Dapr Documentation](https://docs.dapr.io/concepts/?_gl=1*1v6gt5w*_ga*MTc2MjAwNzU0OC4xNzE1NjA5MTE0*_ga_60C6Q1ETC1*MTcxNjM2MTUyOC4zLjAuMTcxNjM2MTUyOC4wLjAuMA..)

### Objective

Learn to efficiently import, configure, and deploy Helm charts within Meshery Playground, enhancing your understanding of managing Kubernetes applications using Helm.

### Steps

#### Access Meshery Playground

- Log in to the [Meshery Playground](https://meshery.layer5.io/) using your credentials. On successful login, you should be at the dashboard. Press the **X** on the _Where do you want to start?_ popup to close it (if required).
- Click **Explore** in the Cloud Native Playground tile to navigate to _MeshMap_.

> **_NOTE:_** MeshMap is still in beta.

#### Choosing Correct Dapr Helm charts

- Navigate to the Dapr Helm charts repository on GitHub or any other trusted source. Here is the one of the source from which you can pick [helm charts](https://github.com/jangocheng/dapr-helm-charts/tree/master).
  ![](./screenshots/H1.png)
- Ensure you select the appropriate chart version compatible with your Kubernetes cluster and Dapr runtime.
  ![](./screenshots/H2.png)
- Copy the repository URL for later use in Meshery.
  ![](./screenshots/H3.png)

### Note

To import Helm charts into Meshery, any direct downloadable link to a file ,follow this format: `https://raw.githubusercontent.com/username/repository/branch/path/to/file.tgz`.

Generally, we don't have GitHub raw links readily available, so you can use the [Git-Rawify tool](https://git-rawify.vercel.app/#convert) to convert standard GitHub URLs into raw links. Ensure you copy the file with the `.tgz` extension only.



### Importing Dapr Helm Chart to Meshery playground

1. Open the Meshery playground                                                                                                                                                                                    
## Note 
   Access to the _Meshery Playground_. If you don't have an account, sign up at [Meshery Playground](https://play.meshery.io/).

2. On the right side, click on the upward arrow symbol to import designs into  Meshery
   ![Import](./screenshots/H5.png)

   ## Note
   Importing Meshery designs can be done using several methods, including:

   - Importing through the Meshery UI configuration page.
   - Importing through the Meshery CLI.
   - Importing through the Meshmap(Which we are currently doing).
   
    For more detailed instructions, you can check out the Meshery documentation on importing designs: [Importing Designs](https://docs.meshery.io/guides/configuration-management/importing-designs).


4. After clicking on that, name this design as "Dapr Helm design".

   ![Name Design](./screenshots/H6.png)

5. In the dropdown on the left side, select "Helm Chart".

   ![Dropdown Selection](./screenshots/H7.png)

6. In the URL section, paste the existing GitHub raw URL link for the helm charts.

   ![URL Section](./screenshots/H8.png)

7. Click on "Import" to import the helm chart.

8. Wait for some time for the helm to be loaded into Meshery.

   ![Loading](./screenshots/H9.png)

9. Successfully, the Dapr helm chart is imported into Meshery.

   ![Success](./screenshots/H10.png)







  


