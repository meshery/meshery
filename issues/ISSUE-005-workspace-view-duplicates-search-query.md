### Current Behavior
`mesheryctl workspace view <workspace-name> --orgId <org-id>` builds a URL that already includes `search=<name>`, then passes the same name again through `DisplayDataAsync.SearchTerm`.

During pagination, `HandlePaginationAsync` appends `search` whenever `SearchTerm` is set, so the final request contains the same search filter twice.

### Expected Behavior
Workspace lookup by name should send one consistent `search` parameter.

The command should either keep `search` in the URL or pass it through `SearchTerm`, but not both.

### Screenshots/Logs
Relevant code path:

```go
// mesheryctl/internal/cli/root/workspaces/view.go
viewUrlValue.Add("orgID", workspaceViewFlagsProvided.OrgID)
viewUrlValue.Add("search", workspaceNameOrID)

urlPath = fmt.Sprintf("%s?%s", workspacesApiPath, viewUrlValue.Encode())
displayData = display.DisplayDataAsync{
	UrlPath:    urlPath,
	SearchTerm: workspaceNameOrID,
}
```

Pagination adds `search` again when `SearchTerm` is non-empty:

```go
// mesheryctl/internal/cli/pkg/display/pagination.go
if displayData.SearchTerm != "" {
	pagesQuerySearch.Set("search", displayData.SearchTerm)
}
```

### Environment

- **Host OS:** Mac Linux Windows
- **Platform:** Docker or Kubernetes
- **Meshery Server Version:** stable-v
- **Meshery Client Version:** stable-v

<!-- Optional 
### To Reproduce
1. Run `mesheryctl workspace view my-workspace --orgId <org-id>`.
2. Capture the outgoing HTTP request in a test or proxy.
3. Inspect the query string.
4. See `search` appended twice in the final URL.
-->

---
### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://meshery.io/community#handbook)
- 🛠 [Meshery Build & Release Strategy](https://docs.meshery.io/project/contributing/build-and-release)
- 📚 [Instructions for contributing to documentation](https://docs.meshery.io/project/contributing/contributing-docs)
   - Meshery documentation [site](https://docs.meshery.io/) and [source](https://github.com/meshery/meshery/tree/master/docs)
- 🎨 Wireframes and [designs for Meshery UI](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI) in Figma [(open invite)](https://www.figma.com/team_invite/redeem/GvB8SudhEOoq3JOvoLaoMs)
- 📺 [Self-paced Contributor Trainings](https://meshery.io/talks-and-trainings#trainings)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://meshery.io/community#community-forums) and [Community Slack](https://slack.meshery.io)
