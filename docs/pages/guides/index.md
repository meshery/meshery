---
layout: default
title: Guides
permalink: guides
redirect_from: 
- guides/
- tasks
language: en
list: exclude
abstract: Using, operating, and troubleshooting Meshery.
---

Guides to using, operating, and troubleshooting Meshery's various features and components.

<!-- {% assign sorted_guides = site.pages | sort: "type" | reverse %} -->

<!-- ### General 

<ul>
    {% for item in sorted_guides %}
    {% if item.type=="guides" and item.category!="mesheryctl" and item.category!="tutorials" and item.list!="exclude" and item.language=="en"  -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

### <a href="{{ site.baseurl }}/guides/tutorials" class="text-black">Tutorials</a>

<ul>
  {% for item in sorted_guides %}
  {% if item.type=="guides" and item.category=="tutorials" and item.list!="exclude" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>

### <a href="{{ site.baseurl }}/guides/mesheryctl" class="text-black">Meshery CLI</a>

<ul>
  {% for item in sorted_guides %}
  {% if item.type=="guides" and item.category=="mesheryctl" and item.list!="exclude" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>

<!-- {% comment %}
#
#  Change date order by adding '| reversed'
#  To sort by title or other variables use {% assign sorted_posts = category[1] | sort: 'title' %}
#
{% endcomment %}

{% for guide in site.adapter %}
<h2 id="{{guide[0] | uri_escape | downcase }}">{{guide[0] | capitalize}}1</h2>

{% endfor %}

{% assign sorted_guides = site.guides | sort %}
{% for guide in sorted_guides %}
<h2 id="{{guide[0] | uri_escape | downcase }}">{{guide[0] | capitalize}}</h2>

{% endfor %} -->
<!-- GUIDES -->
  <div class="section">
    <a href="{{ site.baseurl }}/guides">
        <div class="btn-primary">Guides</div>
    </a>
    <!-- <h6><a href="{{ site.baseurl }}/guides" class="text-black section-title">Guides</a></h6> -->
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/guides/mesheryctl/" class="text-black">Using the CLI</a>
        </p>
      </summary>
      <ul class="section-title">
        {% assign sorted_mesheryctl = site.pages | where: "type","guides" %}
        {% for item in sorted_mesheryctl %}
        {% if item.type=="guides" and item.category=="mesheryctl" and item.language=="en" -%}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
          {% if item.abstract != " " %}
            - {{ item.abstract }}
          {% endif %}
          </li>
          {% endif %}
        {% endfor %}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/guides/infrastructure-management" class="text-black">Infrastructure Management</a>
        </p>
      </summary>
      <ul class="section-title">
       {% assign sorted_infrastructure = site.pages | where: "type","guides" %}
          {% for item in sorted_infrastructure %}
          {% if item.type=="guides" and item.category=="infrastructure" and item.language=="en" -%}
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
          <a href="{{ site.baseurl }}/guides/performance-management" class="text-black">Performance Management</a>
        </p>
      </summary>
      <ul class="section-title">
        {% assign performance = site.pages | where: "type","guides" %}
          {% for item in performance %}
          {% if item.type=="guides" and item.category=="performance" and item.language=="en" -%}
            <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
            {% if item.abstract != " " %}
              - {{ item.abstract }}
            {% endif %}
            </li>
            {% endif %}
          {% endfor %}
      </ul>
    </details>
      <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/guides/infrastructure-management" class="text-black">Configuration Management</a>
        </p>
      </summary>
      <ul class="section-title">
        {% assign configuration = site.pages | where: "type","guides" %}
          {% for item in configuration %}
          {% if item.type=="guides" and item.category=="configuration" and item.language=="en" -%}
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
          <a href="{{ site.baseurl }}/guides/infrastructure-management" class="text-black">Troubleshooting</a>
        </p>
      </summary>
      <ul class="section-title">
          {% assign troubleshooting = site.pages | where: "category","troubleshooting" %}
          {% for item in troubleshooting %}
          {% if item.type=="guides" and item.category=="troubleshooting" and item.language=="en" -%}
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
          <a href="{{ site.baseurl }}/guides/infrastructure-management" class="text-black">🧑‍🔬 Tutorials</a>
        </p>
      </summary>
      <ul class="section-title">
          {% assign tutorials = site.pages | where: "category","tutorials" %}
          {% for item in tutorials %}
          {% if item.type=="guides" and item.category=="tutorials" and item.language=="en" -%}
            <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
            {% if item.abstract != " " %}
              -  {{ item.abstract }}
            {% endif %}
            </li>
            {% endif %}
          {% endfor %}
      </ul>
    </details>
    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="guides" and item.category!="mesheryctl" and item.category!="infrastructure" and item.category!="troubleshooting" and item.category!="performance" and item.category!="configuration" and item.category!="tutorials" and item.language=="en" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>

  </div>