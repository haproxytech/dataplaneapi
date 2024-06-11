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
load "../../libs/get_json_path"
load '../../libs/version'
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'

load 'utils/_helpers'

@test "captures: Replace an existing declare capture" {
    resource_put "$_CAPTURES_BASE_PATH/0" "data/replace.json" "frontend=test_replace"
    assert_equal "$SC" 202
}

@test "captures: Replace an existing declare capture with an empty one" {
    resource_put "$_CAPTURES_BASE_PATH/0" "data/empty.json" "frontend=test_replace"
    assert_equal "$SC" 422
}

@test "captures: Replace a non existing declare capture in a non existing frontend" {
    resource_put "$_CAPTURES_BASE_PATH/0" "data/replace.json" "frontend=fake"
    assert_equal "$SC" 404
}

@test "captures: Replace all existing declare capture" {
    resource_put "$_CAPTURES_BASE_PATH" "data/replace-all.json" "frontend=test_replace"
    assert_equal "$SC" 202

    resource_get "$_CAPTURES_BASE_PATH" "frontend=test_replace"
    assert_equal "$(get_json_path "${BODY}" ". | length")" 3
	assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace-all.json")" ".")"
}
