#!/usr/bin/env bash

_tests_helper() {
    export BATS_LIB_PATH=${BATS_LIB_PATH:-"/usr/lib"}
    echo "Loading BATS libraries from $BATS_LIB_PATH"
    bats_load_library bats-support
    bats_load_library bats-assert

    # Required due to difference between local setup and ci
    if [ -f "$BATS_LIB_PATH/bats-detik/lib/detik.bash" ]
    then
        # For local setup
        bats_load_library bats-detik/lib/detik.bash
    else
        # for ci
        bats_load_library bats-detik/detik.bash
    fi
}
