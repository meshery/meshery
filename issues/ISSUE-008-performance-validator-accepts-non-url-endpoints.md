### Current Behavior
`SMPPerformanceTestConfigValidator` uses `url.Parse` to validate performance test endpoints, but `url.Parse` accepts many non-URL strings as relative paths without returning an error.

As a result, malformed endpoint values can pass validation and only fail later during execution.

### Expected Behavior
Performance test validation should reject endpoint values that are not valid absolute HTTP or HTTPS URLs.

Invalid endpoint configuration should fail fast during validation instead of surfacing later as runtime request errors.

### Screenshots/Logs
Relevant validator logic:

```go
// server/models/validators.go
for _, URL := range testClient.EndpointUrls {
	if _, err := url.Parse(URL); err != nil {
		return ErrValidURL
	}
}
```

`url.Parse` is too permissive for endpoint validation because strings like `not-a-url` parse successfully as relative references.

The validator already treats these values as concrete test endpoints:

```go
// server/models/validators.go
if len(testClient.EndpointUrls) < 1 {
	return ErrTestEndpoint
}
```

### Environment

- **Host OS:** Mac Linux Windows
- **Platform:** Docker or Kubernetes
- **Meshery Server Version:** stable-v
- **Meshery Client Version:** stable-v

<!-- Optional 
### To Reproduce
1. Build a performance test config with an endpoint like `not-a-url`.
2. Pass it through `SMPPerformanceTestConfigValidator`.
3. Observe that validation succeeds.
4. See the invalid endpoint fail only later when the test tries to execute requests.
-->

---
### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://meshery.io/community#handbook)
- 🛠 [Meshery Build & Release Strategy](https://docs.meshery.io/project/contributing/build-and-release)
- 📚 [Instructions for contributing to documentation](https://docs.meshery.io/project/contributing/contributing-docs)
   - Meshery documentation [site](https://docs.meshery.io/) and [source](https://github.com/meshery/meshery/tree/master/docs)
- 🎨 Wireframes and [designs for Meshery UI](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI) in Figma [(open invite)](https://www.figma.com/team_invite/redeem/GvB8SudhEOoq3JOvoLaoMs)
- 📺 [Self-paced Contributor Trainings](https://meshery.io/talks-and-trainings#trainings)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://meshery.io/community#community-forums) and [Community Slack](https://slack.meshery.io)
