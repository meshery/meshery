#!/usr/bin/env python3
"""Propagate the canonical remote-provider list to every install artifact.

`install/providers.env` is the single source of truth (Name=URL per line). This script
derives the per-artifact PROVIDER_BASE_URLS value and rewrites the managed line in each
install file in place, and emits the generated Go fallback constant for mesheryctl.

Usage:
    sync-provider-urls.py            # rewrite all managed artifacts in place
    sync-provider-urls.py --check    # verify artifacts are in sync; exit 1 + diff if not

Profiles:
    full       - every provider in providers.env, comma-joined (standard installs)
    playground - the Meshery provider only (hosted playground / pre-selected demo)

See docs/superpowers/specs/2026-05-22-canonical-provider-urls-design.md.
"""

from __future__ import annotations

import argparse
import difflib
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PROVIDERS_ENV = ROOT / "install" / "providers.env"

# The provider whose URL alone composes the "playground" profile. It must exist in
# providers.env; it is also the value pre-selected via PROVIDER=<name>.
PLAYGROUND_PROVIDER = "Meshery"


def load_providers() -> "list[tuple[str, str]]":
    """Return ordered (name, url) pairs from providers.env."""
    if not PROVIDERS_ENV.is_file():
        die(f"canonical source not found: {PROVIDERS_ENV}")
    pairs: list[tuple[str, str]] = []
    for raw in PROVIDERS_ENV.read_text().splitlines():
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        if "=" not in line:
            die(f"malformed line in providers.env (expected Name=URL): {raw!r}")
        name, url = line.split("=", 1)
        name, url = name.strip(), url.strip()
        if not name or not url:
            die(f"malformed line in providers.env (empty name or url): {raw!r}")
        pairs.append((name, url))
    if not pairs:
        die("providers.env contains no providers")
    return pairs


def die(msg: str) -> "None":
    sys.stderr.write(f"sync-provider-urls: error: {msg}\n")
    raise SystemExit(2)


# ---------------------------------------------------------------------------
# Per-format value replacement. Each returns the new file content, and raises
# on a missing/ambiguous managed line so drift surfaces loudly.
# ---------------------------------------------------------------------------

def replace_kv(text: str, key: str, value: str, path: Path) -> str:
    """Replace the value in a `<prefix>KEY=<value>[<closing-quote>]` assignment line.

    Handles docker-compose list items (quoted or unquoted) and Makefile variables.
    """
    pattern = re.compile(rf'^(?P<pre>.*?{re.escape(key)}=)(?P<val>[^"\n]*)(?P<post>"?\s*)$')
    matches = [i for i, ln in enumerate(text.splitlines()) if pattern.match(ln)]
    if len(matches) != 1:
        die(f"{path}: expected exactly one `{key}=` line, found {len(matches)}")
    lines = text.splitlines(keepends=True)
    idx = matches[0]
    newline = "\n" if lines[idx].endswith("\n") else ""
    stripped = lines[idx].rstrip("\n")
    lines[idx] = pattern.sub(rf'\g<pre>{value}\g<post>', stripped) + newline
    return "".join(lines)


def replace_k8s(text: str, key: str, value: str, path: Path) -> str:
    """Replace the `value:` line following a `- name: KEY` line in a k8s manifest."""
    lines = text.splitlines(keepends=True)
    name_re = re.compile(rf'^\s*-?\s*name:\s*{re.escape(key)}\s*$')
    value_re = re.compile(r'^(?P<pre>\s*value:\s*)(?P<val>.*?)(?P<post>\s*)$')
    hits = 0
    for i, ln in enumerate(lines):
        if not name_re.match(ln.rstrip("\n")):
            continue
        hits += 1
        if i + 1 >= len(lines) or not value_re.match(lines[i + 1].rstrip("\n")):
            die(f"{path}: `name: {key}` not followed by a `value:` line")
        nl = "\n" if lines[i + 1].endswith("\n") else ""
        lines[i + 1] = value_re.sub(rf'\g<pre>{value}', lines[i + 1].rstrip("\n")) + nl
    if hits != 1:
        die(f"{path}: expected exactly one `name: {key}` block, found {hits}")
    return "".join(lines)


def replace_helm(text: str, key: str, value: str, path: Path) -> str:
    """Replace the quoted scalar in a `  KEY: "<value>"` Helm values mapping."""
    pattern = re.compile(rf'^(?P<pre>\s*{re.escape(key)}:\s*")(?P<val>[^"]*)(?P<post>".*)$')
    matches = [i for i, ln in enumerate(text.splitlines()) if pattern.match(ln)]
    if len(matches) != 1:
        die(f"{path}: expected exactly one `{key}:` mapping, found {len(matches)}")
    lines = text.splitlines(keepends=True)
    idx = matches[0]
    nl = "\n" if lines[idx].endswith("\n") else ""
    lines[idx] = pattern.sub(rf'\g<pre>{value}\g<post>', lines[idx].rstrip("\n")) + nl
    return "".join(lines)


REPLACERS = {"kv": replace_kv, "k8s": replace_k8s, "helm": replace_helm}

# (relative path, format, key, profile)
MANAGED_FILES = [
    ("install/docker/docker-compose.yaml", "kv", "PROVIDER_BASE_URLS", "full"),
    ("install/mesheryapp.dockerapp/docker-compose.yml", "kv", "PROVIDER_BASE_URLS", "full"),
    ("install/docker-extension/docker-compose.yaml", "kv", "PROVIDER_BASE_URLS", "full"),
    ("install/deployment_yamls/k8s/meshery-deployment.yaml", "k8s", "PROVIDER_BASE_URLS", "full"),
    ("install/kubernetes/helm/meshery/values.yaml", "helm", "PROVIDER_BASE_URLS", "full"),
    ("install/playground/docker/docker-compose.yaml", "kv", "PROVIDER_BASE_URLS", "playground"),
    ("install/playground/kubernetes/playground-deployment.yaml", "k8s", "PROVIDER_BASE_URLS", "playground"),
    ("install/playground/docker/Makefile", "kv", "REMOTE_PROVIDER", "playground"),
]

GO_FILE = "mesheryctl/pkg/utils/providers_gen.go"


def render_go(full: str) -> str:
    return (
        "// Code generated by install/scripts/sync-provider-urls.py from "
        "install/providers.env; DO NOT EDIT.\n"
        "\n"
        "package utils\n"
        "\n"
        "// DefaultProviderBaseURLs is the canonical comma-joined list of Meshery's\n"
        "// production remote providers, compiled into mesheryctl as the fallback used\n"
        "// when the install docker-compose file cannot be downloaded.\n"
        f'const DefaultProviderBaseURLs = "{full}"\n'
    )


def desired_contents() -> "dict[str, str]":
    providers = load_providers()
    urls = {name: url for name, url in providers}
    full = ",".join(url for _, url in providers)
    if PLAYGROUND_PROVIDER not in urls:
        die(f"playground provider {PLAYGROUND_PROVIDER!r} not present in providers.env")
    profiles = {"full": full, "playground": urls[PLAYGROUND_PROVIDER]}

    out: dict[str, str] = {}
    for rel, fmt, key, profile in MANAGED_FILES:
        path = ROOT / rel
        if not path.is_file():
            die(f"managed file not found: {path}")
        out[rel] = REPLACERS[fmt](path.read_text(), key, profiles[profile], path)
    out[GO_FILE] = render_go(full)
    return out


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--check", action="store_true",
                    help="verify artifacts are in sync; do not write")
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
            "Run `make generate-install`.\n")
        return 1
    if not args.check and not drifted:
        print("all artifacts already in sync")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
