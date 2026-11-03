#!/usr/bin/env bats

load "clean_helper.sh"

setup_file() {
  remove_all "setup_file"
}

setup() {
  export BATS_LIB_PATH=${BATS_LIB_PATH:-"/usr/lib"}
  bats_load_library bats-support
  bats_load_library bats-assert
 
}

teardown() {
  remove_all "teardown"
}


@test "Given Homebrew when installing Meshery CLI then mesheryctl is installed" {
  run brew install mesheryctl
  assert_success
  assert_output --partial "mesheryctl"

  run mesheryctl version
  assert_success
  assert_output --partial "Client"
  assert_output --partial "Server"
}