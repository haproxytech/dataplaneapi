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

@test "log_targets: Return one Log Target from frontend" {
	PARENT_NAME="test_frontend"
  resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/log_targets/0"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".address")" "localhost"
	assert_equal "$(get_json_path "$BODY" ".facility")" "user"
	assert_equal "$(get_json_path "$BODY" ".format")" "raw"
	assert_equal "$(get_json_path "$BODY" ".level")" "warning"

	resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/log_targets/1"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "${BODY}" ".address")" "10.0.0.1"
	assert_equal "$(get_json_path "${BODY}" ".facility")" "user"
	assert_equal "$(get_json_path "${BODY}" ".format")" "raw"
	assert_equal "$(get_json_path "${BODY}" ".level")" "info"
}

@test "log_targets: Return one Log Target from backend" {
	PARENT_NAME="test_backend"
  resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/log_targets/0" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".address")" "localhost"
	assert_equal "$(get_json_path "$BODY" ".facility")" "user"
	assert_equal "$(get_json_path "$BODY" ".format")" "raw"
	assert_equal "$(get_json_path "$BODY" ".level")" "warning"

	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/log_targets/1"
	assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ".address")" "10.0.0.1"
	assert_equal "$(get_json_path "$BODY" ".facility")" "user"
	assert_equal "$(get_json_path "$BODY" ".format")" "raw"
	assert_equal "$(get_json_path "$BODY" ".level")" "info"
}
