#!/usr/bin/env bats
#
# Copyright 2023 HAProxy Technologies
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
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "global: Replace a global configuration (>=2.8)" {
  if haproxy_version_ge "2.8"
  then
	resource_put "$_GLOBAL_BASE_PATH" "data/put_2.8.json" ""
	assert_equal "$SC" 202

	resource_get "$_GLOBAL_BASE_PATH" ""
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" '.maxconn')" "5000"
	assert_equal "$(get_json_path "$BODY" '.daemon')" "enabled"
	assert_equal "$(get_json_path "$BODY" '.pidfile')" "/var/run/haproxy.pid"
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].address')" "/var/lib/haproxy/stats"
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].level')" "admin"
	assert_equal "$(get_json_path "$BODY" '.tune_options.h2_be_initial_window_size')" "20"
    assert_equal "$(get_json_path "$BODY" '.tune_options.h2_be_max_concurrent_streams')" "21"
	assert_equal "$(get_json_path "$BODY" '.tune_options.h2_fe_initial_window_size')" "22"
    assert_equal "$(get_json_path "$BODY" '.tune_options.h2_fe_max_concurrent_streams')" "23"
  fi
}
