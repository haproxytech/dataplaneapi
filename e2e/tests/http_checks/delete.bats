#!/usr/bin/env bats
#
# Copyright 2022 HAProxy Technologies
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

@test "http_checks: Delete a HTTP Check from defaults" {
  #
  # Deleting first
  #
  resource_delete "$_CHECKS_BASE_PATH/0" "parent_type=defaults&parent_name=mydefaults&force_reload=true"
	assert_equal "$SC" 204
	#
  # Deleting second
  #
	resource_delete "$_CHECKS_BASE_PATH/0" "parent_type=defaults&parent_name=mydefaults&force_reload=true"
	assert_equal "$SC" 204
	#
  # Not found
  #
  resource_delete "$_CHECKS_BASE_PATH/0" "parent_type=defaults&parent_name=mydefaults&force_reload=true"
	assert_equal "$SC" 404
}

@test "http_checks: Delete a HTTP Check from backend" {
 	#
  # Deleting second
  #
	resource_delete "$_CHECKS_BASE_PATH/1" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 204
	#
  # Deleting first
  #
	resource_delete "$_CHECKS_BASE_PATH/0" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 204
	#
  # Not found
  #
	resource_delete "$_CHECKS_BASE_PATH/0" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 404
}
