#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  # Download the live installer so this test exercises the same
  # platform-validation logic used by the production installer.

  TEST_DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" && pwd)"
  INSTALL_SCRIPT="$TEST_DIR/install.sh"
  PLATFORM_SCRIPT="$TEST_DIR/platform.sh"
  

  curl -fsSL \
    --connect-timeout 10 \
    --max-time 60 \
    https://meshery.io/install \
    -o "$INSTALL_SCRIPT"


  # Extract only the platform-selection logic from the installer.
  # Everything after COMMON FUNCTIONS is excluded to prevent the test
  # from executing the actual mesheryctl installation or Meshery startup.
  # This allows the test to focus solely on validating accepted platforms
  # and rejecting invalid platform inputs.
  sed '/^####### COMMON FUNCTIONS$/,$d' \
    "$INSTALL_SCRIPT" > "$PLATFORM_SCRIPT"

  chmod +x "$PLATFORM_SCRIPT"
  
}

teardown() {
  rm -f "$INSTALL_SCRIPT" "$PLATFORM_SCRIPT"
}

@test "Given wrong platforms provided when installing then installer exits after 5 invalid attempts" {

    run expect -c "
        set timeout 30
        spawn env PLATFORM=bob bash \"$PLATFORM_SCRIPT\"

        expect \"platform\"
        send \"deomon\r\"

        expect \"platform\"
        send \"youtube\r\"

        expect \"platform\"
        send \"meshery\r\"

        expect \"platform\"
        send \"\r\"

        expect \"platform\"
        send \"bob\r\"

        expect eof
    "

  assert_success
  assert_output --partial "Too many invalid attempts. Please try again."
}


@test "Given docker as the platform when installing then docker is accepted" {

  run env PLATFORM=docker "$PLATFORM_SCRIPT"
  assert_success

}

@test "Given kubernetes as the platform when installing then kubernetes is accepted" {

  run env PLATFORM=kubernetes "$PLATFORM_SCRIPT"
  assert_success

}