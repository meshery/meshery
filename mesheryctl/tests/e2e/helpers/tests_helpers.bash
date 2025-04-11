#!/usr/bin/env bash

_tests_helper() {
    export BATS_LIB_PATH=${BATS_LIB_PATH:-"/usr/lib"}
    
    # Loading BATS libraries
    bats_load_library bats-support
    bats_load_library bats-assert
    bats_load_library bats-file

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
