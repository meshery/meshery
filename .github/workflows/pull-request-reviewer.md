---
name: Pull Request Reviewer
description: Reviews newly opened pull requests targeting master and leaves concise, polite feedback as PR comments
on:
  pull_request:
    branches: [master]
    types: [opened]
permissions:
  contents: read
  issues: read
  pull-requests: read
engine: copilot
imports:
  - .github/agents/meshery-code-contributor.md
timeout-minutes: 20
tools:
  github:
    toolsets: [default]
safe-outputs:
  add-comment:
    max: 1
  create-pull-request-review-comment:
    max: 5
  missing-data: false
  missing-tool: false
  noop: false
  report-failure-as-issue: false
  report-incomplete: false
network:
  allowed:
    - defaults
---

Offer concise, polite, feedback on the following pull request(s). Add comment(s) to pull request. Include any changes needed, if any.
