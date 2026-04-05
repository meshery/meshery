# [mesheryctl] filter view accepts and ignores arbitrary positional arguments

### Current Behavior
`mesheryctl filter view` uses `cobra.ArbitraryArgs`, so commands like `mesheryctl filter view arg1 arg2 arg3` are silently accepted without any validation error.

This is inconsistent with other view commands (`design view` uses `cobra.MaximumNArgs(1)`, `workspace view` validates exactly one argument).

### Expected Behavior
`mesheryctl filter view` should validate positional arguments consistently with other CLI view commands. Extra arguments should be rejected with a usage error.

### Screenshots/Logs
Relevant code path:

```go
// mesheryctl/internal/cli/root/filter/view.go
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View filter(s)",
	Args: cobra.ArbitraryArgs, // ← no validation
	RunE: func(cmd *cobra.Command, args []string) error {
		// ...
	},
}
```

Compare with `design view` which properly validates:

```go
// mesheryctl/internal/cli/root/design/view.go
var viewCmd = &cobra.Command{
	Use:  "view design name",
	Args: cobra.MaximumNArgs(1), // ← proper validation
}
```

### Proposed Fix
1. Replace `cobra.ArbitraryArgs` with proper argument validation
2. Add `Args` function that validates 0 or 1 args (matching `--all` flag behavior)
3. Add corresponding unit tests

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
