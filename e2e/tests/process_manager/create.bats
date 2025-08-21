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
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "process-manager: Create one new program" {
  if haproxy_version_ge "3.3"; then skip "programs have been removed in haproxy 3.3"; fi

  resource_post "$_PROGRAMS_BASE_PATH" "data/program.json" "force_reload=true"
  assert_equal "$SC" 201

  resource_get "$_PROGRAMS_BASE_PATH/test"
  assert_equal "$SC" 200
  assert_equal "test" "$(get_json_path "$BODY" ".name")"
  assert_equal "haproxy" "$(get_json_path "$BODY" ".user")"
  assert_equal "haproxy" "$(get_json_path "$BODY" ".group")"
}

@test "process-manager: Fail creating program with same name" {
  if haproxy_version_ge "3.3"; then skip; fi

  resource_post "$_PROGRAMS_BASE_PATH" "data/program_duplicated.json" "force_reload=true"
  assert_equal "$SC" 409
}

@test "process-manager: Fail creating program that isn't valid" {
  if haproxy_version_ge "3.3"; then skip; fi

  resource_post "$_PROGRAMS_BASE_PATH" "data/program_invalid.json" "force_reload=true"
  assert_equal "$SC" 422
}
