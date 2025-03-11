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
load "../../libs/get_json_path"
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "ssl_crt_lists_runtime: Get all crt lists" {
    resource_post "$_RUNTIME_SSL_CRT_LISTS_BASE_PATH/entries" "data/post.json" "name=/ssl/crt-list.txt"
    assert_equal "$SC" 201
    resource_get "$_RUNTIME_SSL_CRT_LISTS_BASE_PATH"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 1
    assert_equal "$(get_json_path "$BODY" '.[0].file')" "/ssl/crt-list.txt"
}

@test "ssl_crt_lists_runtime: Get a crt list entry" {
    resource_post "$_RUNTIME_SSL_CRT_LISTS_BASE_PATH/entries" "data/post.json" "name=/ssl/crt-list.txt"
    dpa_curl_status_body '$output'
    assert_equal "$SC" 201
    resource_get "$_RUNTIME_SSL_CRT_LISTS_BASE_PATH/entries" "name=/ssl/crt-list.txt"
    dpa_curl_status_body '$output'
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 2
    assert_equal "$(get_json_path "$BODY" '.[0].file')" "/ssl/1.pem"
    assert_equal "$(get_json_path "$BODY" '.[0].line_number')" "1"
    assert_equal "$(get_json_path "$BODY" '.[1].file')" "/ssl/1.pem"
    assert_equal "$(get_json_path "$BODY" '.[1].line_number')" "2"
    assert_equal "$(get_json_path "$BODY" '.[1].ssl_bind_config')" "alpn h2"
    assert_equal "$(get_json_path "$BODY" '.[1].sni_filter[0]')" "mysite.local"
}

@test "ssl_crt_lists_runtime: Delete a crt list entry" {
    resource_post "$_RUNTIME_SSL_CRT_LISTS_BASE_PATH/entries" "data/post.json" "name=/ssl/crt-list.txt"
    dpa_curl_status_body '$output'
    assert_equal "$SC" 201
    resource_delete "$_RUNTIME_SSL_CRT_LISTS_BASE_PATH/entries" "name=/ssl/crt-list.txt&cert_file=/ssl/1.pem&line_number=2"
    assert_equal "$SC" 204
}
