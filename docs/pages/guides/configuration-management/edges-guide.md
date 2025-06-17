---
layout: default
title: Edge Styles Guide
permalink: guides/configuration-management/edges-guide
language: en
abstract: Learn about the different edge styles and how to interpret their visual meaning in Meshery's component relationships.
list: include
type: guides
category: configuration
redirect_from:
  - extensions/edges-shape-guide
  - extensions/edges-shape-guide/
---

In Meshery, the line that connects [components](https://docs.meshery.io/concepts/logical/components) is called an **edge**. Each edge visually represents a [relationship](https://docs.meshery.io/concepts/logical/relationships) and uses a specific style to communicate its nature.

This guide helps you interpret the most common edge styles you will encounter.

### Interpreting Common Edge Styles

Meshery uses a set of default visual styles to provide at-a-glance information about the type of connection an edge represents. While these styles can be customized in the UI, understanding the defaults is key to interpreting component relationships.

#### Line Style: The Primary Indicator

The line style is the most important visual cue for understanding an edge's purpose.

- **Dotted Line**:
  - **What it means**: A **semantic relationship**. This represents a real, functional connection that Meshery understands and can manage, such as a network link or a volume mount.
  - **When you'll see it**: These lines indicate active relationships between components, showing how they interact and communicate. The dotted pattern represents dynamic connections, while arrowheads show the direction of data flow or dependency.

- **Solid Line**:
  - **What it means**: A **non-semantic annotation**. This is a visual note or organizational aid for human interpretation only. Meshery's engine ignores these connections.
  - **When you'll see it**: These lines represent static or conceptual relationships between components. They help visualize structural connections or highlight specific component groupings without implying active data flow.

#### Color: A Secondary Cue

Color provides an additional hint about an edge's nature.

- **Green / Teal**:
  - This is the default color for **non-semantic annotations**, helping them stand out from functional connections.

- **Blue / Grey**:
  - These are the typical default colors for **semantic relationships**.

### Edge Style Gallery

The following gallery showcases the full range of visual styles available for edges in Meshery.

{% include extension-guide.html
  data_file="edges"
  guide_title="Edge"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}