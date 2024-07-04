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
load "../../libs/get_json_path"
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/haproxy_version'
load '../../libs/version'

load 'utils/_helpers'

@test "acl_runtime: Return ACL file entries list" {
  PARENT_NAME="0"
  resource_get "$_RUNTIME_ACL_BASE_PATH/$PARENT_NAME/entries"
  assert_equal "$SC" 200

  assert_equal "$(get_json_path "${BODY}" " .[0].value" )" "/static"
  assert_equal "$(get_json_path "${BODY}" " .[1].value" )" "/images"
  assert_equal "$(get_json_path "${BODY}" " .[2].value" )" "/javascript"
  assert_equal "$(get_json_path "${BODY}" " .[3].value" )" "/stylesheets"
}

@test "acl_runtime: Return ACL file entries by their ID" {
  PARENT_NAME="0"
  resource_get "$_RUNTIME_ACL_BASE_PATH/$PARENT_NAME/entries"
  assert_equal "$SC" 200

  local list; list=$BODY

  for index in {0..3}; do
    local id; id="$(get_json_path "${list}" ".[${index}].id" )"
    local value; value="$(get_json_path "${list}" ".[${index}].value" )"

    resource_get "$_RUNTIME_ACL_BASE_PATH/$PARENT_NAME/entries/$id"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ".value" )" "$value"
  done
}

@test "acl_runtime: Add an ACL file entry" {
  PARENT_NAME="0"
  resource_post "$_RUNTIME_ACL_BASE_PATH/$PARENT_NAME/entries" "data/post.json"
  assert_equal "$SC" 201
  assert_equal "$(get_json_path "${BODY}" " .value" )" "/js"
}

@test "acl_runtime: Delete an ACL file entry by its ID" {
  PARENT_NAME="1"
  if haproxy_version_ge "2.9"; then
    skip "cause: bug in HAPRoxy 2.9"
  fi
	# checking items and retrieving first ID
	resource_get "$_RUNTIME_ACL_BASE_PATH/$PARENT_NAME/entries"
  assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ". | length" )" "5"

	local id; id="$(get_json_path "$BODY" ".[0].id")"

	# deleting the entry file by its ID
	resource_delete "$_RUNTIME_ACL_BASE_PATH/$PARENT_NAME/entries/$id"
  assert_equal "$SC" 204

	# checking the file entry has been deleted counting the items
	resource_get "$_RUNTIME_ACL_BASE_PATH/$PARENT_NAME/entries"
  assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ". | length" )" "4"
}
