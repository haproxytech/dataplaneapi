#!/bin/sh
set -eu

case "$ACTION" in
append|set)
  printf '{"host":"%s","value":"%s"}' "$REC_NAME" "$REC_DATA" |
    curl -v -fsS -d @- "http://challtestsrv:8055/set-txt"
  break;;
get|delete)
  # not implemented
  break;;
*)
  echo "$0: error: invalid ACTION '$ACTION'" >&2
  exit 1;;
esac
