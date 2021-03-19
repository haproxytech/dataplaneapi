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

    refute dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example2.cfg'

    run dpa_curl_file_upload POST "/services/haproxy/spoe/spoe_files" "@${BATS_TEST_DIRNAME}/data/spoefile_example2.cfg;filename=spoefile_example2.cfg"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201
}

teardown() {
    run dpa_docker_exec 'rm -rf /etc/haproxy/spoe/spoefile_example2.cfg'
}

@test "spoe_groups: Delete a spoe group" {
    run dpa_curl DELETE "/services/haproxy/spoe/spoe_groups/newgroup?scope=%5Bip-reputation%5D&spoe=spoefile_example2.cfg&version=1"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 204

    run dpa_curl GET "/services/haproxy/spoe/spoe_groups/newgroup?scope=%5Bip-reputation%5D&spoe=spoefile_example2.cfg"
    dpa_curl_status_body '$output'
    assert_equal $SC 404
}

@test "spoe_groups: Return an error when trying to delete non existing spoe group" {
    run dpa_curl DELETE "/services/haproxy/spoe/spoe_groups/not-exists?scope=%5Bip-reputation%5D&spoe=spoefile_example2.cfg&version=1"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 404
}
