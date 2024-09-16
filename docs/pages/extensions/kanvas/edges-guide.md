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

In Meshery, an edge is a visual representation of a relationship or connection between two or more components within your infrastructure. Occasionally, edges might be self-referential, for example, a network connection over a loopback interface in a Kubernetes Pod. Relationships can be of various types, defining how the components interact, depend on, or influence each other.

Meshery supports multiple types of relationships, such as network connections, bindings, permissions, and hierarchical associations. See [Relationships]({{site.baseurl}}/concepts/logical/relationships) for a set of examples.

From a design perspective, edge styles provide a clear and intuitive way to understand the complex interdependencies in your cloud-native environment. From an operational perspective, edge styles can change dynamically along with your infrastructure, reflecting the real-time state of your components and their relationships. By visualizing edges, Meshery helps you identify potential bottlenecks, troubleshoot issues, and optimize the performance of your applications.

Meshery Designs carefully consider the visual style used for edges, imbuing meaning behind their weight, color, stroke, and arrowhead styles. The following list represents all edge styles and their current meaning in a general context.



{% include extension-guide.html
 data_file="edges"
 guide_title="Edge"
 guide_description="Description"
 guide_svg="SVG"
 guide_assests_folder="shapes"
%}