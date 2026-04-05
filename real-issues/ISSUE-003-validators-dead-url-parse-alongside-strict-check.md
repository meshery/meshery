### Current Behavior
In `server/models/validators.go`, the `SMPPerformanceTestConfigValidator` function has two URL validation paths. Line 37 still uses the permissive `url.Parse` which accepts relative/malformed URLs like `not-a-url`, while lines 37-40 use the stricter `url.ParseRequestURI` with scheme/host checks.

The permissive `url.Parse` call on line 37 is now dead code that will never reject anything the stricter check wouldn't catch first, but its presence is confusing and could be copied into new code.

### Expected Behavior
The validator should only use the strict validation path (`url.ParseRequestURI` with scheme and host checks). The permissive `url.Parse` call should be removed to avoid confusion.

### Screenshots/Logs
Current code:

```go
// server/models/validators.go
for _, rawURL := range testClient.EndpointUrls {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ErrValidURL
	}
}
```

Previously the code used the permissive `url.Parse`:

```go
// Old code (before fix)
for _, URL := range testClient.EndpointUrls {
	if _, err := url.Parse(URL); err != nil {
		return ErrValidURL
	}
}
```

### Proposed Fix
1. Ensure only `url.ParseRequestURI` (or equivalent strict check) is used
2. Verify no other validators in the file use the permissive `url.Parse`
3. Add edge-case tests (empty string, path-only, scheme-only)

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
