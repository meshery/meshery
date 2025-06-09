---
layout: default
title: Component Style Guide
permalink: extensions/component-style-guide
language: en
abstract: Kubernetes architecture deployment and architecture diagramming tool for cloud native applications - Kanvas.
display-title: "false"
list: include
type: extensions
category: kanvas
---

## Component Style Guide

The **Component Style Guide** outlines Meshery’s visual conventions for representing different types of components in Kanvas. It aims to promote consistency, clarity, and intuitive understanding across all designs.

> 🎯 This guide replaces the "Component Shape Guide" and provides richer context and visual distinctions between actively used, generic, and deprecated shapes.

---

### 🔷 Shape Usage & Meaning

When adding components, users should select a shape that reflects its function and significance. This guide categorizes shapes into:

- ✅ **Actively Used** – Shapes currently employed across Meshery visualizations.
- ⚙️ **Generic** – Provided by the visualization library, available for custom or experimental use.
- ❌ **Deprecated** – No longer recommended; retained for backward compatibility.

By default, the **circle** shape is used for new components. However, we encourage using custom shapes when needed to express specific roles.

{% include extension-guide.html 
 data_file="shapes"
 guide_title="Shape"
 guide_description="Description"
 guide_usecase="CommonUseCase"
 guide_svg="SVG"
 guide_assets_folder="shapes"
%}

---

{% include alert.html type="info" title="Note" content="Although some shapes are visually distinct, their functionality may be similar. When in doubt, consult the Edge Style Guide or Interpreting Meshery Designs for additional context." %}

---

### 🧭 Deployment Component Example

Some components, like Kubernetes **Deployment**, use a specific shape for a technical reason.

{% include alert.html type="light" title="Example" content="A Kubernetes Deployment appears as a rounded rectangle. This is a limitation of the Javascript library used for rendering, but a badge helps distinguish it from other shapes." %}

<br/>

<a href="../../../assets/img/deployment-shape.png">
    <img src="../../../assets/img/deployment-shape.png" style="width:50%; height:auto;" alt="Deployment Component Shape">
</a>
<p><em>Deployment component with its distinctive rounded-rectangle shape and badge</em></p>

<br/>

[![Deployment in Icon Set](../../../assets/img/deployment-icon.png)](../../../assets/img/deployment-icon.png)
<em>Deployment icon in the component selection panel</em>

<br/>

[![Deployment in Dashboard](../../../assets/img/deployment-dashboard.png)](../../../assets/img/deployment-dashboard.png)
<em>Deployment component as seen in the cluster resource overview</em>

---

### 🔗 Related Guides

- 📘 [Edge Style Guide](../edges-guide)
- 🧩 [Interpreting Meshery Designs](../interpreting-meshery-designs)

---

💡 Want to suggest new shapes or improvements? Contribute on [GitHub](https://github.com/meshery/meshery).
