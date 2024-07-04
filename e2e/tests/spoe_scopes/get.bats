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
load '../../libs/version_spoe'

load 'utils/_helpers'

setup() {
    SPOE_FILE="spoefile_example2.cfg"

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

@test "spoe_scopes: Get one spoe scope" {
    PARENT_NAME="spoefile_example2.cfg"
    resource_get "$_SPOE_BASE_PATH/${PARENT_NAME}/scopes/%5Bip-reputation%5D"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "${BODY}" ".")" "[ip-reputation]"
}

@test "spoe_scopes: Return an error when trying to get non existing spoe scope" {
    PARENT_NAME="spoefile_example2.cfg"
    resource_get "$_SPOE_BASE_PATH/${PARENT_NAME}/scopes/%5Bnot-exists%5D" "spoe=spoefile_example2.cfg"
    assert_equal "$SC" 404
}
