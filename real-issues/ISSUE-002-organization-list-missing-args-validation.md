### Current Behavior
`mesheryctl organization list` does not validate positional arguments. Extra arguments like `mesheryctl organization list random` are silently accepted and the list request still executes.

### Expected Behavior
`mesheryctl organization list` should reject unexpected positional arguments and return a usage error, matching other list commands that use `cobra.NoArgs` or `cobra.ExactArgs(0)`.

### Screenshots/Logs
Relevant code path:

```go
// mesheryctl/internal/cli/root/organizations/list.go
var listOrgCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered organizations",
	// No Args field — accepts anything silently
	RunE: func(cmd *cobra.Command, args []string) error {
		// ...
	},
}
```

By contrast, other commands explicitly validate:

```go
// mesheryctl/internal/cli/root/system/start.go
Args: cobra.NoArgs,
```

### Proposed Fix
1. Add `Args: cobra.NoArgs` to `listOrgCmd`
2. Add test case for rejected positional arguments

### Environment

- **Host OS:** Mac Linux Windows
- **Platform:** Docker or Kubernetes
- **Meshery Server Version:** stable-v
- **Meshery Client Version:** stable-v

---
### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://meshery.io/community#handbook)
- 🛠 [Meshery Build & Release Strategy](https://docs.meshery.io/project/contributing/build-and-release)
- 📚 [Instructions for contributing to documentation](https://docs.meshery.io/project/contributing/contributing-docs)
   - Meshery documentation [site](https://docs.meshery.io/) and [source](https://github.com/meshery/meshery/tree/master/docs)
- 🎨 Wireframes and [designs for Meshery UI](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI) in Figma [(open invite)](https://www.figma.com/team_invite/redeem/GvB8SudhEOoq3JOvoLaoMs)
- 📺 [Self-paced Contributor Trainings](https://meshery.io/talks-and-trainings#trainings)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://meshery.io/community#community-forums) and [Community Slack](https://slack.meshery.io)
