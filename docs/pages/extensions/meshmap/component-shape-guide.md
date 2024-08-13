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

  /* Adjust SVG column width on large screens */
  @media (min-width: 1024px) {
    td:first-child, th:first-child {
      width: 10%; /* Adjust this value as needed */
    }
  }

  /* Responsive table for mobile screens */
  @media (max-width: 767px) {
    .table-container {
      overflow-x: auto;
    }

    table {
      min-width: 600px; /* Adjust this value based on your content */
    }

    td img {
      width: 50px;
      height: 50px;
    }
  }
</style>

<!-- <style>
  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    text-align: left;
    padding: 10px;
  }

  /* Responsive table */
  @media (max-width: 600px) {
    table, thead, tbody, th, td, tr {
      display: block;
    }

    thead tr {
      display: none; /* Hide table header */
    }

    tr {
      margin-bottom: 10px;
      border-bottom: 1px solid #ddd;
    }

    td {
      display: flex;
      justify-content: space-between;
      align-items: center;
      text-align: left;
      padding-left: 50%;
      position: relative;
    }

    td:before {
      content: attr(data-label);
      position: absolute;
      left: 10px;
      width: calc(50% - 20px);
      font-weight: bold;
      white-space: nowrap;
    }

    td img {
      width: 50px;
      height: 50px;
    }
  }
</style> -->

<div class="table-container">
<table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr>
      <th>Shape</th>
      <!-- <th>Shape</th> -->
      <th>Description</th>
      <th>Common Usecase</th>
      <!-- <th>Dependency</th> -->
    </tr>
  </thead>
  <tbody>
    {% for shape in site.data.shapes %}
      <tr>
        <td style="text-align: center">
            <img src="{{ site.baseurl }}/assets/shapes/{{shape.SVG}}" width="100%" height="75px" alt="Shape"/>
            <div>{{ shape.Shape }}</div>
        </td>
        <!-- <td style="text-align: center">{{ shape.Shape }}</td> -->
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
        <!-- <td>{{ shape.DependOn | default: "No dependencies listed." }}</td> -->
      </tr>
    {% endfor %}
  </tbody>
</table>
</div>