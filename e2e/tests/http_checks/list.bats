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

@test "http_checks: Return an array of all HTTP Checks from defaults" {
    resource_get "$_CHECKS_BASE_PATH" "parent_type=defaults&parent_name=mydefaults"
	assert_equal "$SC" 200
    assert_equal 2 "$(get_json_path "$BODY" ". | length")"

    assert_equal "$(get_json_path "$BODY" ".[0].type")" "send-state"

    assert_equal "$(get_json_path "$BODY" ".[1].type")" "disable-on-404"
}

@test "http_checks: Return an array of HTTP Checks from backend" {
    resource_get "$_CHECKS_BASE_PATH" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
	assert_equal 2 "$(get_json_path "$BODY" ". | length")"

    assert_equal "$(get_json_path "$BODY" ".[0].type")" "send"
    assert_equal 1 "$(get_json_path "$BODY" ".[0].headers | length")"
    assert_equal "$(get_json_path "$BODY" ".[0].headers[0].name")" "host"
    assert_equal "$(get_json_path "$BODY" ".[0].headers[0].fmt")" "haproxy.1wt.eu"

    assert_equal "$(get_json_path "$BODY" ".[1].type")" "expect"
    assert_equal "$(get_json_path "$BODY" ".[1].match")" "status"
    assert_equal "$(get_json_path "$BODY" ".[1].pattern")" "200-399"
}

@test "http_checks: Return 422 for frontend" {
    resource_get "$_CHECKS_BASE_PATH" "parent_type=frontend&parent_name=test_frontend"
	assert_equal "$SC" 422
}
