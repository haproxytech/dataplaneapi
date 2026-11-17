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

_SRV2_WEIGHT_25='s/^\( *server srv2 10\.0\.1\.2:8080\) weight 20$/\1 weight 25/'

@test "runtime_apply: Post a raw configuration changing a server weight without a reload" {
  local pid; pid=$(haproxy_worker_pid)
  local cfg; cfg=$(raw_config_with "$_SRV2_WEIGHT_25")
  assert grep -q 'server srv2 10.0.1.2:8080 weight 25' "$cfg"

  get_version
  dpa_request POST "$_RAW_BASE_PATH?version=$VERSION" "$cfg" text/plain
  assert_equal "$SC" 201
  assert_equal "$RELOAD_ID" ""

  assert_equal "$(runtime_server_field be_plain srv2 weight)" 25
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Post the unchanged raw configuration without a reload" {
  local pid; pid=$(haproxy_worker_pid)
  local cfg; cfg=$(raw_config_with "")

  get_version
  dpa_request POST "$_RAW_BASE_PATH?version=$VERSION" "$cfg" text/plain
  assert_equal "$SC" 201
  assert_equal "$RELOAD_ID" ""

  assert_equal "$(runtime_server_field be_plain srv2 weight)" 20
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Post a raw configuration adding an ACL reloads" {
  local pid; pid=$(haproxy_worker_pid)
  local cfg; cfg=$(raw_config_with '/^backend be_plain/a\  acl is_api path_beg /api')
  assert grep -q 'acl is_api path_beg /api' "$cfg"

  get_version
  dpa_request POST "$_RAW_BASE_PATH?version=$VERSION" "$cfg" text/plain
  assert_equal "$SC" 202
  assert [ -n "$RELOAD_ID" ]

  wait_reload_succeeded "$RELOAD_ID"
  wait_new_worker "$pid"
}

@test "runtime_apply: Post a raw configuration with skip_reload leaves HAProxy untouched" {
  local pid; pid=$(haproxy_worker_pid)
  local cfg; cfg=$(raw_config_with "$_SRV2_WEIGHT_25")

  get_version
  dpa_request POST "$_RAW_BASE_PATH?skip_reload=true&version=$VERSION" "$cfg" text/plain
  assert_equal "$SC" 201
  assert_equal "$RELOAD_ID" ""

  # the file has the new weight, the running HAProxy the old one
  resource_get "$_BACKEND_BASE_PATH/be_plain/servers/srv2"
  assert_equal "$(get_json_path "$BODY" '.weight')" 25
  assert_equal "$(runtime_server_field be_plain srv2 weight)" 20
  assert_equal "$(haproxy_worker_pid)" "$pid"
}
