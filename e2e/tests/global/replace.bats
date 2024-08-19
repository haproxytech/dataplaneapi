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
load '../../libs/get_json_path'
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "global: Replace a global configuration" {
    resource_put "$_GLOBAL_BASE_PATH" "data/put.json" ""
	assert_equal "$SC" 202

	resource_get "$_GLOBAL_BASE_PATH" ""
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" '.performance_options.maxconn')" "5000"
	assert_equal "$(get_json_path "$BODY" '.daemon')" "true"
	assert_equal "$(get_json_path "$BODY" '.pidfile')" "/var/run/haproxy.pid"
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].address')" "/var/lib/haproxy/stats"
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].level')" "admin"
}


@test "global: Replace a global configuration with socket path changed" {
    resource_put "$_GLOBAL_BASE_PATH" "data/put_socket.json" ""
	assert_equal "$SC" 202

	resource_get "$_GLOBAL_BASE_PATH" ""
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" '.performance_options.maxconn')" "5000"
	assert_equal "$(get_json_path "$BODY" '.daemon')" "true"
	assert_equal "$(get_json_path "$BODY" '.pidfile')" "/var/run/haproxy.pid"
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].address')" "/var/lib/haproxy/stats-new"
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].level')" "admin"

	# check that runtime client has been reconfigured with the new socket
	sleep 5
	resource_get "$_RUNTIME_MAP_FILES_BASE_PATH" ""
    assert_equal "$SC" 200
}

@test "global: Replace a global configuration with socket path changed (using transaction)" {
    # create transaction
    resource_post "$_TRANSACTIONS_BASE_PATH" ""
    assert_equal "$SC" 201
    local transaction_id; transaction_id=$(get_json_path "${BODY}" ".id")

    # PUT new configuration
    run dpa_curl PUT "${_GLOBAL_BASE_PATH}?transaction_id=${transaction_id}" "data/put_socket.json"
	assert_success
	dpa_curl_status_body '$output'
	assert_equal "$SC" 202

	# commit transaction
    resource_put "$_TRANSACTIONS_BASE_PATH/$transaction_id" ""
    assert_equal "$SC" 202

    # check configuration has been applied
	resource_get "$_GLOBAL_BASE_PATH" ""
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" '.performance_options.maxconn')" "5000"
	assert_equal "$(get_json_path "$BODY" '.daemon')" "true"
	assert_equal "$(get_json_path "$BODY" '.pidfile')" "/var/run/haproxy.pid"
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].address')" "/var/lib/haproxy/stats-new"
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].level')" "admin"

	# check that runtime client has been reconfigured with the new socket
	sleep 5
	resource_get "$_RUNTIME_MAP_FILES_BASE_PATH" ""
    assert_equal "$SC" 200
}


@test "global: Manually replace a global configuration with socket path changed" {
    # check HAPRoxy is configured with the expected socket
	resource_get "$_GLOBAL_BASE_PATH" ""
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].address')" "/var/lib/haproxy/stats"

    pre_logs_count=$(dpa_docker_exec 'cat /var/log/dataplaneapi.log' | wc -l)

	# manually change configuration
	run dpa_docker_exec "sed -i 's@/var/lib/haproxy/stats@/var/lib/haproxy/stats-new@' /etc/haproxy/haproxy.cfg"

    sleep 5
    # check configuration has been reloaded
    resource_get "$_GLOBAL_BASE_PATH" ""
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].address')" "/var/lib/haproxy/stats-new"

    # check that runtime client has been reconfigured with the new socket
    post_logs_count=$(dpa_docker_exec 'sh /var/log/dataplaneapi.log' | wc -l)
    new_logs_count=$(( $pre_logs_count - $post_logs_count ))
    new_logs=$(dpa_docker_exec 'cat /var/log/dataplaneapi.log' | tail -n $new_logs_count)

    echo "$new_logs" # this will help debugging if the test fails
    assert echo -e "$new_logs" | grep -q "reload callback completed, runtime API reconfigured"

	resource_get "$_RUNTIME_MAP_FILES_BASE_PATH" ""
    assert_equal "$SC" 200
}
