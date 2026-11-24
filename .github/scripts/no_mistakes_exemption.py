#!/usr/bin/env python3
"""Classify whether a Meshery pull request is exempt from the no-mistakes gate.

A pull request is exempt only when every changed path is documentation, the
docs website, or CI. The product website, meshery.io, lives in
meshery/meshery.io, not this repository.
"""

from __future__ import annotations

import json
import os
import sys

CAP = 3000

# Product code trees. Markdown under these paths is not a documentation-only change.
CODE_DIR_PREFIXES = (
    "server/",
    "ui/",
    "provider-ui/",
    "mesheryctl/",
    "install/",
    "models/",
)


def is_docs_website(path: str) -> bool:
    # docs/ is the documentation tree and the docs.meshery.io Hugo site.
    # The root CNAME publishes that website.
    return path == "CNAME" or path.startswith("docs/")


def is_markdown_doc(path: str) -> bool:
    lower = path.lower()
    if not (lower.endswith(".md") or lower.endswith(".markdown")):
        return False
    return not path.startswith(CODE_DIR_PREFIXES)


def is_ci(path: str) -> bool:
    if path.startswith(".github/workflows/"):
        return True
    if path.startswith(".github/actions/"):
        return True
    if path.startswith(".github/scripts/"):
        return True
    # CI config files that live directly under .github/, not in
    # samples/, ISSUE_TEMPLATE/, or other subdirectories.
    if path.startswith(".github/") and path.count("/") == 1:
        lower = path.lower()
        if lower.endswith((".yml", ".yaml")):
            return True
    return False


def is_exempt(path: str) -> bool:
    return is_docs_website(path) or is_markdown_doc(path) or is_ci(path)


def load_records(raw: str) -> list:
    raw = raw.strip()
    if not raw:
        return []
    try:
        parsed = json.loads(raw)
    except json.JSONDecodeError:
        parsed = None
    if isinstance(parsed, list):
        return parsed
    if isinstance(parsed, dict):
        return [parsed]
    return [json.loads(line) for line in raw.splitlines() if line.strip()]


def paths_of(record: dict) -> list[str]:
    found = []
    for key in ("filename", "previous_filename"):
        value = record.get(key)
        if isinstance(value, str) and value:
            found.append(value)
    return found


class ClassificationError(Exception):
    """The changed-file list cannot be trusted, so classification is refused."""


def classify(records: list, changed: int) -> tuple[bool, list[str]]:
    """Return (exempt, non-exempt paths) or raise ClassificationError."""
    if changed >= CAP or len(records) != changed or len(records) >= CAP:
        raise ClassificationError(
            "Refusing to classify this pull request: GitHub reports "
            f"{changed} changed files and the API returned {len(records)}. "
            f"The files API stops at {CAP}, so an incomplete list would hide a code change."
        )
    outside = []
    seen = set()
    for record in records:
        if not isinstance(record, dict):
            raise ClassificationError(
                "Refusing to classify this pull request: a files API record was not an object."
            )
        filename = record.get("filename")
        if not isinstance(filename, str) or not filename:
            raise ClassificationError(
                "Refusing to classify this pull request: a files API record is missing a filename."
            )
        for path in paths_of(record):
            if path not in seen and not is_exempt(path):
                seen.add(path)
                outside.append(path)
    return (not outside, outside)


def main() -> int:
    changed = int(os.environ["CHANGED_FILES"])
    with open(os.environ["FILES_JSON"], encoding="utf-8") as handle:
        records = load_records(handle.read())
    try:
        exempt, outside = classify(records, changed)
    except ClassificationError as err:
        print(f"::error::{err}", file=sys.stderr)
        return 1

    output_path = os.environ.get("GITHUB_OUTPUT")
    if output_path:
        with open(output_path, "a", encoding="utf-8") as handle:
            handle.write("exempt=" + ("true" if exempt else "false") + "\n")

    if exempt:
        print(
            f"Exempt: all {changed} changed file(s) are documentation, the docs website, or CI."
        )
        return 0
    print("This pull request changes a path outside documentation, the docs website, and CI.")
    for path in outside[:50]:
        print(f"requires no-mistakes: {path}")
    if len(outside) > 50:
        print(f"... and {len(outside) - 50} more")
    return 0


if __name__ == "__main__":
    sys.exit(main())
