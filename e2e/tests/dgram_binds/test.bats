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

@test "dgram_binds: Add a dgram bind" {
  if haproxy_version_ge "2.3"
  then
    PARENT_NAME="sylog-loadb"
    resource_post "$_LOGFORWARD_BASE_PATH/$PARENT_NAME/dgram_binds" "data/post.json" "force_reload=true"
    assert_equal "$SC" "201"
  fi
}

@test "dgram_binds: Return a dgram bind" {
  if haproxy_version_ge "2.3"
  then
    PARENT_NAME="sylog-loadb"
    resource_get "$_LOGFORWARD_BASE_PATH/$PARENT_NAME/dgram_binds/test_dgram_bind"
    assert_equal "$SC" 200
    assert_equal "test_dgram_bind" "$(get_json_path "$BODY" '.name')"
  fi
}

@test "dgram_binds: Replace a dgram bind" {
  if haproxy_version_ge "2.3"
  then
    PARENT_NAME="sylog-loadb"
    resource_put "$_LOGFORWARD_BASE_PATH/$PARENT_NAME/dgram_binds/test_dgram_bind" "data/put.json" "log_forward=sylog-loadb&force_reload=true"
    assert_equal "$SC" 200
  fi
}

@test "dgram_binds: Return an array of dgram binds" {
  if haproxy_version_ge "2.3"
  then
    PARENT_NAME="sylog-loadb"
    resource_get "$_LOGFORWARD_BASE_PATH/$PARENT_NAME/dgram_binds"
    assert_equal "$SC" 200
    assert_equal "test_dgram_bind" "$(get_json_path "$BODY" '.[0].name')"
  fi
}

@test "dgram_binds: Delete a dgram bind" {
  if haproxy_version_ge "2.3"
  then
    PARENT_NAME="sylog-loadb"
    resource_delete "$_LOGFORWARD_BASE_PATH/$PARENT_NAME/dgram_binds/test_dgram_bind" "force_reload=true"
    assert_equal "$SC" 204
  fi
}
