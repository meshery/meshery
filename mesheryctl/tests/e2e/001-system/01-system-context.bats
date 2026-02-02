#!/usr/bin/env bats

setup() {
   export ORIGINAL_HOME="$HOME"
   export HOME="$(mktemp -d)"

   load "$E2E_HELPERS_PATH/bats_libraries"
   load "$E2E_HELPERS_PATH/constants"
	_load_bats_libraries

   CONTEXT_URL="http://localhost:9081"

   ENDPOINT_REGEX_MATCH='^[[:space:]]*endpoint:[[:space:]](http|https)://.*:[[:digit:]]+$'
   TOKEN_REGEX_MATCH='^[[:space:]]*token:[[:space:]][[:alnum:]]+$'
   PLATFORM_REGEX_MATCH='^[[:space:]]*platform:[[:space:]](kubernetes|docker)+$'
   PROVIDER_REGEX_MATCH='^[[:space:]]*provider:[[:space:]][[:alnum:]]+$'
   CONTEXT_REGEXP_MATCH='^Current[[:space:]]Context:[[:space:]]+[[:alnum:]_-]+$'
}

teardown() {
   rm -rf "$HOME"
   export HOME="$ORIGINAL_HOME"
}

@test "given a valid context-name is provided as an argument when running mesheryctl system context create then the context is created" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"

   assert_success
   assert_output --partial "Added"
}

@test "given an already existing context-name is provided as an argument when running mesheryctl system context create an error message is displayed" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"
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
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME --url $CONTEXT_URL

   assert_success
   assert_output --partial "Added"
}

@test "given a invalid url as an argument when running mesheryctl system context create --url invalid-url then an error message is displayed" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME --url foo

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "invalid URI"
}

@test "given all requirements met with --set flag when running mesheryctl system context create context-name --url valid-url --set then the new context is be created and set it as current context" {
   CONTEXT_NAME_1="example-context1"
   CONTEXT_NAME_2="example-context2"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME_1
   run $MESHERYCTL_BIN system context create $CONTEXT_NAME_2 --url $CONTEXT_URL --set

   assert_success
   assert_output --partial "Added"

   run $MESHERYCTL_BIN system context view
   assert_success
   assert_line --regexp "Current Context:[[:space:]]+$CONTEXT_NAME_2"
}

@test "given --all flag provided when running mesheryctl system context view --all then all contexts details are displayed" {
   run $MESHERYCTL_BIN system context view --all
   assert_success

   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}

@test "given all requirements met when running mesheryctl system context view then the details of current context is displayed" {
   run $MESHERYCTL_BIN system context view
   assert_success
   
   assert_line --regexp "$CONTEXT_REGEXP_MATCH"
   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}

@test "given no context-name as an argument when running mesheryctl system context view --context then the error message displays" {
   run $MESHERYCTL_BIN system context view  --context

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "flag needs an argument"
}

@test "given an invalid context-name as an argument when running mesheryctl system context view --context then the error message displays" {
   run $MESHERYCTL_BIN system context view  --context foo

   assert_success
   assert_output --partial "doesn't exist"
}

@test "given a valid context-name as an argument when running mesheryctl system context view --context then it displays the detailed context" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"
   run $MESHERYCTL_BIN system context view --context $CONTEXT_NAME

   assert_success
   assert_line --regexp "^Current Context:[[:space:]]+$CONTEXT_NAME"
   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}

@test "given all requirements met when running mesheryctl system context list then the available contexts are displayed" {
   run $MESHERYCTL_BIN system context list

   assert_success
   assert_output --partial "Available contexts"
}

@test "given invalid context-name provided when running mesheryctl system context switch then an error message is displayed" {
   run $MESHERYCTL_BIN system context switch foo

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
   CONTEXT_NAME_1="example-context1"
   CONTEXT_NAME_2="example-context2"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME_1"
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME_2"

   run $MESHERYCTL_BIN system context switch $CONTEXT_NAME_2

   assert_success
   assert_output --partial "switched to context"

   run $MESHERYCTL_BIN system context view
   assert_success
   assert_line --regexp "Current Context:[[:space:]]+$CONTEXT_NAME_2"
}

@test "given no context-name provided when running mesheryctl system context delete then an error message displayed" {
   run $MESHERYCTL_BIN system context delete

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide a context name to delete"
}

@test "given an invalid context-name provided when running mesheryctl system context delete invalid-context-name then an error message displayed" {
   run $MESHERYCTL_BIN system context delete foo

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "no context name found"
}

@test "given a valid context-name provided when running mesheryctl system context delete context-name then the specified context is deleted" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME
   run $MESHERYCTL_BIN system context delete $CONTEXT_NAME

   assert_success
   assert_output --partial "deleted context"

   run $MESHERYCTL_BIN system context view --context $CONTEXT_NAME
   assert_success
   assert_output --partial "doesn't exist"
}