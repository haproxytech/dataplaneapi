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
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "tcp_response_rules: Return one TCP Response Rule from backend" {
  resource_get "$_TCP_RES_RULES_CERTS_BASE_PATH" "backend=test_backend"
	assert_equal "$SC" 200

    if haproxy_version_ge "2.8"; then
        assert_equal "$(get_json_path "${BODY}" ". | length")" 3
	else
	    assert_equal "$(get_json_path "${BODY}" ". | length")" 2
    fi
	assert_equal "$(get_json_path "$BODY" ".[] | select(.action | contains(\"accept\") ).action")" "accept"
	assert_equal "$(get_json_path "$BODY" ".[] | select(.action | contains(\"reject\") ).action")" "reject"
	if haproxy_version_ge "2.8"; then
	    assert_equal "$(get_json_path "$BODY" ".[] | select(.action | contains(\"sc-add-gpc\") ).action")" "sc-add-gpc"
    fi
}
