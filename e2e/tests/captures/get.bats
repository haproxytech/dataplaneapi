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

@test "captures: Get one declare capture request" {
    PARENT_NAME="test"
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/captures/0" "frontend=test"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".type")" "request"
    assert_equal "$(get_json_path "${BODY}" ".length")" 1
}

@test "captures: Get one declare capture response" {
    PARENT_NAME="test"
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/captures/1"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".type")" "response"
    assert_equal "$(get_json_path "${BODY}" ".length")" 2
}

@test "captures: Get one non existing declare capture" {
    PARENT_NAME="test"
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/captures/1000" "frontend=test"
    assert_equal "$SC" 404
}

@test "captures: Get one non existing declare capture from a non existant frontend" {
    PARENT_NAME="fake"
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/captures/1000"
    assert_equal "$SC" 404
}

@test "captures: Get first declare capture request" {
    PARENT_NAME="test_second"
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/captures/0"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".type")" "request"
    assert_equal "$(get_json_path "${BODY}" ".length")" 111
}

@test "captures: Get second declare capture request" {
    PARENT_NAME="test_second"
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/captures/1"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".type")" "response"
    assert_equal "$(get_json_path "${BODY}" ".length")" 222
}

@test "captures: Get third declare capture request" {
    PARENT_NAME="test_second"
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/captures/2"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".type")" "request"
    assert_equal "$(get_json_path "${BODY}" ".length")" 8888
}

@test "captures: Get fourth declare capture request" {
    PARENT_NAME="test_second"
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/captures/3"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".type")" "response"
    assert_equal "$(get_json_path "${BODY}" ".length")" 9999
}

@test "captures: Get fifth declare capture request" {
    resource_get "$_CAPTURES_BASE_PATH/4" "frontend=test_second"
    assert_equal "$SC" 404
}
