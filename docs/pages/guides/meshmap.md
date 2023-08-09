---
layout: default
title: Using MeshMap
description: MeshMap allows us to interact with components and create designs with them
permalink: guides/using-meshmap
type: Guides
language: en
---

MeshMap is a **GitOps-infused** cloud native visual designer for Kubernetes and cloud native applications. Meshmap allows us to interact with components and create designs with components.

MeshMap is an extension of Meshery. It can also be accessed in extension panel in Meshery UI.


### Extensions Menu 


In the Dashboard, the "Extensions" are available via a sidebar option, located second to last in the navigation menu. To access the extensions, click on this option, and the extension panel will appear on your screen. Here, you will find all the extensions offered by Meshery. Click on "Sign Up" for Meshery to proceed with the registration process.


### Modes of MeshMap


MeshMap has two modes: 1. [Designer](https://layer5.io/cloud-native-management/meshmap/design) 2. [Visualizer](https://layer5.io/cloud-native-management/meshmap/visualize)


### **Designer**

Design, deploy, and manage your Kubernetes-based, cloud native deployments.


**Step 1. Access the MeshMap in Meshery UI**


In the Dashboard, the MeshMap visualization tool is accessible via a sidebar option located at the bottom of the navigation menu. To access MeshMap, click on this option, and the MeshMap extension will activate on the screen. From there, you can select the "design" option.


**Step 2. Create a New Design**


To create a new design click on the pencil button and give the design a name accordingly. With each change in the design, the design will automatically get saved.


**Step 3. Creating the design**


Drag and Drop the model the model from the configuration panel. Connect them by dragging the cursor from one component to another. Change the configuration of the model by left click, the settings menu with respect to that model will appear. It can be configured from settings panels, information regarding all the components can be accessed from the [MeshModel Summary Browser]({{site.baseurl}}/guides/accessing-meshmodel).

To access delete, copy, duplicate option right click and drag the cursor towards the option.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-component-option.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-component-option.png" /> </a>


### Other MeshMap's features


There are guides available for MeshMap features:

- Publishing a Design - <a href="{{site.baseurl}}/extensions/publishing-a-design"> Pubish Design</a>
- Merging Designs - <a href="{{site.baseurl}}/extensions/merging-design"> Merge Designs </a>
- Sharing a Design - <a href="{{site.baseurl}}/extensions/sharing-a-design"> Share Design </a>
- Importing a Design - <a href="{{site.baseurl}}/extensions/importing-a-design"> Import Design </a>
- Importing an Application - <a href="{{site.baseurl}}/extensions/importing-an-application"> Import an Application </a>


### **Visualizer**

Visualize Mode manages and operates your deployments and services in real-time.


**Step 1. Access the MeshMap in Meshery UI**


The option to access the MeshMap is available as the last option in the sidebar of Meshery UI. Upon clicking the icon, the extension MeshMap will enable on your screen, and then choose the visual option. The visualize mode of MeshMap will activate on your screen.


<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-desktop.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-desktop.png" /> </a>


**Step 2. Connecting Kubernetes Cluster**


Click the Connect Cluster button and upload the `config` file normally present in `/.kube/` directory. After uploading the file, go back to the MeshMap extension visualize mode, where you would be able to visualize your cluster.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-screen.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-screen.png"/> </a>


**Step 3. Working with visualizer mode**


You can visualize your replicaset, pods, deployment, secrets on the canvas. Also, they can be filtered as per their namespaces. On clicking an item in the canvas, its details will appear in the details section on the right. Additionally, many other options, like opening an "Interactive Terminal (Text-based interface for executing commands)" and "Stream Container Logs (Real-time display of container logs)" for a pod and container, are provided under the "Actions" button in the right panel.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual.png"/> </a>

Click on "Interactive Terminal" button, after clicking the Open Interactive Terminal button, an interactive terminal for the selected componenet will get open. 

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-terminal.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-terminal.png"/> </a>


To learn more about working with designs, click here: <a href="{{site.baseurl}}/extensions/meshmap"> MeshMap - Meshery Extensions </a>
