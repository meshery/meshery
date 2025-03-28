#!/usr/bin/env bash

_tests_helper() {
    export BATS_LIB_PATH=${BATS_LIB_PATH:-"/usr/lib"}
    bats_load_library bats-support
    bats_load_library bats-assert
    bats_load_library bats-detik/lib/detik.bash
}
