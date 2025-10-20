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
load '../../libs/haproxy_version'
if haproxy_version_ge "3.3"; then
    load '../../libs/haproxy_config_setup'
    load '../../libs/acme'
fi

_RUNTIME_ACME_PATH="/services/haproxy/runtime/acme"

@test "acme_runtime: Renew a certificate" {
    haproxy_version_ge "3.3" || skip

    cert_name="/var/lib/haproxy/haproxy.pem"

    # Send an 'acme renew' message to HAProxy.
    run dpa_curl PUT "$_RUNTIME_ACME_PATH?certificate=$cert_name"
	assert_success
	dpa_curl_status_body '$output'
    assert_equal "$SC" 200

    # Wait until the status of our certificate is in state "Scheduled",
    # meaning it was renewed successfully.
    state=unknown trials=10
    while [ "$state" != "Scheduled" ] && (( trials > 0 )); do
        sleep 2
        resource_get "$_RUNTIME_ACME_PATH"
        state="$(jq -r .[0].state <<< "$BODY")"
        : $((trials--))
    done
    assert_equal "$state" "Scheduled"

    # HAProxy will then send an event to dpapi to store the cert to disk.
    timeout=8 elapsed=0 inc=1 found=false
    while ((elapsed < timeout)); do
        sleep $inc && elapsed=$((elapsed + inc))
        if dpa_docker_exec 'ls -l /etc/haproxy/ssl/haproxy.pem'; then
            found=true
            break
        fi
    done
    assert_equal "$found" true
}

@test "acme_runtime: dns-01 challenge" {
    haproxy_version_ge "3.3" || skip

    cert_name="/var/lib/haproxy/haproxy2.pem"

    run dpa_curl PUT "$_RUNTIME_ACME_PATH?certificate=$cert_name"
	assert_success
	dpa_curl_status_body '$output'
    assert_equal "$SC" 200

    timeout=20 elapsed=0 inc=2 found=false
    while ((elapsed < timeout)); do
        sleep $inc && elapsed=$((elapsed + inc))
        if dpa_docker_exec 'ls -l /etc/haproxy/ssl/haproxy2.pem'; then
            found=true
            break
        fi
    done
    assert_equal "$found" true
}
