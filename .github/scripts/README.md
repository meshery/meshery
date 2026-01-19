# Hacktoberfest Badge Automation

Automated system for awarding Hacktoberfest badges to Meshery contributors with merged PRs.

## Overview

This Python script automates the manual process of:
1. Monitoring Hacktoberfest-related GitHub issues
2. Detecting when associated PRs are merged
3. Extracting contributor emails from commit `Signed-off-by:` lines
4. Posting `/award-badge` commands to Slack's `#event` channel
5. Preventing duplicate badge awards

## Requirements

- Python 3.7+
- GitHub Personal Access Token (for production use)
- Slack Bot Token (for production use)

## Quick Start

### Installation

```bash
# Navigate to scripts directory
cd .github/scripts

# Install dependencies
pip install -r requirements.txt
```

### Usage

```bash
# Process specific issues (dry-run mode, no tokens needed)
python hacktoberfest-badge-automation.py --issues 16100,16101 --dry-run

# Process issues from file
python hacktoberfest-badge-automation.py --issues-file issues.txt --dry-run

# Process all issues with a label
python hacktoberfest-badge-automation.py --label hacktoberfest --dry-run
```

## Configuration

### For Testing (No Tokens Required)

Use `--dry-run` flag to test without actual API calls:

```bash
python hacktoberfest-badge-automation.py --issues 16187 --dry-run
```

### For Production Use

1. **Copy environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Set up GitHub Token:**
   - Go to https://github.com/settings/tokens/new
   - Name: "Meshery Hacktoberfest Automation"
   - Scope: `repo` (Full control of private repositories)
   - Generate token and add to `.env`

3. **Set up Slack Bot Token:**
   - Go to https://api.slack.com/apps
   - Create new app: "Hacktoberfest Badge Bot"
   - Add scopes: `chat:write`, `chat:write.public`
   - Install to workspace
   - Copy Bot User OAuth Token to `.env`

4. **Run with tokens:**
   ```bash
   python hacktoberfest-badge-automation.py --issues 16100 --log-level INFO
   ```

## How It Works

### 1. Issue Processing
- Fetches issue details from GitHub API
- Searches issue timeline for linked PRs
- Parses issue body for PR references (e.g., `Fixes #16100`)

### 2. PR Validation
- Checks if PR is merged
- Skips unmerged PRs automatically

### 3. Email Extraction
- Retrieves all commits from merged PR
- Scans commit messages for `Signed-off-by:` lines
- Extracts email addresses: `Signed-off-by: Name <email@example.com>`
- Validates email format

### 4. Badge Award
- Posts to Slack: `/award-badge email@example.com hacktoberfest25`
- Tracks awarded emails in `awarded_emails.json`
- Skips duplicate awards

### 5. Logging
- Console output with colored formatting
- Log file: `hacktoberfest_automation_YYYYMMDD_HHMMSS.log`
- Summary statistics at completion

## Example Output

```
2024-01-18 10:30:15 - INFO - ================================================================================
2024-01-18 10:30:15 - INFO - Hacktoberfest Badge Automation Script
2024-01-18 10:30:15 - INFO - ================================================================================
2024-01-18 10:30:15 - WARNING - ⚠️  Running in DRY RUN mode - no badges will be awarded
2024-01-18 10:30:15 - INFO - Starting Hacktoberfest badge automation for 2 issue(s)
2024-01-18 10:30:15 - INFO - Configuration: Repo=meshery/meshery, Badge=hacktoberfest25, DryRun=True
2024-01-18 10:30:15 - INFO - Processing issue #16100...
2024-01-18 10:30:16 - INFO - Found 1 PR(s) for issue #16100: [12345]
2024-01-18 10:30:17 - INFO - PR #12345 is merged ✓
2024-01-18 10:30:18 - INFO - Extracted 1 email(s) from PR #12345: {'contributor@example.com'}
2024-01-18 10:30:19 - INFO - [DRY RUN] Would post to #event: /award-badge contributor@example.com hacktoberfest25
2024-01-18 10:30:20 - INFO - ================================================================================
2024-01-18 10:30:20 - INFO - EXECUTION SUMMARY
2024-01-18 10:30:20 - INFO - ================================================================================
2024-01-18 10:30:20 - INFO - Issues processed: 2
2024-01-18 10:30:20 - INFO - PRs found: 3
2024-01-18 10:30:20 - INFO - Merged PRs: 2
2024-01-18 10:30:20 - INFO - Emails extracted: 4
2024-01-18 10:30:20 - INFO - Badges awarded: 0 (dry run)
2024-01-18 10:30:20 - INFO - Duplicates skipped: 0
2024-01-18 10:30:20 - INFO - Errors: 0
2024-01-18 10:30:20 - INFO - ================================================================================
2024-01-18 10:30:20 - INFO - ✓ Script execution completed
```

## GitHub Actions Integration

The included workflow file (`.github/workflows/hacktoberfest-badge-award.yml`) allows:
- Scheduled daily execution
- Manual triggering with custom parameters
- Automatic log archiving

**Setup for CI/CD:**
1. Add secrets in repository settings:
   - `GH_TOKEN`: GitHub Personal Access Token
   - `SLACK_BOT_TOKEN`: Slack Bot Token
2. Workflow runs automatically or on-demand

## Troubleshooting

### No output when running script
Ensure you're using `--dry-run` flag for testing without tokens:
```bash
python hacktoberfest-badge-automation.py --issues 16187 --dry-run
```

### "Missing required environment variables"
- Use `--dry-run` flag to skip token validation
- Or set up `.env` file with valid tokens

### "No PRs found for issue"
- Issue may not have linked PRs yet
- Check if issue body contains PR reference
- PR must be linked via GitHub's reference system

### "Invalid email format found"
- Commit missing proper DCO sign-off
- Format must be: `Signed-off-by: Name <email@domain.com>`

## Security

- Never commit `.env` file to git
- Store tokens in GitHub Secrets for CI/CD
- Use dry-run mode for local testing
- Rotate tokens periodically

## Files Generated

- `hacktoberfest_automation_*.log` - Execution logs
- `awarded_emails.json` - Tracks awarded emails

## Contributing

Improvements welcome! Test with `--dry-run` first.

## License

Apache 2.0 (same as Meshery project)