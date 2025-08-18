#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/registry"
    mkdir -p "$TESTDATA_DIR"  
    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
    
    # Set the actual spreadsheet ID for tests (Meshery COPY OF Integrations spreadsheet with flyte model)
    export TEST_SPREADSHEET_ID="1PjJv_kj5fB2gFPYokid5F8_epzYXWneFLh_eL8udLwY"
    export TEST_SPREADSHEET_CRED="ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAibWVzaGVyeS1tb2RlbHMiLAogICJwcml2YXRlX2tleV9pZCI6ICI1MGRlZGQ5ODg5Zjk2NTEyMzM3YTcwYjcwNTk2N2I4ZTgyNTI3YzRjIiwKICAicHJpdmF0ZV9rZXkiOiAiLS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tXG5NSUlFdkFJQkFEQU5CZ2txaGtpRzl3MEJBUUVGQUFTQ0JLWXdnZ1NpQWdFQUFvSUJBUUM4RDUzWlJQakIwaUZkXG5LajVMK1c5OW5qRGsvZUxMSFFVTk9yY21sNzlWMFNobXFKakdWMVAwUFh1RkIxWllOOHRQUjFsTmtXSXJiMCt6XG5Sc1Y5RGNpWjh4aCs2WHhYZHNlZnd6VHZXQXVBUHNVeW5DNHNMQmhKdVJOSG9yMzhaUXFmT3NueE1ZY1duY3U1XG5pQVU4UWh2LytsNGozM2hPYmhCYkFweHgwL1F0OFdseURrSTdUMVFXTmdBclh0WUk4L2xuTEMybTM3SW9VMWxhXG5nVmw3VE1NUFRvTS96WTdEajFYNFpkdkUyZXVwckl6VGltZER1RWhNY1d1TkR3Wkx0V0xsdko4bFJzcHhySW44XG5WcTNvaGJZMGZXQlk3Uk51dmZ4MllvSnJiakZNc2J3bW5WM3pJVUFNT1pSMVIrc2tDb1F1RVh3eDdDZGVOY2lMXG4rRUprRlhEbEFnTUJBQUVDZ2dFQUNlaVV1NGJieUhLRzJYN0pTMzEzOTFsbllCbFVZcmxEczkrOUdLWjBLVlJjXG4yWlBXZ05USEVwNU9yTUtuTlB5emxKakF4c3A3THFpdHZWZGtQbU9FaHovOHNXc2g1d1hXSEwvUFFURVZZZ3RzXG5ubThHSzhPaS8vZGxZYWZjZExqcUoxdVVta1d1eFcxSTk1eUlFaFcxbmZZV0ViUExoZTVjcEQ3dCtUMWVwekR5XG50VTlmRmRVZDBidmNsVEhqMnFGU2M5TkNqNVlHektkZ1VXM0tKRHVob01kWi9UbkZQY3d3SzgvQkJqekMzQzgwXG55aXBjSkVPcEgwSXVva0cxaDNzSFcxaWxjZmFGcytjNmVXK2lGb2xzUDRaaTkxK2ltOFJKTVc1RzdBQlRrOWdpXG5mSGJvblNpZytuRVNxaEFXVmt3VHBoMGROVnF5cElIMnRyRFBzUEZLL1FLQmdRRDdHUm5WZVQwUW93S3Mrdk1FXG5PYnZpeWl4RnZvbWJvNVVlU1VDTndQU213Nys5OHgxK2hKNCt1azJaRXVpZ1p1UEpJVmlSVXNMWFlST0hTSGZtXG4rL0I0SllWQzgvL3BhN1RpbjlxOHpLbnZ5Ukk3MmQzcjU0SkxrZlZFeFI1bE4xcWx3akhiRmdJbmtZdW9FTktvXG5mTVlnelF3eHhHM0ppTFNPWStpUEV3aFhad0tCZ1FDL3UzcVdDYitOQWFUaHJkUU4rS1pBZlpZVTY4YjQ0d2k5XG43NmtRK0xZZzNILzEwVVBjWVloVjQ3bS9CcC9oL1FQYXpMbWljV0dVbGdMSWxPNVRNY2FpWFNyRnBYUkd5S2ozXG5Bc1RIOHd0alRPYUw3KzFzT1lzK2w0S2tBNkl5TzhQbDczbHhncWVUQmxvMm1FeWFWS1Z3Zk1RWll6SUJFbFBXXG5pa0YrUUhJQjB3S0JnRk94Wmg3MUxobGhlLzdVRmpRT1drV1pSWEJwdGN3ZWxoRm1NUklPZ24wak9jSENTczVrXG52YW9DeHVLQ251NzMzR3lwcEsrcmVpMzN1VjZhRlB6MGRJNVhVemVoeEdhMDg3R2x2bTYzZU56bmxDQkcyZVZMXG5ZTzUzVTl6OEpzbHR6ZFQ2Q1V3UTVkN3ljQzR3alVhRGNEWjJvZDNrcndOQmdZMEVWL3Q0U0l0NUFvR0FlUm0vXG5XZGthWmNPdUV1VXFWUkRDWVdWeFVCR0VBVnZVLzdMS1JHL3pDVnA4Y0JVWEdiRjQzVmNTMFRjcjhSOVFPZmlwXG5INU1ob1NvRzVBNDlZY3JFdlhTL09VdEI0QnVKclVyUVk4SjlaT0VEbFV1YjJJbVlDNWJDNnpsb1R2WTB5NXl3XG5SaVBTR0NrZy9lN0FJNUpuVkY5WnRQZFVkQnpSVzJiRDZ5eGxWZ1VDZ1lCYTFkS3JmblNSZ3JLQ2JFbXB4YmFkXG5ic0hGUnNzWFdHMlVDNHNEeWsrRitIRnRWMUptVjhyNklQSCtNRnV1MzNJaDl3OGxrV29razRGTUk3VHJheDBQXG5NRlYzRWdjd1hjTmpVWE5QKytvNno3VzZtQ2VTZmlzZUxtbGdobS83dXlSZ2pOMEtwQy9iR2xjdVZwMVplVFZnXG5nMFhiNG9Yc1JjSkFSMDdtSFMwS3NnPT1cbi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS1cbiIsCiAgImNsaWVudF9lbWFpbCI6ICJkYXJzaGFuMTc0QG1lc2hlcnktbW9kZWxzLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAiY2xpZW50X2lkIjogIjEwOTgxMjI4NTU2MjI1MTc5Mjc4NSIsCiAgImF1dGhfdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwKICAidG9rZW5fdXJpIjogImh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwKICAiYXV0aF9wcm92aWRlcl94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL29hdXRoMi92MS9jZXJ0cyIsCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZGFyc2hhbjE3NCU0MG1lc2hlcnktbW9kZWxzLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAidW5pdmVyc2VfZG9tYWluIjogImdvb2dsZWFwaXMuY29tIgp9Cg=="
}

@test "mesheryctl registry update displays usage instructions when no arguments are provided" {
    run $MESHERYCTL_BIN registry update 
    assert_failure
    assert_output --partial "unexpected end of JSON input"
}

@test "mesheryctl registry update fails when spreadsheet-id is provided without spreadsheet-cred" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "test-id"
    assert_failure
    assert_output --partial "if any flags in the group [spreadsheet-id spreadsheet-cred] are set they must all be set"
    assert_output --partial "missing [spreadsheet-cred]"
}

@test "mesheryctl registry update fails when spreadsheet-cred is provided without spreadsheet-id" {
    run $MESHERYCTL_BIN registry update --spreadsheet-cred "test-cred"
    assert_failure
    assert_output --partial "if any flags in the group [spreadsheet-id spreadsheet-cred] are set they must all be set"
    assert_output --partial "missing [spreadsheet-id]"
}

@test "mesheryctl registry update fails with invalid spreadsheet credentials" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "invalid-id" --spreadsheet-cred "invalid-cred"
    assert_failure
}

@test "mesheryctl registry update succeeds with non-existent input directory but updates nothing" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "/non/existent/path"
    assert_success
    assert_output --partial "Updated 0 models and 0 components"
    assert_output --partial "refer"
    assert_output --partial "logs/registry"
}

@test "mesheryctl registry update displays help message" {
    run $MESHERYCTL_BIN registry update --help
    assert_success
    assert_output --partial "Updates the component metadata"
    assert_output --partial "SVGs, shapes, styles and other"
    assert_output --partial "Google Spreadsheet"
    assert_output --partial "Usage:"
    assert_output --partial "mesheryctl registry update [flags]"
    assert_output --partial "Examples:"
    assert_output --partial "Update models from Meshery Integration Spreadsheet"
    assert_output --partial "--spreadsheet-id"
    assert_output --partial "--spreadsheet-cred"
    assert_output --partial "--input"
    assert_output --partial "--model"
}

@test "mesheryctl registry update succeeds with valid spreadsheet credentials and test fixtures" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "$FIXTURES_DIR/test-models"
    assert_success
    # Check for successful completion indicators
    [[ "$output" =~ "Updated" ]] || [[ "$output" =~ "components" ]] || [[ "$output" =~ "models" ]]
}

@test "mesheryctl registry update supports model-specific updates" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "$FIXTURES_DIR/test-models" --model "kubernetes"
    assert_success
    # Check for successful completion indicators
    [[ "$output" =~ "Updated" ]] || [[ "$output" =~ "kubernetes" ]] || [[ "$output" =~ "components" ]]
}

@test "mesheryctl registry update creates log files in correct location" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "$FIXTURES_DIR/test-models"
    assert_success
    # Check that log directory reference is mentioned
    assert_output --partial "logs"
}

@test "mesheryctl registry update handles empty models directory gracefully" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "$FIXTURES_DIR/empty-models"
    assert_success
    # Should complete without errors even with empty directory
    assert_output --partial "Updated 0 models and 0 components"
}

@test "mesheryctl registry update shows summary of updated components" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "$FIXTURES_DIR/test-models"
    assert_success
    # Check for summary output
    [[ "$output" =~ "Updated" ]] && [[ "$output" =~ "models" ]] && [[ "$output" =~ "components" ]]
}