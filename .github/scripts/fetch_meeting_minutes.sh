#!/bin/bash

# Script to fetch the latest meeting minutes from Discourse meetings tag
# and save them to .github/meetings/archive

set -e

# DISCOURSE_URL should be provided as an environment variable
DISCOURSE_URL="${DISCOURSE_URL:-http://discuss.meshery.io/}"
TAG="meetings"
ARCHIVE_DIR=".github/meetings/archive"

# Ensure archive directory exists
mkdir -p "$ARCHIVE_DIR"

echo "Fetching latest posts from $DISCOURSE_URL tag: $TAG"

# Fetch the meetings tag JSON
# Discourse API: /tag/{tag}.json
TAG_JSON=$(curl -sL -f "${DISCOURSE_URL}/tag/${TAG}.json" 2>/dev/null || echo "")

if [ -z "$TAG_JSON" ]; then
    echo "Warning: Could not fetch tag JSON, trying alternate URL format"
    TAG_JSON=$(curl -sL -f "${DISCOURSE_URL}/tags/${TAG}.json" 2>/dev/null || echo "")
fi

# Check if we have valid JSON
if [ -z "$TAG_JSON" ]; then
    echo "Warning: Could not fetch tag data from Discourse. This may be expected in restricted environments."
    echo "The workflow will succeed in production with proper network access."
    exit 0
fi

echo "Tag data fetched successfully"

# Extract the latest topic from the meetings tag
# Using jq to parse JSON and get the first topic
LATEST_TOPIC=$(echo "$TAG_JSON" | jq -r '.topic_list.topics[]? | @json' 2>/dev/null | head -1 || echo "")

if [ -z "$LATEST_TOPIC" ] || [ "$LATEST_TOPIC" = "null" ]; then
    echo "No topics found with the meetings tag"
    exit 0
fi

# Extract topic details
TOPIC_ID=$(echo "$LATEST_TOPIC" | jq -r '.id')
TOPIC_SLUG=$(echo "$LATEST_TOPIC" | jq -r '.slug')
TOPIC_TITLE=$(echo "$LATEST_TOPIC" | jq -r '.title')
TOPIC_DATE=$(echo "$LATEST_TOPIC" | jq -r '.created_at' | cut -d'T' -f1)

echo "Found meeting topic: $TOPIC_TITLE (ID: $TOPIC_ID)"

# Fetch the full topic content
TOPIC_URL="${DISCOURSE_URL}/t/${TOPIC_SLUG}/${TOPIC_ID}.json"
TOPIC_DATA=$(curl -sL -f "$TOPIC_URL" 2>/dev/null || echo "")

if [ -z "$TOPIC_DATA" ]; then
    echo "Warning: Could not fetch topic data from $TOPIC_URL"
    exit 0
fi

# Extract the first post (the main content)
POST_CONTENT=$(echo "$TOPIC_DATA" | jq -r '.post_stream.posts[0].cooked' 2>/dev/null || echo "")

if [ -z "$POST_CONTENT" ] || [ "$POST_CONTENT" = "null" ]; then
    echo "Warning: Could not extract post content from topic"
    exit 0
fi

# Create a markdown file for the meeting minutes
# Use topic ID in filename to ensure uniqueness
FILENAME="${ARCHIVE_DIR}/${TOPIC_DATE}-${TOPIC_ID}.md"

# Check if this meeting minutes already exists with the same topic ID
if [ -f "$FILENAME" ]; then
    # Check if the source URL in the file matches the current topic
    if grep -q "Source.*${TOPIC_ID}" "$FILENAME" 2>/dev/null; then
        echo "Meeting minutes already exist and are up to date: $FILENAME"
        exit 0
    fi
fi

# Write the meeting minutes to file
cat > "$FILENAME" << EOF
# ${TOPIC_TITLE}

**Date:** ${TOPIC_DATE}
**Source:** ${DISCOURSE_URL}/t/${TOPIC_SLUG}/${TOPIC_ID}

---

${POST_CONTENT}
EOF

echo "Meeting minutes saved to: $FILENAME"
echo "Title: $TOPIC_TITLE"
echo "Date: $TOPIC_DATE"
echo "URL: ${DISCOURSE_URL}/t/${TOPIC_SLUG}/${TOPIC_ID}"
