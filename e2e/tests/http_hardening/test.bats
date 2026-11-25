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

@test "http_hardening: Return all HTTP hardening profiles" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_get "$_HTTP_HARDENING_BASE_PATH"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ". | length")" 2
  assert_equal "$(get_json_path "$BODY" ".[0].name")" "base"
  assert_equal "$(get_json_path "$BODY" ".[0].all")" "on"
  assert_equal "$(get_json_path "$BODY" ".[1].name")" "child"
  assert_equal "$(get_json_path "$BODY" ".[1].from")" "base"
  assert_equal "$(get_json_path "$BODY" ".[1].all_delete")" "on"
  assert_equal "$(get_json_path "$BODY" ".[1].reject_lowercase_method")" "off"
}

@test "http_hardening: Return one HTTP hardening profile" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_get "$_HTTP_HARDENING_BASE_PATH/child"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".name")" "child"
  assert_equal "$(get_json_path "$BODY" ".from")" "base"
  assert_equal "$(get_json_path "$BODY" ".all_delete")" "on"
  assert_equal "$(get_json_path "$BODY" ".reject_lowercase_method")" "off"

  resource_get "$_HTTP_HARDENING_BASE_PATH/missing"
  assert_equal "$SC" 404
}

@test "http_hardening: Create an HTTP hardening profile" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_post "$_HTTP_HARDENING_BASE_PATH" "data/post.json" "force_reload=true"
  assert_equal "$SC" 201
  assert_equal "$(get_json_path "$BODY" ".name")" "created"
  assert_equal "$(get_json_path "$BODY" ".all")" "on"

  resource_get "$_HTTP_HARDENING_BASE_PATH/created"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".name")" "created"
  assert_equal "$(get_json_path "$BODY" ".all")" "on"
  assert_equal "$(get_json_path "$BODY" ".reject_duplicate_host_header")" "off"

  resource_post "$_HTTP_HARDENING_BASE_PATH" "data/post.json" "force_reload=true"
  assert_equal "$SC" 409

  resource_post "$_HTTP_HARDENING_BASE_PATH" "data/post_invalid.json" "force_reload=true"
  assert_equal "$SC" 422
}

@test "http_hardening: Replace an HTTP hardening profile" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_put "$_HTTP_HARDENING_BASE_PATH/base" "data/put.json" "force_reload=true"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".name")" "base"
  assert_equal "$(get_json_path "$BODY" ".all")" "off"

  resource_get "$_HTTP_HARDENING_BASE_PATH/base"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".all")" "off"
  assert_equal "$(get_json_path "$BODY" ".reject_underscore_in_header_name")" "on"

  resource_put "$_HTTP_HARDENING_BASE_PATH/base" "data/put_wrong_name.json" "force_reload=true"
  assert_equal "$SC" 400

  resource_put "$_HTTP_HARDENING_BASE_PATH/missing" "data/put_missing.json" "force_reload=true"
  assert_equal "$SC" 404
}

@test "http_hardening: Delete an HTTP hardening profile" {
  haproxy_version_ge "3.5" || skip "requires HAProxy 3.5+"

  resource_delete "$_HTTP_HARDENING_BASE_PATH/child" "force_reload=true"
  assert_equal "$SC" 204

  resource_get "$_HTTP_HARDENING_BASE_PATH/child"
  assert_equal "$SC" 404

  resource_delete "$_HTTP_HARDENING_BASE_PATH/child" "force_reload=true"
  assert_equal "$SC" 404
}
