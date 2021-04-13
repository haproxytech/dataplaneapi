#!/usr/bin/env bats
#
# Copyright 2019 HAProxy Technologies
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

@test "acls: Return ACL list" {
    run dpa_curl GET "/services/haproxy/configuration/acls?parent_name=fe_acl&parent_type=frontend"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal "$(get_json_path "${BODY}" " .data | .[2].acl_name" )" "local_dst"
    assert_equal "$(get_json_path "${BODY}" " .data | .[2].criterion" )" "hdr(host)"
    assert_equal "$(get_json_path "${BODY}" " .data | .[2].index" )" "2"
    assert_equal "$(get_json_path "${BODY}" " .data | .[2].value" )" "-i localhost"
}

@test "acls: Return ACL list by its name" {
    run dpa_curl GET "/services/haproxy/configuration/acls?parent_name=fe_acl&parent_type=frontend&acl_name=invalid_src"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal "$(get_json_path "${BODY}" " .data | .[1].acl_name" )" "invalid_src"
    assert_equal "$(get_json_path "${BODY}" " .data | .[1].criterion" )" "src_port"
    assert_equal "$(get_json_path "${BODY}" " .data | .[1].index" )" "1"
    assert_equal "$(get_json_path "${BODY}" " .data | .[1].value" )" "0:1023"
}
