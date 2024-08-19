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

  .shapes-container {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    justify-content: space-between;
    margin-top: 2rem;

  }
  .shapes-card {
	 display: flex;
   flex-direction: column;
	 max-width: 20rem;
	 max-height: 30rem;
	 gap: 1rem;
}
  .shapes-svg-container {
	 height: auto;
	 display: flex;
	 flex-direction: column;
	 align-items: center;
	 gap: 0.455rem;
   flex-basis: 30%;
}
  .shapes-svg-container img {
   width: 50%;
   height: auto;
}
  .shapes-details {
	 display: flex;
	 flex-direction: column;
	 flex-basis: 77%;
	 gap: 10px;
   text-align: center;
}
  @media (max-width: 767px) {
    .shapes-container {
      flex-direction: column;
    }
    .shapes-svg-container {
      gap: 0.3rem;
    }
    .shapes-svg-container img {
      width: 40%;
    }
    .shapes-card{
      max-width: 30rem;
      flex-direction: column;
    }
   }
</style>


<div class="shapes-container">
{% assign sorted_shapes = site.data.shapes | sort: 'Shape' %}
    {% for shape in sorted_shapes %}

      <div class="shapes-card">
      <div class="shapes-svg-container">
      {% if shape.SVG %}
        <img src="{{ site.baseurl }}/assets/shapes/{{shape.SVG}}" alt="Shape">
      {% endif %}
       <div style="text-align:center;">{{ shape.Shape }}</div>
      </div>
      <div class="shapes-details">
      <div>{{ shape.Description }}</div>
      <div>{{ shape.CommonUseCase }}</div>
      </div>
      </div>

    {% endfor %}
 </div>
