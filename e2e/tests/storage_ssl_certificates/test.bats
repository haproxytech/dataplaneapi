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

@test "Add a ssl certificate file" {
	read -r SC BODY < <(dpa_curl_file_upload POST "/services/haproxy/storage/ssl_certificates" "@${BATS_TEST_DIRNAME}/1.pem;filename=1.pem")

	[ "${SC}" = 201 ]

    local STORAGE_NAME; STORAGE_NAME=$(get_json_path "$BODY" '.storage_name')
    [ "${STORAGE_NAME}" = "1.pem" ]

    [ ! -z "$(docker exec dataplaneapi-e2e /bin/sh -c 'ls /etc/haproxy/ssl/ | grep 1.pem')" ]
}

@test "Get a list of managed ssl certificate files" {
	read -r SC BODY < <(dpa_curl GET "/services/haproxy/storage/ssl_certificates/")

	[ "${SC}" = 200 ]

    local LENGTH; LENGTH=$(get_json_path "$BODY" '.|length')
    [ "${LENGTH}" == 1 ]

    local STORAGE_NAME; STORAGE_NAME=$(get_json_path "$BODY" '.[0].storage_name')
    [ "${STORAGE_NAME}" = "1.pem" ]
}

@test "Get a ssl certificate file contents" {
    local BODY;
    local SC;

	dpa_curl_download GET "/services/haproxy/storage/ssl_certificates/1.pem"

	[ "${SC}" = 200 ]

    [ -z "$(diff <(echo -e "$BODY") ${E2E_DIR}/fixtures/1.pem)" ]

}

@test "Replace a ssl certificate file contents" {
	read -r SC BODY < <(dpa_curl_text_plain PUT "/services/haproxy/storage/ssl_certificates/1.pem" "@${BATS_TEST_DIRNAME}/2.pem")

    [ "${SC}" = 202 ]
}

@test "Delete a ssl certificate file" {
	read -r SC BODY < <(dpa_curl DELETE "/services/haproxy/storage/ssl_certificates/1.pem")

	[ "${SC}" = 204 ]

    [ -z "$(docker exec dataplaneapi-e2e /bin/sh -c 'ls /etc/haproxy/ssl/ | grep 1.pem')" ]
}
