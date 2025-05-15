#!/usr/bin/env bats
#
# Copyright 2025 HAProxy Technologies
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
load '../../libs/debug'
load '../../libs/get_json_path'
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

cert=/etc/haproxy/ssl/cert.pem

@test "ssl_front_uses: get/create/modify/delete (>=3.2)" {
  haproxy_version_ge "3.2" || skip

  run docker cp "${BATS_TEST_DIRNAME}/data/3.pem" "${DOCKER_CONTAINER_NAME}:$cert"

  resource_post "$(ssl_front_uses_path front1)" "data/post.json"
  assert_equal "$SC" "202"

  resource_get "$(ssl_front_uses_path front1)"
  assert_equal "$SC" "200"
  assert_equal "$(get_json_path "$BODY" '.|length')" 1
  assert_equal "$(get_json_path "$BODY" '.[0].certificate')" "$cert"
  assert_equal "$(get_json_path "$BODY" '.[0].allow_0rtt')" true

  resource_get "$(ssl_front_uses_path front1 0)"
  assert_equal "$SC" "200"
  assert_equal "$(get_json_path "$BODY" '.certificate')" "$cert"
  assert_equal "$(get_json_path "$BODY" '.allow_0rtt')" true

  resource_put "$(ssl_front_uses_path front1 0)" "data/put.json"
  assert_equal "$SC" "202"
  assert_equal "$(get_json_path "$BODY" '.certificate')" "$cert"
  assert_equal "$(get_json_path "$BODY" '.allow_0rtt')" null
  assert_equal "$(get_json_path "$BODY" '.no_alpn')" true
  assert_equal "$(get_json_path "$BODY" '.verify')" none

  resource_delete "$(ssl_front_uses_path front1 0)"
  assert_equal "$SC" "202"
  resource_get "$(ssl_front_uses_path front1)"
  assert_equal "$SC" "200"
  assert_equal "$(get_json_path "$BODY" '.|length')" 0

  run docker exec "${DOCKER_CONTAINER_NAME}" rm $cert
}
