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
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "binds: Return an array of binds" {
  PARENT_NAME="test_frontend"
  resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/binds"
  assert_equal "$SC" 200

  if haproxy_version_ge "2.5"
  then
    assert_equal 3 "$(get_json_path "$BODY" ". | length")"
  else
    assert_equal 2 "$(get_json_path "$BODY" ". | length")"
  fi
  assert_equal "fixture" "$(get_json_path "$BODY" ".[0].name")"
  assert_equal "loopback" "$(get_json_path "$BODY" ".[1].name")"
  if haproxy_version_ge "2.5"
  then
    assert_equal "all" "$(get_json_path "$BODY" ".[0].thread")"
    assert_equal "1/1" "$(get_json_path "$BODY" ".[1].thread")"
    assert_equal "loopback1" "$(get_json_path "$BODY" ".[2].name")"
    assert_equal "1/1-1" "$(get_json_path "$BODY" ".[2].thread")"
  fi
}
