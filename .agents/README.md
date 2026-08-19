# Agent tooling layout

This directory is the LLM-agnostic home for coding-agent configuration in this repo: agent
definitions, packaged skills, and automation hooks. `AGENTS.md` at the repository root states
the rules; this file holds the layout detail behind them.

## Agent definitions

| Agent | File | Purpose |
|-------|------|---------|
| Code Reviewer | `.agents/code-reviewer.md` | Parallel review across Go + frontend |
| Security Reviewer | `.agents/security-reviewer.md` | Security audit |
| Meshery Code Contributor | `.agents/meshery-code-contributor.md` | Full-stack contributions |
| Meshery Docs Contributor | `.agents/meshery-docs-contributor.md` | Hugo docs contributions |
| GitHub Actions Engineer | `.agents/github-actions-engineer.md` | CI/CD design and debugging |
| Relationship Fixture Agent | `.agents/relationship-fixture-agent.md` | Relationship test fixtures |

## Skills

`.agents/skills/` is the single source of truth for every packaged workflow in this repo - one
directory per skill, each with a `SKILL.md`. Do not enumerate them in `AGENTS.md`; list the
directory. Add a new skill only here.

### Per-tool discovery

No skill is ever copied per tool:

| Tool | How it finds these skills |
|---|---|
| Codex | Natively scans `$REPO_ROOT/.agents/skills` - nothing to configure ([docs](https://learn.chatgpt.com/docs/build-skills)) |
| OpenCode | Natively scans `.agents/skills`, one of six roots it searches alongside `.opencode/skills` and `.claude/skills` ([docs](https://opencode.ai/docs/skills/)) |
| Claude Code | Reads `.claude/skills`, which is a relative symlink to `../.agents/skills` |

`.claude/skills` is that symlink and nothing else. **Never replace it with real directories or
copies** - that reintroduces the drift this layout removes.

Neither `.codex/skills` nor `.opencode/skills` is created: both tools already read
`.agents/skills` natively, so a second copy or link would be redundant. `.opencode/skills` is a
real OpenCode search root, just an unnecessary one here; `.codex/skills` is not a path Codex
scans at all.

### Skills address their own files through `.agents/`, never `.claude/`

Skill content must address its own files by their canonical `.agents/skills/...` path.
`iterate-pr` used to invoke its scripts as `.claude/skills/iterate-pr/scripts/<script>.py`,
which resolved only through the symlink and so broke wherever the symlink was absent; it was
corrected to `.agents/` in iterate-pr 2.4.0.

The symlink is therefore a *discovery* path for Claude Code, not a runtime dependency - but it
is still load-bearing for discovery, and the installer hazard below is a further reason not to
touch it.

### AXI installer collision hazard

The four skills tracked in `skills-lock.json` - `chrome-devtools-axi`, `gh-axi`, `lavish`,
`quota-axi` - are installed by the AXI installer, whose layout is skill content at
`.agents/skills/<name>/` *plus* a per-skill symlink at `.claude/skills/<name>`. Those per-skill
symlinks were installer-owned, not hand-made, and this layout removed them as redundant.

The next installer run recreates them, and `.claude/skills/<name>` now resolves *through* the
directory symlink onto `.agents/skills/<name>` - an existing real directory holding the
canonical content. Best case the installer fails with `EEXIST`. Worst case a force-replacing
installer destroys that canonical directory and leaves a self-referential symlink loop.

Which of the two occurs is **not established, and must not be determined by running the
installer against a real checkout**: the failure mode under test is destruction of the
canonical skill content.

### Windows caveat

On a checkout with `core.symlinks=false` - the default outside developer mode - git
materialises `.claude/skills` as a regular text file containing the literal string
`../.agents/skills`. Claude Code then discovers no project skills. Enable Windows developer
mode or set `git config core.symlinks true`, then re-checkout. Skill *scripts* still resolve
there, because skill content addresses them via `.agents/` - only discovery breaks.

## Automation hooks

Scripts in `.agents/hooks/`:

| Hook | Script | Trigger | Purpose |
|------|--------|---------|---------|
| Format Frontend | `.agents/hooks/format-frontend.sh` | Post-edit | Auto-format JS/TS with Prettier |
| Block Lock Files | `.agents/hooks/block-lockfiles.sh` | Pre-edit | Prevent direct edits to lock files |
| Require Sign-off | `.agents/hooks/require-signoff.sh` | Pre-command | Prevent a `git commit` with no `Signed-off-by` |

`block-lockfiles.sh` enforces the no-hand-editing rule by basename, so it covers lock files
that `AGENTS.md` does not enumerate.

### Why the sign-off guard is duplicated

`ui/.husky/commit-msg` already applies the DCO rule, and `require-signoff.sh` applies the same
one. That is not redundancy: the two intercept different things, and only one of them survives
the way this repo is actually worked on.

Git reaches the husky hook **only** through `core.hooksPath`. It is therefore silently absent
on a checkout where `make ui-setup` never ran, and under any coding agent that repoints
`core.hooksPath` at its own hooks directory - which several do, at worktree scope, where it
overrides the repo-local setting husky wrote. Nothing reports the hook as disabled, because an
absent guard and a guard that passed are indistinguishable from the outside.

`require-signoff.sh` reads the proposed command line instead of git's hook plumbing, so
repointing `core.hooksPath` does not disarm it. It is deliberately conservative about what it
lets through - an amend is judged by the message it will reuse, `-F` by the file it names, and a
bundled `-ms` is read the way git reads it, as `-m s` with no sign-off at all.

This matters more than an ordinary lint rule because the failure is not recoverable in place:
the DCO check reads its configuration from the default branch only, so no change on a branch can
relax it, and remediation commits are not enabled for this repository. By the time CI reports a
missing trailer, the only remedy left is rewriting the branch and force-pushing it.

`ui/tests/requireSignoff.test.ts` runs the guard against real repositories and asserts the
verdict it returns for each of those shapes.
