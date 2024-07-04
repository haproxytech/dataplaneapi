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
    PARENT_NAME="test_backend_replace"
    resource_put "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/0" "data/replace/connect.json" "force_reload=true"
    assert_equal "$SC" 200
}

@test "tcp_checks: Replace a TCP check send" {
    PARENT_NAME="test_backend_replace"
    resource_put "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/1" "data/replace/send.json" "force_reload=true"
    assert_equal "$SC" 200
}

@test "tcp_checks: Replace a TCP check expect" {
    PARENT_NAME="test_backend_replace"
    resource_put "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/2" "data/replace/expect.json" "force_reload=true"
    assert_equal "$SC" 200
}

@test "tcp_checks: Replace a TCP check send-binary" {
    PARENT_NAME="test_backend_replace"
    resource_put "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks/3" "data/replace/send_binary.json" "force_reload=true"
    assert_equal "$SC" 200
}

@test "tcp_checks: replace all TCP Checks for a backend" {
    PARENT_NAME="test_backend_replace"
    resource_put "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks" "data/replace/replace-all.json"
    assert_equal "$SC" 202
    resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_checks"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 5
    assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace/replace-all.json")" ".")"

}

@test "tcp_checks: replace all TCP Checks for a defaults" {
    PARENT_NAME="mydefaults"
    resource_put "$_DEFAULTS_BASE_PATH/$PARENT_NAME/tcp_checks" "data/replace/replace-all.json"
    assert_equal "$SC" 202
    resource_get "$_DEFAULTS_BASE_PATH/$PARENT_NAME/tcp_checks"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 5
    assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace/replace-all.json")" ".")"

}
