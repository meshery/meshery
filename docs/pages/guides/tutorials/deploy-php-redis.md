---
layout: default
title: Deploying PHP Guestbook application with Redis in Meshery
abstract: Learn how to deploy a PHP Guestbook application with Redis on Kubernetes using Meshery.
permalink: guides/tutorials/deploy-php-redis
redirect_from: guides/tutorials/deploy-php-redis/
model: kubernetes
kind: deployments
type: guides
category: tutorials
language: en
list: include
abstract: "In this tutorial, we will deploy guest book built with PHP and Redis in Kubernetes. We will make use of Meshery Playground in an interactive live cluster environment."
---

### Introduction

In this tutorial, we will learn how to deploy a **PHP Guestbook application with Redis** using Meshery Playground. Meshery Playground is an interactive and collaborative live cluster environment that simplifies the deployment process and enhances user experience by providing visual tools for managing Kubernetes resources.

> **NOTE:** If this is your first time working with Meshery Playground, consider starting with the [Exploring Kubernetes Pods with Meshery Playground](https://docs.meshery.io/guides/tutorials/kubernetes-pods) tutorial first.

### Prerequisites
- Basic understanding of Kubernetes concepts.
- Meshery Playground access. If you don't have an account, sign up at [Meshery Playground](https://play.meshery.io).

### Lab Scenario
- Import the PHP and Redis manifest files into Meshery Playground.
- Deploy these resources on the playground.
    - Start up a Redis leader.
    - Start up two Redis followers.
    - Start up the guestbook frontend.
- Expose and view the Frontend Service.


### Objective
Learn how to import manifest files, visualize Kubernetes resources, create new resource components, and deploy the application using Meshery Playground.


### Steps

#### Download the Kubernetes Configuration Files

To get started we will need some yaml files that will contain the configurations for the pods and services that will run on Meshery. You can download them to see what the look like.
1. [redis-leader-deployment.yaml](https://k8s.io/examples/application/guestbook/redis-leader-deployment.yaml)
2. [redis-leader-service.yaml](https://k8s.io/examples/application/guestbook/redis-leader-service.yaml)
3. [redis-follower-deployment.yaml](https://k8s.io/examples/application/guestbook/redis-follower-deployment.yaml)
4. [redis-follower-service.yaml](https://k8s.io/examples/application/guestbook/redis-follower-service.yaml)
5. [frontend-deployment.yaml](https://k8s.io/examples/application/guestbook/frontend-deployment.yaml)
6. [frontend-service.yaml](https://k8s.io/examples/application/guestbook/frontend-service.yaml)

These YAML files contain the Service definitions and Deployment configurations for the PHP app with Redis.

#### Accessing Meshery Playground

1. Log in to the [Meshery Playground](https://play.meshery.io) using your credentials. On successful login, you should be at the dashboard. Press the **X** on the _Where do you want to start?_ popup to close it (if required).

2. Click **Explore** in the Cloud Native Playground tile to navigate to _MeshMap_


#### Import the Files to Meshery Playground

1. In the left sidebar, click on the upward arrow symbol(import icon) to import the designs into Meshery.

2. On doing so, a modal appears.In the modal that appears, enter a name for your design in the "Design File Name" field (e.g.`redis-leader-deployment`).

3. Now select `Kubernetes Manifest` from the "Design Type" dropdown menu.

4. Then choose `URL Import` for the upload method.

5. Now input the URL.

6. Then, click on `Import`

![Import redis-leader-deployment](./screenshots/redis-leader-deployment.png)


Now, follow the same steps (1-6) to import the rest of the files.

Under the “Designs” tab, you will see that we have successfully imported the manifest as designs, then you can drag and drop them in the canvas. This will "Merge" all the designs since it's all just one application

![merging all designs](./screenshots/app-canvas.png)

#### 4. **Deploy the files:**

1. Click Actions in the top right corner and click on Deploy.
![Import redis-php-guestbook-deploy](./screenshots/redis-php-guestbook-deploy.png)

1. The design will be validated to make sure there are no errors.
![Import redis-php-guestbook-validate](./screenshots/redis-php-guestbook-validate.png)
1. Choose the Kubernetes cluster you want to deploy to.
![Import redis-php-guestbook-identify](./screenshots/redis-php-guestbook-identify.png)
1. A Dry-Run will be triggered.
![Import redis-php-guestbook-dryrun](./screenshots/redis-php-guestbook-dryrun.png)
1. Finally your deployment and click Deploy to deploy the application to the cluster.
![Import redis-php-guestbook-finalize](./screenshots/redis-php-guestbook-finalize.png)
1. On successful deployment you will see the following modal asking yout o open you deployment in visualizer
![Import redis-php-guestbook-finalize](./screenshots/redis-php-guestbook-finalize.png)


#### 5. **Updating the Deployments:**

```bash
kubectl set image deployment/<deployment-name> <container-name>=new-image:tag
```


#### 6. Open the files in Operate

1. Once deployment is successful user can click on Open in Operate, or click on notification on top right click on the deployment successful notification and click on Open in Operate, or directly click on the Operate tab beside Design on the design and visualize the design inside operate mode.
![Operate](./screenshots/operate.png)

2. Click on the Layers option below to view all available filters.
![Operate layers](./screenshots/operate-layers.png)
Now you can utilize this filter to visualize your design.
![Operate layers expanded](./screenshots/operate-layers-expand.png)

3. After selecting your filters, you should see a view displaying only your relevant resources.


#### 7. Using built-in terminal for logs

1. Right click on a pod, to open the circular-content menu.
![Operate logs menu](./screenshots/operate-logs-menu.png)

2. Select the logs option from circular-context menu, this will start the session for logs

3. Now you can use the in built terminal to view logs
![Operate logs ](./screenshots/operate-logs.png)


#### 8. Deleting the resources

To delete the resources, use the **Undeploy** option from the _Design_ view.


### Conclusion
Congratulations! You've successfully completed the lab on exploring Kubernetes Deployments using Meshery Playground. This hands-on experience has equipped you with practical knowledge on deploying, updating, and monitoring applications in a Kubernetes environment. Continue exploring more scenarios in the Meshery Playground to enhance your skills in container orchestration.