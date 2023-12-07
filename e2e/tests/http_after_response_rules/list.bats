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

@test "http_after_response_rules: Return an array of all HTTP After Response Rules from frontend" {
	if [[ "$HAPROXY_VERSION" == "2.1" ]]; then
		skip "http-after-response is not supported in HAProxy 2.1"
	fi

	resource_get "$_RES_RULES_BASE_PATH" "parent_type=frontend&parent_name=test_frontend"
	assert_equal "$SC" 200
	if haproxy_version_ge "2.8"; then
	    assert_equal "$(get_json_path "$BODY" ".data | length")" 11
    else
        assert_equal "$(get_json_path "$BODY" ".data | length")" 2
    fi
	assert_equal "$(get_json_path "$BODY" ".data[0].type")" "add-header"
	assert_equal "$(get_json_path "$BODY" ".data[0].hdr_name")" "X-Add-Frontend"
	assert_equal "$(get_json_path "$BODY" ".data[0].cond")" "unless"
	assert_equal "$(get_json_path "$BODY" ".data[0].cond_test")" "{ src 192.168.0.0/16 }"
	assert_equal "$(get_json_path "$BODY" ".data[1].type")" "del-header"
	assert_equal "$(get_json_path "$BODY" ".data[1].hdr_name")" "X-Del-Frontend"
	assert_equal "$(get_json_path "$BODY" ".data[1].cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".data[1].cond_test")" "{ src 10.1.0.0/16 }"
	if haproxy_version_ge "2.8"; then
        assert_equal "$(get_json_path "$BODY" ".data[2].type")" "set-map"
        assert_equal "$(get_json_path "$BODY" ".data[2].map_file")" "map.lst"
        assert_equal "$(get_json_path "$BODY" ".data[2].map_keyfmt")" "%[src]"
        assert_equal "$(get_json_path "$BODY" ".data[2].map_valuefmt")" "%[res.hdr(X-Value)]"
        assert_equal "$(get_json_path "$BODY" ".data[3].type")" "del-map"
        assert_equal "$(get_json_path "$BODY" ".data[3].map_file")" "map.lst"
        assert_equal "$(get_json_path "$BODY" ".data[3].map_keyfmt")" "%[src]"
        assert_equal "$(get_json_path "$BODY" ".data[3].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".data[3].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".data[4].type")" "del-acl"
        assert_equal "$(get_json_path "$BODY" ".data[4].acl_file")" "map.lst"
        assert_equal "$(get_json_path "$BODY" ".data[4].acl_keyfmt")" "%[src]"
        assert_equal "$(get_json_path "$BODY" ".data[4].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".data[4].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".data[5].type")" "sc-inc-gpc"
        assert_equal "$(get_json_path "$BODY" ".data[5].sc_id")" "1"
        assert_equal "$(get_json_path "$BODY" ".data[5].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".data[5].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".data[6].type")" "sc-inc-gpc0"
        assert_equal "$(get_json_path "$BODY" ".data[6].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".data[6].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".data[7].type")" "sc-inc-gpc1"
        assert_equal "$(get_json_path "$BODY" ".data[7].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".data[7].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".data[8].type")" "sc-set-gpt0"
        assert_equal "$(get_json_path "$BODY" ".data[8].sc_id")" "1"
        assert_equal "$(get_json_path "$BODY" ".data[8].sc_expr")" "hdr(Host),lower"
        assert_equal "$(get_json_path "$BODY" ".data[8].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".data[8].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".data[9].type")" "sc-set-gpt0"
        assert_equal "$(get_json_path "$BODY" ".data[9].sc_id")" "1"
        assert_equal "$(get_json_path "$BODY" ".data[9].sc_int")" "20"
        assert_equal "$(get_json_path "$BODY" ".data[9].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".data[9].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".data[10].type")" "set-log-level"
        assert_equal "$(get_json_path "$BODY" ".data[10].log_level")" "silent"
        assert_equal "$(get_json_path "$BODY" ".data[10].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".data[10].cond_test")" "FALSE"
    fi
}

@test "http_after_response_rules: Return one HTTP After Response Rule from backend" {
	if [[ "$HAPROXY_VERSION" == "2.1" ]]; then
		skip "http-after-response is not supported in HAProxy 2.1"
	fi

	resource_get "$_RES_RULES_BASE_PATH" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
	assert_equal 2 "$(get_json_path "$BODY" ".data | length")"
	assert_equal "$(get_json_path "$BODY" ".data[0].type")" "add-header"
	assert_equal "$(get_json_path "$BODY" ".data[0].hdr_name")" "X-Add-Backend"
	assert_equal "$(get_json_path "$BODY" ".data[0].cond")" "unless"
	assert_equal "$(get_json_path "$BODY" ".data[0].cond_test")" "{ src 192.168.0.0/16 }"
	assert_equal "$(get_json_path "$BODY" ".data[1].type")" "del-header"
	assert_equal "$(get_json_path "$BODY" ".data[1].hdr_name")" "X-Del-Backend"
	assert_equal "$(get_json_path "$BODY" ".data[1].cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".data[1].cond_test")" "{ src 10.1.0.0/16 }"
}
