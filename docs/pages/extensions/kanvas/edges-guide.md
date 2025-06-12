---
layout: default
title: Edges Style Guide
permalink: extensions/edges-shape-guide
language: en
abstract: Visualize and manage complex cloud-native relationships with Kanvas' dynamic edge styling in Kubernetes architecture.
display-title: "false"
list: include
type: extensions
category: kanvas
---

## Edges Style Guide

In Meshery, an edge is a visual representation of a relationship or connection between two or more components within your infrastructure. Occasionally, edges might be self-referential — for example, a network connection over a loopback interface in a Kubernetes Pod. Relationships can be of various types, defining how components interact, depend on, or influence each other.

Meshery supports multiple types of relationships, such as:

- **Network connections** (e.g., traffic flows)
- **Bindings** (e.g., service-to-deployment mapping)
- **Permissions** (e.g., RBAC relations)
- **Hierarchical associations** (e.g., controller/child relationships)

🔗 For examples of these types in practice, see [Relationships]({{site.baseurl}}/concepts/logical/relationships).

---

## Purpose of Edge Styles

Meshery uses different edge styles (line types, colors, and arrows) to help users quickly identify:

- The **type of relationship**
- Its **state or direction**
- Whether it’s **active, passive, or deprecated**

Edge styles are carefully chosen to reflect meaning:

| Visual Property | Purpose |
|------------------|---------|
| **Color**        | Indicates traffic type or status (e.g., green = active, red = failed, grey = passive) |
| **Thickness**    | Thicker lines may suggest higher bandwidth or importance |
| **Arrowhead**    | Shows direction of flow or control |
| **Line Style**      | Dashed lines indicate indirect, inferred, or optional connections |

---

## Edge Styles and Their Meanings

{% include extension-guide.html
 data_file="edges"
 guide_title="Edge"
 guide_description="Description"
 guide_svg="SVG"
 guide_assests_folder="shapes"
%}

Each edge documented above includes:

- ✅ **Current usage status**
- 📘 **What the style means**
- 👁️ **Visual example**
- 🔄 **Whether it’s dynamic**

If any edge style listed is no longer in use in the latest Meshery UI, it will be marked as **deprecated** or removed in future releases.

---

## Interlinking Resources

To fully understand how edges interact with other Meshery design components:

- 🔍 See [Component Style Guide]({{site.baseurl}}/extensions/component-style-guide)
- 🧠 Read [Interpreting Meshery Designs]({{site.baseurl}}/concepts/interpreting-designs)

---

