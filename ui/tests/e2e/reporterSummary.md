# E2E Playwright Test Results

---

## <%- introMessage %>

Summary:

- ⌛ Duration of test run: <%- minutes %> minutes and <%- seconds %> seconds
- 📦 Tests results:
  - ✅ <%- passed %>
  - ❌ <%- failed %>
  - ⚠️ <%- flaky %>
  - ⏩ <%- skipped %>

<%- failed > 0 ? "👎 Some tests failed!" : "👍 All tests passed successfully!" %>

<% if (testTable) { %>

<details>
    <summary>Click Here to see test table</summary>

<%- testTable %>

</details>
<% } %>

To see the full report, please visit our CI/CD pipeline with reporter.
