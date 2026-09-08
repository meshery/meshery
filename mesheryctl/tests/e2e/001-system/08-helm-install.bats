#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    export HELM_CONFIG_HOME="$BATS_TEST_TMPDIR/helm_config"
    export HELM_CACHE_HOME="$BATS_TEST_TMPDIR/helm_cache"
    export HELM_DATA_HOME="$BATS_TEST_TMPDIR/helm_data"
    export TEST_NS="meshery-helm-test"
    export CREATED_NS=0
}

common_kube_and_helm_assertions() {
    if ! command -v helm >/dev/null 2>&1; then
        skip "helm is not installed"
    fi
    if ! kubectl cluster-info >/dev/null 2>&1; then
        skip "Kubernetes cluster not accessible"
    fi
}

teardown() {
    if [ "$CREATED_NS" -eq 1 ]; then
        helm uninstall meshery -n "$TEST_NS" >/dev/null 2>&1 || true
        kubectl delete ns "$TEST_NS" --timeout=30s >/dev/null 2>&1 || true
    fi
    rm -rf "$HELM_CONFIG_HOME" "$HELM_CACHE_HOME" "$HELM_DATA_HOME" >/dev/null 2>&1 || true
}

@test "[TC-11][cut=Installer / Helm Chart][tg=Installation] Install Meshery into a Kubernetes cluster via Helm chart" {
    common_kube_and_helm_assertions

    run kubectl create ns "$TEST_NS"
    assert_success
    CREATED_NS=1

    run helm repo add meshery https://meshery.io/charts/ --force-update
    assert_success

    run helm repo update meshery
    assert_success

    run helm install meshery meshery/meshery -n "$TEST_NS" --wait --timeout 3m
    assert_success
    assert_output --partial "STATUS: deployed"

    run kubectl rollout status deployment/meshery -n "$TEST_NS" --timeout=180s
    assert_success
}
