#!/bin/bash
set -e
set -o pipefail

TARGET_FOLDER="/tmp/meshery"
TARGET_FILE="$TARGET_FOLDER/kubeconfig.yaml"

mkdir -p $TARGET_FOLDER
kubectl config view --minify --flatten > $TARGET_FILE 
