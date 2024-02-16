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

load 'utils/_helpers'

@test "http_checks: Replace a HTTP Check of defaults" {
    resource_put "$_CHECKS_BASE_PATH/1" "data/put.json" "parent_type=defaults&force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_CHECKS_BASE_PATH/1" "parent_type=defaults"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "send-state"
}

@test "http_checks: Replace a HTTP Check of backend" {
	resource_put "$_CHECKS_BASE_PATH/1" "data/put.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_CHECKS_BASE_PATH/1" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "send-state"
}

@test "http_checks: Fail replacing a HTTP Check of unexisting backend" {
	resource_put "$_CHECKS_BASE_PATH/1" "data/put.json" "parent_type=backend&parent_name=i_am_not_here&force_reload=true"
	assert_equal "$SC" 400
}

@test "http_checks: Fail replacing an unvalid HTTP Check of backend" {
	resource_put "$_CHECKS_BASE_PATH/1" "data/put_unvalid.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 422
}

@test "http_checks: Fail replacing an unexisting HTTP Check of backend" {
	resource_put "$_CHECKS_BASE_PATH/1000" "data/put.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 404
}

@test "http_checks: Fail replacing a HTTP Check of frontend" {
	resource_put "$_CHECKS_BASE_PATH/1000" "data/put.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 422
}
