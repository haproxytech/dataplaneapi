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

@test "ssl_crl_files_runtime: Create a new ssl crl file" {
    haproxy_version_ge "2.6" || skip

    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CRL_FILES_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=crlfile.pem"
    assert_success

    resource_get "$_RUNTIME_SSL_CRL_FILES_BASE_PATH"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 1
    assert_equal "$(get_json_path "$BODY" '.[0].storage_name')" "crlfile.pem"

    resource_get "$_RUNTIME_SSL_CRL_FILES_BASE_PATH/crlfile.pem"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 1
    assert_equal "$(get_json_path "$BODY" '.[0]|length')" 8
    assert_equal "$(get_json_path "$BODY" '.[0].storage_name')" "crlfile.pem"
    assert_equal "$(get_json_path "$BODY" '.[0].issuer')" "/C=CN/O=HEP/CN=Institute of High Energy Physics Certification Authority"
    assert_equal "$(get_json_path "$BODY" '.[0].next_update')" "2023-06-08"
    assert_equal "$(get_json_path "$BODY" '.[0].revoked_certificates[0].serial_number')" "04"
    assert_equal "$(get_json_path "$BODY" '.[0].revoked_certificates[0].revocation_date')" "2013-05-08"
    assert_equal "$(get_json_path "$BODY" '.[0].revoked_certificates[1].serial_number')" "05"
    assert_equal "$(get_json_path "$BODY" '.[0].revoked_certificates[1].revocation_date')" "2013-05-08"
    assert_equal "$(get_json_path "$BODY" '.[0].revoked_certificates[2].serial_number')" "08"
    assert_equal "$(get_json_path "$BODY" '.[0].revoked_certificates[2].revocation_date')" "2013-05-10"
    assert_equal "$(get_json_path "$BODY" '.[0].signature_algorithm')" "sha1WithRSAEncryption"
    assert_equal "$(get_json_path "$BODY" '.[0].status')" "Unused"
    assert_equal "$(get_json_path "$BODY" '.[0].version')" "2"
}

@test "ssl_crl_files_runtime: Update an ssl crl file" {
    haproxy_version_ge "2.6" || skip

    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CRL_FILES_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=crlfile.pem"
    assert_success

    run dpa_curl_file_upload PUT "$_RUNTIME_SSL_CRL_FILES_BASE_PATH/crlfile.pem" "@${BATS_TEST_DIRNAME}/data/2.pem;filename=2.pem"
    assert_success

    resource_get "$_RUNTIME_SSL_CRL_FILES_BASE_PATH/crlfile.pem" "index=1"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.[0]|length')" 7
    assert_equal "$(get_json_path "$BODY" '.[0].storage_name')" "crlfile.pem"
    assert_equal "$(get_json_path "$BODY" '.[0].issuer')" "/C=US/O=DigiCert Inc/CN=DigiCert SHA2 Secure Server CA"
    assert_equal "$(get_json_path "$BODY" '.[0].next_update')" "2023-05-24"
    assert_equal "$(get_json_path "$BODY" '.[0].signature_algorithm')" "sha256WithRSAEncryption"
    assert_equal "$(get_json_path "$BODY" '.[0].status')" "Unused"
}

@test "ssl_crl_files_runtime: Delete an ssl crl file" {
    haproxy_version_ge "2.6" || skip

    run dpa_curl_file_upload POST "$_RUNTIME_SSL_CRL_FILES_BASE_PATH" "@${BATS_TEST_DIRNAME}/data/1.pem;filename=crlfile.pem"
    assert_success

    resource_delete "$_RUNTIME_SSL_CRL_FILES_BASE_PATH/crlfile.pem"
    assert_equal $SC 204

    resource_get "$_RUNTIME_SSL_CRL_FILES_BASE_PATH"
    assert_equal "$SC" 200
    assert_equal "$(get_json_path "$BODY" '.|length')" 0
}
