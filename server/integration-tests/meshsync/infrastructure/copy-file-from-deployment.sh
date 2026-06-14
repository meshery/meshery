#!/usr/bin/env bash

## Expected variables to be defined
# NAMESPACE
# DEPLOYMENT
# REMOTE_PATH
# LOCAL_PATH

# Get the name of one pod from the deployment
POD=$(kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/name="$DEPLOYMENT" -o jsonpath="{.items[0].metadata.name}")

if [ -z "$POD" ]; then
  echo "No pod found for deployment $DEPLOYMENT in namespace $NAMESPACE"
  exit 1
fi

# Copy the file from the pod to local
kubectl cp "$NAMESPACE/$POD:$REMOTE_PATH" "$LOCAL_PATH"
