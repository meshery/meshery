#!/usr/bin/env bash

# Run suite setup
source ./setup_suite.bash

# echo "DEBUG: MESHERYCTL_BIN=$MESHERYCTL_BIN"
# echo "DEBUG: TEMP_TEST_DATA_DIR=$TEMP_TEST_DATA_DIR"

./setup_suite.bash

# Run the tests
# Uncomment the following line to enable junit format output
FORMATTER="--formatter tap"

bats $FORMATTER *-*/*.bats

test_result=$?  # Capture the exit code of bats

# Run suite teardown (always)
./teardown_suite.bash

exit $test_result