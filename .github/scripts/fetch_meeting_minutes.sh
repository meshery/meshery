#!/bin/bash

# Script to fetch the latest meeting minutes from Discourse community category
# and save them to .github/meetings/archive

set -e

DISCOURSE_URL="https://discuss.layer5.io"
CATEGORY="community"
ARCHIVE_DIR=".github/meetings/archive"

# Ensure archive directory exists
mkdir -p "$ARCHIVE_DIR"

echo "Fetching latest posts from $DISCOURSE_URL category: $CATEGORY"

# Fetch the community category JSON
# Discourse API: /c/{category}/{id}.json or /c/{category}.json
CATEGORY_JSON=$(curl -sL "${DISCOURSE_URL}/c/${CATEGORY}.json" 2>/dev/null || echo "{}")

if [ -z "$CATEGORY_JSON" ] || [ "$CATEGORY_JSON" = "{}" ]; then
    echo "Warning: Could not fetch category JSON, trying alternate URL format"
    CATEGORY_JSON=$(curl -sL "${DISCOURSE_URL}/c/${CATEGORY}/5.json" 2>/dev/null || echo "{}")
fi

# Check if we have valid JSON
if [ -z "$CATEGORY_JSON" ] || [ "$CATEGORY_JSON" = "{}" ]; then
    echo "Error: Could not fetch category data from Discourse"
    exit 1
fi

echo "Category data fetched successfully"

# Extract the latest topic that contains "meeting" in the title (case insensitive)
# Using jq to parse JSON and find meeting-related topics
LATEST_TOPIC=$(echo "$CATEGORY_JSON" | jq -r '.topic_list.topics[] | select(.title | test("meeting"; "i")) | @json' | head -1)

if [ -z "$LATEST_TOPIC" ] || [ "$LATEST_TOPIC" = "null" ]; then
    echo "No meeting topics found in the community category"
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
TOPIC_DATA=$(curl -sL "$TOPIC_URL" 2>/dev/null || echo "{}")

if [ -z "$TOPIC_DATA" ] || [ "$TOPIC_DATA" = "{}" ]; then
    echo "Error: Could not fetch topic data"
    exit 1
fi

# Extract the first post (the main content)
POST_CONTENT=$(echo "$TOPIC_DATA" | jq -r '.post_stream.posts[0].cooked' 2>/dev/null || echo "")

if [ -z "$POST_CONTENT" ]; then
    echo "Error: Could not extract post content"
    exit 1
fi

# Create a markdown file for the meeting minutes
FILENAME="${ARCHIVE_DIR}/${TOPIC_DATE}-meeting.md"

# Check if this meeting minutes already exists
if [ -f "$FILENAME" ]; then
    EXISTING_CONTENT=$(cat "$FILENAME")
    if [ "$EXISTING_CONTENT" = "# ${TOPIC_TITLE}

**Date:** ${TOPIC_DATE}
**Source:** ${DISCOURSE_URL}/t/${TOPIC_SLUG}/${TOPIC_ID}

---

${POST_CONTENT}" ]; then
        echo "Meeting minutes already up to date: $FILENAME"
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
