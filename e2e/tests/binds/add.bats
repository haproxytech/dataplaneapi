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

@test "binds: Add a new bind" {
  PARENT_NAME="test_frontend"
  if haproxy_version_ge "2.5"
  then
    resource_post "$_FRONTEND_BASE_PATH/$PARENT_NAME/binds" "data/post_2.5.json" "force_reload=true"
  else
    resource_post "$_FRONTEND_BASE_PATH/$PARENT_NAME/binds" "data/post.json" "force_reload=true"
  fi
  assert_equal "$SC" 201
  if haproxy_version_ge "2.5"
  then
    resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/binds/test_bind" "force_reload=true"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.name')" "test_bind"
    assert_equal "1/all" "$(get_json_path "$BODY" ".thread")"
    assert_equal "$(get_json_path "$BODY" '.ca_verify_file')" "/certs/ca-verify.pem"
  fi
}

@test "binds: Add a new bind for log_forward" {
  PARENT_NAME="sylog-loadb"
  if haproxy_version_ge "2.5"
  then
    resource_post "$_LOGFORWARD_BASE_PATH/$PARENT_NAME/binds" "data/post_2.5.json" "force_reload=true"
    assert_equal "$SC" 201

    resource_get "$_LOGFORWARD_BASE_PATH/$PARENT_NAME/binds/test_bind" "force_reload=true"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.name')" "test_bind"
    assert_equal "1/all" "$(get_json_path "$BODY" ".thread")"
    assert_equal "$(get_json_path "$BODY" '.ca_verify_file')" "/certs/ca-verify.pem"
  fi
}

@test "binds: Add a new bind for peers" {
  PARENT_NAME="fusion"
  # segfault in v3.3
  if haproxy_version_ge "2.5" && !haproxy_version_ge "3.3"
  then
    resource_post "$_PEER_BASE_PATH/$PARENT_NAME/binds" "data/post_2.5.json" "force_reload=true"
    assert_equal "$SC" 201

    resource_get "$_PEER_BASE_PATH/$PARENT_NAME/binds/test_bind" "force_reload=true"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.name')" "test_bind"
    assert_equal "1/all" "$(get_json_path "$BODY" ".thread")"
    assert_equal "$(get_json_path "$BODY" '.ca_verify_file')" "/certs/ca-verify.pem"
  fi
}
