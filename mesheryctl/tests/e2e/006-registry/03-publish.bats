#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/registry"
    mkdir -p "$TESTDATA_DIR"  
    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
    
    # Set the actual spreadsheet ID for tests (Meshery Integrations spreadsheet with flyte model)
    export TEST_SPREADSHEET_ID="1PjJv_kj5fB2gFPYokid5F8_epzYXWneFLh_eL8udLwY"
    export TEST_SPREADSHEET_CRED="ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAibWVzaGVyeS1tb2RlbHMiLAogICJwcml2YXRlX2tleV9pZCI6ICI1MGRlZGQ5ODg5Zjk2NTEyMzM3YTcwYjcwNTk2N2I4ZTgyNTI3YzRjIiwKICAicHJpdmF0ZV9rZXkiOiAiLS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tXG5NSUlFdkFJQkFEQU5CZ2txaGtpRzl3MEJBUUVGQUFTQ0JLWXdnZ1NpQWdFQUFvSUJBUUM4RDUzWlJQakIwaUZkXG5LajVMK1c5OW5qRGsvZUxMSFFVTk9yY21sNzlWMFNobXFKakdWMVAwUFh1RkIxWllOOHRQUjFsTmtXSXJiMCt6XG5Sc1Y5RGNpWjh4aCs2WHhYZHNlZnd6VHZXQXVBUHNVeW5DNHNMQmhKdVJOSG9yMzhaUXFmT3NueE1ZY1duY3U1XG5pQVU4UWh2LytsNGozM2hPYmhCYkFweHgwL1F0OFdseURrSTdUMVFXTmdBclh0WUk4L2xuTEMybTM3SW9VMWxhXG5nVmw3VE1NUFRvTS96WTdEajFYNFpkdkUyZXVwckl6VGltZER1RWhNY1d1TkR3Wkx0V0xsdko4bFJzcHhySW44XG5WcTNvaGJZMGZXQlk3Uk51dmZ4MllvSnJiakZNc2J3bW5WM3pJVUFNT1pSMVIrc2tDb1F1RVh3eDdDZGVOY2lMXG4rRUprRlhEbEFnTUJBQUVDZ2dFQUNlaVV1NGJieUhLRzJYN0pTMzEzOTFsbllCbFVZcmxEczkrOUdLWjBLVlJjXG4yWlBXZ05USEVwNU9yTUtuTlB5emxKakF4c3A3THFpdHZWZGtQbU9FaHovOHNXc2g1d1hXSEwvUFFURVZZZ3RzXG5ubThHSzhPaS8vZGxZYWZjZExqcUoxdVVta1d1eFcxSTk1eUlFaFcxbmZZV0ViUExoZTVjcEQ3dCtUMWVwekR5XG50VTlmRmRVZDBidmNsVEhqMnFGU2M5TkNqNVlHektkZ1VXM0tKRHVob01kWi9UbkZQY3d3SzgvQkJqekMzQzgwXG55aXBjSkVPcEgwSXVva0cxaDNzSFcxaWxjZmFGcytjNmVXK2lGb2xzUDRaaTkxK2ltOFJKTVc1RzdBQlRrOWdpXG5mSGJvblNpZytuRVNxaEFXVmt3VHBoMGROVnF5cElIMnRyRFBzUEZLL1FLQmdRRDdHUm5WZVQwUW93S3Mrdk1FXG5PYnZpeWl4RnZvbWJvNVVlU1VDTndQU213Nys5OHgxK2hKNCt1azJaRXVpZ1p1UEpJVmlSVXNMWFlST0hTSGZtXG4rL0I0SllWQzgvL3BhN1RpbjlxOHpLbnZ5Ukk3MmQzcjU0SkxrZlZFeFI1bE4xcWx3akhiRmdJbmtZdW9FTktvXG5mTVlnelF3eHhHM0ppTFNPWStpUEV3aFhad0tCZ1FDL3UzcVdDYitOQWFUaHJkUU4rS1pBZlpZVTY4YjQ0d2k5XG43NmtRK0xZZzNILzEwVVBjWVloVjQ3bS9CcC9oL1FQYXpMbWljV0dVbGdMSWxPNVRNY2FpWFNyRnBYUkd5S2ozXG5Bc1RIOHd0alRPYUw3KzFzT1lzK2w0S2tBNkl5TzhQbDczbHhncWVUQmxvMm1FeWFWS1Z3Zk1RWll6SUJFbFBXXG5pa0YrUUhJQjB3S0JnRk94Wmg3MUxobGhlLzdVRmpRT1drV1pSWEJwdGN3ZWxoRm1NUklPZ24wak9jSENTczVrXG52YW9DeHVLQ251NzMzR3lwcEsrcmVpMzN1VjZhRlB6MGRJNVhVemVoeEdhMDg3R2x2bTYzZU56bmxDQkcyZVZMXG5ZTzUzVTl6OEpzbHR6ZFQ2Q1V3UTVkN3ljQzR3alVhRGNEWjJvZDNrcndOQmdZMEVWL3Q0U0l0NUFvR0FlUm0vXG5XZGthWmNPdUV1VXFWUkRDWVdWeFVCR0VBVnZVLzdMS1JHL3pDVnA4Y0JVWEdiRjQzVmNTMFRjcjhSOVFPZmlwXG5INU1ob1NvRzVBNDlZY3JFdlhTL09VdEI0QnVKclVyUVk4SjlaT0VEbFV1YjJJbVlDNWJDNnpsb1R2WTB5NXl3XG5SaVBTR0NrZy9lN0FJNUpuVkY5WnRQZFVkQnpSVzJiRDZ5eGxWZ1VDZ1lCYTFkS3JmblNSZ3JLQ2JFbXB4YmFkXG5ic0hGUnNzWFdHMlVDNHNEeWsrRitIRnRWMUptVjhyNklQSCtNRnV1MzNJaDl3OGxrV29razRGTUk3VHJheDBQXG5NRlYzRWdjd1hjTmpVWE5QKytvNno3VzZtQ2VTZmlzZUxtbGdobS83dXlSZ2pOMEtwQy9iR2xjdVZwMVplVFZnXG5nMFhiNG9Yc1JjSkFSMDdtSFMwS3NnPT1cbi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS1cbiIsCiAgImNsaWVudF9lbWFpbCI6ICJkYXJzaGFuMTc0QG1lc2hlcnktbW9kZWxzLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAiY2xpZW50X2lkIjogIjEwOTgxMjI4NTU2MjI1MTc5Mjc4NSIsCiAgImF1dGhfdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwKICAidG9rZW5fdXJpIjogImh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwKICAiYXV0aF9wcm92aWRlcl94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL29hdXRoMi92MS9jZXJ0cyIsCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZGFyc2hhbjE3NCU0MG1lc2hlcnktbW9kZWxzLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAidW5pdmVyc2VfZG9tYWluIjogImdvb2dsZWFwaXMuY29tIgp9Cg=="

    export TEST_MODELS_OUTPUT="$TESTDATA_DIR/models-output"
    export TEST_IMGS_OUTPUT="$TESTDATA_DIR/imgs-output"
    mkdir -p "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
}

common_assertions() {
    assert_output --partial "[ system, google sheet credential, sheet-id, models output path, imgs output path] are required"
}

@test "mesheryctl registry publish displays usage instructions when no arguments are provided" {
    run $MESHERYCTL_BIN registry publish 
    assert_failure
    common_assertions
    assert_output --partial "Usage:"
    assert_output --partial "mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]"
}

@test "mesheryctl registry publish fails when insufficient arguments are provided" {
    run $MESHERYCTL_BIN registry publish website
    assert_failure
    common_assertions
}

@test "mesheryctl registry publish fails when only 4 arguments are provided" {
    run $MESHERYCTL_BIN registry publish website cred id path
    assert_failure
    common_assertions
}

@test "mesheryctl registry publish fails with invalid credentials" {
    run $MESHERYCTL_BIN registry publish website "invalid-cred" "invalid-id" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
    assert_failure
    assert_output --partial "Invalid JWT Token: Ensure the provided token is a base64-encoded, valid Google Spreadsheets API token"
}

@test "mesheryctl registry publish displays help message" {
    run $MESHERYCTL_BIN registry publish --help
    assert_success
    assert_output --partial "Publishes metadata about Meshery Models"
    assert_output --partial "Websites, Remote Provider, or Meshery Server"
    assert_output --partial "Google Spreadsheet"
    assert_output --partial "Usage:"
    assert_output --partial "mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]"
    assert_output --partial "Examples:"
    assert_output --partial "Publish To System"
    assert_output --partial "Publish To Meshery"
    assert_output --partial "Publish To Remote Provider"
    assert_output --partial "Publish To Website"
    assert_output --partial "--output-format"
}

@test "mesheryctl registry publish succeeds with meshery system" {
    run $MESHERYCTL_BIN registry publish meshery "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
    assert_success
}

@test "mesheryctl registry publish succeeds with remote-provider system" {
    run $MESHERYCTL_BIN registry publish remote-provider "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
    assert_success
}

@test "mesheryctl registry publish succeeds with website system and md output format" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o md
    assert_success
}

@test "mesheryctl registry publish succeeds with website system and mdx output format" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o mdx
    assert_success
}

@test "mesheryctl registry publish succeeds with website system and js output format" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o js
    assert_success
}

@test "mesheryctl registry publish fails with website system and invalid output format" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o invalid
    assert_failure
    assert_output --partial "invalid output format: invalid"
}

@test "mesheryctl registry publish creates output directories if they don't exist" {
    local new_models_dir="$TESTDATA_DIR/new-models"
    local new_imgs_dir="$TESTDATA_DIR/new-imgs"
    
    run $MESHERYCTL_BIN registry publish remote-provider "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$new_models_dir" "$new_imgs_dir"
    assert_success
}

@test "mesheryctl registry publish handles non-existent output paths gracefully" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "/non/existent/models" "/non/existent/imgs" -o md
    assert_success
}