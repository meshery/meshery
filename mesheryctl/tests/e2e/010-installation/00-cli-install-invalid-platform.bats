#!/usr/bin/env bats

setup_file() {
  export INSTALL_SCRIPT="${BATS_FILE_TMPDIR}/install.sh"
  curl -fsSL --connect-timeout 10 --max-time 60 https://meshery.io/install -o "$INSTALL_SCRIPT"

  # Extract ONLY the getopts block into a variable for focused testing
  export GETOPTS_BLOCK
  GETOPTS_BLOCK=$(sed -n '/while getopts/,/done/p' "$INSTALL_SCRIPT")
}

teardown_file() {
  rm -f "$INSTALL_SCRIPT"
}

@test "Supported platform (e.g. 'docker') is defined in getopts cases" {
  # Check if 'docker)' case exists specifically inside the getopts block
  run grep -E '^\s*docker\)' <<< "$GETOPTS_BLOCK"

  [ "$status" -eq 0 ]
  [[ "$output" =~ "docker" ]]
}

@test "Supported platform (e.g. 'kubernetes') is defined in getopts cases" {
  # Check if 'docker)' case exists specifically inside the getopts block
  run grep -E '^\s*kubernetes\)' <<< "$GETOPTS_BLOCK"

  [ "$status" -eq 0 ]
  [[ "$output" =~ "kubernetes" ]]
}

@test "Unsupported platform (e.g. 'bob') is NOT defined in getopts cases" {
  # Check that 'bob)' case does NOT exist inside the getopts block
  run grep -E '^\s*bob\)' <<< "$GETOPTS_BLOCK"

  [ "$status" -ne 0 ]
}
@test "Unsupported platform (e.g. 'else') is NOT defined in getopts cases" {
  # Check that 'bob)' case does NOT exist inside the getopts block
  run grep -E '^\s*else\)' <<< "$GETOPTS_BLOCK"

  [ "$status" -ne 0 ]
}