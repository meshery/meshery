#!/usr/bin/env bash
# Zero-tolerance AI-attribution blocker.
#
# Reads PreToolUse JSON from stdin and emits a deny decision if the
# payload (Bash command, Write content, Edit/MultiEdit new_string) contains
# forbidden material as defined in the PATTERNS array below.
#
# Legitimate SDK references such as a bare `import anthropic`,
# `from anthropic import ...`, or the `@anthropic-ai/sdk` package
# identifier are NOT matched and remain usable.

INPUT=$(cat)
[ -z "$INPUT" ] && exit 0

TOOL_NAME=$(printf '%s' "$INPUT" | jq -r '.tool_name // ""' 2>/dev/null)
PAYLOAD=""

case "$TOOL_NAME" in
  Bash)
    PAYLOAD=$(printf '%s' "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null)
    ;;
  Write)
    PAYLOAD=$(printf '%s' "$INPUT" | jq -r '.tool_input.content // ""' 2>/dev/null)
    ;;
  Edit)
    PAYLOAD=$(printf '%s' "$INPUT" | jq -r '.tool_input.new_string // ""' 2>/dev/null)
    ;;
  MultiEdit)
    PAYLOAD=$(printf '%s' "$INPUT" | jq -r '[.tool_input.edits[]?.new_string // ""] | join("\n")' 2>/dev/null)
    ;;
  *)
    exit 0
    ;;
esac

[ -z "$PAYLOAD" ] && exit 0

PATTERNS=(
  'Co-Authored-By:[[:space:]]*(Claude|Anthropic|noreply@anthropic)'
  'Co-authored-by:[[:space:]]*(Claude|Anthropic|noreply@anthropic)'

  '(Generated|Powered|Created|Built|Made|Written|Authored|Assisted|Produced|Crafted|Drafted)[[:space:]]+(with|by)[[:space:]]+(\[)?(Claude|Anthropic)'

  '\[Claude[[:space:]]+(Code|Design)\]'
  'Generated[[:space:]]+with[[:space:]]+\[Claude'

  '🤖[[:space:]]+Generated'

  'claude\.ai/share/'
  'claude\.com/share/'

  '(https?://)?([a-zA-Z0-9-]+\.)*claude\.ai'
  '(https?://)?([a-zA-Z0-9-]+\.)*claude\.com'
  '(https?://)?([a-zA-Z0-9-]+\.)*anthropic\.com'

  '@anthropic\.com'
)

for P in "${PATTERNS[@]}"; do
  if printf '%s' "$PAYLOAD" | grep -iEq "$P"; then
    MATCH=$(printf '%s' "$PAYLOAD" | grep -iEo "$P" | head -1)
    jq -n --arg tool "$TOOL_NAME" --arg match "$MATCH" '{
      hookSpecificOutput: {
        hookEventName: "PreToolUse",
        permissionDecision: "deny",
        permissionDecisionReason: ("Blocked by no-ai-attribution policy. " + $tool + " contained forbidden content: \"" + $match + "\". Zero-tolerance policy: no attribution trailers, no vendor product names in markdown links, no vendor-owned domains, no session links, no co-authoring. Remove the content and retry.")
      }
    }'
    exit 0
  fi
done

exit 0
