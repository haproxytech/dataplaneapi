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

@test "runtime_maps_entries: Replace the value corresponding to each id in a map" {
    resource_put "$_RUNTIME_MAP_ENTRIES_BASE_PATH/key1" "data/put.json" "map=mapfile1.map&force_sync=true"
    assert_equal "$SC" 200
    #
    # verify that entry is actually changed
    #
    resource_get "$_RUNTIME_MAP_ENTRIES_BASE_PATH/key1" "map=mapfile1.map"
    assert_equal "$SC" 200

    assert_equal "$(get_json_path "${BODY}" " select(.value | contains(\"replacedvalue\") ).value" )" "replacedvalue"
}

@test "runtime_maps_entries: Return an error when trying to replace non existing map key" {
    resource_put "$_RUNTIME_MAP_ENTRIES_BASE_PATH/not-exists" "data/put.json" "map=mapfile1.map&force_sync=true"
    assert_equal "$SC" 404
}
