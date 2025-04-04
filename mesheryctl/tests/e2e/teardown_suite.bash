#!/usr/bin/env bash

echo -e "\n\n\nCleaning up test data directory: $TEMP_TEST_DATA_DIR"
rm -rf $TEMP_TEST_DATA_DIR

echo "Cleaning up mesheryctl binaries"
pushd $PWD/../.. > /dev/null
make clean

popd > /dev/null