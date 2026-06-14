#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
}

@test "mesheryctl perf result displays results with correct headers and at least 2 entries" {
    run $MESHERYCTL_BIN perf result test-profile

    assert_success

    # Check column names
    assert_line --partial "NAME"
    assert_line --partial "MESH"
    assert_line --partial "QPS"
    assert_line --partial "DURATION"
    assert_line --partial "P50"
    assert_line --partial "P99.9"
    assert_line --partial "START-TIME"

    # Count lines that look like result rows (exclude header line)
    row_count=$(echo "$output" | tail -n +2 | grep -E '^\s*\w+' | wc -l)

    if [ "$row_count" -lt 2 ]; then
        echo "Expected at least 2 results, but got $row_count"
        return 1
    fi
}
