#!/usr/bin/env bats

setup() {
    load "$SUPPORT_DESTDIR/load"
    load "$ASSERT_DESTDIR/load"
    load "$DETIK_DESTDIR/detik.bash"

    DETIK_CLIENT_NAMESPACE="meshery"
    DETIK_CLIENT_NAME="kubectl"
}

@test "mesheryctl system start on kubernetes is succeeded" {
    run $MESHERYCTL_BIN system start -p kubernetes
    [ "$status" -eq 0 ]

	run verify "there is 1 service named 'meshery'"
	[ "$status" -eq 0 ]
}

