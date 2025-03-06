#!/bin/bash


install_mesheryctl() {
    echo "start: Install mesheryctl"
    curl -L https://meshery.io/install -s | DEPLOY_MESHERY=false bash -
    echo "done: Install mesheryctl"
}



start_meshery() {
    platform="$1"
    echo "start: meshery server start for platform $platform"
    mesheryctl system start -p "$platform"
    echo "in progress: meshery server start for platform $platform"
    sleep 40
    echo "done: meshery server start for platform $platform"
}


create_auth_file() {
    echo "start: authentication configuration"
    if [ ! -d "$HOME/.meshery" ]
    then
        mkdir "$HOME/.meshery"
    fi

    # Generate auth file to comunicate with meshery server
    echo "{\"meshery-provider\": \"Meshery\", \"token\": \"${MESHERY_PROVIDER_TOKEN}\"}" > "$HOME/.meshery/auth.json"
    echo "done: authentication configuration"
}

main() {
    echo -e "### start: Test environment setup ###\n"

    install_mesheryctl
    start_meshery "$MESHERY_PLATFORM"
    create_auth_file 
    
    export MESHERYCTL_BIN="mesheryctl"

    echo -e "\nCreate temp directory for test data"
    TEMP_DATA_DIR=`mktemp -d`
    # Expose the temp directory to the following tests
    export TEMP_TEST_DATA_DIR=$TEMP_DATA_DIR


    echo -e "### done: Test environment setup ###\n"
}


main
