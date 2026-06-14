#!/usr/bin/env python3
"""Propagate the canonical remote-provider roster to every install artifact.

`install/providers.env` is the single source of truth. Each `Name=URL` line is an
ACTIVE provider (registered at startup, emitted into PROVIDER_BASE_URLS); lines prefixed
`#inactive ` are DECLARED-but-inactive (documented, not registered until their DNS
resolves). The first active provider is the PRIMARY, used by single-URL consumers.

This script rewrites the managed value in each consumer, regenerates the docker-extension
provider chooser, and emits the Go constant files for the server and mesheryctl.

Usage:
    sync-provider-urls.py            # rewrite all generated artifacts in place
    sync-provider-urls.py --check    # verify everything is in sync; exit 1 + diff if not
"""

from __future__ import annotations

import argparse
import difflib
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PROVIDERS_ENV = ROOT / "install" / "providers.env"
MARKER = "// AUTO-GENERATED from install/providers.env - run `make providers-propagate`"


def die(msg: str) -> "None":
    sys.stderr.write(f"sync-provider-urls: error: {msg}\n")
    raise SystemExit(2)


def _unquote(s: str) -> str:
    """Strip a single pair of surrounding matching quotes, e.g. `"TCS Labs"` -> `TCS Labs`."""
    s = s.strip()
    if len(s) >= 2 and s[0] == s[-1] and s[0] in ('"', "'"):
        return s[1:-1]
    return s


def load_active() -> "list[tuple[str, str]]":
    """Return ordered (name, url) pairs for ACTIVE providers (uncommented lines).

    Names may be quoted (recommended when they contain whitespace, e.g. "TCS Labs");
    surrounding quotes are stripped. Commented lines are declared-but-inactive.
    """
    if not PROVIDERS_ENV.is_file():
        die(f"canonical source not found: {PROVIDERS_ENV}")
    pairs: list[tuple[str, str]] = []
    for raw in PROVIDERS_ENV.read_text().splitlines():
        line = raw.strip()
        if not line or line.startswith("#"):
            continue  # blank or commented (declared-but-inactive) line
        if "=" not in line:
            die(f"malformed active line (expected Name=URL): {raw!r}")
        name, url = line.split("=", 1)
        name, url = _unquote(name), _unquote(url)
        if not name or not url:
            die(f"malformed active line (empty name or url): {raw!r}")
        pairs.append((name, url))
    if not pairs:
        die("providers.env has no active providers")
    return pairs


# ---------------------------------------------------------------------------
# Format-specific replacers. Each raises on a missing/ambiguous managed line so
# drift surfaces loudly rather than silently no-op'ing.
# ---------------------------------------------------------------------------

def _one(matches: "list[int]", path: Path, what: str) -> int:
    if len(matches) != 1:
        die(f"{path}: expected exactly one {what}, found {len(matches)}")
    return matches[0]


def replace_kv(text: str, key: str, value: str, path: Path) -> str:
    """Replace the value in a `<prefix>KEY=<value>[<closing-quote>][ # comment]` line."""
    pat = re.compile(rf'^(?P<pre>.*?{re.escape(key)}=)(?P<val>[^"\n]*?)(?P<post>"?\s*)$')
    lines = text.splitlines(keepends=True)
    idx = _one([i for i, ln in enumerate(lines) if pat.match(ln.rstrip("\n"))], path, f"`{key}=` line")
    nl = "\n" if lines[idx].endswith("\n") else ""
    lines[idx] = pat.sub(rf'\g<pre>{value}\g<post>', lines[idx].rstrip("\n")) + nl
    return "".join(lines)


def replace_k8s(text: str, key: str, value: str, path: Path) -> str:
    """Replace the `value:` line following a `- name: KEY` line in a k8s manifest."""
    lines = text.splitlines(keepends=True)
    name_re = re.compile(rf'^\s*-?\s*name:\s*{re.escape(key)}\s*$')
    value_re = re.compile(r'^(?P<pre>\s*value:\s*).*$')
    hits = [i for i, ln in enumerate(lines) if name_re.match(ln.rstrip("\n"))]
    idx = _one(hits, path, f"`name: {key}` block")
    if idx + 1 >= len(lines) or not value_re.match(lines[idx + 1].rstrip("\n")):
        die(f"{path}: `name: {key}` not followed by a `value:` line")
    nl = "\n" if lines[idx + 1].endswith("\n") else ""
    lines[idx + 1] = value_re.sub(rf'\g<pre>{value}', lines[idx + 1].rstrip("\n")) + nl
    return "".join(lines)


def replace_helm(text: str, key: str, value: str, path: Path) -> str:
    """Replace the quoted scalar in a `  KEY: "<value>"` Helm values mapping."""
    pat = re.compile(rf'^(?P<pre>\s*{re.escape(key)}:\s*")(?P<val>[^"]*)(?P<post>".*)$')
    lines = text.splitlines(keepends=True)
    idx = _one([i for i, ln in enumerate(lines) if pat.match(ln.rstrip("\n"))], path, f"`{key}:` mapping")
    nl = "\n" if lines[idx].endswith("\n") else ""
    lines[idx] = pat.sub(rf'\g<pre>{value}\g<post>', lines[idx].rstrip("\n")) + nl
    return "".join(lines)


def replace_js(text: str, lhs: str, quote: str, value: str, path: Path) -> str:
    """Replace the quoted string assigned by `<lhs><quote>...<quote>;`, appending MARKER."""
    q = re.escape(quote)
    pat = re.compile(rf'^(?P<pre>{re.escape(lhs)}{q})[^{q}]*(?P<mid>{q}\s*;).*$', re.M)
    matches = list(pat.finditer(text))
    _one([m.start() for m in matches], path, f"assignment `{lhs}`")
    return pat.sub(lambda m: f"{m.group('pre')}{value}{m.group('mid')} {MARKER}", text)


def replace_chooser(text: str, active: "list[tuple[str, str]]", path: Path) -> str:
    """Regenerate the docker-extension REMOTE_PROVIDERS array from the active roster."""
    entries = "".join(
        f'    {{\n        name: "{name}",\n        url: "{url}",\n    }},\n'
        for name, url in active
    )
    block = (
        "export const REMOTE_PROVIDERS = [\n"
        "    // BEGIN AUTO-GENERATED from install/providers.env - run `make providers-propagate`\n"
        f"{entries}"
        "    // END AUTO-GENERATED\n"
        "];"
    )
    pat = re.compile(r'export const REMOTE_PROVIDERS = \[[\s\S]*?\];')
    if not pat.search(text):
        die(f"{path}: REMOTE_PROVIDERS array not found")
    return pat.sub(lambda _m: block, text, count=1)


def render_go_mesheryctl(full: str) -> str:
    return (
        "// Code generated by install/scripts/sync-provider-urls.py from "
        "install/providers.env; DO NOT EDIT.\n\n"
        "package utils\n\n"
        "// DefaultProviderBaseURLs is the canonical comma-joined list of Meshery's active\n"
        "// remote providers, compiled into mesheryctl as the fallback used when the install\n"
        "// docker-compose file cannot be downloaded.\n"
        f'const DefaultProviderBaseURLs = "{full}"\n'
    )


def render_go_server(full: str, primary: str) -> str:
    return (
        "// Code generated by install/scripts/sync-provider-urls.py from "
        "install/providers.env; DO NOT EDIT.\n\n"
        "package models\n\n"
        "// DefaultRemoteProviderURLs is the comma-joined list of active default remote\n"
        "// provider URLs. The server seeds it via SetDefault(\"PROVIDER_BASE_URLS\", ...) so\n"
        "// operators who do not set the env var still register the canonical providers.\n"
        f'const DefaultRemoteProviderURLs = "{full}"\n\n'
        "// PrimaryProviderURL is the single canonical provider host used by single-URL\n"
        "// consumers (the built-in local provider's capability paths, SaaS deep links,\n"
        "// provider-ui return-to fallback, e2e defaults).\n"
        f'const PrimaryProviderURL = "{primary}"\n'
    )


def desired_contents() -> "dict[str, str]":
    active = load_active()
    full = ",".join(url for _, url in active)
    primary = active[0][1]

    # (relative path, builder)
    plan = [
        ("install/docker/docker-compose.yaml",
         lambda t, p: replace_kv(t, "PROVIDER_BASE_URLS", full, p)),
        ("install/mesheryapp.dockerapp/docker-compose.yml",
         lambda t, p: replace_kv(t, "PROVIDER_BASE_URLS", full, p)),
         ("install/deployment_yamls/k8s/meshery-deployment.yaml",
         lambda t, p: replace_k8s(t, "PROVIDER_BASE_URLS", full, p)),
         ("install/deployment_yamls/k8s/meshery-deployment.yaml",
         lambda t, p: replace_k8s(t, "PROVIDER_BASE_URLS", full, p)),
        ("install/kubernetes/helm/meshery/values.yaml",
         lambda t, p: replace_helm(t, "PROVIDER_BASE_URLS", full, p)),
        ("install/docker-extension/docker-compose.yaml",
         lambda t, p: replace_kv(t, "PROVIDER_BASE_URLS", primary, p)),
        ("install/playground/docker/docker-compose.yaml",
         lambda t, p: replace_kv(t, "PROVIDER_BASE_URLS", primary, p)),
        ("install/playground/kubernetes/playground-deployment.yaml",
         lambda t, p: replace_k8s(t, "PROVIDER_BASE_URLS", primary, p)),
        ("install/playground/docker/Makefile",
         lambda t, p: replace_kv(t, "REMOTE_PROVIDER", primary, p)),
        ("ui/constants/endpoints.ts",
         lambda t, p: replace_js(t, "export const MESHERY_CLOUD_PROD = ", "'", primary, p)),
        ("provider-ui/lib/data-fetch.js",
         lambda t, p: replace_js(t, "export const PROVIDER_URL = ", '"', primary, p)),
        ("ui/tests/e2e/env.js",
         lambda t, p: replace_js(
             t, "const REMOTE_PROVIDER_URL = process.env.REMOTE_PROVIDER_URL || ", "'", primary, p)),
        ("install/docker-extension/ui/src/components/utils/constants.js",
         lambda t, p: replace_chooser(t, active, p)),
    ]

    out: dict[str, str] = {}
    for rel, builder in plan:
        path = ROOT / rel
        if not path.is_file():
            die(f"managed file not found: {path}")
        out[rel] = builder(path.read_text(), path)
    out["mesheryctl/pkg/utils/providers_gen.go"] = render_go_mesheryctl(full)
    out["server/models/default_remote_providers.go"] = render_go_server(full, primary)
    return out


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--check", action="store_true", help="verify in sync; do not write")
    args = ap.parse_args()

    desired = desired_contents()
    drifted = []
    for rel, content in desired.items():
        path = ROOT / rel
        current = path.read_text() if path.is_file() else ""
        if current == content:
            continue
        drifted.append(rel)
        if args.check:
            sys.stderr.write(f"\n--- drift in {rel} ---\n")
            sys.stderr.writelines(difflib.unified_diff(
                current.splitlines(keepends=True), content.splitlines(keepends=True),
                fromfile=f"a/{rel}", tofile=f"b/{rel}"))
        else:
            path.write_text(content)
            print(f"updated {rel}")

    if args.check and drifted:
        sys.stderr.write(
            f"\n{len(drifted)} artifact(s) out of sync with install/providers.env. "
            "Run `make providers-propagate`.\n")
        return 1
    if not args.check and not drifted:
        print("all artifacts already in sync")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
