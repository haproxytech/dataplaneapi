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

@test "server_templates: Delete a server template 1" {
  PARENT_NAME="test_backend"
  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/server_templates/srv_google" "force_reload=true"
	assert_equal "$SC" 204
}

@test "server_templates: Delete a server template 2" {
  PARENT_NAME="test_backend"
  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/server_templates/srv_bing" "force_reload=true"
	assert_equal "$SC" 204
}

@test "server_templates: Delete a non existing server template" {
  PARENT_NAME="test_backend"
  resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/server_templates/ghost" "force_reload=true"
	assert_equal "$SC" 404
}
