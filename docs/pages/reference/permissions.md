---
layout: default
title: Permissions
permalink: reference/permissions
redirect_from: reference/permissions
language: en
list: exclude
abstract: List of default permissions.
---
<style>
      .scrollable-container {
        overflow-x: auto; 
      }
      table {
        border-collapse: collapse;
        width: auto; 
      }
    </style>

<!-- <table>
  {% for row in site.data.keys %}
    {% if forloop.first %}
    <tr>
      {% for pair in row %}
        <th>{{ pair[0] }}</th>
      {% endfor %}
    </tr>
    {% endif %}
{% tablerow pair in row %}
      {{ pair[1] }}
    {% endtablerow %}
  {% endfor %}
</table> -->

<!-- Assuming the CSV data is stored in site.data.data -->
<table>
  <thead>
    <!-- If you want to include headers from the first row -->
    <tr>
      {% assign all_rows = site.data.keys %}
      {% assign header_row = all_rows[1] %}
      {% for header in header_row %}
        <th>{{ header }}</th>
      {% endfor %}
    </tr>
  </thead>
  <tbody>
    <!-- Start loop from the second row -->
    {% for row in all_rows %}
      {% unless forloop.index == 1 %}
        <tr>
          <!-- Access each value in the row -->
          {% for value in row %}
            <td>{{ value }}</td>
          {% endfor %}
        </tr>
      {% endunless %}
    {% endfor %}
  </tbody>
</table>
