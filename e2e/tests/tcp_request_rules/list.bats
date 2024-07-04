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
  PARENT_NAME="test_frontend"
  resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/tcp_request_rules"
	assert_equal "$SC" 200

    if haproxy_version_ge "2.8"; then
        assert_equal "$(get_json_path "${BODY}" ". | length")" 5
    else
        assert_equal "$(get_json_path "${BODY}" ". | length")" 2
    fi
	assert_equal "$(get_json_path "$BODY" ".[0].type")" "inspect-delay"
	assert_equal "$(get_json_path "$BODY" ".[1].type")" "content"
	assert_equal "$(get_json_path "$BODY" ".[1].action")" "accept"
	if haproxy_version_ge "2.8"; then
        assert_equal "$(get_json_path "$BODY" ".[2].type")" "connection"
        assert_equal "$(get_json_path "$BODY" ".[2].action")" "sc-add-gpc"
        assert_equal "$(get_json_path "$BODY" ".[3].type")" "session"
        assert_equal "$(get_json_path "$BODY" ".[3].action")" "sc-add-gpc"
        assert_equal "$(get_json_path "$BODY" ".[4].type")" "content"
        assert_equal "$(get_json_path "$BODY" ".[4].action")" "sc-add-gpc"
    fi
}

@test "tcp_request_rules: Return an array of all TCP Request Rules from backend" {
  PARENT_NAME="test_backend"
  resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_request_rules"
	assert_equal "$SC" 200

  assert_equal "$(get_json_path "$BODY" ".[] | select(.type | contains(\"inspect-delay\") ).type")" "inspect-delay"
  assert_equal "$(get_json_path "$BODY" ".[] | select(.type | contains(\"content\") ).type")" "content"
}
