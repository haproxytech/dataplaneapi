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

@test "storage_maps: Add a mapfile" {

    refute dpa_docker_exec 'ls /etc/haproxy/maps/mapfile_example.map'

    run dpa_curl_file_upload POST "/services/haproxy/storage/maps" "@${BATS_TEST_DIRNAME}/mapfile_example.map;filename=mapfile_example.map"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201

    assert_equal $(get_json_path "$BODY" '.storage_name') "mapfile_example.map"

    assert dpa_docker_exec 'ls /etc/haproxy/maps/mapfile_example.map'
}

@test "storage_maps: Get a list of managed mapfiles" {

    # sometimes we can't establish a connection to the haproxy stat socket
    # forcing haproxy to restart seems to fix that
    assert dpa_docker_exec 'kill -SIGUSR2 1'
    sleep 1

    run dpa_curl GET "/services/haproxy/storage/maps/"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal $(get_json_path "$BODY" '.|length') 1

    assert_equal $(get_json_path "$BODY" '.[0].storage_name') "mapfile_example.map"
}

@test "storage_maps: Get a mapfile contents" {

    run dpa_curl_download GET "/services/haproxy/storage/maps/mapfile_example.map"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert dpa_diff_var_file '$BODY' 'mapfile_example.map'

    assert dpa_diff_docker_file '/etc/haproxy/maps/mapfile_example.map' "mapfile_example.map"
}

@test "storage_maps: Try to get unavailable mapfile contents" {

    run dpa_curl GET "/services/haproxy/storage/maps/not_here.map"
    assert_success

    dpa_curl_status_body_safe '$output'
    assert_equal $SC 404
}

@test "storage_maps: Replace a mapfile contents" {

    run dpa_curl_text_plain PUT "/services/haproxy/storage/maps/mapfile_example.map" "@${BATS_TEST_DIRNAME}/mapfile_example2.map"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 202

    run dpa_curl_download GET "/services/haproxy/storage/maps/mapfile_example.map"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert dpa_diff_var_file '$BODY' 'mapfile_example2.map'

    assert dpa_diff_docker_file '/etc/haproxy/maps/mapfile_example.map' "mapfile_example2.map"
}

@test "storage_maps: Delete a mapfile" {

    run dpa_curl DELETE "/services/haproxy/storage/maps/mapfile_example.map"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 204

    refute dpa_docker_exec 'ls /etc/haproxy/maps/mapfile_example.map'
}
