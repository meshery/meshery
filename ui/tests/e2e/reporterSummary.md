Test run results

---

## <%- introMessage %>

Summary:

- ⌛ Duration of test run: <%- minutes %> minutes and <%- seconds %> seconds
- 📦 Tests results:
  - ✅ <%- passed %>
  - ❌ <%- failed %>
  - ⏩ <%- skipped %>

<%- failed > 0 ? "👎 Some tests failed!" : "👍 All tests passed successfully!" %>

<% if (failsMessage) { %>

<details>
    <summary>Click Here for more details</summary>
    <%- failsMessage %>
</details>
<% } %>

To see the full report, please visit our CI/CD pipeline with reporter.`;
