---
# Documentation Server Lifecycle Management
#
# This shared workflow provides instructions for starting, waiting for readiness,
# and cleaning up the Meshery Docs local preview server.
#
# Prerequisites:
# - The docs build step has already been attempted in docs/
# - Bash permissions should include: npm *, npx *, curl *, kill *, echo *, sleep *, cat *, rm *
# - Working directory should be the repository root unless noted otherwise
---

## Starting the Documentation Preview Server

**Context**: Meshery Docs is a Hugo site that is typically previewed with the repo-native `npm run serve` command. Start that first, then fall back to `hugo server` only if needed.

Start the preview server in the background and capture its PID:

```bash
cd ${{ github.workspace }}/docs
npm run serve -- --bind 0.0.0.0 --disableFastRender > /tmp/meshery-docs-preview.log 2>&1 &
echo $! > /tmp/meshery-docs-preview.pid
```

If that process exits quickly or never becomes reachable, stop it and fall back to Hugo's built-in server:

```bash
if ! kill -0 "$(cat /tmp/meshery-docs-preview.pid)" 2>/dev/null; then
  npx hugo server -D -F --bind 0.0.0.0 --baseURL http://localhost:1313/ > /tmp/meshery-docs-preview.log 2>&1 &
  echo $! > /tmp/meshery-docs-preview.pid
fi
```

This will:
- Start the preview server on port `1313`
- Redirect logs to `/tmp/meshery-docs-preview.log`
- Save the process ID to `/tmp/meshery-docs-preview.pid` for cleanup

## Waiting for Server Readiness

Poll the local site until it is reachable:

```bash
for i in {1..30}; do
  curl -fsS http://localhost:1313/ > /dev/null && echo "Server ready!" && break
  echo "Waiting for docs server... ($i/30)" && sleep 2
done
```

If the homepage still does not respond after polling, inspect the log and treat it as a critical issue:

```bash
curl -fsS http://localhost:1313/ > /dev/null || cat /tmp/meshery-docs-preview.log
```

## Verifying Server Accessibility

Optionally verify the rendered content:

```bash
curl -fsS http://localhost:1313/ | head -20
```

## Stopping the Documentation Server

When finished, stop the preview server and remove temp files:

```bash
kill "$(cat /tmp/meshery-docs-preview.pid)" 2>/dev/null || true
rm -f /tmp/meshery-docs-preview.pid /tmp/meshery-docs-preview.log
```

## Usage Notes

- The local preview should be reachable at `http://localhost:1313/`
- Prefer the repo-native `npm run serve` command when it works
- Use the Hugo fallback only if the npm preview command fails
- If startup fails, inspect `/tmp/meshery-docs-preview.log` and capture the exact error
