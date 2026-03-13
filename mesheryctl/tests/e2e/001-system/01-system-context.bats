#!/usr/bin/env bats

setup() {   
   load "$E2E_HELPERS_PATH/bats_libraries"
   load "$E2E_HELPERS_PATH/constants"
	_load_bats_libraries

   CONTEXT_URL="http://localhost:9081"
   CONTEXT_NAME="example-context"
   CONTEXT_NAME_2="example-context2"
}

@test "given a valid context-name is provided as an argument when running mesheryctl system context create then the context is created" {
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"

   assert_success
   assert_output --partial "Added"
}

@test "given an already existing context-name is provided as an argument when running mesheryctl system context create then an error message is displayed" {
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "name already exists"
}

@test "given no context-name as an argument when running mesheryctl system context create then the error message displays" {
   run $MESHERYCTL_BIN system context create

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide a context name"
}

@test "given a valid url as an argument when running mesheryctl system context create --url invalid-url then an the context is displayed" {
   run $MESHERYCTL_BIN system context delete "$CONTEXT_NAME"
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME" --url "$CONTEXT_URL"

   assert_success
   assert_output --partial "Added"
}

@test "given an invalid url as an argument when running mesheryctl system context create --url invalid-url then an error message is displayed" {
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME" --url invalid

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "invalid URI"
}

@test "given all requirements met with --set flag when running mesheryctl system context create context-name --url valid-url --set then the new context is created and set it as current context" {
   skip "Temporarily skipping"
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME_2" --url "$CONTEXT_URL" --set

   assert_success
   assert_output --partial "Added"

   run $MESHERYCTL_BIN system context view
   assert_success
   assert_line --regexp "Current Context:[[:space:]]+$CONTEXT_NAME_2"
}

@test "given --all flag provided when running mesheryctl system context view --all then all contexts details are displayed" {
   run $MESHERYCTL_BIN system context view --all
   assert_success

   assert_output --partial "endpoint"
   assert_output --partial "token"
}

@test "given all requirements met when running mesheryctl system context view then the details of current context is displayed" {
   run $MESHERYCTL_BIN system context view
   assert_success
   
   assert_output --partial "endpoint"
   assert_output --partial "token"
}

@test "given an invalid context-name as an argument when running mesheryctl system context view --context then the error message displays" {
   run $MESHERYCTL_BIN system context view --context invalid

   assert_failure
   assert_output --partial "does not exist"
}

@test "given a valid context-name as an argument when running mesheryctl system context view --context then it displays the detailed context" {
   run $MESHERYCTL_BIN system context view --context "$CONTEXT_NAME"

   assert_success
    assert_output --partial "endpoint"
   assert_output --partial "token"
}

@test "given all requirements met when running mesheryctl system context list then the available contexts are displayed" {
   run $MESHERYCTL_BIN system context list

   assert_success
   assert_output --partial "Available contexts"
}

@test "given an invalid context-name provided when running mesheryctl system context switch then an error message is displayed" {
   run $MESHERYCTL_BIN system context switch invalid

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "context does not exist"
}

@test "given no context-name provided when running mesheryctl system context switch then an error message is displayed" {
   run $MESHERYCTL_BIN system context switch

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide exactly one context name"
}

@test "given a valid context-name provided when running mesheryctl system context switch context-name then the current context is switched to specified context" {
   skip "Temporarily skipping"
   run $MESHERYCTL_BIN system context switch "$CONTEXT_NAME_2"

   assert_success
   assert_output --partial "switched to context"

   run $MESHERYCTL_BIN system context view
   assert_success
   assert_line --regexp "Current Context:[[:space:]]+$CONTEXT_NAME_2"
}

@test "given no context-name provided when running mesheryctl system context delete then an error message is displayed" {
   run $MESHERYCTL_BIN system context delete

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide a context name to delete"
}

@test "given an invalid context-name provided when running mesheryctl system context delete invalid-context-name then an error message is displayed" {
   run $MESHERYCTL_BIN system context delete invalid

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "no context name found"
}

@test "given a valid context-name provided when running mesheryctl system context delete context-name then the specified context is deleted" {
   run $MESHERYCTL_BIN system context delete "$CONTEXT_NAME"

   assert_success
   assert_output --partial "deleted context"

   run $MESHERYCTL_BIN system context view --context "$CONTEXT_NAME"
   assert_failure
   assert_output --partial "does not exist"
}