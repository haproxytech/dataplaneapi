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
load '../../libs/version'

load 'utils/_helpers'

@test "servers: Delete a server" {
  PARENT_NAME="test_backend"
  for name in "server_01" "server_02" "server_03" "server_04"; do
    resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/servers/$name" "force_reload=true"
	  assert_equal "$SC" 204
  done
}

@test "servers: Delete a non existing server" {
  PARENT_NAME=test_backend
  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/servers/ghost"
  assert_equal "$SC" 404
}

@test "servers: Delete a server on a non existing backend" {
  PARENT_NAME=ghost
  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/servers/server_01"
  assert_equal "$SC" 404
}
