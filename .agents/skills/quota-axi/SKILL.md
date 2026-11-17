---
name: quota-axi
description: "Report local Claude, Codex, Cursor, GitHub Copilot, and Grok quota windows via the quota-axi CLI - remaining percentages, reset times, and provider status read from local auth sources, with no routing, recommendation, or provider mutation. Use before deciding whether it is safe to keep spending a provider's quota, when the user asks about usage, rate limits, or remaining quota, or when comparing local provider headroom."
user-invocable: false
author: Kun Chen (kunchenguid)
metadata:
  hermes:
    tags: [quota, rate-limits, claude, codex, cursor, copilot, grok, cli]
    category: observability
---

# quota-axi

Report local agent-provider quota windows for routing-aware agents.

You do not need quota-axi installed globally - invoke it with `npx -y quota-axi`.

quota-axi is data only: it never routes, recommends, proxies, intercepts, logs in, imports
browser cookies, or mutates provider state. It reads local provider auth sources and calls
first-party provider quota, usage, billing, or entitlement endpoints; it never launches the
Claude CLI, so it cannot spend the quota it measures.

## When to use

Use quota-axi whenever you need local quota headroom before deciding whether it is safe to
keep working on a provider, when the user asks about usage, rate limits, or remaining quota,
or when comparing supported local provider headroom side by side.

## Workflow

1. Run `npx -y quota-axi` for compact TOON output covering supported providers' quota windows.
2. Scope to one provider with `--provider claude` or to a subset with `--provider cursor,copilot,grok`.
3. Pass `--json` for the normalized machine-readable model instead of TOON.
4. Pass `--full` to include account identity and per-source attempt details.
5. Run `npx -y quota-axi auth` to check local auth-source availability without printing
   secret values.
6. On macOS, Claude Keychain value reads are skipped by default until the user grants access once.
   If quota output reports `reason: keychain_access_required`, tell your user to run
   `quota-axi --allow-keychain-prompt` once and approve Keychain access ("Always Allow").
   After that successful grant, plain `quota-axi` calls reuse the existing Keychain access
   marker to refresh live Claude quota without requiring the flag.
7. For a managed Codex installation, set `QUOTA_AXI_CODEX_BINARY` to its absolute executable
   path. quota-axi uses that exact executable for auth inspection and the read-only app-server
   fallback, and fails closed if the override is invalid.

## Usage

```
usage: quota-axi [auth] [flags]
commands[2]:
  (none)=quota, auth
flags[6]:
  --provider <claude,codex,cursor,copilot,grok>, --json, --full, --allow-keychain-prompt, --help, -v/--version
examples:
  quota-axi
  quota-axi --provider claude
  quota-axi --provider cursor,copilot,grok
  quota-axi --json
  quota-axi --full
  quota-axi auth
```

## Tips

- Output is TOON-encoded and token-efficient by default; pass `--json` only when you need
  the normalized schema.
- Exit code 0 means at least one provider returned data (fresh or stale); exit code 1 means
  every provider failed; exit code 2 means a usage error.
- Percentages are not comparable across providers - quota-axi never claims one provider's
  percentage equals another's.
- Claude `--full` output exposes the authoritative OAuth profile `account.uuid` as
  `account.accountId` when Anthropic returns one; otherwise the account identity is explicitly
  marked unverified rather than inferred.
- The quota cache at `~/.cache/quota-axi/quotas.json` only ever holds normalized
  non-secret snapshots.
  Fresh provider reports with no windows clear stale provider snapshots instead of caching
  empty quota.
  The Claude Keychain access marker lives alongside it and contains no credential values.

---

## Landing through no-mistakes — mandatory

Every pull request that results from work under this skill must use and pass the
no-mistakes pipeline before it merges. This is the repository landing gate in
`AGENTS.md` ("Worktrees and landing a change"), not an optional final check.

- Work in a treehouse worktree, never in the main clone: `treehouse get --lease -b <branch>`,
  and release it with `treehouse return <path>` once the change has merged.
- Commit, then run `no-mistakes axi run --intent "<goal>"` (a clone without a
  `no-mistakes` remote needs `no-mistakes init` once). `no-mistakes axi status` shows
  the step and its findings, and `no-mistakes axi respond` answers a review gate.
- At most three no-mistakes runs per pull request. A run that stopped on a usage limit
  counts; one cancelled before it did anything does not.
- A pull request merges only after a no-mistakes run passes. If it has still not passed
  after three runs, do not merge it as it stands: open an issue for anything still
  flagged and hand the pull request to a human.
- A green CI run, a passing test suite, or a favourable review from another tool does
  not replace a passing no-mistakes run, and a passing no-mistakes run does not waive
  the "Require Human Review" areas listed in `AGENTS.md`.
