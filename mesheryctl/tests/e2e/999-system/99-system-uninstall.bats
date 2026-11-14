#!/usr/bin/env bats

# Runs in the 999-system suite, last of all suites (`bats *-*/*.bats` expands
# directories in name order and 99- sorts after 11-): `system uninstall` removes
# the whole installation, so every suite that needs a running Meshery has to
# complete first.
#
# Requires a live Kubernetes installation of Meshery, which the e2e CI job
# starts before the suite runs; there is nothing to uninstall on a bare
# checkout.

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    export UNINSTALL_START_MSG="Uninstalling Meshery resources"
    export UNINSTALL_SUCCESS_MSG="Meshery uninstalled"
}

@test "[cut=System][tg=System Management] given an argument when running mesheryctl system uninstall then an error message is displayed" {
    run $MESHERYCTL_BIN system uninstall invalid-arguments
    assert_failure
    assert_output --partial "this command takes no arguments"
}

@test "[cut=System][tg=System Management] given a running installation when running mesheryctl system uninstall -y then Meshery and its namespace are removed" {
    run $MESHERYCTL_BIN system uninstall -y
    assert_success
    assert_output --partial "$UNINSTALL_START_MSG"
    assert_output --partial "$UNINSTALL_SUCCESS_MSG"

    # Nothing is left running afterwards
    run $MESHERYCTL_BIN system status -y
    assert_success
    assert_output --partial "Meshery is not running"

    # ... and the namespace the command deletes is actually gone
    run kubectl get namespace meshery
    assert_failure
}
