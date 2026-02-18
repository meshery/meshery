---
title: Exploring Kubernetes Services with Meshery
params:
  model: kubernetes
  kind: services
categories: [tutorials]
description: Explore Kubernetes Services using Meshery using Meshery Playground, an interactive live environment, through a series of hands-on exercises.

---

### Introduction

In this tutorial, we'll learn to implement Kubernetes **Services**, the resources responsible for exposing applications inside and outside the cluster. Using Meshery Playground, an interactive live cluster environment, we'll perform hands-on labs to gain practical experience with the Kubernetes Services, without writing any YAML.

> **NOTE:** If this is your first time working with Meshery Playground, consider starting with the [Exploring Kubernetes Pods with Meshery Playground](/guides/tutorials/kubernetes-pods) tutorial first or [Exploring Kubernetes Deployments with Meshery](/guides/tutorials/kubernetes-deployments).


### Prerequisites

- Basic understanding of containerization and Kubernetes concepts.
- Access to the _Meshery Playground_. If you don't have an account, sign up at [Meshery Playground](https://playground.meshery.io/).

### Lab Scenario

Deploy a simple application or a simple Pod and expose it using ClusterIP, NodePort, and LoadBalancer services. Each service will route traffic to the pods via a common label selector. We’ll inspect the service details in Operator mode to confirm their types and behavior.

### Objective

Learn how to create, manage, and explore _Kubernetes Services_ to expose applications within the context of a microservices architecture.

### Steps

#### Access Meshery Playground
- Log in to the [Meshery Playground](https://playground.meshery.io/) using your credentials.  
- On successful login, you should be at the dashboard.
- Click **Kanvas** from the left menu to navigate to the [_Kanvas_ design](https://kanvas.new/extension/meshmap) page.

  ![](/guides/tutorials/images/kubernetes-deployments/2025-02-27_16-59.png)

> **_NOTE:_** Kanvas is still in beta.


#### Create a Deployment

1. In the _Kanvas Design_ page, start by renaming the design to a name of your choice for easier identification later.

    ![](/guides/tutorials/images/kubernetes-services/2025-09-04_02.png)

2. From the floating dock below, click the **Kubernetes** icon and then click **Deployment** from the list. This will create the _Deployment_ component on the design canvas. 

    ![](/guides/tutorials/images/kubernetes-services/2025-09-04_03.png)

3. Click or Drag the _Deployment_ component onto the canvas and the **Configure** tab automatically opens.

    ![](/guides/tutorials/images/kubernetes-services/2025-09-04_04.png)
    
4. Change the **Name** of the deployment and the **Namespace** if required. For this demonstration, we will leave them as they are and deploy this to the _default_ namespace.
    
5. Set **Replicas** to `2`. Under **Selector** and **MatchLabels**, Set a _matchLabel_ pair. Here we have set `app:9988110`.
  ![](/guides/tutorials/images/kubernetes-services/2025-09-04_05.png)

6. Under **Template → Metadata → Labels**, add the same label `app:9988110`. 
  ![](/guides/tutorials/images/kubernetes-services/metadata.png)


7. While still under **Template**, click **Spec** to load the _spec_ configuration modal. Then scroll down and click **+ Add Item** next to **Containers**. This will create a container **Containers 1**. Click on it and add:  
- **Image**: `meshery/meshery-milestone:latest`  
- **Name**: `meshery-milestone`  
  ![](/guides/tutorials/images/kubernetes-services/2025-09-04_06.png)

8. Click outside to close the modal. The deployment is now ready and it will look similar to this:
  ![](/guides/tutorials/images/kubernetes-services/2025-09-04_07.png)


9. Validate and Deploy the design: Click Validate (**Actions** toolbar), ensure that there are no errors and then click Deploy. Wait for the deployment to complete (Notifications will appear on bottom right).

##### You have now deployed a Deployment with 2 pods running in the cluster.

---

#### Add a ClusterIP Service

1. From **Components**, search for **Service** and drag it to the canvas, rename the service, here I will go with `service-clusterip`. Click on the service component to open its config modal. 
  ![](/guides/tutorials/images/kubernetes-services/2025-09-04_08.png)

2. In the service configuration modal:  
- Set **Type** to `ClusterIP`.  
- Click on  **+ Add Item**  under Ports to add a port called **Ports 1**. Click on it and add: 
  - **Port**: `80`  
  - **TargetPort**: `80` (These match the container port that our image serves on.)
- Also add the same key value pair as before under **Selector**: `app:9988110`
- We will also add the same label as the deployment for easier identification in Operator Mode.
  ![](/guides/tutorials/images/kubernetes-services/edit-01.png)


3. Connect the Service to the Deployment: Click over the service component until green dots appear, click the arrow and select network. Drag to the deployment. This creates a Network link.  
  ![](/guides/tutorials/images/kubernetes-services/2025-12-13_1.png)

From the Actions Tab, Undeploy the deployment first and then, validate and dry-run the new design, resolve any errors that may arise. Now, deploy the design. A pop up in the bottom right will confirm that the design is successfully configured.

  ![](/guides/tutorials/images/kubernetes-services/2025-12-13_2.png)

Switch to Operator mode, explore the Service details. Select the `service-clusterip` resource to see its details. 
  ![](/guides/tutorials/images/kubernetes-services/2025-09-05_11.png)

Notice the ClusterIP listed under Addresses and that no external IP or NodePort is assigned. This confirms that a ClusterIP service provides an internal IP reachable only within the cluster.
This Service has a ClusterIP (10.98.146.20) and a selector (app=9988110). Any Pod with that label automatically becomes part of the Service’s backend. This label-to-Pod binding is how a ClusterIP Service internally routes traffic to its backing workloads.
  ![](/guides/tutorials/images/kubernetes-services/2025-09-05_12.png)


---

#### Add a NodePort Service
 
To allow external access, we’ll use a NodePort service. For simplicity, I will switch from using deployment to Pod for our next Service. 

1. Back in Design mode, we will drag a Pod from the dock onto the canvas. Scroll down within the Pod configuration modal to the Containers section. Click **+** to add a container. Expand **Containers-1**. Next, fill out some of the required container specifications. Start by entering the container image, we will use _meshery/meshery-milestone:latest_ for this exercise. Give the container a name and a unique label (This unique label will be used by the service selector.). 
  ![](/guides/tutorials/images/kubernetes-services/2025-09-06_13.png)

2. Now, drag a Service component onto the canvas and rename it to `service-nodeport`.

3. Under the config modal, set **Type** to `NodePort` and the same selector as the Pod label, so that our Service is able to connect with our Pod.
  ![](/guides/tutorials/images/kubernetes-services/2025-09-06_14.png)

4. Click on **+ Add Item** under Ports to reveal **Ports 1**, expand **Ports 1** and add: 
- **Port**: `80`  
- **TargetPort**: `80`  
- **NodePort**: `30091` (or leave blank to auto-assign). 
  ![](/guides/tutorials/images/kubernetes-services/2025-09-06_15.png)


5. Validate and deploy from the Action tab at the top right.
  ![](/guides/tutorials/images/kubernetes-services/2025-12-13_3.png)

Now switch to Operator mode, click on any component to view details(like type or selector) about the Service or the Pod.
  ![](/guides/tutorials/images/kubernetes-services/2025-12-13_4.png)

Note that this service is mapped NodePort and is accessible on the **Node’s IP address**.

Expand the details section and you will see a NodePort value (30091), this means the service is exposed on each Node’s IP at port 30091. You can access the app externally via `http://<NodeIP>:30091`.  
  ![](/guides/tutorials/images/kubernetes-services/2025-09-06_17.png)


The Operator mode also provides an interactive terminal, along with other Details about the Pod. Click on the Pod to reveal the `Initiate Terminal Session` option.
  ![](/guides/tutorials/images/kubernetes-services/2025-12-13_5.png)
 
> **_NOTE:_** In Meshery Playground, Node IPs may not be directly reachable from your local machine due to the sandboxed environment. The NodePort value confirms that the Service is correctly exposed by Kubernetes, even if direct browser access is restricted.

---


#### Add a LoadBalancer Service

Finally, we’ll create a LoadBalancer service. In a real cloud environment, this would provision an external load balancer. In Meshery Playground, you will see how the service object is defined, even though a real cloud IP isn’t provided.

1. In Design mode, add another Pod to the canvas (as before). Add a container with Name `meshery-milestone`, Image `meshery/meshery-milestone:latest`, and add label `app:8080` (or any unique label).
2. Drag a Service component onto the canvas, rename it `service-loadbalancer`. 
3. In the service’s Configure panel, set **Type** to `LoadBalancer`. Under Selector, add `app:8080`.
4. Under Ports, click **+ Add Item**. Expand **Ports 1** and set: 
- **Port**: `80`  
- **TargetPort**: `80`  


5. Close the panel. 
  ![](/guides/tutorials/images/kubernetes-services/2025-12-13_6.png)


6. Validate and Deploy (undeploy the old design first).  

In Operator mode, observe the LoadBalancer service. In a real Kubernetes environment, a LoadBalancer provides an **external IP address**. 
  ![](/guides/tutorials/images/kubernetes-services/2025-12-13_7.png)

---

### Cleaning Up

To remove all the resources you created in this tutorial: in Design mode, go to Actions → Undeploy and confirm. This will delete the Deployments/Pods and Services from the cluster.

---

### Conclusion

Congratulations! You've successfully completed the lab on exploring Kubernetes Services with Meshery Playground. You created and deployed a sample application, then exposed it with different Service types (ClusterIP, NodePort, LoadBalancer).

Continue exploring more scenarios in the Meshery Kanvas to enhance your skills.


---
