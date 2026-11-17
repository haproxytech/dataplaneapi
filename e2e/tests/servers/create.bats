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
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/haproxy_version'
load '../../libs/haproxy_worker'
load '../../libs/version'

load 'utils/_helpers'

@test "servers: Add a new server to backend" {
  PARENT_NAME="test_backend"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/servers" "data/post.json" "force_reload=true"
	assert_equal "$SC" 201
}

@test "servers: Add a new server to backend through runtime with deprecated backend param" {
  haproxy_version_ge 2.6 || skip "requires HAProxy 2.6+"

  local pid; pid=$(haproxy_worker_pid)

  PARENT_NAME="test_backend"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/servers" "data/post.json"
  assert_equal "$SC" 201

  # the running HAProxy got the server without a reload
  resource_get "$_RUNTIME_BACKEND_BASE_PATH/$PARENT_NAME/servers/test_server"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.address')" "10.1.0.1"
  assert_equal "$(get_json_path "$BODY" '.port')" 8080
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "servers: Add a new server to backend through runtime with parent_type/ parent_name" {
  haproxy_version_ge 2.6 || skip "requires HAProxy 2.6+"
  local pid; pid=$(haproxy_worker_pid)

  PARENT_NAME="test_backend"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/servers" "data/post.json"
  assert_equal "$SC" 201

  # the running HAProxy got the server without a reload
  resource_get "$_RUNTIME_BACKEND_BASE_PATH/$PARENT_NAME/servers/test_server"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.address')" "10.1.0.1"
  assert_equal "$(get_json_path "$BODY" '.port')" 8080
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "servers: Add a new server to peer" {
  PARENT_NAME="fusion"
  resource_post "$_PEER_BASE_PATH/$PARENT_NAME/servers" "data/post.json"
	assert_equal "$SC" 202
}

@test "servers: Add a new server to ring" {
  haproxy_version_ge 2.2 || skip "requires HAProxy 2.2+"
  PARENT_NAME=logbuffer
  resource_post "$_RING_BASE_PATH/$PARENT_NAME/servers" "data/post.json"
	assert_equal "$SC" 202
}
