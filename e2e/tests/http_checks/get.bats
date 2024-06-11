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

@test "http_checks: Return one HTTP Check from defaults" {
  	resource_get "$_CHECKS_BASE_PATH/0" "parent_type=defaults&parent_name=mydefaults"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "send-state"

	resource_get "$_CHECKS_BASE_PATH/1" "parent_type=defaults&parent_name=mydefaults"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "disable-on-404"
}

@test "http_checks: Return one HTTP Check from backend" {
	resource_get "$_CHECKS_BASE_PATH/0" "parent_type=backend&parent_name=test_backend"
    assert_equal "$(get_json_path "$BODY" ".type")" "send"
    assert_equal 1 "$(get_json_path "$BODY" ".headers | length")"
    assert_equal "$(get_json_path "$BODY" ".headers[0].name")" "host"
    assert_equal "$(get_json_path "$BODY" ".headers[0].fmt")" "haproxy.1wt.eu"

	resource_get "$_CHECKS_BASE_PATH/1" "parent_type=backend&parent_name=test_backend"
    assert_equal "$(get_json_path "$BODY" ".type")" "expect"
    assert_equal "$(get_json_path "$BODY" ".match")" "status"
    assert_equal "$(get_json_path "$BODY" ".pattern")" "200-399"
}

@test "http_checks: Return 422 when fetching HTTP Check from frontend" {
	resource_get "$_CHECKS_BASE_PATH/0" "parent_type=frontend&parent_name=test_frontend"
	assert_equal "$SC" 422
}

@test "http_checks: Return 400 when fetching HTTP Check from unexisting backend" {
    resource_get "$_CHECKS_BASE_PATH/0" "parent_type=backend&parent_name=i_am_not_here"
	assert_equal "$SC" 400
}

@test "http_checks: Return 404 when fetching unexisting HTTP Check" {
    resource_get "$_CHECKS_BASE_PATH/1000" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 404
}
