---
title: Exploring Kubernetes Pods with Meshery
model: kubernetes
params:
    kind: pods
categories: [tutorials]
description: Explore Kubernetes Pods using Meshery Playground in an interactive live cluster environment, through a series of hands-on exercises.
aliases:
- /guides/tutorials/kubernetes-pods
---

### Introduction

In this tutorial, we'll learn the fundamentals of Pods, the smallest deployable units in the Kubernetes ecosystem. Using Meshery Playground, an interactive live cluster environment, we'll perform hands-on labs to gain practical experience in deploying, managing, and understanding some of the concepts related to Pods.

> **_NOTE:_** This tutorial demonstrates stand-alone pods, i.e. pods not managed through deployments.

### Prerequisites

- Basic understanding of containerization and Kubernetes concepts.
- Access to the _Meshery Playground_. If you don't have an account, sign up at [Meshery Playground](https://play.meshery.io/).

### Lab Scenario

Deploy and explore an NGINX pod in a Kubernetes cluster. Additionally, expose the pod through a service.

### Objective

Learn how to create, manage, and explore _Kubernetes Pods and Services_ within the context of a microservices architecture.

### Steps

#### Access Meshery Playground

- Log in to the [Meshery Playground](https://play.meshery.io) using your credentials. On successful login, you should be at the dashboard. Press the **X** on the _Where do you want to start?_ popup to close it (if required).
- Click **Explore** in the Cloud Native Playground tile to open the visual configurator.

#### Create a simple stand-alone Pod

1. In the visual configurator, rename the design from _Untitled Design_ to a name of your choice. This helps in identifying the design later.
2. Click the **Components** tab.
3. Search for **Pod** in the list of components.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-22_18-20.png)
4. Scroll down, _select and drag_ the **Pod** component from the search results onto the design canvas.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-22_18-32.png)
5. You should now have a pod on the design canvas ready for configuration. It has a default name assigned to it.
    Click the pod component to open the configuration modal.
6. Rename the pod if necessary.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-22_18-35.png)

#### Add a container to the Pod

The pod at this stage does not have a container, so we will add one. This is equivalent to defining the containers in the _spec:_ section of a YAML manifest. For this exercise we will make only the basic required configurations.

7. Scroll down within the Pod configuration modal to the Containers section. Click **+** to add a container.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-23_11-54.png)
8. Expand **Containers-0**.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-23_11-54_1.png)
9. Fill in the required container specifications. Start by entering the container image - use _nginx:latest_ for this exercise.
10. Give the container a name.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-23_11-57.png)
11. Add a label by clicking the label icon on the modal.
12. Click **+** next to _Labels_.
13. Set a label of your choice. Use a unique label, since the playground is a shared environment.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-25_18-18.png)
14. Click outside to close the modal.

#### Validating and Deploying the Pod

Validate the design before deploying it.

Click **Validate** in the top toolbar.
![](/guides/tutorials/images/kubernetes-pods/2024-02-23_19-52.png)

A pop-up will appear showing details such as component count, annotations, and any errors.
![](/guides/tutorials/images/kubernetes-pods/2024-02-23_19-54.png)

Ensure there are no errors before proceeding. To deploy, click the **Deploy** (1) tab in the pop-up modal and then click **Deploy** (2) again.
![](/guides/tutorials/images/kubernetes-pods/2024-02-23_19-56.png)

You should see deployment alerts in the bottom-right corner.

#### Visualizing the Pod

To view the deployed resources, switch to the **Visualize** mode. Create a view with filters scoped to your resources.

1. Click **Visualize** to begin.
2. Give the view a name.
3. Click the filter icon.
4. Select appropriate filters. For this exercise, filter by the label set during pod creation to limit the view to your resources only.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-25_18-22.png)

Click the filter icon again to close. The view should now display only your resources, similar to the screenshot below:
![](/guides/tutorials/images/kubernetes-pods/2024-02-25_18-25.png)

#### Pod state and information

Select the **Details** tab, then select the pod in the view to inspect it.
The details tab shows:

- Pod state - a green check indicates **running**.
- The namespace where the pod is running.
- Container image details including image tag and name.
- Number of restarts and uptime.

![](/guides/tutorials/images/kubernetes-pods/2024-02-25_17-38.png)

#### Connecting to the Pods / containers

Meshery Playground provides terminal capabilities to connect to containers and stream logs. The following steps demonstrate this.

1. Select the **Pod** (or a specific container if the pod has multiple containers).
2. Select the **Actions** tab on the right.
3. Click **Open Interactive Terminal**.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-05.png)
    The terminal is displayed on screen.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-10.png)
4. To view streaming logs, click **Stream Container Logs**. A log view is added alongside the terminal.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-11.png)

#### Exposing a Pod with a Service

> **_NOTE:_** Meshery Playground currently supports NodePort only.

1. Return to the _Design_ view.
2. From **Components**, search for _service_ and drag the service component onto the design canvas.
3. Click the _service_ component to open its configuration modal.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-33.png)
4. Rename it if required.
5. Change _Type_ to **NodePort**.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-35.png)
6. Click **+** to add a port.
7. Expand **Ports-0** to configure the service port.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-37.png)
8. Set **Target Port** to 80.
9. Set **Port** to 80.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-45.png)
10. Click **+** next to Selector and add the key-value pair matching the label set on the pod (e.g. **app:9870**).
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-48.png)
11. Add the same key-value label to the service. This enables filtering in the visualization view.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-28_11-53.png)
12. Hover over the service component until green anchor points appear on all four sides. Hover over one anchor point to reveal an arrow, then click it. Select the **Network** relationship type and draw the connection to the pod.
    ![](/guides/tutorials/images/kubernetes-pods/2024-02-29_20-37.png)
13. Click **Validate** in the top toolbar to validate the updated design.

#### Updating the deployment

To apply the updated design (which now includes the service), undeploy and redeploy.

- Select **Undeploy** and confirm by clicking **Undeploy** again.
- Select **Deploy** and confirm by clicking **Deploy**.

#### Deleting and Recreating Pods

To delete deployed resources such as Pods, use the **Undeploy** option from the _Design_ view.

To recreate them, use the **Deploy** option.

### Conclusion

Congratulations! You've successfully completed the lab on exploring Kubernetes Pods using Meshery Playground. This hands-on experience should have provided valuable insights into the deployment, management, and interaction with Pods in a Kubernetes environment. Continue exploring more scenarios in Meshery Playground to deepen your skills.
