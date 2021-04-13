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
load '../../libs/version'
load '../../libs/haproxy_config_setup'

@test "acls: Replace one ACL" {
    run dpa_curl PUT "/services/haproxy/configuration/acls/2?parent_name=fe_acl&parent_type=frontend&version=$(version)" "/data/put.json"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 202

    # verify that ACL is actually changed
    run dpa_curl GET "/services/haproxy/configuration/acls/2?parent_name=fe_acl&parent_type=frontend"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal "$(get_json_path "${BODY}" " .data" )" "$(cat ${BATS_TEST_DIRNAME}/data/put.json)"
}

@test "acls: Return an error when trying to replace non existing ACL" {
    run dpa_curl PUT "/services/haproxy/configuration/acls/200?parent_name=fe_acl&parent_type=frontend&version=$(version)" "/data/put.json"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 404
}
