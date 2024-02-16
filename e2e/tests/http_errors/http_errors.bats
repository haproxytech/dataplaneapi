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

load '../../libs/dataplaneapi'
load '../../libs/get_json_path'
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "http_error_sections: Return all sections" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_get "$_ERR_SECTIONS_BASE_PATH"
  assert_equal "$SC" 200

	assert_equal "$(get_json_path "$BODY" ". | length")" 2
}

@test "http_error_sections: Return one section by name" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_get "$_ERR_SECTIONS_BASE_PATH/website-2"
  assert_equal "$SC" 200

  assert_equal "$(get_json_path "$BODY" ".name")" "website-2"
  assert_equal "$(get_json_path "$BODY" ".error_files | length")" 3
  assert_equal "$(get_json_path "$BODY" ".error_files[0].code")" 500
  assert_equal "$(get_json_path "$BODY" ".error_files[0].file")" "/dev/null"
  assert_equal "$(get_json_path "$BODY" ".error_files[1].code")" 404
  assert_equal "$(get_json_path "$BODY" ".error_files[1].file")" "/dev/null"
  assert_equal "$(get_json_path "$BODY" ".error_files[2].code")" 503
  assert_equal "$(get_json_path "$BODY" ".error_files[2].file")" "/dev/null"
}

@test "http_error_sections: Fail to return a section that does not exist" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_get "$_ERR_SECTIONS_BASE_PATH/i_am_not_here"
  assert_equal "$SC" 404
}

@test "http_error_sections: Delete one section by name" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_delete "$_ERR_SECTIONS_BASE_PATH/website-2"
  assert_equal "$SC" 202

  resource_get "$_ERR_SECTIONS_BASE_PATH/website-2"
  assert_equal "$SC" "404"
}

@test "http_error_sections: Fail to delete a section that does not exist" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_delete "$_ERR_SECTIONS_BASE_PATH/i_am_not_here"
  assert_equal "$SC" 404
}

@test "http_error_sections: Create a new section" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_post "$_ERR_SECTIONS_BASE_PATH" "data/section.json" "force_reload=true"
  assert_equal "$SC" 201

  resource_get "$_ERR_SECTIONS_BASE_PATH/website-3"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".error_files | length")" 2
  assert_equal "$(get_json_path "$BODY" ".error_files[0].code")" 500
  assert_equal "$(get_json_path "$BODY" ".error_files[0].file")" "/dev/null"
  assert_equal "$(get_json_path "$BODY" ".error_files[1].code")" 502
  assert_equal "$(get_json_path "$BODY" ".error_files[1].file")" "/dev/null"
}

@test "http_error_sections: Fail to create a section with an existing name" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_post "$_ERR_SECTIONS_BASE_PATH" "data/section_name_exists.json" "force_reload=true"
  assert_equal "$SC" 409
}

@test "http_error_sections: Fail to create a section with missing data" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_post "$_ERR_SECTIONS_BASE_PATH" "data/section_name_missing.json" "force_reload=true"
  assert_equal "$SC" 422

  resource_post "$_ERR_SECTIONS_BASE_PATH" "data/section_entries_missing.json" "force_reload=true"
  assert_equal "$SC" 422
}

@test "http_error_sections: Fail to create a section with invalid data" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_post "$_ERR_SECTIONS_BASE_PATH" "data/section_errorfile_no_file.json" "force_reload=true"
  assert_equal "$SC" 400

  resource_post "$_ERR_SECTIONS_BASE_PATH" "data/section_errorfile_no_code.json" "force_reload=true"
  assert_equal "$SC" 400

  resource_post "$_ERR_SECTIONS_BASE_PATH" "data/section_errorfile_unsupported_code.json" "force_reload=true"
  assert_equal "$SC" 422
}

@test "http_error_sections: Replace one section" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_put "$_ERR_SECTIONS_BASE_PATH/website-1" "data/section_replace.json" "force_reload=true"
  assert_equal "$SC" 200

  resource_get "$_ERR_SECTIONS_BASE_PATH/website-1"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".name")" "website-1"
  assert_equal "$(get_json_path "$BODY" ".error_files | length")" 2
  assert_equal "$(get_json_path "$BODY" ".error_files[0].code")" 500
  assert_equal "$(get_json_path "$BODY" ".error_files[0].file")" "/dev/null"
  assert_equal "$(get_json_path "$BODY" ".error_files[1].code")" 502
  assert_equal "$(get_json_path "$BODY" ".error_files[1].file")" "/dev/null"
}

@test "http_error_sections: Fail to replace a section that does not exist" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_put "$_ERR_SECTIONS_BASE_PATH/i_am_not_there" "data/section_replace.json" "force_reload=true"
  assert_equal "$SC" 409
}

@test "http_error_sections: Fail to replace a section with missing data" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_put "$_ERR_SECTIONS_BASE_PATH/website-1" "data/section_name_missing.json" "force_reload=true"
  assert_equal "$SC" 422

  resource_put "$_ERR_SECTIONS_BASE_PATH/website-1" "data/section_entries_missing.json" "force_reload=true"
  assert_equal "$SC" 422
}

@test "http_error_sections: Fail to replace a section with invalid data" {
  haproxy_version_ge $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SECTIONS_SUPPORTED_HAPROXY_VERSION+"

  resource_put "$_ERR_SECTIONS_BASE_PATH/website-1" "data/section_replace_errorfile_no_file.json" "force_reload=true"
  assert_equal "$SC" 400

  resource_put "$_ERR_SECTIONS_BASE_PATH/website-1" "data/section_replace_errorfile_no_code.json" "force_reload=true"
  assert_equal "$SC" 400

  resource_put "$_ERR_SECTIONS_BASE_PATH/website-1" "data/section_replace_errorfile_unsupported_code.json" "force_reload=true"
  assert_equal "$SC" 422
}
