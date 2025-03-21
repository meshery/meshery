check_output(){
  local expected_message="$1"

  run grep -Fq "$expected_message" <<< "$actual_output"
  if [ "$status" -ne 0 ]; then
    echo "FAILED: Expected message not found: \"$expected_message\"" >&2
    echo "Actual Output: " >&2
    echo "$actual_output" >&2
  fi

  [ "$status" -eq 0 ]
}