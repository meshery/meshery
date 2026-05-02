---
name: Pull Request Reviewer
description: Reviews newly opened pull requests from trusted repository contributors and leaves concise, polite feedback as PR review comments
on:
  pull_request:
    types: [opened]
  skip-if-match: ${{ !contains(fromJSON('["dependabot[bot]", "copilot", "gemini", "codex", "claude"]'), github.actor) }}
permissions:
  contents: read
  issues: read
  pull-requests: read
engine: copilot
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

# Pull request review

Review the newly opened pull request and leave feedback directly on the PR.

## Activation

- This workflow runs on `pull_request_target`, so the compiled workflow intentionally limits activation to pull requests opened by trusted repository contributors with `admin`, `maintainer`, or `write` access.
- Pull requests from forks and other external contributors are skipped by design so the workflow can safely use repository-scoped secrets and post review feedback.

## Scope

- Review the PR title, description, changed files, and diff.
- Focus on correctness, security, reliability, maintainability, tests, and user-visible regressions.
- Ignore minor style nits, formatting-only issues, and unrelated concerns.
- Skip generated files, lock files, snapshots, vendored content, and other derived artifacts unless they reveal a real defect.

## Commenting rules

- Always leave feedback in the pull request as one or more comments.
- Keep feedback polite, concise, and actionable.
- Prefer inline review comments for specific issues tied to changed lines.
- If you identify a change that should be made, explain exactly what should change and why.
- Limit feedback to the most important issues; do not overwhelm the author with low-value comments.
- If you do not find actionable issues, post a brief PR comment saying no blocking issues were found and optionally note one positive observation.
- Never approve the PR and never request changes.

## Process

1. Read the pull request metadata and inspect the diff.
2. Identify only substantive issues worth telling the author about.
3. Create up to five inline review comments when specific file-level feedback is warranted.
4. Post one concise overall PR comment that references the most important findings, or states that no blocking issues were found.

## Usage

Edit this file and run `gh aw compile pull-request-reviewer` to regenerate the lock file.
