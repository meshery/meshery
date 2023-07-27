---
layout: default
title: Using MeshMap
description: MeshMap allows us to interact with components and create designs with them
permalink: guides/meshmap
type: Guides
language: en
---

MeshMap is a **GitOps-infused** cloud native visual designer for Kubernetes and cloud native applications. Meshmap allows us to interact with components and create designs with components.

MeshMap has two modes:

1. Designer
2. Visualizer

### Designer

**Step 1. Access the MeshMap in Meshery UI**

The option to access the MeshMap is available as the last option inside the sidebar of Meshery UI. Upon clicking the icon, the extension MeshMap will enable on your screen, and then you can choose the "design" option.

**Step 2. Create a New Design**

You can create your first design and give it a new name.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-creation.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-creation.png"/> </a>

**Step 3. Creating the design**

You can create the design by dragging and dropping the components. You can simply drag and drop the component on the MeshMap Canvas using the configuration panel, and while creating the design, it gets automatically saved with each change. You can also connect two components by dragging your cursor from one component to another.

**Step 4. Setting Configuration of Component**

You can left-click on each component to access configuration and settings for each component. To delete the component, you need to right-click on it and select the delete option.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-component-option.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-design-component-option.png"/> </a>

You can also merge two or more designs by dragging and dropping them on the canvas.

### Visualizer

**Step 1. Access the MeshMap in Meshery UI**

The option to access the MeshMap is available as the last option in the sidebar of Meshery UI. Upon clicking the icon, the extension MeshMap will enable on your screen, and then choose the visual option. The visualize mode will pop up on your screen.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-desktop.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-desktop.png" /> </a>

**Step 2. Connecting Kubernetes Cluster**

Click the Connect Cluster button and upload the `config` file normally present in `/.kube/` directory. You would be able to visualize it on your canvas.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-screen.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-screen.png"/> </a>


**Step 3. Working in visualizer mode**

You can visualize your replicaset, pods, deployment, secrets in the canvas. Also, these can be filtered as per the namespaces. On clicking an item in the canvas, its details will appear in the details section on the right. Additionally, many other options, like opening an interactive terminal for a pod and container, are provided under the "Actions" button in the right panel.

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual.png"/> </a>

<a href="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-terminal.png"> <img src="{{ site.baseurl }}/assets/img/meshmap/meshmap-visual-terminal.png"/> </a>

