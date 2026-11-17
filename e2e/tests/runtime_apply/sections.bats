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

@test "runtime_apply: Replace a full backend section changing a server weight without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain?full_section=true&version=$VERSION" data/be_plain_full_section.json
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""

  assert_equal "$(runtime_server_field be_plain srv3 weight)" 35
  assert_equal "$(runtime_server_field be_plain srv2 weight)" 20
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Replace a frontend maxconn without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_FRONTEND_BASE_PATH/fe_main?version=$VERSION" data/fe_main_maxconn.json
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""

  resource_get "$_NATIVE_STATS_PATH" "type=frontend&name=fe_main"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.stats[0].stats.slim')" 500
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Replace the global maxconn without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  get_version
  dpa_request PUT "$_GLOBAL_BASE_PATH?version=$VERSION" data/global_maxconn.json
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""

  resource_get "$_RUNTIME_INFO_PATH"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.info.max_conn')" 3500
  assert_equal "$(haproxy_worker_pid)" "$pid"
}
