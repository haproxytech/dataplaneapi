#!/usr/bin/env bats
#
# Copyright 2024 HAProxy Technologies
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

# We are using "haproxy_config_setup" here because we need the files
# from data/container. HAProxy will actually load those PEM files when
# checking if the configuration is valid.

@test "crt_store: all tests (>=3.0)" {
  haproxy_version_ge "3.0" || skip

  debug "crt_store: create a new section"
  resource_post "$_CRT_STORE_PATH" "data/new_store.json" "force_reload=true"
  assert_equal "$SC" "201"

  debug "crt_store: get a section"
  resource_get "$_CRT_STORE_PATH/$_STORE_NAME"
  assert_equal "$SC" "200"
  assert_equal "$_STORE_NAME" "$(get_json_path "$BODY" .name)"
  assert_equal "/secure/certs" "$(get_json_path "$BODY" .crt_base)"
  assert_equal "/secure/keys" "$(get_json_path "$BODY" .key_base)"

  debug "crt_store: edit a section"
  resource_put "$_CRT_STORE_PATH/$_STORE_NAME" "data/edit_store.json" "force_reload=true"
  assert_equal "$SC" "200"
  resource_get "$_CRT_STORE_PATH/$_STORE_NAME"
  assert_equal "/sec" "$(get_json_path "$BODY" .key_base)"

  debug "crt_store: get a list of sections"
  resource_get "$_CRT_STORE_PATH"
  assert_equal "$SC" "200"
  assert_equal "$_STORE_NAME" "$(get_json_path "$BODY" .[0].name)"

  debug "crt_store: add load entries"
  resource_post "$_CRT_LOAD_PATH" "data/post_entry1.json" "crt_store=$_STORE_NAME"
  assert_equal "$SC" "202"
  resource_post "$_CRT_LOAD_PATH" "data/post_entry2.json" "crt_store=$_STORE_NAME"
  assert_equal "$SC" "202"

  debug "crt_store: get a load entry"
  resource_get "$_CRT_LOAD_PATH/c1.pem" "crt_store=$_STORE_NAME"
  assert_equal "$SC" "200"
  assert_equal "c1.pem" "$(get_json_path "$BODY" .certificate)"
  assert_equal "k1.pem" "$(get_json_path "$BODY" .key)"
  assert_equal "disabled" "$(get_json_path "$BODY" .ocsp_update)"

  debug "crt_store: get all load entries"
  resource_get "$_CRT_LOAD_PATH" "crt_store=$_STORE_NAME"
  assert_equal "$SC" "200"
  assert_equal "2" "$(get_json_path "$BODY" '.|length')"
  assert_equal "c1.pem" "$(get_json_path "$BODY" .[0].certificate)"
  assert_equal "c2.pem" "$(get_json_path "$BODY" .[1].certificate)"

  debug "crt_store: modify a load entry"
  resource_put "$_CRT_LOAD_PATH/c2.pem" "data/put_entry.json" \
    "crt_store=$_STORE_NAME" "force_reload=true"
  assert_equal "$SC" "202"
  resource_get "$_CRT_LOAD_PATH/c2.pem" "crt_store=$_STORE_NAME"
  assert_equal "c2.pem" "$(get_json_path "$BODY" .certificate)"
  assert_equal "disabled" "$(get_json_path "$BODY" .ocsp_update)"
  assert_equal "example.com" "$(get_json_path "$BODY" .alias)"

  debug "crt_store: delete a load entry"
  resource_delete "$_CRT_LOAD_PATH/c1.pem" "crt_store=$_STORE_NAME" "force_reload=true"
  assert_equal "$SC" "202"
  resource_delete "$_CRT_LOAD_PATH/c2.pem" "crt_store=$_STORE_NAME" "force_reload=true"
  assert_equal "$SC" "202"
  resource_get "$_CRT_LOAD_PATH/c2.pem" "crt_store=$_STORE_NAME"
  assert_equal "$SC" "404"

  debug "crt_store: delete a section"
  resource_delete "$_CRT_STORE_PATH/$_STORE_NAME" "force_reload=true"
  assert_equal "$SC" "204"
  resource_get "$_CRT_STORE_PATH/$_STORE_NAME"
  assert_equal "$SC" "404"
}
