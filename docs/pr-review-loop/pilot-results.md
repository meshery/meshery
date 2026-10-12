# Meshery PR Review Loop — Pilot Results

**Status: Not yet run**

> **Note**: No live production pilot has been conducted yet. This document establishes the tracking framework, target evaluation metrics, and validation template for a future live pilot. No metrics or review conversion numbers are fabricated or assumed.

---

## 1. Status Overview

| Stage | Status | Notes |
|---|---|---|
| Implementation Validation | **Completed** | Algorithmic logic, edge-case tests, YAML validation, and security reviews verified locally. |
| Staging / Dry-run Evaluation | **Pending** | Awaiting maintainer review and pilot deployment schedule. |
| Production Pilot Run | **Not yet run** | Live data collection will commence once enabled on `master`. |

---

## 2. Metrics to Measure

When an approved pilot is launched, the following metrics will be tracked:

1. **PRs Processed**: Total number of newly opened or labeled pull requests evaluated by the workflow.
2. **PRs Receiving Suggestions**: Count and percentage of PRs where at least one candidate PR was suggested.
3. **Suggestions Displayed**: Total number of suggested PR links generated across all PR comments.
4. **Subsequent Review / Comment Conversion**: Number of suggested PRs that subsequently receive a code review, approval, or substantive comment from the author who received the recommendation.
5. **False-Positive Feedback**: Instances where suggested PRs were reported or identified as technically irrelevant despite label/keyword matches.
6. **False-Negative Feedback**: Instances where highly relevant open PRs existed in the candidate pool but were missed by the scoring algorithm.
7. **Maintainer & Contributor Feedback**: Qualitative feedback regarding comment readability, notification frequency, and review discovery value.

---

## 3. Pilot Tracking Template

The following template will be populated during the pilot evaluation period:

| Metric | Target / Benchmark | Observed Result | Notes |
|---|---|---|---|
| Total PRs evaluated | Baseline cohort | TBD | Pending live run |
| PRs with ≥1 suggestion | > 60% | TBD | Pending live run |
| Average suggestions per PR | 1–3 | TBD | Pending live run |
| Review conversions initiated | Positive baseline | TBD | Pending live run |
| Comment update / edit rate | < 2 updates per PR | TBD | Pending live run |
| False-positive report rate | < 10% | TBD | Pending live run |
| Workflow execution failures | 0% | TBD | Pending live run |

---

## 4. Completed Implementation Validation

While real-world pilot metrics are pending live deployment, the underlying implementation has passed the following static and synthetic verifications:

- **Syntax & Schema**: Verified valid YAML schema and node/javascript execution.
- **Label Eligibility Enforcement**: Verified that candidates lacking `language/*`, `area/*`, or `component/*` overlap are strictly excluded regardless of text overlap.
- **Candidate Filtering**: Verified exclusion of self-PRs, closed PRs, and draft PRs.
- **Deduplication & Tie-breaking**: Verified ranking order (score descending, PR number descending) and slicing to top 3.
- **Concurrency & Idempotency**: Concurrency group protection and sticky HTML comment marker (`<!-- meshery-review-suggestions -->`) logic validated to prevent duplicate comments.
- **Security Scoping**: Verified no repository code execution, no checkout of untrusted PR branches, and strict least-privilege token permissions (`issues: write`, `pull-requests: read`).
