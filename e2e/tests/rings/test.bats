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
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "rings: Add a ring" {
  if haproxy_version_ge "2.2"
  then
	resource_post "$_RING_BASE_PATH" "data/post.json" "force_reload=true"
	assert_equal "$SC" "201"

	resource_post "$_RING_BASE_PATH" "data/post-ring2.json" "force_reload=true"
    assert_equal "$SC" "201"
  fi
}

@test "rings: Return a ring" {
  if haproxy_version_ge "2.2"
  then
	resource_get "$_RING_BASE_PATH/test_ring"
	assert_equal "$SC" 200
	assert_equal "test_ring" "$(get_json_path "$BODY" '.name')"
  fi
}

@test "rings: Replace a ring" {
  if haproxy_version_ge "2.2"
  then
	resource_put "$_RING_BASE_PATH/test_ring" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200
  fi
}

@test "rings: Return an array of rings" {
  if haproxy_version_ge "2.2"
  then
	resource_get "$_RING_BASE_PATH"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ". | length")" 2
	assert_equal "$(get_json_path "$BODY" ".[] | select(.name | contains(\"test_ring_2\") ).name")" "test_ring_2"
  fi
}

@test "rings: Delete a ring" {
  if haproxy_version_ge "2.2"
  then
	resource_delete "$_RING_BASE_PATH/test_ring" "force_reload=true"
	assert_equal "$SC" 204
  fi
}
