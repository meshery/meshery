"""Regression checks for paths that must trigger the Meshery edge image build."""

from fnmatch import fnmatchcase
from pathlib import Path
import unittest


REPOSITORY_ROOT = Path(__file__).resolve().parents[2]
WORKFLOW_PATH = REPOSITORY_ROOT / ".github/workflows/build-and-release-edge.yml"


def push_paths_ignore() -> list[str]:
    """Read push.paths-ignore from the edge workflow without extra dependencies."""
    lines = WORKFLOW_PATH.read_text(encoding="utf-8").splitlines()
    push_start = lines.index("  push:")
    push_end = next(
        (index for index in range(push_start + 1, len(lines)) if lines[index] and not lines[index].startswith(" ")),
        len(lines),
    )
    push_lines = lines[push_start:push_end]
    ignore_start = push_lines.index("    paths-ignore:")

    patterns = []
    for line in push_lines[ignore_start + 1 :]:
        if not line.startswith("      - "):
            break
        patterns.append(line.removeprefix("      - ").strip("\"'"))
    return patterns


class EdgeImagePathFilterTest(unittest.TestCase):
    def test_dockerfile_change_triggers_edge_image_build(self) -> None:
        self.assertFalse(
            any(fnmatchcase("install/docker/Dockerfile", pattern) for pattern in push_paths_ignore())
        )

    def test_docs_only_changes_remain_excluded(self) -> None:
        self.assertTrue(
            any(fnmatchcase("docs/content/en/index.md", pattern) for pattern in push_paths_ignore())
        )


if __name__ == "__main__":
    unittest.main()
