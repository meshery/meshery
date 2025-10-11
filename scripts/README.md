
Hacktoberfest Badge Award Automation
===================================

This script automates awarding Hacktoberfest badges by detecting merged PRs
that reference Hacktoberfest issues and posting `/award-badge <email> <badge>`
commands to a configured Slack channel.


Configuration
-------------

- GITHUB_TOKEN: GitHub token with read access to the repository
- SLACK_TOKEN: Slack bot token (optional for dry-run)
- REPO: owner/repo (defaults to repository where the workflow runs)
- ISSUE_LIST: comma-separated issue numbers to process (optional)
- HACKTOBERFEST_LABEL: label to search for issues (default: hacktoberfest)
- SLACK_CHANNEL: Slack channel to post award commands (e.g. #event)
- BADGE_NAME: Badge name to award (default: hacktoberfest25)


Usage
-----

Dry run:

    GITHUB_TOKEN=... REPO=owner/repo python3 hacktoberfest_award.py --dry-run


To run for real, provide SLACK_TOKEN and omit --dry-run.


Security
--------

Store tokens as GitHub Actions secrets (SLACK_TOKEN) and do not commit them in
the repository.
