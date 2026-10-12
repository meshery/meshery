# Meshery PR Review Loop — Matching Design

## 1. Goal

When a contributor opens a pull request (PR) in the `meshery/meshery` repository, which existing open PRs are most relevant for them to review?

The purpose of the Meshery PR Review Loop is to encourage reciprocal peer review among contributors working on similar areas of the codebase. By suggesting relevant open PRs to an author at the moment they submit code, contributors can share domain knowledge, review each other's work, and help maintainers move changes forward faster.

Version 1 (V1) must be a simple, deterministic, explainable, no-AI matching system.

## 2. Inputs

The matching algorithm operates entirely on public pull request metadata available through the GitHub API.

For the newly opened PR:
1. PR labels
2. PR title
3. PR description/body

For candidate PRs:
1. PR labels
2. PR title
3. PR description/body

## 3. Label Matching

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

## 4. Text Matching

Deterministic local text matching is used to compare the textual contents of the PR title and description. It does not use any LLM, AI agent, vector embeddings, or external AI services.

### Normalization Process

Before comparison, the title and description of both the new PR and candidate PRs are normalized:
1. **Lowercase**: Convert all characters to lowercase.
2. **Punctuation removal**: Strip punctuation, special characters, and formatting symbols.
3. **Tokenization**: Split text into individual words/tokens.
4. **Stopword filtering**: Ignore generic, low-value words that do not indicate technical or domain relevance.

### Generic Words to Ignore (Stopwords)

The following low-value words are ignored during text matching:
- `the`
- `this`
- `that`
- `fix`
- `add`
- `update`
- `change`
- `changes`
- `issue`
- `pull`
- `request`
- `pr`
- `improve`
- `new`

The V1 implementation may include additional generic stopwords as needed to keep token comparisons focused.

### Technical and Domain-Specific Words

Text matching specifically prioritizes technical and domain-specific words, such as:
- `environment`
- `pagination`
- `kubernetes`
- `graphql`
- `database`
- `authentication`
- `server`
- `mesheryctl`
- `relationship`
- API-related terms (e.g., `endpoint`, `handler`, `schema`, `route`)
- Framework and component names (e.g., `helm`, `operator`, `meshsync`, `sistent`, `relay`)

## 5. Scoring Model

The proposed initial scoring model for V1 calculates relevance using the following point values:

| Criterion | Points | Details |
|---|---|---|
| Same language | +40 | Both PRs share at least one matching language label |
| Same area/component | +40 | Both PRs share at least one matching `area/*` or `component/*` label |
| Technical keyword overlap in PR title | +10 per word | Each matching technical keyword found in both PR titles |
| Technical keyword overlap in PR description/body | +5 per word | Each matching technical keyword found in both PR descriptions |

These weights are initial V1 configuration values and can be adjusted after testing and observing recommendations on real pull requests.

## 6. Candidate Filtering

Before scoring, candidate PRs must pass strict eligibility and filtering criteria:
- **Same repository**: Must belong to `meshery/meshery`.
- **Open state**: Candidate PR must currently be open.
- **Not current PR**: The PR triggering the suggestion is excluded from the candidate pool.
- **No duplicates**: Candidate PRs cannot appear multiple times in the candidate list.
- **Draft PR exclusion**: Draft PRs are excluded from candidate recommendations, ensuring contributors are only directed toward review-ready pull requests.
- **Mandatory label overlap**: The candidate must share at least one `language/*`, `area/*`, or `component/*` label with the triggering PR. Candidates with zero label overlap are discarded immediately before ranking.

## 7. Ranking

After candidates are filtered and scored:
1. **Sort candidates**: Order candidates by descending relevance score.
2. **Tie-breaking**: When scores are identical, order by PR number descending (favoring newer pull requests).
3. **Select top candidates**: Select up to three candidates.
4. **Do not fabricate**: Do not fabricate or invent candidates to reach three. If only one or two PRs qualify, only return those.
5. **No irrelevant recommendations**: If no candidate shares label overlap or qualifies, no recommendations are made.
6. **Label primacy**: Label overlap is the mandatory entry gate, while technical keywords provide ranking refinement among eligible candidates.

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

**Explanation**: Candidate PR A is eligible due to sharing `language/go`, and receives ranking refinement from strong title and description keyword overlap.

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

**Explanation**: Candidate PR B shares both the programming language and the specific server component with the new PR, earning an initial baseline score of 80 points. Even though there is no title or description keyword overlap, the candidate remains highly relevant due to strong label signals.

### Example 3: Text Overlap without Label Overlap (Excluded)

- **Candidate PR C**:
  - Labels: `area/docs`
  - Title: `[Docs] Fix environment pagination guide`
  - Description: `Updates documentation covering environment pagination.`

**Eligibility**: Shares NO `language/*`, `area/*`, or `component/*` labels with the new PR (`language/go`, `component/server`).
**Result**: **EXCLUDED**. Even though titles share `environment` and `pagination`, lack of label overlap disqualifies Candidate PR C from eligibility.

## 9. Comment Integration & Concurrency Safety

The GitHub Actions workflow presents the top recommendations to the contributor via a single sticky PR comment:

- **Marker**: The workflow creates or updates ONE comment using this hidden marker:
  ```html
  <!-- meshery-review-suggestions -->
  ```
- **Update in-place**: If a comment containing the hidden marker already exists on the PR, the workflow updates that comment rather than creating duplicate comments.
- **Concurrency Protection**: Workflow runs are protected by per-PR concurrency grouping:
  ```yaml
  concurrency:
    group: pr-review-suggestions-${{ github.event.pull_request.number }}
    cancel-in-progress: true
  ```
  This cancels obsolete in-flight runs when rapid events (`opened`, `labeled`) trigger in close succession, preventing race conditions and duplicate comment creation.

## 10. V1 Limitations

- **Keyword matching is not semantic understanding**: The matching logic looks for identical tokens and cannot infer semantic relationships (e.g., recognizing that `auth` relates to `authentication`).
- **Label dependency**: Inconsistent or missing PR labels can reduce suggestion relevance.
- **Generic titles**: Very generic PR titles or brief descriptions produce weak textual signals.
- **Score tuning**: Initial scoring weights are heuristic estimates and require future tuning using real results.
- **Suggestions, not reviewer assignment**: The system is an advisory suggestion mechanism, not a formal reviewer assignment.
- **Non-blocking**: Review suggestions must never block PR submission, checks, or merging.

## 11. Out of Scope

V1 explicitly does NOT include:
- AI/LLM matching
- Automatic reviewer assignment
- Merge blocking
- Mandatory review contribution
- Contributor leaderboard/ranking
- Slack notifications
- Dashboard/analytics implementation

## 12. Future Tuning

Following initial V1 testing and feedback:
- Scoring weights (+40, +10, +5) can be calibrated against contributor feedback and review uptake.
- The stopword list and domain dictionary can be refined based on common terminology in `meshery/meshery`.
- Minimum score thresholds can be evaluated to ensure only high-confidence suggestions are posted.
- Component-specific weighting or required language match options can be added as configuration options if needed.
