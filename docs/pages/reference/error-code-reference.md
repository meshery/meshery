---
layout: default
title: Error Code Reference
abstract: "A reference of Error codes in Meshery ecosystem"
permalink: reference/error-codes
type: Reference
---
## Meshery Error Code Reference

Error codes will fall into specific, predetermined number range in order to categorize the error by functional area / component, so that users might easily reference and troubleshoot error codes. Each error code is to be used as a unique identifier of the problem faced. Ranges defined for error codes:

<br>

## Error Codes Categories

<table>
<thead>
  <tr>
    <th>Category</th>
    <th>Error Code Range</th>
  </tr>
</thead>
<tbody>
  {% for err_cat in site.data.error_references.err_categories %}    
    <tr>
      <td >{{ err_cat.category }}</td>
      <td>{{ err_cat.range_start }} to {{ err_cat.range_end }} </td>
    </tr>
  {% endfor %}
</tbody>
</table>

<br>
<hr>

{% comment %} the error ref of each cat code {% endcomment %}
{% for err_cat in site.data.error_references.err_categories %}    

  {% for err_ref in site.data.error_references.err_descriptions %}    
    {% if err_ref.cat_code== err_cat.cat_code %}
      {% assign errorCodesExist = true %}
      {% break %}
    {% else %}
      {% assign errorCodesExist = false %}
    {% endif %}
  {% endfor %}

{% if errorCodesExist == true %}
  
## {{ err_cat.range_start }} {{ err_cat.category }}

  <table>
  <thead>
    <tr>
      <th>Error Code</th>
      <th>Severity</th>
      <th>Short Description</th>
      <th>Long Description</th>
      <th>Probable Cause</th>
      <th>Suggested Remediation</th>
    </tr>
  </thead>
  <tbody>
    {% for err_ref in site.data.error_references.err_descriptions %}    
      {% if err_ref.cat_code== err_cat.cat_code %}
        <tr>
          <td >{{ err_ref.code }}</td>
          <td >{{ err_ref.severity }}</td>
          <td >{{ err_ref.short_description }}</td>
          <td >{{ err_ref.long_description }}</td>
          <td >{{ err_ref.probable_cause }}</td>
          <td >{{ err_ref.suggested_remediation }}</td>
        </tr>
      {% endif %}
    {% endfor %}
  </tbody>
  </table>
  {% endif %}

{% endfor %}

