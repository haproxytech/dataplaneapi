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
	PARENT_NAME="mydefaults"
    resource_put "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks/1" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks/1"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "send-state"
}

@test "http_checks: Replace a HTTP Check of backend" {
	PARENT_NAME="test_backend"
	resource_put "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/1" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/1"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "send-state"
}

@test "http_checks: Fail replacing a HTTP Check of unexisting backend" {
	PARENT_NAME="i_am_not_here"
	resource_put "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/1" "data/put.json" "force_reload=true"
	assert_equal "$SC" 400
}

@test "http_checks: Fail replacing an unvalid HTTP Check of backend" {
	PARENT_NAME="test_backend"
	resource_put "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/1" "data/put_unvalid.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 422
}

@test "http_checks: Fail replacing an unexisting HTTP Check of backend" {
	resource_put "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks/1000" "data/put.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 404
}

@test "http_checks: Fail replacing a HTTP Check of frontend" {
	resource_put "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks/1000" "data/put.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 404
}

@test "http_checks: Replace all HTTP Checks of backend (>=2.2)" {
	if haproxy_version_ge "2.2"
    then
	resource_put "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks" "data/replace-all.json" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 202

	resource_get "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 2
	assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace-all.json")" ".")"
	fi
}

@test "http_checks: Replace all HTTP Checks for defaults (>=2.2)" {
	if haproxy_version_ge "2.2"
    then
	resource_put "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks" "data/replace-all.json" "parent_type=defaults&parent_name=mydefaults"
	assert_equal "$SC" 202

	resource_get "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks" "parent_type=defaults&parent_name=mydefaults"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 2
	assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace-all.json")" ".")"
	fi
}
