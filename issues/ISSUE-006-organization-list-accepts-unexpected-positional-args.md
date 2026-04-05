### Current Behavior
`mesheryctl organization list` does not validate positional arguments, so extra arguments are silently accepted instead of being rejected.

That means commands like `mesheryctl organization list random` still execute the list request, even though `list` is documented as taking no positional arguments.

### Expected Behavior
`mesheryctl organization list` should reject unexpected positional arguments and return a usage error.

The command behavior should match other list commands in the CLI that use `cobra.ExactArgs(0)` or equivalent validation.

### Screenshots/Logs
Relevant command definition:

```go
// mesheryctl/internal/cli/root/organizations/list.go
var listOrgCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered organizations",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		pagesize, _ := cmd.Flags().GetInt("pagesize")
		count, _ := cmd.Flags().GetBool("count")
		// executes without checking args
		return display.ListAsyncPagination(data, processOrgData)
	},
}
```

By contrast, other commands explicitly validate zero args:

```go
// mesheryctl/internal/cli/root/system/token.go
Args: cobra.ExactArgs(0),
```

### Environment

- **Host OS:** Mac Linux Windows
- **Platform:** Docker or Kubernetes
- **Meshery Server Version:** stable-v
- **Meshery Client Version:** stable-v

<!-- Optional 
### To Reproduce
1. Run `mesheryctl organization list unexpected-arg`.
2. Observe that the command still lists organizations.
3. Compare with other CLI commands that reject extra args.
4. See the missing argument validation on the list subcommand.
-->

---
### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://meshery.io/community#handbook)
- 🛠 [Meshery Build & Release Strategy](https://docs.meshery.io/project/contributing/build-and-release)
- 📚 [Instructions for contributing to documentation](https://docs.meshery.io/project/contributing/contributing-docs)
   - Meshery documentation [site](https://docs.meshery.io/) and [source](https://github.com/meshery/meshery/tree/master/docs)
- 🎨 Wireframes and [designs for Meshery UI](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI) in Figma [(open invite)](https://www.figma.com/team_invite/redeem/GvB8SudhEOoq3JOvoLaoMs)
- 📺 [Self-paced Contributor Trainings](https://meshery.io/talks-and-trainings#trainings)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://meshery.io/community#community-forums) and [Community Slack](https://slack.meshery.io)
