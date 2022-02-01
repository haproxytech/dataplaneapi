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
load '../../libs/resource_client'
load '../../libs/haproxy_config_setup'
load '../../libs/version'

load 'utils/_helpers'

@test "transactions: Outdated transactions cannot be committed" {
	# creating 5 transactions
	for _ in {1..5}; do
	  resource_post "$_TRANSACTIONS_BASE_PATH" ""
	  assert_equal "$SC" 201
	done

	# retrieving the first one
	resource_get "$_TRANSACTIONS_BASE_PATH"
	local id; id=$(get_json_path "${BODY}" ".[0].id")

	# commit it, must succeed
	resource_put "$_TRANSACTIONS_BASE_PATH/$id" ""
 	assert_equal "$SC" 202

	# retrieve other transactions
	resource_get "$_TRANSACTIONS_BASE_PATH"
  	assert_equal "$SC" 200
	# iterate over them, should fail with 406 status code
	for tx in $(echo "${BODY}" | jq -r '.[].id'); do
		resource_put "$_TRANSACTIONS_BASE_PATH/${tx}" ""
		assert_equal "$SC" 406
	done
}
