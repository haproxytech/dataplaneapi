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

@test "groups: Return group G1 from a userlist" {
	resource_get "$_GROUPS_BASE_PATH/G1" "userlist=first"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".name")" "G1"
    assert_equal "$(get_json_path "$BODY" ".users")" "tiger,scott"
}

@test "groups: Return group G2 from a userlist" {
	resource_get "$_GROUPS_BASE_PATH/G2" "userlist=first"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".name")" "G2"
    assert_equal "$(get_json_path "$BODY" ".users")" "scott"
}

@test "groups: Return a non existing group from a userlist" {
	resource_get "$_GROUPS_BASE_PATH/fake" "userlist=first"
	assert_equal "$SC" 404
}

@test "groups: Return a non existing group from a non existing userlist" {
	resource_get "$_GROUPS_BASE_PATH/fake" "userlist=fake"
	assert_equal "$SC" 400
}

@test "groups: Return group one from a userlist" {
	resource_get "$_GROUPS_BASE_PATH/one" "userlist=second"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".name")" "one"
    assert_equal "$(get_json_path "$BODY" ".users")" null
}

@test "groups: Return group two from a userlist" {
	resource_get "$_GROUPS_BASE_PATH/two" "userlist=second"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".name")" "two"
    assert_equal "$(get_json_path "$BODY" ".users")" null
}

@test "groups: Return group three from a userlist" {
	resource_get "$_GROUPS_BASE_PATH/three" "userlist=second"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".name")" "three"
    assert_equal "$(get_json_path "$BODY" ".users")" null
}
