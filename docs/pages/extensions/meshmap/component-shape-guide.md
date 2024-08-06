---
layout: default
title: Component Shape Guide
permalink: extensions/component-shape-guide
language: en
abstract: Kubernetes architecture deployment and architecture diagramming tool for cloud native applications - MeshMap.
display-title: "false"
list: include
type: extensions
category: meshmap
---

## Component Shape Guide

Inside MeshMap, the allocation of specific shapes to signify various purposes creates a coherent and intelligible visual representation of intricate designs.
Currently, the circle is used as the default shape for new components. However, if users or contributors have alternative shapes they believe better suit a particular component, they are encouraged to propose them.

Although the usage of the components is divided into categories, some shapes serve as a universal representation of particular components.

Below are all the shapes with their current usage in a general context.

<table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr>
      <th>SVG Image</th>
      <th>Shape</th>
      <th>Description</th>
      <th>Common Usecase</th>
      <th>Dependency</th>
    </tr>
  </thead>
  <tbody>
    {% for shape in site.data.shapes %}
      <tr>
        <td style="text-align: center; padding: 10px;">
          <svg width="75" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            {{ shape.SVG | remove: '<svg' | remove: '</svg>' }}
          </svg>
        </td>
        <td>{{ shape.Shape }}</td>
        <td>{{ shape.Description }}</td>
        <td>
          <ul>
            {% if shape.CommonUsecase %}
              {% for use in shape.CommonUsecase %}
                <li>{{ use }}</li>
              {% endfor %}
            {% else %}
              <li>No common use cases listed.</li>
            {% endif %}
          </ul>
        </td>
        <td>{{ shape.DependOn | default: "No dependencies listed." }}</td>
      </tr>
    {% endfor %}
  </tbody>
</table>