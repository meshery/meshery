#!/usr/bin/env bats

declare FILE_TO_CLEANUP
declare COMPONENT_NAME

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  COMPONENT_REQUIRED_FIELDS='
  .id != null and
  .displayName != null and
  .description != null and
  .schemaVersion != null and
  .format != null and
  .version != null and
  has(\"configuration\") and
  .metadata != null and
  .model.id != null and
  .model != null and
  .component != null
  '
  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/component-search"
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/valid-model
  assert_success
  export TESTDATA_DIR="$BATS_TEST_DIRNAME/testdata/component-view"

  COMPONENT_NAME=$($MESHERYCTL_BIN component list --page 1 --pagesize 1 | tail -n 1 | awk '{print $1}')
}

teardown() {
  if [[ -f "$FILE_TO_CLEANUP" ]]; then
    rm -f "$FILE_TO_CLEANUP"
    FILE_TO_CLEANUP=""
  fi
}

test_component_view_format() {
    local format=$1
    local -A validation_tool
    validation_tool["json"]="jq"
    validation_tool["yaml"]="yq"

    if [[ -z "${validation_tool[$format]+_}" ]]; then
        echo "Unsupported format: $format"
        return 1
    fi

    FILE_TO_CLEANUP="${HOME}/.meshery/component_${COMPONENT_NAME}.${format}"

    printf '\n' | $MESHERYCTL_BIN component view "${COMPONENT_NAME}" -o "${format}" --save

    assert_file_exist "${FILE_TO_CLEANUP}"
    run bash -c "${validation_tool[$format]} -e \"$COMPONENT_REQUIRED_FIELDS\" \"${FILE_TO_CLEANUP}\""

    assert_success
    assert_output "true"
}

test_view_save() {
  local format=$1
  FILE_TO_CLEANUP="${HOME}/.meshery/component_${COMPONENT_NAME}.${format}"

  local expected_success_message="Output saved as ${format} in file: ${FILE_TO_CLEANUP}"
  run bash -c "printf '\n' | $MESHERYCTL_BIN component view \"${COMPONENT_NAME}\" -o ${format} --save"

  assert_success

  assert_output --partial "$expected_success_message"
  assert_file_exist "$FILE_TO_CLEANUP"
}

@test "view command fails with no arguments" {
  run $MESHERYCTL_BIN component view
  assert_failure
  assert_output --partial "Error: [component name] is required but not specified"
  assert_output --partial "Usage: mesheryctl component view [component-name]"
  assert_output --partial "Run 'mesheryctl component view --help' to see detailed help message"
}

@test "view command fails with too many arguments" {
  run $MESHERYCTL_BIN component view comp1 comp2
  assert_failure
  assert_output --partial "Error: too many arguments specified"
  assert_output --partial "Usage: mesheryctl component view [component-name]"
  assert_output --partial "Run 'mesheryctl component view --help' to see detailed help message"
}

@test "view command fails with an invalid output format" {
  run $MESHERYCTL_BIN component view some-component -o xml

  assert_failure
  assert_output --partial 'Error: output-format "xml" is invalid. Available options [json|yaml]'
  assert_output --partial 'See https://docs.meshery.io/reference/mesheryctl/exp/components/view for usage details'
}

@test "view command displays JSON output for a known component" {
    test_component_view_format "json"
}

@test "view command displays YAML output for a known component" {
    test_component_view_format "yaml"
}

@test "view command saves JSON output file" {
  test_view_save "json"
}

@test "view command saves YAML output file" {
  test_view_save "yaml"
}