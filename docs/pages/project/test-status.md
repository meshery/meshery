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
- See the [Compatibility Matrix](/installation/compatibility-matrix)

## Relationships Test Results

{% assign relationship_tests = site.data.relationshiptestresult.relationship-v07167-2 %}
{% assign grouped_relationships = relationship_tests | group_by: "name" %}

<style>
/* General Table Styling */
table {
    width: 100%;
    border-collapse: collapse;
}

th, td {
    border: 1px solid #333;
    padding: 8px;
    text-align: left;
}

.accordion-header {
    cursor: pointer;
    background-color: #444;
    color: white;
    font-weight: bold;
    text-align: center;
    font-size: 18px;
}

.accordion-header:hover {
    background-color: #666;
}

.accordion-content {
    display: none;
    background-color: #222;
}

.accordion-content table {
    margin: 0;
    background-color: #333;
}

.accordion-content td {
    color: white;
    padding: 5px 10px;
}
</style>

<table>
    <thead>
        <tr>
            <th>Model</th>
            <th>Meshery Version</th>
            <th>Relationship Kind</th>
            <th>Relationship Type</th>
            <th>Relationship SubType</th>
            <th>Test Result</th>
        </tr>
    </thead>
    <tbody>
        {% for group in grouped_relationships %}
        <!-- Accordion Header Row -->
        <tr class="accordion-header">
            <td colspan="6">{{ group.name }}</td>
        </tr>
        <!-- Hidden Content: Detailed Table -->
        <tr class="accordion-content">
            <td colspan="6">
                <table>
                    <thead>
                        <tr>
                            <th>Meshery Version</th>
                            <th>Relationship Kind</th>
                            <th>Relationship Type</th>
                            <th>Relationship SubType</th>
                            <th>Test Result</th>
                        </tr>
                    </thead>
                    <tbody>
                        {% for item in group.items %}
                        <tr>
                            <td>{{ item.extensionVersion }}</td>
                            <td>{{ item.kind }}</td>
                            <td>{{ item.type }}</td>
                            <td>{{ item.subType }}</td>
                            <td>
                                {% if item.testResultPassed %}
                                    <img src="/assets/img/passing.svg" alt="Pass" />
                                {% else %}
                                    <img src="/assets/img/failing.svg" alt="Fail" />
                                {% endif %}
                            </td>
                        </tr>
                        {% endfor %}
                    </tbody>
                </table>
            </td>
        </tr>
        {% endfor %}
    </tbody>
</table>

<script>
// JavaScript for toggling the accordion content
document.addEventListener("DOMContentLoaded", function () {
    const headers = document.querySelectorAll(".accordion-header");

    headers.forEach(header => {
        header.addEventListener("click", function () {
            const content = this.nextElementSibling;
            if (content.classList.contains("accordion-content")) {
                content.style.display = content.style.display === "table-row" ? "none" : "table-row";
            }
        });
    });
});
</script>
