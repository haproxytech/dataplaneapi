#!/usr/bin/env bats
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

load '../../libs/dataplaneapi'
load '../../libs/get_json_path'
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_worker'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "runtime_apply: Change not applicable at runtime reloads" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain?version=$VERSION" data/be_plain_leastconn.json
  assert_equal "$SC" 202
  assert [ -n "$RELOAD_ID" ]

  wait_reload_succeeded "$RELOAD_ID"
  local new_pid; new_pid=$(wait_new_worker "$pid")
  assert [ "$new_pid" != "$pid" ]
  resource_get "$_NATIVE_STATS_PATH" "type=backend&name=be_plain"
  assert_equal "$(get_json_path "$BODY" '.stats[0].stats.algo')" "leastconn"
}

@test "runtime_apply: Runtime apply works again once the needed reload is done" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain?version=$VERSION" data/be_plain_leastconn.json
  assert_equal "$SC" 202
  assert [ -n "$RELOAD_ID" ]
  wait_reload_succeeded "$RELOAD_ID"
  pid=$(wait_new_worker "$pid")

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv2?version=$VERSION" data/srv2_weight.json
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""

  assert_equal "$(runtime_server_field be_plain srv2 weight)" 25
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: A runtime change made while a reload is taken over reaches HAProxy" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain?version=$VERSION" data/be_plain_leastconn.json
  assert_equal "$SC" 202
  local reload_id=$RELOAD_ID

  # no wait, not even for get_version: the old worker may still answer the runtime API
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv2?version=$(( VERSION + 1 ))" data/srv2_weight.json
  assert [ "$SC" = 200 -o "$SC" = 202 ]

  wait_reload_succeeded "$reload_id"
  if [[ -n "$RELOAD_ID" ]]; then
    wait_reload_succeeded "$RELOAD_ID"
  fi
  wait_new_worker "$pid"

  resource_get "$_BACKEND_BASE_PATH/be_plain/servers/srv2"
  assert_equal "$(get_json_path "$BODY" '.weight')" 25
  # a second reload may still be taken over: allow it 3s, a lost change never shows up
  local i weight
  for i in $(seq 1 30); do
    weight=$(runtime_server_field be_plain srv2 weight)
    [[ "$weight" == 25 ]] && break
    sleep 0.1
  done
  assert_equal "$weight" 25
}
