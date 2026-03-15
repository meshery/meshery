import logging
import sys
import json
from pathlib import Path
from typing import Optional

import httpx
from github import Github
from pydantic import BaseSettings, SecretStr

github_api = "https://api.github.com"


class Settings(BaseSettings):
    github_repository: str
    github_event_path: Path
    github_event_name: Optional[str] = None
    input_token: SecretStr
    input_deploy_url: str


def get_pr_details(repo, event):
    pull_request = event.get("pull_request")
    if pull_request:
        pr_number = event.get("number") or pull_request.get("number")
        if pr_number:
            return int(pr_number), pull_request.get("head", {}).get("sha")

    workflow_run = event.get("workflow_run")
    if workflow_run:
        head_sha = workflow_run.get("head_commit", {}).get("id")
        if not head_sha:
            logging.error("No head commit found in workflow_run payload")
            return None, None
        commit = repo.get_commit(sha=head_sha)
        pulls = commit.get_pulls()
        if pulls.totalCount > 0:
            # A commit can be in multiple PRs. We'll take the first one found.
            return pulls[0].number, pulls[0].head.sha
        logging.error(f"No PR found for hash: {head_sha}")
        return None, head_sha

    logging.error("Unsupported GitHub event payload")
    return None, None


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    settings = Settings()
    logging.info(f"Using config: {settings.json()}")
    g = Github(settings.input_token.get_secret_value())
    repo = g.get_repo(settings.github_repository)
    try:
        with settings.github_event_path.open() as event_file:
            event = json.load(event_file)
    except json.JSONDecodeError as e:
        logging.error(f"Error parsing event file: {e}")
        sys.exit(0)
    pr_number, pr_head_sha = get_pr_details(repo, event)
    if not pr_number:
        sys.exit(0)
    display_sha = pr_head_sha or "unknown"
    github_headers = {
        "Authorization": f"token {settings.input_token.get_secret_value()}"
    }
    url = f"{github_api}/repos/{settings.github_repository}/issues/{pr_number}/comments"
    logging.info(f"Using comments URL: {url}")
    response = httpx.post(
        url,
        headers=github_headers,
        json={
            "body": f"🚀 Preview for commit {display_sha} at: {settings.input_deploy_url}"
        },
    )
    if not (200 <= response.status_code <= 300):
        logging.error(f"Error posting comment: {response.text}")
        sys.exit(1)
    logging.info("Finished")