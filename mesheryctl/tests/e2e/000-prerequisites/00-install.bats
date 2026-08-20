#!/usr/bin/env bats

# [Test Plan] Test #4: Install mesheryctl but do not deploy Meshery

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    # Create an isolated temporary directory for test artifacts
    TEST_TMP_DIR=$(mktemp -d)
}

teardown() {
    if [ -d "$TEST_TMP_DIR" ]; then
        rm -rf "$TEST_TMP_DIR"
    fi
}

@test "[TC-0004][cut=Installer][tg=Installation] given DEPLOY_MESHERY=false when running install script then mesheryctl is installed without deploying Meshery" {
    # Execute the installer script with DEPLOY_MESHERY=false
    run bash -c "DEPLOY_MESHERY=false bash -c \"\$(curl -sL https://meshery.io/install)\""
    assert_success

    # Assert that the installer prints the expected instruction to start Meshery
    assert_output --partial 'Run "mesheryctl system start" to start Meshery.'

    # Verify that mesheryctl was installed and can execute client version
    if [ -f "$HOME/.meshery/bin/mesheryctl" ]; then
        run "$HOME/.meshery/bin/mesheryctl" version --client
        assert_success
    elif [ -f "/usr/local/bin/mesheryctl" ]; then
        run /usr/local/bin/mesheryctl version --client
        assert_success
    elif [ -n "$MESHERYCTL_BIN" ] && [ -f "$MESHERYCTL_BIN" ]; then
        run "$MESHERYCTL_BIN" version --client
        assert_success
    fi
}
