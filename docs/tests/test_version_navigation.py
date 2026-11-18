"""Regression tests for documentation version navigation (Issue #23925)
and breadcrumb layout spacing (Issue #23937).

Verifies that:
1. docs/data/versions.yml only advertises published versions, avoiding 404 links.
2. docs/layouts/partials/header.html does not match root '/' as a prefix for 'latest',
   which erroneously marks 'latest' selected on every page and causes duplicate selected options.
3. docs/assets/scss/_layout.scss provides vertical breathing room between the site header
   and the breadcrumb navigation.
"""

from pathlib import Path
import re
import sys

try:
    import pytest

    @pytest.fixture
    def root() -> Path:
        return get_repo_root()
except ImportError:
    pass


def get_repo_root() -> Path:
    cur = Path(__file__).resolve()
    while cur != cur.parent:
        if (cur / "docs").is_dir() and (cur / ".github").is_dir():
            return cur
        cur = cur.parent
    raise RuntimeError("Could not determine repository root")


def test_versions_yaml(root: Path | None = None):
    root = root or get_repo_root()
    versions_file = root / "docs" / "data" / "versions.yml"
    assert versions_file.is_file(), f"Missing {versions_file}"
    content = versions_file.read_text(encoding="utf-8")

    # Parse active (non-commented) versions
    active_versions = []
    for line in content.splitlines():
        line = line.strip()
        if line.startswith("-") and "version:" in line:
            ver = line.split("version:")[1].strip().strip('"').strip("'")
            active_versions.append(ver)

    # v0.8 and v0.9 snapshots are trimmed by GitHub Pages deployment (1GB limit)
    # and must not be advertised as active versions until dedicated hosting is available.
    for trimmed_ver in ("v0.8", "v0.9"):
        assert trimmed_ver not in active_versions, (
            f"{trimmed_ver} is advertised in versions.yml but trimmed during deployment, leading to 404s"
        )


def test_header_version_selection_logic(root: Path | None = None):
    root = root or get_repo_root()
    header_file = root / "docs" / "layouts" / "partials" / "header.html"
    assert header_file.is_file(), f"Missing {header_file}"
    content = header_file.read_text(encoding="utf-8")

    # Ensure the flawed prefix check on '/' is removed:
    # `{{ if or (eq $current $linkUrl) (strings.HasPrefix $current $linkUrl) }}`
    # where $linkUrl evaluates to "/" matches every URL, causing 'latest' to always be selected.
    flawed_pattern = r'strings\.HasPrefix\s+\$current\s+\$linkUrl'
    assert not re.search(flawed_pattern, content), (
        "header.html uses strings.HasPrefix with $linkUrl ('/'), causing 'latest' to be selected on all pages"
    )

    # Ensure distinct version matching is implemented
    assert "activeVersionUrl" in content or "isVersioned" in content or "isSelected" in content, (
        "header.html must distinctly identify the active version without treating '/' as a universal prefix"
    )


def test_breadcrumb_spacing_in_layout(root: Path | None = None):
    root = root or get_repo_root()
    layout_file = root / "docs" / "assets" / "scss" / "_layout.scss"
    assert layout_file.is_file(), f"Missing {layout_file}"
    content = layout_file.read_text(encoding="utf-8")

    # The breadcrumb navigation must not have zero margin, which causes it
    # to sit flush against the fixed site header (Issue #23937).
    match = re.search(r'\.breadcrumb\s*\{([^}]+)\}', content)
    assert match, "Could not find .breadcrumb rule in _layout.scss"
    rule_body = match.group(1)

    assert "margin: 0;" not in rule_body, (
        ".breadcrumb in _layout.scss has 'margin: 0;', leaving no vertical whitespace under the header"
    )
    # Verify effective top margin is positive in either shorthand margin or margin-top declaration
    positive_top_margin_pattern = (
        r'(?:margin:\s*(?:[1-9]\d*(?:\.\d+)?|0?\.\d+)\w*|margin-top:\s*(?:[1-9]\d*(?:\.\d+)?|0?\.\d+)\w*)'
    )
    assert re.search(positive_top_margin_pattern, rule_body), (
        ".breadcrumb in _layout.scss must specify a positive top margin (via margin shorthand or margin-top)"
    )


def main():
    root = get_repo_root()
    test_versions_yaml(root)
    test_header_version_selection_logic(root)
    test_breadcrumb_spacing_in_layout(root)
    print("All version navigation and layout regression tests passed!")


if __name__ == "__main__":
    try:
        main()
    except AssertionError as exc:
        print(f"FAILED: {exc}", file=sys.stderr)
        sys.exit(1)
