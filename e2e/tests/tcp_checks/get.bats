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

@test "tcp_checks: Return connect TCP check from a backend" {
	PARENT_NAME="test_backend_get"
	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".action")" "connect"
}

@test "tcp_checks: Return send TCP check from a backend" {
	PARENT_NAME="test_backend_get"
	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/1"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".action")" "send"
}

@test "tcp_checks: Return expect TCP check from a backend" {
	PARENT_NAME="test_backend_get"
	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/2"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".action")" "expect"
}

@test "tcp_checks: Return send-binary TCP check from a backend" {
	PARENT_NAME="test_backend_get"
	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/3"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".action")" "send-binary"
	assert_equal "$(get_json_path "$BODY" ".hex_string")" "50494e470d0a"
}

@test "tcp_checks: Return a non existing TCP check from a backend" {
	PARENT_NAME="test_backend_get"
  resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/1000"
  assert_equal "$SC" 404
}

@test "tcp_checks: Return a TCP check from a non existing backend" {
	PARENT_NAME="unknown"
  resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0"
  assert_equal "$SC" 400
}

@test "tcp_checks: Return a TCP check from a frontend" {
	PARENT_NAME="unknown"
  resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/tcp_checks/0"
  assert_equal "$SC" 404
}
