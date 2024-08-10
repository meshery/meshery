#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status.
set -u  # Treat unset variables as an error and exit immediately.
set -o pipefail  # Consider the exit status of a pipeline to be the exit status of the last command to fail (or zero if no command failed).

# Function to print usage
print_usage() {
  echo "Usage: helm snapshot -f <URI to helm tar.gz file>"
  echo "  -f    URI to the packaged Helm chart (tar.gz)"
}

# Check for required arguments
if [ "$#" -ne 2 ] || [ "$1" != "-f" ]; then
  print_usage
  exit 1
fi


CHART_URI=$2
TEMP_DIR=$(mktemp -d)
CHART_PATH="$TEMP_DIR/chart.tgz"
LOG_FILE="$TEMP_DIR/snapshot.log"

# Log function for better debugging
log() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# Ensure mesheryctl is installed
if ! command -v mesheryctl &> /dev/null; then
  log "mesheryctl could not be found. Please install it before using this plugin."
  exit 1
fi

# Download the provided Helm chart tarball
log "Downloading Helm chart from $CHART_URI..."
if ! curl -L "$CHART_URI" -o "$CHART_PATH"; then
  log "Failed to download Helm chart from $CHART_URI"
  exit 1
fi

# Validate the downloaded chart
if [ ! -f "$CHART_PATH" ]; then
  log "Downloaded file is not found: $CHART_PATH"
  exit 1
fi

log "Helm chart downloaded successfully."

# Convert Helm chart to Meshery Design
log "Converting Helm chart to Meshery Design..."
# Logic to convert Helm chart to Meshery Design 

# Create a visual snapshot
log "Creating visual snapshot..."
# Logic that Meshery Snapshots generate an image as a visual snapshot of Meshery designs

if [ -z "$SNAPSHOT_ID" ]; then
  log "Failed to create Meshery Snapshot."
  exit 1
fi

log "Meshery Snapshot created successfully. Snapshot ID: $SNAPSHOT_ID"

# Clean up
log "Cleaning up temporary files..."
rm -rf "$TEMP_DIR"

log "Operation completed successfully. Snapshot ID: $SNAPSHOT_ID"
