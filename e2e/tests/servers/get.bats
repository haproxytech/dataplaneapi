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

@test "servers: Return one server" {
  for name in "server_01" "server_02" "server_03"; do
    resource_get "$_SERVER_BASE_PATH/$name" "backend=test_backend"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.name')" "$name"
  done
}


@test "servers: Return one IPv6 server" {
    resource_get "$_SERVER_BASE_PATH/server_ipv6" "backend=test_backend"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.address')" "fd00:6:48:c85:deb:f:62:4"
    assert_equal "$(get_json_path "$BODY" '.check')" "enabled"
    assert_equal "$(get_json_path "$BODY" '.name')" "server_ipv6"
    assert_equal "$(get_json_path "$BODY" '.port')" "80"
}
