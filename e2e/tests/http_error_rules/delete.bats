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
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "http_error_rules: Delete a HTTP Error Rule from frontend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	#
	# Deleting first
	#
	PARENT_NAME="test_frontend"
	resource_delete "$_FRONTEND_BASE_PATH/$PARENT_NAME/http_error_rules/0" "force_reload=true"
	assert_equal "$SC" 204
	#
	# Not found
	#
	resource_delete "$_FRONTEND_BASE_PATH/$PARENT_NAME/http_error_rules/0" "force_reload=true"
	assert_equal "$SC" 404
}

@test "http_error_rules: Delete a HTTP Error Rule from backend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	#
	# Deleting first
	#
	PARENT_NAME="test_backend"
	resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/http_error_rules/0" "force_reload=true"
	assert_equal "$SC" 204
	#
	# Deleting second
	#
	resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/http_error_rules/0" "force_reload=true"
	assert_equal "$SC" 204
	#
	# Not found
	#
	resource_delete "$_BACKEND_BASE_PATH/$PARENT_NAME/http_error_rules/0" "force_reload=true"
	assert_equal "$SC" 404
}

@test "http_error_rules: Delete a HTTP Error Rule from defaults" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	#
	# Deleting first
	#
	PARENT_NAME="mydefaults"
	resource_delete "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_error_rules/0" "force_reload=true"
	assert_equal "$SC" 204
	#
	# Not found
	#
	resource_delete "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_error_rules/0" "force_reload=true"
	assert_equal "$SC" 404
}
