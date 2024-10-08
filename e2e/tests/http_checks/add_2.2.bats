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
load '../../libs/get_json_path'
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "http_checks: Add a new HTTP Check (send) to defaults (>=2.2)" {
    if haproxy_version_ge "2.2"
    then
    PARENT_NAME="mydefaults"
    resource_post "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks/2" "data/post_defaults_send.json" "force_reload=true"
	assert_equal "$SC" 201

    resource_get "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks/2"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ".method")" "OPTIONS"
	assert_equal "$(get_json_path "$BODY" ".uri")" "/"
	assert_equal "$(get_json_path "$BODY" ".version")" "HTTP/1.1"
    fi
}


@test "http_checks: Add a new HTTP Check (send) to backend (>=2.2)" {
    if haproxy_version_ge "2.2"
    then
    PARENT_NAME="test_backend_2"
    resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/0" "data/post_send.json" "force_reload=true"
	assert_equal "$SC" 201

    resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/0"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ".method")" "OPTIONS"
	assert_equal "$(get_json_path "$BODY" ".uri")" "/"
	assert_equal "$(get_json_path "$BODY" ".version")" "HTTP/1.1"
    fi
}


@test "http_checks: fail adding an invalid HTTP Check (send method) to backend (>=2.2)" {
    if haproxy_version_ge "2.2"
    then
    PARENT_NAME="test_backend_2"
    resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/0" "data/post_invalid_send_method.json" "force_reload=true"
	assert_equal "$SC" 422
    assert_equal "$(get_json_path "$BODY" ".code")" "606"
    fi
}
