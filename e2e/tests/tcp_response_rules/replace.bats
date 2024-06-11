#!/usr/bin/env bats
#
# Copyright 2021 HAProxy Technologies
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

load 'utils/_helpers'

@test "tcp_response_rules: Replace a TCP Response Rule" {
  resource_put "$_TCP_RES_RULES_CERTS_BASE_PATH/0" "data/put.json" "backend=test_backend&force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_TCP_RES_RULES_CERTS_BASE_PATH/1" "backend=test_backend"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".cond_test")" "{ src 192.168.0.0/16 }"
	assert_equal "$(get_json_path "$BODY" ".type")" "content"
	assert_equal "$(get_json_path "$BODY" ".action")" "reject"
}

@test "tcp_response_rules: Replace all TCP Response Rule" {
  resource_put "$_TCP_RES_RULES_CERTS_BASE_PATH" "data/replace-all.json" "backend=test_backend"
	assert_equal "$SC" 202

  resource_get "$_TCP_RES_RULES_CERTS_BASE_PATH" "backend=test_backend"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 2
    assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace-all.json")" ".")"
}

@test "tcp_response_rules: Replace all TCP Response Rule (>=2.7)" {
  if haproxy_version_ge "2.7"
  then
  resource_put "$_TCP_RES_RULES_CERTS_BASE_PATH" "data/replace-all-27.json" "backend=test_backend"
	assert_equal "$SC" 202

  resource_get "$_TCP_RES_RULES_CERTS_BASE_PATH" "backend=test_backend"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 4
    assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace-all-27.json")" ".")"
  fi
}
