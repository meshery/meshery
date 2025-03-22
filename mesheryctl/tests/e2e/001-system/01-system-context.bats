#!/usr/bin/env bats

setup() {
    load "$SUPPORT_DESTDIR"
    load "$ASSERT_DESTDIR"
}

# bats test_tags=system:context
@test "mesheryctl system context view --all is succeeded" {
   run $MESHERYCTL_BIN system context view --all
   [ "$status" -eq 0 ]
   assert_output --regexp 'endpoint'
   assert_output --regexp 'token'
   assert_output --regexp 'platform'
   assert_output --regexp 'provider'
}

# bats test_tags=system:context
@test "mesheryctl system context view is succedeed" {
   run $MESHERYCTL_BIN system context view
   [ "$status" -eq 0 ]
   assert_output --regexp 'Current Context'
   assert_output --regexp 'endpoint'
   assert_output --regexp 'token'
   assert_output --regexp 'platform'
   assert_output --regexp 'provider'
}
