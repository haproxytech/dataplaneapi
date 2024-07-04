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

@test "runtime_maps_entries: Adds an entry into the map file" {
    PARENT_NAME="mapfile1.map"
    resource_post "$_RUNTIME_MAP_BASE_PATH/$PARENT_NAME/entries" "data/post.json"
    assert_equal "$SC" 201
    #
    # verify that entry is actually added
    #
    resource_get "$_RUNTIME_MAP_BASE_PATH/$PARENT_NAME/entries/newkey"
    assert_equal "$SC" 200
}

@test "runtime_maps_entries: Refuse adding an existing map entry into the map file" {
    PARENT_NAME="mapfile1.map"
    resource_post "$_RUNTIME_MAP_BASE_PATH/$PARENT_NAME/entries" "data/existing.json"
    assert_equal "$SC" 409
}
