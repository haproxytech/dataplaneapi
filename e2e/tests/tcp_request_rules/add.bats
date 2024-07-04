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

load 'utils/_helpers'

@test "tcp_request_rules: Add a new TCP Request Rule to frontend" {
  PARENT_NAME="test_frontend"
  resource_post "$_FRONTEND_BASE_PATH/$PARENT_NAME/tcp_request_rules/0" "data/accept.json" "force_reload=true"
	assert_equal "$SC" 201
}

@test "tcp_request_rules: Add a new TCP Request Rule to backend" {
  PARENT_NAME="test_backend"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_request_rules/0" "data/accept.json" "force_reload=true"
	assert_equal "$SC" 201
}
