---
layout: default
title: MeshModel Summary
description: Accessing the MeshModel Summary in settings in Meshery UI
permalink: guides/meshmodel-summary
type: Guides
language: en
---


### **About MeshModel Summary**

The **MeshModel Browser Summary** exposes information about all the different types of infrastructure and service meshes that Meshery manages. Meshery models all the different types of infrastructures and service meshes that Meshery manages. It is quite convenient to look into the configuration of a Model in the Settings Menu while using them. Some of the core constructs are Components, Designs, Patterns, Traits, Metrics, Actions, Color, and Relationships. The feature offers a convenient way for administrators and developers to gain insights into their service mesh deployments at a glance. These model are registered by another systems. 


#### Component

Components should be Portable, Accessable and easily Managable. A user is able to create his own ComponentDefinition and publish it so that other can make a use of it. Components exposes the capabilities of the underlying platform. 

The core constructs of a Components are: 
1. Schema - the skeletal structure representing a logical view of shape, size, and characteristics of a construct.
2. Definition - implementation of the schema. 
3. Instance - An instantiation of the Definition. 
Inside components summary table, information regarding name, model name, apiVersion and subcategory is provided. 


#### Design 

Designs are how the users describe their infrastructure. A Design Schema consists of the components and patterns. No two components can have same name. 


#### Relationship

Relationships in MeshModel define how interconnected components interact with each other. They encompass different kinds of connections and dependencies among components, such as hierarchical, network, or default relationships. Each relationship is associated with selectors, metadata, and the possibility of having optional parameters.


### Accessing the Summary Table

**Step 1. Access the settings in Meshery UI**

Click on the Settings button in the Meshery UI to access settings menu.

<a href="{{ site.baseurl }}/assets/img/meshmodel/settings-meshmodal.png"><img alt="Settings" style="border-radius: 0.5%;" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmodel/settings-meshmodal.png" /></a>

**Step 2. Accessing Meshmodel Summary**

To view **MeshModel Summary**, select the **MeshModel Summary** tab. The summary will appear on your screen. You can acccess these models and their components through the [MeshMap extension]({{site.baseurl}}/extensions/meshmap).

<a href="{{ site.baseurl }}/assets/img/meshmodel/settings-meshmodel-summary.png"><img alt="MeshModel Summary" style="border-radius: 0.5%;" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/meshmodel/settings-meshmodel-summary.png" /></a>

You can visualize Models, Components, and Relationships in the Summary table. 
