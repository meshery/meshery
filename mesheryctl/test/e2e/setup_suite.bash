#!/bin/bash

echo "Setting up test environment"


# Expose binary path to the following tests
export MESHERYCTL_BIN="mesheryctl"

echo -e "\ncreate temp directory for test data"
TEMP_DATA_DIR=`mktemp -d`
# Expose the temp directory to the following tests
export TEMP_TEST_DATA_DIR=$TEMP_DATA_DIR

echo -e "Setup complete\n\n\n"

# Now the "program" executable is available for all test
