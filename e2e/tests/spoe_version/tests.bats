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

load 'utils/_helpers'

setup() {

    SPOE_FILE="spoefile_example.cfg"
    refute dpa_docker_exec 'ls /etc/haproxy/spoe/spoefile_example.cfg'

    run dpa_curl_file_upload POST "/services/haproxy/spoe/spoe_files" "@${BATS_TEST_DIRNAME}/data/spoefile_example.cfg;filename=spoefile_example.cfg"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201
}

teardown() {
    run dpa_docker_exec 'rm -rf /etc/haproxy/spoe/spoefile_example.cfg'
}

@test "spoe_version: Get a spoe version" {
    PARENT_NAME="spoefile_example.cfg"
    run dpa_curl GET "$_SPOE_BASE_PATH/$PARENT_NAME/version"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal $(get_json_path "${BODY}" ".") 12
}

@test "spoe_version: Return error when getting version for non existing spoe transaction" {
    PARENT_NAME="spoefile_example.cfg"
    run dpa_curl GET "$_SPOE_BASE_PATH/$PARENT_NAME/version?transaction_id=263166c2-3093-40ff-a750-4a4d114dfd99"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 400
}
