#!/usr/bin/env bats
#
# Copyright 2022 HAProxy Technologies
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

# We test both mailers sections and entries in the same file.
# If they were separated in 2 files, bats would run them
# in parrallel, making some tests fail.

load '../../libs/dataplaneapi'
load '../../libs/get_json_path'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "mailers: add a section" {
  resource_post "$_MAILERS_SECTION_PATH" "data/post_section.json" "force_reload=true"
  assert_equal "$SC" "201"
}

@test "mailers: get a section" {
  resource_get "$_MAILERS_SECTION_PATH/$_SECTION_NAME"
  assert_equal "$SC" "200"
  assert_equal "$_SECTION_NAME" "$(get_json_path "$BODY" .name)"
  assert_equal "15000" "$(get_json_path "$BODY" .timeout)"
}

@test "mailers: edit a section" {
  resource_put "$_MAILERS_SECTION_PATH/$_SECTION_NAME" "data/put_section.json" "force_reload=true"
  assert_equal "$SC" "200"
  resource_get "$_MAILERS_SECTION_PATH/$_SECTION_NAME"
  assert_equal "30000" "$(get_json_path "$BODY" .timeout)"
}

@test "mailers: get a list of sections" {
  resource_get "$_MAILERS_SECTION_PATH"
  assert_equal "$SC" "200"
  assert_equal "$_SECTION_NAME" "$(get_json_path "$BODY" .[0].name)"
}

@test "mailers: add entries" {
  resource_post "$_MAILER_ENTRIES_PATH" "data/post_entry1.json" "mailers_section=$_SECTION_NAME"
  assert_equal "$SC" "202"
  resource_post "$_MAILER_ENTRIES_PATH" "data/post_entry2.json" "mailers_section=$_SECTION_NAME"
  assert_equal "$SC" "202"
}

@test "mailers: get an entry" {
  resource_get "$_MAILER_ENTRIES_PATH/smtp1" "mailers_section=$_SECTION_NAME"
  assert_equal "$SC" "200"
  assert_equal "smtp1" "$(get_json_path "$BODY" .name)"
  assert_equal "10.0.10.1" "$(get_json_path "$BODY" .address)"
  assert_equal "587" "$(get_json_path "$BODY" .port)"
}

@test "mailers: get all entries" {
  resource_get "$_MAILER_ENTRIES_PATH" "mailers_section=$_SECTION_NAME"
  assert_equal "$SC" "200"
  assert_equal "2" "$(get_json_path "$BODY" '.|length')"
  assert_equal "smtp1" "$(get_json_path "$BODY" .[0].name)"
  assert_equal "smtp2" "$(get_json_path "$BODY" .[1].name)"
}

@test "mailers: modify an entry" {
  resource_put "$_MAILER_ENTRIES_PATH/smtp2" "data/put_entry.json" \
    "mailers_section=$_SECTION_NAME" "force_reload=true"
  assert_equal "$SC" "202"
  resource_get "$_MAILER_ENTRIES_PATH/smtp2" "mailers_section=$_SECTION_NAME"
  assert_equal "smtp2" "$(get_json_path "$BODY" .name)"
  assert_equal "10.0.10.88" "$(get_json_path "$BODY" .address)"
  assert_equal "8587" "$(get_json_path "$BODY" .port)"
}

@test "mailers: delete an entry" {
  resource_delete "$_MAILER_ENTRIES_PATH/smtp2" "mailers_section=$_SECTION_NAME" "force_reload=true"
  assert_equal "$SC" "202"
  resource_get "$_MAILER_ENTRIES_PATH/smtp2" "mailers_section=$_SECTION_NAME"
  assert_equal "$SC" "404"
}

@test "mailers: delete a section" {
  resource_delete "$_MAILERS_SECTION_PATH/$_SECTION_NAME" "force_reload=true"
  assert_equal "$SC" "204"
  resource_get "$_MAILERS_SECTION_PATH/$_SECTION_NAME"
  assert_equal "$SC" "404"
}
