#!/usr/bin/env bats
#
# Copyright 2022 HAProxy Technologies
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
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/debug'
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "defaults: Create a named defaults configuration" {
  resource_put "$_DEFAULTS_BASE_PATH" "data/post.json" ""
  assert_equal "$SC" 202

  resource_get "$_DEFAULTS_BASE_PATH/created"
  assert_equal "$SC" 200

  assert_equal "$(get_json_path "$BODY" '.name')" "created"
  assert_equal "$(get_json_path "$BODY" '.server_timeout')" "20000"
  assert_equal "$(get_json_path "$BODY" '.client_timeout')" "20000"
  assert_equal "$(get_json_path "$BODY" '.mode')" "http"
}

@test "defaults: Create a named defaults configuration that already exists" {
  resource_put "$_DEFAULTS_BASE_PATH" "data/post_existing.json" ""
  assert_equal "$SC" 409
}

@test "defaults: Create a named defaults configuration with from" {
  haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

  resource_put "$_DEFAULTS_BASE_PATH" "data/post.json" ""
  assert_equal "$SC" 202

  resource_get "$_DEFAULTS_BASE_PATH/created"
  assert_equal "$SC" 200

  resource_put "$_DEFAULTS_BASE_PATH" "data/post_with_from.json" ""
  assert_equal "$SC" 202

  resource_get "$_DEFAULTS_BASE_PATH/created_with_from"
  assert_equal "$SC" 200

  assert_equal "$(get_json_path "$BODY" '.name')" "created_with_from"
  assert_equal "$(get_json_path "$BODY" '.from')" "created"
  assert_equal "$(get_json_path "$BODY" '.server_timeout')" "20000"
  assert_equal "$(get_json_path "$BODY" '.client_timeout')" "20000"
  assert_equal "$(get_json_path "$BODY" '.mode')" "http"
}
