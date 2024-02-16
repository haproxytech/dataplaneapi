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

@test "groups: Return an array of groups from userlist first" {
  resource_get "$_GROUPS_BASE_PATH" "userlist=first"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 2

    assert_equal "$(get_json_path "${BODY}" ".[0].name" )" "G1"
    assert_equal "$(get_json_path "${BODY}" ".[0].users" )" "tiger,scott"

    assert_equal "$(get_json_path "${BODY}" ".[1].name" )" "G2"
    assert_equal "$(get_json_path "${BODY}" ".[1].users" )" "scott"
}

@test "groups: Return an array of groups from userlist second" {
  resource_get "$_GROUPS_BASE_PATH" "userlist=second"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 3

    assert_equal "$(get_json_path "${BODY}" ".[0].name" )" "one"
    assert_equal "$(get_json_path "${BODY}" ".[0].users" )" null

    assert_equal "$(get_json_path "${BODY}" ".[1].name" )" "two"
    assert_equal "$(get_json_path "${BODY}" ".[1].users" )" null

    assert_equal "$(get_json_path "${BODY}" ".[2].name" )" "three"
    assert_equal "$(get_json_path "${BODY}" ".[2].users" )" null
}

@test "groups: Return an array of groups from userlist empty" {
  resource_get "$_GROUPS_BASE_PATH" "userlist=empty"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 0
}

@test "groups: Return an array of groups from non existing userlist" {
  resource_get "$_GROUPS_BASE_PATH" "userlist=fake"
	assert_equal "$SC" 404
}
