#!/usr/bin/env python3
"""
Simple automation to award Hacktoberfest badges when PRs referencing
configured issues are merged. It extracts Signed-off-by emails from PR commits
and posts `/award-badge <email> <badge>` messages to a Slack channel.

This script is intentionally conservative (has a dry-run mode) and expects the
following environment variables:

- GITHUB_TOKEN: personal access token with repo read access
- SLACK_TOKEN: Slack bot token with chat:write scope
- REPO: repository in owner/repo format (e.g. layer5io/meshery)
- ISSUE_LIST: comma-separated issue numbers to monitor (optional)
- HACKTOBERFEST_LABEL: label to filter issues if ISSUE_LIST not provided (optional, default: hacktoberfest)
- SLACK_CHANNEL: channel to post award commands (e.g. #event)
- BADGE_NAME: name of the badge to award (default: hacktoberfest25)

Usage examples:
  # Dry run for issues 15980 and 15981
  ISSUE_LIST=15980,15981 GITHUB_TOKEN=... SLACK_TOKEN=... REPO=owner/repo python3 hacktoberfest_award.py --dry-run

"""
import os
import re
import sys
import argparse
from typing import List, Set

import requests

GITHUB_API = "https://api.github.com"


def get_issues_from_label(owner: str, repo: str, label: str, token: str) -> List[int]:
    url = f"{GITHUB_API}/repos/{owner}/{repo}/issues"
    issues = []
    page = 1
    headers = {"Authorization": f"token {token}", "Accept": "application/vnd.github+json"}
    while True:
        params = {"labels": label, "state": "all", "per_page": 100, "page": page}
        r = requests.get(url, headers=headers, params=params)
        r.raise_for_status()
        data = r.json()
        if not data:
            break
        for it in data:
            if 'pull_request' in it:
                # skip PRs
                continue
            issues.append(int(it["number"]))
        page += 1
    return issues


def get_merged_prs_for_issue(owner: str, repo: str, issue_number: int, token: str) -> List[int]:
    # We'll list closed PRs and check whether their title or body references the issue number.
    url = f"{GITHUB_API}/repos/{owner}/{repo}/pulls"
    prs = []
    page = 1
    headers = {"Authorization": f"token {token}", "Accept": "application/vnd.github+json"}
    while True:
        params = {"state": "closed", "per_page": 100, "page": page}
        r = requests.get(url, headers=headers, params=params)
        r.raise_for_status()
        data = r.json()
        if not data:
            break
        for pr in data:
            if pr.get("merged_at") is None:
                continue
            # check body and title for reference like '#<issue_number>'
            if re.search(rf"#\s*{issue_number}\b", (pr.get("title") or "") + "\n" + (pr.get("body") or "")):
                prs.append(int(pr["number"]))
        page += 1
    return prs


def extract_signed_off_emails(owner: str, repo: str, pr_number: int, token: str) -> Set[str]:
    url = f"{GITHUB_API}/repos/{owner}/{repo}/pulls/{pr_number}/commits"
    headers = {"Authorization": f"token {token}", "Accept": "application/vnd.github+json"}
    emails = set()
    r = requests.get(url, headers=headers)
    r.raise_for_status()
    for c in r.json():
        msg = c.get("commit", {}).get("message", "")
        for m in re.findall(r"Signed-off-by:\s*(.+)", msg, flags=re.IGNORECASE):
            # extract email from the Signed-off-by line if present
            e = re.search(r"<([^>]+)>", m)
            if e:
                emails.add(e.group(1).strip())
            else:
                # fallback: try to find an email in the line
                e2 = re.search(r"([a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+)", m)
                if e2:
                    emails.add(e2.group(1).strip())
    return emails


def post_slack_award(slack_token: str, channel: str, email: str, badge: str, dry_run: bool = True):
    text = f"/award-badge {email} {badge}"
    if dry_run:
        print(f"DRY RUN: would post to {channel}: {text}")
        return True
    url = "https://slack.com/api/chat.postMessage"
    headers = {"Authorization": f"Bearer {slack_token}", "Content-Type": "application/json; charset=utf-8"}
    payload = {"channel": channel, "text": text}
    r = requests.post(url, headers=headers, json=payload)
    r.raise_for_status()
    data = r.json()
    if not data.get("ok"):
        print("Slack API error:", data)
        return False
    return True


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--dry-run", action="store_true", help="Do not post to Slack; just print actions")
    args = parser.parse_args()

    github_token = os.getenv("GITHUB_TOKEN")
    slack_token = os.getenv("SLACK_TOKEN")
    repo = os.getenv("REPO")
    issue_list_env = os.getenv("ISSUE_LIST")
    label = os.getenv("HACKTOBERFEST_LABEL", "hacktoberfest")
    slack_channel = os.getenv("SLACK_CHANNEL", "#event")
    badge_name = os.getenv("BADGE_NAME", "hacktoberfest25")

    if not github_token or not repo:
        print("GITHUB_TOKEN and REPO (owner/repo) are required")
        sys.exit(2)

    owner, repo_name = repo.split("/")

    if issue_list_env:
        issue_nums = [int(x.strip()) for x in issue_list_env.split(",") if x.strip()]
    else:
        issue_nums = get_issues_from_label(owner, repo_name, label, github_token)

    print(f"Processing issues: {issue_nums}")

    processed_emails = set()

    for issue in issue_nums:
        prs = get_merged_prs_for_issue(owner, repo_name, issue, github_token)
        print(f"Issue #{issue} -> merged PRs: {prs}")
        for pr in prs:
            emails = extract_signed_off_emails(owner, repo_name, pr, github_token)
            for e in emails:
                if e in processed_emails:
                    print(f"Already awarded {e}, skipping")
                    continue
                ok = post_slack_award(slack_token, slack_channel, e, badge_name, dry_run=args.dry_run or slack_token is None)
                if ok:
                    processed_emails.add(e)

    print("Done.")


if __name__ == "__main__":
    main()
