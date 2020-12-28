#!/usr/bin/env bats
#
# Copyright 2019 HAProxy Technologies
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
load '../../libs/version'

@test "Add a mapfile" {

    [ -z "$(docker exec dataplaneapi-e2e /bin/sh -c 'ls /etc/haproxy/maps/ | grep mapfile_example.map')" ]

    read -r SC BODY < <(dataplaneapi_file_upload POST "/services/haproxy/storage/maps" "@${BATS_TEST_DIRNAME}/mapfile_example.map;filename=mapfile_example.map")

    [ "${SC}" = 201 ]

    local STORAGE_NAME; STORAGE_NAME=$(get_json_path "$BODY" '.storage_name')
    [ "${STORAGE_NAME}" = "mapfile_example.map" ]

    [ ! -z "$(docker exec dataplaneapi-e2e /bin/sh -c 'ls /etc/haproxy/maps/ | grep mapfile_example.map')" ]
}

@test "Get a list of managed mapfiles" {

    read -r SC BODY < <(dataplaneapi GET "/services/haproxy/storage/maps/")

    [ "${SC}" = 200 ]

    local LENGTH; LENGTH=$(get_json_path "$BODY" '.|length')
    [ "${LENGTH}" == 1 ]

    local STORAGE_NAME; STORAGE_NAME=$(get_json_path "$BODY" '.[0].storage_name')
    [ "${STORAGE_NAME}" = "mapfile_example.map" ]
}

@test "Get a mapfile contents" {
    local BODY;
    local SC;

    dataplaneapi_download GET "/services/haproxy/storage/maps/mapfile_example.map"

    [ "${SC}" = 200 ]

    [ -z "$(diff <(echo -e "$BODY") ${E2E_DIR}/fixtures/mapfile_example.map)" ]
}

@test "Try to get unavailable mapfile contents" {

    read -r SC BODY < <(dataplaneapi GET "/services/haproxy/storage/maps/not_here.map")

    [ "${SC}" = 404 ]
}

@test "Replace a mapfile contents" {

    read -r SC BODY < <(dataplaneapi_text_plain PUT "/services/haproxy/storage/maps/mapfile_example.map" "@${BATS_TEST_DIRNAME}/mapfile_example2.map")

    [ "${SC}" = 202 ]

    dataplaneapi_download GET "/services/haproxy/storage/maps/mapfile_example.map"
    [ -z "$(diff <(echo -e "$BODY") ${BATS_TEST_DIRNAME}/mapfile_example2.map)" ]
}

@test "Delete a mapfile" {

    read -r SC BODY < <(dataplaneapi DELETE "/services/haproxy/storage/maps/mapfile_example.map")

    [ "${SC}" = 204 ]

    [ -z "$(docker exec dataplaneapi-e2e /bin/sh -c 'ls /etc/haproxy/maps/ | grep mapfile_example.map')" ]
}
