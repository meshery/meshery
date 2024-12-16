---
layout: default
title: End-to-End Test Status
permalink: project/contributing/test-status
abstract: Status reports of Meshery's various test results.
language: en
type: project
category: contributing
list: include
---

This page contains results of tests performed in the development of Meshery.

**Resources:**

- See [Contributing to Meshery's End-to-End Tests](/project/contributing/contributing-ui-tests)
- See the [Meshery Test Plan](https://docs.google.com/spreadsheets/d/13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/edit?gid=0#gid=0) for a list of test cases.

## Relationships Test Results

{% assign relationship_tests = site.data.relationshiptestresult.relationship-v07167-2 %}

{% assign grouped_relationships = relationship_tests | group_by: "name" %}

<table class="table table-bordered">
    <thead>
        <tr>
            <th>Name</th>
            <th>Extension Version</th>
            <th>Kind</th>
            <th>Type</th>
            <th>SubType</th>
            <th>Test Result</th>
        </tr>
    </thead>
    <tbody>
        {% for group in grouped_relationships %}
        {% for item in group.items %}
        <tr>
            <td>{{ item.name }}</td>
            <td>{{ item.extensionVersion }}</td>
            <td>{{ item.kind }}</td>
            <td>{{ item.type }}</td>
            <td>{{ item.subType }}</td>
            <td>
                {% if item.testResultPassed %}
                    <img src="/assets/img/passing.svg" />
                {% else %}
                    <img src="/assets/img/failing.svg" />
                {% endif %}
            </td>
        </tr>
        {% endfor %}
        {% endfor %}
    </tbody>
</table>

