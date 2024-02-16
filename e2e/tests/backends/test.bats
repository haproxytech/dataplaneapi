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
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "backends: Add a backend" {
  resource_post "$_BACKEND_BASE_PATH" "data/post.json" "force_reload=true"
  assert_equal "$SC" "201"

  resource_get "$_BACKEND_BASE_PATH/test_backend"  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".adv_check")" "httpchk"
  assert_equal "$(get_json_path "$BODY" ".httpchk_params.method")" "GET"
  assert_equal "$(get_json_path "$BODY" ".httpchk_params.uri")" "/check"
  assert_equal "$(get_json_path "$BODY" ".httpchk_params.version")" "HTTP/1.1"
}

@test "backends: fail adding a backend (invalid send method in httpchk_params)" {
  resource_post "$_BACKEND_BASE_PATH" "data/post_invalid_send_method_1.json" "force_reload=true"
	assert_equal "$SC" 422
    assert_equal "$(get_json_path "$BODY" ".code")" "606"
}

@test "backends: fail adding a backend (invalid send method in http-check)" {
  resource_post "$_BACKEND_BASE_PATH" "data/post_invalid_send_method_2.json" "force_reload=true"
	assert_equal "$SC" 422
    assert_equal "$(get_json_path "$BODY" ".code")" "606"
}

@test "backends: Return a backend" {
	resource_get "$_BACKEND_BASE_PATH/test_backend"
	assert_equal "$SC" 200
	assert_equal "test_backend" "$(get_json_path "$BODY" '.name')"
}

@test "backends: Replace a backend" {
	resource_put "$_BACKEND_BASE_PATH/test_backend" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_BACKEND_BASE_PATH/test_backend"  assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ".adv_check")" "httpchk"
    assert_equal "$(get_json_path "$BODY" ".httpchk_params.method")" "GET"
    assert_equal "$(get_json_path "$BODY" ".httpchk_params.uri")" "/healthz"
    assert_equal "$(get_json_path "$BODY" ".httpchk_params.version")" "HTTP/1.1"
}

@test "backends: Return an array of backends" {
	resource_get "$_BACKEND_BASE_PATH"
	assert_equal "$SC" 200
	assert_equal "test_backend" "$(get_json_path "$BODY" '.[0].name')"
}

@test "backends: Delete a backend" {
	resource_delete "$_BACKEND_BASE_PATH/test_backend" "force_reload=true"
	assert_equal "$SC" 204
}
