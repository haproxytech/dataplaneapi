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

# NOTE: "Replace" and "Delete" currently fail against client-native v6: its
# transaction round-trip (loadDataForChange) drops force-be-switch rules, so
# edit/delete cannot find a rule that GET/list return. These cases are kept to
# document the intended behavior and will pass once client-native is fixed.

@test "force_be_switches: Add a force-be-switch rule" {
  if haproxy_version_ge "3.4"
  then
	resource_post "$_FRONTEND_BASE_PATH/test_frontend/force_be_switches" "data/post.json" "force_reload=true"
	assert_equal "$SC" 201
  fi
}

@test "force_be_switches: Return a force-be-switch rule" {
  if haproxy_version_ge "3.4"
  then
	resource_get "$_FRONTEND_BASE_PATH/test_frontend/force_be_switches/0"
	assert_equal "$SC" 200
	assert_equal "if" "$(get_json_path "$BODY" '.cond')"
	assert_equal "{ src 10.0.0.0/8 }" "$(get_json_path "$BODY" '.cond_test')"
  fi
}

@test "force_be_switches: Return an array of force-be-switch rules" {
  if haproxy_version_ge "3.4"
  then
	resource_get "$_FRONTEND_BASE_PATH/test_frontend/force_be_switches"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ". | length")" 1
  fi
}

@test "force_be_switches: Replace a force-be-switch rule" {
  if haproxy_version_ge "3.4"
  then
	resource_put "$_FRONTEND_BASE_PATH/test_frontend/force_be_switches/0" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200
	assert_equal "unless" "$(get_json_path "$BODY" '.cond')"
  fi
}

@test "force_be_switches: Delete a force-be-switch rule" {
  if haproxy_version_ge "3.4"
  then
	resource_delete "$_FRONTEND_BASE_PATH/test_frontend/force_be_switches/0" "force_reload=true"
	assert_equal "$SC" 204
  fi
}
