#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

   export CHECK_KUBERNETES_API_HEADER="Kubernetes API"
   export CHECK_KUBERNETES_VERSION_HEADER="Kubernetes Version"
   export CHECK_MESHERY_COMPONENTS_HEADER="Meshery Components"
   export CHECK_MESHERY_OPERATORS_HEADER="Meshery Operators"
   export CHECK_DOCKER_HEADER="Docker"
   export CHECK_PREREQUISISTE_RESULT="Meshery prerequisites met"
}

@test "mesheryctl system check succeeds displaying required sections" {
   run $MESHERYCTL_BIN system check
   assert_success

   assert_output --partial "$CHECK_KUBERNETES_API_HEADER"
   assert_output --partial "$CHECK_KUBERNETES_VERSION_HEADER"
   assert_output --partial "$CHECK_MESHERY_COMPONENTS_HEADER"
   assert_output --partial "$CHECK_MESHERY_OPERATORS_HEADER"
}


@test "mesheryctl system check --pre succeeds displaying required sections and expected prerequisites result" {
   run $MESHERYCTL_BIN system check --pre
   assert_success

   assert_output --partial "$CHECK_DOCKER_HEADER"
   assert_output --partial "$CHECK_KUBERNETES_API_HEADER"
   assert_output --partial "$CHECK_KUBERNETES_VERSION_HEADER"
   assert_output --partial "$CHECK_PREREQUISISTE_RESULT"
}

@test "mesheryctl system check --preflight succeeds displaying required sections and expected prerequisites result" {
   run $MESHERYCTL_BIN system check --pre
   assert_success

   assert_output --partial "$CHECK_DOCKER_HEADER"
   assert_output --partial "$CHECK_KUBERNETES_API_HEADER"
   assert_output --partial "$CHECK_KUBERNETES_VERSION_HEADER"
   assert_output --partial "$CHECK_PREREQUISISTE_RESULT"
}

@test "mesheryctl system check --operators succeeds displaying required section and meshery operators in running state" {
   run $MESHERYCTL_BIN system check --operator
   assert_success

   assert_output --partial "$CHECK_MESHERY_OPERATORS_HEADER"
   assert_output --partial "Meshery Operator is running"
   assert_output --partial "Meshery Broker is running"
   assert_output --partial "Meshsync is running"
}