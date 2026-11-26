#!/usr/bin/env bash

append_log_excerpt() {
  local LOG_FILE="$1"
  local LOG_LABEL="$2"
  local NORMALIZED_LOG="${LOG_FILE}.normalized"
  local EXCERPT_FILE="${LOG_FILE}.excerpt"
  local START_LINE_COUNT CURRENT_LINE_COUNT LINES_USED REMAINING_LINES TAIL_LINES GREP_STATUS
  local MATCH_LINE_COUNT MATCH_LINE_LIMIT MAX_EXCERPT_LINES=100 MAX_TAIL_LINES=50
  # Include the trailing-context heading and blank line in the reservation.
  local TAIL_RESERVE_LINES=$((MAX_TAIL_LINES + 2))
  START_LINE_COUNT="$(wc -l < "$EMAIL_FILE")"

  if ! sed 's/\x1b\[[0-9;]*m//g' "$LOG_FILE" | tr -d '\r' > "$NORMALIZED_LOG"; then
    echo "Could not normalize the ${LOG_LABEL}; the complete log is attached." >> "$EMAIL_FILE"
    return 0
  fi

  echo "Log excerpt for ${LOG_LABEL}:" >> "$EMAIL_FILE"
  if grep -niE -C 3 \
    'error|fail(ed|ure)?|fatal|panic|denied|timed out|could not|unable to|exit status|cannot|not found|killed|unexpected|invalid|refused' \
    "$NORMALIZED_LOG" > "$EXCERPT_FILE"; then
    echo "Relevant error/failure lines with surrounding context:" >> "$EMAIL_FILE"
    CURRENT_LINE_COUNT="$(wc -l < "$EMAIL_FILE")"
    LINES_USED=$((CURRENT_LINE_COUNT - START_LINE_COUNT))
    REMAINING_LINES=$((MAX_EXCERPT_LINES - LINES_USED))
    MATCH_LINE_COUNT="$(wc -l < "$EXCERPT_FILE")"
    MATCH_LINE_LIMIT="$REMAINING_LINES"
    if [ "$REMAINING_LINES" -gt "$TAIL_RESERVE_LINES" ] &&
      [ "$MATCH_LINE_COUNT" -gt "$((REMAINING_LINES - TAIL_RESERVE_LINES))" ]; then
      MATCH_LINE_LIMIT=$((REMAINING_LINES - TAIL_RESERVE_LINES))
    fi
    if [ "$MATCH_LINE_LIMIT" -gt 0 ]; then
      head -n "$MATCH_LINE_LIMIT" "$EXCERPT_FILE" >> "$EMAIL_FILE"
    fi
  else
    GREP_STATUS=$?
    if [ "$GREP_STATUS" -eq 1 ]; then
      echo "No matching error/failure lines found in this log." >> "$EMAIL_FILE"
    else
      echo "Could not search this log for error/failure lines (grep exit ${GREP_STATUS})." >> "$EMAIL_FILE"
    fi
  fi

  CURRENT_LINE_COUNT="$(wc -l < "$EMAIL_FILE")"
  LINES_USED=$((CURRENT_LINE_COUNT - START_LINE_COUNT))
  REMAINING_LINES=$((MAX_EXCERPT_LINES - LINES_USED))
  if [ "$REMAINING_LINES" -gt 2 ]; then
    echo "Additional context from the end of the log:" >> "$EMAIL_FILE"
    CURRENT_LINE_COUNT="$(wc -l < "$EMAIL_FILE")"
    LINES_USED=$((CURRENT_LINE_COUNT - START_LINE_COUNT))
    REMAINING_LINES=$((MAX_EXCERPT_LINES - LINES_USED))
    TAIL_LINES=$((REMAINING_LINES - 1))
    if [ "$TAIL_LINES" -gt "$MAX_TAIL_LINES" ]; then
      TAIL_LINES="$MAX_TAIL_LINES"
    fi
    if [ "$TAIL_LINES" -gt 0 ]; then
      if tail -n "$TAIL_LINES" "$NORMALIZED_LOG" > "${EXCERPT_FILE}.tail"; then
        cat "${EXCERPT_FILE}.tail" >> "$EMAIL_FILE"
      else
        echo "Could not read the last lines of this log." >> "$EMAIL_FILE"
      fi
    fi
    echo "" >> "$EMAIL_FILE"
  fi
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  EMAIL_FILE="${1:?EMAIL_FILE is required}"
  LOG_FILE="${2:?LOG_FILE is required}"
  LOG_LABEL="${3:-test step}"
  append_log_excerpt "$LOG_FILE" "$LOG_LABEL"
fi
