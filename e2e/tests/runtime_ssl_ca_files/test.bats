#!/usr/bin/env bats
#
# Copyright 2025 HAProxy Technologies
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
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "ssl_ca_files_runtime: Create and delete an ssl CA file" {
    haproxy_version_ge "2.6" || skip

    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CA_FILES_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=cafile.pem"
    dpa_curl_status_body '$output'
    assert_equal $SC 201

    resource_get "$_RUNTIME_SSL_CA_FILES_BASE_PATH"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 2
    assert_equal "$(get_json_path "$BODY" '.[1].storage_name')" "cafile.pem"
    assert_equal "$(get_json_path "$BODY" '.[1].count')" "3 certificate(s)"

    resource_get "$_RUNTIME_SSL_CA_FILES_BASE_PATH/cafile.pem/entries/1"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 9
    assert_equal "$(get_json_path "$BODY" '.storage_name')" "cafile.pem"
    assert_equal "$(get_json_path "$BODY" '.algorithm')" "RSA2048"

    resource_delete "$_RUNTIME_SSL_CA_FILES_BASE_PATH/cafile.pem"
    assert_equal $SC 204

    resource_get "$_RUNTIME_SSL_CA_FILES_BASE_PATH"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 1
    # The name can change depending on systems.
    #assert_equal "$(get_json_path "$BODY" '.[0].storage_name')" "@system-ca"
}

@test "ssl_ca_files_runtime: Add a ssl certificate entry to a CA file" {
    haproxy_version_ge "2.8" || skip

    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CA_FILES_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=cafile.pem"
    dpa_curl_status_body '$output'
    assert_equal $SC 201

    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CA_FILES_BASE_PATH/cafile.pem/entries" "@${BATS_TEST_DIRNAME}/data/4.pem"
    dpa_curl_status_body '$output'
    assert_equal $SC 201

    resource_get "$_RUNTIME_SSL_CA_FILES_BASE_PATH/cafile.pem"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 2

    resource_get "$_RUNTIME_SSL_CA_FILES_BASE_PATH/cafile.pem/entries/2"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 9
    assert_equal "$(get_json_path "$BODY" '.storage_name')" "cafile.pem"
    assert_equal "$(get_json_path "$BODY" '.algorithm')" "RSA2048"
}

@test "ssl_ca_files_runtime: Replace the contents of a CA file" {
    haproxy_version_ge "2.6" || skip

    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CA_FILES_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=cafile.pem"
    dpa_curl_status_body '$output'
    assert_equal $SC 201

    run dpa_curl_file_upload PUT "$_RUNTIME_SSL_CA_FILES_BASE_PATH/cafile.pem" "@${BATS_TEST_DIRNAME}/data/4.pem;filename=cafile.pem"
    assert_success

    resource_get "$_RUNTIME_SSL_CA_FILES_BASE_PATH/cafile.pem/entries/1"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 9
    assert_equal "$(get_json_path "$BODY" '.storage_name')" "cafile.pem"
    assert_equal "$(get_json_path "$BODY" '.algorithm')" "RSA2048"
    assert_equal "$(get_json_path "$BODY" '.issuers')" "/CN=Interm2."
    assert_equal "$(get_json_path "$BODY" '.serial')" "1002"
}
