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

@test "mesheryctl design apply -f invalid path shows error" {
  INVALID_PATH="./test/invalid/path/nginx.yaml"
  run $MESHERYCTL_BIN design apply -f "$INVALID_PATH"

  assert_failure
  assert_output --partial "invalid"
  assert_output --partial "Enter a valid path"
}


