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
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "runtime_maps_entries: Return one map runtime entries" {
    PARENT_NAME="mapfile1.map"
    resource_get "$_RUNTIME_MAP_BASE_PATH/$PARENT_NAME/entries"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "$BODY" " .[] | select(.key | contains(\"key1\") ).key" )" "key1"
    assert_equal "$(get_json_path "$BODY" " .[] | select(.value | contains(\"value1\") ).value" )" "value1"

    assert_equal "$(get_json_path "$BODY" " .[] | select(.key | contains(\"api.example.com\") ).key" )" "api.example.com"
    assert_equal "$(get_json_path "$BODY" " .[] | select(.value | contains(\"be_api\") ).value" )" "be_api"
}

@test "runtime_maps_entries: https://github.com/haproxytech/dataplaneapi/issues/159" {
    PARENT_NAME="not-exists.map"
    resource_get "$_RUNTIME_MAP_BASE_PATH/$PARENT_NAME/entries"
    assert_equal "$SC" 404
}

@test "runtime_maps_entries: https://github.com/haproxytech/dataplaneapi/issues/234" {
    PARENT_NAME="empty.map"
    resource_get "$_RUNTIME_MAP_BASE_PATH"/"$PARENT_NAME/entries"
    assert_equal "$SC" 200
    assert_equal "$($BODY)" ""
}
