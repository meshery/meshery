#!/usr/bin/env bats

# Setup function to load libraries and prepare test fixtures
setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  mkdir -p "$TEMP_DATA_DIR/design"
  cp "$E2E_TESTDATA_PATH/design-import/nginx.yaml" "$TEMP_DATA_DIR/design/nginx.yaml"
  export TESTDATA_PATH="$TEMP_DATA_DIR/design/nginx.yaml"
}

@test "mesheryctl design apply applies design file" {
  if [ ! -f "$TESTDATA_PATH" ]; then
    skip "Test fixture $TESTDATA_PATH not found"
  fi

  run $MESHERYCTL_BIN design apply -f "$TESTDATA_PATH"

  assert_success
}

@test "mesheryctl design apply design-name applies an imported design" {
  if [ ! -f "$TESTDATA_PATH" ]; then
    skip "Test fixture $TESTDATA_PATH not found"
  fi

  # Import the design file
  run $MESHERYCTL_BIN design import -f "$TESTDATA_PATH" --source-type "Kubernetes Manifest"
  assert_success
  assert_output --partial "imported"

  # Derive design name from file name (nginx.yaml- nginx)
  DESIGN_NAME=$(basename "$TESTDATA_PATH" .yaml)

  # Apply the design using its name
  run $MESHERYCTL_BIN  design apply "$DESIGN_NAME"

  # Assert apply succeeded
  assert_success
  assert_output --partial "applied" || assert_output --partial "success"
}

@test "mesheryctl design apply -f invalid path shows error" {
  INVALID_PATH="./test/invalid/path/nginx.yaml"
  run $MESHERYCTL_BIN design apply -f "$INVALID_PATH"

  assert_failure
  assert_output --partial "invalid"
  assert_output --partial "Enter a valid path"
}


