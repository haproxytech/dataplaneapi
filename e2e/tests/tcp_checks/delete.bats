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

@test "tcp_checks: Delete some TCP checks" {
  PARENT_NAME="test_backend_delete"
  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "force_reload=true"
	assert_equal "$SC" 204

  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/1" "force_reload=true"
	assert_equal "$SC" 204

  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/1000" "force_reload=true"
	assert_equal "$SC" 404

  resource_delete "$_FRONTEND_BASE_PATH/$PARENT_NAME/tcp_checks/1000" "force_reload=true"
	assert_equal "$SC" 404

  resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "${BODY}" ". | length")" 2
}

@test "tcp_checks: Delete a non existant backend TCP check" {
  PARENT_NAME="test_backend_delete"
  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/1000" "force_reload=true"
	assert_equal "$SC" 404
}

@test "tcp_checks: Delete a non existant frontend TCP check" {
  PARENT_NAME="test_frontend"
  resource_delete "$_FRONTEND_BASE_PATH/$PARENT_NAME/tcp_checks/1000" "force_reload=true"
	assert_equal "$SC" 404
}
