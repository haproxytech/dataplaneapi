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

    refute dpa_docker_exec 'ls /etc/haproxy/ssl/1.pem'
    pre_logs_count=$(docker logs dataplaneapi-e2e 2>&1 | wc -l)

    run dpa_curl_file_upload POST "/services/haproxy/storage/ssl_certificates" "@${BATS_TEST_DIRNAME}/1.pem;filename=1.pem"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201

    assert_equal $(get_json_path "$BODY" '.storage_name') "1.pem"

    assert dpa_docker_exec 'ls /etc/haproxy/ssl/1.pem'

    # confirm haproxy wasn't reloaded or restarted
    post_logs_count=$(docker logs dataplaneapi-e2e 2>&1 | wc -l)
    new_logs_count=$(( $pre_logs_count - $post_logs_count ))
    assert [ $new_logs_count = 0 ]
}

@test "Get a list of managed ssl certificate files" {
    run dpa_curl GET "/services/haproxy/storage/ssl_certificates/"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal $(get_json_path "$BODY" '.|length') 1

    assert_equal $(get_json_path "$BODY" '.[0].storage_name') "1.pem"
}

@test "Get a ssl certificate file contents" {

    run dpa_curl_download GET "/services/haproxy/storage/ssl_certificates/1.pem"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_output --partial "1.pem"

    # we opted to not return the certificate contents (i.e. secret keys, just the identifier)
    #assert dpa_diff_var_file '$BODY' "1.pem"
}

@test "Replace a ssl certificate file contents" {

    pre_logs_count=$(docker logs dataplaneapi-e2e 2>&1 | wc -l)

    run dpa_curl_text_plain PUT "/services/haproxy/storage/ssl_certificates/1.pem" "@${BATS_TEST_DIRNAME}/2.pem"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 202

    # confirm haproxy wasn't reloaded or restarted
    post_logs_count=$(docker logs dataplaneapi-e2e 2>&1 | wc -l)
    new_logs_count=$(( $pre_logs_count - $post_logs_count ))
    assert [ $new_logs_count = 0 ]
}

@test "Delete a ssl certificate file" {
    run dpa_curl DELETE "/services/haproxy/storage/ssl_certificates/1.pem"
    assert_success

    dpa_curl_status_body_safe '$output'
    assert_equal $SC 204

    refute dpa_docker_exec 'ls /etc/haproxy/ssl/1.pem'
}

@test "Add a ssl certificate file and reload HAProxy" {

    refute dpa_docker_exec 'ls /etc/haproxy/ssl/1.pem'

    pre_logs_count=$(docker logs dataplaneapi-e2e 2>&1 | wc -l)

    run dpa_curl_file_upload POST "/services/haproxy/storage/ssl_certificates?force_reload=true" "@${BATS_TEST_DIRNAME}/1.pem;filename=1.pem"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201

    assert_equal $(get_json_path "$BODY" '.storage_name') "1.pem"

    assert dpa_docker_exec 'ls /etc/haproxy/ssl/1.pem'

    post_logs_count=$(docker logs dataplaneapi-e2e 2>&1 | wc -l)
    new_logs_count=$(( $pre_logs_count - $post_logs_count ))

    new_logs=$(docker logs dataplaneapi-e2e 2>&1 | tail -n $new_logs_count)

    echo -e "$new_logs" # this will help debugging if the test fails
    assert echo -e "$new_logs" | head -n 1 | grep -q "Reexecuting Master process"
}

@test "Replace a ssl certificate file contents and reload HAPRoxy" {

    pre_logs_count=$(docker logs dataplaneapi-e2e 2>&1 | wc -l)

    run dpa_curl_text_plain PUT "/services/haproxy/storage/ssl_certificates/1.pem?force_reload=true" "@${BATS_TEST_DIRNAME}/2.pem"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 202

    post_logs_count=$(docker logs dataplaneapi-e2e 2>&1 | wc -l)
    new_logs_count=$(( $pre_logs_count - $post_logs_count ))

    new_logs=$(docker logs dataplaneapi-e2e 2>&1 | tail -n $new_logs_count)

    echo -e "$new_logs" # this will help debugging if the test fails
    assert echo -e "$new_logs" | head -n 1 | grep -q "Reexecuting Master process"

    # clean up after the test
    dpa_docker_exec 'rm /etc/haproxy/ssl/1.pem'
}
