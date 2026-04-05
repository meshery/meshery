#!/usr/bin/env bash

create_meshery_config_folder() {
    echo "start: Create meshery config folder"
    if [ ! -d "$HOME/.meshery" ]
    then
        mkdir "$HOME/.meshery"
    fi
    echo "done: Create meshery config folder"
}

# Generate auth file to communicate with meshery server
create_auth_file() {
    echo "start: authentication configuration" 
    echo '{ "meshery-provider": "Layer5", "token": null }' | jq -c '.token = "'$MESHERY_PROVIDER_TOKEN'"' > "${HOME}/.meshery/auth.json"
    echo "done: authentication configuration"
}

main() {
    echo -e "### start: Test environment setup ###\n"
    
    export MESHERYCTL_BIN="../../mesheryctl"
    export MESHERY_CONFIG_FILE_PATH="${HOME}/.meshery/config.yaml"
	export MESHERY_AUTH_FILE="${HOME}/.meshery/auth.json"
    export E2E_HELPERS_PATH="${E2E_HELPERS_PATH}"
    export BATS_LIB_PATH="${BATS_LIB_PATH}"

    echo -e "\nCreate temp directory for test data"
    TEMP_DATA_DIR=`mktemp -d`
    # Expose the temp directory to the following tests
    export TEMP_DATA_DIR=$TEMP_DATA_DIR

    create_meshery_config_folder
    create_auth_file 

    echo -e "### done: Test environment setup ###\n"
}


main
