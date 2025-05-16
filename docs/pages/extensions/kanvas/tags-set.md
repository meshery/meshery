---
layout: default
title: Working with Tags
permalink: extensions/working-with-tags
language: en
abstract: Tags can be used to visually group components.
display-title: "false"
list: include
type: extensions
category: kanvas
---

# Working with Tags

Kubernetes components can be assigned Label and Annotation key/value pairs, known as Tags. When pairs of Labels or Annotations match, a relationship is established and visualized as shown below.

<a href="{{ site.baseurl }}/assets/img/kanvas/tags.gif"><img style="border-radius: 0.5%;" alt="Working-with-Tags" style="width:800px;height:auto;" src="{{ site.baseurl }}/assets/img/kanvas/tags.gif" /></a>

You can group components using tags. Tags are key-value pairs that help you organize and categorize components within your design. Tags can be used to visually group components. You can also use tags to filter components and view only those that match the tag criteria.

## Labels and Annotations

Designs support two different types of tags: Labels and Annotations. Labels are often used to identify components and are visible on the design canvas. Annotations are often used to provide additional information about components.

<div class="alert alert-warning" role="alert">
<h4 class="alert-heading">Performance Consideration</h4>
Tags are indexed and searchable. However, the performance of design operations may degrade as the number of tags increases. To ensure an optimal user experience, we recommend using tags judiciously and limiting the number of tags used in a design.

Upon loading a design exceeds that exceeds 10 tags within a single design, Kanvas will automatically disable grouping by tags. You can manually enable grouping by tags by clicking the “Group Components” button in the Designer dock.
</div>
