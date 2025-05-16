#!/usr/bin/env bats
#
# Copyright 2025 HAProxy Technologies
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
load '../../libs/debug'
load '../../libs/get_json_path'
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

_ACME_PATH="/services/haproxy/configuration/acme"
_ACME_NAME="letsencrypt"

@test "acme: all tests (>=3.2)" {
  haproxy_version_ge "3.2" || skip

  resource_post "$_ACME_PATH" "data/new_acme.json" "force_reload=true"
  assert_equal "$SC" "201"

  resource_get "$_ACME_PATH/$_ACME_NAME"
  assert_equal "$SC" "200"
  assert_equal "$(get_json_path "$BODY" .name)" "$_ACME_NAME"
  assert_equal "$(get_json_path "$BODY" .contact)" "me@example.com"
  assert_equal "$(get_json_path "$BODY" .directory)" "https://acme-staging-v02.api.letsencrypt.org/directory"

  resource_put "$_ACME_PATH/$_ACME_NAME" "data/edit_acme.json" "force_reload=true"
  assert_equal "$SC" "200"
  resource_get "$_ACME_PATH/$_ACME_NAME"
  assert_equal "$(get_json_path "$BODY" .directory)" "https://acme-v02.api.letsencrypt.org/directory"
  assert_equal "$(get_json_path "$BODY" .keytype)" "RSA"
  assert_equal "$(get_json_path "$BODY" .bits)" 4096

  resource_get "$_ACME_PATH"
  assert_equal "$SC" "200"
  assert_equal "$(get_json_path "$BODY" '.|length')" 1
  assert_equal "$(get_json_path "$BODY" .[0].name)" "$_ACME_NAME"

  # back to the original
  resource_put "$_ACME_PATH/$_ACME_NAME" "data/new_acme.json" "force_reload=true"
  assert_equal "$SC" "200"
  resource_get "$_ACME_PATH/$_ACME_NAME"
  assert_equal "$(get_json_path "$BODY" .directory)" "https://acme-staging-v02.api.letsencrypt.org/directory"
  assert_equal "$(get_json_path "$BODY" .keytype)" null
  assert_equal "$(get_json_path "$BODY" .bits)" null

  resource_delete "$_ACME_PATH/$_ACME_NAME" "force_reload=true"
  assert_equal "$SC" "204"
  resource_get "$_ACME_PATH/$_ACME_NAME"
  assert_equal "$SC" "404"
}
