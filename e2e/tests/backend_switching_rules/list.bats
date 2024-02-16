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
load '../../libs/resource_client'
load '../../libs/get_json_path'
load '../../libs/haproxy_config_setup'
load '../../libs/version'

load 'utils/_helpers'

@test "backend_switching_rules: Replace a Backend Switching Rule" {
  resource_get "$_BSR_BASE_PATH" "frontend=test_frontend"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ". | length")" 2
	assert_equal "$(get_json_path "$BODY" ".[0].cond_test")" "{ req_ssl_sni -i first.example.com }"
	assert_equal "$(get_json_path "$BODY" ".[1].cond_test")" "{ req_ssl_sni -i second.example.com }"
}
