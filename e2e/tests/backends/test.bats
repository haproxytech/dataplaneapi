#!/usr/bin/env bats
#
# Copyright 2019 HAProxy Technologies
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
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "backends: Add a backend" {
  resource_post "$_BACKEND_BASE_PATH" "data/post.json" "force_reload=true"
  assert_equal "$SC" "201"
}

@test "backends: Return a backend" {
	resource_get "$_BACKEND_BASE_PATH/test_backend"
	assert_equal "$SC" 200
	assert_equal "test_backend" "$(get_json_path "$BODY" '.data.name')"
}

@test "backends: Replace a backend" {
	resource_put "$_BACKEND_BASE_PATH/test_backend" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200
}

@test "backends: Return an array of backends" {
	resource_get "$_BACKEND_BASE_PATH"
	assert_equal "$SC" 200
	assert_equal "test_backend" "$(get_json_path "$BODY" '.data[0].name')"
}

@test "backends: Delete a backend" {
	resource_delete "$_BACKEND_BASE_PATH/test_backend" "force_reload=true"
	assert_equal "$SC" 204
}
