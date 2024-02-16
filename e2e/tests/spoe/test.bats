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
load '../../libs/resource_client'
load "../../libs/run_only"
load '../../libs/version'

load 'utils/_helpers'

setup() {
    run_only

    refute dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example.cfg'

    # adding a spoe file case handled here
    run dpa_curl_file_upload POST "$_SPOE_FILES_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/spoefile_example.cfg;filename=spoefile_example.cfg"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201

    assert dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example.cfg'
}

teardown() {
    run dpa_docker_exec 'rm -rf /etc/haproxy/spoe/spoefile_example.cfg'
}

@test "spoe: Get a list of spoefiles" {
    resource_get "$_SPOE_FILES_BASE_PATH"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "$BODY" '. | length')" 1
    assert_equal "$(get_json_path "$BODY" '.[0]')" "/etc/haproxy/spoe/spoefile_example.cfg"
}

@test "spoe: Get a spoefile contents" {
    run dpa_curl_download GET "$_SPOE_FILES_BASE_PATH/spoefile_example.cfg"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal "$SC" 200

    assert dpa_diff_docker_file '/etc/haproxy/spoe/spoefile_example.cfg' "data/spoefile_example.cfg"
}

@test "spoe: Try to get unavailable spoefile contents" {
    resource_get "$_SPOE_FILES_BASE_PATH/not_here.cfg"
    dpa_curl_status_body_safe '$output'
    assert_equal "$SC" 404
}

@test "spoe: Delete a spoefile" {
    resource_delete "$_SPOE_FILES_BASE_PATH/spoefile_example.cfg"
    assert_equal "$SC" 204

    refute dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example.cfg'
}
