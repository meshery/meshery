# Meeting Minutes Archive

This directory contains archived meeting minutes from the Meshery community.

## Overview

Meeting minutes are automatically fetched from the community discussion forum's "meetings" tag and committed to this repository daily via GitHub Actions.

## Structure

- `archive/` - Contains meeting minutes in Markdown format
  - Files are named using the pattern: `YYYY-MM-DD-meeting.md`
  - Each file includes:
    - Meeting title
    - Date
    - Link to the original Discourse post
    - Full meeting content

## Automation

The meeting minutes are updated automatically through the GitHub Actions workflow:
- **Workflow**: `.github/workflows/meeting-minutes-update.yml`
- **Script**: `.github/scripts/fetch_meeting_minutes.sh`
- **Schedule**: Daily at midnight UTC
- **Trigger**: Can also be manually triggered via workflow_dispatch

### Required Secrets

The workflow requires the following GitHub secrets to be configured:
- `DISCOURSE_URL`: The base URL of the Discourse forum (e.g., `https://discuss.meshery.io`)
- `MESHERY_CI`: GitHub token for committing changes

## Manual Updates

To manually fetch the latest meeting minutes:

```bash
export DISCOURSE_URL="https://<discussion-forum-URL>"
bash .github/scripts/fetch_meeting_minutes.sh
```

## Source

Meeting minutes are sourced from:
- Forum: https://discuss.meshery.io
- Tag: meetings
- URL: https://discuss.meshery.io/tag/meetings
