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

@test "log_targets: Add a new Log Target to frontend" {
  resource_post "$_LOG_TRAGET_BASE_PATH/0" "data/nolog.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 201
}

@test "log_targets: Add a new Log Target to backend" {
	resource_post "$_LOG_TRAGET_BASE_PATH/0" "data/nolog.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 201
}
