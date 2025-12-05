#!/usr/bin/env bats

load "$BATS_LIB_PATH/bats-support/load.bash"
load "$BATS_LIB_PATH/bats-assert/load.bash"

# ------------------------------------------------------------------ #
#  Global                                                            #
# ------------------------------------------------------------------ #
MESHERY_CRDS=(
  "brokers.meshery.io"
  "meshsyncs.meshery.io"
)

# ------------------------------------------------------------------ #
#  Setup / Teardown                                                  #
# ------------------------------------------------------------------ #
setup() {
  echo "--- Global Setup: Cleaning up before test ---"

  # 1. Delete CRDs from any previous Meshery install
  for crd in "${MESHERY_CRDS[@]}"; do
    run kubectl delete crd "$crd" --ignore-not-found=true
  done

  # 2. Wait until CRDs are really gone (max 60 s)
  for crd in "${MESHERY_CRDS[@]}"; do
    for _ in {1..30}; do
      run kubectl get crd "$crd" 2>/dev/null
      [[ "$status" -ne 0 ]] && break
      sleep 2
    done
    run kubectl get crd "$crd" 2>/dev/null
    [[ "$status" -eq 0 ]] && fail "CRD $crd still present after cleanup"
  done

  # 3. Delete any leftovers from earlier test runs
  run helm uninstall meshery-e2e-no-crds     --namespace meshery-e2e-test-ns-no-crds     --ignore-not-found=true
  run helm uninstall meshery-e2e-with-crds   --namespace meshery-e2e-test-ns-default     --ignore-not-found=true
  run kubectl delete ns meshery-e2e-test-ns-no-crds --timeout=2m --force --grace-period=0 --ignore-not-found=true
  run kubectl delete ns meshery-e2e-test-ns-default --timeout=2m --force --grace-period=0 --ignore-not-found=true

  echo "--- Global Setup complete ---"
}

teardown() {
  echo "--- Teardown complete ---"
}

# ------------------------------------------------------------------ #
#  Test 1 – install WITHOUT CRDs                                     #
# ------------------------------------------------------------------ #
@test "helm install with installCRDs=false (and --skip-crds) should NOT install Meshery CRDs" {
  local release="meshery-e2e-no-crds"
  local ns="meshery-e2e-test-ns-no-crds"

  run helm install "$release" ./install/kubernetes/helm/meshery \
        --namespace "$ns" --create-namespace \
        --set installCRDs=false \
        --skip-crds \
        --wait --timeout 10m
  assert_success "helm install with --skip-crds failed"

  # CRDs must NOT be present
  for crd in "${MESHERY_CRDS[@]}"; do
    run kubectl get crd "$crd" 2>/dev/null
    assert_failure "CRD $crd should NOT exist"
  done

  # Cleanup
  run helm uninstall "$release" --namespace "$ns" || true
  run kubectl delete ns "$ns" --timeout=2m --force --grace-period=0 || true
}

# ------------------------------------------------------------------ #
#  Test 2 – default install WITH CRDs                                #
# ------------------------------------------------------------------ #
@test "default helm install SHOULD install Meshery CRDs" {
  local release="meshery-e2e-with-crds"
  local ns="meshery-e2e-test-ns-default"

  run helm install "$release" ./install/kubernetes/helm/meshery \
        --namespace "$ns" --create-namespace \
        --wait --timeout 10m
  assert_success "default helm install failed"

  # CRDs MUST be present
  for crd in "${MESHERY_CRDS[@]}"; do
    # retry up to 1 min for slower clusters
    for _ in {1..30}; do
      run kubectl get crd "$crd" 2>/dev/null
      [[ "$status" -eq 0 ]] && break
      sleep 2
    done
    assert_success "CRD $crd not found after install"
  done

  # Cleanup
  run helm uninstall "$release" --namespace "$ns" || true
  run kubectl delete ns "$ns" --timeout=2m --force --grace-period=0 || true
}
