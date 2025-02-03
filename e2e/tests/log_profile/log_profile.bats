#!/usr/bin/env bats
#
# Copyright 2025 HAProxy Technologies
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
load '../../libs/debug'
load '../../libs/get_json_path'
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "log_profile: all tests (>=3.1)" {
  haproxy_version_ge "3.1" || skip

  debug "log_profile: create a new section"
  resource_post "$_LOG_PROFILE_PATH" "data/new_profile.json" "force_reload=true"
  assert_equal "$SC" "201"

  debug "log_profile: get a section"
  resource_get "$_LOG_PROFILE_PATH/$_PROFILE_NAME"
  assert_equal "$SC" "200"
  assert_equal "$_PROFILE_NAME" "$(get_json_path "$BODY" .name)"
  assert_equal "tag1" "$(get_json_path "$BODY" .log_tag)"
  assert_equal "enabled" "$(get_json_path "$BODY" .steps.[0].drop)"
  assert_equal "any" "$(get_json_path "$BODY" .steps.[2].step)"

  debug "log_profile: edit a section"
  resource_put "$_LOG_PROFILE_PATH/$_PROFILE_NAME" "data/edit_profile.json" "force_reload=true"
  assert_equal "$SC" "200"
  resource_get "$_LOG_PROFILE_PATH/$_PROFILE_NAME"
  assert_equal "tag2" "$(get_json_path "$BODY" .log_tag)"
  assert_equal "disabled" "$(get_json_path "$BODY" .steps.[0].drop)"
  assert_equal 2 "$(get_json_path "$BODY" '.steps|length')"

  debug "log_profile: get a list of sections"
  resource_get "$_LOG_PROFILE_PATH"
  assert_equal "$SC" "200"
  assert_equal "$_PROFILE_NAME" "$(get_json_path "$BODY" .[0].name)"

  debug "log_profile: delete a section"
  resource_delete "$_LOG_PROFILE_PATH/$_PROFILE_NAME" "force_reload=true"
  assert_equal "$SC" "204"
  resource_get "$_LOG_PROFILE_PATH/$_PROFILE_NAME"
  assert_equal "$SC" "404"
}
