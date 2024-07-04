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
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "http_checks: Add a new HTTP Check to defaults" {
    PARENT_NAME="mydefaults"
    resource_post "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks/2" "data/post_defaults.json" "force_reload=true"
	assert_equal "$SC" 201

    resource_get "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_checks/2"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ".type")" "disable-on-404"
}

@test "http_checks: Add a new HTTP Check to backend" {
    PARENT_NAME="test_backend"
    resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/2" "data/post.json" "force_reload=true"
	assert_equal "$SC" 201

    resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/2"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" ".type")" "disable-on-404"
}

@test "http_checks: fail adding an unvalid HTTP Check to backend" {
    PARENT_NAME="test_backend"
    resource_post "$_BACKEND_BASE_PATH/$PARENT_NAME/http_checks/0" "data/post_unvalid.json" "force_reload=true"
	assert_equal "$SC" 400
}
