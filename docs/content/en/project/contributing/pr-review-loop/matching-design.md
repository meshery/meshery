---
title: PR Review Loop — Matching Design
description: Architectural design, label eligibility rules, text normalization, and scoring model for PR review recommendations.
categories: [contributing]
---

# Meshery PR Review Loop — Matching Design

## 1. Goal

When a contributor opens a pull request (PR) in the `meshery/meshery` repository, which existing open PRs are most relevant for them to review?

The purpose of the Meshery PR Review Loop is to encourage reciprocal peer review among contributors working on similar areas of the codebase. By suggesting relevant open PRs to an author at the moment they submit code, contributors can share domain knowledge, review each other's work, and help maintainers move changes forward faster.

Version 1 (V1) must be a simple, deterministic, explainable, no-AI matching system.

## 2. Inputs

The matching algorithm operates entirely on public pull request metadata available through the GitHub API.

For the newly opened or updated PR:
1. PR labels
2. PR title
3. PR description/body

For candidate PRs:
1. PR labels
2. PR title
3. PR description/body

## 3. Label Matching & Mandatory Eligibility Gate

Labels serve as the primary relevance signal in the matching design.

### Label Categories

- **Language labels**:
  - `language/go`
  - `language/javascript`
  - `language/typescript`
  - `language/python`
  - `language/ruby`
  - `language/shell`
- **Area and Component labels**:
  - Existing `area/*` labels (e.g., `area/docs`, `area/ci`, `area/models`, `area/lifecycle`, `area/tests`)
  - Existing `component/*` labels (e.g., `component/ui`, `component/server`, `component/mesheryctl`, `component/extensions`, `component/database`, `component/policy`)

### Label Matching Rules (Eligibility Gate)

- **Mandatory Eligibility Requirement**: A candidate PR is eligible for consideration **ONLY** when it shares at least one label with the current PR in the following categories:
  - `language/*` label
  - `area/*` label
  - `component/*` label
- **Text-only Exclusion**: A candidate with matching title or body keywords but **no** overlapping language, area, or component label is **strictly excluded**.
- Matching language indicates that both contributors are working in the same programming environment and runtime ecosystem.
- Matching area or component indicates that both contributors are touching related architectural subsystems of Meshery.
- Once eligibility is established via label overlap, title and body keyword matching serve as **ranking refinement** to order candidates.

## 4. Text Matching & Normalization

Deterministic local text matching is used to compare the textual contents of the PR title and description. It does not use any LLM, AI agent, vector embeddings, or external AI services.

### Normalization Process

Before comparison, the title and description of both the new PR and candidate PRs are normalized:
1. **Lowercase**: Convert all characters to lowercase.
2. **Punctuation conversion**: Convert all punctuation—including hyphens, slashes, underscores, and brackets—to spaces before tokenization (regex `/[^a-z0-9]/g`). This ensures compound phrases like `server-side` or `component/server` produce separate tokens (`server`, `side`, `component`) and share expected token overlap with individual words.
3. **Tokenization**: Split text on whitespace into individual tokens.
4. **Token Filtering**: Discard tokens of length ≤ 1.
5. **Stopword filtering**: Ignore generic, low-value words that do not indicate technical or domain relevance.
6. **Set Deduplication**: Tokens are stored in a `Set` to prevent repeated words from inflating scores.

### Generic Words to Ignore (Stopwords)

The following low-value words are ignored during text matching:
- `the`, `a`, `an`, `this`, `that`, `is`, `are`, `was`, `were`
- `and`, `or`, `to`, `of`, `for`, `in`, `on`, `with`, `from`
- `fix`, `add`, `update`, `change`, `changes`, `issue`, `pull`, `request`, `pr`, `improve`, `new`
- `by`, `at`, `it`, `as`, `be`, `all`, `not`, `we`, `you`, `i`, `our`, `my`
- `do`, `does`, `did`, `but`, `so`, `if`, `then`, `into`, `also`, `has`, `have`, `had`, `been`

### Technical and Domain-Specific Words

Text matching specifically prioritizes technical and domain-specific words, such as:
- `environment`, `pagination`, `kubernetes`, `graphql`, `database`, `authentication`, `server`, `mesheryctl`, `relationship`
- API-related terms (e.g., `endpoint`, `handler`, `schema`, `route`)
- Framework and component names (e.g., `helm`, `operator`, `meshsync`, `sistent`, `relay`)

## 5. Scoring Model

The scoring model calculates relevance using the following point values:

| Criterion | Points | Details |
|---|---|---|
| Same language | +40 | Both PRs share at least one matching `language/*` label |
| Same area/component | +40 | Both PRs share at least one matching `area/*` or `component/*` label |
| Technical keyword overlap in PR title | +10 per word | Each matching technical keyword shared between titles |
| Technical keyword overlap in PR description/body | +5 per word | Each matching technical keyword shared between descriptions |

## 6. Candidate Filtering

Before scoring, candidate PRs must pass strict eligibility and filtering criteria:
- **Same repository**: Must belong to `meshery/meshery`.
- **Target branch**: Must target the `master` branch.
- **Open state**: Candidate PR must currently be open (`state === 'open'`). Closed and merged PRs are discarded.
- **Not current PR**: The PR triggering the suggestion is excluded from the candidate pool (`candidate.number !== prNumber`).
- **No duplicates**: Candidate PR numbers are tracked in a `Set` to prevent duplicate evaluations.
- **Draft PR exclusion**: Draft PRs (`draft === true`) are excluded from recommendations, ensuring contributors are only directed toward review-ready pull requests.
- **Mandatory label overlap**: The candidate must share at least one `language/*`, `area/*`, or `component/*` label with the triggering PR. Candidates with zero label overlap are discarded immediately before ranking.

## 7. Ranking

After candidates are filtered and scored:
1. **Sort candidates**: Order candidates primarily by descending relevance score (`b.score - a.score`).
2. **Tie-breaking**: When scores are identical, order by PR number descending (`b.number - a.number`), favoring newer pull requests.
3. **Select top candidates**: Select up to three candidates (`slice(0, 3)`).
4. **Do not fabricate**: Do not fabricate or invent candidates to reach three. If only one or two PRs qualify, only return those.
5. **No irrelevant recommendations**: If no candidate shares label overlap or qualifies, the comment indicates that no closely related open PRs were found.

## 8. Example

### Example 1: Strong Match (Language + Text Overlap)

- **New PR**:
  - Labels: `language/go`, `component/server`
  - Title: `[Server] Fix environment pagination`
  - Description: `Fix pagination when more than 10 environments are returned. Updates environment API handling and tests.`

- **Candidate PR A**:
  - Labels: `language/go`, `component/mesheryctl`
  - Title: `[mesheryctl] Fix environment view pagination`
  - Description: `Fix environment pagination and v1beta3 handling.`

**Eligibility**: Shares `language/go` → **Eligible**
**Score Breakdown**:
- Same language (`language/go`): +40
- Title technical keyword overlap (`environment`, `pagination`): 2 words × +10 = +20
- Description technical keyword overlap (`environment`, `pagination`): 2 words × +5 = +10
- **Total Score**: 70 points

### Example 2: Strong Label Match with Weak Text Overlap

- **Candidate PR B**:
  - Labels: `language/go`, `component/server`
  - Title: `[Server] Add telemetry heartbeat endpoint`
  - Description: `Implements ping interval monitoring for connected agents.`

**Eligibility**: Shares `language/go` and `component/server` → **Eligible**
**Score Breakdown**:
- Same language (`language/go`): +40
- Same area/component (`component/server`): +40
- Title technical keyword overlap: 0 words = 0
- Description technical keyword overlap: 0 words = 0
- **Total Score**: 80 points

### Example 3: Text Overlap without Label Overlap (Excluded)

- **Candidate PR C**:
  - Labels: `area/docs`
  - Title: `[Docs] Fix environment pagination guide`
  - Description: `Updates documentation covering environment pagination.`

**Eligibility**: Shares NO `language/*`, `area/*`, or `component/*` labels with the new PR (`language/go`, `component/server`).
**Result**: **EXCLUDED**. Even though titles share `environment` and `pagination`, lack of label overlap disqualifies Candidate PR C from eligibility.

## 9. Comment Integration & Security Hardening

The workflow communicates recommendations through a single sticky PR comment:

- **Marker**: The comment begins with a hidden HTML comment marker:
  ```html
  <!-- meshery-review-suggestions -->
  ```
- **Bot Ownership Verification**: To prevent comment hijacking, the workflow verifies that existing sticky comments were authored by `github-actions[bot]`. Comments by human participants or unrelated bots containing the marker are never modified.
- **Candidate Title Escaping**: Candidate titles are untrusted user input. All titles are escaped before Markdown rendering to prevent Markdown link injection (`[text](url)`), code formatting spans, HTML tags, or table pipe breaks.
- **Update in-place**: If a bot-authored sticky comment exists:
  - If the content is identical, the update is skipped (no redundant API calls).
  - If labels changed or candidates changed, the comment is updated in place.
  - If all eligible candidates are removed, the comment updates to notify the contributor that no closely related open PRs were found, preventing stale recommendations.
