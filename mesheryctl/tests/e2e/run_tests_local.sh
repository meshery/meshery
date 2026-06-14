#!/bin/bash

# Setup up mesheryctl
MESHERYCTL_FILE="mesheryctl"

# Default build flag value
BUILD_FLAG=false

# Default bats file path value
BATS_FILE=""

# Get parsed build flag
while getopts "f:b" opt; do
   case "$opt" in
      f) FILE="$OPTARG" ;;
      b) BUILD_FLAG=true ;;
      \?) 
      echo -e "\nUsage: $0 [-f filename: to specify single bat file to run] [-b: to rebuild binary] " 
      exit 1 ;;
   esac
done

check_bin_file() {
   if [ -f "../../$MESHERYCTL_FILE" ]; then
      echo "$MESHERYCTL_FILE Binary file found..."
      if $BUILD_FLAG; then
         echo -e "\nBuild flag parsed, rebuilding $MESHERYCTL_FILE binary..."
         build_bin
      fi
   else
      echo -e "$MESHERYCTL_FILE Binary file not found.. Attempting to build binary...\n"
      build_bin
   fi
}

build_bin() {
   cd ../../

   if ! make; then
      echo "X Build failed. Check for errors or dependencies."
      exit 1
   fi

   if [ -f "$MESHERYCTL_FILE" ]; then
      echo "√ $MESHERYCTL_FILE Build successful..."
      cd tests/e2e
   else
      echo "X Build failed. Binary file not found."
      return 1
   fi
}

main() {
   echo -e "### start: Test environment setup ###\n"

   check_bin_file

   export MESHERYCTL_BIN="../../$MESHERYCTL_FILE"
   export MESHERY_CONFIG_FILE_PATH="$HOME/.meshery/config.yaml"
   export MESHERY_AUTH_FILE="$HOME/.meshery/auth.json"
   export E2E_HELPERS_PATH="$(pwd)/helpers"
   export BATS_LIB_PATH="$(pwd)/helpers/bats-libs"
   export BATS_TEST_DIRNAME=""
 
   echo -e "\nCreate temp directory for test data"
   TEMP_DATA_DIR=`mktemp -d`
   # Expose the temp directory to the following tests
   export TEMP_DATA_DIR=$TEMP_DATA_DIR

   echo -e "### done: Test environment setup ###\n"

   # Run tests
   # Uncomment the following line to enable junit format output
   # FORMATTER="--formatter tap"

   if [[ -z "$FILE" ]]; then
      echo "Running all E2E tests."
      bats $FORMATTER *-*/*.bats
   else
      if [[ "$FILE" == *.bats && -f "$FILE" ]]; then
         echo "Running E2E for specified file: $FILE"
         bats $FORMATTER $FILE
      else
         echo "X Invalid file format or file not found"
         exit 1
      fi
   fi

   test_result=$?

   # Run suite teardown (always)

   exit $test_result
}

main
