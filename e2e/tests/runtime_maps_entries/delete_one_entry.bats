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

@test "runtime_maps_entries: Deletes all the map entries from the map by its id" {
    resource_delete "$_RUNTIME_MAP_ENTRIES_BASE_PATH/key1" "map=mapfile1.map"
    assert_equal "$SC" 204
    #
    # verify that entry is actually deleted
    #
    resource_get "$_RUNTIME_MAP_ENTRIES_BASE_PATH/key1" "map=mapfile1.map"
    assert_equal "$SC" 404
}

@test "runtime_maps_entries: https://github.com/haproxytech/dataplaneapi/issues/160 with non existing key" {
    resource_delete "$_RUNTIME_MAP_ENTRIES_BASE_PATH/not-exists" "map=mapfile1.map"
    assert_equal "$SC" 404
}

@test "runtime_maps_entries: https://github.com/haproxytech/dataplaneapi/issues/160 with non existing map file" {
    resource_delete "$_RUNTIME_MAP_ENTRIES_BASE_PATH/key1" "map=not-exists.map"
    assert_equal "$SC" 404
}
