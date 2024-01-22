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
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "tcp_request_rules: Return an array of all TCP Request Rules from frontend" {
  resource_get "$_TCP_REQ_RULES_CERTS_BASE_PATH" "parent_type=frontend&parent_name=test_frontend"
	assert_equal "$SC" 200

    if haproxy_version_ge "2.8"; then
        assert_equal "$(get_json_path "${BODY}" ".data | length")" 5
    else
        assert_equal "$(get_json_path "${BODY}" ".data | length")" 2
    fi
	assert_equal "$(get_json_path "$BODY" ".data[0].type")" "inspect-delay"
	assert_equal "$(get_json_path "$BODY" ".data[1].type")" "content"
	assert_equal "$(get_json_path "$BODY" ".data[1].action")" "accept"
	if haproxy_version_ge "2.8"; then
        assert_equal "$(get_json_path "$BODY" ".data[2].type")" "connection"
        assert_equal "$(get_json_path "$BODY" ".data[2].action")" "sc-add-gpc"
        assert_equal "$(get_json_path "$BODY" ".data[3].type")" "session"
        assert_equal "$(get_json_path "$BODY" ".data[3].action")" "sc-add-gpc"
        assert_equal "$(get_json_path "$BODY" ".data[4].type")" "content"
        assert_equal "$(get_json_path "$BODY" ".data[4].action")" "sc-add-gpc"
    fi
}

@test "tcp_request_rules: Return one TCP Request Rule from backend" {
  resource_get "$_TCP_REQ_RULES_CERTS_BASE_PATH" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200

  assert_equal "$(get_json_path "$BODY" ".data[] | select(.type | contains(\"inspect-delay\") ).type")" "inspect-delay"
  assert_equal "$(get_json_path "$BODY" ".data[] | select(.type | contains(\"content\") ).type")" "content"
}
