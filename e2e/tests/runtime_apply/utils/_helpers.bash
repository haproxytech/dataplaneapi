#!/usr/bin/env bash
#
# Copyright 2026 HAProxy Technologies
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

_BACKEND_BASE_PATH="/services/haproxy/configuration/backends"
_FRONTEND_BASE_PATH="/services/haproxy/configuration/frontends"
_GLOBAL_BASE_PATH="/services/haproxy/configuration/global"
_RAW_BASE_PATH="/services/haproxy/configuration/raw"
_TRANSACTIONS_BASE_PATH="/services/haproxy/transactions"
_RELOADS_BASE_PATH="/services/haproxy/reloads"
_RUNTIME_BACKEND_BASE_PATH="/services/haproxy/runtime/backends"
_RUNTIME_INFO_PATH="/services/haproxy/runtime/info"
_NATIVE_STATS_PATH="/services/haproxy/stats/native"

# sed script dropping "from mydefaults" from section headers: mydefaults is the
# defaults section written last, which a section without from inherits anyway
_DROP_FROM_MYDEFAULTS='s/^\(\(frontend\|backend\) [^ ]*\) from mydefaults$/\1/'

# dpa_request is dpa_curl that also reads the Reload-ID response header. The
# endpoint carries its own query string (version or transaction_id). The data
# file is relative to the test directory unless absolute.
# Sets SC, BODY and RELOAD_ID (empty when the response has no Reload-ID).
#
# Arguments:
# - HTTP verb
# - endpoint with query string
# - optional data file
# - optional content type (default application/json)
function dpa_request() {
  local verb=$1 endpoint=$2 data=${3:-/dev/null} ctype=${4:-application/json}
  [[ $data == /* ]] || data="${BATS_TEST_DIRNAME}/${data}"
  local headers="${BATS_TEST_TMPDIR}/headers"
  run curl -m 10 -s -D "$headers" -H "content-type: ${ctype}" --user dataplaneapi:mypassword \
    "-X${verb}" -w "\n%{http_code}" --data-binary "@${data}" "http://${LOCAL_IP_ADDRESS}:${E2E_PORT}${BASE_PATH}${endpoint}"
  assert_success
  SC=$(tail -n1 <<< "$output")
  BODY=$(head -n -1 <<< "$output")
  RELOAD_ID=$(awk 'tolower($1) == "reload-id:" { print $2 }' "$headers" | tr -d '\r')
  echo "> ${verb} ${endpoint}: ${SC}, Reload-ID '${RELOAD_ID}'"
}

# runtime_server_field echoes one live value of a server, as the runtime API
# reports it.
#
# Example:
# runtime_server_field be_plain srv1 weight
# >>> 10
function runtime_server_field() {
  dpa_curl GET "${_RUNTIME_BACKEND_BASE_PATH}/$1/servers/$2" | head -n -1 | jq -r ".$3"
}

# wait_reload_succeeded waits (up to 15s) for a reload to reach the status
# succeeded; it fails when the reload failed or did not finish in time.
function wait_reload_succeeded() {
  local i status
  for i in $(seq 1 150); do
    status=$(dpa_curl GET "${_RELOADS_BASE_PATH}/$1" | head -n -1 | jq -r '.status // empty')
    case "$status" in
      succeeded) return 0 ;;
      failed) echo "reload $1 failed" >&2; return 1 ;;
    esac
    sleep 0.1
  done
  echo "reload $1 did not succeed in time, last status: '${status}'" >&2
  return 1
}

# wait_new_worker waits (up to 10s) until a reload has replaced the HAProxy
# worker whose PID is given, and echoes the PID of the new worker. The e2e
# reload command returns, and the reload succeeds, before that happens.
function wait_new_worker() {
  local i pid
  for i in $(seq 1 100); do
    pid=$(haproxy_worker_pid) || return 1
    if [[ "$pid" != "$1" ]]; then
      echo "$pid"
      return 0
    fi
    sleep 0.1
  done
  echo "HAProxy worker $1 was not replaced" >&2
  return 1
}

# raw_config_with writes the current raw configuration, edited by the sed
# script given as argument, to a file and echoes its path.
function raw_config_with() {
  local file="${BATS_TEST_TMPDIR}/haproxy.cfg"
  dpa_curl GET "$_RAW_BASE_PATH" | head -n -1 | sed -e "$1" > "$file"
  echo "$file"
}

# new_transaction starts a transaction and sets TRANSACTION_ID.
function new_transaction() {
  get_version
  dpa_request POST "${_TRANSACTIONS_BASE_PATH}?version=${VERSION}"
  assert_equal "$SC" 201
  TRANSACTION_ID=$(get_json_path "$BODY" '.id')
}
