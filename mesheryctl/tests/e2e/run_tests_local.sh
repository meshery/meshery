#!/bin/bash

# Setup up mesheryctl
MESHERYCTL_FILE="mesheryctl"

check_bin_file() {
   if [ -f "../../$MESHERYCTL_FILE" ]; then
      echo "$MESHERYCTL_FILE Binary file found..."
   else
      echo -e "$MESHERYCTL_FILE Binary file not found.. Attempting to build binary...\n"
      build_bin
   fi
}

build_bin() {
   cd ../../
   pwd; ls

   if ! make; then
      echo "X Build failed. Check for errors or dependencies."
      exit 1
   fi

   if [ -f "$MESHERYCTL_FILE" ]; then
      echo "√ $MESHERYCTL_FILE Build successful..."
   else
      echo "X Build failed. Binary file not found."
      return 1
   fi
}

create_auth_file() {
   MESHERY_PROVIDER_TOKEN=$1
   echo "start: authentication configuration"
   if [ ! -d "$HOME/.meshery" ]
   then
      mkdir "$HOME/.meshery"
   fi

   # Generate auth file to comunicate with meshery server
   # if [ ! -f "$HOME/.meshery/auth.json" ]
   # then 
   # fi
   echo "{\"meshery-provider\": \"Meshery\", \"token\": \"${MESHERY_PROVIDER_TOKEN}\"}" > "$HOME/.meshery/auth.json"
   echo "done: authentication configuration"
}

main() {
   echo -e "### start: Test environment setup ###\n"

   check_bin_file
   create_auth_file $1

   export MESHERYCTL_BIN="../../mesheryctl"
   
   echo -e "\nCreate temp directory for test data"
   TEMP_DATA_DIR=`mktemp -d`
   # Expose the temp directory to the following tests
   export TEMP_TEST_DATA_DIR=$TEMP_DATA_DIR

   echo -e "### done: Test environment setup ###\n"
}

main $1

# Run tests
# Uncomment the following line to enable junit format output
# FORMATTER="--formatter tap"
bats $FORMATTER *-*/*.bats

test_result=$?

exit $test_result