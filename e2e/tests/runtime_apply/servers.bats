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

@test "runtime_apply: Replace a server weight without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv2?version=$VERSION" data/srv2_weight.json
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""

  assert_equal "$(runtime_server_field be_plain srv2 weight)" 25
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Replace a server address and port without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv2?version=$VERSION" data/srv2_address.json
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""

  resource_get "$_RUNTIME_BACKEND_BASE_PATH/be_plain/servers/srv2"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.address')" "10.0.1.22"
  assert_equal "$(get_json_path "$BODY" '.port')" 8082
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Put a server into maintenance and back without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv2?version=$VERSION" data/srv2_maintenance.json
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""
  assert_equal "$(runtime_server_field be_plain srv2 admin_state)" "maint"

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv2?version=$VERSION" data/srv2.json
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""
  assert_equal "$(runtime_server_field be_plain srv2 admin_state)" "ready"

  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Add a server without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request POST "$_BACKEND_BASE_PATH/be_plain/servers?version=$VERSION" data/new_server.json
  assert_equal "$SC" 201
  assert_equal "$RELOAD_ID" ""

  resource_get "$_RUNTIME_BACKEND_BASE_PATH/be_plain/servers/srv_new"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.address')" "10.0.9.9"
  assert_equal "$(get_json_path "$BODY" '.port')" 8090
  assert_equal "$(get_json_path "$BODY" '.admin_state')" "ready"
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Add a server inheriting the backend default-server without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request POST "$_BACKEND_BASE_PATH/be_defsrv/servers?version=$VERSION" data/new_server.json
  assert_equal "$SC" 201
  assert_equal "$RELOAD_ID" ""

  # default-server check inter 2s weight 7
  resource_get "$_RUNTIME_BACKEND_BASE_PATH/be_defsrv/servers/srv_new"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.weight')" 7
  assert_equal "$(get_json_path "$BODY" '.admin_state')" "ready"
  # health check configured (0x2) and enabled (0x4)
  local check_state; check_state=$(get_json_path "$BODY" '.check_state')
  assert_equal "$(( check_state & 6 ))" 6
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Add a server inheriting the defaults default-server without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request POST "$_BACKEND_BASE_PATH/be_chain/servers?version=$VERSION" data/new_server.json
  assert_equal "$SC" 201
  assert_equal "$RELOAD_ID" ""

  # defaults d_chain: default-server check port 8999 weight 13
  resource_get "$_RUNTIME_BACKEND_BASE_PATH/be_chain/servers/srv_new"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.weight')" 13
  assert_equal "$(get_json_path "$BODY" '.check_port')" 8999
  local check_state; check_state=$(get_json_path "$BODY" '.check_state')
  assert_equal "$(( check_state & 6 ))" 6
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Delete a server without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request DELETE "$_BACKEND_BASE_PATH/be_plain/servers/srv3?version=$VERSION"
  assert_equal "$SC" 204
  assert_equal "$RELOAD_ID" ""

  resource_get "$_RUNTIME_BACKEND_BASE_PATH/be_plain/servers/srv3"
  assert_equal "$SC" 404
  resource_get "$_BACKEND_BASE_PATH/be_plain/servers/srv3"
  assert_equal "$SC" 404
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Add a server inheriting the implicit defaults default-server without a reload" {
  local pid; pid=$(haproxy_worker_pid)
  # be_plain without from inherits mydefaults, the defaults section written last
  local cfg; cfg=$(raw_config_with "${_DROP_FROM_MYDEFAULTS};/^defaults mydefaults$/a\  default-server check weight 17")
  assert grep -q '^backend be_plain$' "$cfg"
  assert grep -q '^  default-server check weight 17$' "$cfg"

  get_version
  dpa_request POST "$_RAW_BASE_PATH?version=$VERSION" "$cfg" text/plain
  # a new default-server changes the existing servers: it may need a reload
  if [[ "$SC" == 202 ]]; then
    wait_reload_succeeded "$RELOAD_ID"
    pid=$(wait_new_worker "$pid")
  else
    assert_equal "$SC" 201
  fi
  resource_get "$_BACKEND_BASE_PATH/be_plain"
  assert_equal "$(get_json_path "$BODY" '.from // ""')" ""

  get_version
  dpa_request POST "$_BACKEND_BASE_PATH/be_plain/servers?version=$VERSION" data/new_server.json
  assert_equal "$SC" 201
  assert_equal "$RELOAD_ID" ""

  resource_get "$_RUNTIME_BACKEND_BASE_PATH/be_plain/servers/srv_new"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.weight')" 17
  # health check configured (0x2) and enabled (0x4)
  local check_state; check_state=$(get_json_path "$BODY" '.check_state')
  assert_equal "$(( check_state & 6 ))" 6
  assert_equal "$(haproxy_worker_pid)" "$pid"
}
