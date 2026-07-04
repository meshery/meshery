---
name: iterate-pr
description: Iterate on a PR until CI passes. Use when you need to fix CI failures, address review feedback, or continuously push fixes until all checks are green. Automates the feedback-fix-push-wait cycle.
---

# Iterate on PR Until CI Passes

Continuously iterate on the current branch until all CI checks pass and review feedback is addressed.

**Requires**: GitHub CLI (`gh`) authenticated.

**Requires**: The `uv` CLI for python package management, install guide at https://docs.astral.sh/uv/getting-started/installation/

**Important**: All scripts must be run from the repository root directory (where `.git` is located), not from the skill directory. Use the full path to the script via `${HOME}/.agents/skills/iterate-pr`.

## Bundled Scripts

### `scripts/fetch_pr_checks.py`

Fetches CI check status and extracts failure snippets from logs.

```bash
uv run ${HOME}/.agents/skills/iterate-pr/scripts/fetch_pr_checks.py [--pr NUMBER]
```

Returns JSON:
```json
{
  "pr": {"number": 123, "branch": "feat/foo"},
  "summary": {"total": 5, "passed": 3, "failed": 2, "pending": 0},
  "checks": [
    {"name": "tests", "status": "fail", "log_snippet": "...", "run_id": 123},
    {"name": "lint", "status": "pass"}
  ]
}
```

### `scripts/fetch_pr_feedback.py`

Fetches and categorizes PR review feedback using the [LOGAF scale](https://develop.sentry.dev/engineering-practices/code-review/#logaf-scale).

```bash
uv run ${HOME}/.agents/skills/iterate-pr/scripts/fetch_pr_feedback.py [--pr NUMBER]
```

Returns JSON with feedback categorized as:
- `high` - Must address before merge (`h:`, blocker, changes requested)
- `medium` - Should address (`m:`, standard feedback)
- `low` - Optional (`l:`, nit, style, suggestion)
- `bot` - Informational automated comments (Codecov, Dependabot, etc.)
- `resolved` - Already resolved threads

Review bot feedback (from Sentry, Warden, Cursor, Bugbot, CodeQL, etc.) appears in `high`/`medium`/`low` with `review_bot: true` — it is NOT placed in the `bot` bucket.

Each feedback item may also include:
- `thread_id` - GraphQL node ID for inline review comments (used for replies via `reply_to_thread.py`)

### `scripts/reply_to_thread.py`

Replies to PR review threads. Batches multiple replies into a single GraphQL call.

```bash
uv run ${HOME}/.agents/skills/iterate-pr/scripts/reply_to_thread.py THREAD_ID "body" [THREAD_ID "body" ...]
```

Arguments are alternating `(thread_id, body)` pairs. The script sends the reply body without adding signatures, attribution, or sign-off text. Example:
```bash
uv run ${HOME}/.agents/skills/iterate-pr/scripts/reply_to_thread.py \
  PRRT_abc "Fixed the null check." \
  PRRT_def "Replaced with path-segment counting."
```

## Workflow

### 1. Identify PR

```bash
gh pr view --json number,url,headRefName
```

Stop if no PR exists for the current branch.

### 2. Gather Review Feedback

Run `${HOME}/.agents/skills/iterate-pr/scripts/fetch_pr_feedback.py` to get categorized feedback already posted on the PR.

### 3. Handle Feedback by LOGAF Priority

**Auto-fix (no prompt):**
- `high` - must address (blockers, security, changes requested)
- `medium` - should address (standard feedback)

When fixing feedback:
- Understand the root cause, not just the surface symptom
- Check for similar issues in nearby code or related files
- Fix all instances, not just the one mentioned

This includes review bot feedback (items with `review_bot: true`). Treat it the same as human feedback:
- Real issue found → fix it
- False positive → skip, but explain why in a brief comment
- Never silently ignore review bot feedback — always verify the finding

**Prompt user for selection:**
- `low` - present numbered list and ask which to address:

```
Found 3 low-priority suggestions:
1. [l] "Consider renaming this variable" - @reviewer in api.py:42
2. [nit] "Could use a list comprehension" - @reviewer in utils.py:18
3. [style] "Add a docstring" - @reviewer in models.py:55

Which would you like to address? (e.g., "1,3" or "all" or "none")
```

**Skip silently:**
- `resolved` threads
- `bot` comments (informational only — Codecov, Dependabot, etc.)

#### Replying to Comments

After processing each inline review comment, reply on the PR thread to acknowledge the action taken. Only reply to items with a `thread_id` (inline review comments).

**When to reply:**
- `high` and `medium` items — whether fixed or determined to be false positives
- `low` items — whether fixed or declined by the user

**How to reply:** Use `${HOME}/.agents/skills/iterate-pr/scripts/reply_to_thread.py`. Batch all replies for a round into a single call:

```bash
uv run ${HOME}/.agents/skills/iterate-pr/scripts/reply_to_thread.py \
  PRRT_abc "Fixed — description of change." \
  PRRT_def "Not applicable — reason."
```

**Reply format:**
- 1-2 sentences: what was changed, why it's not an issue, or acknowledgment of declined items
- Never add a signature, attribution line, tag, or vendor/model mention in replies
- Keep replies tool-agnostic and identity-free
- If the script fails, log and continue — do not block the workflow

### 4. Check CI Status

Run `${HOME}/.agents/skills/iterate-pr/scripts/fetch_pr_checks.py` to get structured failure data.

**Wait if pending:** If review bot checks (sentry, warden, cursor, bugbot, seer, codeql) are still running, wait before proceeding—they post actionable feedback that must be evaluated. Informational bots (codecov) are not worth waiting for.

### 5. Fix CI Failures

For each failure in the script output:
1. Read the `log_snippet` and trace backwards from the error to understand WHY it failed — not just what failed
2. Read the relevant code and check for related issues (e.g., if a type error in one call site, check other call sites)
3. Fix the root cause with minimal, targeted changes
4. Find existing tests for the affected code and run them. If the fix introduces behavior not covered by existing tests, extend them to cover it (add a test case, not a whole new test file)

Do NOT assume what failed based on check name alone—always read the logs. Do NOT "quick fix and hope" — understand the failure thoroughly before changing code.

### 6. Verify Locally, Then Commit and Push

Before committing, verify your fixes locally:
- If you fixed a test failure: re-run that specific test locally
- If you fixed a lint/type error: re-run the linter or type checker on affected files
- For any code fix: run existing tests covering the changed code

If local verification fails, fix before proceeding — do not push known-broken code.

```bash
git add <files>
gh auth status -a
git commit --signoff -m "fix: <descriptive message>"
git push
```

Always add exactly one sign-off to each commit for the active authenticated GitHub user. Check `gh auth status -a` before committing to confirm the active account, use `git commit --signoff`, and do not add any other trailers, signatures, or tool attribution to commits or PR comments.

### 7. Monitor CI and Address Feedback

Poll CI status and review feedback in a loop instead of blocking:

1. Run `uv run ${HOME}/.agents/skills/iterate-pr/scripts/fetch_pr_checks.py` to get current CI status
2. If all checks passed → proceed to exit conditions
3. If any checks failed (none pending) → return to step 5
4. If checks are still pending:
   a. Run `uv run ${HOME}/.agents/skills/iterate-pr/scripts/fetch_pr_feedback.py` for new review feedback
   b. Address any new high/medium feedback immediately (same as step 3)
   c. If changes were needed, commit and push (this restarts CI), then continue polling
   d. Sleep 30 seconds (don't increase on subsequent iterations), then repeat from sub-step 1
5. After all checks pass, do a final feedback check: `sleep 10`, then run `uv run ${HOME}/.agents/skills/iterate-pr/scripts/fetch_pr_feedback.py`. Address any new high/medium feedback — if changes are needed, return to step 6.

### 8. Repeat

If step 7 required code changes (from new feedback after CI passed), return to step 2 for a fresh cycle. CI failures during monitoring are already handled within step 7's polling loop.

## Exit Conditions

**Success:** All checks pass, post-CI feedback re-check is clean (no new unaddressed high/medium feedback including review bot findings), user has decided on low-priority items.

**Ask for help:** Same failure after 2 attempts, feedback needs clarification, infrastructure issues.

**Stop:** No PR exists, branch needs rebase.

## Fallback

If scripts fail, use `gh` CLI directly:
- `gh pr checks name,state,bucket,link`
- `gh run view <run-id> --log-failed`
- `gh api repos/{owner}/{repo}/pulls/{number}/comments`
