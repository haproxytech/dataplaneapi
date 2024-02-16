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
load '../../libs/version'


load 'utils/_helpers'

@test "acl_runtime: Return ACL files list" {
    resource_get "$_RUNTIME_ACL_BASE_PATH"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "${BODY}" " .[0].storage_name" )" "path_beg"
    assert_equal "$(get_json_path "${BODY}" " .[1].storage_name" )" "path_end"
}

@test "acl_runtime: Return ACL file by its ID" {
    resource_get "$_RUNTIME_ACL_BASE_PATH/0"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" " .storage_name" )" "path_beg"

    resource_get "$_RUNTIME_ACL_BASE_PATH/1"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "${BODY}" " .storage_name" )" "path_end"
}
