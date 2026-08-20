#!/usr/bin/env bats

# [Test Plan] Test #4: Install mesheryctl but do not deploy Meshery

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    # Create isolated temporary directory and pre-create bin dir to target isolated install
    TEST_TMP_DIR=$(mktemp -d)
    mkdir -p "$TEST_TMP_DIR/.meshery/bin"
}

teardown() {
    if [ -d "$TEST_TMP_DIR" ]; then
        rm -rf "$TEST_TMP_DIR"
    fi
}

@test "[TC-0004][cut=Installer][tg=Installation] given DEPLOY_MESHERY=false when running install script then mesheryctl is installed without deploying Meshery" {
    # Execute the installer script with DEPLOY_MESHERY=false targeting isolated HOME
    run bash -c "HOME='$TEST_TMP_DIR' DEPLOY_MESHERY=false bash -c \"\$(curl --proto '=https' --tlsv1.2 -sSfL --max-time 30 https://meshery.io/install)\""
    assert_success

    # Assert that the installer prints the expected instruction to start Meshery
    assert_output --partial 'Run "mesheryctl system start" to start Meshery.'

    # Assert that mesheryctl binary was created in the isolated install directory
    assert_file_exists "$TEST_TMP_DIR/.meshery/bin/mesheryctl"

    # Verify that the installed binary is executable and returns client version
    run "$TEST_TMP_DIR/.meshery/bin/mesheryctl" version --client
    assert_success

    # Assert that Meshery was not deployed or started (no running meshery containers)
    if command -v docker >/dev/null 2>&1; then
        run docker ps --filter "name=meshery" --format "{{.Names}}"
        assert_success
        refute_output --partial "meshery"
    fi
}
