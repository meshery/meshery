#!/usr/bin/env bats

# Create a temporary nginx YAML file for testing with proper .yaml extension
setup() {
  # Create a temp directory
  TEMP_DIR=$(mktemp -d)
  DEMO_YAML="$TEMP_DIR/nginx-k8s.yaml"
  
  cat > "$DEMO_YAML" << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
EOF

  # Create an invalid YAML file for negative testing
  INVALID_YAML="$TEMP_DIR/invalid.yaml"
  cat > "$INVALID_YAML" << EOF
This is not valid YAML
  - item1
    item2:
  random text
EOF
}

# Clean up the temporary file and directory after tests
teardown() {
  if [ -d "$TEMP_DIR" ]; then
    rm -rf "$TEMP_DIR"
  fi
  
  # Keep design ID file for delete tests
}

@test "mesheryctl design import with nginx YAML is succeeded" {
  run $MESHERYCTL_BIN design import -f "$DEMO_YAML" -s "Kubernetes Manifest"
  echo "Command output: $output"
  [ "$status" -eq 0 ]
  
  echo "$output" | grep -q "has been imported"
  [ "$?" -eq 0 ]
  
  echo "$output" | grep -q "Design ID:"
  [ "$?" -eq 0 ]
  
  # Extract design ID - ID appears after "Design ID:" in the output
  DESIGN_ID=$(echo "$output" | grep -o "Design ID: [a-f0-9\-]*" | cut -d' ' -f3 || echo "")
  
  # If we found an ID, save it
  if [ -n "$DESIGN_ID" ]; then
    echo "$DESIGN_ID" > /tmp/meshery_test_design_id
    echo "Extracted design ID: $DESIGN_ID"
  else
    echo "Warning: Could not extract design ID from output"
    # Additional debug output to understand if we're getting a different format
    echo "Full output:"
    echo "$output"
  fi
}

@test "mesheryctl design import with custom name is succeeded" {
  CUSTOM_NAME="custom-nginx-design"
  run $MESHERYCTL_BIN design import -f "$DEMO_YAML" -s "Kubernetes Manifest" -n "$CUSTOM_NAME"
  echo "Command output: $output"
  [ "$status" -eq 0 ]
  
  # Verify the design was named correctly
  echo "$output" | grep -q "The design file '$CUSTOM_NAME' has been imported"
  [ "$?" -eq 0 ]
  
  # Extract the ID of this design for later verification
  DESIGN_ID2=$(echo "$output" | grep -o "Design ID: [a-f0-9\-]*" | cut -d' ' -f3 || echo "")
  if [ -n "$DESIGN_ID2" ]; then
    echo "$DESIGN_ID2" > /tmp/meshery_test_design_id2
  fi
}

@test "mesheryctl design import with invalid YAML should fail" {
  run $MESHERYCTL_BIN design import -f "$INVALID_YAML" -s "Kubernetes Manifest"
  echo "Command output: $output"
  
  # Based on the actual output, look for specific error messages
  echo "$output" | grep -q "Response Status Code 400"
  [ "$?" -eq 0 ]
  
  # Check for YAML parsing error message
  echo "$output" | grep -q "could not be parsed due to invalid YAML syntax"
  [ "$?" -eq 0 ]
  
  if echo "$output" | grep -q "has been imported"; then
    # If "has been imported" is found, fail the test
    [ 1 -eq 0 ]  # This will always fail
  else
    # If "has been imported" is NOT found, pass this check
    [ 0 -eq 0 ]  # This will always pass
  fi
}
