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

load 'utils/_helpers'

@test "acls: Return ACL list" {
    resource_get "$_ACL_BASE_PATH" "parent_name=fe_acl&parent_type=frontend"
    assert_equal $SC "200"

    assert_equal "$(get_json_path "$BODY" " . | .[2].acl_name" )" "local_dst"
    assert_equal "$(get_json_path "$BODY" " . | .[2].criterion" )" "hdr(host)"
    assert_equal "$(get_json_path "$BODY" " . | .[2].value" )" "-i localhost"
}

@test "acls: Return ACL list by its name" {
    resource_get "$_ACL_BASE_PATH" "parent_name=fe_acl&parent_type=frontend&acl_name=invalid_src"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "$BODY" " . | .[1].acl_name" )" "invalid_src"
    assert_equal "$(get_json_path "$BODY" " . | .[1].criterion" )" "src_port"
    assert_equal "$(get_json_path "$BODY" " . | .[1].value" )" "0:1023"
}
