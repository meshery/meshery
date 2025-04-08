#!/usr/bin/env bats

load '../helpers/bats-support/load'
load '../helpers/bats-assert/load'

# Use absolute path for testdata directory
TESTDATA_DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")/../testdata/design-import" && pwd 2>/dev/null || echo '/tmp')"

setup() {
  # Check if Meshery server is running
  if ! curl -s "http://localhost:9081/api/system/version" > /dev/null 2>&1; then
    echo "WARNING: Meshery server is not running, tests may fail"
  fi
  
  # Create testdata directory if it doesn't exist
  mkdir -p "$TESTDATA_DIR"
}

@test "mesheryctl design import successfully imports a design" {
  # Reference the YAML file from testdata directory (show absolute path for debugging)
  DESIGN_FILE="${TESTDATA_DIR}/nginx.yaml"
  echo "Using design file: $DESIGN_FILE"
  
  # Verify file exists
  [ -f "$DESIGN_FILE" ] || (echo "File doesn't exist: $DESIGN_FILE" && exit 1)
  
  # Import the design
  run $MESHERYCTL_BIN design import -f "$DESIGN_FILE" --source-type "Kubernetes Manifest"
  echo "Import command output: $output"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  # Skip if it's a file not found error
  if [[ "$output" =~ "no such file or directory" ]]; then
    skip "Test file not found, check path: $DESIGN_FILE"
  fi
  
  # Verify command succeeded
  assert_success
  assert_output --partial "imported" || assert_output --partial "Design ID" || assert_output --partial "saved"
  
  # Extract design ID from output - adjust regex to match output format which seems to be a shorter ID
  DESIGN_ID=$(echo "$output" | grep -o '[0-9a-f]\{8\}')
  
  # If not found with the short format, try the full UUID format
  if [ -z "$DESIGN_ID" ]; then
    DESIGN_ID=$(echo "$output" | grep -o '[0-9a-f]\{8\}-[0-9a-f]\{4\}-[0-9a-f]\{4\}-[0-9a-f]\{4\}-[0-9a-f]\{12\}')
  fi
  
  # Ensure we got a valid design ID
  [ -n "$DESIGN_ID" ]
  echo "Captured design ID: $DESIGN_ID"
  
  # Instead of writing to TEMP_DIR, just output the ID
  echo "Design ID for subsequent tests: $DESIGN_ID"
  
  # Verify the design exists using view command
  run $MESHERYCTL_BIN design view "$DESIGN_ID"
  assert_success
  assert_output --partial "$DESIGN_ID"
}

@test "mesheryctl design import for invalid file should fail" {
  # Create a path to a file that definitely doesn't exist
  NONEXISTENT_FILE="${TEMP_DIR}/definitely_nonexistent_$(date +%s).yaml"
  
  # Try to import a non-existent file with full path
  run $MESHERYCTL_BIN design import -f "$NONEXISTENT_FILE" --source-type "Kubernetes Manifest"
  echo "Import command output: $output"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  # Verify command failed
  assert_output --partial "Error" || assert_output --partial "no such file"
}
