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
load '../../libs/get_json_path'
load '../../libs/version'

@test "transactions: Outdated transactions cannot be committed" {
	# creating 5 transactions
	for _ in {1..5}; do
		run dpa_curl POST "/services/haproxy/transactions?version=$(version)";
		assert_success;
	done

	# retrieving the first one
	run dpa_curl GET "/services/haproxy/transactions?version=$(version)"
	dpa_curl_status_body '$output'
	local id
	id=$(get_json_path "${BODY}" ".[0].id")

	# commit it, must succeed
	run dpa_curl PUT "/services/haproxy/transactions/${id}?version=$(version)"
	assert_success
	dpa_curl_status_body '$output'
 	assert_equal "$SC" 202

	# retrieve other transactions
 	run dpa_curl GET "/services/haproxy/transactions?version=$(version)"
	dpa_curl_status_body '$output'
	# iterate over them, should fail with 406 status code
	for tx in $(echo "${BODY}" | jq -r '.[].id'); do
		run dpa_curl PUT "/services/haproxy/transactions/${tx}?version=$(version)"
		assert_success
		dpa_curl_status_body '$output'
		assert_equal "$SC" 406
	done
}
