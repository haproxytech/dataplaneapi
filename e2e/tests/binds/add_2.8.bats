#!/usr/bin/env bats
#
# Copyright 2023 HAProxy Technologies
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
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "binds: Add a new bind (>=2.8)" {
  if haproxy_version_ge "2.8"
  then
    resource_post "$_BIND_BASE_PATH" "data/post_2.8.json" "frontend=test_frontend&force_reload=true"
    assert_equal "$SC" 201

    resource_get "$_BIND_BASE_PATH/test_bind" "frontend=test_frontend&force_reload=true"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.data.name')" "test_bind"
    assert_equal "$(get_json_path "$BODY" ".data.no_alpn")" "true"
    assert_equal "$(get_json_path "$BODY" '.data.ca_verify_file')" "/certs/ca-verify.pem"
  fi
}
