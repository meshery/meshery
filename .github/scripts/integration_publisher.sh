#!/bin/bash

# Setup for local testing
# INTEGRATIONS=$(< results.json)
# integration_count=$(jq ' .results[].result.rawData | length' results.json)

# Setup for repository-based run
#results=$(< results.json)
INTEGRATIONS=$(< .github/scripts/integrations.json)
integration_count=$(jq ' .results[].result.rawData | length' $INTEGRATIONS)

for (( i=0; i<integration_count; i++ ))
do 
    integration_name=$(jq " .results[].result.rawData[${i}][0]" -r $INTEGRATIONS)
    integration_data=$(jq " .results[].result.rawData[${i}][1]" -r $INTEGRATIONS)
    mkdir -p ${integration_name}/icon/color
    mkdir -p ${integration_name}/icon/white
    echo $integration_data > ${integration_name}/index.mdx   
done