### END-TO-END TESTS

<%- introMessage %>

**📦 Test Result Summary**

- ✅ <%- passed %> passed
- ❌ <%- failed %> failed
- ⚠️ <%- flaky %> flaked
- ⏩ <%- skipped %> skipped

⌛ _Duration: <%- minutes %> minutes and <%- seconds %> seconds_

<%- failed > 0 ? "**Overall Result**: 👎 Some tests failed." : "**Overall Result**: 👍 All tests passed." %>

<% if (testTable) { %>

<details>
    <summary>[Show/Hide] Test Result Details</summary>
    <div markdown="1">

<%- testTable %>

</div>
</details>
<% } %>

<% if (failedTestCount > 0) { %>

**❌ Failed Test Details**

The following tests failed. Click on each to see the failure reason:

<%= failedTestDetails %>

<% } %>

<% if (relationshipTestCount > 0) { %>

**🔗 Relationship Tests**

<details>
    <summary>[Show/Hide] Relationship Test Details (<%- relationshipTestCount %> tests)</summary>
    <div markdown="1">

<%- relationshipTestTable %>

</div>
</details>
<% } %>

<!-- To see the full report, please visit our CI/CD pipeline with reporter. -->
