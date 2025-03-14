### END-TO-END TESTS

<%- introMessage %>

**📦 Test Result Summary**

- ✅ <%- passed %> passed
- ❌ <%- failed %> failed
- ⚠️ <%- flaky %> flaked
- ⏩ <%- skipped %> skipped

⌛ _Duration: <%- minutes %> minutes and <%- seconds %> seconds_

<% if (failed > 0) { %>
  **Overall Result**: 👎 Some tests failed.
<% } else { %>
  **Overall Result**: 👍 All tests passed.
<% } %>

<% if (testData && testData.length > 0) { %>

<details>
    <summary>[Show/Hide] Test Result Details</summary>

<table>
  <thead>
    <tr>
      <th>Browser</th>
      <th>Spec</th>
      <th>Test Case</th>
      <th>Tags</th>
      <th>Result</th>
    </tr>
  </thead>
  <tbody>
    <% testData.forEach(function(test) { %>
      <tr>        
        <td><%= test.project %></td>
        <td><%= test.spec %></td>
        <td><%= test.title %></td>
        <td><%= test.tags %></td>
        <td><%= test.status %></td>
      </tr>
    <% }); %>
  </tbody>
</table>

</details>
<% } %>

<!-- To see the full report, please visit our CI/CD pipeline with reporter. -->
