#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries
  MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
  export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/model-view"
}

@test "mesheryctl model view displays usage instructions when no model name is provided" {
  run $MESHERYCTL_BIN model view
  
  assert_failure
  
  assert_output --partial "Error: model name isn't specified"
  assert_output --partial "Usage: mesheryctl model view [model-name]"
}

@test "mesheryctl model view displays an existing model" {
  run bash -c "printf '\n' | TERM=dumb $MESHERYCTL_BIN model view model-import_cli-e2e-test \
      | grep -Ev '^(  )?(created_at|updated_at|deleted_at):' \
      | col -b \
      | sed -E \
          -e 's/(id: )[0-9a-fA-F-]{36}[[:space:]]*$/\1<UUID>/' \
          -e 's/(components_count: )[0-9]+[[:space:]]*$/\1<NUMBER>/' "

  assert_success
  assert_output --partial "$(cat "$TESTDATA_DIR/exp_out_existing_model.txt")"
}

@test "mesheryctl model view handles non-existent models gracefully" {
  run $MESHERYCTL_BIN model view non-existent-model

  assert_success
  assert_output --partial "No model(s) found for the given name  non-existent-model"
}

@test "mesheryctl model view supports JSON output" {
  run bash -c "printf '\n' | TERM=dumb $MESHERYCTL_BIN model view model-import_cli-e2e-test -o json \
      | grep -Ev '^[[:space:]]*\"(created_at|updated_at|deleted_at)\"[[:space:]]*:' \
      | col -b \
      | sed -E \
          -e 's/(id\": )\"[0-9a-fA-F-]{36}\"[[:space:]]*(,?)$/\1<UUID>\2/' \
          -e 's/(\"components_count\": )[0-9]+[[:space:]]*(,?)$/\1<NUMBER>\2/' "

  assert_success
  assert_output --partial "$(cat "$TESTDATA_DIR/exp_out_json.txt")"
}

@test "mesheryctl model view supports YAML output" {
  run bash -c "printf '\n' | TERM=dumb $MESHERYCTL_BIN model view model-import_cli-e2e-test -o yaml \
      | grep -Ev '^(  )?(created_at|updated_at|deleted_at):' \
      | col -b \
      | sed -E \
          -e 's/(id: )[0-9a-fA-F-]{36}[[:space:]]*$/\1<UUID>/' \
          -e 's/(components_count: )[0-9]+[[:space:]]*$/\1<NUMBER>/' "

  assert_success
  assert_output --partial "$(cat "$TESTDATA_DIR/exp_out_yaml.txt")"
}