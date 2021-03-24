#!/usr/bin/env bats
#
# Copyright 2019 HAProxy Technologies
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
load "../../libs/run_only"

setup() {
    run_only

    refute dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example.cfg'

    # adding a spoe file case handled here
    run dpa_curl_file_upload POST "/services/haproxy/spoe/spoe_files" "@${BATS_TEST_DIRNAME}/spoefile_example.cfg;filename=spoefile_example.cfg"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201

    assert dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example.cfg'
}

teardown() {
    run dpa_docker_exec 'rm -rf /etc/haproxy/spoe/spoefile_example.cfg'
}

@test "spoe: Get a list of spoefiles" {

    run dpa_curl GET "/services/haproxy/spoe/spoe_files"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal $(get_json_path "$BODY" '.|length') 1

    assert_equal $(get_json_path "$BODY" '.[0]') "/etc/haproxy/spoe/spoefile_example.cfg"
}

@test "spoe: Get a spoefile contents" {

    run dpa_curl_download GET "/services/haproxy/spoe/spoe_files/spoefile_example.cfg"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert dpa_diff_docker_file '/etc/haproxy/spoe/spoefile_example.cfg' "spoefile_example.cfg"
}

@test "spoe: Try to get unavailable spoefile contents" {

    run dpa_curl GET "/services/haproxy/spoe/spoe_files/not_here.cfg"
    assert_success

    dpa_curl_status_body_safe '$output'
    assert_equal $SC 404
}

@test "spoe: Delete a spoefile" {

    run dpa_curl DELETE "/services/haproxy/spoe/spoe_files/spoefile_example.cfg"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 204

    refute dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example.cfg'
}
