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
load "../../libs/run_only"

setup() {
    run_only

    refute dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example.cfg'

    run dpa_curl_file_upload POST "/services/haproxy/spoe/spoe_files" "@${BATS_TEST_DIRNAME}/data/spoefile_example.cfg;filename=spoefile_example.cfg"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201
}

teardown() {
    run dpa_docker_exec 'rm -rf /etc/haproxy/spoe/spoefile_example.cfg'
}

@test "spoe_scopes: Add a spoe scope" {
    run dpa_curl POST "/services/haproxy/spoe/spoe_scopes?spoe=spoefile_example.cfg&version=1" /data/add_scope.txt
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201
}

@test "spoe_scopes: Refuse adding an existing spoe scope" {
    run dpa_curl POST "/services/haproxy/spoe/spoe_scopes?spoe=spoefile_example.cfg&version=2" /data/add_scope.txt
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 409
}

@test "spoe_scopes: Return an error when spoe file doesn't exists" {
    run dpa_curl POST "/services/haproxy/spoe/spoe_scopes?spoe=not_exists.cfg&version=3" /data/add_scope.txt
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 500
}

@test "spoe_scopes: Return an error when version not matched" {
    run dpa_curl POST "/services/haproxy/spoe/spoe_scopes?spoe=spoefile_example.cfg&version=10000" /data/new_scope.txt
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 409
}
