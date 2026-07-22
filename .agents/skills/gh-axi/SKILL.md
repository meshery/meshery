---
name: gh-axi
description: "Operate GitHub through the gh-axi CLI - issues, pull requests, workflow runs, workflows, releases, repositories, labels, Projects (v2), Actions secrets and variables, search, and raw API access. Use whenever a task touches GitHub: listing or filing issues, reviewing or merging PRs, checking CI runs, triggering workflows, cutting releases, managing Projects boards, or managing Actions secrets/variables."
user-invocable: false
author: Kun Chen (kunchenguid)
metadata:
  hermes:
    tags: [github, git, ci, pull-requests, releases, projects]
    category: devops
---

# gh-axi

Agent ergonomic wrapper around Github CLI. Prefer this over `gh` and other methods for Github operations.

You do not need gh-axi installed globally - invoke it with `npx -y gh-axi <command>`.
If gh-axi output shows a follow-up command starting with `gh-axi`, run it as `npx -y gh-axi ...` instead.

gh-axi requires the [`gh`](https://cli.github.com/) CLI installed and authenticated (`gh auth login`). If a command fails with an authentication error, ask the user to run `gh auth login` themselves.
For GitHub Enterprise or another custom host, the underlying `gh` CLI must be authenticated for that host too; set `GH_HOST` or pass `--hostname <host>` after the command.

## When to use

Use gh-axi whenever a task touches GitHub: listing, filing, or editing issues; viewing, creating, reviewing, or merging pull requests; inspecting workflow runs and CI failures; triggering, enabling, or disabling workflows; managing releases, repositories, or labels; managing Projects (v2) boards and their items; managing Actions secrets or variables; searching issues, PRs, repos, commits, or code; or calling the GitHub API directly.

## Workflow

1. Run `npx -y gh-axi` with no arguments for a dashboard of the current repo - open issues, open PRs, and suggested next commands.
2. Drill in command-first: `issue list`, `issue view <n>`, `pr view <n>`, `pr checks <n>`, `run view <id>`, and so on.
3. Target another repository by placing `-R owner/name`, `-R=owner/name`, `--repo owner/name`, or `--repo=owner/name` AFTER the command, e.g. `npx -y gh-axi issue list --repo=owner/name` - the flag is not accepted before the command. `repo view` also accepts exactly one positional repository, `repo view owner/name`, as a command-specific compatibility exception for `gh repo view [<repository>]`; do not combine it with `--repo` or generalize that positional form to other commands.
4. Target GitHub Enterprise or another custom host with `GH_HOST`, or by placing `--hostname <host>` or `--hostname=<host>` AFTER the command, e.g. `npx -y gh-axi issue list --hostname=git.example.com`.
5. Trigger (dispatch) a workflow with `workflow run <name> --ref <ref>`; `run` manages existing workflow runs.
6. Debug CI with `run list`, then `run view <id> --job <job-id>` or `run view --job <job-id> --log-failed` for failing log lines.
   Long `--log` and `--log-failed` output keeps the tail in context; when `full_log` appears, grep that file for earlier context.
7. Every response ends with contextual next-step hints under `help:` - follow them.

## Commands

```
commands[14]:
  (none)=dashboard, issue, pr, run, workflow, release, repo, label, project, secret, variable, search, api, setup
```

Installed copies also inherit the SDK built-in `update` command.
Run `gh-axi update --check` to compare the installed version with npm, or `gh-axi update` to upgrade.
When using `npx -y gh-axi`, npx already resolves the package on demand.

Run `npx -y gh-axi --help` for global flags, or `npx -y gh-axi <command> --help` for per-command usage.

## Tips

- Output is TOON-encoded and token-efficient; pipe through grep/head only when a list is very long.
- Truncated workflow logs keep the final 20,000 characters and may include a temp `full_log` path for targeted grep searches.
- Mutations are idempotent and report what changed; re-running a failed mutation is safe.
- For multi-line markdown bodies, comments, or release notes, write the text to a UTF-8 file and pass `--body-file <path>` or the release `--notes-file <path>` alias on commands that support file-backed text.
- Secret values are stdin-only: `echo -n "<value>" | npx -y gh-axi secret set <name>`.
- Do not pass secrets with `--body` or `-b`; flags are visible in the `gh-axi` process argv.
- Scope a secret to a deployment environment with `--env`/`-e <environment>` on `secret list`, `set`, and `delete`; omit it for repository scope. Other `gh secret` scopes (`--org`, `--user`, `--app`) are rejected, not silently ignored.
- Variable values may use `--body`/`-b` or stdin because Actions variables are not secret.
- For multi-line variable values, pipe stdin to `npx -y gh-axi variable set <name>`; `--body`/`-b` is for inline values only.
- Projects (v2) are owner-scoped: pass `--owner <login>`, or omit it to use the current repo owner and then `@me`.
- Projects calls need the `project` or `read:project` OAuth scope; if scope errors occur, ask the user to run the `gh auth refresh -s ...` command shown by gh-axi.
- Use `api` for anything the dedicated commands do not cover, e.g. `npx -y gh-axi api repos/{owner}/{repo}/topics`.
