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

    PARENT_NAME="test_frontend"
	resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/http_after_response_rules"
	assert_equal "$SC" 200
	if haproxy_version_ge "2.8"; then
	    assert_equal "$(get_json_path "$BODY" ". | length")" 11
    else
        assert_equal "$(get_json_path "$BODY" ". | length")" 2
    fi
	assert_equal "$(get_json_path "$BODY" ".[0].type")" "add-header"
	assert_equal "$(get_json_path "$BODY" ".[0].hdr_name")" "X-Add-Frontend"
	assert_equal "$(get_json_path "$BODY" ".[0].cond")" "unless"
	assert_equal "$(get_json_path "$BODY" ".[0].cond_test")" "{ src 192.168.0.0/16 }"
	assert_equal "$(get_json_path "$BODY" ".[1].type")" "del-header"
	assert_equal "$(get_json_path "$BODY" ".[1].hdr_name")" "X-Del-Frontend"
	assert_equal "$(get_json_path "$BODY" ".[1].cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".[1].cond_test")" "{ src 10.1.0.0/16 }"
	if haproxy_version_ge "2.8"; then
        assert_equal "$(get_json_path "$BODY" ".[2].type")" "set-map"
        assert_equal "$(get_json_path "$BODY" ".[2].map_file")" "map.lst"
        assert_equal "$(get_json_path "$BODY" ".[2].map_keyfmt")" "%[src]"
        assert_equal "$(get_json_path "$BODY" ".[2].map_valuefmt")" "%[res.hdr(X-Value)]"
        assert_equal "$(get_json_path "$BODY" ".[3].type")" "del-map"
        assert_equal "$(get_json_path "$BODY" ".[3].map_file")" "map.lst"
        assert_equal "$(get_json_path "$BODY" ".[3].map_keyfmt")" "%[src]"
        assert_equal "$(get_json_path "$BODY" ".[3].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".[3].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".[4].type")" "del-acl"
        assert_equal "$(get_json_path "$BODY" ".[4].acl_file")" "map.lst"
        assert_equal "$(get_json_path "$BODY" ".[4].acl_keyfmt")" "%[src]"
        assert_equal "$(get_json_path "$BODY" ".[4].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".[4].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".[5].type")" "sc-inc-gpc"
        assert_equal "$(get_json_path "$BODY" ".[5].sc_id")" "1"
        assert_equal "$(get_json_path "$BODY" ".[5].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".[5].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".[6].type")" "sc-inc-gpc0"
        assert_equal "$(get_json_path "$BODY" ".[6].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".[6].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".[7].type")" "sc-inc-gpc1"
        assert_equal "$(get_json_path "$BODY" ".[7].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".[7].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".[8].type")" "sc-set-gpt0"
        assert_equal "$(get_json_path "$BODY" ".[8].sc_id")" "1"
        assert_equal "$(get_json_path "$BODY" ".[8].sc_expr")" "hdr(Host),lower"
        assert_equal "$(get_json_path "$BODY" ".[8].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".[8].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".[9].type")" "sc-set-gpt0"
        assert_equal "$(get_json_path "$BODY" ".[9].sc_id")" "1"
        assert_equal "$(get_json_path "$BODY" ".[9].sc_int")" "20"
        assert_equal "$(get_json_path "$BODY" ".[9].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".[9].cond_test")" "FALSE"
        assert_equal "$(get_json_path "$BODY" ".[10].type")" "set-log-level"
        assert_equal "$(get_json_path "$BODY" ".[10].log_level")" "silent"
        assert_equal "$(get_json_path "$BODY" ".[10].cond")" "if"
        assert_equal "$(get_json_path "$BODY" ".[10].cond_test")" "FALSE"
    fi
}

@test "http_after_response_rules: Return one HTTP After Response Rule from backend" {
	if [[ "$HAPROXY_VERSION" == "2.1" ]]; then
		skip "http-after-response is not supported in HAProxy 2.1"
	fi

    PARENT_NAME="test_backend"
	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/http_after_response_rules"
	assert_equal "$SC" 200
	assert_equal 2 "$(get_json_path "$BODY" ". | length")"
	assert_equal "$(get_json_path "$BODY" ".[0].type")" "add-header"
	assert_equal "$(get_json_path "$BODY" ".[0].hdr_name")" "X-Add-Backend"
	assert_equal "$(get_json_path "$BODY" ".[0].cond")" "unless"
	assert_equal "$(get_json_path "$BODY" ".[0].cond_test")" "{ src 192.168.0.0/16 }"
	assert_equal "$(get_json_path "$BODY" ".[1].type")" "del-header"
	assert_equal "$(get_json_path "$BODY" ".[1].hdr_name")" "X-Del-Backend"
	assert_equal "$(get_json_path "$BODY" ".[1].cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".[1].cond_test")" "{ src 10.1.0.0/16 }"
}
