#!/bin/bash

# helper function to get the latest version from a list of versions
# handles semantic versioning: v1.0.0, release-5.0.19, 5.0.19-rc.23, etc
get_latest_version() {
    local versions="$1"
    echo "$versions" | sed 's/$/|/' | sort -V | sed 's/|$//' | tail -n 1
}


main() {
    local base_dir="server/meshmodel" # source for model definitions: server/meshmodel
    local output_dir="${1:-models-oci-images}" #default: models-oci-images
    
    # init log files (uploaded as artifacts)
    > build.log
    > error.log
    
    echo "Starting model image build process at $(date)" >> build.log
    

    # current working directory is mesheryctl, go up one level to find server/meshmodel 
    if [[ ! -d "../$base_dir" ]]; then
        echo "Error: Directory $base_dir does not exist" | tee -a error.log
        exit 1
    fi
    
    # Iterate through each model directory in server/meshmodel/
    # model directory: <model-name>/<model-version>/v1.0.0 
    # ref: https://github.com/meshery/meshkit/blob/master/registry/model.go#L674

    for model_dir in "../$base_dir"/*; do
        if [[ -d "$model_dir" ]]; then
            model_name=$(basename "$model_dir") #<model-name>
            echo "Processing model: $model_name" | tee -a build.log

            # Collect version directories that contain v1.0.0/ subdirectory
            versions=""
            for version_dir in "$model_dir"/*; do
                if [[ -d "$version_dir" ]]; then
                    version_name=$(basename "$version_dir") #<model-version>
                    #  Only include versions with v1.0.0/ subdirectory
                    if [[ -d "$version_dir/v1.0.0" ]]; then
                        if [[ -z "$versions" ]]; then
                            versions="$version_name"
                        else
                            versions="$versions"$'\n'"$version_name"
                        fi
                    fi
                fi
            done
            
            # get latest from collected versions
            # If no valid versions found, skip to next model
            if [[ -n "$versions" ]]; then
                latest_version=$(get_latest_version "$versions")
                model_path="$model_name/$latest_version" # to be passed to mesheryctl model build
                file_path="../$base_dir/$model_name/$latest_version"
            else
                echo "Warning: No valid versions found for $model_name (no v1.0.0 subdirectories)" | tee -a error.log
                continue
            fi

            
            echo "Latest version for $model_name: $model_path" | tee -a build.log
            
            # check if the target directory exists before running mesheryctl
            if [[ -d "$file_path" ]]; then
                echo "Running: ./mesheryctl exp model build $model_path --path ../$base_dir -o $output_dir --no-version" | tee -a build.log
                ./mesheryctl exp model build "$model_path" --path "../$base_dir" -o "$output_dir" --no-version 2>> error.log

                
                # If command was successful
                if [[ $? -eq 0 ]]; then
                    echo "✓ Successfully built model for $model_path" | tee -a build.log
                else
                    echo "✗ Failed to build model for $model_path" | tee -a error.log
                fi
            else
                echo "Warning: Path $file_path does not exist" | tee -a error.log
            fi
            
            echo "---"
        fi
    done
    
    echo "Build process completed at $(date)" >> build.log
}

# run main func
main "$@"
