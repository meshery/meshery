#!/bin/bash

echo "Setting up test environment\n"

echo "Install mesheryctl"
curl -L https://meshery.io/install -s | DEPLOY_MESHERY=false bash -

# Generate auth file to comunicate with meshery server
echo "{\"meshery-provider\": \"Meshery\", \"token\": \"${MESHERY_PROVIDER_TOKEN}\"}" > "$HOME/.meshery/auth.json"

# Expose binary path to the following tests
export MESHERYCTL_BIN="mesheryctl"

echo -e "\nCreate temp directory for test data"
TEMP_DATA_DIR=`mktemp -d`
# Expose the temp directory to the following tests
export TEMP_TEST_DATA_DIR=$TEMP_DATA_DIR


echo -e "Setup complete\n\n\n"

# Now the "program" executable is available for all test
