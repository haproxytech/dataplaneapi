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

@test "captures: Return an array of all declare captures from the test frontend" {
  resource_get "$_CAPTURES_BASE_PATH" "frontend=test"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "${BODY}" ". | length")" 2
}

@test "captures: Return an array of all declare captures from the test_second frontend" {
  resource_get "$_CAPTURES_BASE_PATH" "frontend=test_second"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "${BODY}" ". | length")" 4
}

@test "captures: Return an array of all declare captures from the test_empty frontend" {
  resource_get "$_CAPTURES_BASE_PATH" "frontend=test_empty"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "${BODY}" ". | length")" 0
}

@test "captures: Return an array of all declare captures from a non existing frontend" {
  resource_get "$_CAPTURES_BASE_PATH" "frontend=fake"
	assert_equal "$SC" 404
}
