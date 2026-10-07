# Meshery PR Review Loop — Implementation

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

*Note: Other labels present in the repository do not automatically participate in Review Loop scoring.*

### `.github/workflows/pr-review-suggestions.yml`
The core GitHub Actions workflow implementing the Review Loop logic:
- Listens for PR events (`opened`, `reopened`, `labeled`).
- Fetches triggering PR details and current open PR candidates.
- Applies candidate filtering (excluding self, drafts, and closed PRs).
- Calculates deterministic relevance scores.
- Ranks candidates and selects up to three top matches.
- Creates or updates a single sticky comment on the PR.

### `docs/pr-review-loop/matching-design.md`
The architectural design specification for V1 matching, scoring heuristics, stopwords, and boundaries.

### `docs/pr-review-loop/implementation.md`
This document — the technical operations and maintenance guide for contributors and maintainers.

---

## 3. Workflow Trigger & Concurrency

The workflow is configured with the following trigger and concurrency controls:

```yaml
on:
  pull_request_target:
    types: [opened, reopened, labeled]
    branches:
      - master

concurrency:
  group: pr-review-suggestions-${{ github.event.pull_request.number }}
  cancel-in-progress: true
```

### Rationale
- **`pull_request_target`**: Executes in the context of the base repository (`meshery/meshery`), allowing the workflow to post comments on fork PRs securely without checking out or executing untrusted PR code.
- **`branches: [master]`**: Meshery development targets `master`.
- **`types: [opened, reopened, labeled]`**:
  - `opened` / `reopened`: Runs immediately upon PR submission.
  - `labeled`: Meshery's PR labeler workflow (`labeler.yml`) applies labels asynchronously after PR creation. Triggering on `labeled` allows the Review Loop to re-evaluate and update its recommendations once labels become available.
- **`concurrency`**: Groups execution by `pr-review-suggestions-${{ github.event.pull_request.number }}` with `cancel-in-progress: true`. If `opened` and `labeled` events trigger in rapid succession, obsolete in-flight runs are cancelled, guaranteeing serial execution and preventing race conditions or duplicate comment creation.

---

## 4. Metadata Inputs

The matching system operates exclusively on public GitHub API metadata:

### Triggering PR Inputs
- **PR Number**: Unique identifier used to filter out self-recommendations and manage comments.
- **PR Title**: Subject line describing the change.
- **PR Description/Body**: Detailed markdown summary provided by the author.
- **PR Labels**: Applied labels categorized into language and area/component sets.

To avoid race conditions where labels are added immediately after opening, the workflow performs an API lookup (`github.rest.pulls.get`) to fetch the most up-to-date PR metadata, falling back to the event payload if necessary.

### Candidate PR Inputs
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
For each candidate, the workflow inspects:
- `number`, `title`, `body`, `html_url`, `labels`, `draft`, and `state`.

---

## 5. Candidate Discovery & Eligibility Filtering

Before scoring begins, candidate PRs must pass six strict filters:

1. **Same Repository**: Only PRs belonging to `meshery/meshery` are queried.
2. **Current PR Exclusion**: The triggering PR is excluded (`candidate.number !== prNumber`).
3. **Open State**: Only pull requests with `state === 'open'` are considered. Closed and merged PRs are discarded.
4. **Draft PR Exclusion**: Draft pull requests (`candidate.draft === true`) are skipped. Contributors should only be directed toward review-ready work.
5. **Deduplication**: Candidate numbers are tracked in a `Set` to prevent duplicate evaluations.
6. **Mandatory Label Overlap Eligibility**: The candidate must share at least one `language/*`, `area/*`, or `component/*` label with the current PR. Candidates with zero label overlap are excluded immediately, preventing text-only false positives.

---

## 6. Scoring Algorithm & Ranking Refinement

Scoring is additive, deterministic, and bounded, where **label overlap guarantees eligibility** and **text matching acts as ranking refinement**:

| Signal | Points | Role | Condition |
| :--- | :--- | :--- | :--- |
| **Language Match** | **+40** | Eligibility & Base | Current PR and candidate share at least one `language/*` label |
| **Area / Component Match** | **+40** | Eligibility & Base | Current PR and candidate share at least one `area/*` or `component/*` label |
| **Title Keyword Overlap** | **+10** | Ranking Refinement | Per unique matching technical keyword shared between titles |
| **Body Keyword Overlap** | **+5** | Ranking Refinement | Per unique matching technical keyword shared between bodies |

### Rules & Normalization
1. **Mandatory Eligibility Gate**: Candidates must score at least `+40` from label overlap (sharing a language, area, or component label). If no label overlap exists, the candidate receives `0` and is discarded.
2. **Flat Category Bonuses**: Multiple matching language labels yield +40 (not +80). Multiple matching area/component labels yield +40. Maximum possible label baseline score is **+80**.
3. **Text Normalization**:
   - Lowercase all characters.
   - Strip punctuation and special characters (`[^a-z0-9_/-]`).
   - Split on whitespace into individual tokens.
   - Filter out tokens with length ≤ 1.
4. **Stopword Filtering**:
   Common non-technical words are ignored to prevent accidental matches:
   `the`, `a`, `an`, `this`, `that`, `is`, `are`, `was`, `were`, `and`, `or`, `to`, `of`, `for`, `in`, `on`, `with`, `from`, `fix`, `add`, `update`, `change`, `changes`, `issue`, `pull`, `request`, `pr`, `improve`, `new`, `by`, `at`, `it`, `as`, `be`, `all`, `not`, `we`, `you`, `i`, `our`, `my`, `do`, `does`, `did`, `but`, `so`, `if`, `then`, `into`, `also`, `has`, `have`, `had`, `been`.
5. **Token Deduplication**: Words are stored in a `Set`. If a word like `kubernetes` appears three times in a title, it contributes +10 exactly once.
6. **Field Separation**: Title tokens are only compared against title tokens; body tokens are only compared against body tokens. No cross-field leakage.

---

## 7. Ranking & Selection

Once scored:
1. **Positive Score Requirement**: Candidates with `score <= 0` are excluded.
2. **Sorting**:
   - **Primary**: Total relevance score descending (`b.score - a.score`).
   - **Secondary (Tie-breaker)**: PR number descending (`b.number - a.number`), favoring more recently created pull requests.
3. **Selection Cap**: Up to three candidates are selected (`scoredCandidates.slice(0, 3)`).
4. **No Artificial Filling**: If only 1 or 2 candidates qualify, only those 1 or 2 are shown. Unrelated PRs are never manufactured to fill a quota.

---

## 8. Sticky Comment Mechanism

The workflow communicates recommendations through a single "sticky" comment on the triggering PR.

### Hidden Marker
All Review Loop comments begin with a hidden HTML comment marker:
```html
<!-- meshery-review-suggestions -->
```

### Idempotent Lifecycle
1. The workflow retrieves existing PR comments using `github.paginate(github.rest.issues.listComments)`.
2. It searches for an existing comment containing `<!-- meshery-review-suggestions -->`.
3. **Update in-place**: If an existing comment is found:
   - If the new content is identical to the existing content (`existingComment.body.trim() === commentBody.trim()`), the workflow logs an informational message and skips updating, avoiding unnecessary comment edit histories.
   - If the content changed (e.g. new labels were applied), it updates the existing comment via `issues.updateComment`.
4. **Create**: If no comment exists, it creates one via `issues.createComment`.
5. Duplicate suggestion comments are **never** created.

### Comment Layouts

#### When Matches Are Found
```markdown
<!-- meshery-review-suggestions -->

### 🔎 PRs you may want to review

Based on the language, area, and context of this PR, here are a few related open PRs that may be useful to review:

- [#1234](https://github.com/meshery/meshery/pull/1234) — [Server] Fix environment pagination
- [#1235](https://github.com/meshery/meshery/pull/1235) — [mesheryctl] Support system context switch

Reviewing related work can help contributors learn the codebase and contribute review capacity back to the community.

_This is only a suggestion and does not block your PR._
```

#### When No Matches Are Found
```markdown
<!-- meshery-review-suggestions -->

### 🔎 PRs you may want to review

Based on the language, area, and context of this PR, no closely related open PRs were found to review at this time.

Reviewing related work can help contributors learn the codebase and contribute review capacity back to the community.

_This is only a suggestion and does not block your PR._
```

---

## 9. Permissions & Security

Because this workflow uses `pull_request_target`, security is enforced through strict least-privilege configurations:

### Minimal Permissions
```yaml
permissions:
  issues: write
  pull-requests: read
```
- `pull-requests: read`: Allows reading triggering PR details and querying open candidate pull requests via `pulls.get` and `pulls.list`.
- `issues: write`: Allows listing, creating, and updating PR comments through GitHub's Issues API (`issues.createComment`, `issues.updateComment`, `issues.listComments`).
- **No `contents: read` / `contents: write`**: The workflow does not read or write repository files.

### Security Guarantees
- **No Code Checkout**: The workflow does **not** use `actions/checkout`. PR branches and forks are never cloned.
- **No Code Execution**: No npm packages from the PR are installed, no scripts are executed, and no tests are run.
- **No Shell Interpolation**: PR titles and descriptions are processed exclusively as JavaScript strings in memory. They are never passed to shell commands or bash scripts.
- **No External Secrets**: Runs strictly on GitHub's built-in `GITHUB_TOKEN`.
- **No Automatic Reviewer Assignment**: Suggestions are purely informational; the workflow never calls reviewer assignment APIs.

---

## 10. Failure Handling

- The entire script body is wrapped in a top-level `try / catch` block:
  ```javascript
  try {
    // Review loop execution
  } catch (error) {
    core.warning(`PR review suggestions failed safely: ${error.message}`);
  }
  ```
- If an API call fails or rate limits are reached, the error is logged as a warning (`core.warning`).
- The step completes without failing the job, ensuring that Review Loop issues **never fail contributor CI checks or block PR submission**.

---

## 11. Current Limitations

1. **Token-based Keyword Matching**: Does not perform semantic stemming or synonym recognition (e.g., `auth` does not automatically match `authentication`).
2. **Label Reliance**: Matching accuracy is highest when PRs are properly labeled with language and component tags.
3. **Recent Candidate Window**: The workflow evaluates the 100 most recently updated open PRs. Open PRs inactive for long periods outside this window are not evaluated.
4. **Heuristic Weights**: Initial weights (+40 for categories, +10 for title tokens, +5 for body tokens) are initial baseline heuristics and may be adjusted based on live repository metrics.

---

## 12. Maintenance & Operations

### Tuning Weights
To modify scoring weights, edit `calculateScore` in `.github/workflows/pr-review-suggestions.yml`:
- Language match bonus (default: `+40`)
- Area/Component match bonus (default: `+40`)
- Title keyword weight (default: `+10`)
- Body keyword weight (default: `+5`)

### Updating Stopwords
To filter out newly observed generic words, add them to `STOP_WORDS` in `.github/workflows/pr-review-suggestions.yml`.

### Verifying Syntax
After modifying the workflow:
```bash
# Validate YAML syntax
python3 -c 'import yaml; yaml.safe_load(open(".github/workflows/pr-review-suggestions.yml"))'

# Validate JavaScript syntax
node -c <(python3 -c 'import yaml; wf = yaml.safe_load(open(".github/workflows/pr-review-suggestions.yml")); print("async () => {\n" + wf["jobs"]["suggest-reviews"]["steps"][0]["with"]["script"] + "\n}")')
```
