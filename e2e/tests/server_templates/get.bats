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

@test "server_templates: Return one server template" {
  resource_get "$_SERVER_TEMPLATE_BASE_PATH/srv_google" "backend=test_backend"
  assert_equal "$SC" 200

  assert_equal "srv_google" "$(get_json_path "$BODY" '.prefix')"
  assert_equal "0-9" "$(get_json_path "$BODY" '.num_or_range')"
  assert_equal "google.com" "$(get_json_path "$BODY" '.fqdn')"
  assert_equal "443" "$(get_json_path "$BODY" '.port')"

  resource_get "$_SERVER_TEMPLATE_BASE_PATH/srv_bing" "backend=test_backend"
  assert_equal "$SC" 200

  assert_equal "srv_bing" "$(get_json_path "$BODY" '.prefix')"
  assert_equal "10-19" "$(get_json_path "$BODY" '.num_or_range')"
  assert_equal "bing.com" "$(get_json_path "$BODY" '.fqdn')"
  assert_equal "8080" "$(get_json_path "$BODY" '.port')"
}

@test "server_templates: Return a non existing server template" {
  resource_get "$_SERVER_TEMPLATE_BASE_PATH/ghost" "backend=test_backend"
  assert_equal "$SC" 404
}
