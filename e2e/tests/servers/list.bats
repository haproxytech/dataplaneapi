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

@test "servers: Return an array of servers" {
  resource_get "$_SERVER_BASE_PATH" "backend=test_backend"
	assert_equal "$SC" 200

	assert_equal "$(get_json_path "$BODY" ".data | length")" "5"

    INDEX=0
	for name in "server_01" "server_02" "server_03" "server_ipv6" "server_04"; do
  	assert_equal "$(get_json_path "$BODY" ".data[] | select(.name | contains(\"$name\") ).name")" "$name"
  	if [[ "$(get_json_path "$BODY" ".data[] | select(.name | contains(\"$name\") ).name")" == "server_04" ]]; then
  	    assert_equal "$(get_json_path "$BODY" ".data[${INDEX}].check")" "enabled"
  	    assert_equal "$(get_json_path "$BODY" ".data[${INDEX}].resolve_opts")" "allow-dup-ip,ignore-weight"
  	    assert_equal "$(get_json_path "$BODY" ".data[${INDEX}].\"resolve-net\"")" "10.0.0.0/8,10.200.200.0/12"
    fi
    let INDEX=${INDEX}+1
  done
}
