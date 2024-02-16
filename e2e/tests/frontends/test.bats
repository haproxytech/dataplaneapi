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

@test "frontends: Add a frontend" {
  resource_post "$_FRONTEND_BASE_PATH" "data/post.json" "force_reload=true"
  assert_equal "$SC" "201"
}

@test "frontends: Return a frontend" {
  resource_get "$_FRONTEND_BASE_PATH/test_frontend"
  assert_equal "$SC" 200
  assert_equal "test_frontend" "$(get_json_path "$BODY" '.name')"
}

@test "frontends: Replace a frontend" {
  resource_put "$_FRONTEND_BASE_PATH/test_frontend" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200
  #
  # Retrieving the current data
  #
	resource_get "$_FRONTEND_BASE_PATH/test_frontend"
	assert_equal "http-keep-alive" "$(get_json_path "$BODY" '.http_connection_mode')"
	assert_equal "http" "$(get_json_path "$BODY" '.mode')"
}

@test "frontends: Return an array of frontends" {
  resource_get "$_FRONTEND_BASE_PATH"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" '.[0].name')" "test_frontend"
}

@test "frontends: Delete a frontend" {
  resource_delete "$_FRONTEND_BASE_PATH/test_frontend" "force_reload=true"
	assert_equal "$SC" 204
	#
	# Deleted frontend should be not found
	#
	resource_get "$_FRONTEND_BASE_PATH/test_frontend"
	assert_equal "$SC" 404
}
