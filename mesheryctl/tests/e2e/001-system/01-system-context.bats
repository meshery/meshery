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

@test "mesheryctl system context create" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"

   assert_success
   assert_output --partial "Added"
}

@test "mesheryctl system context create duplicate" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME"

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "name already exists"
}

@test "mesheryctl system context create empty" {
   run $MESHERYCTL_BIN system context create

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide a context name"
}

@test "mesheryctl system context create correct --url empty" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME --url

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "flag needs an argument"
}

@test "mesheryctl system context create correct --url correct" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME --url $CONTEXT_URL

   assert_success
   assert_output --partial "Added"
}

@test "mesheryctl system context create correct --url foo" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME --url foo

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "invalid URI"
}

@test "mesheryctl system context create correct --url correct --set" {
   CONTEXT_NAME_1="example-context1"
   CONTEXT_NAME_2="example-context2"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME_1
   run $MESHERYCTL_BIN system context create $CONTEXT_NAME_2 --url $CONTEXT_URL --set

   assert_success
   assert_output --partial "Added"
}

# bats test_tags=system:context
@test "mesheryctl system context view --all display valid content" {
   run $MESHERYCTL_BIN system context view --all
   assert_success

   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}

# bats test_tags=system:context
@test "mesheryctl system context view display valid content" {
   run $MESHERYCTL_BIN system context view
   assert_success
   
   assert_line --regexp "$CONTEXT_REGEXP_MATCH"
   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}

@test "mesheryctl system context view --context" {
   run $MESHERYCTL_BIN system context view  --context

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "flag needs an argument"
}

@test "mesheryctl system context view --context invalid" {
   run $MESHERYCTL_BIN system context view  --context foo

   assert_failure
   assert_output --partial "doesn't exist"
}

@test "mesheryctl system context view --context correct" {
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

@test "mesheryctl system context list" {
   run $MESHERYCTL_BIN system context list

   assert_success
   assert_output --partial "Available contexts"
}

@test "mesheryctl system context switch foo" {
   run $MESHERYCTL_BIN system context switch foo

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "context does not exist"
}

@test "mesheryctl system context switch" {
   run $MESHERYCTL_BIN system context switch

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide exactly one context name"
}

@test "mesheryctl system context switch name" {
   CONTEXT_NAME_1="example-context1"
   CONTEXT_NAME_2="example-context2"

   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME_1"
   run $MESHERYCTL_BIN system context create "$CONTEXT_NAME_2"

   run $MESHERYCTL_BIN system context switch $CONTEXT_NAME_2

   assert_success
   assert_output --partial "switched to context"
}

@test "mesheryctl system context delete" {
   run $MESHERYCTL_BIN system context delete

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "provide a context name to delete"
}

@test "mesheryctl system context delete foo" {
   run $MESHERYCTL_BIN system context delete foo

   assert_failure
   assert_output --partial "Error"
   assert_output --partial "no context name found"
}

@test "mesheryctl system context delete correct" {
   CONTEXT_NAME="example-context"

   run $MESHERYCTL_BIN system context create $CONTEXT_NAME
   run $MESHERYCTL_BIN system context delete $CONTEXT_NAME

   assert_success
   assert_output --partial "deleted context"
}