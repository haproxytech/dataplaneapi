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


@test "runtime_maps_entries: Return one map runtime setting" {
    run dpa_curl GET "/services/haproxy/runtime/maps_entries/key1?map=mapfile1.map"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal "$(get_json_path "${BODY}" " select(.key | contains(\"key1\") ).key" )" "key1"
    assert_equal "$(get_json_path "${BODY}" " select(.value | contains(\"value1\") ).value" )" "value1"
}

@test "runtime_maps_entries: Return an error when one map runtime setting doesn't exists" {
    run dpa_curl GET "/services/haproxy/runtime/maps_entries/not-exists?map=mapfile1.map"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 404
}
