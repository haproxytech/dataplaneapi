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
load '../../libs/get_json_path'
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_version'

load 'utils/_helpers'


@test "backends: Add a backend (>=2.2)" {
  if haproxy_version_ge "2.2"
  then
  resource_post "$_BACKEND_BASE_PATH" "data/post_2.2.json" "force_reload=true"
  assert_equal "$SC" "201"

  resource_get "$_BACKEND_BASE_PATH/test_backend"  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".\"http-check\".method")" "OPTIONS"
  assert_equal "$(get_json_path "$BODY" ".\"http-check\".uri")" "/"
  assert_equal "$(get_json_path "$BODY" ".\"http-check\".version")" "HTTP/1.1"
  assert_equal "$(get_json_path "$BODY" ".adv_check")" "httpchk"
  assert_equal "$(get_json_path "$BODY" ".httpchk_params.method")" "GET"
  assert_equal "$(get_json_path "$BODY" ".httpchk_params.uri")" "/check"
  assert_equal "$(get_json_path "$BODY" ".httpchk_params.version")" "HTTP/1.1"
  fi
}

@test "backends: Replace a backend (>=2.2)" {
  if haproxy_version_ge "2.2"
  then
	resource_put "$_BACKEND_BASE_PATH/test_backend" "data/put.json" "force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_BACKEND_BASE_PATH/test_backend"  assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ".\"http-check\".method")" "OPTIONS"
    assert_equal "$(get_json_path "$BODY" ".\"http-check\".uri")" "/"
    assert_equal "$(get_json_path "$BODY" ".\"http-check\".version")" "HTTP/1.1"
    assert_equal "$(get_json_path "$BODY" ".adv_check")" "httpchk"
    assert_equal "$(get_json_path "$BODY" ".httpchk_params.method")" "GET"
    assert_equal "$(get_json_path "$BODY" ".httpchk_params.uri")" "/healthz"
    assert_equal "$(get_json_path "$BODY" ".httpchk_params.version")" "HTTP/1.1"
  fi
}

@test "backends: Delete a backend" {
  if haproxy_version_ge "2.2"
  then
	resource_delete "$_BACKEND_BASE_PATH/test_backend" "force_reload=true"
	assert_equal "$SC" 204
  fi
}
