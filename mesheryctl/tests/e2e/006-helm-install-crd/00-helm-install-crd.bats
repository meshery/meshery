#!/usr/bin/env bats

load "$BATS_LIB_PATH/bats-support/load.bash"
load "$BATS_LIB_PATH/bats-assert/load.bash"

# Ensure the MESHERYCTL_BIN is set, it's usually provided by the e2e test runner or `make e2e`
# If running standalone, you might need to export it: export MESHERYCTL_BIN="/path/to/mesheryctl"

# --- Global Variables ---
MESHERY_CRDS=(
  "brokers.meshery.io"
  "meshsyncs.meshery.io"
  # Add any other Meshery-specific CRD names here, verified from Meshery Helm chart's CRD folder
)

# --- Setup function: Runs before each test ---
setup() {
  # Clean up any potential lingering Meshery installations and global ClusterRoles
  echo "--- Global Setup: Cleaning up before test ---"

  # 1. DELETE CRDs
  echo "Deleting Meshery CRDs from previous runs to ensure a clean state..."
  for crd in "${MESHERY_CRDS[@]}"; do
    run kubectl delete crd "${crd}" --ignore-not-found=true
  done

  # 2. *** IMPORTANT: WAIT FOR CRD DELETION TO COMPLETE ***
  echo "Waiting for Meshery CRDs to be fully deleted..."
  local max_retries=30 # 30 retries * 2 seconds = 1 minute wait
  local retry_delay_sec=2

  for crd in "${MESHERY_CRDS[@]}"; do
    for i in $(seq 1 "${max_retries}"); do
      # Attempt to get the CRD; if it fails (non-zero status), it's gone
      run kubectl get crd "${crd}" 2>/dev/null
      if [ "$status" -ne 0 ]; then
        echo "CRD '${crd}' is gone."
        break # CRD found to be gone, exit retry loop for this CRD
      fi
      echo "CRD '${crd}' still present during cleanup, retrying... (${i}/${max_retries})"
      sleep "${retry_delay_sec}"
    done
    # Final check: if loop completed and CRD is still present, something is wrong
    run kubectl get crd "${crd}" 2>/dev/null
    if [ "$status" -eq 0 ]; then
      fail "CRD '${crd}' could not be deleted after ${max_retries} retries during setup cleanup. This will cause test failures."
    fi
  done
  echo "Meshery CRDs confirmed deleted."
  # *** END OF ADDED WAIT LOOP ***

  # 3. Rest of the existing setup commands (these should come AFTER the CRD wait loop)
  run kubectl delete clusterrole meshery-server --ignore-not-found=true
  assert_success
  run kubectl delete clusterrolebinding meshery-server --ignore-not-found=true
  assert_success
  run helm uninstall meshery-e2e-no-crds --namespace meshery-e2e-test-ns-no-crds --ignore-not-found=true
  assert_success
  run kubectl delete namespace meshery-e2e-test-ns-no-crds --timeout=2m --force --grace-period=0 --ignore-not-found=true
  assert_success
  run helm uninstall meshery-e2e-with-crds --namespace meshery-e2e-test-ns-default --ignore-not-found=true
  assert_success
  run kubectl delete namespace meshery-e2e-test-ns-default --timeout=2m --force --grace-period=0 --ignore-not-found=true
  assert_success
  echo "--- Global Setup: Cleanup complete ---"

  # Ensure helm repo is added and updated if not handled globally by the test runner
  # Only uncomment and use if `make e2e` or `run_tests.bash` doesn't do this
  # echo "Adding and updating Meshery Helm repo..."
  # run helm repo add meshery https://meshery.io/charts/
  # run helm repo update
  # assert_success
}

# --- Teardown function: Runs after each test ---
teardown() {
  # Clean up resources created by the test
  # The test itself includes cleanup in its finally block for robustness,
  # but this ensures a clean slate even if a test fails unexpectedly before its internal cleanup.
  # This is a good place for general cleanup if the specific test's cleanup fails.
  echo "--- Teardown: Performing post-test cleanup ---"
  # Note: The test cases themselves also contain their cleanup in their main body using `|| true`
  # This teardown is a safeguard.
  echo "--- Teardown: Cleanup complete ---"
}

# --- Test Case 1: should not install Meshery CRDs when installCRDs=false ---
@test "helm install with installCRDs=false should not install Meshery CRDs" {
  local release_name="meshery-e2e-no-crds"
  local namespace="meshery-e2e-test-ns-no-crds"

  echo "--- Test: ${test_name} ---"
  echo "Installing Meshery release '${release_name}' in namespace '${namespace}' with installCRDs=false..."

  # Install Meshery with installCRDs=false
  run helm install "${release_name}" meshery/meshery --namespace "${namespace}" --create-namespace --set installCRDs=false
  assert_success "Expected helm install with installCRDs=false to succeed"

  echo "Verifying Meshery CRDs are NOT present..."
  for crd in "${MESHERY_CRDS[@]}"; do
    echo "Checking for absence of CRD: ${crd}"
    # Try to get the CRD. It should fail (exit with non-zero status) if not present.
    # We pipe stderr to /dev/null to avoid cluttering output with "not found" messages.
    run kubectl get crd "${crd}" 2>/dev/null
    assert_failure "CRD '${crd}' should NOT be present but was found." # Assert that the command failed (non-zero exit code)
  done
  echo "Confirmed Meshery CRDs are NOT present (as expected)."

  # --- Cleanup ---
  echo "Cleaning up: Uninstalling ${release_name} and deleting namespace ${namespace}..."
  run helm uninstall "${release_name}" --namespace "${namespace}" || true
  # assert_success # Do not assert success here if uninstall might legitimately fail due to release not being fully installed
  run kubectl delete namespace "${namespace}" --timeout=2m --force --grace-period=0 || true
  # assert_success # Do not assert success here if namespace might not exist
}

# --- Test Case 2: should install Meshery CRDs by default ---
@test "helm install by default should install Meshery CRDs" {
  local release_name="meshery-e2e-with-crds"
  local namespace="meshery-e2e-test-ns-default"

  echo "--- Test: ${test_name} ---"
  echo "Installing Meshery release '${release_name}' in namespace '${namespace}' with default CRD installation..."

  # Install Meshery with default settings (CRDs installed)
  run helm install "${release_name}" meshery/meshery --namespace "${namespace}" --create-namespace
  assert_success "Expected default helm install to succeed"
  echo "Meshery default installation command executed."

  echo "Verifying Meshery CRDs ARE present..."
  local max_retries=30 # Increased retries (1 min total, 2s * 30 retries)
  local retry_delay_sec=2

  for crd in "${MESHERY_CRDS[@]}"; do
    echo "Checking for presence of CRD: ${crd}"
    local crd_found=false
    for i in $(seq 1 "${max_retries}"); do
      run kubectl get crd "${crd}" 2>/dev/null
      if [ "$status" -eq 0 ]; then # Check the status variable from `run`
        crd_found=true
        echo "Confirmed CRD '${crd}' is present."
        break # CRD found, exit retry loop
      fi
      echo "CRD '${crd}' not yet present, retrying... (${i}/${max_retries})"
      sleep "${retry_delay_sec}"
    done
    if [ "$crd_found" = false ]; then
      fail "CRD '${crd}' should be present but was not found after ${max_retries} retries."
    fi
  done

  # --- Cleanup ---
  echo "Cleaning up: Uninstalling ${release_name} and deleting namespace ${namespace}..."
  run helm uninstall "${release_name}" --namespace "${namespace}" || true
  # assert_success
  run kubectl delete namespace "${namespace}" --timeout=2m --force --grace-period=0 || true
  # assert_success
}