#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    export HELM_CONFIG_HOME="$BATS_TEST_TMPDIR/helm_config"
    export HELM_CACHE_HOME="$BATS_TEST_TMPDIR/helm_cache"
    export HELM_DATA_HOME="$BATS_TEST_TMPDIR/helm_data"
    export TEST_NS="meshery-helm-$RANDOM"
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
        kubectl delete ns "$TEST_NS" --timeout=60s --wait=true >/dev/null 2>&1 || true
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

    # Validate deployment rollout
    run kubectl rollout status deployment/meshery -n "$TEST_NS" --timeout=180s
    assert_success

    # Validate Meshery pod is running and ready (readiness probe /healthz/ready passed)
    run kubectl get pods -n "$TEST_NS" -l app.kubernetes.io/name=meshery -o jsonpath='{.items[0].status.phase}'
    assert_success
    assert_equal "$output" "Running"

    run kubectl get pods -n "$TEST_NS" -l app.kubernetes.io/name=meshery -o jsonpath='{.items[0].status.conditions[?(@.type=="Ready")].status}'
    assert_success
    assert_equal "$output" "True"

    # Validate Meshery service is created and endpoints are attached
    run kubectl get svc meshery -n "$TEST_NS"
    assert_success

    run kubectl get endpoints meshery -n "$TEST_NS" -o jsonpath='{.subsets[*].addresses[*].ip}'
    assert_success
    [ -n "$output" ]
}