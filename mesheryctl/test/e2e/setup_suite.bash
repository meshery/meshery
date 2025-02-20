#!/bin/bash

echo "Setting up test environment"

# Moving to mesheryctl folder
# pushd $PWD/../.. > /dev/null  # Redirect output to suppress pushd messages

# Install mesheryctl
curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -
 
# echo  "Build mesheryctl binary"
# make clean
# make > /dev/null
# Expose binary path to the following tests
export MESHERYCTL_BIN="mesheryctl"

# Go back to the tests directory
# popd > /dev/null

echo -e "\ncreate temp directory for test data"
TEMP_DATA_DIR=`mktemp -d`
# Expose the temp directory to the following tests
export TEMP_TEST_DATA_DIR=$TEMP_DATA_DIR

echo -e "Setup complete\n\n\n"

# Now the "program" executable is available for all test
