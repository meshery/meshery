---
layout: default
title: Using Meshmap
description: Meshmap allows us to interact with components and create designs with them
permalink: guides/meshmap
type: Guides
language: en
---

MeshMap is a **GitOps-infused** cloud native visual designer for Kubernetes and cloud native applications. Meshmap allows us to interact with components and create designs with components.

MeshMap has two modes:

1. Designer
2. Visualizer

### Designer

**Step 1. Access the MeshMap in Meshery UI.**

The option to access the MeshMap is present at last option inside the sidebar of Meshery UI. Upon clicking the icon, the extension MeshMap will enable at your screen and then choose design option.

**Step 2. Create a New Design**

You can create your first design and give it a name. The design that you will be creating to create will be saved automatically.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-creation.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-creation.png"/> </a>

**Step 3. Creating the design**

You can create the design by dragging and droping the component. You can simply drag the component and drop the components on the MeshMap Canvas using the configuration panel, and while creating the design the design gets automatically saved with each change. You can also connect two componets by dragging your cursor from one component to another.

**Step 4. Setting Configuration of Component**

You can left click on each component to access configuration and settings for each components. And to delete the component you need to do a right click on the compoent. And select the delete option.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-component-option.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-component-option.png"/> </a>

You can also merge two or more designs by dragging and dropping them on the canvas.

### Visualizer

**Step 1. Access the Meshmap in Meshery UI**

The option to access the MeshMap is present at last option in the sidebar of Meshery UI. On clicking the icon, the extension MeshMap will enable at your screen and then choose visual option. The visualize mode will pop up on your screen.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-desktop.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-desktop.png" /> </a>

**Step 2. Connecting Kubernetes Cluster**

Click Connect Cluster button and upload the `config` file normally present in `/.kube/` directory. Now you can visualize the your kubernetes cluster on the screen.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-screen.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-screen.png"/> </a>

Now you can visualize the cluster in the workspace.

**Step 3. Working in visualizer mode**

You can visualize you replicaset, pods, deployment, secrets in the canvas, also these can be filtered as per the namespaces. On clicking an item in the canvas its detail will appear in the details section in right, also many other option like opening an interactive terminal are provided for an pod and container under Actions button in right panel.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual.png"/> </a>

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-terminal.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-terminal.png"/> </a>

