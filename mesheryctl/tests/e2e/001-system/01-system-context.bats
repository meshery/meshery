#!/usr/bin/env bats

setup() {
    load "$SUPPORT_DESTDIR"
    load "$ASSERT_DESTDIR"

    ENDPOINT_REGEX_MATCH='^[[:space:]]*endpoint:[[:space:]](http|https)://.*:[[:digit:]]+$'
    TOKEN_REGEX_MATCH='^[[:space:]]*token:[[:space:]][a-zA-Z]+$'
    PLATFORM_REGEX_MATCH='^[[:space:]]*platform:[[:space:]][a-zA-Z]+$'
    PROVIDER_REGEX_MATCH='^[[:space:]]*provider:[[:space:]][a-zA-Z]+$'
    CONTEXT_REGEXP_MATCH='^Current[[:space:]]Context:[[:space:]][a-zA-Z]+$'
}

# bats test_tags=system:context
@test "mesheryctl system context view --all is succeeded" {
   run $MESHERYCTL_BIN system context view --all
   [ "$status" -eq 0 ]
   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}

# bats test_tags=system:context
@test "mesheryctl system context view is succedeed" {
   run $MESHERYCTL_BIN system context view
   [ "$status" -eq 0 ]
   assert_line --regexp "$CONTEXT_REGEXP_MATCH"
   assert_line --regexp "$ENDPOINT_REGEX_MATCH"
   assert_line --regexp "$TOKEN_REGEX_MATCH"
   assert_line --regexp "$PLATFORM_REGEX_MATCH"
   assert_line --regexp "$PROVIDER_REGEX_MATCH"
}
