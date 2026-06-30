#!/usr/bin/env bash
# SessionStart hook — two jobs for ephemeral Claude-Code-on-the-web sessions:
#
#   1. Provision the meshery/schemas source-of-truth repo (and, optionally, the
#      other adjacent coordination repos) as siblings of this repo so the
#      MANDATED schema-first workflow and `consumer-audit` are actually
#      runnable.
#   2. Surface the GitHub tooling policy into the session (see the echo near the
#      end) so every web run uses the right tool for each GitHub operation.
#
# Why job 1 exists: meshery-cloud is schema-driven — wire/DB constructs are
# owned by github.com/meshery/schemas and regenerated here. When the schemas
# repo is absent (as it is by default in a fresh web container), that workflow
# is impossible and the local-struct shortcut becomes tempting. Cloning it up
# front removes that excuse.
#
# Persistence: this file is committed to the repo, so every fresh web session
# (which re-clones meshery-cloud) loads and runs it — it is NOT in ~/.claude.
#
# Best-effort and non-fatal: a clone failure prints a warning to stderr but
# never blocks the session.
set -uo pipefail

# Only meaningful in the remote (web) environment, where the siblings start
# unprovisioned. Local dev clones already follow the "adjacent repos" layout.
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

# ---------------------------------------------------------------------------
# GitHub tooling policy — split by capability, NOT by convenience.
#
# This web environment reaches GitHub two ways; each has exactly one job:
#
#   • GitHub MCP server  →  ONLY for subscribing to pull-request activity
#       (mcp__github__subscribe_pr_activity / mcp__github__unsubscribe_pr_activity).
#       Only the MCP server can stream PR events (CI status, reviews, comments)
#       back into the session; gh cannot. Other read-only MCP calls (search /
#       list / get) are fine, but the subscribe/unsubscribe pair is the reason
#       the server is here.
#
#   • gh CLI  →  ALL GitHub writes.
#       Opening pull requests (gh pr create), responding to review comments
#       (gh api .../pulls/comments/{id}/replies, gh pr comment), reviewing
#       (gh pr review), editing issues, merging — every write goes through gh.
#       It keeps writes on one auditable path that the no-attribution hook
#       inspects. Do NOT use the GitHub MCP *write* tools for these.
#
# The line below is echoed to stdout, which the SessionStart hook injects into
# the session context, so the agent sees this rule on every web run.
# ---------------------------------------------------------------------------
echo "[session-start] GitHub tooling: use the GitHub MCP server ONLY to subscribe/unsubscribe to PR activity (subscribe_pr_activity); use gh for ALL writes — opening PRs (gh pr create), responding to comments (gh pr comment / gh api .../replies), reviews, and merges."

# Configure git identity from the authenticated gh user so every commit's
# Signed-off-by reflects the real developer — not the container default.
# Uses gh api to resolve name and a guaranteed-deliverable noreply address
# (ID+login@users.noreply.github.com) so the hook works for every contributor
# without any personal secrets committed to the shared repo.
if command -v gh >/dev/null 2>&1; then
  _gh_user=$(gh api user --jq '{name:.name,login:.login,id:.id}' 2>/dev/null || true)
  _gh_name=$(printf '%s' "$_gh_user" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('name') or d.get('login',''))" 2>/dev/null || true)
  _gh_login=$(printf '%s' "$_gh_user" | python3 -c "import json,sys; print(json.load(sys.stdin).get('login',''))" 2>/dev/null || true)
  _gh_id=$(printf '%s' "$_gh_user" | python3 -c "import json,sys; print(json.load(sys.stdin).get('id',''))" 2>/dev/null || true)
  # Prefer the primary email from the API; fall back to GitHub's verified
  # noreply address which every developer already has on their account.
  _gh_email=$(gh api user/emails --jq '.[] | select(.primary==true) | .email' 2>/dev/null | head -1 || true)
  if [ -z "$_gh_email" ] && [ -n "$_gh_id" ] && [ -n "$_gh_login" ]; then
    _gh_email="${_gh_id}+${_gh_login}@users.noreply.github.com"
  fi
  if [ -n "$_gh_name" ] && [ -n "$_gh_email" ]; then
    git config --global user.name  "$_gh_name"
    git config --global user.email "$_gh_email"
    echo "[session-start] git identity set to: ${_gh_name} <${_gh_email}>"
  else
    echo "[session-start] WARNING: could not resolve git identity from gh — commits will use the container default." >&2
  fi
fi

# Source-of-truth + coordination repos, cloned as siblings (../<name>) per
# CLAUDE.md "Related Repos". `schemas` is the one the schema-first gate and the
# consumer-audit require, so it is on by default. Uncomment the others to
# provision them too (e.g. ../meshery is needed for
# `consumer-audit --meshery-repo ../meshery`).
REPOS=(
  "meshery/schemas"
  # "meshery/meshery"
  # "meshery/meshkit"
  # "layer5io/sistent"
)

repo_root="${CLAUDE_PROJECT_DIR:-$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
parent="$(cd "$repo_root/.." 2>/dev/null && pwd)" || exit 0

# Clone a repo trying multiple methods in order, falling back from one to the
# next, so provisioning works whether or not `gh` is installed/authenticated:
#   1. gh   — honours the session's gh auth (private repos, SSO).
#   2. HTTPS — these repos are public, so an anonymous clone works with no
#              auth and no gh dependency.
#   3. SSH  — for environments that only have SSH key access to GitHub.
# Returns 0 on the first method that succeeds, 1 if all fail.
clone_repo() {
  local slug="$1" dest="$2"
  if command -v gh >/dev/null 2>&1 && gh repo clone "$slug" "$dest" -- --depth 1 >/dev/null 2>&1; then
    echo "[session-start]   cloned via gh"
    return 0
  fi
  if git clone --depth 1 "https://github.com/$slug.git" "$dest" >/dev/null 2>&1; then
    echo "[session-start]   cloned via https"
    return 0
  fi
  if git clone --depth 1 "git@github.com:$slug.git" "$dest" >/dev/null 2>&1; then
    echo "[session-start]   cloned via ssh"
    return 0
  fi
  return 1
}

for slug in "${REPOS[@]}"; do
  name="${slug##*/}"
  dest="$parent/$name"
  if [ -d "$dest/.git" ]; then
    echo "[session-start] $slug already present at $dest"
    continue
  fi
  echo "[session-start] cloning $slug -> $dest (shallow)"
  if ! clone_repo "$slug" "$dest"; then
    echo "[session-start] WARNING: could not clone $slug by any method (gh/https/ssh) — schema-first changes to constructs it owns will need it cloned manually, e.g.: git clone https://github.com/$slug.git ../$name" >&2
  fi
done

exit 0
