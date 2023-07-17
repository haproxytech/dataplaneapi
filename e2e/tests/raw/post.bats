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
load '../../libs/debug'
load "../../libs/get_json_path"
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "raw: Post new configuration" {
  resource_post "$_RAW_BASE_PATH" 'data/haproxy.cfg.json'
  assert_equal "$SC" 202
}

@test "raw: Post new configuration incorrectly (bug 1219)" {
  resource_post_text "$_RAW_BASE_PATH" 'data/haproxy.cfg.json'
  assert_equal "$SC" 400
  assert_equal "$(get_json_path "$BODY" '.message')" "invalid configuration: no newline character found"
}


@test "raw: Post new configuration with socket path changed" {
    resource_post "$_RAW_BASE_PATH" 'data/haproxy_socket.cfg.json'
	assert_equal "$SC" 202

	resource_get "$_RAW_BASE_PATH" ""
	assert_equal "$SC" 200

    local socket; socket='stats socket /var/lib/haproxy/stats-new'
    if [[ "$BODY" != *"$socket"* ]]; then
       batslib_print_kv_single_or_multi 8 \
          'configuration' "$BODY" \
          'expected socket'   "$socket" \
          | batslib_decorate 'configuration does not contains the new socket' \
          | fail
    fi

	# check that runtime client has been reconfigured with the new socket
	sleep 5
	resource_get "$_RUNTIME_MAP_FILES_BASE_PATH" ""
    assert_equal "$SC" 200
}