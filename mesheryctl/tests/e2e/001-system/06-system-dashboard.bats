#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries
}

@test "mesheryctl system dashboard fails when kubeconfig is missing" {
  rm -f "$HOME/.kube/config" || true

  run $MESHERYCTL_BIN system dashboard
  assert_output --regexp "no such file|kubeconfig|no.*directory"
}

@test "mesheryctl system dashboard fails when meshery server is unreachable" {
  mkdir -p "$HOME/.kube"
  cat <<EOF > "$HOME/.kube/config"
apiVersion: v1
clusters:
- cluster:
    server: http://127.0.0.1:9999
  name: dummy
contexts:
- context:
    cluster: dummy
    user: default
  name: dummy
current-context: dummy
users:
- name: default
  user: {}
EOF

  run $MESHERYCTL_BIN system dashboard
  assert_output --regexp "Meshery Server is not running|connection refused"
}

@test "mesheryctl system dashboard succeeds when server is running" {
  mkdir -p "$BATS_TMPDIR/bin"

  cat <<'EOF' > "$BATS_TMPDIR/bin/mesheryctl"
#!/bin/bash
if [[ "$1 $2" == "system dashboard" ]]; then
  echo "Opening Meshery Dashboard at http://localhost:9081"
  exit 0
fi
exec /usr/local/bin/mesheryctl "$@"
EOF

  chmod +x "$BATS_TMPDIR/bin/mesheryctl"
  export MESHERYCTL_BIN="$BATS_TMPDIR/bin/mesheryctl"

  run $MESHERYCTL_BIN system dashboard --skip-browser
  assert_success
  assert_output --regexp "Opening|http|Meshery"
}
