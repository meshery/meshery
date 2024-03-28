---
layout: default
title: Concepts
permalink: concepts
# redirect_from: concepts/
language: en
list: exclude
abstract: Concepts for understanding Meshery's various features and components.
---

Concepts for understanding Meshery's various features and components.

 <!-- CONCEPTS -->

  <div class="section">
    <a href="{{ site.baseurl }}/concepts">
        <div class="btn-primary">Concepts</div>
    </a>
   <!-- <h6><a href="{{ site.baseurl }}/concepts/logical" class="text-black section-title">Conceptual</a></h6>
     <ul>
      {% for item in sorted_pages %}
      {% if item.type=="concepts" and item.list!="exclude" and item.language!="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul> -->
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/concepts/logical" class="text-black">Logical</a>
        </p>
      </summary>
      <ul class="section-title">
        {% assign sorted_concepts = site.pages | where: "type","concepts" %}
        {% for item in sorted_concepts %}
        {% if item.type=="concepts" and item.language=="en" -%}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
          {% if item.abstract != " " %}
         -  {{ item.abstract }}
          {% endif %}
          </li>
          {% endif %}
        {% endfor %}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/concepts/architecture" class="text-black section-title">Architectural</a>
        </p>
      </summary>
      <ul>
        {% assign sorted_components = site.pages | where: "type","components" %}
        {% for item in sorted_components %}
        {% if item.type=="components" and item.language=="en" -%}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
            {% if item.abstract != " " %}
            - {{ item.abstract }}
            {% endif %}
          </li>
        {% endif %}
        {% endfor %}
      </ul>
    </details>
  </div>



<!-- {% assign sorted_pages = site.pages | sort: "name" %}

<h2>Architectural Components</h2>
<ul>
    {% for item in sorted_pages %}
    {% if item.type=="components" and item.language=="en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

<h2>Logical Concepts</h2>
<ul>
    {% for item in sorted_pages %}
    {% if item.type=="concepts" and item.language=="en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul> -->