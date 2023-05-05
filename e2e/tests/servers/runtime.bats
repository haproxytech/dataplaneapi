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
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "servers: Return an array of runtime servers' settings" {
  resource_get "$_RUNTIME_SERVER_BASE_PATH" "backend=test_backend"
	assert_equal "$SC" 200

	for name in "server_01" "server_02" "server_03" "server_04"; do
  	assert_equal "$(get_json_path "$BODY" ".[] | select(.name | contains(\"$name\") ).name")" "$name"
  done
}

@test "servers: Replace server transient settings" {
  resource_put "$_RUNTIME_SERVER_BASE_PATH/server_01" "data/transient.json" "backend=test_backend"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.name')" "server_01"
}

@test "servers: Return one server runtime settings" {
  resource_get "$_RUNTIME_SERVER_BASE_PATH/server_01" "backend=test_backend"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.name')" "server_01"
}

@test "servers: Add a server via runtime" {
  haproxy_version_ge 2.6 || skip "requires HAProxy 2.6+"

  resource_post "$_RUNTIME_SERVER_BASE_PATH" "data/runtime_add_server.json" "backend=test_backend"
  assert_equal "$SC" 201
  assert_equal "$(get_json_path "$BODY" '.name')" "rt_server"

  resource_get "$_RUNTIME_SERVER_BASE_PATH/rt_server" "backend=test_backend"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.name')" "rt_server"
  assert_equal "$(get_json_path "$BODY" '.address')" "10.11.12.13"
  assert_equal "$(get_json_path "$BODY" '.port')" 8088
}

@test "servers: Add an existing server via runtime" {
  haproxy_version_ge 2.6 || skip "requires HAProxy 2.6+"

  resource_post "$_RUNTIME_SERVER_BASE_PATH" "data/runtime_add_server.json" "backend=test_backend"
  assert_equal "$SC" 201

  resource_post "$_RUNTIME_SERVER_BASE_PATH" "data/runtime_add_server.json" "backend=test_backend"
  assert_equal "$SC" 409
}

@test "servers: Add a server to a wrong backend via runtime" {
  haproxy_version_ge 2.6 || skip "requires HAProxy 2.6+"

  resource_post "$_RUNTIME_SERVER_BASE_PATH" "data/runtime_add_server.json" "backend=does_not_exist"
  assert_equal "$SC" 404
}

@test "servers: Delete a server via runtime" {
  haproxy_version_ge 2.6 || skip "requires HAProxy 2.6+"

  resource_post "$_RUNTIME_SERVER_BASE_PATH" "data/runtime_add_server.json" "backend=test_backend"
  assert_equal "$SC" 201

  resource_delete "$_RUNTIME_SERVER_BASE_PATH/rt_server" "backend=test_backend"
  assert_equal "$SC" 204
}

@test "servers: Delete a non-existant server via runtime" {
  haproxy_version_ge 2.6 || skip "requires HAProxy 2.6+"

  resource_delete "$_RUNTIME_SERVER_BASE_PATH/rt_server1" "backend=test_backend"
  assert_equal "$SC" 404

  # wrong backend
  resource_delete "$_RUNTIME_SERVER_BASE_PATH/rt_server2" "backend=does_not_exist"
  assert_equal "$SC" 404
}
