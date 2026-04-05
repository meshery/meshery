**Notes for Reviewers**

- This PR fixes the SMP performance test validator so invalid endpoint strings like `not-a-url` are rejected during validation instead of slipping through to runtime.
- Scope is intentionally limited to `server/models/validators.go` and `server/models/validators_test.go`.
- This PR fixes #<issue-number>

**What changed**

- Replaced permissive endpoint parsing with stricter URL validation in `SMPPerformanceTestConfigValidator`.
- Added regression coverage for:
  - rejecting relative / malformed endpoint strings
  - accepting valid absolute endpoints

**Validation**

- Ran:
  `go test ./server/models -run '^TestSMPPerformanceTestConfigValidator' -count=1`

**[Signed commits](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md#signing-off-on-commits-developer-certificate-of-origin)**
- [x] Yes, I signed my commits.
