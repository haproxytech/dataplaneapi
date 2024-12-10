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
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'


@test "traces: get/create/modify/delete (>=3.1)" {
  haproxy_version_ge "3.1" || skip

  debug "traces: create section"
  resource_post "$_TRACES_PATH" "data/new_section.json" "force_reload=true"
  assert_equal "$SC" "201"

  debug "traces: get section"
  resource_get "$_TRACES_PATH"
  assert_equal "$SC" "200"
  assert_equal "h1 sink stderr level developer verbosity complete start now" \
    "$(get_json_path "$BODY" .entries[0].trace)"
  assert_equal "h2 sink stderr level developer verbosity complete start now" \
    "$(get_json_path "$BODY" .entries[1].trace)"

  debug "traces: replace section"
  resource_put "$_TRACES_PATH" "data/replace_section.json" "force_reload=true"
  assert_equal "$SC" "200"
  resource_get "$_TRACES_PATH"
  assert_equal "fcgi sink stderr level developer verbosity quiet start now" \
    "$(get_json_path "$BODY" .entries[0].trace)"


  debug "traces: add trace entries"
  resource_post "$_TRACE_ENTRIES_PATH" "data/entry1.json"
  assert_equal "$SC" "202"
  resource_post "$_TRACE_ENTRIES_PATH" "data/entry2.json"
  assert_equal "$SC" "202"
  resource_get "$_TRACES_PATH"
  assert_equal "fcgi sink stderr level developer verbosity quiet start now" \
    "$(get_json_path "$BODY" .entries[0].trace)"
  assert_equal "peers sink stderr level developer verbosity quiet start now" \
    "$(get_json_path "$BODY" .entries[1].trace)"
  assert_equal "check sink stderr level developer verbosity quiet start now" \
    "$(get_json_path "$BODY" .entries[2].trace)"

  debug "traces: delete entries"
  resource_delete_body "$_TRACE_ENTRIES_PATH" "data/entry1.json" "force_reload=true"
  assert_equal "$SC" "204"
  resource_delete_body "$_TRACE_ENTRIES_PATH" "data/entry2.json" "force_reload=true"
  assert_equal "$SC" "204"

  debug "traces: delete section"
  resource_delete "$_TRACES_PATH" "force_reload=true"
  assert_equal "$SC" "204"
  resource_get "$_TRACES_PATH"
  assert_equal "$SC" "404"
}
