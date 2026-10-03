#!/usr/bin/env bats
#
# E2E tests for the Meshery installer script.

load "clean_helper.sh"

setup_file() {
  remove_all "setup_file"
}

setup() {
  export BATS_LIB_PATH=${BATS_LIB_PATH:-"/usr/lib"}
  bats_load_library bats-support
  bats_load_library bats-assert
  bats_load_library bats-detik/detik.bash

  DETIK_CLIENT_NAMESPACE="meshery"
  DETIK_CLIENT_NAME="kubectl"

}

teardown() {
  remove_all "teardown"
}

# Helper: run an `expect` script with a hard wall-clock timeout so a hang
# in the installer cannot stall the whole suite.
run_expect() {
  local script="$1"
  local wall_timeout="${2:-300}"
  run timeout "$wall_timeout" expect -c "$script"
}


@test "Given wrong platforms provided when installing then failed after 5 invalid attempts" {
  run_expect '
    set timeout 150
    log_user 1

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=bob bash -}

    expect {
      -re {(?i)select.*platform|platform.*:} { send "deomon\r" }
      timeout { exit 2 }
    }

    expect {
      -re {(?i)select.*platform|platform.*:} { send "youtube\r" }
      timeout { exit 2 }
    }

    expect {
      -re {(?i)select.*platform|platform.*:} { send "meshery\r" }
      timeout { exit 2 }
    }

    expect {
      -re {(?i)select.*platform|platform.*:} { send "\r" }
      timeout { exit 2 }
    }

    expect {
      -re {(?i)select.*platform|platform.*:} { send "bob\r" }
      timeout { exit 2 }
    }

    expect eof
    
    # Capture the exit status of the spawned process
    lassign [wait] pid spawnid os_error_flag exit_code
    exit $exit_code
  ' 180

  assert_failure
  assert_output --partial "Too many invalid attempts. Please try again."
  # The script should reject each bad value before ultimately bailing out,
  # not just print the final message with a lucky earlier match.
  refute_output --partial "Starting Meshery..."
}


@test "Given docker as the platform when installing then succeed" {
  run_expect '
    set timeout 300

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=docker bash -}

    expect eof
  ' 320

  assert_success
  assert_output --partial "Starting Meshery..."
  assert_output --partial "mesheryctl installed."
  assert_output --partial "Meshery endpoint is accessible."

  # Functional check, not just log-message matching: confirm mesheryctl was
  # actually placed on PATH and that Meshery believes itself to be running.
  run bash -lc "command -v mesheryctl"
  assert_success

  run bash -lc "mesheryctl system status"
  assert_success
 
}


@test "Given kubernetes as the platform when installing then succeed" {
  run_expect '
    set timeout 300

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=kubernetes bash -}

    expect eof
  ' 320

  assert_success
  assert_output --partial "Meshery deployed on Kubernetes."

  # Functional check: verify the Meshery pods actually exist and reach
  # Ready state, rather than trusting only the printed success message.

  run kubectl wait --for=condition=Ready pod --all -n meshery --timeout=120s
  assert_success
  
  run verify "there are more than 0 pod named '^meshery-[a-z0-9]+-[a-z0-9]+$'"
	assert_success  

  run verify "'status' is 'running' for pods named '^meshery-[a-z0-9]+-[a-z0-9]+$'"
  assert_success
}


@test "Given DEPLOY_MESHERY=false when installing then mesheryctl binary is present" {
  run_expect '
    set timeout 300

    spawn bash -c {curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -}

    expect eof
  ' 320

  assert_success
  assert_output --partial "mesheryctl installed."
  assert_output --partial "Run \"mesheryctl system start\" to start Meshery."

  # DEPLOY_MESHERY=false should mean nothing was actually deployed — verify
  # that no meshery namespace/resources were created as a side effect.
  run verify "there is 1 namespace named '^meshery$'"
  assert_failure

  # Assert that no Meshery Docker container exists after installation.
  run bash -lc 'test -z "$(docker ps -aq --filter name=meshery)"'
  assert_success
}

