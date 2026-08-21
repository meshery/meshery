#!/usr/bin/env bash
# PreToolUse hook: refuse a `git commit` that would carry no Signed-off-by.
#
# This is the Claude Code wiring for `.agents/hooks/require-signoff.sh`, which
# holds the rules and is deliberately harness-agnostic: it takes a command line
# as an argument and reports its verdict as an exit status. Claude Code speaks
# neither of those - it passes the tool call as JSON on stdin, and it treats a
# plain non-zero exit as a non-blocking error, which would let the commit
# through. This adapter translates between the two, so the guard actually
# blocks here rather than only appearing to.
#
# Without it, both sign-off guards are absent under this harness: git reaches
# `ui/.husky/commit-msg` only through `core.hooksPath`, and an agent harness
# that repoints that at its own hooks directory disarms it silently. That is how
# PR #21515 came to carry two commits with no trailer, and a missing trailer is
# not recoverable in place - the DCO check reads its configuration from the
# default branch only, so by the time CI reports it the only remedy left is
# rewriting the branch.
set -uo pipefail

INPUT=$(cat)
[ -z "$INPUT" ] && exit 0

ADAPTER_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)
GUARD="$ADAPTER_DIR/../../.agents/hooks/require-signoff.sh"
[ -x "$GUARD" ] || exit 0

# The payload is read for two fields: the command to judge, and the directory it
# would have run in. They are read one at a time rather than returned together,
# because a command line may hold anything a separator could be mistaken for -
# and bash silently drops a NUL from a command substitution, so even that is not
# a safe one. A payload that is not a Bash tool call, or carries no command,
# yields nothing and leaves the call alone.
field() {
  printf '%s' "$INPUT" | python3 -c '
import json, sys

key = sys.argv[1]

try:
    payload = json.load(sys.stdin)
except ValueError:
    sys.exit(1)

if payload.get("tool_name") != "Bash":
    sys.exit(1)

value = payload.get("tool_input", {}).get("command") if key == "command" else payload.get(key)
if not value:
    sys.exit(1)

sys.stdout.write(value)
' "$1" 2>/dev/null
}

command_line=$(field command) || exit 0
cwd=$(field cwd) || cwd=""

# The guard reads the repository to judge an amend, a reused message and a merge
# in progress, so it has to run where the command would have run.
if [ -n "$cwd" ] && [ -d "$cwd" ]; then
  cd "$cwd" || exit 0
fi

if reason=$("$GUARD" "$command_line" 2>&1); then
  exit 0
fi

printf '%s' "$reason" | python3 -c '
import json, sys

print(json.dumps({
    "hookSpecificOutput": {
        "hookEventName": "PreToolUse",
        "permissionDecision": "deny",
        "permissionDecisionReason": sys.stdin.read().strip(),
    }
}))
'
