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
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "filter_configs: Return filter configs from all parent types" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  for PARENT in "backends/test_backend" "frontends/test_frontend" "defaults/filter_defaults"; do
    resource_get "/services/haproxy/configuration/$PARENT/filter_configs"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ". | length")" 1
    assert_equal "$(get_json_path "$BODY" ".[0].class")" "trace"
    assert_equal "$(get_json_path "$BODY" ".[0].id")" "A"
  done
}

@test "filter_configs: Return one filter config from all parent types" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  for PARENT in "backends/test_backend" "frontends/test_frontend" "defaults/filter_defaults"; do
    resource_get "/services/haproxy/configuration/$PARENT/filter_configs/0"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ".class")" "trace"
    assert_equal "$(get_json_path "$BODY" ".id")" "A"

    resource_get "/services/haproxy/configuration/$PARENT/filter_configs/1"
    assert_equal "$SC" 404
  done
}

@test "filter_configs: Add a filter config" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_post "$_BACKEND_BASE_PATH/test_backend/filter_configs/1" "data/post.json" "force_reload=true"
  assert_equal "$SC" 201
  assert_equal "$(get_json_path "$BODY" ".class")" "trace"
  assert_equal "$(get_json_path "$BODY" ".id")" "B"

  resource_get "$_BACKEND_BASE_PATH/test_backend/filter_configs/1"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".class")" "trace"
  assert_equal "$(get_json_path "$BODY" ".id")" "B"

  resource_get "$_BACKEND_BASE_PATH/test_backend/filter_configs"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ". | length")" 2
}

@test "filter_configs: Replace a filter config" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_put "$_BACKEND_BASE_PATH/test_backend/filter_configs/0" "data/put.json" "force_reload=true"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".class")" "trace"
  assert_equal "$(get_json_path "$BODY" ".id")" "C"

  resource_get "$_BACKEND_BASE_PATH/test_backend/filter_configs/0"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".id")" "C"
}

@test "filter_configs: Replace all filter configs" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_put "$_BACKEND_BASE_PATH/test_backend/filter_configs" "data/replace_all.json" "force_reload=true"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ". | length")" 2

  resource_get "$_BACKEND_BASE_PATH/test_backend/filter_configs"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ". | length")" 2
  assert_equal "$(get_json_path "$BODY" ".[0].id")" "B"
  assert_equal "$(get_json_path "$BODY" ".[1].id")" "C"
}

@test "filter_configs: Delete a filter config" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_delete "$_BACKEND_BASE_PATH/test_backend/filter_configs/0" "force_reload=true"
  assert_equal "$SC" 204

  resource_get "$_BACKEND_BASE_PATH/test_backend/filter_configs/0"
  assert_equal "$SC" 404

  resource_delete "$_BACKEND_BASE_PATH/test_backend/filter_configs/0" "force_reload=true"
  assert_equal "$SC" 404
}

@test "filter_configs: Handle a missing parent" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_get "$_BACKEND_BASE_PATH/ghost/filter_configs"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ". | length")" 0

  resource_post "$_BACKEND_BASE_PATH/ghost/filter_configs/0" "data/post.json" "force_reload=true"
  assert_equal "$SC" 400
}
