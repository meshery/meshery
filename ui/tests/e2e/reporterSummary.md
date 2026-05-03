### END-TO-END TESTS

<%- introMessage %>

**📦 Summary**

- ✅ <%- passed %> passed
- ❌ <%- failed %> failed
- ⚠️ <%- flaky %> flaked
- ⏩ <%- skipped %> skipped

⌛ _Duration: <%- minutes %> minutes and <%- seconds %> seconds_

<%- failed > 0 ? "**Overall Result**: 👎 Some tests failed." : "**Overall Result**: 👍 All tests passed." %>

<% if (failedTestCount > 0) { %>

**❌ Failed Test Summary**

Review the compact failure summary below. Full diagnostics are kept in workflow artifacts.

<%= failedTestDetails %>

<% } %>

<% if (testTable) { %>

<details>
    <summary>[Show/Hide] Full Test Matrix</summary>
    <div markdown="1">

<%- testTable %>

</div>
</details>
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
