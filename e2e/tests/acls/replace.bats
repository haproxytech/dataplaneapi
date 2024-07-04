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
load '../../libs/version'
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'

load 'utils/_helpers'

@test "acls: Replace one ACL" {
    PARENT_NAME="fe_acl"
    resource_put "$_FRONTEND_BASE_PATH/$PARENT_NAME/acls/2" "data/put.json"
    assert_equal "$SC" 202
    #
    # verify that ACL is actually changed
    #
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/acls/2"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "$BODY" " ." )" "$(cat "$BATS_TEST_DIRNAME/data/put.json")"
}

@test "acls: Return an error when trying to replace non existing ACL" {
    PARENT_NAME="fe_acl"
    resource_put "$_FRONTEND_BASE_PATH/$PARENT_NAME/200" "data/put.json"
    assert_equal "$SC" 404
}

@test "acls: Replace all ACLs" {
    PARENT_NAME="fe_acl"
    resource_put "$_FRONTEND_BASE_PATH/$PARENT_NAME/acls" "data/replace-all.json"
    assert_equal "$SC" 202
    #
    # verify that ACL is actually changed
    #
    #
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/acls"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" ". | length")" 3
    assert_equal "$(get_json_path "$BODY" ".")" "$(get_json_path "$(cat "$BATS_TEST_DIRNAME/data/replace-all.json")" ".")"

}
