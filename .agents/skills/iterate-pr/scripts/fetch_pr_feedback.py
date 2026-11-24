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

Output contract: a single JSON object on stdout carrying a top-level ``status``
field. ``status`` is ``"ok"`` when the fetch completed and ``"error"`` when it
did not. On error the ``summary`` and ``feedback`` keys are ``null`` - never
empty - so that "there is no feedback" and "the feedback was never fetched"
cannot be confused by a caller reading only stdout. Errors also exit non-zero.
"""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from typing import Any, NoReturn


# Bots whose comments are actionable code review. Matched against the login with
# any trailing "[bot]" stripped, because the REST API reports "coderabbitai[bot]"
# while GraphQL reports "coderabbitai" for the same account.
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
    r"(?i)^coderabbit",
    r"(?i)^gemini-code-assist",
    r"(?i)^greptile",
    r"(?i)^sourcery",
    r"(?i)^qodo",
    r"(?i)^codium",
    r"(?i)^ellipsis",
]

# Bots that post informational status reports rather than review findings.
INFO_BOT_PATTERNS = [
    r"(?i)^codecov",
    r"(?i)^dependabot",
    r"(?i)^renovate",
    r"(?i)^github-actions",
    r"(?i)^mergify",
    r"(?i)^semantic-release",
    r"(?i)^sonarcloud",
    r"(?i)^snyk",
    r"(?i)^netlify",
    r"(?i)^vercel",
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

BOT_SUFFIX = re.compile(r"\[bot\]$", re.IGNORECASE)


def fail(message: str) -> NoReturn:
    """Emit a machine-readable failure and exit non-zero.

    ``summary`` and ``feedback`` are explicitly null so a caller that reads only
    stdout cannot mistake a failed lookup for a clean PR.
    """
    print(
        json.dumps(
            {
                "status": "error",
                "error": message,
                "pr": None,
                "summary": None,
                "feedback": None,
                "action_required": (
                    "STOP: PR feedback could not be fetched, so this PR's review "
                    "state is unknown. Do not treat this as 'no feedback' and do "
                    f"not merge. Cause: {message}"
                ),
            },
            indent=2,
        )
    )
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


def run_gh_json_list(args: list[str]) -> list[dict[str, Any]]:
    """Fetch a paginated REST list and flatten it, or raise.

    Every caller of this needs a list. ``gh`` can exit 0 while writing nothing -
    a truncated response, a proxy, a degraded endpoint - and the shape that
    reaches us is then ``None``, not ``[]``. Returning ``[]`` for it is the
    silent-empty inversion this whole script exists to prevent: a whole feedback
    channel disappears and the run reports a clean PR. So anything that is not a
    list is an error, not an absence.
    """
    result = run_gh_json(args)
    if result is None:
        raise RuntimeError(
            f"gh {' '.join(args)} exited 0 but returned no output; "
            "the response was not delivered, so the feedback it carries is unknown"
        )
    if not isinstance(result, list):
        raise RuntimeError(
            f"gh {' '.join(args)} returned {type(result).__name__}, expected a list of pages"
        )
    return flatten_pages(result)


def get_repo_info() -> tuple[str, str]:
    result = run_gh_json(["repo", "view", "--json", "owner,name"])
    if not isinstance(result, dict):
        raise RuntimeError("could not determine repository from current directory")
    owner = result.get("owner", {}).get("login")
    repo = result.get("name")
    if not owner or not repo:
        raise RuntimeError("could not determine repository owner/name")
    return str(owner), str(repo)


def get_current_user() -> str:
    """Login of the account ``gh`` is authenticated as.

    Required, not optional: without it the fetcher cannot tell the caller's own
    replies apart from reviewer feedback, and reports them back as new work.
    """
    result = run_gh_json(["api", "user"])
    if not isinstance(result, dict):
        raise RuntimeError("could not determine the authenticated gh user")
    login = result.get("login")
    if not login:
        raise RuntimeError("authenticated gh user has no login")
    return str(login)


def get_pr_info(pr_number: int | None = None) -> dict[str, Any]:
    args = ["pr", "view", "--json", "number,url,headRefName,author,reviewDecision"]
    if pr_number is not None:
        args.insert(2, str(pr_number))
    result = run_gh_json(args)
    if not isinstance(result, dict):
        raise RuntimeError("unable to determine PR for current branch")
    return result


def normalize_login(username: str) -> str:
    """Strip the REST-only ``[bot]`` suffix so both APIs classify identically."""
    return BOT_SUFFIX.sub("", username or "").strip()


def is_review_bot(username: str) -> bool:
    login = normalize_login(username)
    return any(re.search(p, login) for p in REVIEW_BOT_PATTERNS)


def is_info_bot(username: str) -> bool:
    return any(re.search(p, username or "") for p in INFO_BOT_PATTERNS)


def flatten_pages(result: list[Any]) -> list[dict[str, Any]]:
    """Flatten the list of pages ``gh api --paginate --slurp`` returns.

    Shape validation belongs to ``run_gh_json_list``; by the time we get here the
    payload is known to be a list of pages.
    """
    entries: list[dict[str, Any]] = []
    for page in result:
        if isinstance(page, list):
            entries.extend(entry for entry in page if isinstance(entry, dict))
        elif isinstance(page, dict):
            entries.append(page)
    return entries


def get_issue_comments(owner: str, repo: str, pr_number: int) -> list[dict[str, Any]]:
    return run_gh_json_list([
        "api",
        f"repos/{owner}/{repo}/issues/{pr_number}/comments",
        "--paginate",
        "--slurp",
    ])


def get_reviews(owner: str, repo: str, pr_number: int) -> list[dict[str, Any]]:
    """Submitted reviews, read from REST because only REST carries ``html_url``.

    ``gh pr view --json reviews`` exposes the review's GraphQL node id, while the
    page anchor is keyed on the numeric database id, so a review body sourced
    that way reaches the caller with no link back to the review it came from.
    """
    return run_gh_json_list([
        "api",
        f"repos/{owner}/{repo}/pulls/{pr_number}/reviews",
        "--paginate",
        "--slurp",
    ])


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
              firstComment: comments(first: 1) {
                nodes {
                  body
                  createdAt
                  author {
                    login
                  }
                }
              }
              lastComment: comments(last: 1) {
                nodes {
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
            raise RuntimeError("gh api graphql returned no review-thread payload")

        pull_request = ((result.get("data") or {}).get("repository") or {}).get("pullRequest")
        if not isinstance(pull_request, dict):
            raise RuntimeError(
                f"GraphQL returned no pull request {pr_number} for {owner}/{repo}"
            )

        review_threads = pull_request.get("reviewThreads") or {}
        nodes = review_threads.get("nodes", [])
        if isinstance(nodes, list):
            threads.extend(nodes)

        page_info = review_threads.get("pageInfo") or {}
        if not page_info.get("hasNextPage"):
            break
        cursor = page_info.get("endCursor")
        if not cursor:
            break

    return threads


# Channel name -> the GraphQL connection on PullRequest that counts the same records.
PROBE_CONNECTIONS = {
    "review threads": "reviewThreads",
    "PR conversation comments": "comments",
    "reviews": "reviews",
}


def probe_counts(owner: str, repo: str, pr_number: int) -> dict[str, int]:
    """Server-side record counts per feedback channel.

    One small query, used only to sanity-check an empty channel. It asks GitHub
    how many records it holds rather than trusting that an empty response means
    an empty PR.
    """
    selections = "\n".join(
        f"{alias}: {connection}(first: 1) {{ totalCount }}"
        for alias, connection in enumerate_probe_aliases()
    )
    query = f"""
    query($owner: String!, $repo: String!, $pr: Int!) {{
      repository(owner: $owner, name: $repo) {{
        pullRequest(number: $pr) {{
          {selections}
        }}
      }}
    }}
    """
    result = run_gh_json([
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
    ])
    if not isinstance(result, dict):
        raise RuntimeError("gh api graphql returned no count-probe payload")

    pull_request = ((result.get("data") or {}).get("repository") or {}).get("pullRequest")
    if not isinstance(pull_request, dict):
        raise RuntimeError(
            f"count probe found no pull request {pr_number} for {owner}/{repo}"
        )

    counts: dict[str, int] = {}
    for channel, (alias, _connection) in zip(PROBE_CONNECTIONS, enumerate_probe_aliases()):
        node = pull_request.get(alias)
        total = node.get("totalCount") if isinstance(node, dict) else None
        if not isinstance(total, int):
            raise RuntimeError(f"count probe returned no total for {channel}")
        counts[channel] = total
    return counts


def enumerate_probe_aliases() -> list[tuple[str, str]]:
    """Stable ``(alias, connection)`` pairs for the probe query.

    Aliases are positional (``c0``, ``c1``, ...) so the query never repeats a
    field name, which GraphQL rejects when the same connection is selected twice.
    """
    return [(f"c{index}", connection) for index, connection in enumerate(PROBE_CONNECTIONS.values())]


def assert_empty_channels_are_really_empty(
    owner: str,
    repo: str,
    pr_number: int,
    fetched: dict[str, list[Any]],
) -> None:
    """Fail loudly when a channel came back empty but GitHub says it has records.

    This is the backstop for the defect class this script is built around: every
    known silent-empty path is closed above, but a *new* one would again look
    exactly like a clean PR. Probing only when something is already empty keeps
    this to one extra query on the rare path and none on the common one.
    """
    empty = [channel for channel, records in fetched.items() if not records]
    if not empty:
        return

    counts = probe_counts(owner, repo, pr_number)
    missing = [
        f"{channel} (fetched 0, GitHub reports {counts[channel]})"
        for channel in empty
        if counts.get(channel, 0) > 0
    ]
    if missing:
        raise RuntimeError(
            "feedback channels came back empty while GitHub reports records in them: "
            + "; ".join(missing)
            + ". The fetch was incomplete, so this PR's review state is unknown"
        )


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


def first_node(container: Any) -> dict[str, Any]:
    """Return the first node of a GraphQL connection, or an empty dict."""
    nodes = container.get("nodes") if isinstance(container, dict) else None
    if isinstance(nodes, list) and nodes and isinstance(nodes[0], dict):
        return nodes[0]
    return {}


def author_login(node: Any) -> str:
    """``author.login`` of a GraphQL node, tolerating a null ``author``.

    A deleted account serializes as ``"author": null``, which an unguarded
    ``.get("author", {}).get("login")`` turns into an AttributeError.
    """
    if not isinstance(node, dict):
        return ""
    author = node.get("author")
    if not isinstance(author, dict):
        return ""
    return str(author.get("login") or "")


def user_login(node: Any) -> str:
    """``user.login`` of a REST node, tolerating a null ``user``."""
    if not isinstance(node, dict):
        return ""
    user = node.get("user")
    if not isinstance(user, dict):
        return ""
    return str(user.get("login") or "")


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
    replied: bool = False,
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
    if replied:
        item["replied"] = True

    return item


def main() -> None:
    parser = argparse.ArgumentParser(description="Fetch and categorize PR feedback")
    parser.add_argument("--pr", type=int, help="PR number (defaults to current branch PR)")
    args = parser.parse_args()

    try:
        owner, repo = get_repo_info()
        viewer = get_current_user()
        pr_info = get_pr_info(args.pr)
    except RuntimeError as error:
        fail(str(error))

    pr_number = pr_info["number"]
    pr_author = author_login(pr_info)
    review_decision = pr_info.get("reviewDecision", "")

    self_logins = {
        normalize_login(viewer).casefold(),
        normalize_login(pr_author).casefold(),
    } - {""}

    def is_self(author: str) -> bool:
        """What the caller wrote is not feedback for the caller."""
        return normalize_login(author).casefold() in self_logins

    feedback: dict[str, list[dict[str, Any]]] = {
        "high": [],
        "medium": [],
        "low": [],
        "bot": [],
        "resolved": [],
    }

    try:
        threads = get_review_threads(owner, repo, pr_number)
        issue_comments = get_issue_comments(owner, repo, pr_number)
        reviews = get_reviews(owner, repo, pr_number)
        assert_empty_channels_are_really_empty(
            owner,
            repo,
            pr_number,
            {
                "review threads": threads,
                "PR conversation comments": issue_comments,
                "reviews": reviews,
            },
        )
    except RuntimeError as error:
        fail(str(error))

    for thread in threads:
        if not isinstance(thread, dict):
            continue

        first_comment = first_node(thread.get("firstComment"))
        author = author_login(first_comment)
        body = first_comment.get("body", "")
        if not author or not body or is_self(author) or len(body.strip()) < 3:
            continue

        is_resolved = bool(thread.get("isResolved", False))
        thread_id = thread.get("id")

        # A thread whose newest comment is ours has already been answered. Without
        # this the caller replies, sees the same thread on the next poll, and
        # replies again forever.
        last_author = author_login(first_node(thread.get("lastComment")))
        replied = bool(last_author) and is_self(last_author)

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
            replied=replied and not is_resolved,
        )
        feedback[category].append(item)

    for comment in issue_comments:
        if not isinstance(comment, dict):
            continue
        author = user_login(comment)
        body = comment.get("body", "")
        if not author or not body or is_self(author) or len(body.strip()) < 3:
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

    # Review bodies. Bots that post their findings as one review (CodeRabbit,
    # Copilot, Gemini) submit with state COMMENTED, so filtering on
    # CHANGES_REQUESTED alone discards the entire finding list.
    #
    # A review body is never flagged `replied`. It carries no thread and no
    # resolve button, so the only stateless signal available is "we said
    # something later", which cannot tell answering a review apart from
    # happening to type after it - and a review that lands mid-round would then
    # be dismissed unread. Re-reporting an answered review is visible; dropping
    # an unread one is not.
    for review in reviews:
        if not isinstance(review, dict):
            continue
        author = user_login(review)
        body = review.get("body", "") or ""
        if not author or not body.strip() or is_self(author):
            continue
        state = review.get("state", "")
        changes_requested = state == "CHANGES_REQUESTED"
        category = "high" if changes_requested else categorize_comment(author, body)
        item = extract_feedback_item(
            body=body,
            author=author,
            created_at=str(review.get("submitted_at", "")),
            url=review.get("html_url"),
            review_bot=category in {"high", "medium", "low"} and is_review_bot(author),
        )
        item["type"] = "changes_requested" if changes_requested else "review"
        feedback[category].append(item)

    for bucket in feedback.values():
        bucket.sort(key=item_sort_key)
        for item in bucket:
            item.pop("_created_at", None)

    priority_buckets = ("high", "medium", "low")
    review_bot_count = sum(
        1
        for bucket in priority_buckets
        for item in feedback[bucket]
        if item.get("review_bot")
    )
    already_replied = sum(
        1
        for bucket in priority_buckets
        for item in feedback[bucket]
        if item.get("replied")
    )

    open_high = [item for item in feedback["high"] if not item.get("replied")]
    open_medium = [item for item in feedback["medium"] if not item.get("replied")]
    open_low = [item for item in feedback["low"] if not item.get("replied")]

    output = {
        "status": "ok",
        "viewer": viewer,
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
            "already_replied": already_replied,
            "needs_attention": len(open_high) + len(open_medium),
        },
        "feedback": feedback,
    }

    if open_high:
        output["action_required"] = "Address high-priority feedback before merge"
    elif open_medium:
        output["action_required"] = "Address medium-priority feedback"
    elif open_low:
        output["action_required"] = "Review low-priority suggestions - ask user which to address"
    else:
        output["action_required"] = None

    print(json.dumps(output, indent=2))


if __name__ == "__main__":
    main()
