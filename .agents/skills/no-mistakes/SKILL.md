---
name: no-mistakes
description: Validate your code changes through the no-mistakes pipeline - automated code review, tests, lint, docs, push, PR, and CI - before they reach the configured push target. Use when the user asks to run no-mistakes, gate or ship or validate their changes, push safely, asks you to do a task and then validate it, or invokes /no-mistakes.
user-invocable: true
---

# no-mistakes

`no-mistakes` is a local gate that validates your code changes through a pipeline
(intent, rebase, review, test, document, lint, push, PR, CI) before they reach
the configured push target. You drive it through the `no-mistakes axi` command family, which prints
machine-readable [TOON](https://toonformat.dev) to stdout and progress to stderr.


## Active validation-step boundary

A no-mistakes validation-step agent is already inside an active outer run. It
must inspect, fix, and return only its assigned phase. It must never initialize,
start, reattach, rerun, respond to, synchronize, abort, eject, or directly push
a no-mistakes pipeline. Delivery requirements in user intent remain
acceptance context, but the outer executor alone performs the other validation,
push, PR, and CI phases.

`NO_MISTAKES_GATE` is fast diagnostic evidence, not authorization by
itself. The runtime combines managed Git identity with authenticated process
ancestry. If a pipeline-control command returns
`error.code: nested_gate_context`, stop immediately and
return control to the outer executor. Safe inspection remains available through
`no-mistakes axi status`, `no-mistakes axi logs`, help, and
`no-mistakes doctor`.


When the user invokes `/no-mistakes`, report the outcome at the end. If the user
asks for something specific, translate that request into the matching `axi run`
flags yourself - for example, "skip the lint step" becomes `--skip=lint`. Run
`no-mistakes axi run --help` to see the available flags.

## Two ways to invoke

`/no-mistakes` works in two modes, depending on whether the user hands you a
task along with the command:

- **Validate-only** - bare `/no-mistakes` (optionally with flag-style requests
  like "skip the lint step"). The user's code changes are already committed;
  validate them and report the outcome.
- **Task-first** - `/no-mistakes <task>`, e.g.
  `/no-mistakes add a --json flag to the status command`. First carry out the
  task yourself, then validate the result through the pipeline:
  1. **Check scope.** Inspect `git status` before you change or commit anything.
     Preserve unrelated pre-existing uncommitted changes, and when you commit,
     commit only the changes that belong to the user's task.
  2. **Do the work.** Make the changes the task describes, then **commit them on
     a feature branch**. If the user is on the repository's default branch,
     create a feature branch first - the gate validates committed history on a
     non-default branch, so the work must land there before you run.
  3. **Then validate**, passing the user's task as explicit intent. The task
     text is exactly what the user set out to accomplish, in their own words, so
     it *is* the intent - preserve requirements stated directly by the user,
     including constraints, exclusions, acceptance criteria, and later decisions;
     do not condense them into a diff summary or drop them while adding
     implementation context. Enrich it with the decisions and tradeoffs you
     made while doing the work (see
     [Intent is required](#intent-is-required)).


## Test-quality rule

Never add a test whose only evidence is that it opens, reads, greps, parses, or
snapshots implementation source code and finds or omits particular strings,
tokens, lines, commands, function names, prompt phrases, regex matches, AST
shapes, or incidental snapshots. That does not prove behavior: matching text
can be dead or commented out, and a behavior-preserving refactor can change it.

Instead execute a public or executable interface and assert observable behavior,
state, output, side effects, and failure modes. For machine-consumed declarative
artifacts such as workflow YAML, JSON, policy, .gitignore, or generated
configuration, invoke the real consumer when feasible or parse into a typed or
normalized semantic model and assert meaning. A raw substring or regex over the
file is still the anti-pattern.

Reading a file is legitimate when the file is itself generated public output, a
serialized protocol, persisted state, an intentional snapshot, or another
explicitly owned text or byte contract. Name that contract, and do not use its
contents as a proxy that unrelated code works. A natural-language prompt or
instruction is not proven effective because its source contains a sentence.
Deterministic CI may test the final emitted prompt delivered to an agent as an
intentional generated interface; model interpretation belongs in
development-only evaluation, not live-LLM CI.

Use an independent oracle: the expected result must come from somewhere other
than the code under test, such as a specification, worked example, published
constant, external contract, or independently justified property. Name the
public behavior, the oracle's source, and a plausible wrong behavior the test
would reject. Do not only check your own mocks, compare a result to itself, or
copy the implementation's expected-value rule: a shared mistake can pass both
sides. External-boundary mocks, constants, snapshots, and computed expectations
remain legitimate when they check an independent contract.

For a regression, reproduce the reported failure when feasible: the test should
fail before the fix and pass after it.


Everything below - preconditions, intent, the validate-and-decide loop - applies
the same way once the work is committed on a feature branch.

## Before you start

- The work you want validated must be **committed** on a branch. The gate
  validates committed history, not your uncommitted working tree.
- You must be on a **feature branch**, not the repository's default branch.
- The repository must already be initialized with `no-mistakes init`.
- The daemon must have a runnable configured pipeline agent: a supported native
  agent binary, the `agent: cursor` or `agent: devin` ACP alias, or an explicit
  `acp:<target>` through `acpx`. You are the AXI driver, not an
  implicit pipeline-agent backend. If none is available, the run fails
  before its first step; `no-mistakes doctor` reports the configuration problem.

If any of these is not met, `axi run` returns an `error:` with the exact command
to fix it - read it and act on it (commit your work, or create a branch). If the
repository is not initialized, run `no-mistakes init` first; if the `no-mistakes`
command itself is missing or misbehaving, `no-mistakes doctor` reports what is
wrong.
Before starting, run `no-mistakes axi` (home view).
If it shows an active run on your current branch, inspect it with `no-mistakes axi status`.
If it is parked at a gate, drive it with `no-mistakes axi respond`.
Reattach an in-flight run by re-running `no-mistakes axi run` when it still matches your current `HEAD` - either as the submitted head or as the current pipeline head.
Only `no-mistakes axi abort` it when you mean to discard that run before starting over; aborting is a between-runs action, never a way to take over or bypass a gate while a run is still going (see [Validate and decide](#validate-and-decide)).
If it shows an active run on another branch, leave that run alone and start validation for your current branch with `no-mistakes axi run --intent "..."`.

## Intent is required

When you start a new run you must supply explicit intent through exactly one of
`--intent TEXT`, `--intent-file PATH`, or `--intent -` (stdin until EOF):
**what the user set out to
accomplish** - the goal or request behind this work, in their terms. This is not
a description of the diff or the files you changed; it is the objective the
change is meant to achieve. You know it from the conversation, so pass it
directly - no-mistakes uses the supplied intent instead of inferring it from local
agent transcripts. Empty or whitespace-only input is rejected, never inferred.
For input handling, whitespace preservation, and reattachment semantics, see the
[CLI intent input reference](https://kunchenguid.github.io/no-mistakes/reference/cli/#intent-input).

Err on the side of completeness, not brevity. The review step uses explicit intent
to tell a deliberate decision apart from a mistake, so a thin one-line summary
makes it flag things the user already chose. Capture the nuance: the user's
goal, the specific decisions and tradeoffs they made along the way, any
constraints or approaches they ruled in or out, and anything they explicitly
asked for that might otherwise look surprising in the diff. A few sentences to a
short paragraph is normal - write down what you learned from the conversation
that a reviewer reading only the diff would not know.

## Validate and decide

Run the pipeline and decide on its findings as they come up:

1. Start the run. It blocks until the first decision point or the end:
   ```sh
   no-mistakes axi run --intent "<what the user set out to accomplish>"
   ```
   `axi run` and every `axi respond` block synchronously - the review, test,
   and CI steps can each take **several minutes**, so a single call may not
   return for a while. That is normal; do not cancel or re-issue the command
   because it seems slow. Both commands default to `--wait 8m` so a harness
   with a 10-minute tool cap gets a structured return instead of an unbounded
   hang. If the command returns because that wait elapsed, it is not a failed
   run and does not mean the daemon is dead: inspect with `no-mistakes axi status`
   and re-run `axi run` or `axi respond` to reattach. A slow live daemon is
   retried after a health probe rather than treated as I/O failure. To check
   progress without disturbing the run, use `no-mistakes axi status` from a
   separate call.
   A long-running call is working, not stalled - background it if your harness
   needs to, but the run **never advances past a gate on its own**. Read every
   return; on a `gate:`, respond; loop until an `outcome:`. Never idle-wait
   for the run to move forward by itself.
   When that status output includes `awaiting_agent: parked <duration>` under the run,
   the run is parked at an approval or fix-review gate and waiting for you to
   send `axi respond`. The field is observability only: it does not change
   gate resolution, auto-resume the run, or make `--yes` the default.
   While a step is actively `running` or `fixing`, `axi status` may include
   `active_steps` with step-scoped `active_for`, current-round `round_active_for`,
   `last_activity`, a native `agent_pid` when a subprocess agent is running, and the current round such as `round 1`,
   `auto-fix 1/3`, or `fix 2`. If `last_activity` is prefixed with
   `quiet`, no step log or native-agent lifecycle activity has arrived for
   longer than `step_quiet_warning`. Treat that as a liveness clue, not as
   permission to cancel, rerun, or edit the worktree yourself.
2. If the output contains a `gate:` object, the pipeline is waiting on you.
   Read its `findings` table. Each finding has an `id`, `severity`,
   `file`, `description`, and an `action` that tells you how the
   pipeline classified it:
   - `auto-fix` - mechanical and low-risk; you can authorize the fix on
     your own judgment by responding with `--action fix`.
   - `no-op` - informational only; nothing to do.
   - `ask-user` - the finding challenges the user's deliberate intent or
     touches product behavior. This is a call only the user can make - see
     [Escalate `ask-user` findings](#escalate-ask-user-findings) below.

   **Review auto-fix is disabled by default** (`auto_fix.review: 0`; a repo
   or global `auto_fix.review > 0` override re-enables it), so blocking and
   ask-user review findings park for your decision rather than being silently
   self-fixed. (Other steps such as test and lint may auto-fix within the
   pipeline and re-run before they ever gate.)

   Choose one response:
   ```sh
   # accept the step as-is and continue
   no-mistakes axi respond --action approve

   # have the pipeline fix specific findings, then continue
   no-mistakes axi respond --action fix --findings <id1,id2> --instructions "<optional guidance>"

   # fix some findings and explicitly decline the rest
   no-mistakes axi respond --action fix --findings <id1,id2> --ignore <id3>

   # skip this step
   no-mistakes axi respond --action skip

   # continue static review when coverage is incomplete and findings are empty
   no-mistakes axi respond --step review --action complete-coverage

   # reject this gate (only when the user has decided it must fail)
   no-mistakes axi respond --action abort --reason "<why the gate is rejected>"
   ```
   `complete-coverage` runs a fresh reviewer over only the missing paths.
   It reuses stored coverage only for the same run, step, round, head, base, and trusted scope.
   Missing or invalid provenance keeps Review gated. Source changes during the turn also keep it gated.
   It never approves incomplete coverage, runs a code fixer, or answers a question.
   It refuses outstanding findings, finding overrides, and `--yes`.
   With `--action fix`, declines are explicit: every finding the gate shows
   must be listed in `--findings` or `--ignore`, and a response that
   leaves one out is refused with the missing ids named while the gate stays
   parked. A finding an earlier round of this step already decided may be
   omitted to keep that decision - omitting it does **not** decline it - and a
   finding that round chose to fix cannot be declined by a later response at
   all: reverting an applied fix is out of scope for a gate response, so list
   it in `--findings` to fix it again or leave it out to keep it. Every
   successful `fix` response echoes what it recorded in a `recorded:` object
   (`fixed`, `ignored`, `kept`), so read it to confirm the decision
   you actually made.
   While a run is active, never fix findings by editing the code yourself -
   the pipeline owns both the findings and the fixes. Your job at a gate is to
   decide and respond; `--action fix` has the pipeline apply the fix and
   re-review the result. For the same reason, while a run is active do **not**
   `abort` or `rerun` to go fix a finding yourself - even a real bug in
   your own code - because that discards the pipeline's in-flight work and
   forces a full re-validation. `abort` and `rerun` are for *between*
   runs (after a `failed` or `cancelled` outcome), never to circumvent a
   gate.

    Each `respond` blocks until the next `gate:`, `checks-passed` decision point, or final outcome, subject to the same default `--wait 8m` hold.

    A review gate whose findings are `question-<id>` rows is waiting on answers to its reviewer's questions, not on a verdict: answer each with `no-mistakes axi answer --question <id> --answer "<one of its options>"` instead of approving or fixing. The answer that closes the last open question blocks exactly like `respond` and returns the next `gate:` or outcome; any other answer returns at once.

    Extra flags on `respond`:
    - `--wait` bounds the hold (default 8m).
    - `--reason "the operator's explanation"` records an explicitly authorized Test exception with `--step test --action approve`.
      This does not grant approval authority; escalate ask-user findings as before.
      Without a reason, Test approval remains effective; an approval past a failing command, `no-go`, or `inconclusive` verdict is reported as an exception with no operator reason supplied.
    - `--action abort --reason "<why>"` records the user's decision that the parked gate fails. The reason is required.
      The round is recorded as declined, no fix round starts, nothing is approved, pushed or merged, and the run ends `failed` with the reason in the step and run errors. Saved commits, findings and evidence stay available, but the run's worktree is removed, so uncommitted edits in it are lost.
      Follow the output's `branch_sync.next_action` to get the branch back. It is the user's call, never yours: do not send it to get past a finding, and it is not `axi abort`, which cancels the run instead.
    - Prefer `--instructions-file PATH` and `--add-finding-file PATH` for text containing backticks, quotes, or newlines. Write the exact text to a file and pass its path, without shell substitution. Alternatively `--instructions -` or `--add-finding -` reads stdin until EOF, cancellation, or `--wait` expiry (only one stdin input per response). Each file flag is exclusive with its inline flag; empty or whitespace-only file/stdin input is rejected.
    - `--add-finding '<json>'` (with `--action fix`) folds a finding you
      spotted yourself - one the pipeline did not surface - into the fix round,
      as a JSON finding object. Use it for a problem you noticed that is not in
      the gate's own `findings` table.
    - `--ignore <id1,id2>` (with `--action fix`) declines the named findings.
      Every finding in the gate's `findings` table must appear in
      `--findings` or `--ignore`; anything left out is refused, naming the
      ids, so an omission can never be recorded as a decline. A finding an
      earlier round of this step already chose to fix cannot be listed here:
      reverting an applied fix is out of scope for a gate response.
    - `--step <name>` responds to a specific step instead of the one currently
      awaiting approval. You rarely need this; omit it to answer the active gate.
3. Repeat step 2 until the output has an `outcome:` instead of a `gate:`. The
   outcomes are:
   - `checks-passed` - the change is validated and CI is green (or the
     trusted default-branch config declares `no_ci: true` and no checks are
     registered - the help line names that declaration when it applies), but
     the PR is not merged yet. **You are done driving the pipeline.** Do not
     wait for the merge: tell the user the PR is ready and ask them to review
     and merge it (the PR link is in the `help` line). A generic empty forge
     check list without that declaration is not ready. no-mistakes keeps
     monitoring the PR in the background while the monitor is active; use
     `no-mistakes axi handoff` to end green monitoring after handing
     the open PR to human reviewers, or let the configured idle timeout elapse.
     Under the opt-in `ci.release_after_checks_pass` the run instead finishes at
     the verdict: `run.status` is `checks_passed`, the run released its worktree,
     and nothing is left watching the PR.
     An update it needs later is a fresh `no-mistakes axi run` on the updated head;
     this run will not rebase or re-push it.
   - `passed` - the pipeline completed under the requested steps, including any
     explicit per-run skips. This alone is not evidence that a PR was merged.
   - `passed-with-override` - the pipeline completed with an explicitly approved Test exception or CI failure.
     Report the exception, not a clean pass.
     Test evidence is in `run.test_override_reason`, including when CI readiness returns `checks-passed`; do not omit it from the summary.
   - `passed-with-skips` - publication or CI verification automatically skipped.
     Report the missing evidence and its cause from `run.automatic_skips`,
     bound to the full `run.head_sha`. This is neither CI readiness nor a
     failing code verdict. Explicit per-run skips retain their existing behavior.
   - `failed` or `cancelled` - they did not; read the output and address it.
     Follow the custody guidance below before fixing whatever the output
     points at (a failing test, a lint error, a finding you skipped). Commit the
     fix on the same feature branch, then submit it with
     `no-mistakes axi run --intent "..."`. A fresh run or `rerun` is a
     *between-runs* action, correct only after a terminal outcome like this -
     never mid-run to circumvent a gate. Do not leave the user at a `failed`
     outcome without either retrying or explaining what blocks it.

`no-mistakes rerun` keeps its existing head selection: the gate head, or the
latest terminal run's verified unpublished preserved head while custody remains
outstanding. If a known clean caller `HEAD` differs from that selected head,
it refuses before starting or superseding any run and reports both full SHAs.
It never substitutes the caller head or moves either branch to make them match.
On refusal, inspect `no-mistakes axi status` and follow the custody guidance
below. Dirty callers and callers without clean-head evidence retain existing
selection behavior.

When a `failed` outcome came from the agent invocation itself - a provider
processing error, a usage or quota limit, a transport failure, or the agent
process exiting - in Test or an earlier step, and the run never pushed,
`no-mistakes axi retry --run <id>` continues that run instead of repeating it:
a new run starts at the preserved head, carries the steps before the failed
one, and restarts the failed step from the beginning, so a failed Review
reviews the resulting head in full before Test. First adopt the preserved head
with the reported `no-mistakes axi sync --recover` so your clean `HEAD` is
exactly that head. Retry refuses, changing nothing, and names the failed proof
whenever continuing is not provably safe - a failure about the change itself,
a moved or dirty branch, another run, or changed configuration (repository
settings, the agent, or its model since the carried Review ran); then use
`rerun` or address the failure. Like `respond`, it blocks until the next
`gate:` or outcome.

Before any post-pipeline local commit or fresh run, read the structured `branch_sync` object returned by AXI home, status, or a drive result.
Only when its `next_action.code` is `sync`, run `no-mistakes axi sync` first.
That guarded sync may be a strict fast-forward or a content-equivalent diverged advance that anchors the pre-sync head before moving the branch with reset semantics; genuine divergence stays blocked and names `git checkout --detach <pushed-sha>`, which moves no branch ref.
If it reports `next_action.code` is `continue_active_run`, the pipeline still owns the branch: run the reported command, keep driving the active run, and do not make local follow-up commits.
When `next_action.code` is `recover_custody`, run its exact `next_action.command` rather than reconstructing one. That is `no-mistakes axi sync --recover` to take a still-available preserved pipeline head, or `no-mistakes axi sync --recover --keep-local` in two keep-local cases: when an accessible gate confirms the verified preserved head is missing and you are explicitly discarding those unpublished commits, or when a bound archive proves divergent later work remains preserved while recovery keeps the branch at the exact reported required head and never selects, merges, or replays the archive. Do not substitute plain `--recover` or `rerun` for a reported keep-local action. `no-mistakes rerun` can resume validating a still-available ordinary preserved head instead, subject to the clean-head check above.
Ordinary recovery takes that head by fast-forward, or by adopting a diverged preserved head proven to carry every local change - the ordinary result of the pipeline rebasing your commits onto a newer base - after anchoring your pre-recovery head under `refs/no-mistakes/recover-local/<run>`.
The ordinary containment proof is deliberately narrow, so a rebase whose fix rounds also rewrote your own lines refuses instead of being adopted: when nothing can tell a deliberate pipeline fix from a dropped change, the decision is yours.
When a run failed in a step after Test without touching the code (a Document or Lint fetch error) and its recorded head was never verified, recovery still returns custody at that head once your clean checkout sits exactly on it and the gate proves the head is the reviewed result plus the pipeline's own later commits on top of the live gate branch; recovery never moves a behind worktree onto an unverified head, so check that exact head out first.
When `next_action.code` is `recover_remote_rewritten`, the configured push target was force-rewritten outside the pipeline after a terminal run: run exact `no-mistakes axi sync --recover`. It re-verifies the live target, anchors the superseded pipeline head under `refs/no-mistakes/recover-rewritten/<run>/<generation>`, and rebinds only the recorded push binding to the verified live head; it never moves your branch, the gate, or the remote. It refuses `--keep-local`, a live head or target that changes during recovery, a head it cannot anchor, and a merged or closed PR. Afterwards follow the ordinary `next_action` it reports.
When `next_action.code` is `adopt_published`, a custody-returned branch was rebased after its gate lane stopped moving: run `no-mistakes axi sync --adopt-published`. It verifies the configured push target already has the exact rebased local head, preserves the old lane head, and updates only that stale gate lane. If the target differs or changes during verification, it refuses without replacing the lane.
A `branch_sync.state` of `user_owned` means the run went terminal before changing the submitted head and cancellation released the branch: the exact branch and head are yours and immediately usable for whichever delivery path is authorized - no sync action is needed, and a repeated `--recover` there is a harmless no-op.
A dirty worktree, or divergence that cannot be proven contained, makes the recovery refuse with explicit choices; `--keep-local` keeps your current head while the preserved commits stay anchored under `refs/no-mistakes/recover/<run>`. The same flag is the recovery when an accessible gate confirms that the verified preserved head is missing and recovery refs are compatible: it returns custody at the current local head without requiring that object.
If synchronization is blocked, process that structured state instead of improvising reset, stash, merge, rebase, force, or branch replacement.
After synchronization, commit the follow-up on top and re-run `no-mistakes axi run --intent "..."` with the original user intent.
This preserves every prior gate-fix commit regardless of its configured subject.

The CI step deliberately keeps watching the PR after checks pass, so
`axi run` returns `checks-passed` the moment checks are green (or a trusted
`no_ci: true` declaration covers a zero-check repository) rather than
blocking on the human merge. Never poll or re-run waiting for the merge yourself.
That monitor is live by default; only the opt-in
`ci.release_after_checks_pass` finishes the run at the verdict instead.
Never treat "no CI checks reported" alone as green.

Because that monitor stays live by default, a PR that falls behind the default branch or
hits a merge conflict after checks pass - commonly because another PR merged
first - needs **no command from you**: never hand-rebase. When the CI monitor
sees an actual conflict it **rebases onto the base, resolves it, revalidates from Review
because rebasing cannot prove continuity with the reviewed head, and re-pushes
the branch through Push**. On GitHub, a clean PR whose base requires branches
to be up to date (`mergeStateStatus: BEHIND`) still reports `checks-passed` and
waits for a human to bring it up to date; only when the repository opts
into `ci.rebase_behind` does the monitor give it the same rebase, once its
checks are green and the base has held still for one poll. A PR that is merely behind a base without that rule needs nothing, since
the platform merges it. The one exception is when that monitor is
no longer running - the PR was closed, the run was aborted, superseded, or
explicitly handed off, it idle-timed-out, or its auto-fix attempts were
exhausted - in which case recover with `no-mistakes rerun`, subject to the
clean-head check above. An accepted
rerun cancels the stale monitor and re-runs the full pipeline including a
deterministic rebase step. Do **not** reach for
`no-mistakes axi run` to refresh a still-active PR: with the default
monitor, after `checks-passed` it
reattaches to the running monitor (HEAD unchanged) and returns its output
without rebasing.

On a successful outcome (`checks-passed`, `passed`, or `handed-off`), close the loop with the
user: summarize what happened during the pipeline in a concise, easily readable
format - what was validated and what was found. If the output includes a
`fixes` table, the pipeline fixed findings your original change missed:
acknowledge those misses and explicitly list each fix so the user can easily
review them.

## Escalate `ask-user` findings

A gate whose findings are all `auto-fix` or `no-op` is safe to drive on your
own judgment: respond with `--action fix` or `--action approve` as
appropriate. But a finding marked
`ask-user` is a decision that belongs to the user, not you - the pipeline
flagged it because it challenges their deliberate intent or changes product
behavior. Do not approve, fix, or skip it on your own. Instead, stop and bring
it to the user before you respond:

- Relay each `ask-user` finding to them as the pipeline wrote it - its
  `id`, `file`, and full `description` verbatim. Do not paraphrase,
  summarize away the detail, or pre-judge the answer.
- Ask how they want to proceed, then translate their decision into the matching
  `respond` call: `--action fix` (pass their guidance through
  `--instructions`), `--action approve`, or `--action skip`.

The exception is `--yes` (below): it is the user's standing consent to
drive eligible gates unattended, so under `--yes` you resolve ordinary
`ask-user` findings automatically instead of stopping to ask.

If you have clear consent to drive the run automatically, pass `--yes` to `axi run`
or `axi respond`. For eligible gates, it treats actionable findings - `auto-fix` and
`ask-user` alike - as consent to fix it, selects every current finding for one
fix round, accepts the resulting fix review, and approves gates with only
`no-op` findings. Only use it when the user has asked you to drive the whole
run without checking back.

A `protected-path-refusal` gate still requires an explicit operator response
under `--yes`. Relay its path and rule; do not automatically fix, approve,
or skip it. Approval is rejected. Have the operator inspect and resolve the
reported edit, then send `--action fix` to retry the unfinished step.
The [protected-path reference](https://kunchenguid.github.io/no-mistakes/reference/repo-config/#protected_paths)
owns the staging guard's scope and limitations.

A gate holding open `question-<id>` findings is not eligible either: a question
is settled by `axi answer`, never by a verdict, so `--yes` stands aside
there - as it does when the question history cannot be read end to end, which
also parks for a human.

A `test-agent-unvalidated-work` finding means a timed-out Test agent left
commits or changes no Test turn validated. Approval is rejected, so `--yes`
stops at that gate without responding. Relay what the finding names and do
not skip Test, which would publish that work. Ask the operator to choose:
`--action fix` spends another agent budget to validate the work, or
`no-mistakes axi abort` stops and cancels a discarded run; its behavior is unchanged.

At a green open-PR CI monitor, `no-mistakes axi handoff` ends monitoring when the PR is handed to human reviewers. It requires authoritative `checks-passed` readiness and records terminal `handed_off`, without merging or closing the PR. Before green, resolve the checks or gates with `axi run` instead.

## Inspecting state

```sh
no-mistakes axi               # home view: current branch, active runs, next steps
no-mistakes axi status        # full detail plus cached branch_sync when relevant
no-mistakes axi sync --check  # freshly verify an offered synchronization plan
no-mistakes axi sync          # apply only an offered guarded synchronization
no-mistakes axi sync --recover  # return custody after a terminal run left unpublished pipeline commits
no-mistakes axi sync --adopt-published  # adopt an exactly published rebased head into its stale gate lane
no-mistakes axi logs --step <name> --full   # one step's recorded findings, complete summary, and full log
no-mistakes axi handoff       # end green CI monitoring for human review (handed_off, PR stays open)
no-mistakes axi abort         # cancel the current-branch active run
no-mistakes axi abort --run <id>   # cancel a specific run by id (works outside its worktree)
no-mistakes axi respond --run <id> --action approve   # answer a specific run's parked gate by id (works outside its worktree)
```

## Reading the output

- Output is TOON: `key: value` pairs, `name[N]{cols}:` tables, and `help[N]:` hints.
- `axi status` is scoped to your current branch when `--run` is omitted: with a known current branch, an implicitly resolved `run:` is this branch's. A run under `other_branch_run:` is one you named with `--run <id>` that belongs to another branch - never read its status or outcome as your own work. An explicit `--run <id>` rendered under `run:` while the current branch is unknown (detached `HEAD` or a branch-lookup failure) encodes no branch relationship. In a successful status response, no run object at all means this branch has no run yet, whatever the recent-runs table lists; an `error:` response proves nothing about run ownership, so act on the error instead of concluding the branch is idle.
- A non-terminal run object may include `awaiting_agent: parked <duration>` immediately after `status`; that means the run is parked at a gate. Only an implicitly resolved current-branch gate offers bare `axi respond`; an explicit `--run <id>` status is inspection-only even when its branch matches, because the branch may have a newer active run, and names `axi respond --run <id>` as its response command. Follow the response's `help`.
- A run object with a `running` or `fixing` step may include an `active_steps` table. `active_for` is the enclosing step duration; `round_active_for` is the displayed execution or fix round duration and resets for a fix round. Older runs without round timing leave `round_active_for` empty.
- The `help` list at the bottom of most responses tells you the next commands to run.
- Errors are printed as `error: ...` on stdout with a `help` list; act on the suggestion.
- Exit codes: `0` success, no-op, or normal decision gates, `1` failed or cancelled final outcomes, `2` bad usage.

A `gate:` waiting on you looks roughly like this - a `gate:` line naming the step, optional step-specific fields such as `note`, a `findings[N]{...}:` table with one row per finding, and a `help[N]:` list of next commands:

```
gate: review
note: Review auto-fix is disabled by default (auto_fix.review: 0; a repo or global auto_fix.review > 0 override re-enables it), so blocking and ask-user review findings park for your decision rather than being silently self-fixed.
findings[2]{id,severity,file,line,action,description}:
  r1,warning,internal/pipeline/executor.go,,auto-fix,Error from os.Remove is ignored
  r2,error,cmd/no-mistakes/main.go,,ask-user,New --force flag bypasses the confirm prompt
help[6]:
  Run `no-mistakes axi respond --action approve` to accept this step and continue
  Run `no-mistakes axi respond --action fix --findings <ids> [--ignore <ids>]` to have the pipeline fix the selected findings (do not edit files yourself); list every finding below in `--findings` or `--ignore`, unless an earlier round of this step already decided it, and a finding that round chose to fix cannot be declined by a later response
  Run `no-mistakes axi respond --action skip` to skip this step
  Run `no-mistakes axi logs --step review --full` to read the complete step summary and log
  A long-running call is working, not stalled - background it if your harness needs to, but the run never advances past a gate on its own. Read every return; on a `gate:`, respond; loop until an `outcome:`.
  Commit post-pipeline follow-up work on top of the existing branch so every pipeline fix commit remains present. Never abort-and-restart, reset, or replace the branch in a way that drops prior gate-fix commits.
```

Read the `action` column per row: decide `r1` (auto-fix) on your own
judgment - `respond --action fix --findings r1` hands it to the pipeline to
fix - but stop and escalate `r2` (ask-user) to the user before responding. A
final state
instead shows `outcome: <checks-passed|passed|passed-with-override|passed-with-skips|handed-off|failed|cancelled>` with no
`findings` table. Field names and exact columns can vary by step and version,
so read the actual `findings` header rather than assuming this layout.
