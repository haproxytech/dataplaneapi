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

@test "log_targets: Return an array of all Log Targets from frontend" {
	PARENT_NAME="test_frontend"
  resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/log_targets"
	assert_equal "$SC" 200

	assert_equal "$(get_json_path "${BODY}" ". | length")" 2

	assert_equal "$(get_json_path "${BODY}" ".[0].address")" "localhost"
	assert_equal "$(get_json_path "${BODY}" ".[0].facility")" "user"
	assert_equal "$(get_json_path "${BODY}" ".[0].format")" "raw"
	assert_equal "$(get_json_path "${BODY}" ".[0].level")" "warning"

	assert_equal "$(get_json_path "${BODY}" ".[1].address")" "10.0.0.1"
	assert_equal "$(get_json_path "${BODY}" ".[1].facility")" "user"
	assert_equal "$(get_json_path "${BODY}" ".[1].format")" "raw"
	assert_equal "$(get_json_path "${BODY}" ".[1].level")" "info"
}

@test "log_targets: Return an array of all Log Targets from backend" {
	PARENT_NAME="test_backend"
  resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/log_targets" "parent_type=backend&parent_name=test_backend"
  assert_equal "$SC" 200

  assert_equal "$(get_json_path "${BODY}" ". | length")" 2

	assert_equal "$(get_json_path "${BODY}" ".[0].address")" "localhost"
	assert_equal "$(get_json_path "${BODY}" ".[0].facility")" "user"
	assert_equal "$(get_json_path "${BODY}" ".[0].format")" "raw"
	assert_equal "$(get_json_path "${BODY}" ".[0].level")" "warning"

	assert_equal "$(get_json_path "${BODY}" ".[1].address")" "10.0.0.1"
	assert_equal "$(get_json_path "${BODY}" ".[1].facility")" "user"
	assert_equal "$(get_json_path "${BODY}" ".[1].format")" "raw"
	assert_equal "$(get_json_path "${BODY}" ".[1].level")" "info"
}
