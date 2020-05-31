#!/bin/bash
set -e
set -o pipefail

TARGET_FILE="/tmp/meshery/config-minikube.yaml"

kubectl config view --minify --flatten > $TARGET_FILE 
