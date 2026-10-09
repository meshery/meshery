---
name: iterate-pr
description: Iterate on a PR until CI passes. Optionally merge or merge and publish a release of the repo. Automates the feedback-fix-push-wait cycle.
argument-hint: "[--merge] [--release]"
metadata:
  author: leecalcote
  version: "2.4.0"
---

# Iterate on PR Until CI Passes

> **2.4.0** - `fetch_pr_feedback.py` no longer reports an empty result as a successful
> fetch. Every path where a feedback channel could come back empty because the request
> did not actually deliver now fails loudly, and a server-side count probe backstops any
> path not yet known. See [Fetch Integrity](#fetch-integrity).

Continuously iterate on the current branch until all CI checks pass and review feedback is addressed.

**Requires**: GitHub CLI (`gh`) authenticated.

**Requires**: Python 3.9+.

**Important**: Run scripts from the repository root directory (where `.git` is located).

### Running the Bundled Scripts

Skills live in `.agents/skills/` - that is the canonical path, and the only one used in
this document. Some repos also expose `.claude/skills` as a symlink to it, but that
symlink is not guaranteed: it does not survive a Windows checkout with
`core.symlinks=false`, and not every repo carries it. Always invoke via `.agents/`.

Every script below accepts either runner - they only use the standard library, so no
dependency installation is needed either way:

- **Preferred**: [`uv`](https://docs.astral.sh/uv/getting-started/installation/) - `uv run <script>`. Faster startup, isolated from system Python.
- **Fallback**: `python3 <script>` - used automatically when `uv` is not on `PATH`.

Before running any script command shown in this skill, check for `uv` first and prefer it;
fall back to `python3` only if `uv` is unavailable:

```bash
if command -v uv >/dev/null 2>&1; then
  uv run .agents/skills/iterate-pr/scripts/<script>.py [args]
else
  python3 .agents/skills/iterate-pr/scripts/<script>.py [args]
fi
```

The rest of this document shows invocations in the shorter `python3 <script>` form for
readability - substitute `uv run` per the rule above whenever `uv` is available.

### Reading Script Output - Mandatory Gate

`fetch_pr_feedback.py` and `fetch_pr_checks.py` both emit a single JSON object with a
top-level `status` field, and both exit non-zero on failure.

| `status` | Meaning | What you must do |
|---|---|---|
| `"ok"` | The fetch completed. `summary` and `feedback`/`checks` are trustworthy. | Proceed. |
| `"error"` | The fetch did **not** complete. `summary` and `feedback`/`checks` are `null`. | **Stop.** Surface `error` to the user. Never merge. |

**Check `status` before reading any other field.** A failed lookup is not a clean PR.
Do not pipe these scripts through anything that discards the exit status, and never infer
"no feedback" or "checks passed" from an absent count - on failure the counts are `null`,
not `0`, precisely so that the two cases cannot be confused.

### Fetch Integrity

The failure this skill guards against hardest is the *silent* one: reporting a clean PR
because the feedback never arrived, rather than because there was none. A run that merges
on that basis merges unreviewed code, and nothing in the output says so.

`fetch_pr_feedback.py` reads three independent channels - inline review threads (GraphQL),
PR conversation comments (REST), and review bodies (REST). Each is treated as a channel
that must *deliver*, not merely return:

- `gh` exiting non-zero is an error, never an empty channel.
- `gh` exiting 0 while writing nothing is an error too. An undelivered response and an
  empty one are indistinguishable downstream, so they are not allowed to look alike here.
- A payload of unexpected shape is an error, not an absence.
- **Count probe backstop:** if any channel does come back empty, the script asks GitHub how
  many records that channel actually holds. If the server says there are records and the
  fetch produced none, the script fails with the channel named and both numbers shown.

The probe costs one small extra query and only runs when a channel is already empty, so a
genuinely quiet PR reports `status: "ok"` with zero counts and no extra cost on the common
path. A PR with real feedback can no longer report zero.

## Invocation and Modes

Invoke as:

- `/iterate-pr`
- `/iterate-pr --merge` (alias: `full`)
- `/iterate-pr --release`
- `/iterate-pr --merge --release` (alias: `full-release`)

In Claude Desktop, arguments are **hints**, not strict CLI parsing. Treat whatever follows `/iterate-pr` as a mode hint string.

`full` and `full-release` are the positional spellings this skill used before 2.3.0 and are
still what several callers' standing instructions say. They are exact synonyms for `--merge`
and `--merge --release`. Both spellings are supported and must stay supported: a caller who
types `full` is asking for an autonomous merge run, and silently giving them default mode
means the PR they expected to be merged just sits there.

### Argument Hint Interpretation (Claude Desktop)

Use this deterministic precedence:

1. If hints include `--release`, `full-release`, or terms like `release`, `publish`, `ship` → run **release mode**.
2. Else if hints include `--merge`, `full`, or terms like `autonomous`, `merge` → run **merge mode**.
3. Else run **default mode**.

Do not fail because hints are missing or unrecognized; default safely.

| Mode | Behavior |
|---|---|
| Default (`/iterate-pr`) | Iterates on CI + high/medium feedback, asks user about low-priority items, then exits without merging. |
| `--merge` / `full` | Fully autonomous: handles every new feedback item and replies to each one once, re-requests review (Gemini or Copilot) after each push, iterates until no new feedback and CI is green, then administratively merges the PR without waiting for a required-approval status. |
| `--release` / `full-release` | Does everything in `--merge`, then cuts/publishes a release for the repository. |

## Bundled Scripts

### `scripts/fetch_pr_checks.py`

Fetches CI check status and extracts failure snippets from logs.

```bash
python3 .agents/skills/iterate-pr/scripts/fetch_pr_checks.py [--pr NUMBER]
```

Returns JSON:
```json
{
  "status": "ok",
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
python3 .agents/skills/iterate-pr/scripts/fetch_pr_feedback.py [--pr NUMBER]
```

Sources covered, all paginated: review bodies (**every** review state, not only
`CHANGES_REQUESTED`), inline review threads, and PR conversation comments. Each is a
channel that must deliver - see [Fetch Integrity](#fetch-integrity) - so a zero here means
the PR really has nothing, never that a request quietly returned nothing.

Returns JSON with feedback categorized as:
- `high` - Must address before merge (`h:`, blocker, changes requested)
- `medium` - Should address (`m:`, standard feedback)
- `low` - Optional (`l:`, nit, style, suggestion)
- `bot` - Informational automated comments (Codecov, Dependabot, etc.)
- `resolved` - Already resolved threads

Review bot feedback (CodeRabbit, Gemini Code Assist, Copilot, Sentry, Warden, Cursor,
Bugbot, CodeQL, etc.) appears in `high`/`medium`/`low` with `review_bot: true` — it is NOT
placed in the `bot` bucket. `REVIEW_BOT_PATTERNS` in `scripts/fetch_pr_feedback.py` owns
the full roster; do not maintain a second copy of it here. Logins are matched with any
trailing `[bot]` stripped, because REST reports `coderabbitai[bot]` where GraphQL reports
`coderabbitai` for the same account; without that normalization the generic `[bot]` suffix
rule files the fleet's main reviewer under `bot` and it gets skipped silently.

Each feedback item may also include:
- `thread_id` - GraphQL node ID for inline review comments (used for replies via `reply_to_thread.py`)
- `replied` - the newest comment in that thread is yours, so you have already answered it

Top-level fields:
- `status` - see the mandatory gate above
- `viewer` - the login `gh` is authenticated as
- `summary.needs_attention` - open `high` + `medium` items, **excluding** those flagged `replied`
- `summary.already_replied` - how many priority items were excluded on that basis

**Self-authored feedback is never reported.** Comments and reviews written by the PR author
or by `viewer` are dropped, and a thread whose newest comment is yours is flagged `replied`
and leaves `needs_attention`. Without both rules a `--merge` run answers an item, sees its
own answer as new feedback, and never converges.

`replied` is decided per source, by the strongest evidence each one offers:

| Source | `replied` when |
|---|---|
| Inline review thread | the newest comment in that thread is yours - exact |
| Review body | never - a review body has no thread and no resolve button, so track in-session which you have answered |
| PR conversation comment | never - GitHub gives these no resolution state, so track in-session which you have answered |

Only the inline-thread signal is exact evidence, so it is the only one used. Anything
weaker - "the review predates something else we said" - cannot tell answering a review
apart from merely typing after it, and would silently dismiss a review that arrives
mid-round. Re-reporting a review you already answered is visible to you; dropping one you
never read is not.

### `scripts/reply_to_thread.py`

Replies to PR review threads. Batches multiple replies into a single GraphQL call.

```bash
python3 .agents/skills/iterate-pr/scripts/reply_to_thread.py THREAD_ID "body" [THREAD_ID "body" ...]
```

Arguments are alternating `(thread_id, body)` pairs. The script sends the reply body without adding signatures, attribution, or sign-off text. Example:
```bash
python3 .agents/skills/iterate-pr/scripts/reply_to_thread.py \
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

Run `python3 .agents/skills/iterate-pr/scripts/fetch_pr_feedback.py` to get categorized feedback already posted on the PR.

Check `status` first. If it is not `"ok"`, stop and report the `error` — do not continue as
though the PR had no feedback.

### 3. Handle Feedback by Priority and Mode

Determine mode from invocation (`/iterate-pr`, `/iterate-pr --merge`, `/iterate-pr --release`).

#### Default mode (`/iterate-pr`)

**Auto-fix (no prompt):**
- `high` - must address (blockers, security, changes requested)
- `medium` - should address (standard feedback)

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
- items flagged `replied` (you already answered them; act only when a reviewer follows up)
- `bot` comments (informational only — Codecov, Dependabot, etc.)

#### Merge modes (`/iterate-pr --merge`, `/iterate-pr --merge --release`)

Operate autonomously. Process every **new** feedback item returned by `fetch_pr_feedback.py` (`high`, `medium`, `low`, and `bot`).

An item is **new** when it is not `resolved`, not flagged `replied`, and not one you already
answered in an earlier round of this same session. Only inline review threads carry `replied`;
review bodies and PR conversation comments have no resolution state, so the fetcher re-returns
them verbatim on every poll and cannot make that last distinction for you. Keep your own
in-session record of which ones you have replied to and do not answer them twice.

For each item:
- Decide whether the feedback is justified
- If justified, implement the change
- If not justified, reject it with a concise technical reason
- Never leave an item without an explicit decision

When fixing feedback (all modes):
- Understand the root cause, not just the surface symptom
- Check for similar issues in nearby code or related files
- Fix all instances, not just the one mentioned

This includes review bot feedback (items with `review_bot: true`). Treat it the same as human feedback:
- Real issue found → fix it
- False positive → reject with explanation
- Never silently ignore review bot feedback — always verify the finding

#### Replying to Comments

After processing feedback, reply to PR comments/threads to acknowledge the action taken.

**Scope by mode:**
- Default mode: reply to `high`/`medium`; reply to `low` only when fixed or declined by the user
- Merge modes: reply to every new feedback item, including informational bot feedback - one reply per item per session, not one per poll

**How to reply:**
- If `thread_id` exists (inline review thread), use `python3 .agents/skills/iterate-pr/scripts/reply_to_thread.py`
- If no `thread_id` exists, post a PR comment with `gh pr comment <PR_NUMBER> --body "..."`
- In merge modes, a round is incomplete until every **new** item has a corresponding reply; an item you answered in an earlier round of this session already has one, so leave it alone rather than replying again

Batch inline replies for a round into a single call:

```bash
python3 .agents/skills/iterate-pr/scripts/reply_to_thread.py \
  PRRT_abc "Fixed — description of change." \
  PRRT_def "Not applicable — reason."
```

**Reply format:**
- 1-2 sentences: what was changed, why it's not an issue, or acknowledgment of declined items
- Never add a signature, attribution line, tag, or vendor/model mention in replies
- Keep replies tool-agnostic and identity-free
- If the script fails, log and continue — do not block the workflow

### 4. Check CI Status

Run `python3 .agents/skills/iterate-pr/scripts/fetch_pr_checks.py` to get structured failure data.
Check `status` first; if it is not `"ok"`, stop and report the `error` rather than treating
the PR as green.

**Wait if pending:** If review bot checks (sentry, warden, cursor, bugbot, seer, codeql, coderabbit, gemini) are still running, wait before proceeding—they post actionable feedback that must be evaluated. Informational bots (codecov) are not worth waiting for.

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

1. Run `python3 .agents/skills/iterate-pr/scripts/fetch_pr_checks.py` to get current CI status
2. If all checks passed → proceed to exit conditions
3. If any checks failed (none pending) → return to step 5
4. If checks are still pending:
   a. Run `python3 .agents/skills/iterate-pr/scripts/fetch_pr_feedback.py` for new review feedback
   b. Address feedback based on mode:
      - Default mode: new `high`/`medium`
      - Merge modes: every new item (`high`/`medium`/`low`/`bot`) and reply to each item - "new" as defined in step 3, so a review body or PR comment you already answered this session is skipped, not re-answered
   c. If changes were needed, commit and push (this restarts CI)
   d. In merge modes, after each push, explicitly re-request review:
      - Prefer Gemini review command (`/gemini review`) when available
      - Otherwise request Copilot review by commenting `@copilot review` on the PR
   e. Sleep 30 seconds (don't increase on subsequent iterations), then repeat from sub-step 1
5. After all checks pass, do a final feedback check: `sleep 10`, then run `python3 .agents/skills/iterate-pr/scripts/fetch_pr_feedback.py`.
   - Default mode: address any new `high`/`medium` feedback; if changes are needed, return to step 6
   - Merge modes: address any new item; if changes are needed, return to step 6 and re-request review

### 8. Repeat

If step 7 required code changes (from new feedback after CI passed), return to step 2 for a fresh cycle. CI failures during monitoring are already handled within step 7's polling loop.

In merge modes, continue looping until both conditions are true:
- CI checks are green
- No new feedback remains after the latest review request - that is, every item the fetcher
  returned is `resolved`, flagged `replied`, or one you already answered this session. Do not
  wait for the fetcher to return an empty list; it never will while an unresolved review body
  or PR conversation comment exists.

### 9. Finish by Mode

- Default mode: stop after success conditions are met (do not merge automatically)
- `--merge`: once CI is green and no new feedback remains, merge administratively - do not wait for a required-approval status to clear:

```bash
gh pr merge <PR_NUMBER> --admin --delete-branch
```

`--admin` uses the authenticated account's admin/maintainer permissions to bypass the
required-approving-review branch-protection rule; it does not bypass or skip the CI-green
and feedback-resolved checks this skill already enforces in steps 4-7, and it never fabricates
or requests an approving review from another account. If the command fails (e.g. the
authenticated account lacks admin/maintainer rights on the repo, or a check can't be
bypassed), stop and surface the `gh` error to the user - do not retry under a different
identity and do not fall back silently.

- `--release`: complete `--merge` mode merge, then cut a release.

  **First check whether this repository ships its own release skill** (look for
  `cut-release` or similar under `.agents/skills/`). If it does, use it instead of the
  generic steps below - it knows this repo's drafter, tag and post-publish workflow
  conventions, and the generic path does not.

  Otherwise, resolve the repo from the checkout rather than hardcoding it:

  ```bash
  REPO=$(gh repo view --json nameWithOwner --jq .nameWithOwner)
  ```

  1. Find the draft release tag:
     ```bash
     gh release list --repo "$REPO" --json tagName,isDraft --jq '.[] | select(.isDraft) | .tagName'
     ```
  2. Review draft notes:
     ```bash
     gh release view vX.Y.Z --repo "$REPO"
     ```
  3. Publish the release:
     ```bash
     gh release edit vX.Y.Z --repo "$REPO" --draft=false --latest
     ```

  Publishing a GitHub release does not prove the artifact reached its registry. Verify at
  the real surface - the package registry, the version endpoint - not the release page.

## Exit Conditions

**Success (default):** All checks pass, post-CI feedback re-check is clean (no new unaddressed high/medium feedback including review bot findings), user has decided on low-priority items.

**Success (`--merge`):** All checks pass, no new feedback remains after the latest review request (step 8's definition - not an empty fetcher result), every new feedback item has a reply (one per item per session, not one per poll), and the PR is administratively merged (no wait on a required-approval status).

**Success (`--release`):** `--merge` success criteria are met and the draft GitHub release has been published.

**Ask for help:** Same failure after 2 attempts, feedback needs clarification, infrastructure issues, or `gh pr merge --admin` fails (insufficient permissions or an unbypassable check) - surface the error rather than retrying under a different identity or falling back silently.

**Stop:** No PR exists, branch needs rebase, or either fetch script returned `status: "error"`.
A `status: "error"` is never a success condition and never a reason to merge.

## Fallback

If scripts fail, use `gh` CLI directly:
- `gh pr checks --json name,state,bucket,link`
- `gh run view <run-id> --log-failed`
- `gh api repos/{owner}/{repo}/pulls/{number}/comments`
- `gh api repos/{owner}/{repo}/pulls/{number}/reviews`
