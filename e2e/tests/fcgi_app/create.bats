#!/usr/bin/env bats
#
# Copyright 2022 HAProxy Technologies
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

@test "fcgi-app: Create one new app" {
  resource_post "$_FCGIAPP_BASE_PATH" "data/app.json" "force_reload=true"
  assert_equal "$SC" 201

  resource_get "$_FCGIAPP_BASE_PATH/app_created"
  assert_equal "$SC" 200
  assert_equal "app_created" "$(get_json_path "$BODY" ".name")"
  assert_equal "/path/to/chroot" "$(get_json_path "$BODY" ".docroot")"
}

@test "fcgi-app: Fail creating app with same name" {
  resource_post "$_FCGIAPP_BASE_PATH" "data/app_duplicated.json" "force_reload=true"
  assert_equal "$SC" 409
}

@test "fcgi-app: Fail creating app that isn't valid" {
  resource_post "$_FCGIAPP_BASE_PATH" "data/app_unvalid.json" "force_reload=true"
  assert_equal "$SC" 422
}
