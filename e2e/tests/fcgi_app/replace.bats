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

@test "fcgi-app: Replace one app" {
  resource_put "$_FCGIAPP_BASE_PATH/test_1" "data/app_duplicated.json" "force_reload=true"
  assert_equal "$SC" 200

  resource_get "$_FCGIAPP_BASE_PATH/test_1"
  assert_equal "$SC" 200
  assert_equal "test_1" "$(get_json_path "$BODY" ".name")"
  assert_equal "$(get_json_path "${BODY}" ".log_stder | length")" 0
}

@test "fcgi-app: Fail replacing app that doesn't exist" {
  resource_put "$_FCGIAPP_BASE_PATH/i_am_not_here" "data/app_duplicated.json" "force_reload=true"
  assert_equal "$SC" 409
}
