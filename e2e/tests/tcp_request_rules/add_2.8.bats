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
load '../../libs/haproxy_version'

load 'utils/_helpers'


## tune.stick-counters is a haproxy >= 2.8 global parameter
@test "tcp_request_rules: Add a new TCP Request Rule track-sc to backend" {
  if haproxy_version_ge "2.8"
  then
  # Using new track-sc with track_sc_stick_counter
  PARENT_NAME="test_sticksc"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_request_rules/0" "data/post-track-sc.json" "force_reload=true"
	assert_equal "$SC" 201

  resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_request_rules/0"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".action")" "track-sc"
	assert_equal "$(get_json_path "$BODY" ".type")" "content"
	assert_equal "$(get_json_path "$BODY" ".cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".cond_test")" "TRUE"
	assert_equal "$(get_json_path "$BODY" ".track_key")" "src"
	assert_equal "$(get_json_path "$BODY" ".track_table")" "test_sticksc"
	assert_equal "$(get_json_path "$BODY" ".track_stick_counter")" 5
  fi
}

@test "tcp_request_rules: Fail - Add a new TCP Request Rule track-sc to backend - when track_sc_stick_counter is missing" {
  if haproxy_version_ge "2.8"
  then
  # Using new track-sc with track_sc_stick_counter
  # Fail due to missing track_sc_stick_counter
  PARENT_NAME="test_sticksc"
  resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/tcp_request_rules/0" "data/post-track-sc-fail.json" "force_reload=true"
	assert_equal "$SC" 400
  fi
}
