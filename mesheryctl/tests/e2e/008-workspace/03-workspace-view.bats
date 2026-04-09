#!/usr/bin/env bats

setup() {
load "$E2E_HELPERS_PATH/bats_libraries"
load "$E2E_HELPERS_PATH/constants"
_load_bats_libraries
}

@test "given no argument provided when running mesheryctl workspace view then an error message is displayed" {
run $MESHERYCTL_BIN workspace view

assert_failure
assert_output --partial "Error"
assert_output --partial "please provide exactly one workspace name or ID"
}

@test "given a name without --orgId when running mesheryctl workspace view foo then an error message is displayed" {
run $MESHERYCTL_BIN workspace view foo

assert_failure
assert_output --partial "Error"
assert_output --partial "--orgId is required when searching by name"
}

@test "given an invalid orgId when running mesheryctl workspace view somename --orgId foo then an error message is displayed" {
run $MESHERYCTL_BIN workspace view somename --orgId foo

assert_failure
assert_output --partial "Error"
assert_output --partial "Invalid value for --orgId 'foo': must be a valid UUID"
}

@test "given non-existent UUID when running mesheryctl workspace view then workspace not found error is displayed" {
run $MESHERYCTL_BIN workspace view 00000000-0000-0000-0000-000000000000

assert_failure
assert_output --partial "Error"
}

@test "given a valid workspace id when running mesheryctl workspace view then it should pass validation" {
FETCHED_ID=$($MESHERYCTL_BIN workspace list --orgId "2d2c0b60-076a-4f0a-8a63-de538570a553" | grep -E '[a-f0-9-]{36}' | awk '{print $1}' | head -n 1)

if [ -z "$FETCHED_ID" ]; then
FETCHED_ID="2d2c0b60-076a-4f0a-8a63-de538570a553"
fi

run $MESHERYCTL_BIN workspace view "$FETCHED_ID"

refute_output --partial "must be a valid UUID"

if [ "$status" -eq 0 ]; then
assert_output --partial "id: $FETCHED_ID"
else
assert_output --partial "Unable to connect"
fi
}