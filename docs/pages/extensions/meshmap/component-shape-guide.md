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


<style>
  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    text-align: center;
    padding: 10px;
    border: 1px solid #ddd;
  }

  @media (min-width: 1024px) {
    td:first-child, th:first-child {
      width: 10%; 
    }
  }

  @media (max-width: 767px) {
    .table-container {
      overflow-x: auto;
    }

    table {
      min-width: 600px; 
    }

    td img {
      width: 50px;
      height: 50px;
    }
  }
</style>

<div class="table-container">
<table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr>
      <th>Shape</th>
      <th>Description</th>
      <th>Common Usecase</th>
    </tr>
  </thead>
  <tbody>
    {% for shape in site.data.shapes %}
      <tr>
        <td style="text-align: center">
            <img src="{{ site.baseurl }}/assets/shapes/{{shape.SVG}}" width="100%" height="75px" alt="Shape"/>
            <div>{{ shape.Shape }}</div>
        </td>
        <td>{{ shape.Description }}</td>
        <td style="text-align: left">
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
      </tr>
    {% endfor %}
  </tbody>
</table>
</div>