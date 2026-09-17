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
    # Build the token in with jq --arg so the value is passed as data, never
    # spliced into the jq program. The previous form,
    # jq '.token = "'$MESHERY_PROVIDER_TOKEN'"', interpolated the token into the
    # filter text: any '"' or '\' in the token would break the JSON, and an
    # unset var would silently yield an empty token. --arg is injection-safe for
    # any token shape.
    jq -cn --arg token "$MESHERY_PROVIDER_TOKEN" \
        '{ "meshery-provider": "Meshery", token: $token }' > "$MESHERY_AUTH_FILE"
    echo "done: authentication configuration"
}

set_context_to_meshery() {
    echo "start: set context to Meshery"
    yq -i '.contexts.local.provider = "Meshery"' "$MESHERY_CONFIG_FILE_PATH"
    echo "done: set context to Meshery"
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
    set_context_to_meshery

    echo -e "### done: Test environment setup ###\n"
}


main
