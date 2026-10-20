---
title: PR Review Loop — Implementation & Operations
description: Technical operations guide for the PR Review Suggestions GitHub Actions workflow.
categories: [contributing]
---

# Meshery PR Review Loop — Implementation & Operations

## 1. Overview

> **"When you open a PR, show you a PR you can help review."**

The Meshery PR Review Loop is an automated community collaboration mechanism designed to encourage reciprocal peer review among contributors. When an author opens a pull request (PR) in `meshery/meshery`, the system surfaces up to three currently open PRs that share technical context, programming language, or architectural components with their change.

Version 1 (V1) is:
- **A lightweight GitHub Actions workflow**: Runs natively in GitHub Actions with minimal resource overhead.
- **Deterministic and explainable**: Relies on simple mathematical scoring rules rather than black-box algorithms.
- **Metadata-driven**: Compares public PR metadata (labels, titles, descriptions).
- **A pure suggestion mechanism**: Friendly, advisory guidance for contributors.
- **Completely non-blocking**: Never blocks PR submission, CI checks, or merging.

What V1 does **NOT** do:
- Does **not** automatically assign or request reviewers.
- Does **not** block PR creation or merges.
- Does **not** require contributors to review another PR before proceeding.
- Does **not** use AI, large language models (LLMs), or vector embeddings.

---

## 2. Files Involved

The PR Review Loop feature consists of the following components:

### `.github/labeler.yml`
Defines automated PR label rules for the repository. The Review Loop relies on two specific label categories as primary relevance signals:
- **Language labels**:
  - `language/go`
  - `language/javascript`
  - `language/typescript`
  - `language/python`
  - `language/ruby`
  - `language/shell`
- **Area & Component labels**:
  - `area/*` (e.g., `area/docs`, `area/ci`, `area/models`, `area/lifecycle`, `area/tests`)
  - `component/*` (e.g., `component/ui`, `component/server`, `component/mesheryctl`, `component/extensions`, `component/database`, `component/policy`)

### `.github/workflows/pr-review-suggestions.yml`
The core GitHub Actions workflow implementing the Review Loop logic:
- Listens for PR events (`opened`, `reopened`, `labeled`, `unlabeled`) on `master`.
- Listens for `workflow_run` completion of the `Pull Request Labeler` workflow.
- Resolves and validates PR metadata through the GitHub REST API.
- Applies candidate filtering (excluding self, drafts, and closed PRs).
- Calculates deterministic relevance scores.
- Ranks candidates and selects up to three top matches.
- Creates or updates a single bot-owned sticky comment on the PR.

### `.github/scripts/pr-review-suggestions.test.js`
Independent, dependency-free test suite utilizing Node's built-in `node:test` runner. Directly exercises the workflow script against unit and mocked API integration scenarios.

---

## 3. Workflow Trigger & Concurrency

```yaml
on:
  # Runs in the base-repository context. Never check out or execute code from the PR head.
  pull_request_target:
    types: [opened, reopened, labeled, unlabeled]
    branches:
      - master
  workflow_run:
    workflows: ["Pull Request Labeler"]
    types: [completed]

concurrency:
  group: pr-review-suggestions-${{ github.event.pull_request.number || (github.event.workflow_run && github.event.workflow_run.pull_requests[0] && github.event.workflow_run.pull_requests[0].number) || github.run_id }}
  cancel-in-progress: true
```

### Trigger Architecture & Event Semantics

1. **`pull_request_target`**:
   - Executes in the base repository context (`meshery/meshery`). This allows the workflow to authenticate with sufficient permissions to create and update comments on PRs submitted from external forks.
   - Listens on `[opened, reopened, labeled, unlabeled]`:
     - `opened` / `reopened`: Surfaces initial recommendations immediately upon PR submission.
     - `labeled` / `unlabeled`: Recomputes suggestions whenever labels are manually added or removed. If a contributor removes the only qualifying label, the sticky comment updates to remove obsolete recommendations.

2. **Automated Labeling Handoff (`workflow_run`)**:
   - The repository's `Pull Request Labeler` workflow (`.github/workflows/labeler.yml`) applies automated labels using `GITHUB_TOKEN`.
   - By GitHub Actions design, events generated using `GITHUB_TOKEN` do **not** trigger downstream workflow runs (preventing accidental infinite loops). Consequently, relying solely on `pull_request_target: types: [labeled]` leaves suggestions stale after automated triage.
   - To resolve this, the workflow triggers on `workflow_run` when the labeler completes.
   - **PR Disambiguation & Validation**:
     - When `workflow_run` fires, the workflow checks `run.conclusion === 'success'`.
     - It identifies the target PR from `workflow_run.pull_requests[0]` or queries open PRs matching `workflow_run.head_sha`.
     - It validates that the PR targets `master`, is currently `open`, and retrieves verified metadata directly from `github.rest.pulls.get`.
     - If the PR cannot be resolved unambiguously, it logs an informational message and exits safely without modifying any comments.

3. **Concurrency & Race Condition Prevention**:
   - Concurrency groups cancel redundant in-flight runs when events fire in rapid succession.
   - The workflow also compares the computed comment content against existing comment text (`existingComment.body.trim() === commentBody.trim()`). If no recommendations changed, redundant API updates are skipped.

---

## 4. Permissions & Security Boundary

### Least-Privilege Scoping

```yaml
permissions:
  pull-requests: write
```

- **`pull-requests: write`**: Grants the narrowest necessary permission to read PR metadata and list, create, or update comments on pull requests (`issues.createComment`, `issues.updateComment`, `issues.listComments` are scoped to PR issues).
- **No `issues: write`**: Unnecessary repository-wide issue modification permissions are omitted.
- **No `contents: read` / `contents: write`**: The workflow does not read or write repository files.

### Security Invariants for `pull_request_target`

Security scanners (such as `zizmor`) flag `pull_request_target` as a `dangerous-triggers` finding because unhardened workflows frequently check out untrusted contributor code with privileged repository tokens.

Meshery enforces strict architectural invariants to eliminate this threat vector:

1. **Zero Code Checkout**: The workflow does **not** invoke `actions/checkout`. PR branches and external forks are never cloned to the runner.
2. **Zero Code Execution**: No scripts, dependencies, npm packages, or build steps from the contributor's PR are imported, sourced, or executed.
3. **No Artifact Ingestion**: No workflow artifacts or files generated by previous jobs are downloaded or evaluated.
4. **Candidate Title Sanitization**: Candidate PR titles are untrusted user input. Titles are processed with `escapeMarkdown()` to escape Markdown link syntax (`[text](url)`), backticks, asterisks, underscores, HTML tags, and table pipes. This prevents malicious titles from injecting phishing links or breaking comment formatting.
5. **Sticky Comment Ownership**: When searching for an existing comment to update, the workflow asserts that `comment.user.login === 'github-actions[bot]'` in addition to checking the `<!-- meshery-review-suggestions -->` marker. Comments authored by human contributors or other bots that happen to contain the marker are strictly ignored.
6. **Graceful Error Containment**: All API calls are wrapped in top-level error handling. Network or rate-limiting failures log a non-fatal `core.warning` and exit cleanly, never failing contributor CI checks.

---

## 5. Candidate Discovery & Scoring Algorithm

The workflow queries the GitHub API for candidate PRs:
```javascript
const response = await github.rest.pulls.list({
  owner: context.repo.owner,
  repo: context.repo.repo,
  state: 'open',
  per_page: 100,
  sort: 'updated',
  direction: 'desc'
});
```

### Filtering Pipeline
1. **Self-PR Exclusion**: `candidate.number !== prNumber`
2. **Open PRs Only**: `candidate.state === 'open'`
3. **Draft Exclusion**: `!candidate.draft`
4. **Deduplication**: `!seenNumbers.has(candidate.number)`
5. **Mandatory Label Gate**: Must share at least one `language/*`, `area/*`, or `component/*` label with current PR.

### Scoring Breakdown
- **Language Match**: `+40` (shares at least one `language/*` label)
- **Area / Component Match**: `+40` (shares at least one `area/*` or `component/*` label)
- **Title Keyword Overlap**: `+10` per unique matching technical keyword
- **Body Keyword Overlap**: `+5` per unique matching technical keyword

### Ranking & Selection
- **Primary Sort**: Score descending (`b.score - a.score`)
- **Tie-breaker**: PR number descending (`b.number - a.number`), prioritizing newer pull requests.
- **Top 3**: Returns at most 3 candidates (`slice(0, 3)`).

---

## 6. Testing & Local Validation

The implementation includes a zero-dependency test suite located at `.github/scripts/pr-review-suggestions.test.js`:

```bash
# Run unit and integration tests
node --test .github/scripts/pr-review-suggestions.test.js

# Validate workflow schema
actionlint .github/workflows/pr-review-suggestions.yml

# Verify workflow security posture
zizmor .github/workflows/pr-review-suggestions.yml
```
