#!/usr/bin/env bats

setup() {
    INSTALLER="$BATS_TEST_TMPDIR/meshery-install.sh"
    HARNESS="$BATS_TEST_TMPDIR/platform-prompt.sh"

    curl -fsSL --connect-timeout 10 --max-time 60 https://meshery.io/install -o "$INSTALLER"

    awk '
    /^####### COMMON FUNCTIONS/ { found = 1; exit }
    { print }
    END { exit !found }
' "$INSTALLER" > "$HARNESS" || return 1

    cat >> "$HARNESS" <<'EOF'
printf 'SELECTED_PLATFORM=%s\n' "$PLATFORM"
EOF
}

@test "installer accepts docker platform selection" {
    run expect -c "
        set timeout 20
        spawn bash \"$HARNESS\"

        expect {
            \"Enter a platform to deploy Meshery\" {
                send \"docker\r\"
            }
            timeout {
                exit 1
            }
        }

        expect {
            \"SELECTED_PLATFORM=docker\" {
                exit 0
            }
            timeout {
                exit 1
            }
        }
    "

    [ "$status" -eq 0 ]
    [[ "$output" == *"Enter a platform to deploy Meshery"* ]]
    [[ "$output" == *"SELECTED_PLATFORM=docker"* ]]
}

@test "installer rejects invalid platform and accepts kubernetes" {
    run expect -c "
        set timeout 20
        spawn bash \"$HARNESS\"

        expect {
            \"Enter a platform to deploy Meshery\" {
                send \"invalid\r\"
            }
            timeout {
                exit 1
            }
        }

        expect {
            \"Invalid platform\" {
                send \"kubernetes\r\"
            }
            timeout {
                exit 1
            }
        }

        expect {
            \"SELECTED_PLATFORM=kubernetes\" {
                exit 0
            }
            timeout {
                exit 1
            }
        }
    "

    [ "$status" -eq 0 ]
    [[ "$output" == *"Invalid platform"* ]]
    [[ "$output" == *"SELECTED_PLATFORM=kubernetes"* ]]
}
