#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

setup_file() {
  export INSTALL_SCRIPT="${BATS_FILE_TMPDIR}/install.sh"
  curl -fsSL --connect-timeout 10 --max-time 60 https://meshery.io/install -o "$INSTALL_SCRIPT"

  export ALLOWED_CASES
  export ENV_ALLOWED_CASES

  # 1. Extracts platform cases from the getopts flag parsing block (-p)
  ALLOWED_CASES=$(sed -n '/case $OPTARG in/,/esac/p' "$INSTALL_SCRIPT" | grep -E '^\s+[a-zA-Z0-9_-]+\)' | tr -d ' )' | grep -v '^\*$' | sort -u)

  # 2. Extracts platform cases from the environment variable check block ($PLATFORM)
  ENV_ALLOWED_CASES=$(awk '/case "\$PLATFORM" in/,/esac/' "$INSTALL_SCRIPT" | grep -E '\)' | tr -d ' );' | tr '|' '\n' | grep -v '^\*' | grep -v '^$' | sort -u)
}

teardown_file() {
  rm -f "$INSTALL_SCRIPT"
}


@test "[TC-0100][tg=installation]verify getopts flag -p contains only expected supported platforms (docker or kubernetes only)" {
  run echo "$ALLOWED_CASES"

  assert_success
  assert_line --index 0 "docker"
  assert_line --index 1 "kubernetes"
  assert_equal "${#lines[@]}" 2
}

@test "[TC-0101][tg=installation]verify PLATFORM environment variable block contains only expected supported platforms (docker or kubernetes only)" {
  run echo "$ENV_ALLOWED_CASES"

  assert_success
  assert_line --index 0 "docker"
  assert_line --index 1 "kubernetes"
  assert_equal "${#lines[@]}" 2
}