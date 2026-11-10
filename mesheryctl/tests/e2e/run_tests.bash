#!/usr/bin/env bash

# Run suite setup
source ./setup_suite.bash

# echo "DEBUG: MESHERYCTL_BIN=$MESHERYCTL_BIN"
# echo "DEBUG: TEMP_DATA_DIR=$TEMP_DATA_DIR"

# Run the tests
# Uncomment the following line to enable junit format output
FORMATTER="--formatter tap"

# --print-output-on-failure makes BATS emit each failing test's captured $output
# / $stderr as TAP `#` diagnostics. bats-to-allure.js turns those into the
# Allure failure trace + a text attachment, so a failed test is debuggable in
# the Connection Lifecycle report instead of showing only the assertion line.
bats $FORMATTER --print-output-on-failure *-*/*.bats

test_result=$?  # Capture the exit code of bats

# Run suite teardown (always)
./teardown_suite.bash

exit $test_result