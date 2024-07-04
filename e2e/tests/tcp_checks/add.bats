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
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_version'
load '../../libs/get_json_path'

load 'utils/_helpers'

@test "tcp_checks: Add a new connect TCP check to a backend" {
  PARENT_NAME="test_backend_add"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/connect.json" "force_reload=true"
	assert_equal "$SC" 201
}

@test "tcp_checks: Add a new send TCP check to a backend" {
  PARENT_NAME="test_backend_add"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/send.json" "force_reload=true"
	assert_equal "$SC" 201
}

@test "tcp_checks: Add a new expect TCP check to a backend" {
  PARENT_NAME="test_backend_add"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/expect.json" "force_reload=true"
	assert_equal "$SC" 201
}

@test "tcp_checks: Add a new send-binary TCP check to a backend" {
  PARENT_NAME="test_backend_add"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/send_binary.json" "force_reload=true"
	assert_equal "$SC" 201
}

@test "tcp_checks: Add a new comment TCP check to a backend" {
    if haproxy_version_ge "2.2"
    then
    PARENT_NAME="test_backend_add"
    resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/comment.json" "force_reload=true"
    assert_equal "$SC" 201
  fi
}

@test "tcp_checks: Add a new send-lf TCP check to a backend" {
  PARENT_NAME="test_backend_add"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/send_lf.json" "force_reload=true"
  if haproxy_version_ge "2.2"
  then
    assert_equal "$SC" 201
  else
    assert_equal "$SC" 400
  fi
}

@test "tcp_checks: Add a new send-binary-lf TCP check to a backend" {
  PARENT_NAME="test_backend_add"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/send_binary_lf.json" "force_reload=true"
  if haproxy_version_ge "2.2"
  then
    assert_equal "$SC" 201
  else
    assert_equal "$SC" 400
  fi
}

@test "tcp_checks: Add a new set-var and uset-var TCP check to a backend" {
  PARENT_NAME="test_backend_add"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/set_var.json" "force_reload=true"
  if haproxy_version_ge "2.2"
  then
    assert_equal "$SC" 201
  else
    assert_equal "$SC" 400
  fi
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/unset_var.json" "force_reload=true"
  if haproxy_version_ge "2.2"
  then
    assert_equal "$SC" 201
  else
    assert_equal "$SC" 400
  fi
}

@test "tcp_checks: Add an empty TCP check to a backend/0" {
  PARENT_NAME="test_backend_add"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/empty.json" "force_reload=true"
  if haproxy_version_ge "2.2"
  then
    assert_equal "$SC" 422
  else
    assert_equal "$SC" 422
  fi
}
