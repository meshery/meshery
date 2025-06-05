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
  - extensions/component-shape-guide
---

## Component Style Guide

When you're designing and visualizing cloud-native systems in Kanvas, you'll interact with a rich library of visual elements. This guide is here to help you understand what those visuals mean.

We'll cover two main types of visual elements you'll encounter:
-   **System-Defined Visuals:** These are icons and styles that Kanvas uses to represent specific functionalities, component types (like Kubernetes resources), or integrated technologies. They have predefined meanings within the system.
-   **General Diagramming Elements:** These are tools like generic shapes, arrows, and flowchart symbols provided for your custom diagramming and annotation needs. Their meaning is primarily defined by you or by common diagramming conventions.

### System-Defined Visual Representations

This section focuses on visual elements that have a specific, predefined meaning within Kanvas. Understanding these will help you accurately interpret designs and leverage Meshery's capabilities.

#### Kubernetes Components

When you're working with Kubernetes designs in Kanvas, you'll notice that components related to Kubernetes have a distinct and consistent visual style. This deliberate approach is designed to help you instantly recognize and understand the various Kubernetes resources within your infrastructure designs.

Kanvas employs its own thoughtful design system to represent Kubernetes resources, built on a few key principles for clarity:

-   **Uniform Color Scheme:** Kubernetes component icons in Kanvas typically use a **distinctive blue background**. This blue often serves as a standard identifier for Kubernetes-related elements, helping you differentiate them at a glance.
-   **Standardized Icon Structure:** While the exact outer container shape can vary, the fundamental structure is consistent: an outer container shape with the blue background, encompassing a unique inner white symbol.
-   **Outer Container Shape Variations:** You'll observe that the blue background is framed by different outer shapes. While a **blue rounded square** is the most versatile container, specific patterns do emerge:
    -   **Triangles:** A prominent pattern is the use of triangular outer shells for core networking resources like `Service` and `API Service`.
    -   **Hexagons:** You might observe hexagonal shapes for some foundational workload controllers like `StatefulSet`.
    -   **Unique Polygons:** Several Kinds feature highly unique shapes tailored to their function, such as those for `Endpoints`, `PriorityClass` (often resembling a gauge), or `ValidatingWebhookConfiguration` (which might use a shield shape, suggesting security).
-   **Meaningful Inner Symbols:** The white symbol inside each icon is the crucial unique identifier for that specific Kubernetes Kind. These symbols are often inspired by the core function of the resource (e.g., a symbol suggesting containerization for a Pod) and may occasionally incorporate elements of the official Kubernetes logo.

> This systematic approach ensures that once you learn the basics of this visual language, you can easily identify any Kubernetes component in your Kanvas designs. While these aren't official icons mandated by the Kubernetes project, they are Meshery's carefully considered way of providing a clear and consistent visual representation.

{% raw %}{% include alert.html type="light" title="Note on a Specific Component: Deployment Shape" content="Due to characteristics of the rendering library used in Kanvas, a Kubernetes **Deployment** component might sometimes display with a rounded-rectangle outline as its primary shape. You might also see a small badge placed over its lefthand-side border; this badge aids in quickly identifying the component type in certain views." %}{% endraw %}

#### Icons for Integrated Technologies and Their Components

Ever looked at a complex design and seen various software logos, wondering what they represent? Or needed to use official icons for cloud services or other tools in your own designs? Kanvas makes this straightforward.

Beyond Kubernetes-native resources, your designs frequently bring together a wide array of services and technologies. Kanvas visually represents these integrations—and the specific **components** they offer—using their widely recognized, official icons.

**Visual Organization in the Palette**

In Kanvas, you'll find these integrations thoughtfully organized:
-   They are grouped into logical **Categories** (such as "Cloud Native Network," "Database," etc.).
-   Within each category, Meshery lists specific **integration models**. Each model (e.g., "AWS App Mesh," "Prometheus") is typically represented by the **official logo** of that technology.

**Designable Components from Integration Models**

Each integration model serves as a source for one or more **designable components** – the actual building blocks you'll drag onto your canvas. When an integration model in the palette is clicked, it reveals the specific components it offers. These individual components also feature their own icons, typically derived from the parent integration's logo.

> This visual mapping is a **system-defined representation**: Kanvas uses that specific icon to denote that exact service, technology, or the functionality it provides. When you incorporate one into your design, it precisely signifies that your architecture includes or interacts with that external service or functional unit.

This structured visual approach allows you to build rich, unambiguous diagrams where you can easily see which parts of your system depend on specific cloud resources, networking functions, or other managed services.

*(An animated GIF here, showing an integration model's components being revealed, one dragged onto the canvas, and its settings being opened, would be very illustrative.)*

{% raw %}{% include alert.html type="info" title="Official Iconography and Further Details" content='Meshery leverages these official or standardized icons to clearly identify the technology or service being integrated. If you require detailed guidelines on the usage of the original logos themselves, it\'s always best to refer to the official branding or iconography resources provided by the respective vendors. For specifics on what designable components each Meshery integration model provides, exploring them within the Meshery UI palette is your most direct reference.' %}{% endraw %}

### General Diagramming Elements

Now, let's explore the visual elements that are all about giving you creative freedom. This part of the guide focuses on elements designed for your flexible diagramming, annotation, and custom visual communication needs.

{% raw %}{% include alert.html type="light" title="A Note on General Elements" content="A key thing to remember is that, generally, the **Meshery system itself does not assign specific operational or semantic functions to the elements in this section**. Their power lies in how *you* choose to use them to enrich your designs, illustrate workflows, or create custom visual representations tailored to your team's conventions." %}{% endraw %}

#### Generic Shapes

The "Shapes" palette in Kanvas offers a diverse collection of common geometric figures and symbolic graphics. These are your go-to tools for general-purpose diagramming, especially when you need to represent concepts not covered by system-defined component icons. For every shape listed, we explain its general use, but remember, its specific meaning in your design is up to you.

##### Foundational Geometric Shapes
These are the quintessential building blocks for many diagrams – your circles, squares, triangles, and basic polygons. They offer simple, clean, and universally understood forms for a wide range of uses.

{% include extension-guide.html
  data_file="foundational_geometric_shapes"
  guide_title="Shape"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}

##### Common Representational & Symbolic Shapes
This group includes shapes that, by common convention, often evoke a more specific symbolic meaning, like using a "Barrel" for data storage or an "Actor" for a user role. Leveraging these established visual metaphors can make your custom diagrams more intuitive.

{% include extension-guide.html
  data_file="common_representational_shapes"
  guide_title="Shape"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}

##### Specialized & Decorative Geometric Shapes
When your diagrams require a more distinct visual style or an element for emphasis, this collection offers a variety of options, from complex polygons to stylized forms like "Crescent," "Star," or "XWing."

{% include extension-guide.html
  data_file="specialized_decorative_shapes"
  guide_title="Shape"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}

#### Arrows and Connectors

Arrows and connector lines are fundamental tools for linking elements and showing direction. Kanvas provides both a palette of **static arrow shapes** for quick visual annotations and a more sophisticated **Edge system** for creating meaningful, styled connections between components.

This section focuses on the static arrow shapes available in the "Arrows" palette. Think of these primarily as visual aids for annotation. They function similarly to the generic shapes—their meaning is defined by you in your design.

**For Meaningful, Dynamic Connections: The Edge System**

When you need to represent **actual, functional relationships** between your components—such as network connections, dependencies, or data flows—Kanvas provides a dedicated and more powerful "Edge" system. These Edges are not just static lines; they are an integral part of your design's model and can be styled with specific colors, thicknesses, and arrowheads to convey rich semantic information.

{% raw %}{% include alert.html type="info" title="All About Edges: Styling, Meaning, and Usage" content='For comprehensive details on how to create and interpret these meaningful relationships using the Edge system, please consult our dedicated **Edge Style Guide**. That guide is the definitive resource for understanding how Kanvas visually represents and interprets the connections between your components.' %}{% endraw %}

#### Flowchart Shapes

To help you visually map out processes, Kanvas includes a dedicated palette of standard flowchart shapes. If you've ever created a flowchart, these symbols will be very familiar to you.

> **Meshery System Meaning:** Crucially, beyond their universally understood roles in flowcharting (e.g., a diamond for a decision), these shapes have **no predefined operational meaning within the Meshery system itself.** Kanvas provides them as a powerful diagramming aid. The interpretation of any flowchart you create relies on these standard conventions and the specific context you provide.

#### Simple Line Icons

Kanvas also provides a comprehensive library of **Simple Line Icons**, representing common objects, actions, and concepts.

**Purpose: Enhancing Your Designs with Visual Cues**

These icons are intended for **user-driven annotations and visual enhancement**. You might use a "lightbulb" icon to signify an idea, or a "warning" icon to flag a component for attention.

> **Meshery System Meaning: User-Defined.** When you use these icons from the general palette, they do **not** have a predefined operational meaning in Meshery. Their interpretation is flexible and is defined by **you**, the user, within the context of your specific design. Their strength lies in their generic nature and broad applicability for custom annotations.

For example, you might place a "lightbulb" icon next to a note to signify an idea, or a "warning" icon to visually flag a component that requires special attention based on your team's conventions. Their strength lies in their generic nature and broad applicability for custom annotations.


#### Where 
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