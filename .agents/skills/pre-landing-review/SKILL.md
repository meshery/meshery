---
name: pre-landing-review
description: 'Gate-oriented safety audit for a Meshery change before it lands, using a checklist with two-pass severity triage (blocking CRITICAL, advisory INFORMATIONAL). Triggers on: "is this safe to land", "pre-landing review", "safety check before merge", "gate check", "/pre-landing-review". NOT a general diff review - use /code-review for that.'
metadata:
  version: 1.0.0-meshery
  upstream: https://github.com/Mathews-Tom/armory/tree/main/skills/pre-landing-review
  upstream-version: 1.0.1
  license: MIT (see .agents/skills/ATTRIBUTIONS.md)
---

# Pre-Landing Review

A safety gate run before a change lands. It checks the change against
`.agents/skills/pre-landing-review/references/checklist.md` in two passes and
returns **CLEAR TO LAND** or **BLOCKED**.

It is not a substitute for `/code-review` (correctness bugs) or
`/security-review`. It is a short, opinionated gate that encodes the failure
classes Meshery has actually shipped: SQL injection through gorm `.Order()`,
nil providers at boot, ungated pages, wire-casing drift from schemas, and
unregenerated MeshKit error codes.

## Workflow

### 1. Determine the diff

- On a feature branch: `git diff "$(git merge-base HEAD origin/master)"...HEAD`
- Given a PR number: `gh-axi pr diff <number>`
- Given specific files: review those files

### 2. Load the checklist

Read `.agents/skills/pre-landing-review/references/checklist.md`. This is
mandatory - if it is unreadable, STOP and report the error.

### 3. Pass 1 - CRITICAL (blocking)

Check every changed hunk against the Pass 1 sections of the checklist:
SQL & data safety, concurrency, trust boundaries, and Meshery contracts.
Where the checklist cites an `AGENTS.md` rule, the rule is the authority -
re-read it rather than paraphrasing from memory.

For each CRITICAL finding:

1. Cite the exact `file:line`.
2. Explain the concrete failure (inputs/state -> wrong outcome), not the
   category name.
3. Ask with `AskUserQuestion`, exactly three options:
   **Fix now** / **Acknowledge risk** / **False positive**.
4. **Fix now**: make the fix, then re-check that hunk.
5. **Acknowledge** or **False positive**: record it and continue.

### 4. Pass 2 - INFORMATIONAL (non-blocking)

Check the remaining checklist sections: hidden side effects, magic numbers,
dead code, test gaps, crypto, time handling, type coercion. Present every
finding in one table (`file:line`, category, description).

### 5. Mechanical checks

Run the gates that apply to the touched components and report their real
output. Do not claim a pass you did not observe:

| Touched | Run |
|---|---|
| `server/` | `make golangci` (includes the `orderby` analyzer) |
| `ui/` | `make ui-lint` |
| MeshKit errors added | `make error`, then confirm `next_error_code` bumped and `docs/data/errorref/` regenerated |
| API / schema consumers | `cd ../schemas && make validate-schemas && make consumer-audit` |

A failing gate is a Pass 1 finding.

### 6. Gate classification

- All Pass 1 findings fixed or acknowledged, mechanical checks green -> **CLEAR TO LAND**
- Any unresolved Pass 1 finding or red mechanical check -> **BLOCKED**
- Pass 2 findings never block.

If the change falls under AGENTS.md "Require Human Review" (auth/secrets,
database migrations, breaking API changes, Helm templates, CI workflows), say
so in the verdict even when it is CLEAR TO LAND. For migrations, also run
`/migration-risk-analyzer`.

### 7. Suppressions

Do NOT flag:

- Test files using fixtures or factories
- `//nolint:orderby` with a justification comment (see the lint docs)
- Code whose inline comment explains why a flagged pattern is intentional
- Generated files (schemas output, `docs/data/errorref/`, GraphQL generated code)

## Output

Gate verdict (CLEAR TO LAND / BLOCKED), the CRITICAL findings with their
resolution, the mechanical-check results, and the INFORMATIONAL table.

**Read-only by default** - modifies code only when the user picks "Fix now".
