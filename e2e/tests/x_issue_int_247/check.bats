#!/usr/bin/env bats
#
# Copyright 2020 HAProxy Technologies
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

@test "x_issue_int_247: runtime admin-state of backend app2" {
    #run dpa_docker_exec 'started" >> /etc/haproxy/log.log'
    #assert_success

    run dpa_curl GET "/services/haproxy/runtime/servers/app2?backend=bug_int_247"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal $(get_json_path "$BODY" '.admin_state') "maint"
}

@test "x_issue_int_247: runtime admin-state of backend app1" {
    run dpa_curl GET "/services/haproxy/runtime/servers/app1?backend=bug_int_247"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal $(get_json_path "$BODY" '.admin_state') "ready"
}

@test "x_issue_int_247: admin-state always reports admin_state of maint if disabled keyword is used" {
    run dpa_curl PUT "/services/haproxy/runtime/servers/app2?backend=bug_int_247" enable.json
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal $(get_json_path "$BODY" '.admin_state') "ready"
}
