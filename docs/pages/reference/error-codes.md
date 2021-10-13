---
layout: default
title: Error Code Reference
abstract: "Meshery Error Code Reference"
permalink: reference/error-codes
type: Reference
---
<style>

.title::first-letter {
  text-transform:capitalize;
}
</style>
## Error Codes and Troubleshooting

Meshery and it's components use a common framework (defined within MeshKit) to generate and document an error with a unique identifier - an error code. Each error code identifies the source component for the error and a standard set of information to describe the error and provide helpful details for troubleshooting the situation surrounding the specific error.

## Error Code Categories by Component

<table style="margin:auto;padding-right:25%; padding-left:20%;">
<thead>
  <tr>
    <th align="right">Component Type</th>
    <th>Component Name</th>
  </tr>
</thead>
<tbody>
  {% for files in site.data.errorref %}
    {% for eachFile in files %}
      {% for component in eachFile  %}
          {% comment %} <tr><td colspan="2">{{component}}</td></tr> {% endcomment %}
           {% capture thecycle %}{% cycle 'odd', 'even' %}{% endcapture %}
            {% if thecycle == 'even' %}
            <tr>
              <td align="right">{{ component[1].component_type }}</td>
              <td class="title"><a href="#{{ component[1].component_name  | camelcase }}-{{ component[1].component_type }}">{{ component[1].component_name }}</a></td>
            </tr>
        {% endif %}
      {% endfor %}
    {% endfor %}
  {% endfor %}
</tbody>
</table>

<hr />

  {% for files in site.data.errorref %}    
    {% for eachFile in files %}
      {% for component in eachFile  %}
          {% comment %} <tr><td colspan="2">{{component}}</td></tr> {% endcomment %}
           {% capture thecycle %}{% cycle 'odd', 'even' %}{% endcapture %}
            {% if thecycle == 'even' %}
            {% if component[1].component_type == 'adapter' %}
              {% capture heading %}
               Meshery Adapter for {{ component[1].component_name }}
              {% endcapture %}
            {% endif %}
            {% if component[1].component_type == 'client' %}
              {% capture heading %}
               {{ component[1].component_name }} client
              {% endcapture %}
            {% endif %}
            {% if component[1].component_type == 'library' %}
              {% capture heading %}
                {{ component[1].component_name }} {{ component[1].component_type | camelcase }}
              {% endcapture %}
            {% endif %}
            {% if component[1].component_name == 'meshery-server' %}
              {% capture heading %}
                Meshery Server
              {% endcapture %}
            {% endif %}

<h2> {{ heading }} </h2>
  <table>
  <thead>
    <tr>
      <th>Error Name</th>
      <th>Error Code</th>
      <th style="white-space:nowrap; transform-origin:30% 70%; transform: rotate(-90deg);padding:0px;">Severity</th>
      <th>Short Description</th>
      <th>Long Description</th>
      <th>Probable Cause</th>
      <th>Suggested Remediation</th>
    </tr>
  </thead>
  <tbody>
  
    {% for err_code in component[1].errors %}    
        <tr>
          <td >
            <a id="{{component[1].component_name}}-{{err_code[1]["name"]}}">
            {{ err_code[1]["name"] | xml_escape }}
            </a>
          </td>
          <td >{{ err_code[1]["code"] }}</td>
          <td >{{ err_code[1]["severity"]}}</td>
          <td style="max-width:125px;">{{ err_code[1]["short_description"] | xml_escape}}</td>
          <td style="min-width:200px;">{{ err_code[1]["long_description"] | xml_escape }}</td>
          <td style="min-width:200px;">{{ err_code[1]["probable_cause"] | xml_escape }}</td>
          <td style="min-width:200px;">{{ err_code[1]["suggested_remediation"] }}</td>
        </tr>
    {% endfor %}

  </tbody>
  </table>
  <a href="#{{ component[1].component_name }}">Top</a>
  <hr>
  <br>
{% endif %}
{% endfor %}
{% endfor %}
{% endfor %}
