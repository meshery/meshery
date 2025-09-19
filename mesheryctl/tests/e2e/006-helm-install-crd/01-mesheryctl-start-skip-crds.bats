#!/usr/bin/env bats

load "$BATS_LIB_PATH/bats-support/load.bash"
load "$BATS_LIB_PATH/bats-assert/load.bash"

MESHERY_CRDS=(
  "brokers.meshery.io"
  "meshsyncs.meshery.io"
)

setup() {
  echo "--- Setup: Cleaning CRDs and old Meshery installs ---"

  # Delete Meshery CRDs (if any)
  for crd in "${MESHERY_CRDS[@]}"; do
    run kubectl delete crd "$crd" --ignore-not-found=true
  done

  # Stop Meshery if running
  run mesheryctl system stop --yes || true

  # Delete Meshery namespace
  run kubectl delete ns meshery --timeout=2m --force --grace-period=0 --ignore-not-found=true || true
}

teardown() {
  echo "--- Teardown: Cleanup after test ---"
  run mesheryctl system stop --yes || true
  run kubectl delete ns meshery --timeout=2m --force --grace-period=0 || true
}

@test "mesheryctl start with --skip-crds should NOT install CRDs" {
  # Start Meshery with skip-crds
  run mesheryctl system start --platform kubernetes --skip-crds --yes
  assert_success "mesheryctl system start failed with --skip-crds"

  # Allow some time for installation
  sleep 20

  # Assert CRDs are NOT present
  for crd in "${MESHERY_CRDS[@]}"; do
    run kubectl get crd "$crd" 2>/dev/null
    assert_failure "CRD $crd should NOT be installed when --skip-crds is used"
  done
}
