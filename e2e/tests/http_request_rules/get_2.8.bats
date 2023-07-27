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

@test "http_request_rules: Return one track-sc HTTP Request Rule from frontend" {
  if haproxy_version_ge "2.8"
  then
  resource_get "$_REQ_RULES_BASE_PATH/0" "parent_type=frontend&parent_name=test_sticksc"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".data.type")" "track-sc"
	assert_equal "$(get_json_path "$BODY" ".data.cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".data.cond_test")" "TRUE"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_key")" "src"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_table")" "test_sticksc"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_stick_counter")" 0

  resource_get "$_REQ_RULES_BASE_PATH/1" "parent_type=frontend&parent_name=test_sticksc"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".data.type")" "track-sc"
	assert_equal "$(get_json_path "$BODY" ".data.cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".data.cond_test")" "TRUE"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_key")" "src"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_table")" "test_sticksc"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_stick_counter")" 5
  fi
}

@test "http_request_rules: Return one track-sc HTTP Request Rule from backend" {
  if haproxy_version_ge "2.8"
  then
  resource_get "$_REQ_RULES_BASE_PATH/0" "parent_type=backend&parent_name=test_sticksc"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".data.type")" "track-sc"
	assert_equal "$(get_json_path "$BODY" ".data.cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".data.cond_test")" "TRUE"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_key")" "src"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_table")" "test_sticksc"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_stick_counter")" 0

  resource_get "$_REQ_RULES_BASE_PATH/1" "parent_type=backend&parent_name=test_sticksc"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".data.type")" "track-sc"
	assert_equal "$(get_json_path "$BODY" ".data.cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".data.cond_test")" "TRUE"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_key")" "src"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_table")" "test_sticksc"
	assert_equal "$(get_json_path "$BODY" ".data.track_sc_stick_counter")" 5
  fi
}
