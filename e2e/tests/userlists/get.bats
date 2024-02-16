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

@test "userlists: Return userlist first" {
    resource_get "$_USERLISTS_BASE_PATH/first"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".name")" "first"
}

@test "userlists: Return userlist second" {
    resource_get "$_USERLISTS_BASE_PATH/second"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".name")" "second"
}

@test "userlists: Return userlist empty" {
    resource_get "$_USERLISTS_BASE_PATH/empty"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".name")" "empty"
}

@test "userlists: Return userlist add_test" {
    resource_get "$_USERLISTS_BASE_PATH/add_test"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".name")" "add_test"
}

@test "userlists: Return userlist replace_test" {
    resource_get "$_USERLISTS_BASE_PATH/replace_test"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".name")" "replace_test"
}

@test "userlists: Return userlist delete_test" {
    resource_get "$_USERLISTS_BASE_PATH/delete_test"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".name")" "delete_test"
}

@test "userlists: Return userlist 0" {
    resource_get "$_USERLISTS_BASE_PATH/0"
    assert_equal "$SC" 404
}

@test "userlists: Return userlist 1000" {
    resource_get "$_USERLISTS_BASE_PATH/1000"
    assert_equal "$SC" 404
}

@test "userlists: Return userlist fake" {
    resource_get "$_USERLISTS_BASE_PATH/fake"
    assert_equal "$SC" 404
}
