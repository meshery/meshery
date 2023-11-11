---
layout: default
title: Installation
type: installation
abstract: Installation procedures for deploying Meshery with mesheryctl.
permalink: installation
redirect_from: 
- platforms
- platforms/
- installation/
- installation/platforms
- installation/platforms/
language: en
list: exclude
---

## Supported Platforms

Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster. See the complete list of supported platforms in the table below. With service meshes having sprung to life in the context of Kubernetes, so too, can Meshery’s deployment models be characterized in the context of Kubernetes. A given deployment of Meshery can be described as either an _in-cluster_ or an _out-of-cluster_ deployment. Meshery deploys as a stand-alone, management plane on a Docker host (_out-of-cluster_) or as a management plane in a Kubernetes cluster (_in-cluster_).

{% assign sorted_index = site.pages | sort: "name" | alphabetical %}

### Install `mesheryctl`

<ul>
    {% for item in sorted_index %}
    {% if item.type=="installation" and item.category=="mesheryctl" and item.list=="include" and item.language == "en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

### Install on Kubernetes

<ul>
    {% for item in sorted_index %}
    {% if item.type=="installation" and item.category=="kubernetes" and item.list=="include" and item.language == "en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

### Install on Docker

<ul>
    {% for item in sorted_index %}
    {% if item.type=="installation" and item.category=="docker" and item.list=="include" and item.language == "en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>


<!-- {% include toc.html page=reference %} -->
