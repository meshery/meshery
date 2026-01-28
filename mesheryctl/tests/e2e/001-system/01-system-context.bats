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

@test "given a valid context-name is provided as an argument when running mesheryctl system context create then the context will be created" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"

   assert_success
   assert_output --partial "Added"
}

@test "given a already existing context-name is provided as an argument when running mesheryctl system context create then the duplication error message displays" {
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

@test "given no url as an argument when running mesheryctl system context create --url then the error message displays" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME --url

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "flag needs an argument"
}

@test "given a valid url as an argument when running mesheryctl system context create --url then the context will be created" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME --url $CONTEXT_URL

   assert_success
   assert_output --partial "Added"
}

@test "given a invalid url as an argument when running mesheryctl system context create --url then the error message displays" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME --url foo

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "invalid URI"
}

@test "given all required arguments with --set flag when running mesheryctl system context create --url --set then the new context will be created and set it as current context" {
   CONTEXT_NAME_1="example-context1"
   CONTEXT_NAME_2="example-context2"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME_1
   run $MESHERYCTL_BIN system context create $CONTEXT_NAME_2 --url $CONTEXT_URL --set

   assert_success
   assert_output --partial "Added"
}

# bats test_tags=system:context
@test "given an --all flag to view when running mesheryctl system context view displays all the context with it's details" {
   run $MESHERYCTL_BIN system context view --all
   assert_success

   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}

# bats test_tags=system:context
@test "running mesheryctl system context view displays the current context" {
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

   assert_failure
   assert_output --partial "does not exist"
}

@test "given a valid context-name as an argument when running mesheryctl system context view --context then it displays the detailed context" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"
   run $MESHERYCTL_BIN system context view  --context $CONTEXT_NAME

   assert_success
   assert_line --regexp "^Current Context:[[:space:]]+$CONTEXT_NAME"
   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}

@test "running mesheryctl system context list displays all the available context" {
   run $MESHERYCTL_BIN system context list

   assert_success
   assert_output --partial "Available contexts"
}

@test "given invalid context-name as an argument when running mesheryctl system context switch then the error message displays" {
   run $MESHERYCTL_BIN system context switch foo

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "context does not exist"
}

@test "given no context-name as an argument when running mesheryctl system context switch then the error message displays" {
   run $MESHERYCTL_BIN system context switch

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide exactly one context name"
}

@test "given a valid context-name as an argument when running mesheryctl system context switch then the current context will switch to specified context" {
   CONTEXT_NAME_1="example-context1"
   CONTEXT_NAME_2="example-context2"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME_1"
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME_2"

   run $MESHERYCTL_BIN system context switch $CONTEXT_NAME_2

   assert_success
   assert_output --partial "switched to context"
}

@test "given no context-name as an argument when running mesheryctl system context delete then the error message displays" {
   run $MESHERYCTL_BIN system context delete

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide a context name to delete"
}

@test "given an invalid context-name as an argument when running mesheryctl system context delete then the error message displays" {
   run $MESHERYCTL_BIN system context delete foo

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "no context name found"
}

@test "given a valid context-name as an argument when running mesheryctl system context delete then the specified context will be deleted" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME
   run $MESHERYCTL_BIN system context delete $CONTEXT_NAME

   assert_success
   assert_output --partial "deleted context"
}