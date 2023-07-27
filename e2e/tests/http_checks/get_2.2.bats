#!/usr/bin/env bats
#
# Copyright 2023 HAProxy Technologies
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
load '../../libs/haproxy_version'

load 'utils/_helpers'


@test "http_checks: Return one HTTP Check from backend (>=2.2)" {
	if haproxy_version_ge "2.2"
    then
	resource_get "$_CHECKS_BASE_PATH/0" "parent_type=backend&parent_name=test_backend"
    assert_equal "$(get_json_path "$BODY" ".data.type")" "send"
    assert_equal 1 "$(get_json_path "$BODY" ".data.headers | length")"
    assert_equal "$(get_json_path "$BODY" ".data.headers[0].name")" "host"
    assert_equal "$(get_json_path "$BODY" ".data.headers[0].fmt")" "haproxy.1wt.eu"
    assert_equal "$(get_json_path "$BODY" ".data.method")" "OPTIONS"
	assert_equal "$(get_json_path "$BODY" ".data.uri")" "/"
	assert_equal "$(get_json_path "$BODY" ".data.version")" "HTTP/1.1"

	resource_get "$_CHECKS_BASE_PATH/1" "parent_type=backend&parent_name=test_backend"
    assert_equal "$(get_json_path "$BODY" ".data.type")" "expect"
    assert_equal "$(get_json_path "$BODY" ".data.match")" "status"
    assert_equal "$(get_json_path "$BODY" ".data.pattern")" "200-399"
	fi
}
