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
DONT_RESTART_DPAPI=1
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'
load '../../libs/acme'

_RUNTIME_ACME_PATH="/services/haproxy/runtime/acme"
_CERT_NAME="/var/lib/haproxy/haproxy.pem"

@test "acme_runtime: Renew a certificate" {
    haproxy_version_ge "3.3" || skip

    sleep 2 # wait for haproxy to create the acme account key
    run dpa_curl PUT "$_RUNTIME_ACME_PATH?certificate=$_CERT_NAME"
	assert_success
	dpa_curl_status_body '$output'
    assert_equal "$SC" 200

    # Wait until the status of our certificate is in state "Scheduled",
    # meaning it was renewed successfully.
    state=unknown
    trials=10
    while [ "$state" != "Scheduled" ] && (( trials > 0 )); do
        sleep 2
        resource_get "$_RUNTIME_ACME_PATH"
        state="$(jq -r .[0].state <<< "$BODY")"
        : $((trials--))
    done
    assert_equal "$state" "Scheduled"

    # HAProxy will then send an event to dpapi to store the cert to disk.
    sleep 1
    run dpa_docker_exec 'ls /etc/haproxy/ssl/haproxy.pem'
}
