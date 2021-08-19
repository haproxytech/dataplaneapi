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

@test "storage_ssl_certificates: Refuse to delete still used ssl certificate file" {
    run docker cp "${BATS_TEST_DIRNAME}/data/3.pem" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/ssl/"
    assert_success

    resource_delete "$_STORAGE_SSL_CERTS_BASE_PATH/3.pem"
    assert_equal $SC 409

    assert dpa_docker_exec 'ls /etc/haproxy/ssl/3.pem'

    # clean up this test
    assert dpa_docker_exec 'rm /etc/haproxy/ssl/3.pem'
}

@test "storage_ssl_certificates: Allow to delete ssl certificate file referenced in comments" {
    run docker cp "${BATS_TEST_DIRNAME}/data/4.pem" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/ssl/"
    assert_success

    resource_delete "$_STORAGE_SSL_CERTS_BASE_PATH/4.pem"
    assert_equal "$SC" 202

    refute dpa_docker_exec 'ls /etc/haproxy/ssl/4.pem'

    # clean up this test
    run dpa_docker_exec 'rm /etc/haproxy/ssl/4.pem'
}
