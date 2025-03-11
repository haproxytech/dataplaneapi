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
load "../../libs/get_json_path"
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

# Doing all the tests in a single @tests runs faster.

@test "storage crt-list files: all tests" {
    # This is needed because docker cp cannot copy stuff
    # to /etc/haproxy/ssl during setup for some reason.
    dpa_docker_exec 'mv /ssl/* /etc/haproxy/ssl'

    debug "crt-list: upload a new list"
    assert dpa_docker_exec 'ls /etc/haproxy/ssl/crt-list.txt'
    refute dpa_docker_exec 'ls /etc/haproxy/ssl/crt-list2.txt'
    run dpa_curl_file_upload POST "$_STORAGE_SSL_CRT_LIST_PATH" "@$BATS_TEST_DIRNAME/data/crt-list2.txt;filename=crt-list2.txt"
    assert_success
    dpa_curl_status_body '$output'
    assert_equal "$SC" 202
    assert_equal "$(get_json_path "$BODY" '.storage_name')" 'crt-list2.txt'
    assert dpa_docker_exec 'ls /etc/haproxy/ssl/crt-list2.txt'

    debug "crt-list files: replace a crt-list"
    run dpa_curl_text_plain PUT "$_STORAGE_SSL_CRT_LIST_PATH/crt-list2.txt" "@$BATS_TEST_DIRNAME/data/crt-list2.txt"
    assert_success
    dpa_curl_status_body '$output'
    assert_equal "$SC" 202

    debug "crt-list files: get a list of crt-list files"
    resource_get "$_STORAGE_SSL_CRT_LIST_PATH"
    assert_equal "$SC" 200
    # [{"description":"managed certificate list","file":"/etc/haproxy/ssl/crt-list2.txt","storage_name":"crt-list2.txt"}]
    assert_equal "$(get_json_path "$BODY" '.[0].file')" '/etc/haproxy/ssl/crt-list.txt'
    assert_equal "$(get_json_path "$BODY" '.|length')" 2

    debug "crt-list files: get a single crt-list"
    resource_get "$_STORAGE_SSL_CRT_LIST_PATH/crt-list2.txt"
    test "$BODY" = "$(cat $BATS_TEST_DIRNAME/data/crt-list2.txt)" || fail

    debug "crt-list files: delete a list"
    resource_delete "$_STORAGE_SSL_CRT_LIST_PATH/crt-list2.txt"
    assert_equal "$SC" 202

    debug "crt-list files: ensure the file has been deleted"
    resource_get "$_STORAGE_SSL_CRT_LIST_PATH"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 1
}

@test "storage crt-list entries: all tests" {
    dpa_docker_exec 'mv /ssl/* /etc/haproxy/ssl'

    debug "crt-list entries: get the list of entries (1)"
    resource_get "$_STORAGE_SSL_CRT_LIST_PATH/crt-list.txt/entries"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')"  1

    debug "crt-list entries: add a new entry"
    resource_post "$_STORAGE_SSL_CRT_LIST_PATH/crt-list.txt/entries" "data/post.json"
    assert_equal "$SC" 201
    assert_equal "$(get_json_path "$BODY" '.file')" '/etc/haproxy/ssl/1.pem'

    debug "crt-list entries: get the list of entries (2)"
    resource_get "$_STORAGE_SSL_CRT_LIST_PATH/crt-list.txt/entries"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')"  2
    
    debug "crt-list entries: delete an entry"
    resource_delete "$_STORAGE_SSL_CRT_LIST_PATH/crt-list.txt/entries" "certificate=/etc/haproxy/ssl/1.pem&line_number=1"
    assert_equal "$SC" 204

    debug "crt-list entries: get the list of entries (3)"
    resource_get "$_STORAGE_SSL_CRT_LIST_PATH/crt-list.txt/entries"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')"  1
}
