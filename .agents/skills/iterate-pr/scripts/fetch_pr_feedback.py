#!/usr/bin/env python3
# /// script
# requires-python = ">=3.9"
# ///
"""
Fetch and categorize PR review feedback.

Usage:
    uv run fetch_pr_feedback.py [--pr PR_NUMBER]
    python3 fetch_pr_feedback.py [--pr PR_NUMBER]

If --pr is not specified, uses the PR for the current branch.
"""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from typing import Any, NoReturn


REVIEW_BOT_PATTERNS = [
    r"(?i)^sentry",
    r"(?i)^warden",
    r"(?i)^cursor",
    r"(?i)^bugbot",
    r"(?i)^seer",
    r"(?i)^copilot",
    r"(?i)^codex",
    r"(?i)^claude",
    r"(?i)^codeql",
]

INFO_BOT_PATTERNS = [
    r"(?i)^codecov",
    r"(?i)^dependabot",
    r"(?i)^renovate",
    r"(?i)^github-actions",
    r"(?i)^mergify",
    r"(?i)^semantic-release",
    r"(?i)^sonarcloud",
    r"(?i)^snyk",
    r"(?i)bot$",
    r"(?i)\[bot\]$",
]

LOGAF_PATTERNS = [
    (re.compile(r"^\s*(?:h:|h\s*:|high:|\[h\])", re.IGNORECASE), "high"),
    (re.compile(r"^\s*(?:m:|m\s*:|medium:|\[m\])", re.IGNORECASE), "medium"),
    (re.compile(r"^\s*(?:l:|l\s*:|low:|\[l\])", re.IGNORECASE), "low"),
]

HIGH_PATTERNS = [
    re.compile(pattern, re.IGNORECASE)
    for pattern in (
        r"must\s+(fix|change|update|address)",
        r"this\s+(is\s+)?(wrong|incorrect|broken|buggy)",
        r"security\s+(issue|vulnerability|concern)",
        r"will\s+(break|cause|fail)",
        r"\bcritical\b",
        r"\bblocker\b",
    )
]

LOW_PATTERNS = [
    re.compile(pattern, re.IGNORECASE)
    for pattern in (
        r"nit[:\s]",
        r"nitpick",
        r"suggestion[:\s]",
        r"consider\s+",
        r"could\s+(also\s+)?",
        r"might\s+(want\s+to|be\s+better)",
        r"optional[:\s]",
        r"minor[:\s]",
        r"style[:\s]",
        r"prefer\s+",
        r"what\s+do\s+you\s+think",
        r"up\s+to\s+you",
        r"take\s+it\s+or\s+leave",
        r"\bfwiw\b",
    )
]


def fail(message: str) -> NoReturn:
    print(json.dumps({"error": message}))
    sys.exit(1)


def run_gh_json(args: list[str]) -> Any:
    try:
        result = subprocess.run(
            ["gh", *args],
            capture_output=True,
            text=True,
        )
    except OSError as exc:
        raise RuntimeError(f"failed to run gh CLI: {exc}") from exc
    if result.returncode != 0:
        stderr = (result.stderr or result.stdout).strip() or "unknown gh error"
        raise RuntimeError(f"gh {' '.join(args)} failed: {stderr}")

    stdout = result.stdout.strip()
    if not stdout:
        return None

    try:
        return json.loads(stdout)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"gh {' '.join(args)} returned non-JSON output") from exc


def get_repo_info() -> tuple[str, str]:
    result = run_gh_json(["repo", "view", "--json", "owner,name"])
    if not isinstance(result, dict):
        raise RuntimeError("could not determine repository from current directory")
    owner = result.get("owner", {}).get("login")
    repo = result.get("name")
    if not owner or not repo:
        raise RuntimeError("could not determine repository owner/name")
    return str(owner), str(repo)


def get_pr_info(pr_number: int | None = None) -> dict[str, Any]:
    args = ["pr", "view", "--json", "number,url,headRefName,author,reviews,reviewDecision"]
    if pr_number is not None:
        args.insert(2, str(pr_number))
    result = run_gh_json(args)
    if not isinstance(result, dict):
        raise RuntimeError("unable to determine PR for current branch")
    return result


def is_review_bot(username: str) -> bool:
    return any(re.search(p, username) for p in REVIEW_BOT_PATTERNS)


def is_info_bot(username: str) -> bool:
    return any(re.search(p, username) for p in INFO_BOT_PATTERNS)


def get_issue_comments(owner: str, repo: str, pr_number: int) -> list[dict[str, Any]]:
    result = run_gh_json([
        "api",
        f"repos/{owner}/{repo}/issues/{pr_number}/comments",
        "--paginate",
        "--slurp",
    ])
    if not isinstance(result, list):
        return []

    comments: list[dict[str, Any]] = []
    for page in result:
        if isinstance(page, list):
            comments.extend(entry for entry in page if isinstance(entry, dict))
        elif isinstance(page, dict):
            comments.append(page)
    return comments


def get_review_threads(owner: str, repo: str, pr_number: int) -> list[dict[str, Any]]:
    query = """
    query($owner: String!, $repo: String!, $pr: Int!, $after: String) {
      repository(owner: $owner, name: $repo) {
        pullRequest(number: $pr) {
          reviewThreads(first: 100, after: $after) {
            pageInfo {
              hasNextPage
              endCursor
            }
            nodes {
              id
              isResolved
              isOutdated
              path
              line
              comments(first: 1) {
                nodes {
                  body
                  createdAt
                  author {
                    login
                  }
                }
              }
            }
          }
        }
      }
    }
    """
    threads: list[dict[str, Any]] = []
    cursor: str | None = None
    while True:
        args = [
            "api",
            "graphql",
            "-f",
            f"query={query}",
            "-F",
            f"owner={owner}",
            "-F",
            f"repo={repo}",
            "-F",
            f"pr={pr_number}",
        ]
        if cursor is not None:
            args.extend(["-F", f"after={cursor}"])

        result = run_gh_json(args)
        if not isinstance(result, dict):
            break

        review_threads = (
            result.get("data", {})
            .get("repository", {})
            .get("pullRequest", {})
            .get("reviewThreads", {})
        )
        nodes = review_threads.get("nodes", [])
        if isinstance(nodes, list):
            threads.extend(nodes)

        page_info = review_threads.get("pageInfo", {})
        if not page_info.get("hasNextPage"):
            break
        cursor = page_info.get("endCursor")
        if not cursor:
            break

    return threads


def detect_logaf(body: str) -> str | None:
    for pattern, level in LOGAF_PATTERNS:
        if pattern.search(body):
            return level
    return None


def categorize_comment(author: str, body: str) -> str:
    if is_info_bot(author) and not is_review_bot(author):
        return "bot"

    logaf_level = detect_logaf(body)
    if logaf_level:
        return logaf_level

    if any(pattern.search(body) for pattern in HIGH_PATTERNS):
        return "high"
    if any(pattern.search(body) for pattern in LOW_PATTERNS):
        return "low"
    return "medium"


def item_sort_key(item: dict[str, Any]) -> str:
    return item.get("_created_at", "")


def extract_feedback_item(
    body: str,
    author: str,
    *,
    created_at: str = "",
    path: str | None = None,
    line: int | None = None,
    url: str | None = None,
    is_resolved: bool = False,
    is_outdated: bool = False,
    review_bot: bool = False,
    thread_id: str | None = None,
) -> dict[str, Any]:
    summary = body[:200] + "..." if len(body) > 200 else body
    summary = summary.replace("\n", " ").strip()

    item: dict[str, Any] = {
        "author": author,
        "body": summary,
        "full_body": body,
        "_created_at": created_at,
    }

    if path:
        item["path"] = path
    if line is not None:
        item["line"] = line
    if url:
        item["url"] = url
    if is_resolved:
        item["resolved"] = True
    if is_outdated:
        item["outdated"] = True
    if review_bot:
        item["review_bot"] = True
    if thread_id:
        item["thread_id"] = thread_id

    return item


def main() -> None:
    parser = argparse.ArgumentParser(description="Fetch and categorize PR feedback")
    parser.add_argument("--pr", type=int, help="PR number (defaults to current branch PR)")
    args = parser.parse_args()

    try:
        owner, repo = get_repo_info()
        pr_info = get_pr_info(args.pr)
    except RuntimeError as error:
        fail(str(error))

    pr_number = pr_info["number"]
    pr_author = pr_info.get("author", {}).get("login", "")
    review_decision = pr_info.get("reviewDecision", "")

    feedback: dict[str, list[dict[str, Any]]] = {
        "high": [],
        "medium": [],
        "low": [],
        "bot": [],
        "resolved": [],
    }

    reviews = pr_info.get("reviews", [])
    if isinstance(reviews, list):
        for review in reviews:
            if not isinstance(review, dict):
                continue
            if review.get("state") != "CHANGES_REQUESTED":
                continue
            author = review.get("author", {}).get("login", "")
            body = review.get("body", "")
            if not author or not body or author == pr_author:
                continue
            item = extract_feedback_item(
                body=body,
                author=author,
                created_at=str(review.get("submittedAt", "")),
            )
            item["type"] = "changes_requested"
            feedback["high"].append(item)

    try:
        threads = get_review_threads(owner, repo, pr_number)
        issue_comments = get_issue_comments(owner, repo, pr_number)
    except RuntimeError as error:
        fail(str(error))

    if isinstance(threads, list):
        for thread in threads:
            if not isinstance(thread, dict):
                continue
            comments = thread.get("comments", {}).get("nodes", [])
            if not comments or not isinstance(comments, list):
                continue

            first_comment = comments[0] if comments and isinstance(comments[0], dict) else {}
            author = first_comment.get("author", {}).get("login", "")
            body = first_comment.get("body", "")
            if not author or not body or author == pr_author or len(body.strip()) < 3:
                continue

            is_resolved = bool(thread.get("isResolved", False))
            thread_id = thread.get("id")
            category = "resolved" if is_resolved else categorize_comment(author, body)

            item = extract_feedback_item(
                body=body,
                author=author,
                created_at=str(first_comment.get("createdAt", "")),
                path=thread.get("path"),
                line=thread.get("line"),
                is_resolved=is_resolved,
                is_outdated=bool(thread.get("isOutdated", False)),
                thread_id=thread_id if isinstance(thread_id, str) else None,
                review_bot=category in {"high", "medium", "low"} and is_review_bot(author),
            )
            feedback[category].append(item)

    if isinstance(issue_comments, list):
        for comment in issue_comments:
            if not isinstance(comment, dict):
                continue
            author = comment.get("user", {}).get("login", "")
            body = comment.get("body", "")
            if not author or not body or author == pr_author or len(body.strip()) < 3:
                continue

            category = categorize_comment(author, body)
            item = extract_feedback_item(
                body=body,
                author=author,
                created_at=str(comment.get("created_at", "")),
                url=comment.get("html_url"),
                review_bot=category in {"high", "medium", "low"} and is_review_bot(author),
            )
            feedback[category].append(item)

    for bucket in feedback.values():
        bucket.sort(key=item_sort_key)
        for item in bucket:
            item.pop("_created_at", None)

    review_bot_count = sum(
        1
        for bucket in ("high", "medium", "low")
        for item in feedback[bucket]
        if item.get("review_bot")
    )

    output = {
        "pr": {
            "number": pr_number,
            "url": pr_info.get("url", ""),
            "author": pr_author,
            "review_decision": review_decision,
        },
        "summary": {
            "high": len(feedback["high"]),
            "medium": len(feedback["medium"]),
            "low": len(feedback["low"]),
            "bot_comments": len(feedback["bot"]),
            "resolved": len(feedback["resolved"]),
            "review_bot_feedback": review_bot_count,
            "needs_attention": len(feedback["high"]) + len(feedback["medium"]),
        },
        "feedback": feedback,
    }

    if feedback["high"]:
        output["action_required"] = "Address high-priority feedback before merge"
    elif feedback["medium"]:
        output["action_required"] = "Address medium-priority feedback"
    elif feedback["low"]:
        output["action_required"] = "Review low-priority suggestions - ask user which to address"
    else:
        output["action_required"] = None

    print(json.dumps(output, indent=2))


if __name__ == "__main__":
    main()
