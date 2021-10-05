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

@test "global: Replace a global configuration" {
    resource_put "$_GLOBAL_BASE_PATH" "data/put.json" ""
	assert_equal "$SC" 202

	resource_get "$_GLOBAL_BASE_PATH" ""
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" '.data.maxconn')" "5000"
	assert_equal "$(get_json_path "$BODY" '.data.daemon')" "enabled"
	assert_equal "$(get_json_path "$BODY" '.data.pidfile')" "/var/run/haproxy/haproxy.pid"
	assert_equal "$(get_json_path "$BODY" '.data.runtime_apis[0].address')" "/var/run/haproxy/admin.sock"
	assert_equal "$(get_json_path "$BODY" '.data.runtime_apis[0].level')" "admin"
	assert_equal "$(get_json_path "$BODY" '.data.runtime_apis[0].mode')" "660"
	assert_equal "$(get_json_path "$BODY" '.data.runtime_apis[0].allow_0rtt')" "true"
	assert_equal "$(get_json_path "$BODY" '.data.runtime_apis[0].name')" "runtime_api_1"
}
