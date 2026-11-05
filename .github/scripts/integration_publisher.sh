#!/bin/bash

# Setup for local testing
# INTEGRATIONS=$(< results.json)
# integration_count=$(jq ' .results[].result.rawData | length' results.json)

# Setup for repository-based run
#results=$(< results.json)
curl -o integrations.csv https://docs.google.com/spreadsheets/d/e/2PACX-1vSgOXuiqbhUgtC9oNbJlz9PYpOEaFVoGNUFMIk4NZciFfQv1ewZg8ahdrWHKI79GkKK9TbmnZx8CqIe/pub?gid=0&single=true&output=csv
brew install xsv
xsv slice -i 7 integrations.csv | xsv select 28 | xsv flatten > a.txt
echo "$PATH"
echo "{path}" 
echo $GITHUB_PATH
echo "{path}" >> $GITHUB_PATH
pwd
ls
echo $GITHUB_WORKSPACE

#INTEGRATIONS=$(< $GITHUB_WORKSPACE/integrations.json)
integration_count=$(jq ' .results[].result.rawData | length' $GITHUB_WORKSPACE/integrations.json)

for (( i=0; i<integration_count; i++ ))
do 
    integration_name=$(jq " .results[].result.rawData[${i}][0]" -r $GITHUB_WORKSPACE/integrations.json)
    integration_data=$(jq " .results[].result.rawData[${i}][1]" -r $GITHUB_WORKSPACE/integrations.json)
    mkdir -p ${integration_name}/icon/color
    mkdir -p ${integration_name}/icon/white
    echo $integration_data > ${integration_name}/index.mdx   
done