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

@test "ssl_certs_runtime: Create a new ssl certificate entry" {
    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CERTS_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=cert.pem"
    assert_success

    resource_get "$_RUNTIME_SSL_CERTS_BASE_PATH/cert.pem"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.storage_name')" "cert.pem"
    assert_equal "$(get_json_path "$BODY" '.algorithm')" "RSA2048"
    assert_equal "$(get_json_path "$BODY" '.issuers')" "/CN=Interm2."
    assert_equal "$(get_json_path "$BODY" '.serial')" "1002"
    assert_equal "$(get_json_path "$BODY" '.chain_issuer')" "/CN=Interm1."
    assert_equal "$(get_json_path "$BODY" '.chain_subject')" "/CN=Interm2."
    assert_equal "$(get_json_path "$BODY" '.not_after')" "2021-11-25T12:12:04.000Z"
    assert_equal "$(get_json_path "$BODY" '.not_before')" "2020-11-25T12:12:04.000Z"
    assert_equal "$(get_json_path "$BODY" '.sha1_finger_print')" "BC215F40136157964B69114E3D934934749798D5"
    assert_equal "$(get_json_path "$BODY" '.status')" "Unused"
    assert_equal "$(get_json_path "$BODY" '.subject')" "/CN=1.example.com"
}

@test "ssl_certs_runtime: Set a payload to the ssl certificate entry" {
    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CERTS_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=cert.pem"
    assert_success

    run dpa_curl_file_upload PUT "$_RUNTIME_SSL_CERTS_BASE_PATH/cert.pem" "@${BATS_TEST_DIRNAME}/data/2.pem"
    assert_success

    resource_get "$_RUNTIME_SSL_CERTS_BASE_PATH/cert.pem"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 11
    assert_equal "$(get_json_path "$BODY" '.storage_name')" "cert.pem"
    assert_equal "$(get_json_path "$BODY" '.algorithm')" "RSA2048"
    assert_equal "$(get_json_path "$BODY" '.issuers')" "/CN=Interm2."
    assert_equal "$(get_json_path "$BODY" '.serial')" "1002"
    assert_equal "$(get_json_path "$BODY" '.chain_issuer')" "/CN=Interm1."
    assert_equal "$(get_json_path "$BODY" '.chain_subject')" "/CN=Interm2."
    assert_equal "$(get_json_path "$BODY" '.not_after')" "2021-11-25T12:12:04.000Z"
    assert_equal "$(get_json_path "$BODY" '.not_before')" "2020-11-25T12:12:04.000Z"
    assert_equal "$(get_json_path "$BODY" '.sha1_finger_print')" "BC215F40136157964B69114E3D934934749798D5"
    assert_equal "$(get_json_path "$BODY" '.status')" "Unused"
    assert_equal "$(get_json_path "$BODY" '.subject')" "/CN=1.example.com"
}

@test "ssl_certs_runtime: Get all the certs" {
    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CERTS_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=cert1.pem"
    assert_success

    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CERTS_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/3.pem;filename=cert2.pem"
    assert_success

    resource_get "$_RUNTIME_SSL_CERTS_BASE_PATH"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 2
    assert_equal "$(get_json_path "$BODY" '.[0].storage_name')" "cert1.pem"
    assert_equal "$(get_json_path "$BODY" '.[1].storage_name')" "cert2.pem"
}

@test "ssl_certs_runtime: Delete a ssl certificate entry" {
    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CERTS_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=cert.pem"
    assert_success

    resource_delete "$_RUNTIME_SSL_CERTS_BASE_PATH/cert.pem"
    assert_equal $SC 204

    resource_get "$_RUNTIME_SSL_CERTS_BASE_PATH/cert.pem"
    assert_equal "$SC" 404
}
