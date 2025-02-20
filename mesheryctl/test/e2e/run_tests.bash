#!/bin/bash

# Run suite setup
source ./setup_suite.bash

# echo "DEBUG: MESHERYCTL_BIN=$MESHERYCTL_BIN"
# echo "DEBUG: TEMP_TEST_DATA_DIR=$TEMP_TEST_DATA_DIR"

# Run the tests

# Uncomment the following line to enable junit format output
FORMATTER="--formatter junit"

# Run the prerequisites first
## cli is available
## server is started
#bats $FORMATTER prerequisites/*.bats

# make server available for the following tests
#kubectl port-forward -n meshery service/meshery 9081:9081 &

# Run the tests
bats $FORMATTER *-*/*.bats

test_result=$?  # Capture the exit code of bats

# Run suite teardown (always)
./teardown_suite.bash

exit $test_result