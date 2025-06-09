---
layout: default
title: Edges Style Guide
permalink: extensions/edges-shape-guide
language: en
abstract: Learn how to interpret the different lines (edges) in your Meshery designs, from their style and color to their semantic meaning.
list: include
type: extensions
category: kanvas
---

When you connect components in [Kanvas](https://kanvas.new/), the line you draw is called an **edge**. Each edge visually represents a [relationship]({{site.baseurl}}/concepts/logical/relationships) between your components. Just as each [component has its own shape]({{site.baseurl}}/guides/configuration-management/identifying-components) to signify its type, each edge has a style that tells you about the nature of that connection.

{% include alert.html type="info" title="The Two Types of Edge Styles" content="Meshery treats edges in two fundamentally different ways:
<br><strong>Semantic Edges:</strong> These represent actual infrastructure relationships that Meshery manages and deploys. While you can style them freely, their true power lies in their ability to define and manage real connections between your components.
<br><strong>Annotation Edges:</strong> These are purely for documentation and visual organization. They help you document and explain your designs without creating any actual infrastructure connections." %}

## Understanding Edge Types

When you connect components in [Kanvas](https://kanvas.new/), you can create two types of connections:

### Semantic Edges: Managed Infrastructure Connections

Semantic edges represent real relationships in your infrastructure that Meshery will actually manage and deploy. These include critical connections like:
- Network links between services
- Volume mounts between containers
- Permission relationships between components

While you have complete freedom to style these edges (just like annotation edges), their true importance lies in their functional purpose:

#### Common Semantic Edge Types
| Edge Type | Purpose | Configuration Options |
| :--- | :--- | :--- |
| `Edge-Network` | Defines network connectivity between services | Protocol, Port, Security settings |
| `Edge-Mount` | Establishes volume mounts between components | Mount path, Read/Write permissions |
| `Edge-Permission` | Sets up access control relationships | Access levels, Scope, Conditions |

### Annotation Edges: Your Documentation Tools

Annotation edges are your personal documentation tools. They help you:
- Explain complex relationships
- Mark important paths
- Create visual guides

The key difference is that these edges exist only in your design - they have no impact on your actual infrastructure.

## Styling Your Edges

Both semantic and annotation edges share the same rich set of styling options. Select any edge to reveal the styling toolbar

![]()

### 1. Adding Meaning with Labels
The **Rename Edge** tool (pencil icon) lets you add descriptive text directly on the edge.
- **Use cases**:
  - Document traffic flows ("User Traffic → API Gateway")
  - Mark process steps ("Requires Manual Review")
  - Note important details ("Backup Route")

### 2. Choosing the Right Line Style
The line style controls how your edge appears on the canvas.

**Basic Styles** (first style menu):
- `Solid Line` - Clear, direct connections
- `Dashed Line` - For optional or secondary paths
- `Dotted Line` - For less prominent connections

**Advanced Shapes** (wavy line icon):
- `Zigzag Line` - Perfect for showing indirect or complex paths
- `Wave Line` - Great for flexible or non-critical connections
- `Bezier Curve Line` - Creates smooth, aesthetic layouts that avoid overlapping
- `Line with Circles` - Highlights specific points or stages along a path

### 3. Defining Connection Points
The **Arrow Style** palette lets you customize how your edge starts and ends:

| Style | Best For |
| :--- | :--- |
| `Arrowhead` (Default) | Clear direction indication |
| `Filled Diamond Head` | Decision points or conditional flows |
| `Filled Square Head` | Termination points or specific interfaces |
| `Filled Triangle Head` | Alternative directional indicators |

### 4. Using Color Effectively
The **Colors** tool (palette icon) lets you create your own visual language:
- **Suggested color coding**:
  - Red: Requires attention or review
  - Green: Approved or active paths
  - Yellow: Connections with known issues
  - Blue: Standard documentation
  - Purple: Special or custom processes

### 5. Drawing Attention
The **Marching-ants** effect (three dots icon) adds an animated stroke to your edge.
- **Perfect for**:
  - Presentations and demos
  - Team reviews
  - Highlighting critical paths
  - Screen sharing sessions

## Edge Style Reference

This gallery shows all available edge styles in Meshery. Each style can be applied to both semantic and annotation edges, helping you create clear and meaningful diagrams.

{% include extension-guide.html
  data_file="edges"
  guide_title="Edge"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
%}