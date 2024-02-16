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

@test "users: Return user tiger from a userlist" {
	resource_get "$_USERS_BASE_PATH/tiger" "userlist=first"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".username")" "tiger"
    #assert_equal "$(get_json_path "$BODY" ".secure_password")" `$6$k6y3o.eP$JlKBx9za9667qe4xHSwRv6J.C0/D7cV91`
    assert_equal "$(get_json_path "$BODY" ".secure_password")" true
}

@test "users: Return user scott from a userlist" {
	resource_get "$_USERS_BASE_PATH/scott" "userlist=first"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".username")" "scott"
    assert_equal "$(get_json_path "$BODY" ".secure_password")" false
    assert_equal "$(get_json_path "$BODY" ".password")" "elgato"
}

@test "users: Return a non existing user from a userlist" {
	resource_get "$_USERS_BASE_PATH/fake" "userlist=first"
	assert_equal "$SC" 404
}

@test "users: Return user scott from a non existing userlist" {
	resource_get "$_USERS_BASE_PATH/scott" "userlist=fake"
	assert_equal "$SC" 404
}

@test "users: Return user neo from a userlist" {
	resource_get "$_USERS_BASE_PATH/neo" "userlist=second"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".username")" "neo"
    assert_equal "$(get_json_path "$BODY" ".secure_password")" true
    assert_equal "$(get_json_path "$BODY" ".password")" "JlKBxxHSwRv6J.C0/D7cV91"
    assert_equal "$(get_json_path "$BODY" ".groups")" "one"
}

@test "users: Return user thomas from a userlist" {
	resource_get "$_USERS_BASE_PATH/thomas" "userlist=second"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".username")" "thomas"
    assert_equal "$(get_json_path "$BODY" ".secure_password")" false
    assert_equal "$(get_json_path "$BODY" ".password")" "white-rabbit"
    assert_equal "$(get_json_path "$BODY" ".groups")" "one,two"
}

@test "users: Return user anderson from a userlist" {
	resource_get "$_USERS_BASE_PATH/anderson" "userlist=second"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".username")" "anderson"
    assert_equal "$(get_json_path "$BODY" ".secure_password")" false
    assert_equal "$(get_json_path "$BODY" ".password")" "hello"
    assert_equal "$(get_json_path "$BODY" ".groups")" "two"
}
