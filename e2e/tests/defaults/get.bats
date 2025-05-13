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


@test "defaults: Return a list of defaults configurations" {
  resource_get "$_DEFAULTS_BASE_PATH"
  assert_equal "$SC" 200
}

@test "defaults: Return a defaults configuration" {
  resource_get "$_DEFAULTS_BASE_PATH/unnamed_defaults_1"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.name')" "unnamed_defaults_1"
}

@test "defaults: Return a named defaults configuration that does not exist" {
  resource_get "$_DEFAULTS_BASE_PATH/nothing_to_see_here"
  assert_equal "$SC" 404
}
