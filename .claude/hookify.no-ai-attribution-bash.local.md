---
name: no-ai-attribution-bash
enabled: true
event: bash
action: block
conditions:
  - field: command
    operator: regex_match
    pattern: (Co-Authored-By|Co-authored-by):\s*(Claude|Anthropic|noreply@anthropic)|(Generated|Powered|Created|Built|Made|Written|Authored|Assisted|Produced|Crafted|Drafted)\s+(with|by)\s+\[?(Claude|Anthropic)|\[Claude\s+(Code|Design)\]|🤖\s+Generated|claude\.ai/share/|claude\.com/share/|(https?://)?([a-zA-Z0-9-]+\.)*claude\.ai|(https?://)?([a-zA-Z0-9-]+\.)*claude\.com|(https?://)?([a-zA-Z0-9-]+\.)*anthropic\.com|@anthropic\.com
---

Zero-tolerance AI-attribution policy violated in a Bash command (typically a commit message).

Companion to ~/.claude/hooks/no-ai-attribution.sh. Fires when this directory is the active working directory.
