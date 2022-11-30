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

@test "groups: Replace a group" {
    resource_put "$_GROUPS_BASE_PATH/zion" "data/replace.json" "userlist=replace_test&force_reload=true"
    assert_equal "$SC" 200
}

@test "groups: Replace a group with malformed data" {
    resource_put "$_GROUPS_BASE_PATH/io" "data/empty.json" "userlist=replace_test&force_reload=true"
    assert_equal "$SC" 422
}

@test "groups: Replace a non existing group" {
    resource_put "$_GROUPS_BASE_PATH/fake" "data/replace.json" "userlist=replace_test&force_reload=true"
    assert_equal "$SC" 404
}

@test "groups: Replace a non existing group in a non existing user list" {
    resource_put "$_GROUPS_BASE_PATH/1000" "data/replace.json" "userlist=fake&force_reload=true"
    assert_equal "$SC" 404
}
