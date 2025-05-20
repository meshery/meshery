---
layout: default
title: Component Shape Guide
permalink: extensions/component-shape-guide
language: en
abstract: Kubernetes architecture deployment and architecture diagramming tool for cloud native applications - Kanvas.
display-title: "false"
list: include
type: extensions
category: kanvas
---

## Component Shape Guide

As a best practice, users are encouraged to select an existing or create acustom shape for their components to best visually signify the function of their component.

Currently, the circle is used as the default shape for new components. However, if users or contributors have alternative shapes they believe better suit a particular component, they are encouraged to create new custom components.

<!-- Tool not ready for use

using the Shape Builder extension at [https://shapes.meshery.io](https://shapes.meshery.io).

-->

Although the usage of the components is divided into categories, some shapes serve as a universal representation of particular components.

Below are all the shapes with their current usage in a general context.



{% include extension-guide.html 
 data_file="shapes"
 guide_title="Shape"
 guide_description="Description"
 guide_usecase="CommonUseCase"
 guide_svg="SVG"
 guide_assests_folder="shapes"
%}
 
{% include alert.html type="light" title="Note" content="Because a Kubernetes Deployment can be a parent of other components, it displays as a rounded-rectangle (a limitation of the Javascript library being used to render Designs). There is a small badge, placed over the lefthand-side border that aids in identifying the component type." %}

Here's how the Deployment component appears in different contexts:

<br/>

<a href="../../../assets/img/deployment-shape.png">
    <img src="../../../assets/img/deployment-shape.png" style="width:50%; height:auto;" alt="Deployment Component Shape">
</a>
<p>Deployment component with its distinctive rounded-rectangle shape and badge</p>

<br/>

[![Deployment in Icon Set](../../../assets/img/deployment-icon.png)](../../../assets/img/deployment-icon.png)
Deployment icon in the component selection panel

<br/>

[![Deployment in Dashboard](../../../assets/img/deployment-dashboard.png)](../../../assets/img/deployment-dashboard.png)
Deployment component as seen in the cluster resource overview
