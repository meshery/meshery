#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

}

@test "Given wrong platforms provided when installing then installer exits after 5 invalid attempts" {
  run expect -c '
      set timeout 30
      spawn bash -c "curl -L https://meshery.io/install | PLATFORM=bob bash -"

      expect "platform"
      send "deomon\r"
      expect "platform"
      send "youtube\r"
      expect "platform"
      send "meshery\r"
      expect "platform"
      send "\r"
      expect "platform"
      send "bob\r"

      expect eof
  '

  assert_success

  assert_output --partial "Invalid platform. Enter a valid platform to deploy Meshery [docker, kubernetes]:"
}