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

load 'utils/_helpers'

@test "table: Return one table from peers" {
  resource_get "$_REQ_RULES_BASE_PATH/t1" "peer_section=mycluster"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".data.name")" "t1"
	assert_equal "$(get_json_path "$BODY" ".data.type")" "string"
	assert_equal "$(get_json_path "$BODY" ".data.type_len")" "1000"
	assert_equal "$(get_json_path "$BODY" ".data.size")" "200k"
	assert_equal "$(get_json_path "$BODY" ".data.expire")" "5m"
	assert_equal "$(get_json_path "$BODY" ".data.no_purge")" "true"
	assert_equal "$(get_json_path "$BODY" ".data.store")" "gpc0,conn_rate(30s)"
}

