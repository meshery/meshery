---
title: PR Review Loop
description: Automated peer-review recommendations connecting contributors to related open pull requests.
categories: [contributing]
display-suggested-reading: false
---

# Meshery PR Review Loop

The **PR Review Loop** is an automated community collaboration mechanism designed to encourage reciprocal peer review among contributors. When an author opens or updates a pull request (PR) in `meshery/meshery`, the system surfaces up to three currently open PRs that share technical context, programming language, or architectural components with their change.

---

## 1. Architecture & Design Principles

The PR Review Loop adheres to five core principles:

- **Deterministic**: Suggestions are generated through explainable mathematical scoring based on public PR metadata. No AI/LLM models, embeddings, or external inference APIs are used.
- **Advisory & Non-blocking**: Strictly friendly suggestions. The workflow never blocks PR submission, CI checks, or merges.
- **No Automatic Reviewer Assignment**: Suggestions are informational; contributors are never automatically assigned or requested as formal reviewers.
- **Zero-Noise Notification Policy**: When no eligible candidate PRs exist, no comment is posted. If a previous recommendation becomes stale, the bot cleanly deletes its comment.
- **Zero Untrusted Execution**: Operates in the base-repository context without ever checking out or executing untrusted code from contributor PR branches.

---

## 2. Event Lifecycle & Workflow Orchestration

The PR Review Loop is implemented within the unified repository triage workflow (`.github/workflows/labeler.yml`). This ensures that automated label classification completes before review suggestions are evaluated, eliminating race conditions.

```
pull_request_target event
          │
          ▼
   Check PR Context
          │
    Is triage needed?
    (opened, reopened, synchronize)
       /        \
     Yes         No (labeled, unlabeled, edited, ready_for_review, converted_to_draft)
      │           │
      ▼           │
  Triage Job      │
 (actions/labeler)│
      │           │
      ▼           │
   Success        │
      │           │
      └───┬───────┘
          │
          ▼
 Suggestion Job
  (needs: triage)
          │
          ▼
 Checkout trusted base script
 (ref: base.sha, sparse-checkout)
          │
          ▼
 Fetch latest PR metadata
          │
          ▼
 Filter candidate pool
          │
          ▼
 Score & rank candidates (top 3)
          │
          ▼
 Update / Create / Retire sticky comment
```

### Event Matrix

| Event Action | Triage Job (`actions/labeler`) | Suggestions Job (`suggest-reviews`) | Purpose |
|---|---|---|---|
| `opened` | **Runs** | Runs after triage | Labels newly submitted PR; generates initial recommendations. |
| `reopened` | **Runs** | Runs after triage | Re-applies missing labels; recalculates suggestions. |
| `synchronize` | **Runs** | Runs after triage | Labels new files; refreshes suggestions for updated code. |
| `labeled` | *Skipped* | **Runs directly** | Recalculates suggestions when maintainers manually add labels. |
| `unlabeled` | *Skipped* | **Runs directly** | Recalculates suggestions when labels are removed; retires stale suggestions. |
| `edited` | *Skipped* | **Runs directly** | Recomputes title token matching when PR title is updated. |
| `ready_for_review` | *Skipped* | **Runs directly** | Surfaces suggestions when a draft PR is marked ready for review. |
| `converted_to_draft` | *Skipped* | **Runs directly** | Retires/deletes existing sticky comment while PR is in draft status. |

### Concurrency Control

Execution serialization is decoupled at the job level, with separate per-PR concurrency groups for automated labeling and review suggestions:

```yaml
jobs:
  triage:
    concurrency:
      group: pr-triage-${{ github.event.pull_request.number }}
      cancel-in-progress: false

  suggest-reviews:
    concurrency:
      group: pr-review-suggestions-${{ github.event.pull_request.number }}
      cancel-in-progress: false
```

Decoupling concurrency into dedicated job-level groups guarantees:
- **Triage Isolation**: The `triage` labeler job serializes independently (`pr-triage-${{ github.event.pull_request.number }}`), preventing race conditions from rapid pushes while ensuring suggestion-only events (such as `labeled` or `unlabeled`) cannot cancel or replace pending triage runs.
- **Suggestion Serialization**: The `suggest-reviews` job serializes independently (`pr-review-suggestions-${{ github.event.pull_request.number }}`), ensuring that in-flight comment creation, update, and retirement operations complete sequentially without partial update truncation.
- **State Freshness**: Setting `cancel-in-progress: false` across both groups ensures that every queued job completes, re-querying the GitHub API for the freshest repository metadata upon execution.

---

## 3. Matching Heuristics & Scoring Model

### Mandatory Label Eligibility Gate

Labels serve as the primary relevance signal. A candidate PR is eligible for recommendation **only** if it shares at least one qualifying label with the triggering PR:

- **Language labels**: `language/go`, `language/javascript`, `language/typescript`, `language/shell`
- **Area and Component labels**: `area/*` (e.g. `area/ci`, `area/docs`), `component/*` (e.g. `component/server`, `component/ui`, `component/mesheryctl`)

> [!IMPORTANT]
> Candidate PRs with zero label overlap are **strictly ineligible**, regardless of any title similarity.

### Scoring Formula (V1)

| Criterion | Points | Details |
|---|---|---|
| Shared language label | **+40** | Both PRs share at least one matching `language/*` label |
| Shared area or component label | **+40** | Both PRs share at least one matching `area/*` or `component/*` label |
| Shared unique title keyword | **+10 per word** | Deduplicated technical tokens shared between PR titles |

> [!NOTE]
> **Body Scoring Omission in V1**: PR description text is deliberately excluded from V1 scoring. Because open-source PR templates share standardized boilerplate (conventions, checkboxes, contribution links), body token comparison introduces substantial false similarity (+240 points of template overlap) that swamps architectural labels (+80 max). Meshery PR titles follow the convention `[<component>] descriptive title`, providing clean, high-signal technical tokens without boilerplate contamination.

### Text Normalization

Titles are normalized before token comparison:
1. Converted to lowercase.
2. Punctuation (including hyphens and slashes) converted to whitespace delimiters (`/[^a-z0-9]/g`).
3. Split on whitespace; tokens with length $\le 1$ or matching standard stopwords (`fix`, `update`, `pr`, `add`, `issue`) are discarded.
4. Deduplicated using a `Set` so repeated words do not inflate scores.

---

## 4. Candidate Filtering & Ranking

### Candidate Exclusion Rules

Candidates are retrieved from the 100 most recently updated open PRs targeting `master` (`per_page: 100`, `sort: 'updated'`, `direction: 'desc'`). In an active repository with hundreds of open PRs, capping the candidate search pool to the top 100 prevents API quota exhaustion, avoids iterating through abandoned PRs, and focuses suggestions on actively maintained, responsive pull requests.

Each candidate must pass the following filter pipeline:

1. **Triggering PR Exclusion**: A PR is never recommended to itself (`candidate.number !== prNumber`).
2. **Author Self-PR Exclusion**: PRs authored by the same contributor are excluded (`candidate.user.login !== currentPRAuthor`).
3. **Bot Exclusion**: Automated bot PRs (such as Dependabot, model sync bots, or any login ending in `[bot]`) are excluded.
4. **Draft Exclusion**: Draft PRs (`candidate.draft === true`) are excluded.
5. **State & Branch**: Closed PRs and PRs targeting non-`master` release branches are excluded.
6. **Deduplication**: Candidate numbers are tracked in a `Set` to prevent duplicate evaluations.
7. **Triggering Bot PR Exclusion**: If the triggering PR is bot-authored, the workflow skips execution entirely.

### Deterministic Ranking

Eligible candidates are ordered by:
1. **Relevance score** descending (`b.score - a.score`).
2. **PR number** descending (`b.number - a.number`) as a deterministic tie-breaker favoring newer work.
3. Sliced to at most **three** top recommendations.

---

## 5. Sticky Comment Lifecycle

Recommendations are presented via a single bot-owned comment:

- **Marker**: Comments begin with `<!-- meshery-review-suggestions -->`.
- **Bot Identity Verification**: The workflow strictly checks `comment.user.login === 'github-actions[bot]'`. Participant comments or third-party bot comments containing the marker are never modified or deleted.
- **Zero-Candidate Handling**: If zero eligible recommendations exist and no comment is present, no comment is posted. If recommendations drop to zero (e.g. after label removal or draft conversion), any existing bot comment is deleted.
- **Preserved on Closure / Merge**: When a PR is closed or merged, historical suggestion comments are intentionally preserved in the PR discussion thread rather than deleted. The workflow does not trigger on `closed`, and if an in-flight or queued run evaluates a closed PR, it exits cleanly without deleting the comment.
- **Idempotency**: If the generated comment body matches the existing comment, API update calls are skipped.
- **Duplicate Cleanup**: If multiple bot comments are detected, duplicates are pruned automatically.
- **Untrusted Title Sanitization**: Candidate titles are untrusted input. Titles are truncated to 200 characters, stripped of ASCII control codes and Unicode Bidi override characters (`\u202A`–`\u202E`), escaped for Markdown/HTML syntax, and neutralized for `@` mentions (`@\u200d`) and bare URLs (`:\u200b//`). Only workflow-generated PR links are rendered clickable.

---

## 6. Privileged Workflow Security Boundary

The workflow runs on `pull_request_target` to allow posting suggestion comments on PRs originating from external contributor forks. To prevent arbitrary code execution, the following security invariants are strictly enforced:

1. **Zero Untrusted Checkout**: Contributor branches and fork repositories are **never** checked out.
2. **Trusted Base Script Checkout**: The recommendation script is loaded exclusively from the trusted base repository commit:
   ```yaml
   - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
     with:
       repository: ${{ github.repository }}
       ref: ${{ github.event.pull_request.base.sha }}
       sparse-checkout: |
         .github/scripts/pr-review-suggestions.js
       sparse-checkout-cone-mode: false
       persist-credentials: false
   ```
3. **No Contributor Code Execution**: No dependencies, build steps, npm scripts, or test harnesses from the pull request are executed in the privileged workflow context.
4. **Least-Privilege Scoping**: Permissions are limited strictly to `contents: read` and `pull-requests: write`.
5. **Full Action Pinning**: All GitHub Actions are pinned to full-length commit SHAs.

---

## 7. Testing & Local Validation

The implementation is verified by an automated regression suite located at `.github/scripts/pr-review-suggestions.test.js`:

```bash
# Run unit and regression tests (38 tests passing)
node --test .github/scripts/pr-review-suggestions.test.js

# Validate workflow schemas
actionlint .github/workflows/labeler.yml
actionlint .github/workflows/pr-review-suggestions-tests.yml

# Scan security posture
zizmor .github/workflows/labeler.yml .github/workflows/pr-review-suggestions-tests.yml
```

The test suite runs automatically in CI on pull requests via `.github/workflows/pr-review-suggestions-tests.yml`.
