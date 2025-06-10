---
layout: default
title: Identifying Meshery Components
permalink: guides/configuration-management/identifying-components
language: en
abstract: A guide to help you identify and understand the various component icons, shapes, and visual styles used across the Meshery UI.
list: include
type: guides
category: configuration
redirect_from:
  - extensions/component-shape-guide
---

Ever wondered what the different icons and shapes in Meshery represent? Whether you're looking at a dashboard, a settings page, or a design, you'll encounter a rich library of visual elements. This guide is here to help you understand what they mean.

The [components](https://docs.meshery.io/concepts/logical/components) in Meshery fall into two fundamental categories, distinguished by whether they can be orchestrated (managed) by Meshery during deployment:

- **Semantic Components (Orchestratable):** These represent actual infrastructure resources that Meshery can understand and manage during deployment. Examples include Kubernetes resources (like Pods and Services), databases, and other infrastructure components. Meshery will actively manage their lifecycle during deployment.

- **Non-semantic Components (Annotation):** These are visual elements used for documentation and organization, such as text boxes, arrows, shapes, and comments. Meshery ignores these during deployment as they don't represent actual infrastructure.

{% include alert.html type="info" title="Visual Customization" content="All components, whether semantic or non-semantic, support rich visual customization. For example, you can change the color of a Kubernetes Pod icon, modify its shape, or customize its background - it's all configurable!" %}

## Semantic Components

These components represent real infrastructure that Meshery can manage. They can be either built-in (like Kubernetes components) or custom components that you [create](https://docs.meshery.io/guides/configuration-management/creating-models).

### Kubernetes Components

While Kubernetes components are commonly used, they follow the same principles as all other semantic components. They have a default distinct visual style to help you instantly recognize them:

- **Uniform Color Scheme:** Kubernetes component icons typically use a **distinctive blue background** as a standard identifier.
- **Standardized Icon Structure:** The fundamental structure is consistent: an outer container shape with the blue background, encompassing a unique inner white symbol.
- **Meaningful Inner Symbols:** The white symbol inside each icon is the crucial unique identifier for that specific Kubernetes Kind, often inspired by the core function of the resource.

<a href="../../../assets/img/shapes/k8s_component.gif" target="_blank">
  <img src="../../../assets/img/shapes/k8s_component.gif" style="width:50%; height:auto;" alt="Kubernetes components in Meshery">
</a>

### Integrated Technologies

Meshery supports various technologies (like AWS, Prometheus, Istio, KEDA, etc.) with their official icons. These components have the same orchestratable capabilities as Kubernetes components.

<a href="../../../assets/img/shapes/AWS-models.png" target="_blank">
  <img src="../../../assets/img/shapes/AWS-models.png" style="width:50%; height:auto;" alt="AWS models in Meshery">
</a>

{% include alert.html type="info" title="Exploring All Integrations" content="This guide covers the visual style of components. For a complete catalog of all technologies that Meshery integrates, visit the integrations directory. <strong><a href='https://docs.meshery.io/extensions/integrations'>Explore All Integrations</a></strong>" %}

## Non-semantic Components

These components help you document and organize your designs without affecting the actual infrastructure. They include:

- Text boxes and comments for documentation
- Shapes and containers for visual grouping
- Lines and arrows for showing relationships
- Labels and tags for organization

While these components are ignored during deployment, they support the same visual customization options as semantic components.

{% include alert.html type="info" title="Edge Components" content="To learn more about edge components and their visual styles, visit the <strong><a href='https://docs.meshery.io/extensions/edges-shape-guide'>Edge Components Guide</a></strong>" %}

## Foundational Elements

In addition to the rich library of predefined icons, Meshery also provides a core set of foundational geometric shapes. These are not tied to any specific technology but are intended for representing abstract concepts or for simple annotations within your designs where a specific component icon doesn't apply.

<a href="../../../assets/img/shapes/shapes.png" target="_blank">
  <img src="../../../assets/img/shapes/shapes.png" style="width:50%; height:auto;" alt="Generic shapes palette in Meshery">
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
When your diagrams require a more distinct visual style or an element for emphasis, this collection offers a variety of options, from complex polygons to stylized forms like "Crescent," "Star," or "XWing." These shapes are less commonly used in Meshery and have no universally accepted meaning.

{% include extension-guide.html
  data_file="specialized_decorative_shapes"
  guide_title="Shape"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}

</details>

## Component Visuals in Different Contexts

To see how this works in practice, let's take the Kubernetes `Deployment` component as an example. Its appearance adapts to different views in the Meshery UI:

1.The full component shape as it appears in a design:

<a href="../../../assets/img/shapes/deployment-shape.png">
    <img src="../../../assets/img/shapes/deployment-shape.png" style="width:30%; height:auto;" alt="Deployment Component Shape">
</a>

2.The simplified icon as it appears in a component selection panel:

<a href="../../../assets/img/shapes/deployment-icon.png">
    <img src="../../../assets/img/shapes/deployment-icon.png" style="width:50%; height:auto;" alt="Deployment icon in a component selection panel">
</a>

3.The icon as seen in a cluster resource overview:

<a href="../../../assets/img/shapes/deployment-dashboard.png">
    <img src="../../../assets/img/shapes/deployment-dashboard.png" style="width:50%; height:auto;" alt="Deployment component in a cluster resource overview">
</a>