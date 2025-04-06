#!/bin/bash

# Setup up mesheryctl
MESHERYCTL_FILE="mesheryctl"

# Default flag value
BUILD_FLAG=false

# Get parsed build flag
while getopts "b" opt; do
   case "$opt" in
      b) BUILD_FLAG=true ;;
      \?) 
      echo -e "\nUsage: $0 [-b: to rebuild binary] " 
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

   export SUPPORT_DESTDIR="../helpers/bats-support/load.bash"

	export ASSERT_DESTDIR="../helpers/bats-assert/load.bash"
   
	export DETIK_DESTDIR="../helpers/bats-detik//lib/detik.bash"
   
   echo -e "\nCreate temp directory for test data"
   TEMP_DATA_DIR=`mktemp -d`
   # Expose the temp directory to the following tests
   export TEMP_TEST_DATA_DIR=$TEMP_DATA_DIR

   echo -e "### done: Test environment setup ###\n"

   # Run tests
   # Uncomment the following line to enable junit format output
   # FORMATTER="--formatter tap"
   bats $FORMATTER *-*/*.bats

   test_result=$?

   # Run suite teardown (always)

   exit $test_result
}

main
