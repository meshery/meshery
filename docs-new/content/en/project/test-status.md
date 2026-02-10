---
title: End-to-End Test Status
# permalink: project/contributing/test-status
description: Status reports of Meshery's various test results.
categories: [contributing]
---

This page contains results of tests performed in the development of Meshery.

**Resources:**

- See [Contributing to Meshery's End-to-End Tests](/project/contributing/contributing-ui-tests)
- See the [Meshery Test Plan](https://docs.google.com/spreadsheets/d/13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/edit?gid=0#gid=0) for a list of test cases.
- See the [Compatibility Matrix](/installation/compatibility-matrix)

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
    font-weight: normal;
    text-align: center;
    font-size: 1rem;
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


## Latest E2E Test Report

[🔗 View Build UI and Server Workflow](https://github.com/meshery/meshery/actions/workflows/build-ui-and-server.yml?query=branch%3Amaster)


{{ partial "test-report.md" . }}


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
