#!/usr/bin/env python3
"""Behavior tests for the no-mistakes exemption classifier."""

from __future__ import annotations

import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

import no_mistakes_exemption as exemption

SCRIPT = Path(__file__).resolve().with_name("no_mistakes_exemption.py")


def records(*paths, renames=()):
    """Build GitHub-like file records.

    Each positional path is a current filename. Each item in renames is
    (previous_filename, filename).
    """
    out = [{"filename": path} for path in paths]
    for previous, current in renames:
        out.append({"filename": current, "previous_filename": previous})
    return out


class IsExemptTests(unittest.TestCase):
    def test_docs_tree_and_cname_are_website_or_docs(self):
        self.assertTrue(exemption.is_exempt("docs/content/en/index.md"))
        self.assertTrue(exemption.is_exempt("docs/hugo.toml"))
        self.assertTrue(exemption.is_exempt("docs/assets/style.css"))
        self.assertTrue(exemption.is_exempt("CNAME"))
        self.assertFalse(exemption.is_exempt("docs"))  # not the docs/ tree
        self.assertFalse(exemption.is_exempt("install/CNAME"))

    def test_root_and_non_code_markdown_are_docs(self):
        for path in (
            "README.md",
            "CONTRIBUTING.md",
            "AGENTS.md",
            "CLAUDE.md",
            "SECURITY.md",
            "VISION.answers.md",
            ".github/copilot-instructions.md",
            ".github/PULL_REQUEST_TEMPLATE.md",
            ".github/ISSUE_TEMPLATE/documentation.md",
            ".github/security-self-assessment.md",
            ".github/agents/meshery-code-contributor.md",
            ".agents/skills/no-mistakes/SKILL.md",
            "CHANGELOG.MD",
            "notes.markdown",
        ):
            with self.subTest(path=path):
                self.assertTrue(exemption.is_exempt(path), path)

    def test_markdown_inside_product_code_trees_is_not_exempt(self):
        for path in (
            "server/README.md",
            "ui/docs/ARCHITECTURE.md",
            "provider-ui/README.md",
            "mesheryctl/README.md",
            "install/kubernetes/helm/meshery/README.md",
            "models/some-model/README.md",
        ):
            with self.subTest(path=path):
                self.assertFalse(exemption.is_exempt(path), path)

    def test_ci_workflows_actions_scripts_and_root_config(self):
        for path in (
            ".github/workflows/no-mistakes-required.yml",
            ".github/workflows/shared/docs-server-lifecycle.md",
            ".github/actions/comment-preview-on-pr/action.yml",
            ".github/actions/comment-preview-on-pr/app/main.py",
            ".github/scripts/no_mistakes_exemption.py",
            ".github/scripts/no_mistakes_exemption_test.py",
            ".github/dependabot.yml",
            ".github/.golangci.yml",
            ".github/.goreleaser.yml",
            ".github/labeler.yml",
            ".github/release-drafter.yml",
            ".github/codeql-config.yml",
            ".github/dco.yml",
            ".github/integration-test-config.yaml",
            ".github/label-commenter-config.yml",
            ".github/pvtr-config.yml",
            ".github/config.yml",
        ):
            with self.subTest(path=path):
                self.assertTrue(exemption.is_exempt(path), path)

    def test_non_ci_github_and_code_paths_are_not_exempt(self):
        for path in (
            ".github/samples/emojivoto.yaml",
            ".github/assets/images/fork.png",
            ".github/ISSUE_TEMPLATE/config.yml",
            ".github/CODEOWNERS",
            ".github/aw/actions-lock.json",
            "Makefile",
            "go.mod",
            "server/meshmodel/help.go",
            "ui/components/Header.js",
            "mesheryctl/internal/cli/root/root.go",
            "install/docker/Dockerfile",
            "skills-lock.json",
            ".agents/hooks/block-lockfiles.sh",
        ):
            with self.subTest(path=path):
                self.assertFalse(exemption.is_exempt(path), path)


class ClassifyTests(unittest.TestCase):
    def test_docs_only_pr_is_exempt(self):
        exempt, outside = exemption.classify(
            records("docs/content/en/index.md", "CNAME"), 2
        )
        self.assertTrue(exempt)
        self.assertEqual(outside, [])

    def test_ci_only_pr_is_exempt(self):
        exempt, outside = exemption.classify(
            records(
                ".github/workflows/no-mistakes-required.yml",
                ".github/actions/upload-test-artifacts/action.yml",
                ".github/scripts/no_mistakes_exemption.py",
                ".github/dependabot.yml",
            ),
            4,
        )
        self.assertTrue(exempt)
        self.assertEqual(outside, [])

    def test_root_markdown_plus_docs_and_ci_is_exempt(self):
        # This pull request's own shape: CI + agent/docs markdown + contributing docs.
        exempt, outside = exemption.classify(
            records(
                ".github/workflows/no-mistakes-required.yml",
                ".github/scripts/no_mistakes_exemption.py",
                ".github/scripts/no_mistakes_exemption_test.py",
                ".github/copilot-instructions.md",
                "AGENTS.md",
                "CONTRIBUTING.md",
                "docs/content/en/project/contributing/_index.md",
                "docs/content/en/project/contributing/build-and-release.md",
                "docs/content/en/project/contributing/contributing-gitflow.md",
            ),
            9,
        )
        self.assertTrue(exempt)
        self.assertEqual(outside, [])

    def test_any_non_exempt_path_blocks_exemption(self):
        exempt, outside = exemption.classify(
            records("docs/README.md", "server/main.go"), 2
        )
        self.assertFalse(exempt)
        self.assertEqual(outside, ["server/main.go"])

    def test_rename_from_code_into_docs_is_not_exempt(self):
        exempt, outside = exemption.classify(
            records(renames=[("server/help.go", "docs/help.md")]), 1
        )
        self.assertFalse(exempt)
        self.assertEqual(outside, ["server/help.go"])

    def test_rename_among_exempt_paths_is_exempt(self):
        exempt, outside = exemption.classify(
            records(renames=[("README.md", "docs/README.md")]), 1
        )
        self.assertTrue(exempt)
        self.assertEqual(outside, [])

    def test_empty_change_list_is_exempt(self):
        exempt, outside = exemption.classify([], 0)
        self.assertTrue(exempt)
        self.assertEqual(outside, [])

    def test_count_mismatch_fails_closed(self):
        with self.assertRaises(exemption.ClassificationError) as ctx:
            exemption.classify(records("README.md"), 2)
        self.assertIn("incomplete list", str(ctx.exception))

    def test_at_or_over_cap_fails_closed(self):
        with self.assertRaises(exemption.ClassificationError):
            exemption.classify([], exemption.CAP)
        with self.assertRaises(exemption.ClassificationError):
            exemption.classify(
                [{"filename": f"docs/{i}.md"} for i in range(exemption.CAP)],
                exemption.CAP,
            )

    def test_missing_or_invalid_filename_fails_closed(self):
        cases = (
            [{"previous_filename": "README.md"}],
            [{"filename": ""}],
            [{"filename": None}],
            ["README.md"],
        )
        for records in cases:
            with self.subTest(records=records):
                with self.assertRaises(exemption.ClassificationError):
                    exemption.classify(records, 1)


class CliTests(unittest.TestCase):
    def _run(self, payload, changed, output_name="github-output"):
        with tempfile.TemporaryDirectory() as tmp:
            files = Path(tmp) / "files.json"
            output = Path(tmp) / output_name
            files.write_text(payload, encoding="utf-8")
            env = os.environ.copy()
            env["CHANGED_FILES"] = str(changed)
            env["FILES_JSON"] = str(files)
            env["GITHUB_OUTPUT"] = str(output)
            proc = subprocess.run(
                [sys.executable, str(SCRIPT)],
                env=env,
                capture_output=True,
                text=True,
                check=False,
            )
            output_text = output.read_text(encoding="utf-8") if output.exists() else ""
            return proc, output_text

    def test_cli_exempt_json_array_writes_output_and_exits_zero(self):
        proc, output = self._run(
            json.dumps([{"filename": "AGENTS.md"}, {"filename": "docs/hugo.toml"}]),
            2,
        )
        self.assertEqual(proc.returncode, 0, proc.stderr)
        self.assertIn("Exempt:", proc.stdout)
        self.assertEqual(output, "exempt=true\n")

    def test_cli_accepts_ndjson_and_reports_non_exempt_paths(self):
        payload = (
            json.dumps({"filename": "docs/index.md"})
            + "\n"
            + json.dumps({"filename": "Makefile"})
            + "\n"
        )
        proc, output = self._run(payload, 2)
        self.assertEqual(proc.returncode, 0, proc.stderr)
        self.assertIn("requires no-mistakes: Makefile", proc.stdout)
        self.assertEqual(output, "exempt=false\n")

    def test_cli_refuses_untrusted_lists(self):
        proc, output = self._run(json.dumps([{"filename": "README.md"}]), 2)
        self.assertEqual(proc.returncode, 1)
        self.assertIn("::error::", proc.stderr)
        self.assertEqual(output, "")


if __name__ == "__main__":
    unittest.main()
