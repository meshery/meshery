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


### Access the Extensions


In the Dashboard, the "Extensions" are available via a sidebar option, located second to last in the navigation menu. To access the extensions, click on this option, and the extension panel will appear on your screen. Here, you will find all the extensions offered by Meshery. Click on "Sign Up" for Meshery to proceed with the registration process.


### Modes of MeshMap


MeshMap has two modes:

###### 1. [Designer](https://layer5.io/cloud-native-management/meshmap/design)
###### 2. [Visualizer](https://layer5.io/cloud-native-management/meshmap/visualize)


### **Designer**


**Step 1. Access the MeshMap in Meshery UI**


In the Dashboard, the MeshMap visualization tool is accessible via a sidebar option located at the bottom of the navigation menu. To access MeshMap, click on this option, and the MeshMap extension will activate on the screen. From there, you can select the "design" option.


**Step 2. Create a New Design**


To create a new design click on the pencil button and give the design a name accordingly.


**Step 3. Creating the design**


You can create the design by dragging and dropping the model. You can simply drag and drop the model on the MeshMap Canvas using the configuration panel, and while creating the design, it gets automatically saved with each change. You can also connect two components by dragging your cursor from one component to another.


**Step 4. Setting Configuration of model**


To configure it left-click on the model. The settings menu, would appear where all there configuration options are present. And to delete the model right click on the component and drag the cursor to the delete option. The information about the configuration of the model and its component can also be accessed in [MeshModel Summary Browser]({{site.baseurl}}/guides/accessing-meshmodel).

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-component-option.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-component-option.png" /> </a>


### Guide to use MeshMap's features


There are guides available for MeshMap features:

- Publishing a Design - <a href="{{site.baseurl}}/extensions/publishing-a-design"> Pubish Design</a>
- Merging Designs - <a href="{{site.baseurl}}/extensions/merging-design"> Merge Designs </a>
- Sharing a Design - <a href="{{site.baseurl}}/extensions/sharing-a-design"> Share Design </a>
- Importing a Design - <a href="{{site.baseurl}}/extensions/importing-a-design"> Import Design </a>
- Importing an Application - <a href="{{site.baseurl}}/extensions/importing-an-application"> Import an Application </a>


### **Visualizer**


**Step 1. Access the MeshMap in Meshery UI**


The option to access the MeshMap is available as the last option in the sidebar of Meshery UI. Upon clicking the icon, the extension MeshMap will enable on your screen, and then choose the visual option. The visualize mode of MeshMap will activate on your screen.


<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-desktop.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-desktop.png" /> </a>


**Step 2. Connecting Kubernetes Cluster**


Click the Connect Cluster button and upload the `config` file normally present in `/.kube/` directory. After uploading the file, go back to the MeshMap extension visual mode, where you would be able to visualize your cluster.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-screen.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-screen.png"/> </a>


**Step 3. Working with visualizer mode**


You can visualize your replicaset, pods, deployment, secrets on the canvas. Also, they can be filtered as per their namespaces. On clicking an item in the canvas, its details will appear in the details section on the right. Additionally, many other options, like opening an interactive terminal for a pod and container, are provided under the "Actions" button in the right panel.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual.png"/> </a>

Click on "Interactive Terminal" button, after clicking the Open Interactive Terminal button, an interactive terminal for the selected componenet will get open.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-terminal.png"> <img style="border-radius: 0.5%;" style="width:700px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-terminal.png"/> </a>


To learn more about working with designs, click here: <a href="{{site.baseurl}}/extensions/meshmap"> MeshMap - Meshery Extensions Documentation </a>