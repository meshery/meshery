---
layout: default
title: Component Style Guide
permalink: extensions/component-style-guide
language: en
abstract: "A guide to understanding the visual styles of components used in Kanvas, including Kubernetes resources, integrated technologies, and general diagramming elements."
display-title: "false"
list: include
type: extensions
category: kanvas
redirect_from:
  - /extensions/component-shape-guide
  - extensions/component-shape-guide
---

# Component Style Guide

When you're designing and visualizing cloud-native systems in [Kanvas](https://kanvas.new/), you'll interact with a rich library of visual elements. This guide is here to help you understand what those visuals mean.

We'll cover two main types of visual elements you'll encounter:
- **System-Defined Visuals:** These are icons and styles that Kanvas uses to represent specific functionalities, component types (like Kubernetes resources), or integrated technologies. They have **predefined** meanings within the system.
- **General Diagramming Elements:** These are tools like generic shapes, arrows, and flowchart symbols provided for your custom diagramming and annotation needs. Their meaning is primarily defined by you or by common diagramming conventions.

## System-Defined Visual Representations

This section focuses on visual elements that have a specific, predefined meaning within Kanvas. Understanding these will help you accurately interpret designs.

### Kubernetes Components

When you're working with Kubernetes designs, you'll notice that components related to Kubernetes have a distinct and consistent visual style. This deliberate approach is designed to help you instantly recognize and understand the various Kubernetes resources within your infrastructure designs.

<a href="../../../assets/img/shapes/k8s_style.gif" target="_blank">
  <img src="../../../assets/img/shapes/k8s_style.gif" style="width:50%; height:auto;" alt="Example of using kubernetes components in Kanvas">
</a>

Kanvas employs its own thoughtful design system to represent Kubernetes resources, built on a few key principles for clarity:

- **Uniform Color Scheme:** Kubernetes component icons in Kanvas typically use a **distinctive blue background**. This blue often serves as a standard identifier for Kubernetes-related elements, helping you differentiate them at a glance.
- **Standardized Icon Structure:** While the exact outer container shape can vary, the fundamental structure is consistent: an outer container shape with the blue background, encompassing a unique inner white symbol.
- **Outer Container Shape Variations:** You'll observe that the blue background is framed by different outer shapes. While a **blue rounded square** is the most versatile container, specific patterns do emerge:
    - **Triangles:** A prominent pattern is the use of triangular outer shells for core networking resources like `Service` and `API Service`.
    - **Hexagons:** You might observe hexagonal shapes for some foundational workload controllers like `StatefulSet`.
    - **Unique Polygons:** Several Kinds feature highly unique shapes tailored to their function, such as those for `Endpoints`, `PriorityClass` (often resembling a gauge), or `ValidatingWebhookConfiguration` (which might use a shield shape, suggesting security).
- **Meaningful Inner Symbols:** The white symbol inside each icon is the crucial unique identifier for that specific Kubernetes Kind. These symbols are often inspired by the core function of the resource (e.g., a symbol suggesting containerization for a Pod) and may occasionally incorporate elements of the official Kubernetes logo.

> This systematic approach ensures that once you learn the basics of this visual language, you can easily identify any Kubernetes component in your designs.

### Icons for Integrated Technologies and Their Components

Ever looked at a complex design in Kanvas, perhaps one shared by a colleague, and seen various software logos wondering what they represent? Or, when crafting your own designs, have you needed to represent specific cloud services, databases, or other third-party tools using their official, recognizable visuals? Kanvas is designed to make this straightforward.

Beyond its rich support for Kubernetes-native resources, Kanvas allows you to visually represent a wide array of integrated services and technologies. It does this by using their widely recognized, official icons, helping you create clear and immediately understandable designs for your hybrid, multi-cloud, and microservice architectures.

<a href="../../../assets/img/shapes/component_style.gif" target="_blank">
  <img src="../../../assets/img/shapes/component_style.gif" style="width:50%; height:auto;" alt="Example of using integrated components in Kanvas">
</a>

**How Integrated Components are Represented**

To help you navigate its extensive library of integrations, Kanvas organizes these components in a clear hierarchy:

-   **Categories:** At the highest level, integrations are grouped into logical Categories like "Cloud Native Network," "Database," or "Orchestration & Management." This helps you quickly find the type of technology you're looking for.
-   **Integration Models:** Within each category, you'll find specific integration models (e.g., "AWS App Mesh," "Prometheus," "HashiCorp Consul"). Each model is represented by the official logo or a standardized symbol of that technology, leveraging familiar branding for easy identification.
-   **Designable Components:** Each integration model serves as a source for one or more designable components – these are the actual, functional building blocks you'll drag onto your canvas. These individual components have their own icons, which are typically derived from the parent integration's logo to maintain a clear visual link to their origin.

**What These Icons Mean in Your Design**

When you incorporate an icon representing a specific cloud service (like an AWS S3 bucket) or a component from an integrated technology, you are making a clear statement: your architecture includes, interacts with, or plans to use that specific external service or functional unit. This visual mapping is a system-defined representation.

This approach allows you to build rich, unambiguous designs where you and your team can easily see which parts of your system depend on specific cloud resources, networking functions, observability tools, or other managed services.

{% include alert.html type="info" title="Exploring All Integrations" content="This guide covers the visual style of components. For a complete catalog of all technologies that Meshery integrates, visit the integrations directory. <strong><a href='https://docs.meshery.io/extensions/integrations'>Explore All Integrations</a></strong>" %}

## General Diagramming Elements

Now, let's explore the visual elements that give you creative freedom. This part of the guide focuses on elements designed for your flexible diagramming, annotation, and custom visual communication needs.

<a href="../../../assets/img/shapes/shape_style.gif" target="_blank">
  <img src="../../../assets/img/shapes/shape_style.gif" style="width:50%; height:auto;" alt="Overview of general diagramming elements">
</a>

{% include alert.html type="info" title="Key Principle: User-Defined Meaning" content="The key principle for all elements in this section is that the Meshery system itself does not assign them specific operational or semantic functions." %}

### Generic Shapes

The "Shapes" palette in Kanvas offers a diverse collection of common geometric figures and symbolic graphics. These are your go-to tools for general-purpose diagramming, especially when you need to represent concepts not covered by system-defined component icons. For every shape listed, we will explain its general use, but remember, its specific meaning in your design is up to you.

<a href="/../../../assets/img/shapes/shapes.png" target="_blank">
  <img src="../../../assets/img/shapes/shapes.png" style="width:50%; height:auto;" alt="Generic shapes palette in Kanvas">
</a>

<details>
<summary><strong>Foundational Geometric Shapes</strong></summary>
<br>
These are the quintessential building blocks for many diagrams – your circles, squares, triangles, and basic polygons. They offer simple, clean, and universally understood forms for a wide range of uses.

{% include extension-guide.html
  data_file="foundational_geometric_shapes"
  guide_title="Shape"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}

</details>

<details>
<summary><strong>Common Representational & Symbolic Shapes</strong></summary>
<br>
This group includes shapes that, by common convention, often evoke a more specific symbolic meaning, like using a "Barrel" for data storage or an "Actor" for a user role. Leveraging these established visual metaphors can make your custom diagrams more intuitive.

{% include extension-guide.html
  data_file="common_representational_shapes"
  guide_title="Shape"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}

</details>

<details>
<summary><strong>Specialized & Decorative Geometric Shapes</strong></summary>
<br>
When your diagrams require a more distinct visual style or an element for emphasis, this collection offers a variety of options, from complex polygons to stylized forms like "Crescent," "Star," or "XWing." These shapes are less commonly used in Kanvas and have no universally accepted meaning.

{% include extension-guide.html
  data_file="specialized_decorative_shapes"
  guide_title="Shape"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}

</details>

### Arrows

In Kanvas, arrows are fundamental tools for showing direction or creating simple visual annotations. They are static shapes intended for illustration.

<a href="/../../../assets/img/shapes/arrows.png" target="_blank">
  <img src="../../../assets/img/shapes/arrows.png" style="width:50%; height:auto;" alt="Static arrow shapes palette in Kanvas">
</a>

{% include alert.html type="info" title="Looking for Dynamic Connections?" content="The arrows shown here are static visual aids. To represent actual, functional relationships between components (like network traffic or dependencies), you should use the Edge system instead. <strong><a href='https://docs.meshery.io/extensions/edges-shape-guide'>Learn more</a></strong>" %}

### Flowchart Shapes

To help you visually map out processes, Kanvas includes a dedicated palette of standard flowchart shapes. If you've ever created a flowchart, these symbols will be very familiar to you.

<a href="../../../assets/img/shapes/flowchart.png" target="_blank">
  <img src="../../../assets/img/shapes/flowchart.png" style="width:50%; height:auto;" alt="Flowchart shapes palette in Kanvas">
</a>

### Simple Line Icons

Kanvas also provides a comprehensive library of **Simple Line Icons**, representing common objects, actions, and concepts. These icons are intended for user-driven annotations and visual enhancement. You might use a "lightbulb" icon to signify an idea, or a "warning" icon to flag a component for attention.

<a href="../../../assets/img/shapes/simple_line_icons.png" target="_blank">
  <img src="../../../assets/img/shapes/simple_line_icons.png" style="width:50%; height:auto;" alt="Simple line icons palette in Kanvas">
</a>

## Component Visuals in Different Contexts

A single component will be visually represented differently depending on where you encounter it in the Meshery UI. Let's take the Deployment component as an example to see how its appearance adapts to these different views:

1.**Deployment component with its distinctive shape and badge:**

<a href="../../../assets/img/shapes/deployment-shape.png">
    <img src="../../../assets/img/shapes/deployment-shape.png" style="width:30%; height:auto;" alt="Deployment Component Shape in Kanvas">
</a>

2.**Deployment icon as it might appear in a component selection panel:**

<a href="../../../assets/img/shapes/deployment-icon.png">
    <img src="../../../assets/img/shapes/deployment-icon.png" style="width:50%; height:auto;" alt="Deployment icon in a component selection panel">
</a>

3.**Deployment component as seen in a cluster resource overview:**

<a href="../../../assets/img/shapes/deployment-dashboard.png">
    <img src="../../../assets/img/shapes/deployment-dashboard.png" style="width:50%; height:auto;" alt="Deployment component in a cluster resource overview">
</a>