#!/usr/bin/env bash
MESHERYCTL_PORT_FORWARDING=9081

install_mesheryctl() {
    echo "start: Install mesheryctl"
    curl -L https://meshery.io/install -s | PLATFORM=$1 bash -
    echo "done: Install mesheryctl"
}

create_meshery_config_folder() {
    echo "start: Create meshery config folder"
    if [ ! -d "$HOME/.meshery" ]
    then
        mkdir "$HOME/.meshery"
    fi
    echo "done: Create meshery config folder"
}

# Generate auth file to communicate with meshery server
configure_auth_file() {
    echo "start: authentication configuration" 
    
    echo "done: authentication configuration"
}

port_forwarding() {
    echo "start: Port forwarding"
    context="$(yq '.current-context' "${HOME}/.meshery/config.yaml")"
    MESHERYCTL_PORT_FORWARDING=$(yq ".contexts.\"${context}\".endpoint" ${HOME}/.meshery/config.yaml | cut -d':' -f3)

    nohup kubectl -n meshery port-forward svc/meshery ${MESHERYCTL_PORT_FORWARDING}:$(kubectl -n meshery get svc/meshery -o jsonpath='{.spec.ports[0].port}') &
    echo "done: Port forwarding"
}

main() {
    echo -e "### start: Test environment setup ###\n"

    # install_mesheryctl "$MESHERY_PLATFORM"
    # create_meshery_config_folder
    # configure_auth_file 
    port_forwarding
    # config_mesheryctl_port_forwarding_endpoint
    
    export MESHERYCTL_BIN="../../mesheryctl"
    export MESHERY_CONFIG_FILE_PATH="${HOME}/.meshery/config.yaml"
	export MESHERY_AUTH_FILE="${HOME}/.meshery/auth.json"
    export E2E_HELPERS_PATH="${E2E_HELPERS_PATH}"
    export BATS_LIB_PATH="${BATS_LIB_PATH}"


    echo -e "\nCreate temp directory for test data"
    TEMP_DATA_DIR=`mktemp -d`
    # Expose the temp directory to the following tests
    export TEMP_DATA_DIR=$TEMP_DATA_DIR


    echo -e "### done: Test environment setup ###\n"
}


main
