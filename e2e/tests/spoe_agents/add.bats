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
load '../../libs/debug'

setup() {
    SPOE_FILE="spoefile_example.cfg"

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

@test "spoe_agents: Add a spoe agent" {
    resource_post "$_SPOE_AGENTS_BASE_PATH" "data/post.json" "spoe=spoefile_example.cfg&scope=\[ip-reputation\]"
    assert_equal "$SC" 201

    resource_get "$_SPOE_AGENTS_BASE_PATH/post_agent1" "scope=\[ip-reputation\]&spoe=spoefile_example.cfg"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME"/data/post.json)" ".")"

    # refuse adding an existing spoe agent
    resource_post "$_SPOE_AGENTS_BASE_PATH" "/data/post.json" "spoe=spoefile_example.cfg&scope=%5Bip-reputation%5D"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 409
}
