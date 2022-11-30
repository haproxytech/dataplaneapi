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

@test "groups: Add a group" {
  resource_post "$_GROUPS_BASE_PATH" "data/post.json" "userlist=first&=force_reload=true"
  assert_equal "$SC" 202
}

@test "groups: Add a group to a non existing user list" {
  resource_post "$_GROUPS_BASE_PATH" "data/post.json" "userlist=fake&=force_reload=true"
  assert_equal "$SC" 404
}

@test "groups: Add a malformed group" {
  resource_post "$_GROUPS_BASE_PATH" "data/empty.json" "userlist=first&force_reload=true"
  assert_equal "$SC" 422
}
