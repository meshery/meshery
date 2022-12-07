#!/bin/bash

results=$(< results.json)

for i in $results
  do
    echo $i | jq ".results[].result.rawData[$i][0]"
  done