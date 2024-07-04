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
load '../../libs/resource_client'
load '../../libs/haproxy_config_setup'
load '../../libs/version'

@test "x_issue_int_247: runtime admin-state of backend app2" {
    PARENT_NAME="bug_int_247"
    resource_get "/services/haproxy/runtime/backends/$PARENT_NAME/servers/app2"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "$BODY" '.admin_state')" "maint"
}

@test "x_issue_int_247: runtime admin-state of backend app1" {
    PARENT_NAME="bug_int_247"
    resource_get "/services/haproxy/runtime/backends/$PARENT_NAME/servers/app1"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "$BODY" '.admin_state')" "ready"
}

@test "x_issue_int_247: admin-state always reports admin_state of maint if disabled keyword is used" {
    PARENT_NAME="bug_int_247"
    resource_put "/services/haproxy/runtime/backends/$PARENT_NAME/servers/app2" "data/enable.json"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "$BODY" '.admin_state')" "ready"
}
