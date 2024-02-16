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

@test "users: Return an array of users from userlist first" {
  resource_get "$_USERS_BASE_PATH" "userlist=first"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 2

    assert_equal "$(get_json_path "${BODY}" ".[0].username" )" "tiger"
    assert_equal "$(get_json_path "${BODY}" ".[0].secure_password" )" true
    assert_equal "$(get_json_path "${BODY}" ".[0].groups" )" null

    assert_equal "$(get_json_path "${BODY}" ".[1].username" )" "scott"
    assert_equal "$(get_json_path "${BODY}" ".[1].secure_password" )" false
    assert_equal "$(get_json_path "${BODY}" ".[1].password" )" "elgato"
}

@test "users: Return an array of users from userlist second" {
  resource_get "$_USERS_BASE_PATH" "userlist=second"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 3

    assert_equal "$(get_json_path "${BODY}" ".[0].username" )" "neo"
    assert_equal "$(get_json_path "${BODY}" ".[0].secure_password" )" true
    assert_equal "$(get_json_path "${BODY}" ".[0].password" )" "JlKBxxHSwRv6J.C0/D7cV91"
    assert_equal "$(get_json_path "${BODY}" ".[0].groups" )" "one"

    assert_equal "$(get_json_path "${BODY}" ".[1].username" )" "thomas"
    assert_equal "$(get_json_path "${BODY}" ".[1].secure_password" )" false
    assert_equal "$(get_json_path "${BODY}" ".[1].password" )" "white-rabbit"
    assert_equal "$(get_json_path "${BODY}" ".[1].groups" )" "one,two"

    assert_equal "$(get_json_path "${BODY}" ".[2].username" )" "anderson"
    assert_equal "$(get_json_path "${BODY}" ".[2].secure_password" )" false
    assert_equal "$(get_json_path "${BODY}" ".[2].password" )" "hello"
    assert_equal "$(get_json_path "${BODY}" ".[2].groups" )" "two"
}

@test "users: Return an array of users from userlist empty" {
  resource_get "$_USERS_BASE_PATH" "userlist=empty"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 0
}

@test "users: Return an array of users from non existing userlist" {
  resource_get "$_USERS_BASE_PATH" "userlist=fake"
	assert_equal "$SC" 404
}
