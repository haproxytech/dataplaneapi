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

@test "tcp_checks: Replace a TCP check connect" {
    resource_put "$_TCP_CHECKS_BASE_PATH/0" "data/replace/connect.json" "parent_type=backend&parent_name=test_backend_replace&force_reload=true"
    assert_equal "$SC" 200
}

@test "tcp_checks: Replace a TCP check send" {
    resource_put "$_TCP_CHECKS_BASE_PATH/1" "data/replace/send.json" "parent_type=backend&parent_name=test_backend_replace&force_reload=true"
    assert_equal "$SC" 200
}

@test "tcp_checks: Replace a TCP check expect" {
    resource_put "$_TCP_CHECKS_BASE_PATH/2" "data/replace/expect.json" "parent_type=backend&parent_name=test_backend_replace&force_reload=true"
    assert_equal "$SC" 200
}

@test "tcp_checks: Replace a TCP check send-binary" {
    resource_put "$_TCP_CHECKS_BASE_PATH/3" "data/replace/send_binary.json" "parent_type=backend&parent_name=test_backend_replace&force_reload=true"
    assert_equal "$SC" 200
}

@test "tcp_checks: replace all TCP Checks for a backend" {
    resource_put "$_TCP_CHECKS_BASE_PATH" "data/replace/replace-all.json" "parent_name=test_backend_replace&parent_type=backend"
    assert_equal "$SC" 202
    resource_get "$_TCP_CHECKS_BASE_PATH" "parent_name=test_backend_replace&parent_type=backend"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 5
    assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace/replace-all.json")" ".")"

}

@test "tcp_checks: replace all TCP Checks for a defaults" {
    resource_put "$_TCP_CHECKS_BASE_PATH" "data/replace/replace-all.json" "parent_name=mydefaults&parent_type=defaults"
    assert_equal "$SC" 202
    resource_get "$_TCP_CHECKS_BASE_PATH" "parent_name=mydefaults&parent_type=defaults"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 5
    assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace/replace-all.json")" ".")"

}
