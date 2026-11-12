---
name: gnhf
description: Use when the user asks to run GNHF, says they are going to sleep or leaving and wants an agent-managed coding run, asks to supervise, steer, or review an active GNHF run, or gives feedback on GNHF results.
---

# GNHF

## Overview

GNHF is an agent orchestrator: it repeatedly calls another coding agent until a natural-language stop condition is met. This skill teaches the host agent to prepare one durable run and, in Companion mode, steer or review it.

Core rule: the host agent orchestrates; GNHF executes. Do not manually implement inside the same scope while a GNHF worker is responsible for it unless the user explicitly changes the delegation.

In Companion mode, GNHF completion is not user acceptance. "Stop condition met" only means the worker stopped; the host still compares the result to the user's latest requirements and fresh verification.

## Modes

Choose exactly one mode for the run.

### Hands-Off

Use when the task is bounded, verification is clear, and the user wants one configured run to proceed without steering.

- Prepare a precise prompt with constraints, non-goals, verification, and stop condition.
- Launch GNHF and wait for completion.
- Intervene early only for hard failure, runaway scope, destructive behavior, or impossible prerequisites.
- Report the final GNHF status after exit.

Examples:

- English: "I'm going to bed. Use GNHF with Copilot to keep working on this branch and stop when the test suite passes."

### Companion

Use when the task is uncertain, exploratory, design-heavy, research-heavy, or likely to need course correction.

Default to Companion when the user asks to iterate until satisfied, requests multi-round work, provides review findings, asks for design/skill/documentation improvement, or asks for supervision.

- Keep a note of original intent, branch, session id, and last known result.
- Poll the active GNHF process until exit or until an intervention point appears.
- Intervene when the worker optimizes the wrong thing, repeats failed fixes, skips requested research, drifts scope, or claims success without evidence.
- Treat review findings as the next acceptance criteria.
- Prefer a new bounded GNHF prompt over manually taking over implementation.

Examples:

- English: "Run GNHF for a few rounds on this onboarding flow. Check the diff between rounds and tighten the next prompt if it starts polishing the wrong thing."

## Launch

Check the installed CLI before relying on flags:

```bash
gnhf --help
```

Known shape:

```bash
gnhf \
  --agent <agent> \
  --max-iterations <n> \
  --stop-when "<observable completion condition>" \
  --prevent-sleep on \
  "<worker prompt>"
```

If GNHF has no `--model` flag, put model requirements in the worker prompt or backend config. Do not invent unsupported flags.

Before launch:

```bash
git status --short
git branch --show-current
git log --oneline --max-count=5
```

Prompt skeleton:

```text
Objective: <one concrete outcome>.

Use <agent/model requirement>. Work in this repo. Treat this as a long-running GNHF task.

Before coding, inspect the current repo, relevant docs, and recent commits. Preserve user changes. Do not make unrelated refactors.

After each meaningful slice, run relevant verification. If blocked, commit no fake success; leave notes with the blocker and evidence.

Stop only when: <observable completion condition>.
```

## Steer

In Companion mode, evaluate after each iteration or meaningful output chunk:

| Signal                                                        | Action                                                               |
| ------------------------------------------------------------- | -------------------------------------------------------------------- |
| Worker found a real blocker                                   | Stop or relaunch with blocker-specific instructions                  |
| Good partial slice                                            | Let it continue or tighten the next stop condition                   |
| Skipped requested research                                    | Relaunch with research as explicit first deliverable                 |
| Worker changes unrelated files                                | Stop and review before continuing                                    |
| Worker claims success without verification                    | Review immediately; relaunch only with evidence-based stop condition |
| Reviewer finds a blocking issue or user says not satisfactory | Relaunch with that finding as the sole bounded correction            |

Steering prompt:

```text
Continue from the current repo state. The previous run partially succeeded: <evidence>.

Do not redo completed work. Focus only on <bounded correction>.

The issue to fix now is <specific observed issue>. Verify with <commands/checks>.

Stop only when <observable condition>.
```

## Companion Review

Use only in Companion mode, when the host is supervising quality or deciding whether to continue with another bounded run.

1. Inspect branch, status, commits, changed files, and diff.
2. Read GNHF notes/logs as claims, not evidence.
3. Run independent verification: tests, lint, build, typecheck, manual QA, or domain-specific checks.
4. Compare the result to the stop condition and the user's latest feedback.
5. Decide: **Mergeable**, **Needs follow-up GNHF run**, or **Do not merge**.

If the result needs follow-up, continue in Companion mode instead of presenting the run as complete. Do not merge unless explicitly authorized.

## Findings

Use when the user provides findings such as "not preserved", "scope drift", "missing requirement", or "why did you stop".

1. Treat the run as Companion mode.
2. Convert each finding into an observable correction. Preserve severity, file/line scope, and the user's wording.
3. Relaunch on the same candidate branch when salvageable.
4. Prompt the worker to fix only the bounded finding, preserve completed valid work, verify, and stop only when the finding is no longer true.
5. Review again after the follow-up.
6. Repeat until no blocking findings remain, verification passes, or a real blocker is found.

## Morning Review

Use when the user returns with "good morning", "how did last night's run go?", or similar after a GNHF run.

Do not ask what to review first. Reconstruct state:

```bash
git status --short
git branch --show-current
git log --oneline --decorate --max-count=20
pgrep -fl 'gnhf|claude|codex|copilot|opencode|rovodev' || true
```

Inspect likely GNHF branches, notes, logs, terminal sessions, and changed files. If a GNHF process is still running, report that first.

Report mode, agent, branch, status, changes, verification, stop-condition result, quality assessment, and recommended next action. Nev[118;1:3uer summarize an overnight run from memory.

## Agent

The supported `--agent` roster comes from `gnhf --help`; the [Agents table](../../README.md#agents) in the GNHF README owns per-agent requirements. Do not hard-code the roster.

- Default to the agent the user explicitly requested, or the one already configured and authenticated locally.
- `codex`: repo-aware code work or review-heavy tasks.
- `claude`: reasoning-heavy implementation or prose-heavy planning when configured.
- `acp:<target>`: when the user wants to drive a custom ACP-compatible agent through GNHF.

## Safety

- Preserve user changes. Never run destructive git commands to clean up a GNHF branch.
- In Companion mode, do not trust a worker's success summary without fresh verification.
- Keep prompts outcome-based and evidence-based.
- Use concrete stop conditions. Bad: "looks good". Good: "the target workflow succeeds, relevant checks pass, and no unrelated files changed."
- If the user is away, produce branches and a status report, not irreversible changes, unless explicitly authorized.
