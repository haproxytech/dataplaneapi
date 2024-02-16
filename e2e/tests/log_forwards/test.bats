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

@test "log_forwards: Add a log forward" {
  if haproxy_version_ge "2.3"
  then
    resource_post "$_LOG_FORWARD_BASE_PATH" "data/post.json" "force_reload=true"
    assert_equal "$SC" "201"
  fi
}

@test "log_forwards: Return a log forward" {
  if haproxy_version_ge "2.3"
  then
    resource_get "$_LOG_FORWARD_BASE_PATH/sylog-loadb"
    assert_equal "$SC" 200
  	assert_equal "sylog-loadb" "$(get_json_path "$BODY" '.name')"
  fi
}

@test "log_forwards: Replace a log forward" {
  if haproxy_version_ge "2.3"
  then
    resource_put "$_LOG_FORWARD_BASE_PATH/sylog-loadb" "data/put.json" "force_reload=true"
    assert_equal "$SC" 200
  fi
}

@test "log_forwards: Return an array of log_forwards" {
  if haproxy_version_ge "2.3"
  then
    resource_get "$_LOG_FORWARD_BASE_PATH"
    assert_equal "$SC" 200
    assert_equal "sylog-loadb" "$(get_json_path "$BODY" '.[0].name')"
  fi
}

@test "log_forwards: Delete a log forward" {
  if haproxy_version_ge "2.3"
  then
    resource_delete "$_LOG_FORWARD_BASE_PATH/sylog-loadb" "force_reload=true"
    assert_equal "$SC" 204
  fi
}
