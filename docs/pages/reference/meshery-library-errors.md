---
layout: default
title: Meshery Lib Errors Reference
abstract: "Meshery's library Errors"
permalink: reference/meshkit-lib-errors
type: Reference
---
## Meshkit Library Introduction

Paragraph about Meshkit Library Errors.

<br>

## Meshkit Library Errors

<table>
<thead>
  <tr>
    <th>Component</th>
    <th>Type</th>
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
              <td><a href="#{{ component[1].component_name }}-{{ component[1].component_type }}">{{ component[1].component_name }}</a></td>
              <td>{{ component[1].component_type }}</td>
            </tr>
        {% endif %}
      {% endfor %}      
    {% endfor %}
  {% endfor %}
</tbody>
</table>

<br>
<hr>


  {% for files in site.data.errorref %}    
    {% for eachFile in files %}
      {% for component in eachFile  %}
          {% comment %} <tr><td colspan="2">{{component}}</td></tr> {% endcomment %}
           {% capture thecycle %}{% cycle 'odd', 'even' %}{% endcapture %}
            {% if thecycle == 'even' %}    
## {{ component[1].component_name }} {{ component[1].component_type }}

  <table>
  <thead>
    <tr>
      <th>Error Name</th>
      <th>Error Code</th>
      <th>Severity</th>
      <th>Short Description</th>
      <th>Long Description</th>
      <th>Probable Cause</th>
      <th>Suggested Remediation</th>
    </tr>
  </thead>
  <tbody>
  

    {% for err_code in component[1].errors %}    
        <tr>
          <td >{{ err_code[1]["name"] }}</td>
          <td >{{ err_code[1]["code"] }}</td>
          <td >{{ err_code[1]["severity"] }}</td>
          <td >{{ err_code[1]["short_description"] }}</td>
          <td >{{ err_code[1]["long_description"] }}</td>
          <td >{{ err_code[1]["probable_cause"] }}</td>
          <td >{{ err_code[1]["suggested_remediation"] }}</td>
        </tr>
    {% endfor %}

  </tbody>
  </table>
  <a href="#meshkit-library-errors">Top</a>
  <hr>
  <br>
{% endif %}
{% endfor %}
{% endfor %}
{% endfor %}
