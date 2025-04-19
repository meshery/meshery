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

<!-- To see the full report, please visit our CI/CD pipeline with reporter. -->