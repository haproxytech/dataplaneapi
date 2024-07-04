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

@test "backend_switching_rules: Replace a Backend Switching Rule" {
  PARENT_NAME="test_frontend"
  resource_put "$_FRONTEND_BASE_PATH/$PARENT_NAME/backend_switching_rules/0" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/put.json")" ".")"
}

@test "backend_switching_rules: Replace all Backend Switching Rule" {
  PARENT_NAME="test_frontend"
  resource_put "$_FRONTEND_BASE_PATH/$PARENT_NAME/backend_switching_rules" "data/replace-all.json" "force_reload=true"
	assert_equal "$SC" 200

  resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/backend_switching_rules"
    assert_equal "$(get_json_path "${BODY}" ". | length")" 3
	assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace-all.json")" ".")"
}
