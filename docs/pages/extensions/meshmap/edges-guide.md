---
layout: default
title: Edges Shape Guide
permalink: extensions/edges-shape-guide
language: en
abstract: Kubernetes architecture deployment and architecture diagramming tool for cloud native applications - MeshMap.
display-title: "false"
list: include
type: extensions
category: meshmap
---

## Edges Shape Guide

Inside MeshMap, the allocation of specific edge-shapes to signify various purposes creates a coherent and intelligible visual representation of intricate designs.

Below are all the edges with their current usage in a general context.


<style>

  .edges-container {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    justify-content: space-between;
    margin-top: 2rem;

  }
  .edges-card {
	 display: flex;
   flex-direction: column;
	 max-width: 20rem;
	 max-height: 30rem;
	 gap: 1rem;
}
  .edges-svg-container {
	 height: auto;
	 display: flex;
	 flex-direction: column;
	 align-items: center;
	 gap: 0.455rem;
   flex-basis: 30%;
}
  .edges-svg-container img {
   width: 50%;
   height: auto;
}
  .edges-details {
	 display: flex;
	 flex-direction: column;
	 flex-basis: 77%;
	 gap: 10px;
   text-align: center;
}
  @media (max-width: 767px) {
    .esges-container {
      flex-direction: column;
    }
    .edges-svg-container {
      gap: 0.3rem;
    }
    .edges-svg-container img {
      width: 40%;
    }
    .edges-card{
      max-width: 30rem;
      flex-direction: column;
    }
   }
</style>


<div class="edges-container">
{% assign sorted_edges = site.data.edges | sort: 'Edge' %}
    {% for edge in sorted_edges %}

      <div class="edges-card">
      <div class="edges-svg-container">
      {% if edge.SVG %}
        <img src="{{ site.baseurl }}/assets/shapes/edgesAndArrowHeads/{{edge.SVG}}" alt="Edge">
      {% endif %}
       <div style="text-align:center;">{{ edge.Edge }}</div>
      </div>
      <div class="edges-details">
      <div>{{ edge.Description }}</div>
      </div>
      </div>

    {% endfor %}
 </div>